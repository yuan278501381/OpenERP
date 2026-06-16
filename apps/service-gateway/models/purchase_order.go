package models

import (
	"gorm.io/gorm"
	"gorm.io/datatypes"
)

// SysPurchaseOrder 采购申请单实体 (领域驱动设计)
type SysPurchaseOrder struct {
	gorm.Model
	OrderNo     string                 `gorm:"type:varchar(64);uniqueIndex;comment:采购单号(旧)"`
	DocNo       string                 `gorm:"type:varchar(64);uniqueIndex;comment:采购单号"`
	Title       string                 `gorm:"type:varchar(255);comment:申请标题"`
	ApplicantID string                 `gorm:"type:varchar(64);comment:申请人ID"`
	DeptID      string                 `gorm:"type:varchar(64);comment:申请部门ID"`
	BpID        string                 `gorm:"type:varchar(64);index;comment:业务伙伴ID"`
	Amount      float64                `gorm:"type:decimal(10,2);comment:预估总金额"`
	Status      string                 `gorm:"type:varchar(32);default:'DRAFT';comment:单据状态"`
	Lines       []SysPurchaseOrderLine `gorm:"foreignKey:DocID;comment:采购订单行"`
	
	// 面向未来的设计：无限扩展的自定义字段 (动态表单机制)
	// 在 PostgreSQL 中将映射为原生 JSONB，支持高效索引与多维数据查询。
	ExtData     datatypes.JSON `gorm:"type:json;comment:动态表单扩展字段"`
	
	// 面向未来的设计：MRP (物料需求计划) 关联预留
	// 未来 MRP 计算出的缺口推出来的采购单，会携带此追溯 ID。
	MrpRunID    string         `gorm:"type:varchar(64);index;comment:关联的MRP运算ID"`
}
