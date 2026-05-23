package scanner_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/depwatch/internal/scanner"
)

func TestScan_RequirementsTxt(t *testing.T) {
	dir := t.TempDir()
	content := `# comment
requests==2.31.0
flask==3.0.0
`
	writeFile(t, filepath.Join(dir, "requirements.txt"), content)

	s := scanner.New(dir)
	results, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if len(results[0].Deps) != 2 {
		t.Errorf("expected 2 deps, got %d", len(results[0].Deps))
	}
}

func TestScan_PackageJSON(t *testing.T) {
	dir := t.TempDir()
	content := `{"dependencies":{"express":"^4.18.2"},"devDependencies":{"jest":"^29.0.0"}}`
	writeFile(t, filepath.Join(dir, "package.json"), content)

	s := scanner.New(dir)
	results, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if len(results[0].Deps) != 2 {
		t.Errorf("expected 2 deps, got %d", len(results[0].Deps))
	}
}

func TestScan_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	s := scanner.New(dir)
	results, err := s.Scan()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestScan_IgnoresUnknownFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "Makefile"), "all:\n\techo hi")
	s := scanner.New(dir)
	results, _ := s.Scan()
	if len(results) != 0 {
		t.Errorf("expected 0 results for unknown file, got %d", len(results))
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
}
