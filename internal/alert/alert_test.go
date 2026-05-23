package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/depwatch/internal/alert"
	"github.com/user/depwatch/internal/checker"
)

func TestFromResults_FiltersClean(t *testing.T) {
	results := []checker.Result{
		{Package: "lodash", Current: "4.17.21", Latest: "4.17.21", Outdated: false, Vulnerable: false},
	}
	alerts := alert.FromResults(results)
	if len(alerts) != 0 {
		t.Fatalf("expected 0 alerts for clean deps, got %d", len(alerts))
	}
}

func TestFromResults_OutdatedIsWarn(t *testing.T) {
	results := []checker.Result{
		{Package: "express", Current: "4.17.0", Latest: "4.18.2", Outdated: true, Vulnerable: false},
	}
	alerts := alert.FromResults(results)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Level != alert.LevelWarn {
		t.Errorf("expected WARN level, got %s", alerts[0].Level)
	}
}

func TestFromResults_VulnerableIsCritical(t *testing.T) {
	results := []checker.Result{
		{Package: "axios", Current: "0.21.0", Latest: "1.4.0", Outdated: true, Vulnerable: true},
	}
	alerts := alert.FromResults(results)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Level != alert.LevelCritical {
		t.Errorf("expected CRITICAL level, got %s", alerts[0].Level)
	}
}

func TestAlerter_PrintNoIssues(t *testing.T) {
	var buf bytes.Buffer
	a := alert.New(&buf)
	a.Print([]alert.Alert{})
	if !strings.Contains(buf.String(), "up to date") {
		t.Errorf("expected 'up to date' message, got: %s", buf.String())
	}
}

func TestAlerter_PrintIssues(t *testing.T) {
	var buf bytes.Buffer
	a := alert.New(&buf)
	alerts := []alert.Alert{
		{Package: "react", Current: "17.0.0", Latest: "18.2.0", Level: alert.LevelWarn, Message: "outdated package react@17.0.0 — latest is 18.2.0"},
	}
	a.Print(alerts)
	out := buf.String()
	if !strings.Contains(out, "1 issue(s)") {
		t.Errorf("expected issue count in output, got: %s", out)
	}
	if !strings.Contains(out, "[WARN]") {
		t.Errorf("expected WARN tag in output, got: %s", out)
	}
}
