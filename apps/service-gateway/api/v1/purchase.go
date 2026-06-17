package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"openerp/apps/service-gateway/db"
	"openerp/apps/service-gateway/models"
	"openerp/apps/service-gateway/workflow"
	"openerp/packages/pkg-logger"
)

type PurchaseOrderLineRequest struct {
	ItemCode  string  `json:"itemCode" binding:"required"`
	Qty       float64 `json:"qty" binding:"required,gt=0"`
	UnitPrice float64 `json:"unitPrice" binding:"gte=0"`
	BaseType  int     `json:"baseType"`
	BaseEntry uint    `json:"baseEntry"`
	BaseLine  uint    `json:"baseLine"`
}

type PurchaseOrderRequest struct {
	OrderNo          string                     `json:"orderNo"`
	DocNo            string                     `json:"docNo"`
	Title            string                     `json:"title"`
	ApplicantID      string                     `json:"applicantId"`
	DeptID           string                     `json:"deptId"`
	SoldToBpID       string                     `json:"soldToBpId"`
	ShipToBpID       string                     `json:"shipToBpId"`
	BillToBpID       string                     `json:"billToBpId"`
	PayerBpID        string                     `json:"payerBpId"`
	Amount           float64                    `json:"amount"`
	Status           string                     `json:"status"`
	RelatedPartyCode string                     `json:"relatedPartyCode"`
	Lines            []PurchaseOrderLineRequest `json:"lines"`
	ExtData          map[string]interface{}     `json:"extData"`
	MrpRunID         string                     `json:"mrpRunId"`
}

type PurchaseReceiptLineRequest struct {
	BaseLine uint    `json:"baseLine"`
	Qty      float64 `json:"qty" binding:"required,gt=0"`
}

type PurchaseReceiptRequest struct {
	DocNo           string                       `json:"docNo"`
	PurchaseOrderID uint                         `json:"purchaseOrderId" binding:"required"`
	CostCenter      string                       `json:"costCenter"`
	Lines           []PurchaseReceiptLineRequest `json:"lines"`
	ExtData         map[string]interface{}       `json:"extData"`
}

func CreatePurchaseOrder(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	var req PurchaseOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("采购订单参数解析失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数解析失败"})
		return
	}

	order := buildPurchaseOrder(req)
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return fmt.Errorf("采购订单保存失败: %w", err)
		}
		if err := workflow.TriggerPurchaseApprovalWorkflow(tx, order); err != nil {
			return fmt.Errorf("工作流触发失败: %w", err)
		}
		return nil
	})
	if err != nil {
		log.Error("采购订单创建失败", zap.Error(err), zap.String("docNo", order.DocNo))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": order})
}

func GetPurchaseOrders(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	var orders []models.SysPurchaseOrder
	if err := db.DB.Preload("Lines").Order("created_at desc").Limit(100).Find(&orders).Error; err != nil {
		log.Error("查询采购订单失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": orders})
}

func CreatePurchaseReceipt(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	var req PurchaseReceiptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("采购收货参数解析失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数解析失败"})
		return
	}

	var receipt models.SysPurchaseReceipt
	var journalEntries []*models.SysJournalEntry
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var order models.SysPurchaseOrder
		if err := tx.Preload("Lines", func(db *gorm.DB) *gorm.DB {
			return db.Order("id asc")
		}).First(&order, req.PurchaseOrderID).Error; err != nil {
			return fmt.Errorf("采购订单不存在: %w", err)
		}
		if len(order.Lines) == 0 {
			return fmt.Errorf("采购订单没有可收货行")
		}

		receipt = buildPurchaseReceiptHeader(req, order)
		if err := tx.Create(&receipt).Error; err != nil {
			return fmt.Errorf("创建采购收货单失败: %w", err)
		}

		linesByBaseLine, err := buildReceiptLinePlan(req, order.Lines)
		if err != nil {
			return err
		}

		amountByCategory := map[string]float64{}
		for baseLine, plannedQty := range linesByBaseLine {
			orderLine := order.Lines[baseLine]
			receiptLine, category, err := receivePurchaseOrderLine(tx, receipt.ID, order.ID, uint(baseLine), orderLine, plannedQty, req.CostCenter)
			if err != nil {
				return err
			}
			receipt.TotalAmount += receiptLine.LineTotal
			amountByCategory[category] += receiptLine.LineTotal
		}

		if err := tx.Model(&receipt).Update("total_amount", receipt.TotalAmount).Error; err != nil {
			return fmt.Errorf("更新采购收货单总金额失败: %w", err)
		}

		for category, amount := range amountByCategory {
			entry, err := postJournalEntry(tx, AutoPostRequest{
				EventType:    PostingEventGoodsReceipt,
				ItemCategory: category,
				Amount:       amount,
				CostCenter:   req.CostCenter,
				Reference:    receipt.DocNo,
				PayerBpID:    order.PayerBpID,
				BillToBpID:   order.BillToBpID,
			})
			if err != nil {
				return fmt.Errorf("采购收货自动过账失败: %w", err)
			}
			journalEntries = append(journalEntries, entry)
		}

		return tx.Preload("Lines").First(&receipt, receipt.ID).Error
	})
	if err != nil {
		log.Error("采购收货失败", zap.Error(err), zap.Uint("purchaseOrderId", req.PurchaseOrderID))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": gin.H{
		"receipt":        receipt,
		"journalEntries": journalEntries,
	}})
}

func GetPurchaseReceipts(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	var receipts []models.SysPurchaseReceipt
	if err := db.DB.Preload("Lines").Order("created_at desc").Limit(100).Find(&receipts).Error; err != nil {
		log.Error("查询采购收货单失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": receipts})
}

func buildPurchaseOrder(req PurchaseOrderRequest) models.SysPurchaseOrder {
	docNo := firstNonEmpty(req.DocNo, req.OrderNo, fmt.Sprintf("PO-%d", time.Now().UnixNano()))
	status := firstNonEmpty(req.Status, "OPEN")
	lines := make([]models.SysPurchaseOrderLine, 0, len(req.Lines))
	amount := req.Amount
	if amount == 0 {
		for _, line := range req.Lines {
			amount += line.Qty * line.UnitPrice
		}
	}
	for _, line := range req.Lines {
		lines = append(lines, models.SysPurchaseOrderLine{
			ItemCode:  line.ItemCode,
			Qty:       line.Qty,
			UnitPrice: line.UnitPrice,
			BaseType:  line.BaseType,
			BaseEntry: line.BaseEntry,
			BaseLine:  line.BaseLine,
		})
	}

	return models.SysPurchaseOrder{
		OrderNo:          docNo,
		DocNo:            docNo,
		Title:            req.Title,
		ApplicantID:      firstNonEmpty(req.ApplicantID, "USER-1024"),
		DeptID:           firstNonEmpty(req.DeptID, "DEPT-PUR"),
		SoldToBpID:       req.SoldToBpID,
		ShipToBpID:       req.ShipToBpID,
		BillToBpID:       req.BillToBpID,
		PayerBpID:        req.PayerBpID,
		Amount:           amount,
		Status:           status,
		RelatedPartyCode: req.RelatedPartyCode,
		Lines:            lines,
		ExtData:          marshalExtData(req.ExtData),
		MrpRunID:         req.MrpRunID,
	}
}

func buildPurchaseReceiptHeader(req PurchaseReceiptRequest, order models.SysPurchaseOrder) models.SysPurchaseReceipt {
	return models.SysPurchaseReceipt{
		DocNo:            firstNonEmpty(req.DocNo, fmt.Sprintf("GRPO-%d", time.Now().UnixNano())),
		PurchaseOrderID:  order.ID,
		SoldToBpID:       order.SoldToBpID,
		ShipToBpID:       order.ShipToBpID,
		BillToBpID:       order.BillToBpID,
		PayerBpID:        order.PayerBpID,
		RelatedPartyCode: order.RelatedPartyCode,
		Status:           "POSTED",
		ExtData:          marshalExtData(req.ExtData),
	}
}

func buildReceiptLinePlan(req PurchaseReceiptRequest, orderLines []models.SysPurchaseOrderLine) (map[int]float64, error) {
	plan := map[int]float64{}
	if len(req.Lines) == 0 {
		for index, line := range orderLines {
			plan[index] = line.Qty
		}
		return plan, nil
	}

	for _, line := range req.Lines {
		if int(line.BaseLine) >= len(orderLines) {
			return nil, fmt.Errorf("采购订单行 %d 不存在", line.BaseLine)
		}
		if _, exists := plan[int(line.BaseLine)]; exists {
			return nil, fmt.Errorf("采购订单行 %d 重复收货", line.BaseLine)
		}
		plan[int(line.BaseLine)] = line.Qty
	}
	return plan, nil
}

func receivePurchaseOrderLine(
	tx *gorm.DB,
	receiptID uint,
	orderID uint,
	baseLine uint,
	orderLine models.SysPurchaseOrderLine,
	receivedQty float64,
	costCenter string,
) (models.SysPurchaseReceiptLine, string, error) {
	if receivedQty > orderLine.Qty {
		return models.SysPurchaseReceiptLine{}, "", fmt.Errorf("收货数量 %.4f 超过采购订单行数量 %.4f", receivedQty, orderLine.Qty)
	}
	if orderLine.UnitPrice <= 0 {
		return models.SysPurchaseReceiptLine{}, "", fmt.Errorf("采购订单行 %s 单价必须大于 0", orderLine.ItemCode)
	}

	var material models.SysMaterial
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("material_code = ?", orderLine.ItemCode).First(&material).Error; err != nil {
		return models.SysPurchaseReceiptLine{}, "", fmt.Errorf("物料 %s 不存在: %w", orderLine.ItemCode, err)
	}

	totalValue := (material.Stock * material.MovingAverageCost) + (receivedQty * orderLine.UnitPrice)
	material.Stock += receivedQty
	if material.Stock > 0 {
		material.MovingAverageCost = totalValue / material.Stock
	}
	if err := tx.Save(&material).Error; err != nil {
		return models.SysPurchaseReceiptLine{}, "", fmt.Errorf("更新物料库存失败: %w", err)
	}

	lineTotal := receivedQty * orderLine.UnitPrice
	receiptLine := models.SysPurchaseReceiptLine{
		DocID:     receiptID,
		ItemCode:  orderLine.ItemCode,
		Qty:       receivedQty,
		UnitPrice: orderLine.UnitPrice,
		LineTotal: lineTotal,
		BaseType:  models.DocumentTypePurchaseOrder,
		BaseEntry: orderID,
		BaseLine:  baseLine,
	}
	if err := tx.Create(&receiptLine).Error; err != nil {
		return models.SysPurchaseReceiptLine{}, "", fmt.Errorf("创建采购收货行失败: %w", err)
	}

	movementExtData := marshalExtData(map[string]interface{}{
		"sourceDocType": models.DocumentTypePurchaseReceipt,
		"sourceDocId":   receiptID,
		"baseType":      models.DocumentTypePurchaseOrder,
		"baseEntry":     orderID,
		"baseLine":      baseLine,
		"costCenter":    costCenter,
	})
	if err := tx.Create(&models.SysGoodsMovement{
		MvmtNo:   fmt.Sprintf("GRPO-%d-%d", receiptID, baseLine),
		MvmtType: "Receipt",
		ItemCode: orderLine.ItemCode,
		Qty:      receivedQty,
		ExtData:  movementExtData,
	}).Error; err != nil {
		return models.SysPurchaseReceiptLine{}, "", fmt.Errorf("创建库存移动失败: %w", err)
	}

	return receiptLine, normalizeItemCategory(material.Category), nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func marshalExtData(value map[string]interface{}) datatypes.JSON {
	if value == nil {
		return nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	return datatypes.JSON(data)
}
