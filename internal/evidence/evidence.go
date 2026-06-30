// Package evidence runs a knowledge base's evidence program and returns its
// results, so claims can be verified by *re-execution* rather than trust.
//
// The framework is deliberately runner-agnostic: a content directory may contain
// an `evidence/` Go module (package main) that computes whatever it likes — query
// a database, shell out to DuckDB, read a CSV, call an API — and prints a JSON
// object of results to stdout:
//
//	{ "midwest-regions": { "value": "4" },
//	  "revenue-by-division": { "table": [["division","total"],["Midwest","11400000"]] } }
//
// starbase compiles and runs that program (like `go test` builds and runs tests)
// and diffs each result against what the article claims. The build is the trust
// anchor: the author can't fake a number the build recomputes.
package evidence

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Check is one computed result emitted by the evidence program.
type Check struct {
	Value string     `json:"value,omitempty"`
	Table [][]string `json:"table,omitempty"`
	Error string     `json:"error,omitempty"`
}

// Run compiles and executes <contentDir>/evidence (a Go main package) and parses
// its stdout. The bool reports whether an evidence program is present at all.
func Run(contentDir string, timeout time.Duration) (map[string]Check, bool, error) {
	dir := filepath.Join(contentDir, "evidence")
	if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
		return nil, false, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", ".")
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		msg := stderr.String()
		if msg == "" {
			msg = err.Error()
		}
		return nil, true, fmt.Errorf("evidence program failed:\n%s", msg)
	}

	var checks map[string]Check
	if err := json.Unmarshal(bytes.TrimSpace(stdout.Bytes()), &checks); err != nil {
		return nil, true, fmt.Errorf("evidence program did not print a JSON object of results: %w", err)
	}
	return checks, true, nil
}
