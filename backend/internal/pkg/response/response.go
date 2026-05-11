package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Body 是所有 HTTP JSON 响应的统一外层结构。
type Body struct {
	Code    int    `json:"code"`              // 业务错误码，0 表示成功。
	Message string `json:"message"`           // 面向前端展示或调试的结果说明。
	Data    any    `json:"data,omitempty"`    // 成功响应数据，失败时通常为空。
	TraceID string `json:"traceId,omitempty"` // 请求链路 ID，用于前后端联合排查问题。
}

// Page 是列表接口统一分页响应结构。
type Page[T any] struct {
	List     []T   `json:"list"`     // 当前页数据。
	Total    int64 `json:"total"`    // 符合查询条件的数据总数。
	Page     int   `json:"page"`     // 当前页码，从 1 开始。
	PageSize int   `json:"pageSize"` // 当前页容量。
}

// OK 返回 200 成功响应，并自动携带 traceId。
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{Code: 0, Message: "ok", Data: data, TraceID: traceID(c)})
}

// Created 返回 201 创建成功响应，并自动携带 traceId。
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Body{Code: 0, Message: "ok", Data: data, TraceID: traceID(c)})
}

// Fail 返回统一失败响应，HTTP 状态码和业务错误码由调用方明确传入。
func Fail(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Body{Code: code, Message: message, TraceID: traceID(c)})
}

func traceID(c *gin.Context) string {
	if v, ok := c.Get("traceId"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
