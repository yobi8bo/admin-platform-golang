package file

import (
	"admin-platform/backend/internal/pkg/timex"

	"gorm.io/gorm"
)

// File 记录对象存储中文件的元数据，文件内容本身保存在 RustFS。
type File struct {
	ID           uint           `gorm:"primaryKey" json:"id"`                             // 文件元数据主键。
	CreatedAt    timex.DateTime `gorm:"type:timestamptz;autoCreateTime" json:"createdAt"` // 上传记录创建时间。
	UpdatedAt    timex.DateTime `gorm:"type:timestamptz;autoUpdateTime" json:"updatedAt"` // 元数据最近更新时间。
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`                                   // 软删除标记，不对前端暴露。
	OriginalName string         `gorm:"size:255;not null" json:"originalName"`            // 用户上传时的原始文件名，仅用于展示和下载文件名。
	Bucket       string         `gorm:"size:128;not null" json:"bucket"`                  // 对象存储桶名称。
	ObjectKey    string         `gorm:"size:512;not null" json:"objectKey"`               // 对象存储键，下载和删除都依赖该值。
	ContentType  string         `gorm:"size:128" json:"contentType"`                      // 上传时声明的 MIME 类型，头像上传会校验 image/*。
	Size         int64          `json:"size"`                                             // 对象大小，单位字节。
	CreatedBy    uint           `json:"createdBy"`                                        // 上传用户 ID，来自认证上下文。
}

// TableName 固定文件元数据表名，必须与 migrations 中的 sys_files 保持一致。
func (File) TableName() string { return "sys_files" }
