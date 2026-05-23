package version_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/depwatch/depwatch/internal/version"
)

func TestCompare_IsLatest(t *testing.T) {
	info := version.Compare("1.2.3", "1.2.3")
	if !info.IsLatest {
		t.Errorf("expected IsLatest=true for equal versions")
	}
}

func TestCompare_Outdated(t *testing.T) {
	info := version.Compare("1.0.0", "2.0.0")
	if info.IsLatest {
		t.Errorf("expected IsLatest=false for outdated version")
	}
	if info.Current != "1.0.0" || info.Latest != "2.0.0" {
		t.Errorf("unexpected version fields: %+v", info)
	}
}

func TestFetchNPM_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"version": "3.1.4"})
	}))
	defer server.Close()

	f := version.NewWithBase(server.URL)
	v, err := f.FetchNPM("lodash")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "3.1.4" {
		t.Errorf("expected 3.1.4, got %s", v)
	}
}

func TestFetchNPM_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	f := version.NewWithBase(server.URL)
	_, err := f.FetchNPM("nonexistent-pkg-xyz")
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

func TestFetchNPM_EmptyVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"version": ""})
	}))
	defer server.Close()

	f := version.NewWithBase(server.URL)
	_, err := f.FetchNPM("some-pkg")
	if err == nil {
		t.Fatal("expected error for empty version, got nil")
	}
}

func TestNew_ReturnsFetcher(t *testing.T) {
	f := version.New()
	if f == nil {
		t.Fatal("expected non-nil Fetcher")
	}
}
