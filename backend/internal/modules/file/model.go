package file

import (
	"admin-platform/backend/internal/pkg/timex"

	"gorm.io/gorm"
)

type File struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CreatedAt    timex.DateTime `gorm:"type:timestamptz;autoCreateTime" json:"createdAt"`
	UpdatedAt    timex.DateTime `gorm:"type:timestamptz;autoUpdateTime" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	OriginalName string         `gorm:"size:255;not null" json:"originalName"`
	Bucket       string         `gorm:"size:128;not null" json:"bucket"`
	ObjectKey    string         `gorm:"size:512;not null" json:"objectKey"`
	ContentType  string         `gorm:"size:128" json:"contentType"`
	Size         int64          `json:"size"`
	CreatedBy    uint           `json:"createdBy"`
}

func (File) TableName() string { return "sys_files" }
