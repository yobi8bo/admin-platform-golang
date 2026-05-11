package audit

import (
	"net/http"
	"strconv"
	"strings"

	"admin-platform/backend/internal/pkg/errs"
	"admin-platform/backend/internal/pkg/response"
	"admin-platform/backend/internal/pkg/timex"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler 承载审计日志查询接口，日志写入由认证和中间件模块完成。
type Handler struct {
	db *gorm.DB
}

// NewHandler 创建审计 handler。
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// Register 注册审计查询接口，所有接口都需要对应审计查看权限。
func (h *Handler) Register(rg *gin.RouterGroup, require func(string) gin.HandlerFunc) {
	audit := rg.Group("/audit")
	audit.GET("/login-logs", require("audit:login-log:list"), h.LoginLogs)
	audit.GET("/operation-logs", require("audit:operation-log:list"), h.OperationLogs)
}

// LoginLogs 分页返回登录审计日志。
func (h *Handler) LoginLogs(c *gin.Context) {
	page, pageSize := pageParams(c)
	var total int64
	var list []LoginLog
	if err := h.db.Model(&LoginLog{}).Count(&total).Offset((page - 1) * pageSize).Limit(pageSize).Order("id desc").Find(&list).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, errs.CodeInternal, err.Error())
		return
	}
	response.OK(c, response.Page[LoginLog]{List: list, Total: total, Page: page, PageSize: pageSize})
}

// OperationLogs 分页返回操作审计日志，并补充面向前端展示的中文操作描述。
func (h *Handler) OperationLogs(c *gin.Context) {
	page, pageSize := pageParams(c)
	var total int64
	var list []operationLogItem
	if err := h.db.Model(&OperationLog{}).Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, errs.CodeInternal, err.Error())
		return
	}
	// 左连接用户表用于展示操作者名称；用户被删除时仍保留原始 user_id 和日志记录。
	query := h.db.Table("sys_operation_logs AS logs").
		Select("logs.id, logs.created_at, logs.user_id, COALESCE(users.nickname, users.username, '') AS username, logs.method, logs.path, logs.status").
		Joins("LEFT JOIN sys_users AS users ON users.id = logs.user_id")
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("logs.created_at desc").Scan(&list).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, errs.CodeInternal, err.Error())
		return
	}
	for i := range list {
		list[i].Operation = describeOperation(list[i].Method, list[i].Path)
	}
	response.OK(c, response.Page[operationLogItem]{List: list, Total: total, Page: page, PageSize: pageSize})
}

type operationLogItem struct {
	ID        uint           `json:"id"`
	CreatedAt timex.DateTime `json:"createdAt"`
	UserID    uint           `json:"userId"`
	Username  string         `json:"username"`
	Operation string         `json:"operation"`
	Method    string         `json:"-"`
	Path      string         `json:"-"`
	Status    int            `json:"status"`
}

func describeOperation(method, path string) string {
	// 优先使用精确路径标签，无法识别时再按 HTTP 方法和模块推导通用描述。
	if label := operationLabel(method, path); label != "" {
		return label
	}
	resource := operationResource(path)
	action := operationAction(method)
	if resource == "" {
		if action == "" {
			return "系统操作"
		}
		return action
	}
	if action == "" {
		return resource
	}
	return action + resource
}

func operationLabel(method, path string) string {
	cleaned := strings.TrimPrefix(path, "/api/")
	switch {
	case method == http.MethodPost && cleaned == "auth/logout":
		return "退出登录"
	case method == http.MethodGet && cleaned == "auth/profile":
		return "查看个人资料"
	case method == http.MethodPut && cleaned == "auth/profile":
		return "更新个人资料"
	case method == http.MethodPut && cleaned == "auth/password":
		return "修改密码"
	case method == http.MethodPost && cleaned == "files/avatar":
		return "上传头像"
	case method == http.MethodPost && cleaned == "files/upload":
		return "上传文件"
	default:
		return ""
	}
}

func operationAction(method string) string {
	switch method {
	case http.MethodGet:
		return "查看"
	case http.MethodPost:
		return "新增"
	case http.MethodPut:
		return "更新"
	case http.MethodDelete:
		return "删除"
	default:
		return "操作"
	}
}

func operationResource(path string) string {
	cleaned := strings.TrimPrefix(path, "/api/")
	parts := strings.Split(cleaned, "/")
	if len(parts) == 0 || parts[0] == "" {
		return ""
	}
	if len(parts) < 2 {
		return operationResourceByModule(parts[0])
	}
	switch parts[0] + "/" + parts[1] {
	case "auth/logout":
		return "退出登录"
	case "auth/profile":
		return "个人资料"
	case "auth/password":
		return "登录密码"
	case "files/upload":
		return "文件"
	case "files/avatar":
		return "头像"
	case "system/users":
		return "用户"
	case "system/roles":
		return "角色"
	case "system/menus":
		return "菜单"
	case "system/depts":
		return "部门"
	case "audit/login-logs":
		return "登录日志"
	case "audit/operation-logs":
		return "操作日志"
	default:
		return operationResourceByModule(parts[0])
	}
}

func operationResourceByModule(module string) string {
	switch module {
	case "auth":
		return "认证"
	case "files":
		return "文件"
	case "system":
		return "系统"
	case "audit":
		return "审计日志"
	default:
		return ""
	}
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
