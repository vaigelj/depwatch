package ignore_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/depwatch/internal/ignore"
)

func writeIgnoreFile(t *testing.T, l ignore.List) string {
	t.Helper()
	data, err := json.Marshal(l)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	p := filepath.Join(t.TempDir(), ".depignore")
	if err := os.WriteFile(p, data, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	return p
}

func TestLoad_EmptyPath_ReturnsEmpty(t *testing.T) {
	l, err := ignore.Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(l.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(l.Entries))
	}
}

func TestLoad_MissingFile_ReturnsEmpty(t *testing.T) {
	l, err := ignore.Load("/nonexistent/.depignore")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(l.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(l.Entries))
	}
}

func TestLoad_ValidFile_ReturnsEntries(t *testing.T) {
	l := ignore.List{
		Entries: []ignore.Entry{
			{Ecosystem: "npm", Package: "lodash", Reason: "accepted risk"},
		},
	}
	path := writeIgnoreFile(t, l)
	got, err := ignore.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got.Entries))
	}
	if got.Entries[0].Package != "lodash" {
		t.Errorf("expected lodash, got %s", got.Entries[0].Package)
	}
}

func TestContains_Match_ReturnsTrue(t *testing.T) {
	l := &ignore.List{
		Entries: []ignore.Entry{
			{Ecosystem: "pypi", Package: "requests"},
		},
	}
	if !l.Contains("pypi", "requests") {
		t.Error("expected Contains to return true")
	}
}

func TestContains_Expired_ReturnsFalse(t *testing.T) {
	l := &ignore.List{
		Entries: []ignore.Entry{
			{Ecosystem: "npm", Package: "chalk", Expires: time.Now().Add(-time.Hour)},
		},
	}
	if l.Contains("npm", "chalk") {
		t.Error("expected Contains to return false for expired entry")
	}
}

func TestActive_FiltersExpired(t *testing.T) {
	l := &ignore.List{
		Entries: []ignore.Entry{
			{Ecosystem: "npm", Package: "active", Expires: time.Now().Add(time.Hour)},
			{Ecosystem: "npm", Package: "expired", Expires: time.Now().Add(-time.Hour)},
			{Ecosystem: "go", Package: "noexpiry"},
		},
	}
	active := l.Active()
	if len(active) != 2 {
		t.Errorf("expected 2 active entries, got %d", len(active))
	}
}
