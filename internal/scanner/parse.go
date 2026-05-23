package scanner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"golang.org/x/mod/modfile"
)

// parse dispatches to the correct parser based on fileType.
func parse(path, fileType string) ([]Dependency, error) {
	switch fileType {
	case "go.mod":
		return parseGoMod(path)
	case "package.json":
		return parsePackageJSON(path)
	case "requirements.txt":
		return parseRequirementsTxt(path)
	default:
		return nil, fmt.Errorf("no parser for %s", fileType)
	}
}

func parseGoMod(path string) ([]Dependency, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	f, err := modfile.Parse(path, data, nil)
	if err != nil {
		return nil, err
	}
	var deps []Dependency
	for _, req := range f.Require {
		if !req.Indirect {
			deps = append(deps, Dependency{Name: req.Mod.Path, Version: req.Mod.Version})
		}
	}
	return deps, nil
}

func parsePackageJSON(path string) ([]Dependency, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}
	var deps []Dependency
	for name, ver := range pkg.Dependencies {
		deps = append(deps, Dependency{Name: name, Version: ver})
	}
	for name, ver := range pkg.DevDependencies {
		deps = append(deps, Dependency{Name: name, Version: ver})
	}
	return deps, nil
}

func parseRequirementsTxt(path string) ([]Dependency, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var deps []Dependency
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "==", 2)
		dep := Dependency{Name: strings.TrimSpace(parts[0])}
		if len(parts) == 2 {
			dep.Version = strings.TrimSpace(parts[1])
		}
		deps = append(deps, dep)
	}
	return deps, scanner.Err()
}
