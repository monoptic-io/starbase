// Package render turns parsed topics into finished HTML pages: markdown to
// HTML (goldmark), heading anchors + table of contents, inline math, shortcode
// expansion, wiki-link resolution, and the surrounding page layout with its
// related-topics and backlinks panels.
package render

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"io/fs"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/monoptic-io/starbase/internal/claim"
	"github.com/monoptic-io/starbase/internal/graph"
	"github.com/monoptic-io/starbase/internal/model"
	"github.com/monoptic-io/starbase/internal/parse"
	"github.com/monoptic-io/starbase/internal/registry"
	"github.com/monoptic-io/starbase/internal/tmpl"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	gmrender "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// Site holds project-wide settings exposed to templates.
type Site struct {
	Title        string
	BaseURL      string
	AssetVersion string // short content hash, appended to asset URLs for cache-busting
	KaTeXBase    string // base URL for KaTeX assets; empty means use the locally vendored copy
}

type Renderer struct {
	site   Site
	eng    *tmpl.Engine
	reg    *registry.Registry
	graph  *graph.Graph
	topics map[string]*model.Topic
	md     goldmark.Markdown
	layout *template.Template
	nav    *NavNode
	checks map[string]CheckResult // evidence outputs, injected by {{< val >}} / {{< data >}}
}

// CheckResult is one evidence check's computed output, supplied to the renderer
// so val/data shortcodes inject real values instead of the author transcribing
// them. Err (a non-zero exit) turns any reference into a build error.
type CheckResult struct {
	Output string
	Err    string
}

// SetChecks provides the evidence outputs the injection shortcodes read from.
func (r *Renderer) SetChecks(m map[string]CheckResult) { r.checks = m }

// NavNode is one entry in the sidebar navigation tree.
type NavNode struct {
	Title    string
	URL      string // root-relative output path; "" for a pure group
	Slug     string
	Weight   int
	Children []*NavNode
	Active   bool
	Open     bool // section is expanded (contains the active page)
	Divider  bool // render a divider before this node (separates disjoint clusters)
	RootRel  string
}

type crumb struct{ Title, URL string }
type linkView struct{ Title, URL string }
type relatedView struct{ Title, URL, Summary string }

type pageData struct {
	Site        Site
	Title       string
	Summary     string
	RootRel     string
	Breadcrumbs []crumb
	Tags        []string
	WordCount   int
	ReadingTime int
	ShowTOC     bool
	TOC         []model.Heading
	ContentHTML template.HTML
	Related     []relatedView
	Backlinks   []linkView
	Nav         *NavNode
}

// New builds a Renderer. layoutFS supplies the page templates (built-in plus
// any project override merged in by the caller).
func New(site Site, eng *tmpl.Engine, reg *registry.Registry, g *graph.Graph, topics []*model.Topic, layoutSrcs map[string]string) (*Renderer, error) {
	r := &Renderer{
		site: site, eng: eng, reg: reg, graph: g,
		topics: make(map[string]*model.Topic, len(topics)),
	}
	for _, t := range topics {
		r.topics[t.Slug] = t
	}
	r.md = goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.Footnote, extension.DefinitionList),
		goldmark.WithParserOptions(parser.WithASTTransformers(
			util.Prioritized(headingTransformer{}, 100),
		)),
		goldmark.WithRendererOptions(gmrender.WithUnsafe()),
	)
	r.nav = buildNav(topics)

	layoutFuncs := template.FuncMap{
		"slugify": parse.Slugify,
		"lower":   strings.ToLower,
		"add":     func(a, b int) int { return a + b },
	}
	r.layout = template.New("layout").Funcs(layoutFuncs)
	for name, src := range layoutSrcs {
		if _, err := r.layout.New(name).Parse(src); err != nil {
			return nil, fmt.Errorf("layout %q: %w", name, err)
		}
	}
	return r, nil
}

// LoadLayout reads layout templates from an fs directory into a name->source map.
func LoadLayout(fsys fs.FS, dir string) (map[string]string, error) {
	out := map[string]string{}
	err := fs.WalkDir(fsys, dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || path.Ext(p) != ".html" {
			return err
		}
		b, e := fs.ReadFile(fsys, p)
		if e != nil {
			return e
		}
		out[strings.TrimSuffix(path.Base(p), ".html")] = string(b)
		return nil
	})
	return out, err
}

// Page renders a topic to a complete HTML document.
func (r *Renderer) Page(t *model.Topic) ([]byte, []model.Diagnostic) {
	processed, repl, headings, diags := r.renderBody(t, t.Body, 0)
	content := substitute(processed, repl)

	toc := tocEntries(headings)
	data := pageData{
		Site:        r.site,
		Title:       t.Title,
		Summary:     t.Summary,
		RootRel:     rootRel(t.OutPath),
		Breadcrumbs: r.breadcrumbs(t),
		Tags:        t.Tags,
		WordCount:   t.WordCount,
		ReadingTime: (t.WordCount + 199) / 200,
		ShowTOC:     len(toc) >= 3,
		TOC:         toc,
		ContentHTML: template.HTML(content),
		Related:     r.relatedViews(t),
		Backlinks:   r.backlinkViews(t),
		Nav:         r.navFor(t.Slug, rootRel(t.OutPath)),
	}
	var buf bytes.Buffer
	if err := r.layout.ExecuteTemplate(&buf, "page", data); err != nil {
		diags = append(diags, model.Diagnostic{
			Severity: model.SevError, File: t.SourcePath, Line: 1,
			Message: fmt.Sprintf("layout render failed: %v", err),
		})
		return nil, diags
	}
	return buf.Bytes(), diags
}

// Listing renders an auto-generated index/section/tag page from a set of cards.
func (r *Renderer) Listing(title, slug, outPath, intro string, cards []Card) ([]byte, error) {
	var b strings.Builder
	if intro != "" {
		b.WriteString("<p>" + html.EscapeString(intro) + "</p>\n")
	}
	b.WriteString(`<ul class="sg-cards">`)
	for _, c := range cards {
		b.WriteString(fmt.Sprintf(
			`<li class="sg-card"><a href="%s%s"><h3>%s</h3><p>%s</p></a></li>`,
			rootRel(outPath), c.URL, html.EscapeString(c.Title), html.EscapeString(c.Summary)))
	}
	b.WriteString(`</ul>`)

	data := pageData{
		Site: r.site, Title: title, RootRel: rootRel(outPath),
		ContentHTML: template.HTML(b.String()),
		Nav:         r.navFor(slug, rootRel(outPath)),
	}
	var buf bytes.Buffer
	if err := r.layout.ExecuteTemplate(&buf, "page", data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Card is an entry on a listing page.
type Card struct{ Title, URL, Summary string }

// Nav exposes the navigation tree (for building listing cards in order).
func (r *Renderer) Nav() *NavNode { return r.nav }

// --- body transformation: shortcodes, wiki links, math ---

const phStart, phEnd = '', ''

func placeholder(i int) string { return string(phStart) + fmt.Sprint(i) + string(phEnd) }

var rePH = regexp.MustCompile("\\d+")

type segment struct {
	start, end int
	prio       int // higher wins overlaps: shortcode > link > math
	htmlOut    template.HTML
}

func (r *Renderer) renderBody(t *model.Topic, body string, depth int) (string, map[string]template.HTML, []model.Heading, []model.Diagnostic) {
	var diags []model.Diagnostic
	masked := parse.MaskCode(body)
	var segs []segment
	rr := rootRel(t.OutPath)

	page := tmpl.PageContext{Title: t.Title, Slug: t.Slug, URL: t.OutPath}

	if depth < 4 {
		for i, sc := range parse.ScanShortcodes(body) {
			id := fmt.Sprintf("sg-%s-%d-%d", sanitizeID(t.Slug), depth, i)
			if sc.Name == "claim" {
				out, ds := r.renderClaim(t, sc, id, depth)
				diags = append(diags, ds...)
				segs = append(segs, segment{sc.Start, sc.End, 3, out})
				continue
			}
			if sc.Name == "val" {
				out, ds := r.renderVal(t, sc)
				diags = append(diags, ds...)
				segs = append(segs, segment{sc.Start, sc.End, 3, out})
				continue
			}
			if sc.Name == "data" {
				out, ds := r.renderData(t, sc, page, id)
				diags = append(diags, ds...)
				segs = append(segs, segment{sc.Start, sc.End, 3, out})
				continue
			}
			innerHTML := ""
			if sc.Inner != "" && !r.eng.RawInner(sc.Name) {
				ip, ir, _, id2 := r.renderBody(t, sc.Inner, depth+1)
				innerHTML = substitute(ip, ir)
				diags = append(diags, id2...)
			}
			out, ds := r.eng.Render(sc, innerHTML, sc.Inner, page, id, t.SourcePath)
			diags = append(diags, ds...)
			segs = append(segs, segment{sc.Start, sc.End, 3, out})
		}
	}
	for _, l := range parse.ScanLinks(body) {
		segs = append(segs, segment{l.Start, l.End, 2, r.linkHTML(l, rr)})
	}
	for _, m := range findMath(masked, body) {
		segs = append(segs, m)
	}

	segs = resolveOverlaps(segs)

	var b strings.Builder
	repl := map[string]template.HTML{}
	pos := 0
	for i, s := range segs {
		if s.start < pos {
			continue
		}
		b.WriteString(body[pos:s.start])
		ph := placeholder(i)
		b.WriteString(ph)
		repl[ph] = s.htmlOut
		pos = s.end
	}
	b.WriteString(body[pos:])

	// Convert markdown, collecting headings.
	var out bytes.Buffer
	headings := []model.Heading{}
	ctx := parser.NewContext()
	ctx.Set(tocKey, &headings)
	if err := r.md.Convert([]byte(b.String()), &out, parser.WithContext(ctx)); err != nil {
		diags = append(diags, model.Diagnostic{Severity: model.SevError, File: t.SourcePath, Message: err.Error()})
	}
	return out.String(), repl, headings, diags
}

// renderClaim turns a {{< claim >}} into prose with an evidence disclosure:
// the result and the implementation (verbatim) the agent computed, so a reader
// can dig down to how a value was derived. starbase runs nothing — it surfaces.
func (r *Renderer) renderClaim(t *model.Topic, sc model.Shortcode, id string, depth int) (template.HTML, []model.Diagnostic) {
	info := claim.Parse(sc)
	pp, pr, _, diags := r.renderBody(t, info.Prose, depth+1)
	proseHTML := stripOuterP(substitute(pp, pr))

	// A bound check that ran cleanly is verified; a pasted result/source is
	// attested; nothing is unsupported.
	computed, hasComputed := CheckResult{}, false
	if info.Check != "" {
		if cr, ok := r.checks[info.Check]; ok && cr.Err == "" {
			computed, hasComputed = cr, strings.TrimSpace(cr.Output) != ""
		}
	}

	var b strings.Builder
	b.WriteString(`<div class="sg-claim">`)
	b.WriteString(`<div class="sg-claim-text">` + proseHTML)
	switch {
	case hasComputed:
		b.WriteString(` <span class="sg-claim-badge sg-claim-verified" title="recomputed at build time">verified</span>`)
	case info.HasImpl || info.Source != "":
		b.WriteString(` <span class="sg-claim-badge" title="evidence attached">evidence</span>`)
	default:
		b.WriteString(` <span class="sg-claim-badge sg-claim-unsupported" title="no evidence yet">unsupported</span>`)
	}
	b.WriteString(`</div>`)

	// The displayed result is the pasted block if present, else the check's own
	// (verified) output — so tables need not be transcribed into the article.
	resultHTML := ""
	if info.Result != "" {
		resultHTML = claimResultHTML(info)
	} else if hasComputed {
		resultHTML = rowsTableHTML(claim.Rows(strings.TrimSpace(computed.Output), "csv"))
	}

	if info.HasImpl || info.Source != "" || resultHTML != "" || hasComputed {
		b.WriteString(`<details class="sg-evidence"><summary>How we know this</summary>`)
		var meta []string
		if info.Check != "" && hasComputed {
			meta = append(meta, "computed by <code>evidence/"+html.EscapeString(info.Check)+"</code>")
		}
		if info.Source != "" {
			meta = append(meta, "source <code>"+html.EscapeString(info.Source)+"</code>")
		}
		if info.AsOf != "" {
			meta = append(meta, "as of "+html.EscapeString(info.AsOf))
		}
		if len(meta) > 0 {
			b.WriteString(`<div class="sg-ev-meta">` + strings.Join(meta, " · ") + `</div>`)
		}
		if resultHTML != "" {
			b.WriteString(resultHTML)
		}
		if info.HasImpl {
			lang := info.Lang
			if lang == "" {
				lang = "text"
			}
			b.WriteString(`<pre class="sg-ev-code"><code class="language-` +
				html.EscapeString(lang) + `">` + html.EscapeString(info.Code) + `</code></pre>`)
		}
		b.WriteString(`</details>`)
	}
	b.WriteString(`</div>`)
	return template.HTML(b.String()), diags
}

func claimResultHTML(info claim.Info) string {
	return rowsTableHTML(claim.Rows(info.Result, info.ResultFmt))
}

func rowsTableHTML(rows [][]string) string {
	if len(rows) == 0 {
		return ""
	}
	if len(rows) == 1 && len(rows[0]) == 1 {
		return `<div class="sg-ev-result">= <b>` + html.EscapeString(rows[0][0]) + `</b></div>`
	}
	var b strings.Builder
	b.WriteString(`<table class="sg-ev-table">`)
	for ri, row := range rows {
		b.WriteString("<tr>")
		tag := "td"
		if ri == 0 {
			tag = "th"
		}
		for _, cell := range row {
			b.WriteString("<" + tag + ">" + html.EscapeString(cell) + "</" + tag + ">")
		}
		b.WriteString("</tr>")
	}
	b.WriteString("</table>")
	return b.String()
}

// renderVal injects a check's scalar stdout inline, so a number in prose comes
// from the computation instead of being transcribed by the author.
func (r *Renderer) renderVal(t *model.Topic, sc model.Shortcode) (template.HTML, []model.Diagnostic) {
	cr, msg := r.lookupCheck(sc)
	if msg != "" {
		return valErr("[val: " + msg + "]"), evDiag(t, sc, "val "+msg)
	}
	name := strings.TrimSpace(sc.Args["check"])
	v := strings.TrimSpace(cr.Output)
	switch {
	case strings.TrimSpace(sc.Args["col"]) != "":
		col, row := strings.TrimSpace(sc.Args["col"]), strings.TrimSpace(sc.Args["row"])
		got, ok := pickCell(cr.Output, col, row)
		if !ok {
			return valErr("[val: no cell " + col + "]"),
				evDiag(t, sc, fmt.Sprintf("val: check %q has no column %q (row %q)", name, col, row))
		}
		v = got
	case strings.TrimSpace(sc.Args["field"]) != "":
		field := strings.TrimSpace(sc.Args["field"])
		got, ok := pickField(cr.Output, field)
		if !ok {
			return valErr("[val: no field " + field + "]"),
				evDiag(t, sc, fmt.Sprintf("val: check %q output has no field %q", name, field))
		}
		v = got
	}
	return template.HTML(`<span class="sg-val" title="computed by evidence/` + html.EscapeString(name) + `">` +
		html.EscapeString(v) + `</span>`), nil
}

// pickCell selects a cell from a check's tabular output by column header name.
// row is a 1-indexed data row (default 1) or a value matched against column 0.
func pickCell(output, col, row string) (string, bool) {
	rows := claim.Rows(strings.TrimSpace(output), "")
	if len(rows) < 2 {
		return "", false
	}
	ci := -1
	for i, h := range rows[0] {
		if strings.EqualFold(strings.TrimSpace(h), col) {
			ci = i
			break
		}
	}
	if ci < 0 {
		return "", false
	}
	data := rows[1:]
	ri := 0
	if row != "" {
		if n, err := strconv.Atoi(row); err == nil {
			ri = n - 1
		} else {
			ri = -1
			for i, r := range data {
				if len(r) > 0 && strings.EqualFold(strings.TrimSpace(r[0]), row) {
					ri = i
					break
				}
			}
		}
	}
	if ri < 0 || ri >= len(data) || ci >= len(data[ri]) {
		return "", false
	}
	return strings.TrimSpace(data[ri][ci]), true
}

// pickField extracts a named field from a check's output, matching a
// `field = value`, `field: value`, or `field ~ value` token (whole word).
func pickField(output, field string) (string, bool) {
	re := regexp.MustCompile(`(?mi)\b` + regexp.QuoteMeta(field) + `\b\s*[=:~]\s*([^\s,;]+)`)
	if m := re.FindStringSubmatch(output); m != nil {
		return m[1], true
	}
	return "", false
}

// renderData renders a check's structured stdout as a table (default) or a chart,
// reusing the chart runtime — no data is transcribed into the article.
func (r *Renderer) renderData(t *model.Topic, sc model.Shortcode, page tmpl.PageContext, id string) (template.HTML, []model.Diagnostic) {
	cr, msg := r.lookupCheck(sc)
	if msg != "" {
		return valErr("[data: " + msg + "]"), evDiag(t, sc, "data "+msg)
	}
	switch as := strings.ToLower(strings.TrimSpace(sc.Args["as"])); as {
	case "", "table":
		name := html.EscapeString(strings.TrimSpace(sc.Args["check"]))
		return template.HTML(`<div class="sg-data" data-check="` + name + `">` +
			rowsTableHTML(claim.Rows(strings.TrimSpace(cr.Output), "csv")) + `</div>`), nil
	case "bar", "line", "scatter":
		chart := model.Shortcode{Name: "chart", Line: sc.Line, Args: map[string]string{
			"type": as, "data": csvToPairs(cr.Output), "height": firstArg(sc.Args["height"], "300"),
			"title": sc.Args["title"], "caption": sc.Args["caption"],
			"xlabel": sc.Args["xlabel"], "ylabel": sc.Args["ylabel"], "labels": sc.Args["labels"],
		}}
		return r.eng.Render(chart, "", "", page, id, t.SourcePath)
	default:
		return valErr("[data: bad as]"), evDiag(t, sc, fmt.Sprintf("data: unknown as=%q (use table, bar, line, or scatter)", as))
	}
}

// lookupCheck resolves the check= argument, returning a human message on failure.
func (r *Renderer) lookupCheck(sc model.Shortcode) (CheckResult, string) {
	name := strings.TrimSpace(sc.Args["check"])
	if name == "" {
		return CheckResult{}, `missing required argument "check"`
	}
	cr, ok := r.checks[name]
	if !ok {
		return CheckResult{}, fmt.Sprintf("references unknown evidence check %q", name)
	}
	if cr.Err != "" {
		return CheckResult{}, fmt.Sprintf("check %q failed: %s", name, evFirstLine(cr.Err))
	}
	return cr, ""
}

func evDiag(t *model.Topic, sc model.Shortcode, msg string) []model.Diagnostic {
	return []model.Diagnostic{{Severity: model.SevError, File: t.SourcePath, Line: sc.Line, Message: msg}}
}

func valErr(text string) template.HTML {
	return template.HTML(`<span class="sg-val sg-val-err" title="unresolved evidence reference">` + html.EscapeString(text) + `</span>`)
}

// csvToPairs turns a check's CSV stdout into the chart template's "label:value"
// data string (skipping a non-numeric header row).
func csvToPairs(out string) string {
	rows := claim.Rows(strings.TrimSpace(out), "csv")
	start := 0
	if len(rows) > 0 && len(rows[0]) >= 2 {
		if _, err := strconv.ParseFloat(strings.TrimSpace(rows[0][len(rows[0])-1]), 64); err != nil {
			start = 1 // header row
		}
	}
	var pairs []string
	for _, row := range rows[start:] {
		if len(row) < 2 {
			continue
		}
		pairs = append(pairs, strings.TrimSpace(row[0])+":"+strings.TrimSpace(row[len(row)-1]))
	}
	return strings.Join(pairs, ", ")
}

func firstArg(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func evFirstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

func stripOuterP(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "<p>") && strings.HasSuffix(s, "</p>") && strings.Count(s, "<p>") == 1 {
		return s[3 : len(s)-4]
	}
	return s
}

func substitute(htmlStr string, repl map[string]template.HTML) string {
	for ph, frag := range repl {
		htmlStr = strings.ReplaceAll(htmlStr, "<p>"+ph+"</p>", string(frag))
		htmlStr = strings.ReplaceAll(htmlStr, ph, string(frag))
	}
	return htmlStr
}

func (r *Renderer) linkHTML(l model.Link, rootRel string) template.HTML {
	disp := l.Display
	if disp == "" {
		disp = l.Target
	}
	if dst, ok := r.reg.Resolve(l.Target); ok {
		frag := ""
		if i := strings.IndexByte(l.Target, '#'); i >= 0 {
			frag = "#" + parse.Slugify(l.Target[i+1:])
		}
		return template.HTML(fmt.Sprintf(`<a href="%s%s%s">%s</a>`,
			rootRel, template.HTMLEscapeString(dst.OutPath), frag, template.HTMLEscapeString(disp)))
	}
	return template.HTML(fmt.Sprintf(`<a class="sg-deadlink" title="missing topic: %s" href="#">%s</a>`,
		template.HTMLEscapeString(l.Target), template.HTMLEscapeString(disp)))
}

func findMath(masked, orig string) []segment {
	var segs []segment
	display := regexp.MustCompile(`(?s)\$\$(.+?)\$\$`)
	used := make([]bool, len(masked))
	for _, loc := range display.FindAllStringSubmatchIndex(masked, -1) {
		tex := strings.TrimSpace(orig[loc[2]:loc[3]])
		segs = append(segs, segment{loc[0], loc[1], 1,
			template.HTML(`<div class="sg-math sg-math-display">` + html.EscapeString(tex) + `</div>`)})
		for i := loc[0]; i < loc[1]; i++ {
			used[i] = true
		}
	}
	inline := regexp.MustCompile(`\$([^$\n]+?)\$`)
	for _, loc := range inline.FindAllStringSubmatchIndex(masked, -1) {
		if used[loc[0]] {
			continue
		}
		tex := strings.TrimSpace(orig[loc[2]:loc[3]])
		segs = append(segs, segment{loc[0], loc[1], 1,
			template.HTML(`<span class="sg-math-inline">` + html.EscapeString(tex) + `</span>`)})
	}
	return segs
}

func resolveOverlaps(segs []segment) []segment {
	sort.Slice(segs, func(i, j int) bool {
		if segs[i].start != segs[j].start {
			return segs[i].start < segs[j].start
		}
		return segs[i].prio > segs[j].prio
	})
	var out []segment
	last := -1
	for _, s := range segs {
		if s.start < last {
			continue // overlaps a kept segment
		}
		out = append(out, s)
		last = s.end
	}
	return out
}

// --- headings / TOC via goldmark transformer ---

var tocKey = parser.NewContextKey()

type headingTransformer struct{}

func (headingTransformer) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	coll, _ := pc.Get(tocKey).(*[]model.Heading)
	seen := map[string]int{}
	src := reader.Source()
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		h, ok := n.(*ast.Heading)
		if !ok || !entering {
			return ast.WalkContinue, nil
		}
		txt := stripPH(nodeText(n, src))
		slug := parse.Slugify(txt)
		if slug == "" {
			slug = "section"
		}
		if k := seen[slug]; k > 0 {
			seen[slug] = k + 1
			slug = fmt.Sprintf("%s-%d", slug, k)
		} else {
			seen[slug] = 1
		}
		h.SetAttributeString("id", []byte(slug))
		if coll != nil {
			*coll = append(*coll, model.Heading{Level: h.Level, Text: txt, Slug: slug})
		}
		return ast.WalkContinue, nil
	})
}

func nodeText(n ast.Node, src []byte) string {
	var b strings.Builder
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		switch v := c.(type) {
		case *ast.Text:
			b.Write(v.Segment.Value(src))
		case *ast.String:
			b.Write(v.Value)
		default:
			b.WriteString(nodeText(c, src))
		}
	}
	return b.String()
}

func stripPH(s string) string { return strings.TrimSpace(rePH.ReplaceAllString(s, "")) }

func tocEntries(hs []model.Heading) []model.Heading {
	var out []model.Heading
	for _, h := range hs {
		if h.Level == 2 || h.Level == 3 {
			out = append(out, h)
		}
	}
	return out
}

// --- navigation / breadcrumbs / related ---

func (r *Renderer) breadcrumbs(t *model.Topic) []crumb {
	if t.Slug == "index" {
		return nil
	}
	crumbs := []crumb{{Title: "Home", URL: "index.html"}}
	segs := strings.Split(path.Dir(t.Slug), "/")
	if segs[0] == "." || segs[0] == "" {
		return crumbs
	}
	cur := ""
	for _, s := range segs {
		if cur == "" {
			cur = s
		} else {
			cur += "/" + s
		}
		if idx, ok := r.topics[cur]; ok {
			crumbs = append(crumbs, crumb{Title: idx.Title, URL: idx.OutPath})
		} else {
			crumbs = append(crumbs, crumb{Title: humanize(path.Base(cur))})
		}
	}
	return crumbs
}

func (r *Renderer) relatedViews(t *model.Topic) []relatedView {
	var out []relatedView
	for _, rel := range r.graph.Related(t.Slug, 6) {
		dst, ok := r.topics[rel.Slug]
		if !ok || dst.Draft {
			continue
		}
		out = append(out, relatedView{Title: dst.Title, URL: dst.OutPath, Summary: dst.Summary})
	}
	return out
}

func (r *Renderer) backlinkViews(t *model.Topic) []linkView {
	seen := map[string]bool{}
	var out []linkView
	for _, slug := range r.graph.Backlinks(t.Slug) {
		if seen[slug] {
			continue
		}
		seen[slug] = true
		if dst, ok := r.topics[slug]; ok && !dst.Draft {
			out = append(out, linkView{Title: dst.Title, URL: dst.OutPath})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Title < out[j].Title })
	return out
}

func buildNav(topics []*model.Topic) *NavNode {
	root := &NavNode{}
	nodes := map[string]*NavNode{"": root}
	var group func(p string) *NavNode
	group = func(p string) *NavNode {
		if n, ok := nodes[p]; ok {
			return n
		}
		parent := group(parentDir(p))
		n := &NavNode{Title: humanize(path.Base(p)), Slug: p}
		nodes[p] = n
		parent.Children = append(parent.Children, n)
		return n
	}
	for _, t := range topics {
		if t.Draft || t.Slug == "index" {
			continue
		}
		if isSectionIndex(t) {
			n := group(t.Slug)
			n.Title, n.URL, n.Weight = t.Title, t.OutPath, t.Weight
			continue
		}
		parent := group(parentDir(t.Slug))
		leaf := &NavNode{Title: t.Title, URL: t.OutPath, Slug: t.Slug, Weight: t.Weight}
		nodes[t.Slug] = leaf
		parent.Children = append(parent.Children, leaf)
	}
	sortNav(root)
	return root
}

func sortNav(n *NavNode) {
	sort.SliceStable(n.Children, func(i, j int) bool {
		a, b := n.Children[i], n.Children[j]
		if a.Weight != b.Weight {
			return a.Weight < b.Weight
		}
		// sections (with children) after leaves? keep groups last for readability
		if (len(a.Children) > 0) != (len(b.Children) > 0) {
			return len(a.Children) == 0
		}
		return a.Title < b.Title
	})
	for _, c := range n.Children {
		sortNav(c)
	}
}

// navFor builds the per-page sidebar. A topic shows only its own connected
// component (its reachable world); the home page — and any page outside the
// graph — shows everything, with disjoint clusters separated by a divider. The
// section containing the active page is expanded.
func (r *Renderer) navFor(curSlug, rootRel string) *NavNode {
	comp, known := r.graph.Components[curSlug]
	showAll := !known || curSlug == graph.HubSlug
	inComp := func(slug string) bool {
		if showAll {
			return true
		}
		c, ok := r.graph.Components[slug]
		return ok && c == comp
	}
	root := pruneNav(r.nav, curSlug, rootRel, inComp)
	if showAll && r.graph.ComponentCount() > 1 {
		r.clusterOrder(root)
	}
	return root
}

func pruneNav(n *NavNode, curSlug, rootRel string, inComp func(string) bool) *NavNode {
	c := &NavNode{
		Title: n.Title, URL: n.URL, Slug: n.Slug, Weight: n.Weight, RootRel: rootRel,
		Active: n.Slug != "" && n.Slug == curSlug,
	}
	for _, ch := range n.Children {
		cc := pruneNav(ch, curSlug, rootRel, inComp)
		keep := len(cc.Children) > 0 || (cc.Slug != "" && inComp(cc.Slug))
		if !keep {
			continue
		}
		if cc.Active {
			c.Active = true
		}
		c.Children = append(c.Children, cc)
	}
	if len(c.Children) > 0 && c.Slug != "" &&
		(curSlug == c.Slug || strings.HasPrefix(curSlug, c.Slug+"/")) {
		c.Open = true
	}
	return c
}

// clusterOrder (home page only) groups top-level sections by component, largest
// cluster first, and marks a divider where one cluster ends and the next begins.
func (r *Renderer) clusterOrder(root *NavNode) {
	comp := func(n *NavNode) int {
		if c, ok := r.graph.Components[n.Slug]; ok {
			return c
		}
		return 1 << 30
	}
	sort.SliceStable(root.Children, func(i, j int) bool {
		a, b := root.Children[i], root.Children[j]
		if ca, cb := comp(a), comp(b); ca != cb {
			return ca < cb
		}
		if a.Weight != b.Weight {
			return a.Weight < b.Weight
		}
		return a.Title < b.Title
	})
	last := -1
	for _, ch := range root.Children {
		c := comp(ch)
		if last != -1 && c != last {
			ch.Divider = true
		}
		last = c
	}
}

func isSectionIndex(t *model.Topic) bool {
	return strings.HasSuffix(t.OutPath, "/index.html") && t.Slug != "index"
}

func parentDir(p string) string {
	d := path.Dir(p)
	if d == "." || d == "/" {
		return ""
	}
	return d
}

func rootRel(outPath string) string {
	depth := strings.Count(outPath, "/")
	if depth == 0 {
		return ""
	}
	return strings.Repeat("../", depth)
}

func sanitizeID(s string) string { return strings.ReplaceAll(parse.Slugify(s), "/", "-") }

var reHumanSep = regexp.MustCompile(`[-_]+`)

func humanize(s string) string {
	s = reHumanSep.ReplaceAllString(s, " ")
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
