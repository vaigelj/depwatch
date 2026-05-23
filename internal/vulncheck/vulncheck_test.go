package vulncheck_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"depwatch/internal/vulncheck"
)

func mockOSVServer(t *testing.T, payload any, status int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(payload)
	}))
}

func TestCheck_NoVulnerabilities(t *testing.T) {
	srv := mockOSVServer(t, map[string]any{"vulns": []any{}}, http.StatusOK)
	defer srv.Close()

	c := vulncheck.NewWithBase(srv.URL)
	res, err := c.Check("express", "npm")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Vulnerabilities) != 0 {
		t.Errorf("expected 0 vulns, got %d", len(res.Vulnerabilities))
	}
}

func TestCheck_WithVulnerabilities(t *testing.T) {
	payload := map[string]any{
		"vulns": []any{
			map[string]any{"id": "GHSA-1234", "summary": "Remote code execution"},
			map[string]any{"id": "GHSA-5678", "summary": "SQL injection"},
		},
	}
	srv := mockOSVServer(t, payload, http.StatusOK)
	defer srv.Close()

	c := vulncheck.NewWithBase(srv.URL)
	res, err := c.Check("lodash", "npm")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Vulnerabilities) != 2 {
		t.Fatalf("expected 2 vulns, got %d", len(res.Vulnerabilities))
	}
	if res.Vulnerabilities[0].ID != "GHSA-1234" {
		t.Errorf("unexpected vuln ID: %s", res.Vulnerabilities[0].ID)
	}
	if res.Vulnerabilities[0].Severity != "critical" {
		t.Errorf("expected severity critical, got %s", res.Vulnerabilities[0].Severity)
	}
}

func TestCheck_ServerError(t *testing.T) {
	srv := mockOSVServer(t, map[string]any{}, http.StatusInternalServerError)
	defer srv.Close()

	c := vulncheck.NewWithBase(srv.URL)
	_, err := c.Check("requests", "PyPI")
	if err == nil {
		t.Fatal("expected error for non-200 status, got nil")
	}
}

func TestNew_ReturnsChecker(t *testing.T) {
	c := vulncheck.New()
	if c == nil {
		t.Fatal("expected non-nil checker")
	}
}
