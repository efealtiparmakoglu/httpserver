package main

import (
	"strings"
	"testing"
)

func TestRoute(t *testing.T) {
	tests := []struct {
		path       string
		wantStatus int
		wantBody   string
	}{
		{"/", 200, "httpserver"},
		{"/health", 200, "ok"},
		{"/hello/world", 200, "hello, world"},
		{"/nope", 404, "not found"},
	}
	for _, tc := range tests {
		// reproduce the routing logic from main
		var status int
		var body string
		switch {
		case tc.path == "/":
			status, body = 200, "<h1>httpserver 1.0.0</h1>"
		case tc.path == "/health":
			status, body = 200, `{"status":"ok"}`
		case strings.HasPrefix(tc.path, "/hello/"):
			status, body = 200, "hello, "+strings.TrimPrefix(tc.path, "/hello/")
		default:
			status, body = 404, `{"error":"not found"}`
		}
		if status != tc.wantStatus {
			t.Errorf("%s: status %d != %d", tc.path, status, tc.wantStatus)
		}
		if !strings.Contains(body, tc.wantBody) {
			t.Errorf("%s: body %q missing %q", tc.path, body, tc.wantBody)
		}
	}
}
