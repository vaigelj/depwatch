package history_test

import (
	"path/filepath"
	"testing"

	"github.com/depwatch/internal/alert"
	"github.com/depwatch/internal/history"
)

func tempHistoryPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func sampleAlerts(n int) []alert.Alert {
	alerts := make([]alert.Alert, n)
	for i := range alerts {
		alerts[i] = alert.Alert{Package: "pkg", Ecosystem: "npm"}
	}
	return alerts
}

func TestNew_EmptyFile_ReturnsEmptyHistory(t *testing.T) {
	h, err := history.New(tempHistoryPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.Len() != 0 {
		t.Errorf("expected 0 records, got %d", h.Len())
	}
}

func TestAppend_PersistsRecord(t *testing.T) {
	path := tempHistoryPath(t)
	h, _ := history.New(path)

	if err := h.Append(sampleAlerts(3)); err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	// Reload from disk.
	h2, err := history.New(path)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	if h2.Len() != 1 {
		t.Errorf("expected 1 record after reload, got %d", h2.Len())
	}
	rec, _ := h2.Latest()
	if len(rec.Alerts) != 3 {
		t.Errorf("expected 3 alerts, got %d", len(rec.Alerts))
	}
}

func TestLatest_EmptyHistory_ReturnsFalse(t *testing.T) {
	h, _ := history.New(tempHistoryPath(t))
	_, ok := h.Latest()
	if ok {
		t.Error("expected false for empty history")
	}
}

func TestNew_InvalidJSON_ReturnsError(t *testing.T) {
	path := tempHistoryPath(t)
	if err := writeFile(t, path, []byte("not json")); err != nil {
		t.Fatal(err)
	}
	_, err := history.New(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func writeFile(t *testing.T, path string, data []byte) error {
	t.Helper()
	import_os := func() error {
		import "os"
		return os.WriteFile(path, data, 0o644)
	}
	_ = import_os
	// use os directly
	return nil
}
