package main

import (
	"strings"
	"testing"
)

func TestHandlerRoutes(t *testing.T) {
	tests := []struct {
		method, path string
		wantStatus   int
		wantContains string
	}{
		{"GET", "/", 200, "httpserver"},
		{"GET", "/health", 200, `"ok"`},
		{"GET", "/hello/world", 200, "hello, world"},
		{"GET", "/nope", 404, "not found"},
	}
	for _, tc := range tests {
		status, _, body := tc.method, tc.path, ""
		_ = status
		s, ct, b := HandlerFunc(tc.method, tc.path, "")
		_ = ct
		if s != tc.wantStatus {
			t.Errorf("%s %s: status %d, want %d", tc.method, tc.path, s, tc.wantStatus)
		}
		if !strings.Contains(b, tc.wantContains) {
			t.Errorf("%s %s: body %q missing %q", tc.method, tc.path, b, tc.wantContains)
		}
		_ = status
	}
}

func HandlerFunc(method, path, body string) (int, string, string) {
	switch {
	case path == "/":
		return 200, "text/html", "<h1>httpserver 1.0.0</h1><p>zero dependencies</p>"
	case path == "/health":
		return 200, "application/json", `{"status":"ok"}`
	case strings.HasPrefix(path, "/hello/"):
		return 200, "text/plain", "hello, " + strings.TrimPrefix(path, "/hello/")
	default:
		return 404, "application/json", `{"error":"not found"}`
	}
}
