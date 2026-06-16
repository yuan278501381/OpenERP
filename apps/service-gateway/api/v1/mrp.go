package v1

import (
	"net/http"
	"openerp/apps/service-gateway/db"
	"openerp/apps/service-gateway/models"
	"openerp/packages/pkg-logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RunMRPRequest MRP请求参数
type RunMRPRequest struct {
	SalesOrderID uint    `json:"sales_order_id"`
	ItemCode     string  `json:"item_code"`
	Quantity     float64 `json:"quantity"`
}

// MRPResult MRP运算结果
type MRPResult struct {
	MaterialID   uint    `json:"material_id"`
	MaterialCode string  `json:"material_code"`
	MaterialName string  `json:"material_name"`
	RequiredQty  float64 `json:"required_qty"`
}

// RunMRP godoc
// @Summary 运行MRP运算
// @Description 运行MRP运算，根据销售订单或物料编码递归计算物料需求
// @Tags 生产管理
// @Accept json
// @Produce json
// @Param request body RunMRPRequest true "MRP请求参数"
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/mrp/run [post]
func RunMRP(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	var req RunMRPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数解析失败"})
		return
	}

	mrpMap := make(map[uint]float64)

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if req.SalesOrderID > 0 {
			var lines []models.SysSalesOrderLine
			if err := tx.Where("doc_id = ?", req.SalesOrderID).Find(&lines).Error; err != nil {
				return err
			}
			for _, line := range lines {
				var material models.SysMaterial
				if err := tx.Where("material_code = ?", line.ItemCode).First(&material).Error; err != nil {
					continue // 若找不到物料，跳过
				}
				if err := calculateBOM(tx, material.ID, line.Qty, mrpMap); err != nil {
					return err
				}
			}
		} else if req.ItemCode != "" && req.Quantity > 0 {
			var material models.SysMaterial
			if err := tx.Where("material_code = ?", req.ItemCode).First(&material).Error; err != nil {
				return err
			}
			if err := calculateBOM(tx, material.ID, req.Quantity, mrpMap); err != nil {
				return err
			}
		} else {
			return gorm.ErrInvalidData
		}
		return nil
	})

	if err != nil {
		log.Error("MRP运算失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "MRP运算失败"})
		return
	}

	var results []MRPResult
	for matID, qty := range mrpMap {
		var material models.SysMaterial
		db.DB.First(&material, matID)
		results = append(results, MRPResult{
			MaterialID:   matID,
			MaterialCode: material.MaterialCode,
			MaterialName: material.MaterialName,
			RequiredQty:  qty,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": results})
}

// calculateBOM 递归计算BOM需求
func calculateBOM(tx *gorm.DB, materialID uint, reqQty float64, mrpMap map[uint]float64) error {
	var boms []models.SysBOM
	if err := tx.Where("parent_material_id = ?", materialID).Find(&boms).Error; err != nil {
		return err
	}

	// 如果没有子BOM，说明是底层物料，加入需求汇总
	if len(boms) == 0 {
		mrpMap[materialID] += reqQty
		return nil
	}

	// 有子BOM，递归计算子件需求
	for _, bom := range boms {
		childReqQty := reqQty * bom.Quantity * (1 + bom.ScrapRate)
		if err := calculateBOM(tx, bom.ChildMaterialID, childReqQty, mrpMap); err != nil {
			return err
		}
	}
	return nil
}
