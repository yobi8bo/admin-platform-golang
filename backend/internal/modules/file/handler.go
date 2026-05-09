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

type Handler struct {
	db     *gorm.DB
	client *minio.Client
	cfg    config.RustFSConfig
}

func NewHandler(db *gorm.DB, client *minio.Client, cfg config.RustFSConfig) *Handler {
	return &Handler{db: db, client: client, cfg: cfg}
}

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

func (h *Handler) Upload(c *gin.Context) {
	h.upload(c, false)
}

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

func (h *Handler) URL(c *gin.Context) {
	h.presignedURL(c, false)
}

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
		reqParams.Set("response-content-disposition", `attachment; filename="`+escapeFilename(record.OriginalName)+`"`)
	}
	url, err := h.client.PresignedGetObject(context.Background(), record.Bucket, record.ObjectKey, 15*time.Minute, reqParams)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, errs.CodeInternal, err.Error())
		return
	}
	response.OK(c, gin.H{"url": url.String(), "expiresIn": 900})
}

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
