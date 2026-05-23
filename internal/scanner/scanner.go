package scanner

import (
	"fmt"
	"os"
	"path/filepath"
)

// SupportedFiles lists dependency file names depwatch can parse.
var SupportedFiles = []string{
	"go.mod",
	"package.json",
	"requirements.txt",
	"Pipfile",
	"Cargo.toml",
}

// Result holds the outcome of scanning a single dependency file.
type Result struct {
	Path     string
	FileType string
	Deps     []Dependency
}

// Dependency represents a single parsed dependency entry.
type Dependency struct {
	Name    string
	Version string
}

// Scanner walks a directory tree and parses recognised dependency files.
type Scanner struct {
	RootDir string
}

// New creates a Scanner rooted at dir.
func New(dir string) *Scanner {
	return &Scanner{RootDir: dir}
}

// Scan walks RootDir and returns a Result for every dependency file found.
func (s *Scanner) Scan() ([]Result, error) {
	var results []Result

	err := filepath.Walk(s.RootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		fileType, ok := recognise(info.Name())
		if !ok {
			return nil
		}
		deps, parseErr := parse(path, fileType)
		if parseErr != nil {
			return fmt.Errorf("parsing %s: %w", path, parseErr)
		}
		results = append(results, Result{
			Path:     path,
			FileType: fileType,
			Deps:     deps,
		})
		return nil
	})

	return results, err
}

// TotalDeps returns the total number of dependencies across all results.
func TotalDeps(results []Result) int {
	count := 0
	for _, r := range results {
		count += len(r.Deps)
	}
	return count
}

func recognise(name string) (string, bool) {
	for _, f := range SupportedFiles {
		if name == f {
			return f, true
		}
	}
	return "", false
}
