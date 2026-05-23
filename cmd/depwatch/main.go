package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/depwatch/internal/alert"
	"github.com/depwatch/internal/checker"
	"github.com/depwatch/internal/reporter"
	"github.com/depwatch/internal/scanner"
)

func main() {
	dir := flag.String("dir", ".", "directory to scan")
	fmt := flag.String("format", "text", "output format: text or json")
	flag.Parse()

	if err := run(*dir, *fmt); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(dir, format string) error {
	s := scanner.New()
	deps, err := s.Scan(dir)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	c := checker.New()
	results, err := c.Check(deps)
	if err != nil {
		return fmt.Errorf("check failed: %w", err)
	}

	alerts := alert.FromResults(results)

	fmt := reporter.FormatText
	if format == "json" {
		fmt = reporter.FormatJSON
	}

	r := reporter.New(fmt)
	return r.Write(alerts)
}
