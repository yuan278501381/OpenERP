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

// CreateProductionOrder godoc
// @Summary 创建生产订单
// @Description 创建一个新的生产订单
// @Tags 生产管理
// @Accept json
// @Produce json
// @Param order body models.SysProductionOrder true "生产订单数据"
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/production-orders [post]
func CreateProductionOrder(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	var order models.SysProductionOrder
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数解析失败"})
		return
	}
	if err := db.DB.Create(&order).Error; err != nil {
		log.Error("创建生产订单失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": order})
}

// GetProductionOrders godoc
// @Summary 获取生产订单列表
// @Description 获取生产订单列表
// @Tags 生产管理
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/production-orders [get]
func GetProductionOrders(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	var orders []models.SysProductionOrder
	if err := db.DB.Find(&orders).Error; err != nil {
		log.Error("查询生产订单失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": orders})
}

// UpdateProductionOrder godoc
// @Summary 更新生产订单
// @Description 更新生产订单
// @Tags 生产管理
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Param order body models.SysProductionOrder true "生产订单数据"
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/production-orders/{id} [put]
func UpdateProductionOrder(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	id := c.Param("id")
	var order models.SysProductionOrder
	if err := db.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "未找到生产订单"})
		return
	}
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数解析失败"})
		return
	}
	if err := db.DB.Save(&order).Error; err != nil {
		log.Error("更新生产订单失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": order})
}

// DeleteProductionOrder godoc
// @Summary 删除生产订单
// @Description 删除生产订单
// @Tags 生产管理
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/production-orders/{id} [delete]
func DeleteProductionOrder(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	id := c.Param("id")
	if err := db.DB.Delete(&models.SysProductionOrder{}, id).Error; err != nil {
		log.Error("删除生产订单失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}

// ReportCompletionRequest 报工请求参数
type ReportCompletionRequest struct {
	CompletedQty float64 `json:"completed_qty"`
}

// ReportCompletion godoc
// @Summary 生产订单报工
// @Description 生产订单报工，增加成品库存，扣减BOM组件库存
// @Tags 生产管理
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Param request body ReportCompletionRequest true "报工参数"
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/production-orders/{id}/report-completion [post]
func ReportCompletion(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	id := c.Param("id")
	var req ReportCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数解析失败"})
		return
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 获取生产订单
		var order models.SysProductionOrder
		if err := tx.First(&order, id).Error; err != nil {
			return err
		}

		// 2. 获取成品物料
		var fgMaterial models.SysMaterial
		if err := tx.Where("material_code = ?", order.ItemCode).First(&fgMaterial).Error; err != nil {
			return err
		}

		// 3. 增加成品库存
		if err := tx.Model(&models.SysMaterial{}).Where("id = ?", fgMaterial.ID).Update("stock", gorm.Expr("stock + ?", req.CompletedQty)).Error; err != nil {
			return err
		}

		// 4. 获取BOM，扣减组件库存 (仅限单层BOM扣减)
		var boms []models.SysBOM
		if err := tx.Where("parent_material_id = ?", fgMaterial.ID).Find(&boms).Error; err != nil {
			return err
		}

		for _, bom := range boms {
			consumedQty := req.CompletedQty * bom.Quantity * (1 + bom.ScrapRate)
			if err := tx.Model(&models.SysMaterial{}).Where("id = ?", bom.ChildMaterialID).Update("stock", gorm.Expr("stock - ?", consumedQty)).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Error("报工失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "报工失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}
