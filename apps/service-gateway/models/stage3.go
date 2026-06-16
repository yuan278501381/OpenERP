package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// SysBatch 批次主数据
type SysBatch struct {
	gorm.Model
	BatchNo  string         `gorm:"type:varchar(64);index;comment:批次号"`
	ItemCode string         `gorm:"type:varchar(64);index;comment:物料编码"`
	MfgDate  time.Time      `gorm:"type:datetime;comment:生产日期"`
	ExpDate  time.Time      `gorm:"type:datetime;comment:过期日期"`
	ExtData  datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}

// SysSerialNumber 序列号主数据
type SysSerialNumber struct {
	gorm.Model
	SerialNo string         `gorm:"type:varchar(64);uniqueIndex;comment:序列号"`
	ItemCode string         `gorm:"type:varchar(64);index;comment:物料编码"`
	Status   string         `gorm:"type:varchar(32);comment:状态"`
	ExtData  datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}

// SysBusinessPartner 业务伙伴 (客户/供应商)
type SysBusinessPartner struct {
	gorm.Model
	BpCode  string         `gorm:"type:varchar(64);uniqueIndex;comment:业务伙伴编码"`
	BpName  string         `gorm:"type:varchar(255);comment:业务伙伴名称"`
	BpType  string         `gorm:"type:varchar(32);comment:业务伙伴类型(Customer, Vendor)"`
	TaxID   string         `gorm:"type:varchar(64);comment:税号"`
	ExtData datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}

// SysPricingCondition 定价条件
type SysPricingCondition struct {
	gorm.Model
	ConditionType string         `gorm:"type:varchar(64);index;comment:条件类型"`
	ItemCode      string         `gorm:"type:varchar(64);index;comment:物料编码"`
	BpCode        string         `gorm:"type:varchar(64);index;comment:业务伙伴编码"`
	ValidFrom     time.Time      `gorm:"type:datetime;comment:生效日期"`
	ValidTo       time.Time      `gorm:"type:datetime;comment:失效日期"`
	Price         float64        `gorm:"type:decimal(19,4);comment:价格"`
	Discount      float64        `gorm:"type:decimal(10,4);comment:折扣"`
	ExtData       datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}

// SysSalesQuotation 销售报价单
type SysSalesQuotation struct {
	gorm.Model
	DocNo       string         `gorm:"type:varchar(64);uniqueIndex;comment:单据编号"`
	BpID        string         `gorm:"type:varchar(64);index;comment:业务伙伴ID"`
	TotalAmount float64        `gorm:"type:decimal(19,4);comment:总金额"`
	ExtData     datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}

// SysSalesOrder 销售订单
type SysSalesOrder struct {
	gorm.Model
	DocNo       string              `gorm:"type:varchar(64);uniqueIndex;comment:单据编号"`
	BpID        string              `gorm:"type:varchar(64);index;comment:业务伙伴ID"`
	TotalAmount float64             `gorm:"type:decimal(19,4);comment:总金额"`
	Lines       []SysSalesOrderLine `gorm:"foreignKey:DocID;comment:订单行"`
	ExtData     datatypes.JSON      `gorm:"type:json;comment:扩展数据"`
}

// SysSalesOrderLine 销售订单行
type SysSalesOrderLine struct {
	gorm.Model
	DocID     uint           `gorm:"type:bigint;index;comment:销售订单ID"`
	ItemCode  string         `gorm:"type:varchar(64);index;comment:物料编码"`
	Qty       float64        `gorm:"type:decimal(19,4);comment:数量"`
	UnitPrice float64        `gorm:"type:decimal(19,4);comment:单价"`
	ExtData   datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}

// SysDelivery 交货单
type SysDelivery struct {
	gorm.Model
	DeliveryNo string         `gorm:"type:varchar(64);uniqueIndex;comment:交货单号"`
	OrderID    string         `gorm:"type:varchar(64);index;comment:关联订单ID"`
	BpID       string         `gorm:"type:varchar(64);index;comment:业务伙伴ID"`
	ExtData    datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}

// SysSalesReturnRequest 销售退货请求
type SysSalesReturnRequest struct {
	gorm.Model
	DocNo   string         `gorm:"type:varchar(64);uniqueIndex;comment:单据编号"`
	OrderID string         `gorm:"type:varchar(64);index;comment:关联订单ID"`
	ExtData datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}

// SysSalesReturn 销售退货单
type SysSalesReturn struct {
	gorm.Model
	DocNo   string         `gorm:"type:varchar(64);uniqueIndex;comment:单据编号"`
	OrderID string         `gorm:"type:varchar(64);index;comment:关联订单ID"`
	ExtData datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}

// SysPurchaseRequisition 采购申请
type SysPurchaseRequisition struct {
	gorm.Model
	ReqNo    string         `gorm:"type:varchar(64);uniqueIndex;comment:申请编号"`
	ItemCode string         `gorm:"type:varchar(64);index;comment:物料编码"`
	ReqQty   float64        `gorm:"type:decimal(19,4);comment:申请数量"`
	ExtData  datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}

// SysPurchaseOrderLine 采购订单行
type SysPurchaseOrderLine struct {
	gorm.Model
	DocID     uint           `gorm:"type:bigint;index;comment:采购订单ID"`
	ItemCode  string         `gorm:"type:varchar(64);index;comment:物料编码"`
	Qty       float64        `gorm:"type:decimal(19,4);comment:数量"`
	UnitPrice float64        `gorm:"type:decimal(19,4);comment:单价"`
	ExtData   datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}

// SysLandedCost 到岸成本
type SysLandedCost struct {
	gorm.Model
	DocNo           string         `gorm:"type:varchar(64);uniqueIndex;comment:单据编号"`
	ReceiptID       string         `gorm:"type:varchar(64);index;comment:收货单ID"`
	Freight         float64        `gorm:"type:decimal(19,4);comment:运费"`
	Customs         float64        `gorm:"type:decimal(19,4);comment:关税"`
	AllocatedAmount float64        `gorm:"type:decimal(19,4);comment:分摊金额"`
	ExtData         datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}

// SysPurchaseReturnRequest 采购退货请求
type SysPurchaseReturnRequest struct {
	gorm.Model
	DocNo   string         `gorm:"type:varchar(64);uniqueIndex;comment:单据编号"`
	OrderID string         `gorm:"type:varchar(64);index;comment:关联订单ID"`
	ExtData datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}

// SysPurchaseReturn 采购退货单
type SysPurchaseReturn struct {
	gorm.Model
	DocNo   string         `gorm:"type:varchar(64);uniqueIndex;comment:单据编号"`
	OrderID string         `gorm:"type:varchar(64);index;comment:关联订单ID"`
	ExtData datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}

// SysGoodsMovement 货物移动 (收发货转储)
type SysGoodsMovement struct {
	gorm.Model
	MvmtNo   string         `gorm:"type:varchar(64);uniqueIndex;comment:移动单号"`
	MvmtType string         `gorm:"type:varchar(32);index;comment:移动类型(Receipt, Issue, Transfer)"`
	ItemCode string         `gorm:"type:varchar(64);index;comment:物料编码"`
	Qty      float64        `gorm:"type:decimal(19,4);comment:数量"`
	ExtData  datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}

// SysInventoryTransferRequest 库存转储请求
type SysInventoryTransferRequest struct {
	gorm.Model
	DocNo   string         `gorm:"type:varchar(64);uniqueIndex;comment:单据编号"`
	FromWhs string         `gorm:"type:varchar(64);index;comment:源仓库"`
	ToWhs   string         `gorm:"type:varchar(64);index;comment:目标仓库"`
	ExtData datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}

// SysInventoryTransfer 库存转储单
type SysInventoryTransfer struct {
	gorm.Model
	DocNo   string         `gorm:"type:varchar(64);uniqueIndex;comment:单据编号"`
	FromWhs string         `gorm:"type:varchar(64);index;comment:源仓库"`
	ToWhs   string         `gorm:"type:varchar(64);index;comment:目标仓库"`
	ExtData datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}

// SysInventoryCounting 库存盘点
type SysInventoryCounting struct {
	gorm.Model
	CountNo    string         `gorm:"type:varchar(64);uniqueIndex;comment:盘点单号"`
	WhsCode    string         `gorm:"type:varchar(64);index;comment:仓库编码"`
	ItemCode   string         `gorm:"type:varchar(64);index;comment:物料编码"`
	SystemQty  float64        `gorm:"type:decimal(19,4);comment:系统数量"`
	CountedQty float64        `gorm:"type:decimal(19,4);comment:盘点数量"`
	Variance   float64        `gorm:"type:decimal(19,4);comment:差异数量"`
	ExtData    datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}

// SysWarehouseBin 库位数据
type SysWarehouseBin struct {
	gorm.Model
	BinCode string         `gorm:"type:varchar(64);uniqueIndex;comment:库位编码"`
	WhsCode string         `gorm:"type:varchar(64);index;comment:仓库编码"`
	Level   string         `gorm:"type:varchar(32);comment:层级"`
	Status  string         `gorm:"type:varchar(32);comment:状态"`
	ExtData datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}

// SysInventoryValuation 库存估值
type SysInventoryValuation struct {
	gorm.Model
	ItemCode        string         `gorm:"type:varchar(64);index;comment:物料编码"`
	WhsCode         string         `gorm:"type:varchar(64);index;comment:仓库编码"`
	ValuationMethod string         `gorm:"type:varchar(32);comment:计价方法(MovingAverage, Standard)"`
	Cost            float64        `gorm:"type:decimal(19,4);comment:成本"`
	ExtData         datatypes.JSON `gorm:"type:json;comment:扩展数据"`
}
