package audit

import "admin-platform/backend/internal/pkg/timex"

// LoginLog 记录登录尝试结果，用于安全审计和排查异常登录。
type LoginLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`                             // 登录日志主键。
	CreatedAt timex.DateTime `gorm:"type:timestamptz;autoCreateTime" json:"createdAt"` // 登录尝试发生时间。
	Username  string         `gorm:"size:64" json:"username"`                          // 用户提交的账号，失败时也保留便于审计。
	IP        string         `gorm:"size:64" json:"ip"`                                // 客户端 IP。
	Status    string         `gorm:"size:32" json:"status"`                            // 登录结果，例如 success 或 failed。
	Message   string         `gorm:"size:255" json:"message"`                          // 面向审计展示的结果说明。
}

// TableName 固定登录日志表名。
func (LoginLog) TableName() string { return "sys_login_logs" }

// OperationLog 记录私有写操作请求，用于追踪谁在什么模块执行了什么操作。
type OperationLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`                             // 操作日志主键。
	CreatedAt timex.DateTime `gorm:"type:timestamptz;autoCreateTime" json:"createdAt"` // 操作发生时间。
	UserID    uint           `json:"userId"`                                           // 操作人用户 ID。
	Module    string         `gorm:"size:64" json:"module"`                            // 从请求路径提取的业务模块。
	Action    string         `gorm:"size:64" json:"action"`                            // 当前记录为 HTTP 方法，保留给后续业务动作扩展。
	Method    string         `gorm:"size:16" json:"method"`                            // HTTP 方法。
	Path      string         `gorm:"size:255" json:"path"`                             // 请求路径。
	IP        string         `gorm:"size:64" json:"ip"`                                // 客户端 IP。
	Status    int            `json:"status"`                                           // HTTP 响应状态码。
}

// TableName 固定操作日志表名。
func (OperationLog) TableName() string { return "sys_operation_logs" }
