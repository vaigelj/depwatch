// Package history records scan results over time so that trends and
// regressions can be detected across successive depwatch runs.
package history

import (
	"encoding/json"
	"os"
	"time"

	"github.com/depwatch/internal/alert"
)

// Record holds a snapshot of alerts produced during a single scan.
type Record struct {
	Timestamp time.Time    `json:"timestamp"`
	Alerts    []alert.Alert `json:"alerts"`
}

// History is an ordered list of scan records, newest last.
type History struct {
	Records []Record `json:"records"`
	path    string
}

// New loads existing history from path, or returns an empty History when the
// file does not yet exist.
func New(path string) (*History, error) {
	h := &History{path: path}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return h, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, h); err != nil {
		return nil, err
	}
	return h, nil
}

// Append adds a new record for the given alerts and persists the history.
func (h *History) Append(alerts []alert.Alert) error {
	h.Records = append(h.Records, Record{
		Timestamp: time.Now().UTC(),
		Alerts:    alerts,
	})
	return h.save()
}

// Latest returns the most recent record, or false when history is empty.
func (h *History) Latest() (Record, bool) {
	if len(h.Records) == 0 {
		return Record{}, false
	}
	return h.Records[len(h.Records)-1], true
}

// Len returns the number of stored records.
func (h *History) Len() int { return len(h.Records) }

func (h *History) save() error {
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.path, data, 0o644)
}
