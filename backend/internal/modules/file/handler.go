package file

import (
	"context"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"admin-platform/backend/internal/config"
	"admin-platform/backend/internal/pkg/contextx"
	"admin-platform/backend/internal/pkg/errs"
	"admin-platform/backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

// Handler 承载文件元数据查询、上传、预签名 URL 和删除接口。
type Handler struct {
	db     *gorm.DB
	client *minio.Client
	cfg    config.RustFSConfig
}

// NewHandler 创建文件 handler，对象存储客户端由 bootstrap 统一初始化。
func NewHandler(db *gorm.DB, client *minio.Client, cfg config.RustFSConfig) *Handler {
	return &Handler{db: db, client: client, cfg: cfg}
}

// Register 注册文件私有接口；头像相关接口依赖登录态但不额外要求文件管理权限。
func (h *Handler) Register(rg *gin.RouterGroup, require func(string) gin.HandlerFunc) {
	files := rg.Group("/files")
	files.GET("", require("file:read"), h.List)
	files.POST("/upload", require("file:upload"), h.Upload)
	files.POST("/avatar", h.UploadAvatar)
	files.GET("/avatar/:id/url", h.URL)
	files.GET("/:id/download-url", require("file:read"), h.DownloadURL)
	files.GET("/:id/url", require("file:read"), h.URL)
	files.DELETE("/:id", require("file:delete"), h.Delete)
}

// List 分页查询文件元数据，支持按原始文件名和 MIME 类型搜索。
func (h *Handler) List(c *gin.Context) {
	page, pageSize := pageParams(c)
	var total int64
	var list []File
	query := h.db.Model(&File{})
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("original_name ILIKE ? OR content_type ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if err := query.Count(&total).Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at desc").Find(&list).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, errs.CodeInternal, err.Error())
		return
	}
	response.OK(c, response.Page[File]{List: list, Total: total, Page: page, PageSize: pageSize})
}

// Upload 上传普通文件，调用方需要具备 file:upload 权限。
func (h *Handler) Upload(c *gin.Context) {
	h.upload(c, false)
}

// UploadAvatar 上传头像文件，只允许 image/* 内容类型。
func (h *Handler) UploadAvatar(c *gin.Context) {
	h.upload(c, true)
}

func (h *Handler) upload(c *gin.Context, imageOnly bool) {
	header, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, "missing file")
		return
	}
	src, err := header.Open()
	if err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}
	defer src.Close()

	// objectKey 使用日期前缀便于对象存储分目录管理，Base 防止客户端文件名携带路径穿越。
	objectKey := time.Now().Format("2006/01/02/150405000") + "-" + filepath.Base(header.Filename)
	contentType := header.Header.Get("Content-Type")
	if imageOnly && !strings.HasPrefix(contentType, "image/") {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, "avatar must be an image")
		return
	}
	info, err := h.client.PutObject(c.Request.Context(), h.cfg.Bucket, objectKey, src, header.Size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, errs.CodeInternal, err.Error())
		return
	}
	record := File{
		OriginalName: header.Filename,
		Bucket:       h.cfg.Bucket,
		ObjectKey:    objectKey,
		ContentType:  contentType,
		Size:         info.Size,
		CreatedBy:    contextx.UserID(c),
	}
	if err := h.db.Create(&record).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, errs.CodeInternal, err.Error())
		return
	}
	response.Created(c, record)
}

// URL 生成内联访问的短期预签名 URL。
func (h *Handler) URL(c *gin.Context) {
	h.presignedURL(c, false)
}

// DownloadURL 生成下载用途的短期预签名 URL，并设置下载文件名。
func (h *Handler) DownloadURL(c *gin.Context) {
	h.presignedURL(c, true)
}

func (h *Handler) presignedURL(c *gin.Context, attachment bool) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, "invalid id")
		return
	}
	var record File
	if err := h.db.First(&record, uint(id64)).Error; err != nil {
		response.Fail(c, http.StatusNotFound, errs.CodeNotFound, "file not found")
		return
	}
	reqParams := url.Values{}
	if attachment {
		// 下载文件名只使用安全的 basename，避免响应头注入或路径信息泄露。
		reqParams.Set("response-content-disposition", `attachment; filename="`+escapeFilename(record.OriginalName)+`"`)
	}
	// 预签名 URL 只开放 15 分钟，避免数据库中的私有对象被长期公开访问。
	url, err := h.client.PresignedGetObject(context.Background(), record.Bucket, record.ObjectKey, 15*time.Minute, reqParams)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, errs.CodeInternal, err.Error())
		return
	}
	response.OK(c, gin.H{"url": url.String(), "expiresIn": 900})
}

// Delete 先删除对象存储中的文件，再软删除数据库元数据。
func (h *Handler) Delete(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, "invalid id")
		return
	}
	var record File
	if err := h.db.First(&record, uint(id64)).Error; err != nil {
		response.Fail(c, http.StatusNotFound, errs.CodeNotFound, "file not found")
		return
	}
	if err := h.client.RemoveObject(c.Request.Context(), record.Bucket, record.ObjectKey, minio.RemoveObjectOptions{}); err != nil {
		response.Fail(c, http.StatusInternalServerError, errs.CodeInternal, err.Error())
		return
	}
	if err := h.db.Delete(&record).Error; err != nil {
		response.Fail(c, http.StatusBadRequest, errs.CodeBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"deleted": true})
}

func escapeFilename(name string) string {
	return strings.ReplaceAll(filepath.Base(name), `"`, "")
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
