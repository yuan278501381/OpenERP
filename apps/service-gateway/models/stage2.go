package models

import (
	"gorm.io/gorm"
	"gorm.io/datatypes"
	"time"
)

// SysGLAccount 总账科目
type SysGLAccount struct {
	gorm.Model
	AccountCode string         `gorm:"type:varchar(64);uniqueIndex;comment:科目代码"`
	AccountName string         `gorm:"type:varchar(255);comment:科目名称"`
	AccountType string         `gorm:"type:varchar(32);comment:科目类型(Asset, Liability, Equity, Revenue, Expense)"`
	ExtData     datatypes.JSON `gorm:"type:json;comment:动态表单扩展字段"`
}

// SysJournalEntry 日记账分录头
type SysJournalEntry struct {
	gorm.Model
	EntryNo     string                `gorm:"type:varchar(64);uniqueIndex;comment:分录号"`
	PostingDate time.Time             `gorm:"comment:过账日期"`
	DocDate     time.Time             `gorm:"comment:单据日期"`
	TotalAmount float64               `gorm:"type:decimal(15,2);comment:总金额"`
	Lines       []SysJournalEntryLine `gorm:"foreignKey:EntryID;comment:分录行"`
	ExtData     datatypes.JSON        `gorm:"type:json;comment:动态表单扩展字段"`
}

// SysJournalEntryLine 日记账分录行
type SysJournalEntryLine struct {
	gorm.Model
	EntryID      uint           `gorm:"comment:分录头ID"`
	AccountCode  string         `gorm:"type:varchar(64);comment:科目代码"`
	PayerBpID    string         `gorm:"type:varchar(64);index;comment:付款方ID"`
	BillToBpID   string         `gorm:"type:varchar(64);index;comment:收票方ID"`
	DebitAmount  float64        `gorm:"type:decimal(15,2);comment:借方金额"`
	CreditAmount float64        `gorm:"type:decimal(15,2);comment:贷方金额"`
	CostCenter   string         `gorm:"type:varchar(64);comment:成本中心"`
	ExtData      datatypes.JSON `gorm:"type:json;comment:动态表单扩展字段"`
}

// SysCostCenter 成本中心
type SysCostCenter struct {
	gorm.Model
	CenterCode string         `gorm:"type:varchar(64);uniqueIndex;comment:成本中心代码"`
	CenterName string         `gorm:"type:varchar(255);comment:成本中心名称"`
	ManagerID  string         `gorm:"type:varchar(64);comment:负责人ID"`
	ExtData    datatypes.JSON `gorm:"type:json;comment:动态表单扩展字段"`
}

// SysBudget 预算
type SysBudget struct {
	gorm.Model
	DepartmentID   string         `gorm:"type:varchar(64);comment:部门ID"`
	AccountCode    string         `gorm:"type:varchar(64);comment:科目代码"`
	Period         string         `gorm:"type:varchar(32);comment:期间"`
	BudgetAmount   float64        `gorm:"type:decimal(15,2);comment:预算金额"`
	ConsumedAmount float64        `gorm:"type:decimal(15,2);comment:已耗用金额"`
	ExtData        datatypes.JSON `gorm:"type:json;comment:动态表单扩展字段"`
}

// SysARInvoice 应收发票
type SysARInvoice struct {
	gorm.Model
	InvoiceNo   string         `gorm:"type:varchar(64);uniqueIndex;comment:发票号"`
	SoldToBpID  string         `gorm:"type:varchar(64);index;comment:售达方ID"`
	ShipToBpID  string         `gorm:"type:varchar(64);index;comment:送达方ID"`
	BillToBpID  string         `gorm:"type:varchar(64);index;comment:收票方ID"`
	PayerBpID   string         `gorm:"type:varchar(64);index;comment:付款方ID"`
	TotalAmount float64        `gorm:"type:decimal(15,2);comment:总金额"`
	TaxAmount   float64        `gorm:"type:decimal(15,2);comment:税额"`
	Status      string         `gorm:"type:varchar(32);comment:状态"`
	ExtData     datatypes.JSON `gorm:"type:json;comment:动态表单扩展字段"`
}

// SysAPInvoice 应付发票
type SysAPInvoice struct {
	gorm.Model
	InvoiceNo   string         `gorm:"type:varchar(64);uniqueIndex;comment:发票号"`
	SoldToBpID  string         `gorm:"type:varchar(64);index;comment:售达方ID"`
	ShipToBpID  string         `gorm:"type:varchar(64);index;comment:送达方ID"`
	BillToBpID  string         `gorm:"type:varchar(64);index;comment:收票方ID"`
	PayerBpID   string         `gorm:"type:varchar(64);index;comment:付款方ID"`
	TotalAmount float64        `gorm:"type:decimal(15,2);comment:总金额"`
	TaxAmount   float64        `gorm:"type:decimal(15,2);comment:税额"`
	Status      string         `gorm:"type:varchar(32);comment:状态"`
	ExtData     datatypes.JSON `gorm:"type:json;comment:动态表单扩展字段"`
}

// SysPayment 收付款
type SysPayment struct {
	gorm.Model
	PaymentNo string         `gorm:"type:varchar(64);uniqueIndex;comment:支付单号"`
	BpID      string         `gorm:"type:varchar(64);comment:业务伙伴ID"`
	Amount    float64        `gorm:"type:decimal(15,2);comment:金额"`
	Direction string         `gorm:"type:varchar(32);comment:方向(Inbound/Outbound)"`
	ExtData   datatypes.JSON `gorm:"type:json;comment:动态表单扩展字段"`
}

// SysReconciliation 对账单
type SysReconciliation struct {
	gorm.Model
	ReconNo          string         `gorm:"type:varchar(64);uniqueIndex;comment:对账单号"`
	InvoiceID        uint           `gorm:"comment:发票ID"`
	PaymentID        uint           `gorm:"comment:支付ID"`
	ReconciledAmount float64        `gorm:"type:decimal(15,2);comment:核销金额"`
	ExtData          datatypes.JSON `gorm:"type:json;comment:动态表单扩展字段"`
}

// SysExpenseClaim 费用报销单
type SysExpenseClaim struct {
	gorm.Model
	ClaimNo    string         `gorm:"type:varchar(64);uniqueIndex;comment:报销单号"`
	EmployeeID string         `gorm:"type:varchar(64);comment:员工ID"`
	Amount     float64        `gorm:"type:decimal(15,2);comment:报销金额"`
	Reason     string         `gorm:"type:text;comment:报销事由"`
	Status     string         `gorm:"type:varchar(32);comment:状态(Draft, Approved, Paid)"`
	ExtData    datatypes.JSON `gorm:"type:json;comment:动态表单扩展字段"`
}

// SysFixedAsset 固定资产
type SysFixedAsset struct {
	gorm.Model
	AssetCode               string         `gorm:"type:varchar(64);uniqueIndex;comment:资产代码"`
	AssetName               string         `gorm:"type:varchar(255);comment:资产名称"`
	AcquisitionValue        float64        `gorm:"type:decimal(15,2);comment:原值"`
	AccumulatedDepreciation float64        `gorm:"type:decimal(15,2);comment:累计折旧"`
	NetValue                float64        `gorm:"type:decimal(15,2);comment:净值"`
	ExtData                 datatypes.JSON `gorm:"type:json;comment:动态表单扩展字段"`
}

// SysBillOfExchange 汇票
type SysBillOfExchange struct {
	gorm.Model
	BillNo       string         `gorm:"type:varchar(64);uniqueIndex;comment:汇票号"`
	BpID         string         `gorm:"type:varchar(64);comment:业务伙伴ID"`
	Amount       float64        `gorm:"type:decimal(15,2);comment:金额"`
	MaturityDate time.Time      `gorm:"comment:到期日"`
	Status       string         `gorm:"type:varchar(32);comment:状态(Received, Endorsed, Discounted)"`
	ExtData      datatypes.JSON `gorm:"type:json;comment:动态表单扩展字段"`
}

// SysExchangeRate 汇率
type SysExchangeRate struct {
	gorm.Model
	Currency       string         `gorm:"type:varchar(16);comment:本位币"`
	TargetCurrency string         `gorm:"type:varchar(16);comment:目标币"`
	Rate           float64        `gorm:"type:decimal(10,6);comment:汇率"`
	ValidFrom      time.Time      `gorm:"comment:生效日期"`
	ExtData        datatypes.JSON `gorm:"type:json;comment:动态表单扩展字段"`
}

// SysTaxGroup 税组
type SysTaxGroup struct {
	gorm.Model
	TaxCode     string         `gorm:"type:varchar(64);uniqueIndex;comment:税码"`
	Rate        float64        `gorm:"type:decimal(5,4);comment:税率"`
	Description string         `gorm:"type:varchar(255);comment:描述"`
	ExtData     datatypes.JSON `gorm:"type:json;comment:动态表单扩展字段"`
}

// SysTaxInvoice 税务发票
type SysTaxInvoice struct {
	gorm.Model
	TaxInvoiceNo string         `gorm:"type:varchar(64);uniqueIndex;comment:税务发票号"`
	Amount       float64        `gorm:"type:decimal(15,2);comment:金额"`
	TaxAmount    float64        `gorm:"type:decimal(15,2);comment:税额"`
	Type         string         `gorm:"type:varchar(32);comment:类型(VatIn, VatOut)"`
	Status       string         `gorm:"type:varchar(32);comment:状态"`
	ExtData      datatypes.JSON `gorm:"type:json;comment:动态表单扩展字段"`
}
