package audit

import "admin-platform/backend/internal/pkg/timex"

type LoginLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt timex.DateTime `gorm:"type:timestamptz;autoCreateTime" json:"createdAt"`
	Username  string         `gorm:"size:64" json:"username"`
	IP        string         `gorm:"size:64" json:"ip"`
	Status    string         `gorm:"size:32" json:"status"`
	Message   string         `gorm:"size:255" json:"message"`
}

func (LoginLog) TableName() string { return "sys_login_logs" }

type OperationLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt timex.DateTime `gorm:"type:timestamptz;autoCreateTime" json:"createdAt"`
	UserID    uint           `json:"userId"`
	Module    string         `gorm:"size:64" json:"module"`
	Action    string         `gorm:"size:64" json:"action"`
	Method    string         `gorm:"size:16" json:"method"`
	Path      string         `gorm:"size:255" json:"path"`
	IP        string         `gorm:"size:64" json:"ip"`
	Status    int            `json:"status"`
}

func (OperationLog) TableName() string { return "sys_operation_logs" }
