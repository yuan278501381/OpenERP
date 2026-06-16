package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// SysMaterial 物料主数据
type SysMaterial struct {
	gorm.Model
	MaterialCode string         `gorm:"type:varchar(64);uniqueIndex;not null;comment:物料编码"`
	MaterialName string         `gorm:"type:varchar(255);not null;comment:物料名称"`
	Category     string         `gorm:"type:varchar(64);index;comment:物料分类(如:原材料、半成品、成品)"`
	Unit         string         `gorm:"type:varchar(32);comment:基本计量单位"`
	IsActive     bool           `gorm:"default:true;comment:是否启用"`

	// 面向未来的设计：无限扩展的自定义字段 (JSONB)
	ExtData      datatypes.JSON `gorm:"type:json;comment:动态扩展属性"`
}
