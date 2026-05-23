package filter

import "github.com/user/depwatch/internal/alert"

// Severity represents the minimum severity level to include in output.
type Severity int

const (
	SeverityAll      Severity = iota
	SeverityWarn
	SeverityCritical
)

// ParseSeverity converts a string to a Severity level.
func ParseSeverity(s string) (Severity, bool) {
	switch s {
	case "all", "":
		return SeverityAll, true
	case "warn":
		return SeverityWarn, true
	case "critical":
		return SeverityCritical, true
	}
	return SeverityAll, false
}

// Options holds filtering criteria applied to alerts.
type Options struct {
	MinSeverity Severity
	Ecosystem   string // empty means all
}

// Apply returns only the alerts that match the given options.
func Apply(alerts []alert.Alert, opts Options) []alert.Alert {
	out := make([]alert.Alert, 0, len(alerts))
	for _, a := range alerts {
		if opts.Ecosystem != "" && a.Ecosystem != opts.Ecosystem {
			continue
		}
		if !meetsSeverity(a.Level, opts.MinSeverity) {
			continue
		}
		out = append(out, a)
	}
	return out
}

func meetsSeverity(level alert.Level, min Severity) bool {
	switch min {
	case SeverityCritical:
		return level == alert.LevelCritical
	case SeverityWarn:
		return level == alert.LevelWarn || level == alert.LevelCritical
	default:
		return true
	}
}
