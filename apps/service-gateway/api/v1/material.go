package v1

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"openerp/apps/service-gateway/db"
	"openerp/apps/service-gateway/models"
	"openerp/packages/pkg-logger"
	"go.uber.org/zap"
)

// GetMaterials godoc
// @Summary 获取物料主数据列表
// @Description 从数据库拉取全部的物料主数据，支持包含基于 JSONB 的动态扩展字段。
// @Tags 物料管理 (Material Master)
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/materials [get]
func GetMaterials(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	
	var materials []models.SysMaterial
	result := db.DB.Find(&materials)
	if result.Error != nil {
		log.Error("物料数据查询失败", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": materials,
	})
}
