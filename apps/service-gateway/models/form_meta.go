package models

import (
	"gorm.io/gorm"
)

// SysFormField 动态表单元数据定义
type SysFormField struct {
	gorm.Model
	FormType    string `gorm:"type:varchar(64);index;not null;comment:表单类型(如:PurchaseOrder, Material)"`
	FieldKey    string `gorm:"type:varchar(64);not null;comment:字段英文标识(存入JSONB的Key)"`
	Label       string `gorm:"type:varchar(128);not null;comment:前端展示标签名称"`
	FieldType   string `gorm:"type:varchar(32);not null;comment:字段类型(String, Number, Date, Select)"`
	IsRequired  bool   `gorm:"default:false;comment:是否必填"`
	OptionsJSON string `gorm:"type:json;comment:下拉框选项(如有)"`
	SortOrder   int    `gorm:"default:0;comment:排序权重"`
}
