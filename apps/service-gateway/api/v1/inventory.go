package v1

import (
	"fmt"
	"net/http"
	"openerp/apps/service-gateway/db"
	"openerp/apps/service-gateway/models"
	"openerp/packages/pkg-logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GoodsMovementRequest struct {
	MvmtNo   string  `json:"mvmtNo"`
	MvmtType string  `json:"mvmtType" binding:"required"` // Receipt, Issue
	ItemCode string  `json:"itemCode" binding:"required"`
	Qty      float64 `json:"qty" binding:"required,gt=0"`
	UnitCost float64 `json:"unitCost"`
}

// CreateGoodsMovement godoc
// @Summary 创建货物移动单
// @Description 记录收发货，更新库存及计算移动平均成本
// @Tags 库存管理
// @Accept json
// @Produce json
// @Param movement body GoodsMovementRequest true "货物移动数据"
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/goods-movement [post]
func CreateGoodsMovement(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	var req GoodsMovementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数解析失败: " + err.Error()})
		return
	}

	mvmt := models.SysGoodsMovement{
		MvmtNo:   req.MvmtNo,
		MvmtType: req.MvmtType,
		ItemCode: req.ItemCode,
		Qty:      req.Qty,
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var material models.SysMaterial
		// 使用 FOR UPDATE 悲观锁保证并发情况下的库存和成本计算准确性
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("material_code = ?", req.ItemCode).First(&material).Error; err != nil {
			return fmt.Errorf("查询物料 %s 失败: %w", req.ItemCode, err)
		}

		if req.MvmtType == "Receipt" {
			// 收货：增加库存，计算移动平均成本
			totalValue := (material.Stock * material.MovingAverageCost) + (req.Qty * req.UnitCost)
			newStock := material.Stock + req.Qty
			material.Stock = newStock
			if newStock > 0 {
				material.MovingAverageCost = totalValue / newStock
			}
		} else if req.MvmtType == "Issue" {
			// 发货：扣减库存
			if material.Stock < req.Qty {
				return fmt.Errorf("物料 %s 库存不足，当前库存: %.2f，发货: %.2f", material.MaterialCode, material.Stock, req.Qty)
			}
			material.Stock -= req.Qty
		} else {
			return fmt.Errorf("不支持的移动类型: %s", req.MvmtType)
		}

		// 更新物料主数据
		if err := tx.Save(&material).Error; err != nil {
			return fmt.Errorf("更新物料库存失败: %w", err)
		}

		// 记录货物移动历史
		if err := tx.Create(&mvmt).Error; err != nil {
			return fmt.Errorf("记录货物移动单失败: %w", err)
		}

		return nil
	})

	if err != nil {
		log.Error("处理货物移动失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": mvmt})
}
