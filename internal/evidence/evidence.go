// Package evidence runs a knowledge base's checks and returns their output, so
// claims can be verified by re-execution rather than trust.
//
// A check is a directory under `evidence/` containing an executable `run`:
//
//	evidence/midwest-regions/run   (chmod +x, any #! interpreter)
//
// The check's result is whatever `run` prints to stdout; a non-zero exit is an
// error (its stderr is surfaced). The program can be anything — a shell one-liner
// over DuckDB, a Python script, a Go program — because the contract is just
// (stdout, exit code), like a golden test. `verify` compares that stdout, trimmed,
// against the result a claim embeds, and fails the build when they disagree.
//
// It is incremental: each check is a cache unit keyed by a hash of its `run`
// script plus the files it declares with `starbase:inputs`. A check re-runs only
// when its script or a declared input changes. For inputs a static file list
// can't name (a URL, a database, a clock), a check may add an executable `stamp`
// whose stdout is a cheap fingerprint that also drives re-runs. The cache lives
// in the user cache dir, so CI starts cold and is authoritative.
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

// Check is one check's result: the stdout of its `run`, or an error if `run`
// exited non-zero.
type Check struct {
	Output string `json:"output"`
	Err    string `json:"err,omitempty"`
}

// Unit reports how one check was resolved this run.
type Unit struct {
	Name   string
	Cached bool
	Err    string
}

// Result is every check's output plus per-check status.
type Result struct {
	Checks map[string]Check
	Units  []Unit
}

const perRunTimeout = 20 * time.Minute

// Run resolves every check under evidence/, reusing cached output whose inputs
// are unchanged and re-running the rest. The bool reports whether evidence exists.
//
// A check re-runs when its `run` script or a declared starbase:inputs file
// changed (the static key, computed without executing anything), or — for checks
// with a `stamp` — when their stamp differs from the cached one. Checks with a
// stamp always pay for the cheap stamp; only those whose stamp (or static key)
// moved pay for the expensive run.
func Run(contentDir string, force bool) (Result, bool, error) {
	evDir := filepath.Join(contentDir, "evidence")
	if fi, err := os.Stat(evDir); err != nil || !fi.IsDir() {
		return Result{}, false, nil
	}
	absRoot, _ := filepath.Abs(contentDir)
	entries, err := os.ReadDir(evDir)
	if err != nil {
		return Result{}, true, err
	}

	cache := loadCache(contentDir)
	res := Result{Checks: map[string]Check{}}

	useCached := func(name string, c cachedUnit) {
		res.Checks[name] = c.Check
		res.Units = append(res.Units, Unit{Name: name, Cached: true})
	}

	for _, e := range entries {
		name := e.Name()
		if !e.IsDir() || strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") {
			continue
		}
		runPath := filepath.Join(evDir, name, "run")
		if !isExecutable(runPath) {
			continue // a directory without an executable run/ is not a check
		}
		staticKey := unitKey(runPath, contentDir)
		stampPath := filepath.Join(evDir, name, "stamp")
		hasStamp := isExecutable(stampPath)

		c, cached := cache[name]
		staticFresh := cached && c.StaticKey == staticKey
		if !force && staticFresh && !hasStamp {
			useCached(name, c) // nothing to execute
			continue
		}

		stamp := ""
		if hasStamp {
			out, serr := execScript(stampPath, absRoot)
			if serr != "" {
				stamp = "stamp error: " + serr // forces a (re)run and differs from any prior stamp
			} else {
				stamp = strings.TrimSpace(out)
			}
		}
		if !force && staticFresh && hasStamp && stamp == c.Stamp {
			useCached(name, c) // stamp matched: skip the expensive run
			continue
		}

		out, rerr := execScript(runPath, absRoot)
		ck := Check{Output: out, Err: rerr}
		cache[name] = cachedUnit{StaticKey: staticKey, Stamp: stamp, Check: ck}
		res.Checks[name] = ck
		res.Units = append(res.Units, Unit{Name: name, Err: ck.Err})
	}
	saveCache(contentDir, cache)
	return res, true, nil
}

func isExecutable(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.Mode().IsRegular() && fi.Mode().Perm()&0o111 != 0
}

// execScript runs an executable with the knowledge-base root as its working
// directory (so it references data as "data/sales.csv"). It returns stdout and,
// on a non-zero exit, the stderr (or the error) — empty stderr means success.
func execScript(path, cwd string) (stdout, errMsg string) {
	ctx, cancel := context.WithTimeout(context.Background(), perRunTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	cmd.Dir = cwd
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(errb.String())
		if msg == "" {
			msg = err.Error()
		}
		return out.String(), msg
	}
	return out.String(), ""
}

// --- cache key ---

// reInputs matches a declaration line in any comment style, e.g.
//
//	# starbase:inputs data/sales.csv other/*.csv
//	// starbase:inputs data/sales.csv
var reInputs = regexp.MustCompile(`(?m)^.*?starbase:inputs[ \t]+(.+)$`)

func unitKey(runPath, contentDir string) string {
	h := sha256.New()
	b, err := os.ReadFile(runPath)
	if err != nil {
		return ""
	}
	h.Write(b)

	depSet := map[string]bool{}
	for _, m := range reInputs.FindAllStringSubmatch(string(b), -1) {
		for _, pat := range strings.FieldsFunc(m[1], func(r rune) bool { return r == ',' || r == ' ' || r == '\t' }) {
			// inputs are relative to the KB root (the run script's CWD)
			matches, _ := filepath.Glob(filepath.Join(contentDir, pat))
			for _, mt := range matches {
				depSet[mt] = true
			}
		}
	}
	deps := make([]string, 0, len(depSet))
	for d := range depSet {
		deps = append(deps, d)
	}
	sort.Strings(deps)
	for _, dep := range deps {
		fb, err := os.ReadFile(dep)
		if err != nil {
			continue
		}
		rel, _ := filepath.Rel(contentDir, dep)
		fmt.Fprintf(h, "\x00input:%s\x00", rel)
		h.Write(fb)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// --- persisted cache (user cache dir; CI starts cold) ---

type cachedUnit struct {
	StaticKey string `json:"static_key"` // hash of run script + declared inputs
	Stamp     string `json:"stamp,omitempty"`
	Check     Check  `json:"check"`
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
