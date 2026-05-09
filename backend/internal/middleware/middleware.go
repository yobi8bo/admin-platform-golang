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

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("traceId", time.Now().Format("20060102150405.000000000"))
		c.Next()
	}
}

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

func OperationAudit(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		path := c.Request.URL.Path
		if shouldSkipOperationAudit(c.Request.Method, path) {
			return
		}

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
