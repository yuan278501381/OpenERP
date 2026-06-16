package models

import (
	"gorm.io/gorm"
	"time"
)

// SysTask 统一待办中心模型
type SysTask struct {
	ID        uint           `gorm:"primaryKey" json:"-"`
	TaskID    string         `gorm:"uniqueIndex;type:varchar(50)" json:"task_id"`
	Title     string         `gorm:"type:varchar(200)" json:"title"`
	Type      string         `gorm:"type:varchar(50)" json:"type"`
	Node      string         `gorm:"type:varchar(100)" json:"node"`
	SlaStatus string         `gorm:"type:varchar(20)" json:"sla_status"` // Normal, Warning, Critical
	TenantID  string         `gorm:"index;type:varchar(50)" json:"-"`    // 多租户隔离
	UserID    string         `gorm:"index;type:varchar(50)" json:"-"`    // 归属人
	ExtraData string         `gorm:"type:jsonb" json:"extra_data"`       // 世界级 ERP 必须具备的 JSONB 动态扩展字段
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
