package system

import (
	"admin-platform/backend/internal/pkg/timex"

	"gorm.io/gorm"
)

// BaseModel 是系统业务表的通用审计字段，软删除字段不对前端暴露。
type BaseModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`                             // 系统内通用自增主键。
	CreatedAt timex.DateTime `gorm:"type:timestamptz;autoCreateTime" json:"createdAt"` // 记录创建时间，按统一 DateTime 格式输出。
	UpdatedAt timex.DateTime `gorm:"type:timestamptz;autoUpdateTime" json:"updatedAt"` // 记录最近更新时间，供列表排序和审计展示使用。
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                                   // GORM 软删除标记，查询默认排除已删除数据。
	CreatedBy uint           `json:"createdBy"`                                        // 创建人用户 ID，来自认证上下文。
	UpdatedBy uint           `json:"updatedBy"`                                        // 最近更新人用户 ID，来自认证上下文。
}

// User 表示后台账号，密码哈希只用于服务端校验，禁止序列化给前端。
type User struct {
	BaseModel
	Username string `gorm:"size:64;uniqueIndex;not null" json:"username"`     // 登录账号，系统内唯一。
	Nickname string `gorm:"size:64;not null" json:"nickname"`                 // 展示名称，用于审计日志和页面展示。
	Password string `gorm:"size:255;not null" json:"-"`                       // bcrypt 密码哈希，永远不返回给前端。
	Email    string `gorm:"size:128" json:"email"`                            // 联系邮箱，当前不作为登录凭据。
	Mobile   string `gorm:"size:32" json:"mobile"`                            // 联系手机号，当前不作为登录凭据。
	AvatarID *uint  `json:"avatarId"`                                         // 关联 sys_files 的头像文件 ID，允许为空。
	Status   string `gorm:"size:16;not null;default:enabled" json:"status"`   // 账号状态，enabled 表示允许登录。
	DeptID   *uint  `json:"deptId"`                                           // 所属部门 ID，允许未分配部门。
	Roles    []Role `gorm:"many2many:sys_user_roles;" json:"roles,omitempty"` // 通过 sys_user_roles 维护的角色关系。
}

// TableName 固定用户表名，必须与 migrations 中的 sys_users 保持一致。
func (User) TableName() string { return "sys_users" }

// Role 表示权限角色，菜单关系决定该角色可访问的接口和前端菜单。
type Role struct {
	BaseModel
	Code      string `gorm:"size:64;uniqueIndex;not null" json:"code"`         // 角色编码，系统内唯一，适合做稳定标识。
	Name      string `gorm:"size:64;not null" json:"name"`                     // 角色展示名称。
	Sort      int    `gorm:"default:0" json:"sort"`                            // 列表排序，数值越小越靠前。
	Status    string `gorm:"size:16;not null;default:enabled" json:"status"`   // 角色状态，禁用后不可再分配给用户。
	DataScope string `gorm:"size:32;not null;default:self" json:"dataScope"`   // 数据权限范围，当前默认仅本人数据。
	Menus     []Menu `gorm:"many2many:sys_role_menus;" json:"menus,omitempty"` // 通过 sys_role_menus 维护的菜单和接口权限。
}

// TableName 固定角色表名，必须与 migrations 中的 sys_roles 保持一致。
func (Role) TableName() string { return "sys_roles" }

// Menu 同时承载前端菜单节点和后端权限字符串。
type Menu struct {
	BaseModel
	ParentID   uint   `gorm:"index;default:0" json:"parentId"`  // 父级菜单 ID，0 表示根节点。
	Name       string `gorm:"size:64;not null" json:"name"`     // 前端路由名称或权限节点名称。
	Title      string `gorm:"size:64;not null" json:"title"`    // 前端展示标题。
	Type       string `gorm:"size:16;not null" json:"type"`     // 节点类型：catalog、menu、button 或 api。
	Path       string `gorm:"size:255" json:"path"`             // 前端路由路径或接口分组路径。
	Component  string `gorm:"size:255" json:"component"`        // 前端组件路径，接口权限节点可为空。
	Icon       string `gorm:"size:64" json:"icon"`              // 前端菜单图标标识。
	Permission string `gorm:"size:128;index" json:"permission"` // 后端权限字符串，RequirePermission 依赖该字段。
	Hidden     bool   `gorm:"default:false" json:"hidden"`      // 是否在前端菜单中隐藏。
	Sort       int    `gorm:"default:0" json:"sort"`            // 同级排序，数值越小越靠前。
	Children   []Menu `gorm:"-" json:"children,omitempty"`      // 接口返回树形结构时临时填充，不入库。
}

// TableName 固定菜单表名，必须与 migrations 中的 sys_menus 保持一致。
func (Menu) TableName() string { return "sys_menus" }

// Dept 表示组织部门，用于用户归属和后续数据权限扩展。
type Dept struct {
	BaseModel
	ParentID uint   `gorm:"index;default:0" json:"parentId"`                // 父级部门 ID，0 表示根部门。
	Name     string `gorm:"size:64;not null" json:"name"`                   // 部门名称。
	Sort     int    `gorm:"default:0" json:"sort"`                          // 同级排序，数值越小越靠前。
	Status   string `gorm:"size:16;not null;default:enabled" json:"status"` // 部门状态，enabled 表示可用。
	Children []Dept `gorm:"-" json:"children,omitempty"`                    // 接口返回树形结构时临时填充，不入库。
}

// TableName 固定部门表名，必须与 migrations 中的 sys_depts 保持一致。
func (Dept) TableName() string { return "sys_depts" }

// UserRole 是用户与角色的多对多关联表。
type UserRole struct {
	UserID uint `gorm:"primaryKey"` // 用户 ID。
	RoleID uint `gorm:"primaryKey"` // 角色 ID。
}

// TableName 固定用户角色关联表名。
func (UserRole) TableName() string { return "sys_user_roles" }

// RoleMenu 是角色与菜单权限的多对多关联表。
type RoleMenu struct {
	RoleID uint `gorm:"primaryKey"` // 角色 ID。
	MenuID uint `gorm:"primaryKey"` // 菜单或权限节点 ID。
}

// TableName 固定角色菜单关联表名。
func (RoleMenu) TableName() string { return "sys_role_menus" }
