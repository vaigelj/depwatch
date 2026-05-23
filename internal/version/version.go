package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Info holds version metadata for a package.
type Info struct {
	Current string
	Latest  string
	IsLatest bool
}

// Fetcher retrieves the latest version of a package from a registry.
type Fetcher struct {
	client  *http.Client
	npmBase string
}

// New creates a new Fetcher with default settings.
func New() *Fetcher {
	return &Fetcher{
		client:  &http.Client{Timeout: 10 * time.Second},
		npmBase: "https://registry.npmjs.org",
	}
}

// NewWithBase creates a Fetcher with a custom registry base URL (useful for testing).
func NewWithBase(base string) *Fetcher {
	return &Fetcher{
		client:  &http.Client{Timeout: 10 * time.Second},
		npmBase: base,
	}
}

// FetchNPM returns the latest version string for an npm package.
func (f *Fetcher) FetchNPM(pkg string) (string, error) {
	url := fmt.Sprintf("%s/%s/latest", f.npmBase, pkg)
	resp, err := f.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("version: npm request failed for %q: %w", pkg, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("version: package %q not found on npm", pkg)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("version: npm returned status %d for %q", resp.StatusCode, pkg)
	}

	var payload struct {
		Version string `json:"version"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("version: failed to decode npm response for %q: %w", pkg, err)
	}
	if payload.Version == "" {
		return "", fmt.Errorf("version: empty version returned for npm package %q", pkg)
	}
	return payload.Version, nil
}

// Compare returns an Info struct comparing current vs latest.
func Compare(current, latest string) Info {
	return Info{
		Current:  current,
		Latest:   latest,
		IsLatest: current == latest,
	}
}
