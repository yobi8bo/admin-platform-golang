package system

import (
	"admin-platform/backend/internal/pkg/timex"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt timex.DateTime `gorm:"type:timestamptz;autoCreateTime" json:"createdAt"`
	UpdatedAt timex.DateTime `gorm:"type:timestamptz;autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy uint           `json:"createdBy"`
	UpdatedBy uint           `json:"updatedBy"`
}

type User struct {
	BaseModel
	Username string `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Nickname string `gorm:"size:64;not null" json:"nickname"`
	Password string `gorm:"size:255;not null" json:"-"`
	Email    string `gorm:"size:128" json:"email"`
	Mobile   string `gorm:"size:32" json:"mobile"`
	AvatarID *uint  `json:"avatarId"`
	Status   string `gorm:"size:16;not null;default:enabled" json:"status"`
	DeptID   *uint  `json:"deptId"`
	Roles    []Role `gorm:"many2many:sys_user_roles;" json:"roles,omitempty"`
}

func (User) TableName() string { return "sys_users" }

type Role struct {
	BaseModel
	Code      string `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Name      string `gorm:"size:64;not null" json:"name"`
	Sort      int    `gorm:"default:0" json:"sort"`
	Status    string `gorm:"size:16;not null;default:enabled" json:"status"`
	DataScope string `gorm:"size:32;not null;default:self" json:"dataScope"`
	Menus     []Menu `gorm:"many2many:sys_role_menus;" json:"menus,omitempty"`
}

func (Role) TableName() string { return "sys_roles" }

type Menu struct {
	BaseModel
	ParentID   uint   `gorm:"index;default:0" json:"parentId"`
	Name       string `gorm:"size:64;not null" json:"name"`
	Title      string `gorm:"size:64;not null" json:"title"`
	Type       string `gorm:"size:16;not null" json:"type"`
	Path       string `gorm:"size:255" json:"path"`
	Component  string `gorm:"size:255" json:"component"`
	Icon       string `gorm:"size:64" json:"icon"`
	Permission string `gorm:"size:128;index" json:"permission"`
	Hidden     bool   `gorm:"default:false" json:"hidden"`
	Sort       int    `gorm:"default:0" json:"sort"`
	Children   []Menu `gorm:"-" json:"children,omitempty"`
}

func (Menu) TableName() string { return "sys_menus" }

type Dept struct {
	BaseModel
	ParentID uint   `gorm:"index;default:0" json:"parentId"`
	Name     string `gorm:"size:64;not null" json:"name"`
	Sort     int    `gorm:"default:0" json:"sort"`
	Status   string `gorm:"size:16;not null;default:enabled" json:"status"`
	Children []Dept `gorm:"-" json:"children,omitempty"`
}

func (Dept) TableName() string { return "sys_depts" }

type UserRole struct {
	UserID uint `gorm:"primaryKey"`
	RoleID uint `gorm:"primaryKey"`
}

func (UserRole) TableName() string { return "sys_user_roles" }

type RoleMenu struct {
	RoleID uint `gorm:"primaryKey"`
	MenuID uint `gorm:"primaryKey"`
}

func (RoleMenu) TableName() string { return "sys_role_menus" }
