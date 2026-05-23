package filter_test

import (
	"testing"

	"github.com/user/depwatch/internal/alert"
	"github.com/user/depwatch/internal/filter"
)

func sampleAlerts() []alert.Alert {
	return []alert.Alert{
		{Package: "lodash", Ecosystem: "npm", Level: alert.LevelWarn, Message: "outdated"},
		{Package: "axios", Ecosystem: "npm", Level: alert.LevelCritical, Message: "vulnerable"},
		{Package: "requests", Ecosystem: "pypi", Level: alert.LevelWarn, Message: "outdated"},
		{Package: "flask", Ecosystem: "pypi", Level: alert.LevelCritical, Message: "vulnerable"},
	}
}

func TestApply_AllSeverity_AllEcosystems(t *testing.T) {
	result := filter.Apply(sampleAlerts(), filter.Options{MinSeverity: filter.SeverityAll})
	if len(result) != 4 {
		t.Fatalf("expected 4 alerts, got %d", len(result))
	}
}

func TestApply_WarnSeverity(t *testing.T) {
	result := filter.Apply(sampleAlerts(), filter.Options{MinSeverity: filter.SeverityWarn})
	if len(result) != 4 {
		t.Fatalf("expected 4 alerts (warn+critical), got %d", len(result))
	}
}

func TestApply_CriticalOnly(t *testing.T) {
	result := filter.Apply(sampleAlerts(), filter.Options{MinSeverity: filter.SeverityCritical})
	if len(result) != 2 {
		t.Fatalf("expected 2 critical alerts, got %d", len(result))
	}
	for _, a := range result {
		if a.Level != alert.LevelCritical {
			t.Errorf("expected critical, got %v", a.Level)
		}
	}
}

func TestApply_EcosystemFilter(t *testing.T) {
	result := filter.Apply(sampleAlerts(), filter.Options{Ecosystem: "npm"})
	if len(result) != 2 {
		t.Fatalf("expected 2 npm alerts, got %d", len(result))
	}
}

func TestApply_EcosystemAndSeverity(t *testing.T) {
	result := filter.Apply(sampleAlerts(), filter.Options{
		MinSeverity: filter.SeverityCritical,
		Ecosystem:   "pypi",
	})
	if len(result) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(result))
	}
	if result[0].Package != "flask" {
		t.Errorf("expected flask, got %s", result[0].Package)
	}
}

func TestParseSeverity(t *testing.T) {
	cases := []struct {
		input string
		ok    bool
		want  filter.Severity
	}{
		{"all", true, filter.SeverityAll},
		{"", true, filter.SeverityAll},
		{"warn", true, filter.SeverityWarn},
		{"critical", true, filter.SeverityCritical},
		{"unknown", false, filter.SeverityAll},
	}
	for _, c := range cases {
		got, ok := filter.ParseSeverity(c.input)
		if ok != c.ok || got != c.want {
			t.Errorf("ParseSeverity(%q) = (%v, %v), want (%v, %v)", c.input, got, ok, c.want, c.ok)
		}
	}
}
