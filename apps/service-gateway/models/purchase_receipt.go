package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// SysPurchaseReceipt 采购收货单 (GRPO)
type SysPurchaseReceipt struct {
	gorm.Model
	DocNo            string                   `gorm:"type:varchar(64);uniqueIndex;comment:采购收货单号"`
	PurchaseOrderID  uint                     `gorm:"type:bigint;index;comment:来源采购订单ID"`
	SoldToBpID       string                   `gorm:"type:varchar(64);index;comment:售达方ID"`
	ShipToBpID       string                   `gorm:"type:varchar(64);index;comment:送达方ID"`
	BillToBpID       string                   `gorm:"type:varchar(64);index;comment:收票方ID"`
	PayerBpID        string                   `gorm:"type:varchar(64);index;comment:付款方ID"`
	RelatedPartyCode string                   `gorm:"type:varchar(64);comment:关联方编码"`
	Status           string                   `gorm:"type:varchar(32);default:'POSTED';comment:单据状态"`
	TotalAmount      float64                  `gorm:"type:decimal(19,4);comment:收货总金额"`
	Lines            []SysPurchaseReceiptLine `gorm:"foreignKey:DocID;comment:采购收货行"`
	ExtData          datatypes.JSON           `gorm:"type:json;comment:扩展数据"`
}

// SysPurchaseReceiptLine 采购收货行
type SysPurchaseReceiptLine struct {
	gorm.Model
	DocID     uint           `gorm:"type:bigint;index;comment:采购收货单ID"`
	ItemCode  string         `gorm:"type:varchar(64);index;comment:物料编码"`
	Qty       float64        `gorm:"type:decimal(19,4);comment:收货数量"`
	UnitPrice float64        `gorm:"type:decimal(19,4);comment:采购单价"`
	LineTotal float64        `gorm:"type:decimal(19,4);comment:行金额"`
	BaseType  int            `gorm:"type:int;comment:源单据类型"`
	BaseEntry uint           `gorm:"type:bigint;comment:源单据内部ID"`
	BaseLine  uint           `gorm:"type:bigint;comment:源单据行号"`
	ExtData   datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}
