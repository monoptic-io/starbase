// Package evidence runs a knowledge base's evidence functions and returns their
// results, so claims can be verified by re-execution rather than trust.
//
// It works like `go test`: the author writes plain exported functions in an
// `evidence/` Go module — no main, no JSON, no boilerplate —
//
//	package sales
//	//starbase:deps ../../data/sales.csv
//	func MidwestRegions() (int, error) { ... return n, nil }
//	func RevenueByDivision() ([][]string, error) { ... return table, nil }
//
// starbase discovers those functions (an exported func with no parameters
// returning a value, optionally with an error), generates a tiny runner that
// calls them — exactly as `go test` generates a test main — runs it, and binds
// each result to a claim via `check="midwest-regions"` (the kebab-cased name).
//
// It is incremental: each package is a cache unit, keyed by a hash of its Go
// sources plus the data files it declares with `//starbase:deps`. A package is
// re-run only when its own code or data changes. The cache lives in the user
// cache dir, so CI starts cold and is authoritative.
package evidence

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
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
	Name   string
	Cached bool
	Err    string
}

// Result is the merged output of every package plus per-package status.
type Result struct {
	Checks map[string]Check
	Units  []Unit
}

const perRunTimeout = 20 * time.Minute

type fn struct {
	GoName  string // exported Go function name
	Name    string // kebab-cased check name
	Results int    // 1 or 2 (value, or value+error)
}

type pkg struct {
	ImportPath string
	Dir        string
	RelName    string // path relative to evidence/, for display
	GoFiles    []string
	funcs      []fn
	key        string
}

// Run resolves every evidence package, reusing cached results whose inputs are
// unchanged and re-running the rest. The bool reports whether evidence exists.
func Run(contentDir string, force bool) (Result, bool, error) {
	evDir := filepath.Join(contentDir, "evidence")
	if fi, err := os.Stat(evDir); err != nil || !fi.IsDir() {
		return Result{}, false, nil
	}
	pkgs, modPath, err := listPackages(evDir)
	if err != nil {
		return Result{}, true, err
	}

	cache := loadCache(contentDir)
	res := Result{Checks: map[string]Check{}}
	var toRun []*pkg

	for i := range pkgs {
		p := &pkgs[i]
		if len(p.funcs) == 0 {
			continue // not an evidence package
		}
		p.key = unitKey(p, contentDir)
		if !force {
			if c, ok := cache[p.ImportPath]; ok && c.Key == p.key {
				res.Units = append(res.Units, Unit{Name: p.RelName, Cached: true})
				for n, ck := range c.Checks {
					res.Checks[n] = ck
				}
				continue
			}
		}
		toRun = append(toRun, p)
	}

	if len(toRun) > 0 {
		produced, runErr := runGenerated(evDir, contentDir, modPath, toRun)
		for _, p := range toRun {
			if runErr != "" {
				res.Units = append(res.Units, Unit{Name: p.RelName, Err: runErr})
				continue
			}
			pkgChecks := map[string]Check{}
			for _, f := range p.funcs {
				if ck, ok := produced[f.Name]; ok {
					pkgChecks[f.Name] = ck
					res.Checks[f.Name] = ck
				}
			}
			cache[p.ImportPath] = cachedUnit{Key: p.key, Checks: pkgChecks}
			res.Units = append(res.Units, Unit{Name: p.RelName})
		}
	}
	saveCache(contentDir, cache)
	return res, true, nil
}

// --- package + function discovery ---

func listPackages(evDir string) ([]pkg, string, error) {
	cmd := exec.Command("go", "list", "-json", "./...")
	cmd.Dir = evDir
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return nil, "", fmt.Errorf("evidence/ must be a Go module: %s", msg)
	}
	type listed struct {
		ImportPath string
		Dir        string
		GoFiles    []string
		Module     struct{ Path string }
	}
	var pkgs []pkg
	var modPath string
	dec := json.NewDecoder(&stdout)
	for dec.More() {
		var l listed
		if err := dec.Decode(&l); err != nil {
			return nil, "", err
		}
		if l.Module.Path != "" {
			modPath = l.Module.Path
		}
		rel, _ := filepath.Rel(evDir, l.Dir)
		p := pkg{ImportPath: l.ImportPath, Dir: l.Dir, RelName: filepath.ToSlash(rel)}
		for _, g := range l.GoFiles {
			p.GoFiles = append(p.GoFiles, filepath.Join(l.Dir, g))
		}
		p.funcs = discoverFuncs(p.GoFiles)
		pkgs = append(pkgs, p)
	}
	return pkgs, modPath, nil
}

func discoverFuncs(goFiles []string) []fn {
	var out []fn
	fset := token.NewFileSet()
	for _, gf := range goFiles {
		f, err := parser.ParseFile(fset, gf, nil, 0)
		if err != nil {
			continue
		}
		for _, decl := range f.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok || fd.Recv != nil || !fd.Name.IsExported() {
				continue
			}
			if fd.Type.Params != nil && len(fd.Type.Params.List) > 0 {
				continue
			}
			n := results(fd.Type)
			if n == 1 || n == 2 {
				out = append(out, fn{GoName: fd.Name.Name, Name: kebab(fd.Name.Name), Results: n})
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].GoName < out[j].GoName })
	return out
}

// results returns the count if the signature is a check: (T) or (T, error).
func results(t *ast.FuncType) int {
	if t.Results == nil {
		return 0
	}
	var names []string
	for _, r := range t.Results.List {
		// a field may declare multiple results sharing a type
		count := 1
		if len(r.Names) > 1 {
			count = len(r.Names)
		}
		for i := 0; i < count; i++ {
			if id, ok := r.Type.(*ast.Ident); ok {
				names = append(names, id.Name)
			} else {
				names = append(names, "?")
			}
		}
	}
	switch len(names) {
	case 1:
		return 1
	case 2:
		if names[1] == "error" {
			return 2
		}
	}
	return 0
}

// --- generated runner ---

func runGenerated(evDir, contentDir, modPath string, pkgs []*pkg) (map[string]Check, string) {
	absRoot, _ := filepath.Abs(contentDir)
	genDir := filepath.Join(evDir, ".sbgen")
	_ = os.RemoveAll(genDir)
	if err := os.MkdirAll(genDir, 0o755); err != nil {
		return nil, err.Error()
	}
	defer os.RemoveAll(genDir)

	var imports, calls strings.Builder
	for i, p := range pkgs {
		alias := fmt.Sprintf("p%d", i)
		fmt.Fprintf(&imports, "\t%s %q\n", alias, p.ImportPath)
		for _, f := range p.funcs {
			if f.Results == 2 {
				fmt.Fprintf(&calls, "\t{ v, err := %s.%s(); emit(%q, v, err) }\n", alias, f.GoName, f.Name)
			} else {
				fmt.Fprintf(&calls, "\t{ v := %s.%s(); emit(%q, v, nil) }\n", alias, f.GoName, f.Name)
			}
		}
	}
	// Evidence runs with the knowledge-base root as its working directory, so
	// functions reference data relative to the KB (e.g. "data/sales.csv").
	chdir := fmt.Sprintf("\tif err := os.Chdir(%q); err != nil { panic(err) }\n", absRoot)
	src := "package main\n\nimport (\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"os\"\n\t\"reflect\"\n" +
		imports.String() + ")\n\n" + runnerBody + "\nfunc main() {\n" + chdir + calls.String() +
		"\tjson.NewEncoder(os.Stdout).Encode(out)\n}\n"
	if err := os.WriteFile(filepath.Join(genDir, "main.go"), []byte(src), 0o644); err != nil {
		return nil, err.Error()
	}

	ctx, cancel := context.WithTimeout(context.Background(), perRunTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "run", "./.sbgen")
	cmd.Dir = evDir
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return nil, msg
	}
	var out map[string]Check
	if err := json.Unmarshal(bytes.TrimSpace(stdout.Bytes()), &out); err != nil {
		return nil, "evidence runner produced invalid output: " + err.Error()
	}
	return out, ""
}

const runnerBody = `type check struct {
	Value string     ` + "`json:\"value,omitempty\"`" + `
	Table [][]string ` + "`json:\"table,omitempty\"`" + `
	Error string     ` + "`json:\"error,omitempty\"`" + `
}
var out = map[string]check{}
func emit(name string, v interface{}, err error) {
	c := check{}
	if err != nil {
		c.Error = err.Error()
	} else if rv := reflect.ValueOf(v); rv.IsValid() {
		if s, ok := v.([][]string); ok {
			c.Table = s
		} else if rv.Kind() == reflect.Slice && rv.Len() > 0 && deref(rv.Index(0)).Kind() == reflect.Slice {
			var t [][]string
			for i := 0; i < rv.Len(); i++ {
				row := deref(rv.Index(i))
				var r []string
				for j := 0; j < row.Len(); j++ {
					r = append(r, fmt.Sprintf("%v", row.Index(j).Interface()))
				}
				t = append(t, r)
			}
			c.Table = t
		} else {
			c.Value = fmt.Sprintf("%v", v)
		}
	}
	out[name] = c
}
func deref(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Interface {
		return v.Elem()
	}
	return v
}
`

// --- cache key ---

var reDeps = regexp.MustCompile(`(?m)^\s*//\s*starbase:deps\s+(.+)$`)

func unitKey(p *pkg, contentDir string) string {
	h := sha256.New()
	depSet := map[string]bool{}
	files := append([]string(nil), p.GoFiles...)
	sort.Strings(files)
	for _, gf := range files {
		b, err := os.ReadFile(gf)
		if err != nil {
			continue
		}
		fmt.Fprintf(h, "go:%s\n", filepath.Base(gf))
		h.Write(b)
		for _, m := range reDeps.FindAllStringSubmatch(string(b), -1) {
			for _, pat := range strings.FieldsFunc(m[1], func(r rune) bool { return r == ',' || r == ' ' || r == '\t' }) {
				// dep paths are relative to the KB root (the runner's CWD)
				matches, _ := filepath.Glob(filepath.Join(contentDir, pat))
				for _, mt := range matches {
					depSet[mt] = true
				}
			}
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
		rel, _ := filepath.Rel(p.Dir, dep)
		fmt.Fprintf(h, "dep:%s\n", rel)
		h.Write(b)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// --- kebab-case ---

func kebab(s string) string {
	var b strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				b.WriteByte('-')
			}
			b.WriteRune(r + 32)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
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
