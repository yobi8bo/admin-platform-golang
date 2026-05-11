package middleware

import (
	"net/http"
	"strings"
	"time"

	"admin-platform/backend/internal/config"
	"admin-platform/backend/internal/modules/audit"
	"admin-platform/backend/internal/modules/system"
	"admin-platform/backend/internal/pkg/contextx"
	"admin-platform/backend/internal/pkg/errs"
	jwtx "admin-platform/backend/internal/pkg/jwt"
	"admin-platform/backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Trace 为每个请求写入 traceId，响应和日志依赖该值串联一次请求链路。
func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("traceId", time.Now().Format("20060102150405.000000000"))
		c.Next()
	}
}

// RequestLogger 记录 HTTP 请求摘要，避免在业务 handler 中重复打访问日志。
func RequestLogger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Info("http request",
			zap.String("traceId", traceID(c)),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("ip", c.ClientIP()),
		)
	}
}

// Recovery 统一兜底 panic，防止内部异常直接暴露给前端。
func Recovery(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("panic recovered", zap.Any("panic", r), zap.String("traceId", traceID(c)))
				response.Fail(c, http.StatusInternalServerError, errs.CodeInternal, "internal server error")
				c.Abort()
			}
		}()
		c.Next()
	}
}

// Auth 校验 Bearer access token，并把用户 ID 和角色 ID 写入 Gin 上下文供后续权限判断使用。
func Auth(cfg config.JWTConfig, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			response.Fail(c, http.StatusUnauthorized, errs.CodeUnauthorized, "missing token")
			c.Abort()
			return
		}
		claims, err := jwtx.Parse(cfg.Secret, strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, errs.CodeUnauthorized, "invalid token")
			c.Abort()
			return
		}
		// 角色关系每次请求从数据库读取，确保角色调整后权限能立即生效。
		var roleIDs []uint
		if err := db.Model(&system.UserRole{}).Where("user_id = ?", claims.UserID).Pluck("role_id", &roleIDs).Error; err != nil {
			response.Fail(c, http.StatusUnauthorized, errs.CodeUnauthorized, "invalid user")
			c.Abort()
			return
		}
		c.Set(contextx.UserIDKey, claims.UserID)
		c.Set(contextx.RoleIDsKey, roleIDs)
		c.Next()
	}
}

// RequirePermission 校验当前用户角色是否拥有指定权限字符串，权限数据来自菜单表。
func RequirePermission(db *gorm.DB, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleIDs := contextx.RoleIDs(c)
		if len(roleIDs) == 0 {
			response.Fail(c, http.StatusForbidden, errs.CodeForbidden, "permission denied")
			c.Abort()
			return
		}
		var count int64
		err := db.Model(&system.Menu{}).
			Joins("JOIN sys_role_menus ON sys_role_menus.menu_id = sys_menus.id").
			Where("sys_role_menus.role_id IN ? AND sys_menus.permission = ?", roleIDs, permission).
			Where("sys_menus.deleted_at IS NULL").
			Count(&count).Error
		if err != nil || count == 0 {
			response.Fail(c, http.StatusForbidden, errs.CodeForbidden, "permission denied")
			c.Abort()
			return
		}
		c.Next()
	}
}

// OperationAudit 记录会改变系统状态的私有接口调用，审计失败不阻断主业务响应。
func OperationAudit(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		path := c.Request.URL.Path
		if shouldSkipOperationAudit(c.Request.Method, path) {
			return
		}

		// 审计日志是旁路记录，写入失败不能影响用户操作结果。
		_ = db.Create(&audit.OperationLog{
			UserID: contextx.UserID(c),
			Module: moduleName(path),
			Action: c.Request.Method,
			Method: c.Request.Method,
			Path:   path,
			IP:     c.ClientIP(),
			Status: c.Writer.Status(),
		}).Error
	}
}

func shouldSkipOperationAudit(method, path string) bool {
	if path == "/api/audit/operation-logs" {
		return true
	}
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return false
	default:
		return true
	}
}

func moduleName(path string) string {
	trimmed := strings.TrimPrefix(path, "/api/")
	if trimmed == path {
		return "unknown"
	}
	parts := strings.Split(trimmed, "/")
	if len(parts) == 0 || parts[0] == "" {
		return "unknown"
	}
	return parts[0]
}

func traceID(c *gin.Context) string {
	if v, ok := c.Get("traceId"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
