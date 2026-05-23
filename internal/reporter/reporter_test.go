package reporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/depwatch/internal/alert"
	"github.com/depwatch/internal/reporter"
)

func sampleAlerts() []alert.Alert {
	return []alert.Alert{
		{Package: "lodash", Version: "4.17.11", Severity: "critical", Message: "prototype pollution"},
		{Package: "express", Version: "4.17.1", Severity: "warn", Message: "outdated package"},
	}
}

func TestWrite_TextFormat_NoAlerts(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.NewWithWriter(reporter.FormatText, &buf)
	if err := r.Write(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No issues found.") {
		t.Errorf("expected 'No issues found.' in output, got: %s", buf.String())
	}
}

func TestWrite_TextFormat_WithAlerts(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.NewWithWriter(reporter.FormatText, &buf)
	if err := r.Write(sampleAlerts()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "lodash") {
		t.Errorf("expected 'lodash' in output")
	}
	if !strings.Contains(out, "Total alerts: 2") {
		t.Errorf("expected total alerts count in output")
	}
}

func TestWrite_JSONFormat_Structure(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.NewWithWriter(reporter.FormatJSON, &buf)
	if err := r.Write(sampleAlerts()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var rep reporter.Report
	if err := json.Unmarshal(buf.Bytes(), &rep); err != nil {
		t.Fatalf("failed to parse JSON output: %v", err)
	}
	if rep.TotalAlerts != 2 {
		t.Errorf("expected TotalAlerts=2, got %d", rep.TotalAlerts)
	}
	if len(rep.Alerts) != 2 {
		t.Errorf("expected 2 alerts in JSON, got %d", len(rep.Alerts))
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	r := reporter.New(reporter.FormatText)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}
