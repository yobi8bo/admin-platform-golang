package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"

	"admin-platform/backend/internal/config"
	"admin-platform/backend/internal/modules/audit"
	"admin-platform/backend/internal/modules/system"
	"admin-platform/backend/internal/pkg/contextx"
	"admin-platform/backend/internal/pkg/crypto"
	"admin-platform/backend/internal/pkg/errs"
	jwtx "admin-platform/backend/internal/pkg/jwt"
	"admin-platform/backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Handler 承载认证与个人中心接口，刷新令牌状态通过 Redis 管理。
type Handler struct {
	db     *gorm.DB
	redis  *redis.Client
	cfg    config.JWTConfig
	system *system.Handler
}

// NewHandler 创建认证 handler，systemHandler 用于复用用户权限查询逻辑。
func NewHandler(db *gorm.DB, redis *redis.Client, cfg config.JWTConfig, systemHandler *system.Handler) *Handler {
	return &Handler{db: db, redis: redis, cfg: cfg, system: systemHandler}
}

// RegisterPublic 注册无需登录即可访问的认证接口。
func (h *Handler) RegisterPublic(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	auth.POST("/login", h.Login)
	auth.POST("/refresh", h.Refresh)
}

// RegisterPrivate 注册需要登录访问的个人认证接口。
func (h *Handler) RegisterPrivate(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	auth.POST("/logout", h.Logout)
	auth.GET("/profile", h.Profile)
	auth.PUT("/profile", h.UpdateProfile)
	auth.PUT("/password", h.UpdatePassword)
	auth.GET("/permissions", h.Permissions)
}

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 校验账号密码并签发 access token 和 refresh token。
func (h *Handler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.writeLoginLog(c, "", "failed", "无效请求")
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}

	var user system.User
	if err := h.db.Preload("Roles").Where("username = ?", req.Username).First(&user).Error; err != nil {
		h.writeLoginLog(c, req.Username, "failed", "用户名或密码错误")
		response.Fail(c, http.StatusUnauthorized, errs.CodeUnauthorized, "用户名或密码错误")
		return
	}
	// 禁用账号和密码错误统一返回相同信息，避免泄露账号状态。
	if user.Status != "enabled" || !crypto.CheckPassword(user.Password, req.Password) {
		h.writeLoginLog(c, req.Username, "failed", "用户名或密码错误")
		response.Fail(c, http.StatusUnauthorized, errs.CodeUnauthorized, "用户名或密码错误")
		return
	}

	access, refresh, err := h.issueTokens(c.Request.Context(), user.ID)
	if err != nil {
		h.writeLoginLog(c, req.Username, "failed", "令牌签发失败")
		response.Fail(c, http.StatusInternalServerError, errs.CodeInternal, err.Error())
		return
	}

	h.writeLoginLog(c, req.Username, "success", "登录成功")
	response.OK(c, gin.H{
		"accessToken":  access,
		"refreshToken": refresh,
		"expiresIn":    int(h.cfg.AccessTTL().Seconds()),
	})
}

func (h *Handler) writeLoginLog(c *gin.Context, username, status, message string) {
	// 登录日志是审计辅助信息，写入失败不能影响登录接口的主要结果。
	_ = h.db.Create(&audit.LoginLog{
		Username: username,
		IP:       c.ClientIP(),
		Status:   status,
		Message:  message,
	}).Error
}

type refreshReq struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// Refresh 轮换刷新令牌，旧 refresh token 成功使用后立即失效以降低重放风险。
func (h *Handler) Refresh(c *gin.Context) {
	var req refreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}
	claims, err := jwtx.Parse(h.cfg.Secret, req.RefreshToken)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, errs.CodeUnauthorized, "invalid refresh token")
		return
	}
	key := refreshKey(claims.UserID, claims.TokenID)
	if exists := h.redis.Exists(c.Request.Context(), key).Val(); exists == 0 {
		response.Fail(c, http.StatusUnauthorized, errs.CodeUnauthorized, "refresh token expired")
		return
	}
	// 刷新令牌采用一次性使用语义，删除旧 key 后再签发新的一组令牌。
	_ = h.redis.Del(c.Request.Context(), key).Err()
	access, refresh, err := h.issueTokens(c.Request.Context(), claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, errs.CodeInternal, err.Error())
		return
	}
	response.OK(c, gin.H{"accessToken": access, "refreshToken": refresh, "expiresIn": int(h.cfg.AccessTTL().Seconds())})
}

// Logout 删除当前用户的全部刷新令牌，使该账号在本服务端失去续期能力。
func (h *Handler) Logout(c *gin.Context) {
	userID := contextx.UserID(c)
	pattern := refreshKey(userID, "*")
	iter := h.redis.Scan(c.Request.Context(), 0, pattern, 100).Iterator()
	for iter.Next(c.Request.Context()) {
		_ = h.redis.Del(c.Request.Context(), iter.Val()).Err()
	}
	response.OK(c, gin.H{"logout": true})
}

// Profile 返回当前登录用户的基础资料，密码哈希等敏感字段不会出现在响应中。
func (h *Handler) Profile(c *gin.Context) {
	var user system.User
	if err := h.db.Preload("Roles").First(&user, contextx.UserID(c)).Error; err != nil {
		response.Fail(c, http.StatusNotFound, errs.CodeNotFound, "user not found")
		return
	}
	response.OK(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
		"email":    user.Email,
		"mobile":   user.Mobile,
		"avatarId": user.AvatarID,
		"roles":    user.Roles,
	})
}

type updateProfileReq struct {
	Nickname string `json:"nickname" binding:"required"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	AvatarID *uint  `json:"avatarId"`
}

// UpdateProfile 只允许用户修改自己的展示资料和头像关联。
func (h *Handler) UpdateProfile(c *gin.Context) {
	var req updateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}

	var user system.User
	if err := h.db.First(&user, contextx.UserID(c)).Error; err != nil {
		response.Fail(c, http.StatusNotFound, errs.CodeNotFound, "user not found")
		return
	}

	updates := map[string]any{
		"nickname":  req.Nickname,
		"email":     req.Email,
		"mobile":    req.Mobile,
		"avatar_id": req.AvatarID,
	}
	if err := h.db.Model(&user).Updates(updates).Error; err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}

	response.OK(c, gin.H{"updated": true})
}

type updatePasswordReq struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

// UpdatePassword 校验旧密码后更新当前用户密码。
func (h *Handler) UpdatePassword(c *gin.Context) {
	var req updatePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}
	if len(req.NewPassword) < 6 {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, "新密码至少需要 6 位")
		return
	}

	var user system.User
	if err := h.db.First(&user, contextx.UserID(c)).Error; err != nil {
		response.Fail(c, http.StatusNotFound, errs.CodeNotFound, "user not found")
		return
	}
	if !crypto.CheckPassword(user.Password, req.OldPassword) {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, "原密码不正确")
		return
	}

	hash, err := crypto.HashPassword(req.NewPassword)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, errs.CodeInternal, err.Error())
		return
	}
	if err := h.db.Model(&user).Update("password", hash).Error; err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"updated": true})
}

// Permissions 返回当前用户拥有的权限字符串列表。
func (h *Handler) Permissions(c *gin.Context) {
	perms, err := h.system.UserPermissions(contextx.UserID(c))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, errs.CodeInternal, err.Error())
		return
	}
	response.OK(c, perms)
}

func (h *Handler) issueTokens(ctx context.Context, userID uint) (string, string, error) {
	accessID := randomID()
	refreshID := randomID()
	access, err := jwtx.Sign(h.cfg.Secret, userID, accessID, h.cfg.AccessTTL())
	if err != nil {
		return "", "", err
	}
	refresh, err := jwtx.Sign(h.cfg.Secret, userID, refreshID, h.cfg.RefreshTTL())
	if err != nil {
		return "", "", err
	}
	// refresh token 必须在 Redis 中存在才可续期，便于登出和轮换时主动失效。
	if err := h.redis.Set(ctx, refreshKey(userID, refreshID), "1", h.cfg.RefreshTTL()).Err(); err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

func randomID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// 随机源异常时使用纳秒时间兜底，保证令牌 ID 仍能生成但不作为强随机保证。
		return strings.ReplaceAll(time.Now().Format(time.RFC3339Nano), ":", "")
	}
	return hex.EncodeToString(b)
}

func refreshKey(userID uint, tokenID string) string {
	return "auth:refresh:" + strconvUint(userID) + ":" + tokenID
}

func strconvUint(v uint) string {
	return strconv.FormatUint(uint64(v), 10)
}
