package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "cache.json")
}

func TestNew_CreatesEmptyCache(t *testing.T) {
	c, err := New(tempPath(t), time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil cache")
	}
}

func TestSetAndGet_ReturnsValue(t *testing.T) {
	c, _ := New(tempPath(t), time.Hour)
	if err := c.Set("npm:lodash", "4.17.21"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	v, ok := c.Get("npm:lodash")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if v != "4.17.21" {
		t.Errorf("got %q, want %q", v, "4.17.21")
	}
}

func TestGet_MissingKey_ReturnsFalse(t *testing.T) {
	c, _ := New(tempPath(t), time.Hour)
	_, ok := c.Get("npm:nonexistent")
	if ok {
		t.Fatal("expected cache miss")
	}
}

func TestGet_ExpiredEntry_ReturnsFalse(t *testing.T) {
	c, _ := New(tempPath(t), -time.Second) // TTL already expired
	_ = c.Set("npm:react", "18.0.0")
	_, ok := c.Get("npm:react")
	if ok {
		t.Fatal("expected expired entry to be a miss")
	}
}

func TestPersistence_ReloadsFromDisk(t *testing.T) {
	path := tempPath(t)
	c1, _ := New(path, time.Hour)
	_ = c1.Set("pypi:requests", "2.31.0")

	c2, err := New(path, time.Hour)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	v, ok := c2.Get("pypi:requests")
	if !ok {
		t.Fatal("expected value after reload")
	}
	if v != "2.31.0" {
		t.Errorf("got %q, want %q", v, "2.31.0")
	}
}

func TestNew_InvalidJSON_ReturnsError(t *testing.T) {
	path := tempPath(t)
	_ = os.WriteFile(path, []byte("not-json"), 0o644)
	_, err := New(path, time.Hour)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
