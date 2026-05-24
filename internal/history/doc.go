// Package history provides persistence and trend analysis for depwatch scan
// results. Each time depwatch runs it can append the produced alerts to a
// JSON file; the Compare function then surfaces whether the dependency health
// of the repository is improving or worsening over time.
//
// Typical usage:
//
//	h, err := history.New(".depwatch-history.json")
//	if err != nil { /* handle */ }
//
//	if err := h.Append(alerts); err != nil { /* handle */ }
//
//	if trend, ok := history.Compare(h); ok {
//		fmt.Println(trend.Direction)
//	}
package history
