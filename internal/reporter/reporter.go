package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/depwatch/internal/alert"
)

// Format defines the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report holds the full scan report data.
type Report struct {
	GeneratedAt time.Time     `json:"generated_at"`
	TotalAlerts int           `json:"total_alerts"`
	Alerts      []alert.Alert `json:"alerts"`
}

// Reporter writes scan reports to an output destination.
type Reporter struct {
	format Format
	out    io.Writer
}

// New creates a new Reporter with the given format.
func New(format Format) *Reporter {
	return &Reporter{
		format: format,
		out:    os.Stdout,
	}
}

// NewWithWriter creates a new Reporter writing to the provided writer.
func NewWithWriter(format Format, w io.Writer) *Reporter {
	return &Reporter{format: format, out: w}
}

// Write renders the report to the configured output.
func (r *Reporter) Write(alerts []alert.Alert) error {
	rep := Report{
		GeneratedAt: time.Now().UTC(),
		TotalAlerts: len(alerts),
		Alerts:      alerts,
	}

	switch r.format {
	case FormatJSON:
		return r.writeJSON(rep)
	default:
		return r.writeText(rep)
	}
}

func (r *Reporter) writeJSON(rep Report) error {
	enc := json.NewEncoder(r.out)
	enc.SetIndent("", "  ")
	return enc.Encode(rep)
}

func (r *Reporter) writeText(rep Report) error {
	if _, err := fmt.Fprintf(r.out, "depwatch report — %s\n", rep.GeneratedAt.Format(time.RFC3339)); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(r.out, "Total alerts: %d\n\n", rep.TotalAlerts); err != nil {
		return err
	}
	if rep.TotalAlerts == 0 {
		_, err := fmt.Fprintln(r.out, "No issues found.")
		return err
	}
	for _, a := range rep.Alerts {
		if _, err := fmt.Fprintf(r.out, "[%s] %s @ %s — %s\n", a.Severity, a.Package, a.Version, a.Message); err != nil {
			return err
		}
	}
	return nil
}
