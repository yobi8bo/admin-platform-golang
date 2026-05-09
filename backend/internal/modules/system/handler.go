package system

import (
	"net/http"
	"strconv"

	"admin-platform/backend/internal/pkg/contextx"
	"admin-platform/backend/internal/pkg/crypto"
	"admin-platform/backend/internal/pkg/errs"
	"admin-platform/backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

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
	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		if len(req.RoleIDs) > 0 {
			var roles []Role
			if err := tx.Find(&roles, req.RoleIDs).Error; err != nil {
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
	updates := map[string]any{"nickname": req.Nickname, "email": req.Email, "mobile": req.Mobile, "status": req.Status, "updated_by": contextx.UserID(c)}
	if err := h.db.Model(&user).Updates(updates).Error; err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}
	if req.RoleIDs != nil {
		var roles []Role
		if err := h.db.Find(&roles, req.RoleIDs).Error; err == nil {
			_ = h.db.Model(&user).Association("Roles").Replace(roles)
		}
	}
	response.OK(c, user)
}

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

func (h *Handler) MyMenuTree(c *gin.Context) {
	menus, err := h.userMenus(contextx.UserID(c), []string{"catalog", "menu"})
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, errs.CodeInternal, err.Error())
		return
	}
	response.OK(c, buildMenuTree(menus, 0))
}

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
