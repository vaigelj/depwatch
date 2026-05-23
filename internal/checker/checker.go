package checker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/depwatch/internal/scanner"
)

// Result holds the outcome of checking a single dependency.
type Result struct {
	Dep       scanner.Dependency
	Latest    string
	Outdated  bool
	VulnCount int
	Error     error
}

// Checker fetches version and vulnerability data for dependencies.
type Checker struct {
	client *http.Client
}

// New returns a Checker with a sensible default HTTP client.
func New() *Checker {
	return &Checker{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Check evaluates a slice of dependencies and returns results.
func (c *Checker) Check(deps []scanner.Dependency) []Result {
	results := make([]Result, 0, len(deps))
	for _, d := range deps {
		r := Result{Dep: d}
		latest, err := c.latestVersion(d)
		if err != nil {
			r.Error = err
			results = append(results, r)
			continue
		}
		r.Latest = latest
		r.Outdated = latest != d.Version
		results = append(results, r)
	}
	return results
}

// latestVersion queries the appropriate registry for the latest version.
func (c *Checker) latestVersion(d scanner.Dependency) (string, error) {
	switch d.Ecosystem {
	case "npm":
		return c.npmLatest(d.Name)
	case "pypi":
		return c.pypiLatest(d.Name)
	default:
		return "", fmt.Errorf("unsupported ecosystem: %s", d.Ecosystem)
	}
}

func (c *Checker) npmLatest(pkg string) (string, error) {
	url := fmt.Sprintf("https://registry.npmjs.org/%s/latest", pkg)
	resp, err := c.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var data struct {
		Version string `json:"version"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	return data.Version, nil
}

func (c *Checker) pypiLatest(pkg string) (string, error) {
	url := fmt.Sprintf("https://pypi.org/pypi/%s/json", pkg)
	resp, err := c.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var data struct {
		Info struct {
			Version string `json:"version"`
		} `json:"info"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	return data.Info.Version, nil
}
