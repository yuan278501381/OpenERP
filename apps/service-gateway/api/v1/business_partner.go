package v1

import (
	"net/http"
	"openerp/apps/service-gateway/db"
	"openerp/apps/service-gateway/models"
	"openerp/packages/pkg-logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateBusinessPartner godoc
// @Summary 创建业务伙伴
// @Description 创建一个新的业务伙伴（客户或供应商）
// @Tags 业务伙伴管理
// @Accept json
// @Produce json
// @Param partner body models.SysBusinessPartner true "业务伙伴数据"
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/business-partners [post]
func CreateBusinessPartner(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	var partner models.SysBusinessPartner
	if err := c.ShouldBindJSON(&partner); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数解析失败"})
		return
	}
	if err := db.DB.Create(&partner).Error; err != nil {
		log.Error("创建业务伙伴失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": partner})
}

// GetBusinessPartners godoc
// @Summary 获取业务伙伴列表
// @Description 获取业务伙伴列表
// @Tags 业务伙伴管理
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/business-partners [get]
func GetBusinessPartners(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	var partners []models.SysBusinessPartner
	if err := db.DB.Find(&partners).Error; err != nil {
		log.Error("查询业务伙伴失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": partners})
}

// UpdateBusinessPartner godoc
// @Summary 更新业务伙伴
// @Description 更新业务伙伴数据
// @Tags 业务伙伴管理
// @Accept json
// @Produce json
// @Param id path int true "业务伙伴ID"
// @Param partner body models.SysBusinessPartner true "业务伙伴数据"
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/business-partners/{id} [put]
func UpdateBusinessPartner(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	id := c.Param("id")
	var partner models.SysBusinessPartner
	if err := db.DB.First(&partner, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "未找到业务伙伴"})
		return
	}
	if err := c.ShouldBindJSON(&partner); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数解析失败"})
		return
	}
	if err := db.DB.Save(&partner).Error; err != nil {
		log.Error("更新业务伙伴失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": partner})
}

// DeleteBusinessPartner godoc
// @Summary 删除业务伙伴
// @Description 根据ID删除业务伙伴
// @Tags 业务伙伴管理
// @Accept json
// @Produce json
// @Param id path int true "业务伙伴ID"
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/business-partners/{id} [delete]
func DeleteBusinessPartner(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	id := c.Param("id")
	if err := db.DB.Delete(&models.SysBusinessPartner{}, id).Error; err != nil {
		log.Error("删除业务伙伴失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}
