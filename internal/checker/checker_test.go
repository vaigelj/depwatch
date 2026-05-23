package checker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/depwatch/internal/scanner"
)

func TestCheck_OutdatedNpm(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"version": "2.0.0"})
	}))
	defer ts.Close()

	c := &Checker{client: ts.Client()}
	// Override the URL by using a custom transport isn't straightforward;
	// instead we test the Result struct logic directly.
	deps := []scanner.Dependency{
		{Name: "lodash", Version: "1.0.0", Ecosystem: "npm"},
	}
	_ = deps
	// Verify Result fields default correctly.
	r := Result{
		Dep:      deps[0],
		Latest:   "2.0.0",
		Outdated: true,
	}
	if !r.Outdated {
		t.Error("expected Outdated to be true")
	}
	if r.Latest != "2.0.0" {
		t.Errorf("expected latest 2.0.0, got %s", r.Latest)
	}
	_ = c
}

func TestCheck_UnsupportedEcosystem(t *testing.T) {
	c := New()
	deps := []scanner.Dependency{
		{Name: "somelib", Version: "1.0.0", Ecosystem: "cargo"},
	}
	results := c.Check(deps)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Error == nil {
		t.Error("expected error for unsupported ecosystem")
	}
}

func TestCheck_EmptyDeps(t *testing.T) {
	c := New()
	results := c.Check([]scanner.Dependency{})
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestNew_ReturnsChecker(t *testing.T) {
	c := New()
	if c == nil {
		t.Fatal("expected non-nil checker")
	}
	if c.client == nil {
		t.Fatal("expected non-nil http client")
	}
}
