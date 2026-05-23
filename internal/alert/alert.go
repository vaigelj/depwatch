// Package alert defines the Alert type and helpers for producing alerts from
// dependency check results.
package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/user/depwatch/internal/checker"
)

// Level constants describe alert severity.
const (
	LevelClean    = "clean"
	LevelWarn     = "warn"
	LevelCritical = "critical"
)

// Alert represents a single actionable finding for one dependency.
type Alert struct {
	Package   string `json:"package"`
	Ecosystem string `json:"ecosystem"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

// Alerter prints alerts to a configurable writer.
type Alerter struct {
	w io.Writer
}

// New returns an Alerter that writes to os.Stdout.
func New() *Alerter {
	return &Alerter{w: os.Stdout}
}

// FromResults converts checker results into a slice of Alerts, omitting clean
// results unless they carry a vulnerability.
func FromResults(results []checker.Result) []Alert {
	var alerts []Alert
	for _, r := range results {
		switch {
		case r.Vulnerable:
			alerts = append(alerts, Alert{
				Package:   r.Package,
				Ecosystem: r.Ecosystem,
				Level:     LevelCritical,
				Message:   fmt.Sprintf("vulnerable: %s", r.VulnID),
			})
		case r.Outdated:
			alerts = append(alerts, Alert{
				Package:   r.Package,
				Ecosystem: r.Ecosystem,
				Level:     LevelWarn,
				Message:   fmt.Sprintf("outdated: current=%s latest=%s", r.Current, r.Latest),
			})
		}
	}
	return alerts
}

// Print writes all alerts to the Alerter's writer.
func (a *Alerter) Print(alerts []Alert) {
	if len(alerts) == 0 {
		fmt.Fprintln(a.w, "No issues found.")
		return
	}
	for _, al := range alerts {
		fmt.Fprintf(a.w, "[%s] %s (%s): %s\n", al.Level, al.Package, al.Ecosystem, al.Message)
	}
}
