package history

// Direction indicates whether the alert count improved, worsened, or stayed
// the same between two consecutive scans.
type Direction string

const (
	Improved  Direction = "improved"
	Worsened  Direction = "worsened"
	Unchanged Direction = "unchanged"
)

// Trend summarises the change between the two most recent records.
type Trend struct {
	Previous  int
	Current   int
	Delta     int
	Direction Direction
}

// Compare derives a Trend by comparing the two most recent records in h.
// It returns false when there are fewer than two records.
func Compare(h *History) (Trend, bool) {
	if h.Len() < 2 {
		return Trend{}, false
	}
	prev := len(h.Records[h.Len()-2].Alerts)
	curr := len(h.Records[h.Len()-1].Alerts)
	delta := curr - prev

	var dir Direction
	switch {
	case delta < 0:
		dir = Improved
	case delta > 0:
		dir = Worsened
	default:
		dir = Unchanged
	}

	return Trend{
		Previous:  prev,
		Current:   curr,
		Delta:     delta,
		Direction: dir,
	}, true
}
