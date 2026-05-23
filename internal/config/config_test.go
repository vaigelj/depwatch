package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/depwatch/depwatch/internal/config"
)

func TestDefault_Values(t *testing.T) {
	cfg := config.Default()
	if cfg.OutputFormat != "text" {
		t.Errorf("expected output_format \"text\", got %q", cfg.OutputFormat)
	}
	if !cfg.FailOnVulnerable {
		t.Error("expected fail_on_vulnerable to be true by default")
	}
	if cfg.FailOnOutdated {
		t.Error("expected fail_on_outdated to be false by default")
	}
	if len(cfg.Paths) == 0 || cfg.Paths[0] != "." {
		t.Errorf("expected default path \".\", got %v", cfg.Paths)
	}
}

func TestLoad_EmptyPath_ReturnsDefault(t *testing.T) {
	cfg, err := config.Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OutputFormat != "text" {
		t.Errorf("expected default output_format, got %q", cfg.OutputFormat)
	}
}

func TestLoad_ValidFile(t *testing.T) {
	data := map[string]interface{}{
		"paths":             []string{"/srv/app", "/srv/lib"},
		"output_format":     "json",
		"fail_on_vulnerable": false,
		"fail_on_outdated":  true,
	}
	path := writeJSON(t, data)

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OutputFormat != "json" {
		t.Errorf("expected \"json\", got %q", cfg.OutputFormat)
	}
	if !cfg.FailOnOutdated {
		t.Error("expected fail_on_outdated true")
	}
	if len(cfg.Paths) != 2 {
		t.Errorf("expected 2 paths, got %d", len(cfg.Paths))
	}
}

func TestLoad_InvalidFormat(t *testing.T) {
	data := map[string]interface{}{
		"paths":         []string{"."},
		"output_format": "xml",
	}
	path := writeJSON(t, data)

	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for unsupported output_format, got nil")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/depwatch.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func writeJSON(t *testing.T, v interface{}) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	path := filepath.Join(t.TempDir(), "depwatch.json")
	if err := os.WriteFile(path, b, 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	return path
}
