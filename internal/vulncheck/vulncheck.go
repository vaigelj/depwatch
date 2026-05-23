// Package vulncheck queries vulnerability databases for known CVEs
// affecting scanned dependencies.
package vulncheck

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Vulnerability holds details about a known CVE for a package.
type Vulnerability struct {
	ID       string `json:"id"`
	Summary  string `json:"summary"`
	Severity string `json:"severity"`
}

// Result pairs a package name with any vulnerabilities found.
type Result struct {
	Package         string
	Ecosystem       string
	Vulnerabilities []Vulnerability
}

// Checker queries an OSV-compatible endpoint for vulnerability data.
type Checker struct {
	client  *http.Client
	baseURL string
}

// New returns a Checker using the public OSV API.
func New() *Checker {
	return &Checker{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: "https://api.osv.dev/v1",
	}
}

// NewWithBase returns a Checker with a custom base URL (useful for testing).
func NewWithBase(baseURL string) *Checker {
	return &Checker{
		client:  &http.Client{Timeout: 5 * time.Second},
		baseURL: baseURL,
	}
}

type osvQuery struct {
	Package osvPkg `json:"package"`
}

type osvPkg struct {
	Name      string `json:"name"`
	Ecosystem string `json:"ecosystem"`
}

type osvResponse struct {
	Vulns []struct {
		ID      string `json:"id"`
		Summary string `json:"summary"`
		Database string `json:"database_specific"`
	} `json:"vulns"`
}

// Check queries the vulnerability database for the given package and ecosystem.
func (c *Checker) Check(pkg, ecosystem string) (Result, error) {
	body, err := json.Marshal(osvQuery{Package: osvPkg{Name: pkg, Ecosystem: ecosystem}})
	if err != nil {
		return Result{}, fmt.Errorf("vulncheck: marshal: %w", err)
	}

	resp, err := c.client.Post(
		c.baseURL+"/query",
		"application/json",
		bytesReader(body),
	)
	if err != nil {
		return Result{}, fmt.Errorf("vulncheck: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Result{}, fmt.Errorf("vulncheck: unexpected status %d", resp.StatusCode)
	}

	var osv osvResponse
	if err := json.NewDecoder(resp.Body).Decode(&osv); err != nil {
		return Result{}, fmt.Errorf("vulncheck: decode: %w", err)
	}

	result := Result{Package: pkg, Ecosystem: ecosystem}
	for _, v := range osv.Vulns {
		result.Vulnerabilities = append(result.Vulnerabilities, Vulnerability{
			ID:       v.ID,
			Summary:  v.Summary,
			Severity: "critical",
		})
	}
	return result, nil
}
