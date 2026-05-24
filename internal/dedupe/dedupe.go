// Package dedupe provides utilities for deduplicating alerts
// so that repeated findings across multiple scans are collapsed
// into a single representative entry.
package dedupe

import (
	"fmt"

	"github.com/user/depwatch/internal/alert"
)

// key uniquely identifies an alert by ecosystem, package, and severity.
func key(a alert.Alert) string {
	return fmt.Sprintf("%s::%s::%s", a.Ecosystem, a.Package, a.Severity)
}

// ByKey deduplicates a slice of alerts, keeping the first occurrence of each
// unique (ecosystem, package, severity) combination. Order is preserved.
func ByKey(alerts []alert.Alert) []alert.Alert {
	seen := make(map[string]struct{}, len(alerts))
	out := make([]alert.Alert, 0, len(alerts))

	for _, a := range alerts {
		k := key(a)
		if _, exists := seen[k]; exists {
			continue
		}
		seen[k] = struct{}{}
		out = append(out, a)
	}

	return out
}

// Count returns the number of duplicate alerts that would be removed.
func Count(alerts []alert.Alert) int {
	return len(alerts) - len(ByKey(alerts))
}
