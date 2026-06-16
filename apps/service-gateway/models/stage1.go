package models

import (
	"gorm.io/gorm"
	"gorm.io/datatypes"
	"time"
)

// SysOrganization 系统组织架构
type SysOrganization struct {
	gorm.Model
	OrgCode  string         `gorm:"type:varchar(64);uniqueIndex;comment:组织代码"`
	OrgName  string         `gorm:"type:varchar(255);comment:组织名称"`
	OrgType  string         `gorm:"type:varchar(32);comment:组织类型(Legal, Business, Profit Center)"`
	ParentID uint           `gorm:"comment:父级组织ID"`
	ExtData  datatypes.JSON `gorm:"type:json;comment:动态表单扩展字段"`
}

// SysFinancialPeriod 财务会计期间
type SysFinancialPeriod struct {
	gorm.Model
	Year      int            `gorm:"comment:会计年度"`
	Period    int            `gorm:"comment:会计期间"`
	StartDate time.Time      `gorm:"comment:开始日期"`
	EndDate   time.Time      `gorm:"comment:结束日期"`
	IsClosed  bool           `gorm:"comment:是否已结账"`
	ExtData   datatypes.JSON `gorm:"type:json;comment:动态表单扩展字段"`
}

// SysPaymentTerm 付款条件
type SysPaymentTerm struct {
	gorm.Model
	TermCode        string         `gorm:"type:varchar(64);uniqueIndex;comment:条款代码"`
	Description     string         `gorm:"type:varchar(255);comment:条款描述"`
	DueDays         int            `gorm:"comment:到期天数"`
	DiscountDays    int            `gorm:"comment:折扣天数"`
	DiscountPercent float64        `gorm:"type:decimal(5,2);comment:折扣百分比"`
	ExtData         datatypes.JSON `gorm:"type:json;comment:动态表单扩展字段"`
}

// SysUoMGroup 计量单位组
type SysUoMGroup struct {
	gorm.Model
	GroupCode        string         `gorm:"type:varchar(64);uniqueIndex;comment:单位组代码"`
	BaseUoM          string         `gorm:"type:varchar(32);comment:基础单位"`
	AltUoM           string         `gorm:"type:varchar(32);comment:替代单位"`
	ConversionFactor float64        `gorm:"type:decimal(10,4);comment:转换因子"`
	ExtData          datatypes.JSON `gorm:"type:json;comment:动态表单扩展字段"`
}

// SysAccountDetermination 科目决定规则
type SysAccountDetermination struct {
	gorm.Model
	TransactionType string         `gorm:"type:varchar(64);comment:交易类型"`
	ItemCategory    string         `gorm:"type:varchar(64);comment:物料类别"`
	DebitAccount    string         `gorm:"type:varchar(64);comment:借方科目"`
	CreditAccount   string         `gorm:"type:varchar(64);comment:贷方科目"`
	ExtData         datatypes.JSON `gorm:"type:json;comment:动态表单扩展字段"`
}

// SysNotificationTemplate 通知模板
type SysNotificationTemplate struct {
	gorm.Model
	TemplateCode string         `gorm:"type:varchar(64);uniqueIndex;comment:模板代码"`
	Platform     string         `gorm:"type:varchar(32);comment:平台(WeCom, DingTalk)"`
	Content      string         `gorm:"type:text;comment:模板内容"`
	IsActive     bool           `gorm:"comment:是否启用"`
	ExtData      datatypes.JSON `gorm:"type:json;comment:动态表单扩展字段"`
}
