package audit

import (
	"net/http"
	"strings"
	"testing"
)

func TestDescribeOperationUsesChineseLabels(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
		want   string
	}{
		{name: "list files", method: http.MethodGet, path: "/api/files", want: "查看文件"},
		{name: "delete file", method: http.MethodDelete, path: "/api/files/2", want: "删除文件"},
		{name: "file url", method: http.MethodGet, path: "/api/files/2/url", want: "查看文件"},
		{name: "download url", method: http.MethodGet, path: "/api/files/2/download-url", want: "查看文件"},
		{name: "upload file", method: http.MethodPost, path: "/api/files/upload", want: "上传文件"},
		{name: "delete user", method: http.MethodDelete, path: "/api/system/users/2", want: "删除用户"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := describeOperation(tt.method, tt.path); got != tt.want {
				t.Fatalf("describeOperation(%q, %q) = %q, want %q", tt.method, tt.path, got, tt.want)
			}
		})
	}
}

func TestDescribeOperationDoesNotExposeUnknownAPIPath(t *testing.T) {
	got := describeOperation(http.MethodGet, "/api/internal/raw/42")
	if strings.Contains(got, "internal") || strings.Contains(got, "raw") || strings.Contains(got, "/") {
		t.Fatalf("describeOperation exposed API path: %q", got)
	}
	if got != "查看" {
		t.Fatalf("describeOperation() = %q, want %q", got, "查看")
	}
}
