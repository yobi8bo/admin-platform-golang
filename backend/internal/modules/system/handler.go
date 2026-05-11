package system

import (
	"fmt"
	"net/http"
	"strconv"

	"admin-platform/backend/internal/pkg/contextx"
	"admin-platform/backend/internal/pkg/crypto"
	"admin-platform/backend/internal/pkg/errs"
	"admin-platform/backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler 承载系统管理模块的 HTTP 入口，包括用户、角色、菜单和部门管理。
type Handler struct {
	db *gorm.DB
}

// NewHandler 创建系统管理 handler，数据库依赖由 bootstrap 注入。
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// Register 注册系统管理私有接口，调用方必须传入权限中间件以保护敏感操作。
func (h *Handler) Register(rg *gin.RouterGroup, require func(string) gin.HandlerFunc) {
	system := rg.Group("/system")
	system.GET("/users", require("system:user:list"), h.ListUsers)
	system.POST("/users", require("system:user:create"), h.CreateUser)
	system.PUT("/users/:id", require("system:user:update"), h.UpdateUser)
	system.DELETE("/users/:id", require("system:user:delete"), h.DeleteUser)

	system.GET("/roles", require("system:role:list"), h.ListRoles)
	system.POST("/roles", require("system:role:create"), h.CreateRole)
	system.PUT("/roles/:id", require("system:role:update"), h.UpdateRole)
	system.DELETE("/roles/:id", require("system:role:delete"), h.DeleteRole)

	system.GET("/menus/tree", require("system:menu:list"), h.MenuTree)
	system.POST("/menus", require("system:menu:list"), h.CreateMenu)
	system.PUT("/menus/:id", require("system:menu:list"), h.UpdateMenu)
	system.DELETE("/menus/:id", require("system:menu:list"), h.DeleteMenu)
	system.GET("/menus/my-tree", h.MyMenuTree)
	system.GET("/depts/tree", require("system:dept:list"), h.DeptTree)
	system.POST("/depts", require("system:dept:list"), h.CreateDept)
	system.PUT("/depts/:id", require("system:dept:list"), h.UpdateDept)
	system.DELETE("/depts/:id", require("system:dept:list"), h.DeleteDept)
}

// ListUsers 分页查询后台账号，支持按用户名或昵称模糊搜索。
func (h *Handler) ListUsers(c *gin.Context) {
	page, pageSize := pageParams(c)
	var total int64
	var users []User
	query := h.db.Model(&User{}).Preload("Roles")
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("username ILIKE ? OR nickname ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if err := query.Count(&total).Offset((page - 1) * pageSize).Limit(pageSize).Order("id desc").Find(&users).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, errs.CodeInternal, err.Error())
		return
	}
	response.OK(c, response.Page[User]{List: users, Total: total, Page: page, PageSize: pageSize})
}

type createUserReq struct {
	Username string `json:"username" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	RoleIDs  []uint `json:"roleIds"`
}

// CreateUser 创建后台账号，并在同一事务内写入用户和角色关系。
func (h *Handler) CreateUser(c *gin.Context) {
	var req createUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}
	hash, err := crypto.HashPassword(req.Password)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, errs.CodeInternal, err.Error())
		return
	}
	user := User{
		BaseModel: BaseModel{CreatedBy: contextx.UserID(c)},
		Username:  req.Username,
		Nickname:  req.Nickname,
		Password:  hash,
		Email:     req.Email,
		Mobile:    req.Mobile,
	}
	// 用户主表和角色关联必须同时成功，避免出现账号已创建但无角色的半成品数据。
	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		if len(req.RoleIDs) > 0 {
			roles, err := loadAssignableRoles(tx, req.RoleIDs)
			if err != nil {
				return err
			}
			return tx.Model(&user).Association("Roles").Replace(roles)
		}
		return nil
	})
	if err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}
	response.Created(c, user)
}

type updateUserReq struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Status   string `json:"status"`
	RoleIDs  []uint `json:"roleIds"`
}

// UpdateUser 更新账号基础信息；传入 roleIds 时同步替换角色关系。
func (h *Handler) UpdateUser(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req updateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}
	var user User
	if err := h.db.First(&user, id).Error; err != nil {
		response.Fail(c, http.StatusNotFound, errs.CodeNotFound, "user not found")
		return
	}
	// 显式列出允许更新的字段，避免请求体扩展后意外覆盖敏感字段。
	updates := map[string]any{"nickname": req.Nickname, "email": req.Email, "mobile": req.Mobile, "status": req.Status, "updated_by": contextx.UserID(c)}
	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&user).Updates(updates).Error; err != nil {
			return err
		}
		// roleIds 未传表示不调整角色；空数组表示清空全部角色。
		if req.RoleIDs == nil {
			return nil
		}
		roles, err := loadAssignableRoles(tx, req.RoleIDs)
		if err != nil {
			return err
		}
		return tx.Model(&user).Association("Roles").Replace(roles)
	})
	if err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}
	response.OK(c, user)
}

// DeleteUser 删除后台账号，当前登录用户不能删除自己以避免会话失去主体。
func (h *Handler) DeleteUser(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if id == contextx.UserID(c) {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, "cannot delete current user")
		return
	}
	if err := h.db.Delete(&User{}, id).Error; err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"deleted": true})
}

// ListRoles 查询角色列表，并预加载菜单以便前端展示角色授权状态。
func (h *Handler) ListRoles(c *gin.Context) {
	var roles []Role
	if err := h.db.Preload("Menus").Order("sort asc,id asc").Find(&roles).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, errs.CodeInternal, err.Error())
		return
	}
	response.OK(c, roles)
}

type roleReq struct {
	Code      string `json:"code" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Sort      int    `json:"sort"`
	Status    string `json:"status"`
	DataScope string `json:"dataScope"`
	MenuIDs   []uint `json:"menuIds"`
}

// CreateRole 创建角色并绑定菜单权限。
func (h *Handler) CreateRole(c *gin.Context) {
	var req roleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}
	role := Role{Code: req.Code, Name: req.Name, Sort: req.Sort, Status: valueOr(req.Status, "enabled"), DataScope: valueOr(req.DataScope, "self")}
	if err := h.db.Create(&role).Error; err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}
	h.replaceRoleMenus(role.ID, req.MenuIDs)
	response.Created(c, role)
}

// UpdateRole 更新角色基础信息并替换菜单权限。
func (h *Handler) UpdateRole(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req roleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}
	if err := h.db.Model(&Role{}).Where("id = ?", id).Updates(map[string]any{
		"code": req.Code, "name": req.Name, "sort": req.Sort, "status": req.Status, "data_scope": req.DataScope,
	}).Error; err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}
	h.replaceRoleMenus(id, req.MenuIDs)
	response.OK(c, gin.H{"updated": true})
}

// DeleteRole 删除角色；调用方应确保删除前已评估用户授权影响。
func (h *Handler) DeleteRole(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.db.Delete(&Role{}, id).Error; err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"deleted": true})
}

// MenuTree 返回完整菜单树，用于系统菜单维护和角色授权。
func (h *Handler) MenuTree(c *gin.Context) {
	var menus []Menu
	if err := h.db.Order("sort asc,id asc").Find(&menus).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, errs.CodeInternal, err.Error())
		return
	}
	response.OK(c, buildMenuTree(menus, 0))
}

type menuReq struct {
	ParentID   uint   `json:"parentId"`
	Name       string `json:"name" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Type       string `json:"type" binding:"required"`
	Path       string `json:"path"`
	Component  string `json:"component"`
	Icon       string `json:"icon"`
	Permission string `json:"permission"`
	Hidden     bool   `json:"hidden"`
	Sort       int    `json:"sort"`
}

// CreateMenu 创建菜单或权限节点，permission 字段会参与后端接口鉴权。
func (h *Handler) CreateMenu(c *gin.Context) {
	var req menuReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}
	menu := Menu{
		BaseModel:  BaseModel{CreatedBy: contextx.UserID(c)},
		ParentID:   req.ParentID,
		Name:       req.Name,
		Title:      req.Title,
		Type:       req.Type,
		Path:       req.Path,
		Component:  req.Component,
		Icon:       req.Icon,
		Permission: req.Permission,
		Hidden:     req.Hidden,
		Sort:       req.Sort,
	}
	if err := h.db.Create(&menu).Error; err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}
	response.Created(c, menu)
}

// UpdateMenu 更新菜单或权限节点，必须保持 permission 与路由注册中的权限字符串一致。
func (h *Handler) UpdateMenu(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req menuReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}
	if err := h.db.Model(&Menu{}).Where("id = ?", id).Updates(map[string]any{
		"parent_id": req.ParentID, "name": req.Name, "title": req.Title, "type": req.Type,
		"path": req.Path, "component": req.Component, "icon": req.Icon, "permission": req.Permission,
		"hidden": req.Hidden, "sort": req.Sort, "updated_by": contextx.UserID(c),
	}).Error; err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"updated": true})
}

// DeleteMenu 删除菜单节点；存在子节点时拒绝删除以保护树结构完整性。
func (h *Handler) DeleteMenu(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var children int64
	if err := h.db.Model(&Menu{}).Where("parent_id = ?", id).Count(&children).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, errs.CodeInternal, err.Error())
		return
	}
	if children > 0 {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, "please delete child menus first")
		return
	}
	if err := h.db.Delete(&Menu{}, id).Error; err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"deleted": true})
}

// MyMenuTree 返回当前用户可见的目录和菜单，不包含 button/api 权限节点。
func (h *Handler) MyMenuTree(c *gin.Context) {
	menus, err := h.userMenus(contextx.UserID(c), []string{"catalog", "menu"})
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, errs.CodeInternal, err.Error())
		return
	}
	response.OK(c, buildMenuTree(menus, 0))
}

// DeptTree 返回部门树，用于组织架构维护和用户归属选择。
func (h *Handler) DeptTree(c *gin.Context) {
	var depts []Dept
	if err := h.db.Order("sort asc,id asc").Find(&depts).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, errs.CodeInternal, err.Error())
		return
	}
	response.OK(c, buildDeptTree(depts, 0))
}

type deptReq struct {
	ParentID uint   `json:"parentId"`
	Name     string `json:"name" binding:"required"`
	Sort     int    `json:"sort"`
	Status   string `json:"status"`
}

// CreateDept 创建部门节点，默认启用以便新部门可立即被选择。
func (h *Handler) CreateDept(c *gin.Context) {
	var req deptReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}
	dept := Dept{
		BaseModel: BaseModel{CreatedBy: contextx.UserID(c)},
		ParentID:  req.ParentID,
		Name:      req.Name,
		Sort:      req.Sort,
		Status:    valueOr(req.Status, "enabled"),
	}
	if err := h.db.Create(&dept).Error; err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}
	response.Created(c, dept)
}

// UpdateDept 更新部门节点，显式限制可更新字段以保护审计字段。
func (h *Handler) UpdateDept(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req deptReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}
	if err := h.db.Model(&Dept{}).Where("id = ?", id).Updates(map[string]any{
		"parent_id": req.ParentID, "name": req.Name, "sort": req.Sort,
		"status": valueOr(req.Status, "enabled"), "updated_by": contextx.UserID(c),
	}).Error; err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"updated": true})
}

// DeleteDept 删除部门节点；存在子部门时拒绝删除以保护组织树完整性。
func (h *Handler) DeleteDept(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var children int64
	if err := h.db.Model(&Dept{}).Where("parent_id = ?", id).Count(&children).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, errs.CodeInternal, err.Error())
		return
	}
	if children > 0 {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, "please delete child departments first")
		return
	}
	if err := h.db.Delete(&Dept{}, id).Error; err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"deleted": true})
}

// UserPermissions 查询用户拥有的 button/api 权限字符串，供前端按钮控制和后端鉴权复用。
func (h *Handler) UserPermissions(userID uint) ([]string, error) {
	menus, err := h.userMenus(userID, []string{"button", "api"})
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(menus))
	for _, menu := range menus {
		if menu.Permission != "" {
			out = append(out, menu.Permission)
		}
	}
	return out, nil
}

func (h *Handler) userMenus(userID uint, types []string) ([]Menu, error) {
	var menus []Menu
	// DISTINCT 用于消除用户多个角色命中同一菜单时产生的重复节点。
	err := h.db.Table("sys_menus").
		Select("DISTINCT sys_menus.*").
		Joins("JOIN sys_role_menus ON sys_role_menus.menu_id = sys_menus.id").
		Joins("JOIN sys_user_roles ON sys_user_roles.role_id = sys_role_menus.role_id").
		Where("sys_user_roles.user_id = ? AND sys_menus.type IN ?", userID, types).
		Where("sys_menus.deleted_at IS NULL").
		Order("sys_menus.sort asc, sys_menus.id asc").
		Scan(&menus).Error
	return menus, err
}

func (h *Handler) replaceRoleMenus(roleID uint, menuIDs []uint) {
	// 角色授权以请求中的 menuIds 为准，先清空后重建可以表达“取消全部授权”。
	_ = h.db.Where("role_id = ?", roleID).Delete(&RoleMenu{}).Error
	for _, menuID := range menuIDs {
		_ = h.db.Create(&RoleMenu{RoleID: roleID, MenuID: menuID}).Error
	}
}

func buildMenuTree(all []Menu, parentID uint) []Menu {
	out := make([]Menu, 0)
	for _, item := range all {
		if item.ParentID == parentID {
			item.Children = buildMenuTree(all, item.ID)
			out = append(out, item)
		}
	}
	return out
}

func buildDeptTree(all []Dept, parentID uint) []Dept {
	out := make([]Dept, 0)
	for _, item := range all {
		if item.ParentID == parentID {
			item.Children = buildDeptTree(all, item.ID)
			out = append(out, item)
		}
	}
	return out
}

func pageParams(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		// 限制单页最大数量，避免管理端列表接口被大 pageSize 拖慢。
		pageSize = 20
	}
	return page, pageSize
}

func parseID(c *gin.Context) (uint, bool) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, "invalid id")
		return 0, false
	}
	return uint(id64), true
}

func valueOr(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func loadAssignableRoles(db *gorm.DB, roleIDs []uint) ([]Role, error) {
	ids := uniqueUintIDs(roleIDs)
	if len(ids) == 0 {
		return []Role{}, nil
	}

	var roles []Role
	if err := db.Where("id IN ? AND status = ?", ids, "enabled").Find(&roles).Error; err != nil {
		return nil, err
	}
	if len(roles) != len(ids) {
		// 禁止把不存在或禁用角色写入用户关系，避免绕过角色状态约束。
		return nil, fmt.Errorf("角色不存在或已禁用")
	}
	return roles, nil
}

func uniqueUintIDs(ids []uint) []uint {
	out := make([]uint, 0, len(ids))
	seen := make(map[uint]struct{}, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
