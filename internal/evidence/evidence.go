// Package evidence runs a knowledge base's evidence programs and returns their
// results, so claims can be verified by re-execution rather than trust.
//
// It is runner-agnostic and incremental, in the spirit of `go test`:
//
//   - A content directory may hold an `evidence/` tree. Each directory under it
//     that contains Go files is a *unit* — a `main` package that computes some
//     facts and prints them as JSON. A unit may be the `evidence/` dir itself,
//     or a sub-package per expensive check (e.g. `evidence/orbit-sim/`).
//   - A unit prints either a single result, `{"value": "..."}` /
//     `{"table": [...]}`, named after its directory, or a map of named results,
//     `{"midwest-regions": {"value": "4"}, ...}`.
//   - starbase caches each unit's result keyed by a content hash of its Go
//     sources plus any data files it declares with a `//starbase:deps <glob>`
//     comment. A unit is re-run only when that hash changes — so editing an
//     unrelated page never re-runs a minutes-long simulation. The cache lives in
//     the user cache dir; CI starts cold, so CI is always authoritative.
//
// starbase ships no language runners: the unit's Go does the work (pure Go, or
// shelling out to DuckDB, a driver, an API). starbase builds, runs, caches, and
// hands the results to `verify` for comparison.
package evidence

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Check is one computed result.
type Check struct {
	Value string     `json:"value,omitempty"`
	Table [][]string `json:"table,omitempty"`
	Error string     `json:"error,omitempty"`
}

// Unit is one evidence package and how it was resolved this run.
type Unit struct {
	Name   string // default result name (its directory)
	Dir    string
	Cached bool
	Err    string
}

// Result is the merged output of every unit plus per-unit status.
type Result struct {
	Checks map[string]Check
	Units  []Unit
}

const perUnitTimeout = 20 * time.Minute

// Run resolves every evidence unit, reusing cached results whose inputs are
// unchanged and re-running the rest. The bool reports whether evidence exists.
func Run(contentDir string, force bool) (Result, bool, error) {
	evDir := filepath.Join(contentDir, "evidence")
	if fi, err := os.Stat(evDir); err != nil || !fi.IsDir() {
		return Result{}, false, nil
	}
	units := discoverUnits(evDir)
	cache := loadCache(contentDir)
	res := Result{Checks: map[string]Check{}}

	for _, u := range units {
		key, err := unitKey(u.Dir)
		if err != nil {
			u.Err = err.Error()
			res.Units = append(res.Units, u)
			continue
		}
		if !force {
			if c, ok := cache[u.Name]; ok && c.Key == key {
				u.Cached = true
				merge(res.Checks, c.Checks, u.Name)
				res.Units = append(res.Units, u)
				continue
			}
		}
		checks, runErr := runUnit(u.Dir)
		if runErr != "" {
			u.Err = runErr
			res.Units = append(res.Units, u)
			continue
		}
		cache[u.Name] = cachedUnit{Key: key, Checks: checks}
		merge(res.Checks, checks, u.Name)
		res.Units = append(res.Units, u)
	}
	saveCache(contentDir, cache)
	return res, true, nil
}

func discoverUnits(evDir string) []Unit {
	var units []Unit
	if hasGo(evDir) {
		units = append(units, Unit{Name: "evidence", Dir: evDir})
	}
	entries, _ := os.ReadDir(evDir)
	for _, e := range entries {
		if e.IsDir() {
			d := filepath.Join(evDir, e.Name())
			if hasGo(d) {
				units = append(units, Unit{Name: e.Name(), Dir: d})
			}
		}
	}
	return units
}

func hasGo(dir string) bool {
	m, _ := filepath.Glob(filepath.Join(dir, "*.go"))
	return len(m) > 0
}

var reDeps = regexp.MustCompile(`(?m)^\s*//\s*starbase:deps\s+(.+)$`)

// unitKey hashes a unit's Go sources, go.mod/go.sum, and declared data deps.
func unitKey(dir string) (string, error) {
	goFiles, _ := filepath.Glob(filepath.Join(dir, "*.go"))
	sort.Strings(goFiles)
	h := sha256.New()
	depSet := map[string]bool{}
	for _, gf := range goFiles {
		b, err := os.ReadFile(gf)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(h, "go:%s\n", filepath.Base(gf))
		h.Write(b)
		for _, m := range reDeps.FindAllStringSubmatch(string(b), -1) {
			for _, pat := range strings.FieldsFunc(m[1], func(r rune) bool { return r == ',' || r == ' ' || r == '\t' }) {
				matches, _ := filepath.Glob(filepath.Join(dir, pat))
				for _, mt := range matches {
					depSet[mt] = true
				}
			}
		}
	}
	for _, extra := range []string{"go.mod", "go.sum"} {
		if b, err := os.ReadFile(filepath.Join(dir, extra)); err == nil {
			fmt.Fprintf(h, "%s\n", extra)
			h.Write(b)
		}
	}
	deps := make([]string, 0, len(depSet))
	for d := range depSet {
		deps = append(deps, d)
	}
	sort.Strings(deps)
	for _, dep := range deps {
		b, err := os.ReadFile(dep)
		if err != nil {
			continue
		}
		rel, _ := filepath.Rel(dir, dep)
		fmt.Fprintf(h, "dep:%s\n", rel)
		h.Write(b)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func runUnit(dir string) (map[string]Check, string) {
	ctx, cancel := context.WithTimeout(context.Background(), perUnitTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "run", ".")
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return nil, msg
	}
	checks, err := parseOutput(stdout.Bytes())
	if err != nil {
		return nil, "did not print a JSON result: " + err.Error()
	}
	return checks, ""
}

// parseOutput accepts either a single Check (keyed by ""), or a map of named
// Checks. The "" name is later replaced by the unit's directory name.
func parseOutput(b []byte) (map[string]Check, error) {
	b = bytes.TrimSpace(b)
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, err
	}
	single := len(raw) > 0
	for k := range raw {
		if k != "value" && k != "table" && k != "error" {
			single = false
			break
		}
	}
	if single {
		var c Check
		if err := json.Unmarshal(b, &c); err != nil {
			return nil, err
		}
		return map[string]Check{"": c}, nil
	}
	out := map[string]Check{}
	for k, v := range raw {
		var c Check
		if json.Unmarshal(v, &c) == nil {
			out[k] = c
		}
	}
	return out, nil
}

func merge(dst, src map[string]Check, unitName string) {
	for k, v := range src {
		name := k
		if name == "" {
			name = unitName
		}
		dst[name] = v
	}
}

// --- persisted cache (user cache dir; CI starts cold) ---

type cachedUnit struct {
	Key    string           `json:"key"`
	Checks map[string]Check `json:"checks"`
}

func cachePath(contentDir string) string {
	base, err := os.UserCacheDir()
	if err != nil {
		base = os.TempDir()
	}
	abs, _ := filepath.Abs(contentDir)
	sum := sha256.Sum256([]byte(abs))
	return filepath.Join(base, "starbase", "evidence", hex.EncodeToString(sum[:])[:16]+".json")
}

func loadCache(contentDir string) map[string]cachedUnit {
	b, err := os.ReadFile(cachePath(contentDir))
	if err != nil {
		return map[string]cachedUnit{}
	}
	var m map[string]cachedUnit
	if json.Unmarshal(b, &m) != nil || m == nil {
		return map[string]cachedUnit{}
	}
	return m
}

func saveCache(contentDir string, m map[string]cachedUnit) {
	p := cachePath(contentDir)
	if os.MkdirAll(filepath.Dir(p), 0o755) != nil {
		return
	}
	if b, err := json.MarshalIndent(m, "", "  "); err == nil {
		_ = os.WriteFile(p, b, 0o644)
	}
}
