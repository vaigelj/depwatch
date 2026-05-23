package alert

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/depwatch/internal/checker"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo     Level = "INFO"
	LevelWarn     Level = "WARN"
	LevelCritical Level = "CRITICAL"
)

// Alert holds information about a single dependency issue.
type Alert struct {
	Package   string
	Current   string
	Latest    string
	Level     Level
	Message   string
}

// Alerter formats and writes alerts to an output destination.
type Alerter struct {
	out io.Writer
}

// New returns an Alerter that writes to the given writer.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Alerter {
	if w == nil {
		w = os.Stdout
	}
	return &Alerter{out: w}
}

// FromResults converts checker results into Alerts.
func FromResults(results []checker.Result) []Alert {
	alerts := make([]Alert, 0, len(results))
	for _, r := range results {
		if !r.Outdated && !r.Vulnerable {
			continue
		}
		a := Alert{
			Package: r.Package,
			Current: r.Current,
			Latest:  r.Latest,
		}
		switch {
		case r.Vulnerable:
			a.Level = LevelCritical
			a.Message = fmt.Sprintf("vulnerable package %s@%s — upgrade to %s", r.Package, r.Current, r.Latest)
		case r.Outdated:
			a.Level = LevelWarn
			a.Message = fmt.Sprintf("outdated package %s@%s — latest is %s", r.Package, r.Current, r.Latest)
		}
		alerts = append(alerts, a)
	}
	return alerts
}

// Print writes all alerts to the Alerter's output in a human-readable format.
func (a *Alerter) Print(alerts []Alert) {
	if len(alerts) == 0 {
		fmt.Fprintln(a.out, "[depwatch] all dependencies are up to date.")
		return
	}
	fmt.Fprintf(a.out, "[depwatch] %d issue(s) found:\n", len(alerts))
	fmt.Fprintln(a.out, strings.Repeat("-", 60))
	for _, al := range alerts {
		fmt.Fprintf(a.out, "[%s] %s\n", al.Level, al.Message)
	}
}
