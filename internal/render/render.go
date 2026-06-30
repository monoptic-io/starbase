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
	"strings"

	"github.com/ryannedolan/sitegen/internal/graph"
	"github.com/ryannedolan/sitegen/internal/model"
	"github.com/ryannedolan/sitegen/internal/parse"
	"github.com/ryannedolan/sitegen/internal/registry"
	"github.com/ryannedolan/sitegen/internal/tmpl"
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
}

// NavNode is one entry in the sidebar navigation tree.
type NavNode struct {
	Title    string
	URL      string // root-relative output path; "" for a pure group
	Slug     string
	Weight   int
	Children []*NavNode
	Active   bool
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
		Nav:         cloneNavFor(r.nav, t.Slug, rootRel(t.OutPath)),
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
		Nav:         cloneNavFor(r.nav, slug, rootRel(outPath)),
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
			innerHTML := ""
			if sc.Inner != "" {
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

func cloneNavFor(n *NavNode, curSlug, rootRel string) *NavNode {
	c := &NavNode{
		Title: n.Title, URL: n.URL, Slug: n.Slug, Weight: n.Weight,
		Active: n.Slug != "" && n.Slug == curSlug, RootRel: rootRel,
	}
	for _, ch := range n.Children {
		cc := cloneNavFor(ch, curSlug, rootRel)
		if cc.Active {
			c.Active = true // keep ancestor path highlighted
		}
		c.Children = append(c.Children, cc)
	}
	return c
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
