package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/depwatch/internal/checker"
)

// Severity levels for alerts.
const (
	SeverityWarn     = "warn"
	SeverityCritical = "critical"
)

// Alert represents a single actionable finding for a dependency.
type Alert struct {
	Package  string `json:"package"`
	Version  string `json:"version"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

// Alerter prints alerts to an output destination.
type Alerter struct {
	out io.Writer
}

// New creates a new Alerter writing to stdout.
func New() *Alerter {
	return &Alerter{out: os.Stdout}
}

// FromResults converts checker results into Alert values, filtering clean deps.
func FromResults(results []checker.Result) []Alert {
	var alerts []Alert
	for _, r := range results {
		switch {
		case r.Vulnerable:
			alerts = append(alerts, Alert{
				Package:  r.Package,
				Version:  r.Version,
				Severity: SeverityCritical,
				Message:  r.Advisory,
			})
		case r.Outdated:
			alerts = append(alerts, Alert{
				Package:  r.Package,
				Version:  r.Version,
				Severity: SeverityWarn,
				Message:  fmt.Sprintf("outdated: current=%s latest=%s", r.Version, r.Latest),
			})
		}
	}
	return alerts
}

// Print writes all alerts to the Alerter's output.
func (a *Alerter) Print(alerts []Alert) {
	if len(alerts) == 0 {
		fmt.Fprintln(a.out, "All dependencies look good.")
		return
	}
	for _, al := range alerts {
		fmt.Fprintf(a.out, "[%s] %s@%s: %s\n", al.Severity, al.Package, al.Version, al.Message)
	}
}
