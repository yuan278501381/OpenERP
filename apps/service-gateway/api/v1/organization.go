package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"openerp/apps/service-gateway/db"
	"openerp/apps/service-gateway/models"
	"openerp/packages/pkg-logger"
)

// @Summary 创建组织架构
// @Description 创建一个新的组织架构节点
// @Tags Organization
// @Accept json
// @Produce json
// @Param organization body models.SysOrganization true "组织架构信息"
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/organizations [post]
func CreateOrganization(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	var org models.SysOrganization
	if err := c.ShouldBindJSON(&org); err != nil {
		log.Error("参数解析失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数解析失败"})
		return
	}

	if err := db.DB.Create(&org).Error; err != nil {
		log.Error("创建组织失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "创建组织失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": org.ID})
}

// @Summary 获取组织架构列表
// @Description 获取组织架构列表
// @Tags Organization
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/organizations [get]
func GetOrganizations(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	var orgs []models.SysOrganization
	if err := db.DB.Find(&orgs).Error; err != nil {
		log.Error("查询组织失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询组织失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": orgs})
}

// @Summary 更新组织架构
// @Description 更新现有的组织架构节点
// @Tags Organization
// @Accept json
// @Produce json
// @Param id path int true "组织ID"
// @Param organization body models.SysOrganization true "组织架构信息"
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/organizations/{id} [put]
func UpdateOrganization(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	id := c.Param("id")
	var org models.SysOrganization
	if err := db.DB.First(&org, id).Error; err != nil {
		log.Error("未找到组织", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "未找到组织"})
		return
	}

	var req models.SysOrganization
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("参数解析失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数解析失败"})
		return
	}

	if err := db.DB.Model(&org).Updates(req).Error; err != nil {
		log.Error("更新组织失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新组织失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}

// @Summary 删除组织架构
// @Description 删除指定的组织架构节点
// @Tags Organization
// @Produce json
// @Param id path int true "组织ID"
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/organizations/{id} [delete]
func DeleteOrganization(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	id := c.Param("id")
	if err := db.DB.Delete(&models.SysOrganization{}, id).Error; err != nil {
		log.Error("删除组织失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "删除组织失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}
