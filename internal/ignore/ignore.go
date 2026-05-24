// Package ignore provides functionality to skip dependencies
// that have been explicitly acknowledged or suppressed by the user.
package ignore

import (
	"encoding/json"
	"os"
	"time"
)

// Entry represents a single suppressed dependency.
type Entry struct {
	Ecosystem string    `json:"ecosystem"`
	Package   string    `json:"package"`
	Reason    string    `json:"reason"`
	Expires   time.Time `json:"expires,omitempty"`
}

// List holds all ignored dependency entries.
type List struct {
	Entries []Entry `json:"ignore"`
}

// Load reads an ignore file from the given path.
// If the path is empty or the file does not exist, an empty List is returned.
func Load(path string) (*List, error) {
	if path == "" {
		return &List{}, nil
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &List{}, nil
	}
	if err != nil {
		return nil, err
	}
	var l List
	if err := json.Unmarshal(data, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// Contains reports whether the given ecosystem+package combination
// is present in the ignore list and has not yet expired.
func (l *List) Contains(ecosystem, pkg string) bool {
	now := time.Now()
	for _, e := range l.Entries {
		if e.Ecosystem == ecosystem && e.Package == pkg {
			if e.Expires.IsZero() || e.Expires.After(now) {
				return true
			}
		}
	}
	return false
}

// Active returns only the entries that have not expired.
func (l *List) Active() []Entry {
	now := time.Now()
	out := make([]Entry, 0, len(l.Entries))
	for _, e := range l.Entries {
		if e.Expires.IsZero() || e.Expires.After(now) {
			out = append(out, e)
		}
	}
	return out
}
