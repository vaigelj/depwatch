package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/user/depwatch/internal/checker"
)

// Level represents the severity of an alert.
type Level int

const (
	LevelInfo     Level = iota
	LevelWarn
	LevelCritical
)

func (l Level) String() string {
	switch l {
	case LevelWarn:
		return "WARN"
	case LevelCritical:
		return "CRITICAL"
	default:
		return "INFO"
	}
}

// Alert represents a single actionable finding for a dependency.
type Alert struct {
	Package   string
	Ecosystem string
	Level     Level
	Message   string
}

// Alerter writes alerts to an output writer.
type Alerter struct {
	w io.Writer
}

// New returns an Alerter that writes to stdout.
func New() *Alerter {
	return &Alerter{w: os.Stdout}
}

// FromResults converts checker results into alerts, skipping clean deps.
func FromResults(results []checker.Result) []Alert {
	var alerts []Alert
	for _, r := range results {
		switch {
		case r.Vulnerable:
			alerts = append(alerts, Alert{
				Package:   r.Package,
				Ecosystem: r.Ecosystem,
				Level:     LevelCritical,
				Message:   fmt.Sprintf("vulnerable: %s", r.Advisory),
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
