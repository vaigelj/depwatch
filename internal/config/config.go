package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the runtime configuration for depwatch.
type Config struct {
	// Paths is the list of directories to scan.
	Paths []string `json:"paths"`
	// OutputFormat controls how results are reported: "text" or "json".
	OutputFormat string `json:"output_format"`
	// FailOnVulnerable causes a non-zero exit when vulnerabilities are found.
	FailOnVulnerable bool `json:"fail_on_vulnerable"`
	// FailOnOutdated causes a non-zero exit when outdated packages are found.
	FailOnOutdated bool `json:"fail_on_outdated"`
}

const defaultOutputFormat = "text"

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		Paths:            []string{"."},
		OutputFormat:     defaultOutputFormat,
		FailOnVulnerable: true,
		FailOnOutdated:   false,
	}
}

// Load reads a JSON config file from path and merges it over the defaults.
// If path is empty, Default() is returned.
func Load(path string) (*Config, error) {
	cfg := Default()
	if path == "" {
		return cfg, nil
	}

	clean := filepath.Clean(path)
	f, err := os.Open(clean)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", clean, err)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", clean, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	switch c.OutputFormat {
	case "text", "json":
		// valid
	default:
		return fmt.Errorf("config: unsupported output_format %q (want \"text\" or \"json\")", c.OutputFormat)
	}
	if len(c.Paths) == 0 {
		return fmt.Errorf("config: paths must not be empty")
	}
	return nil
}
