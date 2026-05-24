package history_test

import (
	"testing"

	"github.com/depwatch/internal/history"
)

func buildHistory(t *testing.T, counts []int) *history.History {
	t.Helper()
	h, err := history.New(filepath.Join(t.TempDir(), "h.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, n := range counts {
		if err := h.Append(sampleAlerts(n)); err != nil {
			t.Fatal(err)
		}
	}
	return h
}

func TestCompare_Worsened(t *testing.T) {
	h := buildHistory(t, []int{2, 5})
	trend, ok := history.Compare(h)
	if !ok {
		t.Fatal("expected trend to be available")
	}
	if trend.Direction != history.Worsened {
		t.Errorf("expected Worsened, got %s", trend.Direction)
	}
	if trend.Delta != 3 {
		t.Errorf("expected delta 3, got %d", trend.Delta)
	}
}

func TestCompare_Improved(t *testing.T) {
	h := buildHistory(t, []int{5, 2})
	trend, ok := history.Compare(h)
	if !ok {
		t.Fatal("expected trend")
	}
	if trend.Direction != history.Improved {
		t.Errorf("expected Improved, got %s", trend.Direction)
	}
}

func TestCompare_Unchanged(t *testing.T) {
	h := buildHistory(t, []int{3, 3})
	trend, _ := history.Compare(h)
	if trend.Direction != history.Unchanged {
		t.Errorf("expected Unchanged, got %s", trend.Direction)
	}
}

func TestCompare_SingleRecord_ReturnsFalse(t *testing.T) {
	h := buildHistory(t, []int{4})
	_, ok := history.Compare(h)
	if ok {
		t.Error("expected false with only one record")
	}
}

func TestCompare_EmptyHistory_ReturnsFalse(t *testing.T) {
	h := buildHistory(t, nil)
	_, ok := history.Compare(h)
	if ok {
		t.Error("expected false for empty history")
	}
}
