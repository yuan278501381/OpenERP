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

// CreateSalesOrder godoc
// @Summary 创建销售订单
// @Description 创建一个新的销售订单及其行项目
// @Tags 销售管理
// @Accept json
// @Produce json
// @Param order body models.SysSalesOrder true "销售订单数据"
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/sales-orders [post]
func CreateSalesOrder(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	var order models.SysSalesOrder
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数解析失败"})
		return
	}
	if err := db.DB.Create(&order).Error; err != nil {
		log.Error("创建销售订单失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": order})
}

// GetSalesOrders godoc
// @Summary 获取销售订单列表
// @Description 获取包含行项目的销售订单列表
// @Tags 销售管理
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/sales-orders [get]
func GetSalesOrders(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	var orders []models.SysSalesOrder
	if err := db.DB.Preload("Lines").Find(&orders).Error; err != nil {
		log.Error("查询销售订单失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": orders})
}

// UpdateSalesOrder godoc
// @Summary 更新销售订单
// @Description 更新销售订单及其关联项目
// @Tags 销售管理
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Param order body models.SysSalesOrder true "销售订单数据"
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/sales-orders/{id} [put]
func UpdateSalesOrder(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	id := c.Param("id")
	var order models.SysSalesOrder
	if err := db.DB.Preload("Lines").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "未找到销售订单"})
		return
	}
	// 在复杂ERP场景中，更安全的做法是分步更新主表和子表，或使用专门的事务逻辑，
	// 此处为简化演示，直接更新并替换
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数解析失败"})
		return
	}
	
	// 更新主数据及关联数据
	if err := db.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&order).Error; err != nil {
		log.Error("更新销售订单失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": order})
}

// DeleteSalesOrder godoc
// @Summary 删除销售订单
// @Description 根据ID删除销售订单（及其行项目）
// @Tags 销售管理
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/sales-orders/{id} [delete]
func DeleteSalesOrder(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	id := c.Param("id")
	// 注意配置好软删除或级联删除
	if err := db.DB.Delete(&models.SysSalesOrder{}, id).Error; err != nil {
		log.Error("删除销售订单失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}
