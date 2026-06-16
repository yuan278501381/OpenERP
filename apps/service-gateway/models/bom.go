package models

import (
	"gorm.io/gorm"
)

// SysBOM 物料清单 (Bill of Materials)
type SysBOM struct {
	gorm.Model
	ParentMaterialID uint    `gorm:"index;not null;comment:父物料ID"`
	ChildMaterialID  uint    `gorm:"index;not null;comment:子物料ID"`
	Quantity         float64 `gorm:"type:decimal(10,4);not null;comment:标准用量"`
	ScrapRate        float64 `gorm:"type:decimal(5,4);default:0;comment:损耗率"`
	Version          string  `gorm:"type:varchar(32);default:'v1.0';comment:BOM版本"`
	IsActive         bool    `gorm:"default:true;comment:是否生效"`
}
