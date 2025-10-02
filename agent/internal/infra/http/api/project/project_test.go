package project

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInitProject(t *testing.T) {
	req := httptest.NewRequest("GET", "/init", nil)
	w := httptest.NewRecorder()

	initProject(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body := w.Body.String()
	expected := "init project called"
	if body != expected {
		t.Errorf("Expected %q, got %q", expected, body)
	}
}

func TestHandler(t *testing.T) {
	handler := Handler()
	
	req := httptest.NewRequest("GET", "/api/project/init", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body := w.Body.String()
	if !strings.Contains(body, "init project called") {
		t.Errorf("Expected response to contain 'init project called', got %q", body)
	}
}

func TestAPIPath(t *testing.T) {
	expected := "/api/project"
	if APIPath != expected {
		t.Errorf("Expected APIPath to be %q, got %q", expected, APIPath)
	}
}