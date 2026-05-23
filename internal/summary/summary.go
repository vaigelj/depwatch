// Package summary aggregates scan, check, and alert results into a
// structured report summary suitable for display or downstream processing.
package summary

import (
	"fmt"
	"time"

	"github.com/user/depwatch/internal/alert"
)

// Severity level constants mirror alert severities for summary bucketing.
const (
	LevelClean    = "clean"
	LevelWarn     = "warn"
	LevelCritical = "critical"
)

// Summary holds aggregated statistics from a depwatch run.
type Summary struct {
	ScannedAt  time.Time `json:"scanned_at"`
	TotalDeps  int       `json:"total_deps"`
	Clean      int       `json:"clean"`
	Warn       int       `json:"warn"`
	Critical   int       `json:"critical"`
	Ecosystems []string  `json:"ecosystems"`
}

// Build constructs a Summary from a slice of alerts and the total dependency
// count recorded by the scanner.
func Build(alerts []alert.Alert, totalDeps int) Summary {
	s := Summary{
		ScannedAt: time.Now().UTC(),
		TotalDeps: totalDeps,
	}

	seen := map[string]struct{}{}
	for _, a := range alerts {
		switch a.Level {
		case alert.LevelWarn:
			s.Warn++
		case alert.LevelCritical:
			s.Critical++
		}
		if _, ok := seen[a.Ecosystem]; !ok {
			seen[a.Ecosystem] = struct{}{}
			s.Ecosystems = append(s.Ecosystems, a.Ecosystem)
		}
	}
	s.Clean = totalDeps - s.Warn - s.Critical
	if s.Clean < 0 {
		s.Clean = 0
	}
	return s
}

// OverallLevel returns the highest severity level present in the summary.
func (s Summary) OverallLevel() string {
	if s.Critical > 0 {
		return LevelCritical
	}
	if s.Warn > 0 {
		return LevelWarn
	}
	return LevelClean
}

// String returns a human-readable one-line representation of the summary.
func (s Summary) String() string {
	return fmt.Sprintf(
		"[%s] total=%d clean=%d warn=%d critical=%d ecosystems=%v",
		s.OverallLevel(), s.TotalDeps, s.Clean, s.Warn, s.Critical, s.Ecosystems,
	)
}
