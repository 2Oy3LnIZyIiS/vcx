package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealth(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	health(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	body := w.Body.String()
	expected := `{"status":"ok"}`
	if body != expected {
		t.Errorf("Expected %q, got %q", expected, body)
	}
}

func TestPing(t *testing.T) {
	tests := []struct {
		method   string
		expected string
	}{
		{"GET", "GET"},
		{"POST", "POST"},
		{"CUSTOM", "UNKNOWN"},
		{"", "GET"},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/ping", nil)
			w := httptest.NewRecorder()

			ping(w, req)

			resp := w.Result()
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status 200, got %d", resp.StatusCode)
			}

			body := w.Body.String()
			if !strings.Contains(body, `"method":"`+tt.expected+`"`) {
				t.Errorf("Expected method %q in response, got %q", tt.expected, body)
			}
		})
	}
}
