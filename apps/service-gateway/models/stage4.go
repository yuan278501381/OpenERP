package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// SysWorkCenter 工作中心
type SysWorkCenter struct {
	gorm.Model
	CenterCode string         `gorm:"type:varchar(64);uniqueIndex;comment:工作中心编码"`
	CenterName string         `gorm:"type:varchar(255);comment:工作中心名称"`
	Capacity   float64        `gorm:"type:decimal(19,4);comment:产能"`
	ExtData    datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}

// SysResource 资源 (机器、人工)
type SysResource struct {
	gorm.Model
	ResourceCode string         `gorm:"type:varchar(64);uniqueIndex;comment:资源编码"`
	ResourceType string         `gorm:"type:varchar(32);comment:资源类型(Machine, Labor)"`
	CostPerHour  float64        `gorm:"type:decimal(19,4);comment:每小时成本"`
	ExtData      datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}

// SysRouting 工艺路线
type SysRouting struct {
	gorm.Model
	RoutingCode string         `gorm:"type:varchar(64);uniqueIndex;comment:工艺路线编码"`
	ItemCode    string         `gorm:"type:varchar(64);index;comment:物料编码"`
	Operations  datatypes.JSON `gorm:"type:json;comment:工序列表"`
	ExtData     datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}

// SysProductionOrder 生产订单
type SysProductionOrder struct {
	gorm.Model
	OrderNo    string         `gorm:"type:varchar(64);uniqueIndex;comment:生产订单号"`
	ItemCode   string         `gorm:"type:varchar(64);index;comment:物料编码"`
	PlannedQty float64        `gorm:"type:decimal(19,4);comment:计划数量"`
	Status     string         `gorm:"type:varchar(32);comment:状态"`
	ExtData    datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}

// SysInspectionLot 检验批
type SysInspectionLot struct {
	gorm.Model
	LotNo        string         `gorm:"type:varchar(64);uniqueIndex;comment:检验批号"`
	ItemCode     string         `gorm:"type:varchar(64);index;comment:物料编码"`
	InspectedQty float64        `gorm:"type:decimal(19,4);comment:检验数量"`
	Status       string         `gorm:"type:varchar(32);comment:状态"`
	ExtData      datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}

// SysEquipment 设备
type SysEquipment struct {
	gorm.Model
	EquipCode string         `gorm:"type:varchar(64);uniqueIndex;comment:设备编码"`
	EquipName string         `gorm:"type:varchar(255);comment:设备名称"`
	Status    string         `gorm:"type:varchar(32);comment:状态"`
	ExtData   datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}

// SysEmployee 员工
type SysEmployee struct {
	gorm.Model
	EmpNo      string         `gorm:"type:varchar(64);uniqueIndex;comment:员工号"`
	EmpName    string         `gorm:"type:varchar(255);comment:员工姓名"`
	Department string         `gorm:"type:varchar(128);comment:部门"`
	ExtData    datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}
