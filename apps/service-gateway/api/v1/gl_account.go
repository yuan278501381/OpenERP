package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"openerp/apps/service-gateway/db"
	"openerp/apps/service-gateway/models"
	"openerp/packages/pkg-logger"
)

// @Summary 创建总账科目
// @Description 创建一个新的总账科目
// @Tags GLAccount
// @Accept json
// @Produce json
// @Param account body models.SysGLAccount true "总账科目信息"
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/gl-accounts [post]
func CreateGLAccount(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	var account models.SysGLAccount
	if err := c.ShouldBindJSON(&account); err != nil {
		log.Error("参数解析失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数解析失败"})
		return
	}

	if err := db.DB.Create(&account).Error; err != nil {
		log.Error("创建总账科目失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "创建总账科目失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": account.ID})
}

// @Summary 获取总账科目列表
// @Description 获取所有的总账科目
// @Tags GLAccount
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/gl-accounts [get]
func GetGLAccounts(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	var accounts []models.SysGLAccount
	if err := db.DB.Find(&accounts).Error; err != nil {
		log.Error("查询总账科目失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询总账科目失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": accounts})
}

// @Summary 更新总账科目
// @Description 更新现有的总账科目
// @Tags GLAccount
// @Accept json
// @Produce json
// @Param id path int true "科目ID"
// @Param account body models.SysGLAccount true "总账科目信息"
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/gl-accounts/{id} [put]
func UpdateGLAccount(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	id := c.Param("id")
	var account models.SysGLAccount
	if err := db.DB.First(&account, id).Error; err != nil {
		log.Error("未找到总账科目", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "未找到总账科目"})
		return
	}

	var req models.SysGLAccount
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("参数解析失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数解析失败"})
		return
	}

	if err := db.DB.Model(&account).Updates(req).Error; err != nil {
		log.Error("更新总账科目失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新总账科目失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}

// @Summary 删除总账科目
// @Description 删除指定的总账科目
// @Tags GLAccount
// @Produce json
// @Param id path int true "科目ID"
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/gl-accounts/{id} [delete]
func DeleteGLAccount(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	id := c.Param("id")
	if err := db.DB.Delete(&models.SysGLAccount{}, id).Error; err != nil {
		log.Error("删除总账科目失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "删除总账科目失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}
