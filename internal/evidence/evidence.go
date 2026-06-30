// Package evidence runs a knowledge base's checks and returns their output, so
// claims can be verified by re-execution rather than trust.
//
// A check is a directory under `evidence/` containing an executable `run` and an
// optional `inputs` manifest:
//
//	evidence/midwest-regions/run      (chmod +x, any #! interpreter)
//	evidence/midwest-regions/inputs   data/sales.csv
//
// Each line of `inputs` names a source resolved by a provider — a local file or
// an http(s) URL — and staged into a fresh working directory under its basename
// (override with `source -> localname`). `run` executes there with that directory
// as its CWD, so it reads inputs by name (`sales.csv`), never reaching into the
// repo. Its stdout is the result; a non-zero exit is an error (stderr surfaced).
// The program can be anything — a shell one-liner over DuckDB, Python, a binary —
// because the contract is just (stdout, exit code), like a golden test. `verify`
// compares that stdout, trimmed, against the result a claim embeds.
//
// It is incremental: each check is a cache unit keyed by a hash of its `run`
// script, its `inputs` manifest, and the resolved content of every input. A check
// re-runs only when one of those changes. The cache lives in the user cache dir,
// so CI starts cold and is authoritative.
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
	"sort"
	"strings"
	"time"
)

// Check is one check's result: the stdout of its `run`, or an error if `run`
// exited non-zero (or an input could not be resolved).
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
func Run(contentDir string, force bool) (Result, bool, error) {
	evDir := filepath.Join(contentDir, "evidence")
	if fi, err := os.Stat(evDir); err != nil || !fi.IsDir() {
		return Result{}, false, nil
	}
	entries, err := os.ReadDir(evDir)
	if err != nil {
		return Result{}, true, err
	}

	cache := loadCache(contentDir)
	res := Result{Checks: map[string]Check{}}
	record := func(name string, ck Check, cached bool) {
		res.Checks[name] = ck
		res.Units = append(res.Units, Unit{Name: name, Cached: cached, Err: ck.Err})
	}

	for _, e := range entries {
		name := e.Name()
		if !e.IsDir() || strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") {
			continue
		}
		dir := filepath.Join(evDir, name)
		runPath := filepath.Join(dir, "run")
		if !isExecutable(runPath) {
			continue // a directory without an executable run is not a check
		}
		runBytes, err := os.ReadFile(runPath)
		if err != nil {
			record(name, Check{Err: err.Error()}, false)
			continue
		}
		manifest, _ := os.ReadFile(filepath.Join(dir, "inputs"))

		// Resolve inputs (the http provider may fetch). A resolution failure —
		// missing file, 404 — is the check's error.
		inputs, rerr := resolveInputs(parseInputs(manifest), contentDir, force)
		if rerr != "" {
			record(name, Check{Err: rerr}, false)
			continue
		}

		key := computeKey(runBytes, manifest, inputs)
		if c, ok := cache[name]; ok && c.Key == key && !force {
			record(name, c.Check, true)
			continue
		}

		out, runErr := runStaged(runPath, inputs)
		ck := Check{Output: out, Err: runErr}
		cache[name] = cachedUnit{Key: key, Check: ck}
		record(name, ck, false)
	}
	saveCache(contentDir, cache)
	return res, true, nil
}

// parseInputs reads an `inputs` manifest: one source per line, comments with #,
// optional rename via `source -> localname` (or `-->`).
func parseInputs(b []byte) []inputSpec {
	var specs []inputSpec
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		src, name := line, ""
		for _, sep := range []string{"-->", "->"} {
			if i := strings.Index(line, sep); i >= 0 {
				src = strings.TrimSpace(line[:i])
				name = strings.TrimSpace(line[i+len(sep):])
				break
			}
		}
		specs = append(specs, inputSpec{Source: src, LocalName: name})
	}
	return specs
}

func resolveInputs(specs []inputSpec, contentDir string, force bool) ([]resolved, string) {
	var out []resolved
	for _, s := range specs {
		r, err := providerFor(s.Source).resolve(s, contentDir, force)
		if err != nil {
			return nil, fmt.Sprintf("input %q: %s", s.Source, err)
		}
		out = append(out, r)
	}
	return out, ""
}

// runStaged materializes the inputs into a fresh temp dir and runs the check
// there, so it sees only its declared inputs (by local name).
func runStaged(runPath string, inputs []resolved) (stdout, errMsg string) {
	stage, err := os.MkdirTemp("", "starbase-ev-")
	if err != nil {
		return "", err.Error()
	}
	defer os.RemoveAll(stage)
	for _, r := range inputs {
		if err := r.realize(filepath.Join(stage, r.LocalName)); err != nil {
			return "", fmt.Sprintf("staging %s: %s", r.LocalName, err)
		}
	}
	return execScript(runPath, stage)
}

// computeKey is the content-addressed cache key: the run script, the manifest,
// and every input's resolved content hash.
func computeKey(runBytes, manifest []byte, inputs []resolved) string {
	h := sha256.New()
	h.Write(runBytes)
	h.Write([]byte{0})
	h.Write(manifest)
	sorted := append([]resolved(nil), inputs...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].LocalName < sorted[j].LocalName })
	for _, r := range sorted {
		fmt.Fprintf(h, "\x00%s\x00%s", r.LocalName, r.Hash)
	}
	return hex.EncodeToString(h.Sum(nil))
}

func isExecutable(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.Mode().IsRegular() && fi.Mode().Perm()&0o111 != 0
}

// execScript runs an executable with the given working directory, returning
// stdout and, on a non-zero exit, the stderr (or the error).
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

// --- persisted cache (user cache dir; CI starts cold) ---

type cachedUnit struct {
	Key   string `json:"key"`
	Check Check  `json:"check"`
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
