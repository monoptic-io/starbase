// Package tmpl is starbase's embedded-content engine. Authors invoke templates
// from markdown with shortcodes:
//
//	{{< chart type="line" height="320" data="0,1,4,9,16" >}}
//
// Each template declares its parameters with directives in template comments:
//
//	{{/* @doc A line/bar chart drawn to a <canvas>. */}}
//	{{/* @param type default="line" */}}
//	{{/* @param data required */}}
//
// Invoking a template without a required parameter is a hard ERROR, so an agent
// gets immediate, precise feedback instead of a silently broken page.
package tmpl

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"io/fs"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/monoptic-io/starbase/internal/model"
)

type Param struct {
	Name     string
	Required bool
	Default  string
	HasDef   bool
}

type def struct {
	Name   string
	Params map[string]*Param
	order  []string
	Open   bool // allow undeclared args without warning
	Doc    string
	Hash   string
}

type Engine struct {
	root *template.Template
	defs map[string]*def
}

func New() *Engine {
	return &Engine{
		root: template.New("starbase").Funcs(FuncMap()),
		defs: map[string]*def{},
	}
}

// Names returns the registered template names, sorted.
func (e *Engine) Names() []string {
	out := make([]string, 0, len(e.defs))
	for n := range e.defs {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}

// Has reports whether a template is registered.
func (e *Engine) Has(name string) bool { _, ok := e.defs[name]; return ok }

// Doc returns a template's documentation string.
func (e *Engine) Doc(name string) string {
	if d, ok := e.defs[name]; ok {
		return d.Doc
	}
	return ""
}

// Catalog describes one template for the `starbase templates` command.
type Catalog struct {
	Name   string
	Doc    string
	Open   bool
	Params []Param
}

// Catalogs returns a sorted description of every registered template.
func (e *Engine) Catalogs() []Catalog {
	var out []Catalog
	for _, name := range e.Names() {
		d := e.defs[name]
		c := Catalog{Name: name, Doc: d.Doc, Open: d.Open}
		for _, p := range d.order {
			c.Params = append(c.Params, *d.Params[p])
		}
		out = append(out, c)
	}
	return out
}

// Hash returns a combined hash of every loaded template, for incremental builds.
func (e *Engine) Hash() string {
	h := sha256.New()
	for _, n := range e.Names() {
		fmt.Fprintf(h, "%s:%s\n", n, e.defs[n].Hash)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// LoadFS registers every *.html / *.tmpl file under dir in the given fs.
// Later sources override earlier ones with the same name, so user templates can
// shadow built-ins.
func (e *Engine) LoadFS(fsys fs.FS, dir string) error {
	return fs.WalkDir(fsys, dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			if _, statErr := fs.Stat(fsys, dir); statErr != nil {
				return nil // optional directory
			}
			return err
		}
		if d.IsDir() {
			return nil
		}
		ext := path.Ext(p)
		if ext != ".html" && ext != ".tmpl" {
			return nil
		}
		content, err := fs.ReadFile(fsys, p)
		if err != nil {
			return err
		}
		name := strings.TrimSuffix(path.Base(p), ext)
		return e.register(name, string(content))
	})
}

func (e *Engine) register(name, content string) error {
	t := e.root.New(name)
	if _, err := t.Parse(content); err != nil {
		return fmt.Errorf("template %q: %w", name, err)
	}
	sum := sha256.Sum256([]byte(content))
	d := &def{
		Name:   name,
		Params: map[string]*Param{},
		Hash:   hex.EncodeToString(sum[:]),
	}
	parseDirectives(content, d)
	e.defs[name] = d
	return nil
}

var (
	reDirective = regexp.MustCompile(`(?s)\{\{/\*\s*(@\w+[^*]*?)\s*\*/\}\}`)
	reKV        = regexp.MustCompile(`(\w+)\s*=\s*"([^"]*)"`)
)

func parseDirectives(content string, d *def) {
	for _, m := range reDirective.FindAllStringSubmatch(content, -1) {
		body := strings.TrimSpace(m[1])
		switch {
		case strings.HasPrefix(body, "@params"):
			if strings.Contains(body, "open") {
				d.Open = true
			}
		case strings.HasPrefix(body, "@param"):
			parseParam(strings.TrimSpace(body[len("@param"):]), d)
		case strings.HasPrefix(body, "@doc"):
			d.Doc = strings.TrimSpace(body[len("@doc"):])
		}
	}
}

func parseParam(s string, d *def) {
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return
	}
	p := &Param{Name: fields[0]}
	if strings.Contains(s, "required") {
		p.Required = true
	}
	if m := reKV.FindStringSubmatch(s); m != nil && m[1] == "default" {
		p.Default = m[2]
		p.HasDef = true
	}
	d.Params[p.Name] = p
	d.order = append(d.order, p.Name)
}

// PageContext is the page-level data exposed to templates as `.Page`.
type PageContext struct {
	Title string
	Slug  string
	URL   string
}

// Validate checks a shortcode invocation without rendering: unknown template,
// missing required args, unknown args. Used by the fast `check` path.
func (e *Engine) Validate(sc model.Shortcode, file string) []model.Diagnostic {
	var diags []model.Diagnostic
	d, ok := e.defs[sc.Name]
	if !ok {
		return []model.Diagnostic{{
			Severity: model.SevError, File: file, Line: sc.Line,
			Message: fmt.Sprintf("unknown template %q", sc.Name),
		}}
	}
	for _, p := range d.order {
		param := d.Params[p]
		if param.Required {
			if _, given := sc.Args[p]; !given {
				diags = append(diags, model.Diagnostic{
					Severity: model.SevError, File: file, Line: sc.Line,
					Message: fmt.Sprintf("template %q: missing required argument %q", sc.Name, p),
				})
			}
		}
	}
	if !d.Open {
		for arg := range sc.Args {
			if _, declared := d.Params[arg]; !declared {
				diags = append(diags, model.Diagnostic{
					Severity: model.SevWarn, File: file, Line: sc.Line,
					Message: fmt.Sprintf("template %q: unknown argument %q", sc.Name, arg),
				})
			}
		}
	}
	return diags
}

// Render validates and executes a shortcode, returning the HTML it produces.
// innerHTML is the already-rendered inner content (for block shortcodes); id is
// a stable unique id the template can use for DOM elements.
func (e *Engine) Render(sc model.Shortcode, innerHTML, innerRaw string, page PageContext, id, file string) (template.HTML, []model.Diagnostic) {
	diags := e.Validate(sc, file)
	if hasError(diags) {
		return "", diags
	}
	d := e.defs[sc.Name]

	data := map[string]any{}
	for name, p := range d.Params {
		if p.HasDef {
			data[name] = p.Default
		}
	}
	for k, v := range sc.Args {
		data[k] = v
	}
	data["Args"] = sc.Args
	data["Inner"] = template.HTML(innerHTML)
	data["InnerRaw"] = innerRaw
	data["Page"] = page
	data["ID"] = id

	var buf bytes.Buffer
	if err := e.root.ExecuteTemplate(&buf, sc.Name, data); err != nil {
		diags = append(diags, model.Diagnostic{
			Severity: model.SevError, File: file, Line: sc.Line,
			Message: fmt.Sprintf("template %q failed: %v", sc.Name, cleanErr(err)),
		})
		return "", diags
	}
	return template.HTML(buf.String()), diags
}

func hasError(diags []model.Diagnostic) bool {
	for _, d := range diags {
		if d.Severity == model.SevError {
			return true
		}
	}
	return false
}

func cleanErr(err error) string {
	s := err.Error()
	if i := strings.LastIndex(s, ": "); i >= 0 && strings.Contains(s, "executing") {
		return s[i+2:]
	}
	return s
}
