package v1

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"openerp/apps/service-gateway/db"
	"openerp/apps/service-gateway/models"
	"openerp/packages/pkg-logger"
	"go.uber.org/zap"
)

// GetFormMeta godoc
// @Summary 获取特定业务单据的动态表单元数据
// @Description 根据表单类型（如: Material, PurchaseOrder）拉取注册在后端的所有自定义扩展字段配置。
// @Tags 动态表单引擎 (Form Meta Engine)
// @Accept json
// @Produce json
// @Param type path string true "表单类型 (e.g., Material)"
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/forms/{type}/meta [get]
func GetFormMeta(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	formType := c.Param("type")

	var metaFields []models.SysFormField
	result := db.DB.Where("form_type = ?", formType).Order("sort_order asc").Find(&metaFields)
	if result.Error != nil {
		log.Error("获取表单元数据失败", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": metaFields,
	})
}

// AddFormMeta godoc
// @Summary 为特定表单新增动态字段
// @Description 允许企业在运行时给单据增加新字段（如：长度、颜色、等级），基于 JSONB 落地。
// @Tags 动态表单引擎 (Form Meta Engine)
// @Accept json
// @Produce json
// @Param type path string true "表单类型 (e.g., Material)"
// @Param field body models.SysFormField true "字段配置参数"
// @Success 200 {object} map[string]interface{}
// @Router /openerp/v1/forms/{type}/meta [post]
func AddFormMeta(c *gin.Context) {
	log := logger.Ctx(c.Request.Context())
	formType := c.Param("type")

	var newField models.SysFormField
	if err := c.ShouldBindJSON(&newField); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数解析失败"})
		return
	}

	newField.FormType = formType

	if err := db.DB.Create(&newField).Error; err != nil {
		log.Error("新增动态字段失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "保存失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "动态字段配置成功",
		"data": newField,
	})
}
