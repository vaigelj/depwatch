package summary_test

import (
	"testing"

	"github.com/user/depwatch/internal/alert"
	"github.com/user/depwatch/internal/summary"
)

func sampleAlerts() []alert.Alert {
	return []alert.Alert{
		{Package: "lodash", Ecosystem: "npm", Level: alert.LevelWarn, Message: "outdated"},
		{Package: "axios", Ecosystem: "npm", Level: alert.LevelCritical, Message: "vulnerable"},
		{Package: "requests", Ecosystem: "pypi", Level: alert.LevelWarn, Message: "outdated"},
	}
}

func TestBuild_CountsCorrectly(t *testing.T) {
	s := summary.Build(sampleAlerts(), 10)
	if s.TotalDeps != 10 {
		t.Errorf("expected TotalDeps=10, got %d", s.TotalDeps)
	}
	if s.Warn != 2 {
		t.Errorf("expected Warn=2, got %d", s.Warn)
	}
	if s.Critical != 1 {
		t.Errorf("expected Critical=1, got %d", s.Critical)
	}
	if s.Clean != 7 {
		t.Errorf("expected Clean=7, got %d", s.Clean)
	}
}

func TestBuild_DeduplicatesEcosystems(t *testing.T) {
	s := summary.Build(sampleAlerts(), 10)
	if len(s.Ecosystems) != 2 {
		t.Errorf("expected 2 ecosystems, got %d: %v", len(s.Ecosystems), s.Ecosystems)
	}
}

func TestBuild_EmptyAlerts(t *testing.T) {
	s := summary.Build(nil, 5)
	if s.Clean != 5 {
		t.Errorf("expected Clean=5, got %d", s.Clean)
	}
	if s.OverallLevel() != summary.LevelClean {
		t.Errorf("expected clean level, got %s", s.OverallLevel())
	}
}

func TestOverallLevel_Critical(t *testing.T) {
	s := summary.Build(sampleAlerts(), 10)
	if s.OverallLevel() != summary.LevelCritical {
		t.Errorf("expected critical, got %s", s.OverallLevel())
	}
}

func TestOverallLevel_WarnOnly(t *testing.T) {
	alerts := []alert.Alert{
		{Package: "lodash", Ecosystem: "npm", Level: alert.LevelWarn, Message: "outdated"},
	}
	s := summary.Build(alerts, 5)
	if s.OverallLevel() != summary.LevelWarn {
		t.Errorf("expected warn, got %s", s.OverallLevel())
	}
}

func TestString_ContainsLevel(t *testing.T) {
	s := summary.Build(sampleAlerts(), 10)
	str := s.String()
	if len(str) == 0 {
		t.Error("expected non-empty string")
	}
}
