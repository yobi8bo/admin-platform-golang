package middleware

import (
	"net/http"
	"testing"
)

func TestShouldSkipOperationAudit(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
		want   bool
	}{
		{name: "skip file list", method: http.MethodGet, path: "/api/files", want: true},
		{name: "skip file url", method: http.MethodGet, path: "/api/files/2/url", want: true},
		{name: "skip download url", method: http.MethodGet, path: "/api/files/2/download-url", want: true},
		{name: "skip operation log list", method: http.MethodGet, path: "/api/audit/operation-logs", want: true},
		{name: "record file upload", method: http.MethodPost, path: "/api/files/upload", want: false},
		{name: "record profile update", method: http.MethodPut, path: "/api/auth/profile", want: false},
		{name: "record file delete", method: http.MethodDelete, path: "/api/files/2", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldSkipOperationAudit(tt.method, tt.path); got != tt.want {
				t.Fatalf("shouldSkipOperationAudit(%q, %q) = %v, want %v", tt.method, tt.path, got, tt.want)
			}
		})
	}
}
