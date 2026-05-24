package dedupe_test

import (
	"testing"

	"github.com/user/depwatch/internal/alert"
	"github.com/user/depwatch/internal/dedupe"
)

func makeAlert(ecosystem, pkg, severity string) alert.Alert {
	return alert.Alert{
		Ecosystem: ecosystem,
		Package:   pkg,
		Severity:  severity,
		Message:   "test message",
	}
}

func TestByKey_NoDuplicates_ReturnsSameLength(t *testing.T) {
	alerts := []alert.Alert{
		makeAlert("npm", "lodash", "warn"),
		makeAlert("npm", "express", "critical"),
		makeAlert("pypi", "requests", "warn"),
	}

	result := dedupe.ByKey(alerts)
	if len(result) != 3 {
		t.Fatalf("expected 3 alerts, got %d", len(result))
	}
}

func TestByKey_WithDuplicates_RemovesExtras(t *testing.T) {
	alerts := []alert.Alert{
		makeAlert("npm", "lodash", "warn"),
		makeAlert("npm", "lodash", "warn"), // duplicate
		makeAlert("npm", "express", "critical"),
	}

	result := dedupe.ByKey(alerts)
	if len(result) != 2 {
		t.Fatalf("expected 2 alerts after dedup, got %d", len(result))
	}
}

func TestByKey_PreservesOrder(t *testing.T) {
	alerts := []alert.Alert{
		makeAlert("npm", "lodash", "warn"),
		makeAlert("pypi", "requests", "critical"),
		makeAlert("npm", "lodash", "warn"), // duplicate of first
	}

	result := dedupe.ByKey(alerts)
	if result[0].Package != "lodash" || result[1].Package != "requests" {
		t.Errorf("order not preserved: got %v", result)
	}
}

func TestByKey_EmptySlice_ReturnsEmpty(t *testing.T) {
	result := dedupe.ByKey(nil)
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d", len(result))
	}
}

func TestCount_ReportsDuplicateCount(t *testing.T) {
	alerts := []alert.Alert{
		makeAlert("npm", "lodash", "warn"),
		makeAlert("npm", "lodash", "warn"),
		makeAlert("npm", "lodash", "warn"),
		makeAlert("pypi", "requests", "critical"),
	}

	count := dedupe.Count(alerts)
	if count != 2 {
		t.Errorf("expected 2 duplicates, got %d", count)
	}
}

func TestCount_NoDuplicates_ReturnsZero(t *testing.T) {
	alerts := []alert.Alert{
		makeAlert("npm", "lodash", "warn"),
		makeAlert("pypi", "requests", "critical"),
	}

	if dedupe.Count(alerts) != 0 {
		t.Error("expected zero duplicates")
	}
}
