package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"openerp/apps/service-gateway/db"
	"openerp/apps/service-gateway/models"
	"openerp/packages/pkg-logger"
)

// @Summary 创建费用报销单
// @Description 创建一个新的费用报销单
// @Tags ExpenseClaim
// @Accept json
// @Produce json
// @Param claim body models.SysExpenseClaim true "报销单信息"
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/expense-claims [post]
func CreateExpenseClaim(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	var claim models.SysExpenseClaim
	if err := c.ShouldBindJSON(&claim); err != nil {
		log.Error("参数解析失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数解析失败"})
		return
	}

	if err := db.DB.Create(&claim).Error; err != nil {
		log.Error("创建报销单失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "创建报销单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": claim.ID})
}

// @Summary 获取费用报销单列表
// @Description 获取所有的费用报销单
// @Tags ExpenseClaim
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/expense-claims [get]
func GetExpenseClaims(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	var claims []models.SysExpenseClaim
	if err := db.DB.Find(&claims).Error; err != nil {
		log.Error("查询报销单失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询报销单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": claims})
}

// @Summary 更新费用报销单
// @Description 更新现有的费用报销单
// @Tags ExpenseClaim
// @Accept json
// @Produce json
// @Param id path int true "报销单ID"
// @Param claim body models.SysExpenseClaim true "报销单信息"
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/expense-claims/{id} [put]
func UpdateExpenseClaim(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	id := c.Param("id")
	var claim models.SysExpenseClaim
	if err := db.DB.First(&claim, id).Error; err != nil {
		log.Error("未找到报销单", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "未找到报销单"})
		return
	}

	var req models.SysExpenseClaim
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("参数解析失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数解析失败"})
		return
	}

	if err := db.DB.Model(&claim).Updates(req).Error; err != nil {
		log.Error("更新报销单失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新报销单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}

// @Summary 删除费用报销单
// @Description 删除指定的费用报销单
// @Tags ExpenseClaim
// @Produce json
// @Param id path int true "报销单ID"
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/expense-claims/{id} [delete]
func DeleteExpenseClaim(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	id := c.Param("id")
	if err := db.DB.Delete(&models.SysExpenseClaim{}, id).Error; err != nil {
		log.Error("删除报销单失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "删除报销单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}
