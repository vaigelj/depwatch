// Package summary provides utilities for aggregating the results of a
// depwatch scan into a concise Summary value.
//
// A Summary records the total number of dependencies examined alongside
// per-severity counts (clean, warn, critical) and the set of ecosystems
// encountered. It is intended to be embedded in reports or printed to the
// terminal as a final status line after all alerts have been emitted.
//
// Example usage:
//
//	s := summary.Build(alerts, scanner.TotalDeps())
//	fmt.Println(s)
package summary
