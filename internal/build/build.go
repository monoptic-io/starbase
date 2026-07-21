// Package build orchestrates the two starbase commands:
//
//   - Check: fast validation — parse, resolve links, validate shortcode args.
//     Reports missing topics and bad template calls without rendering anything.
//   - Build: the full static-site generation, with incremental rendering driven
//     by per-page fingerprints.
package build

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/monoptic-io/starbase/internal/assets"
	"github.com/monoptic-io/starbase/internal/cache"
	"github.com/monoptic-io/starbase/internal/claim"
	"github.com/monoptic-io/starbase/internal/evidence"
	"github.com/monoptic-io/starbase/internal/graph"
	"github.com/monoptic-io/starbase/internal/model"
	"github.com/monoptic-io/starbase/internal/parse"
	"github.com/monoptic-io/starbase/internal/registry"
	"github.com/monoptic-io/starbase/internal/render"
	"github.com/monoptic-io/starbase/internal/tmpl"
	"github.com/monoptic-io/starbase/internal/vendor"
)

type Config struct {
	ContentDir string
	OutDir     string
	SiteTitle  string
	BaseURL    string
	Drafts     bool
	Force      bool // ignore cache, rebuild everything
	Trust      bool // never execute evidence checks: require committed attestations
	Vendor     bool // download + bundle third-party assets locally instead of linking a CDN
	Offline    bool // with Vendor, use only cached downloads (never hit the network)
}

type Result struct {
	Topics       int
	Rendered     int
	Skipped      int
	Verified     int             // claims re-executed and matched
	Attested     int             // claims with evidence attached but no re-runnable check
	UnitsRun     int             // evidence units executed this run
	UnitsCached  int             // evidence units served from cache (inputs unchanged)
	UnitsTrusted int             // evidence units served from committed attestations
	Units        []evidence.Unit // per-check status (for -v)
	Claims       []ClaimStatus   // per-claim outcome (for -v)
	Diagnostics  []model.Diagnostic
}

// ClaimStatus is one claim's verification outcome, for the verbose listing.
type ClaimStatus struct {
	Check string
	File  string
	Line  int
	State string // "verified" | "attested" | "failed"
}

// ShowCheck runs the evidence checks and returns the named check's raw stdout —
// for capturing into a claim's result block. It reuses the cache (a check re-runs
// only if its inputs changed), so it's cheap to call in an authoring loop.
func ShowCheck(cfg Config, name string) (string, error) {
	rr, present, err := evidence.Run(cfg.ContentDir, evidence.Options{})
	if err != nil {
		return "", err
	}
	if !present {
		return "", fmt.Errorf("no evidence/ directory under %s", cfg.ContentDir)
	}
	ck, ok := rr.Checks[name]
	if !ok {
		return "", fmt.Errorf("no evidence check named %q (looked for evidence/%s/run)", name, name)
	}
	if ck.Err != "" {
		return "", fmt.Errorf("check %q failed: %s", name, ck.Err)
	}
	return ck.Output, nil
}

func (r Result) Errors() (n int) {
	for _, d := range r.Diagnostics {
		if d.Severity == model.SevError {
			n++
		}
	}
	return n
}

func (r Result) Warnings() (n int) {
	for _, d := range r.Diagnostics {
		if d.Severity == model.SevWarn {
			n++
		}
	}
	return n
}

var ignoredDirs = map[string]bool{
	"templates": true, "theme": true, "layout": true,
	"node_modules": true, "_site": true, "dist": true,
}

// index runs the shared parse + resolve + validate phase used by both commands.
func index(cfg Config) ([]*model.Topic, *registry.Registry, *tmpl.Engine, []model.Diagnostic, error) {
	files, err := collectMarkdown(cfg)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Load templates first: the parser needs to know which templates have
	// opaque (raw) inner blocks so it can skip them when scanning for links.
	eng, engDiags, err := loadEngine(cfg)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	rawInner := eng.RawInnerNames()

	var topics []*model.Topic
	var diags []model.Diagnostic
	if isDir(filepath.Join(cfg.ContentDir, ".git")) {
		diags = append(diags, model.Diagnostic{Severity: model.SevWarn, File: ".",
			Message: "content dir looks like a repo root (.git present); put the KB in a subdirectory and build that (e.g. `starbase build content`), so repo files like README.md aren't rendered as topics"})
	}
	diags = append(diags, engDiags...)
	for _, rel := range files {
		t, ds, err := parse.File(cfg.ContentDir, rel, rawInner)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		diags = append(diags, ds...)
		if t.Draft && !cfg.Drafts {
			continue
		}
		topics = append(topics, t)
	}

	reg, regDiags := registry.New(topics)
	diags = append(diags, regDiags...)
	diags = append(diags, reg.ResolveLinks(topics)...)

	for _, t := range topics {
		for _, sc := range t.Shortcodes {
			switch sc.Name {
			case "claim":
				diags = append(diags, claim.Validate(claim.Parse(sc), t.SourcePath, sc.Line)...)
			case "val", "data":
				diags = append(diags, validateEvidenceRef(cfg, sc, t.SourcePath)...)
				if sc.Name == "data" {
					diags = append(diags, validateDataAs(eng, sc, t.SourcePath)...)
				}
			default:
				diags = append(diags, eng.Validate(sc, t.SourcePath)...)
			}
		}
	}
	return topics, reg, eng, diags, nil
}

// Catalog returns the available template catalog (built-in plus project
// overrides), for the `starbase templates` command.
func Catalog(cfg Config) ([]tmpl.Catalog, error) {
	eng, _, err := loadEngine(cfg)
	if err != nil {
		return nil, err
	}
	return eng.Catalogs(), nil
}

// Verify re-executes the content directory's evidence/ Go program and diffs each
// claim's embedded evidence against the freshly computed result. Mismatches are
// errors, so verify can gate CI: a fabricated number breaks the build. Claims
// without a `check` are "attested" by their author and are not re-run.
func Verify(cfg Config) (Result, error) {
	topics, _, _, diags, err := index(cfg)
	if err != nil {
		return Result{}, err
	}
	res := Result{Topics: len(topics), Diagnostics: diags}

	rr, present, rerr := evidence.Run(cfg.ContentDir, evidence.Options{Force: cfg.Force, Trust: cfg.Trust})
	if rerr != nil {
		res.Diagnostics = append(res.Diagnostics, model.Diagnostic{
			Severity: model.SevError, File: "evidence", Message: rerr.Error()})
		res.Diagnostics = dedupeDiags(res.Diagnostics)
		return res, nil
	}
	res.Units = rr.Units
	for _, u := range rr.Units {
		if u.Err != "" {
			res.Diagnostics = append(res.Diagnostics, model.Diagnostic{
				Severity: model.SevError, File: "evidence/" + u.Name + "/run",
				Message: fmt.Sprintf("check failed: %s", firstLine(u.Err))})
		} else if u.Cached {
			res.UnitsCached++
		} else if u.Trusted {
			res.UnitsTrusted++
		} else {
			res.UnitsRun++
		}
	}

	for _, t := range topics {
		for _, sc := range t.Shortcodes {
			if sc.Name != "claim" {
				continue
			}
			info := claim.Parse(sc)
			if info.Check == "" {
				if info.HasImpl || info.Source != "" {
					res.Attested++
					res.Claims = append(res.Claims, ClaimStatus{File: t.SourcePath, Line: sc.Line, State: "attested"})
				}
				continue
			}
			ck, ok := rr.Checks[info.Check]
			if present && ok && ck.Err != "" {
				// the failing run is already reported at evidence/<name>/run
				res.Claims = append(res.Claims, ClaimStatus{Check: info.Check, File: t.SourcePath, Line: sc.Line, State: "failed"})
				continue
			}
			err := func() string {
				if !present {
					return fmt.Sprintf("claim references check %q but there is no evidence/ directory", info.Check)
				}
				if !ok {
					return fmt.Sprintf("no evidence check named %q", info.Check)
				}
				return compareClaim(info, ck)
			}()
			if err != "" {
				res.Diagnostics = append(res.Diagnostics, model.Diagnostic{
					Severity: model.SevError, File: t.SourcePath, Line: sc.Line, Message: err})
				res.Claims = append(res.Claims, ClaimStatus{Check: info.Check, File: t.SourcePath, Line: sc.Line, State: "failed"})
			} else {
				res.Verified++
				res.Claims = append(res.Claims, ClaimStatus{Check: info.Check, File: t.SourcePath, Line: sc.Line, State: "verified"})
			}
		}
	}
	res.Diagnostics = dedupeDiags(res.Diagnostics)
	return res, nil
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

// compareClaim returns "" if everything a claim embeds is backed by the check's
// output, else a description of the discrepancy. Both the result block and the
// asserted value are checked, so an honest result block can't launder a
// fabricated headline:
//   - a result block (if present) must match the output exactly (trimmed); and
//   - the asserted value (if present) must appear, as a whole word, in the output.
func compareClaim(info claim.Info, ck evidence.Check) string {
	result := strings.TrimSpace(info.Result)
	value := strings.TrimSpace(info.Value)
	if result == "" && value == "" {
		// Nothing transcribed to compare: the claim references a check that ran
		// cleanly and its output is injected/shown. That is verified by
		// construction — there is no author-typed number to disagree with.
		return ""
	}
	if result != "" && normalizeText(ck.Output) != normalizeText(result) {
		return fmt.Sprintf("claim check %q: the article's result does not match the computation\n%s",
			info.Check, firstDiff(result, ck.Output))
	}
	if value != "" && !valueInText(value, ck.Output) {
		return fmt.Sprintf("claim check %q: the asserted value %q does not appear in the computation: %s",
			info.Check, value, oneLine(ck.Output))
	}
	return ""
}

// valueInText reports whether value occurs as a whole word in text, after
// normalizing both: lower-cased, with everything but letters, digits, '.', '%'
// and '-' turned to spaces. So "maxdev 1.7%" matches output "maxdev=1.7%", and
// "29.1%" matches "(29.1%)", but "4" does not match "42".
func valueInText(value, text string) bool {
	nv, nt := normWords(value), normWords(text)
	return nv != "" && strings.Contains(" "+nt+" ", " "+nv+" ")
}

// firstDiff returns the first line where the embedded result and the computed
// output diverge, with spaces made visible (·) when the only difference is
// whitespace — the common trap with hand-padded result tables.
func firstDiff(expected, actual string) string {
	el := strings.Split(normalizeText(expected), "\n")
	al := strings.Split(normalizeText(actual), "\n")
	n := len(el)
	if len(al) > n {
		n = len(al)
	}
	for i := 0; i < n; i++ {
		e, a := lineAt(el, i), lineAt(al, i)
		if e == a {
			continue
		}
		if strings.Join(strings.Fields(e), " ") == strings.Join(strings.Fields(a), " ") {
			e, a = strings.ReplaceAll(e, " ", "·"), strings.ReplaceAll(a, " ", "·") // whitespace-only
		}
		return fmt.Sprintf("           line %d  article:  %s\n                   computed: %s", i+1, e, a)
	}
	return fmt.Sprintf("           article:  %s\n           computed: %s", oneLine(expected), oneLine(actual))
}

func lineAt(lines []string, i int) string {
	if i < len(lines) {
		return lines[i]
	}
	return "(missing)"
}

func normWords(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '%' || r == '-' {
			b.WriteRune(r)
		} else {
			b.WriteRune(' ')
		}
	}
	return strings.Join(strings.Fields(b.String()), " ")
}

// normalizeText trims trailing whitespace per line and overall, and normalizes
// newlines — the one leniency golden-file comparisons universally keep.
func normalizeText(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], " \t")
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func oneLine(s string) string {
	s = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(s, "\r\n", "\n"), "\n", " | "))
	if len(s) > 90 {
		s = s[:87] + "..."
	}
	return s
}

// LabelRow is one (label, topic) pair for the `starbase labels` worklist.
type LabelRow struct {
	Label string
	File  string // source path relative to the content dir
	Title string
}

// Labels lists every labeled topic, sorted by label then title — the worklist
// behind conventions like `open-problem`.
func Labels(cfg Config) ([]LabelRow, error) {
	topics, _, _, _, err := index(cfg)
	if err != nil {
		return nil, err
	}
	var rows []LabelRow
	for _, t := range topics {
		for _, l := range t.Labels {
			rows = append(rows, LabelRow{Label: l, File: t.SourcePath, Title: t.Title})
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Label != rows[j].Label {
			return rows[i].Label < rows[j].Label
		}
		return rows[i].Title < rows[j].Title
	})
	return rows, nil
}

// Check performs fast validation only.
func Check(cfg Config) (Result, error) {
	topics, _, _, diags, err := index(cfg)
	if err != nil {
		return Result{}, err
	}
	return Result{Topics: len(topics), Diagnostics: diags}, nil
}

// Build generates the full site.
func Build(cfg Config) (Result, error) {
	topics, reg, eng, diags, err := index(cfg)
	if err != nil {
		return Result{}, err
	}
	g := graph.Build(topics)

	katexBase := vendor.KaTeXCDN
	if cfg.Vendor {
		katexBase = "" // empty tells the page to use the locally vendored copy
	}
	site := render.Site{
		Title: cfg.SiteTitle, BaseURL: cfg.BaseURL,
		AssetVersion: staticVersion(cfg), KaTeXBase: katexBase,
	}
	layout, err := loadLayout(cfg)
	if err != nil {
		return Result{}, err
	}
	rd, err := render.New(site, eng, reg, g, topics, layout)
	if err != nil {
		return Result{}, err
	}

	if err := os.MkdirAll(cfg.OutDir, 0o755); err != nil {
		return Result{}, err
	}
	c := cache.New()
	if !cfg.Force {
		c = cache.Load(cfg.OutDir)
	}

	res := Result{Topics: len(topics), Diagnostics: diags}
	bySlug := map[string]*model.Topic{}
	for _, t := range topics {
		bySlug[t.Slug] = t
	}

	// Build injects computed values (val/data shortcodes), so it depends on the
	// evidence checks — cached, so it stays cheap in an authoring loop.
	rr, _, evErr := evidence.Run(cfg.ContentDir, evidence.Options{Force: cfg.Force, Trust: cfg.Trust})
	if evErr != nil {
		res.Diagnostics = append(res.Diagnostics, model.Diagnostic{
			Severity: model.SevWarn, File: "evidence", Message: evErr.Error()})
	}
	checks := make(map[string]render.CheckResult, len(rr.Checks))
	for n, ck := range rr.Checks {
		checks[n] = render.CheckResult{Output: ck.Output, Err: ck.Err}
	}
	rd.SetChecks(checks)

	// Render topic pages (incremental).
	for _, t := range topics {
		fp := fingerprint(cfg, t, g, bySlug, eng.Hash(), layoutHash(layout), site.AssetVersion, evidenceHash(t, rr.Checks))
		outFile := filepath.Join(cfg.OutDir, filepath.FromSlash(t.OutPath))
		if !cfg.Force && c.PageFresh(t.OutPath, fp) && fileExists(outFile) {
			res.Skipped++
			c.PutPage(t.OutPath, fp)
			continue
		}
		htmlBytes, pdiags := rd.Page(t)
		res.Diagnostics = append(res.Diagnostics, pdiags...)
		if err := writeFile(outFile, htmlBytes); err != nil {
			return res, err
		}
		c.PutPage(t.OutPath, fp)
		res.Rendered++
	}

	// Auto-generated listing pages (index, sections, tags).
	listings := buildListings(rd, topics, bySlug, cfg)
	for _, lp := range listings {
		fp := hashStrings("listing", lp.fingerprint)
		if !cfg.Force && c.PageFresh(lp.out, fp) && fileExists(filepath.Join(cfg.OutDir, filepath.FromSlash(lp.out))) {
			res.Skipped++
			c.PutPage(lp.out, fp)
			continue
		}
		htmlBytes, err := rd.Listing(lp.title, lp.slug, lp.out, lp.intro, lp.cards)
		if err != nil {
			return res, err
		}
		if err := writeFile(filepath.Join(cfg.OutDir, filepath.FromSlash(lp.out)), htmlBytes); err != nil {
			return res, err
		}
		c.PutPage(lp.out, fp)
		res.Rendered++
	}

	// Static assets (built-in, then project theme overrides).
	if err := copyStatic(cfg, c); err != nil {
		return res, err
	}

	// Vendored third-party assets (KaTeX), only with --vendor. A failure here
	// degrades gracefully: pages still build, math just falls back to its raw TeX.
	if cfg.Vendor {
		if d := vendorAssets(cfg, c); d != nil {
			res.Diagnostics = append(res.Diagnostics, *d)
		}
	}

	if err := writeSearchIndex(cfg, topics, bySlug); err != nil {
		return res, err
	}

	res.Diagnostics = dedupeDiags(res.Diagnostics)
	if err := c.Save(cfg.OutDir); err != nil {
		return res, err
	}
	return res, nil
}

// --- fingerprint ---

// validateEvidenceRef checks a val/data shortcode without executing anything:
// the check= argument is required, and evidence/<check>/run must exist (a dead
// reference is a warning, like a dead wiki link).
func validateEvidenceRef(cfg Config, sc model.Shortcode, file string) []model.Diagnostic {
	name := strings.TrimSpace(sc.Args["check"])
	if name == "" {
		return []model.Diagnostic{{Severity: model.SevError, File: file, Line: sc.Line,
			Message: fmt.Sprintf("%s: missing required argument %q", sc.Name, "check")}}
	}
	runPath := filepath.Join(cfg.ContentDir, "evidence", name, "run")
	if fi, err := os.Stat(runPath); err != nil || fi.IsDir() {
		return []model.Diagnostic{{Severity: model.SevWarn, File: file, Line: sc.Line,
			Message: fmt.Sprintf("%s references evidence check %q, but evidence/%s/run does not exist", sc.Name, name, name)}}
	}
	return nil
}

// validateDataAs checks that a data shortcode's as= names something renderable:
// a built-in shape (table/bar/line/scatter) or a known template. It does not
// execute the check — it only catches a typo'd or missing template at check time,
// the way eng.Validate catches an unknown template invocation.
func validateDataAs(eng *tmpl.Engine, sc model.Shortcode, file string) []model.Diagnostic {
	as := strings.ToLower(strings.TrimSpace(sc.Args["as"]))
	switch as {
	case "", "table", "bar", "line", "scatter":
		return nil
	}
	if eng.Has(as) {
		return nil
	}
	return []model.Diagnostic{{Severity: model.SevError, File: file, Line: sc.Line,
		Message: fmt.Sprintf("data: as=%q is not a built-in shape or a known template", as)}}
}

// evidenceHash folds the outputs of the checks a page references into its
// fingerprint, so the page re-renders when a check's computed value changes even
// though its own markdown did not.
func evidenceHash(t *model.Topic, checks map[string]evidence.Check) string {
	var names []string
	for _, sc := range t.Shortcodes {
		if sc.Name == "val" || sc.Name == "data" {
			if n := strings.TrimSpace(sc.Args["check"]); n != "" {
				names = append(names, n)
			}
		}
	}
	if len(names) == 0 {
		return ""
	}
	sort.Strings(names)
	h := sha256.New()
	for _, n := range names {
		ck := checks[n]
		fmt.Fprintf(h, "E:%s|%s|%s\n", n, ck.Output, ck.Err)
	}
	return hex.EncodeToString(h.Sum(nil))
}

func fingerprint(cfg Config, t *model.Topic, g *graph.Graph, bySlug map[string]*model.Topic, engHash, layHash, assetVer, evHash string) string {
	h := sha256.New()
	fmt.Fprintf(h, "v3|%s|%s|%s|site=%s|assets=%s|ev=%s\n", t.ContentHash, engHash, layHash, cfg.SiteTitle, assetVer, evHash)
	// outbound resolved targets affect rendering (path + dead state)
	for _, l := range t.Links {
		fmt.Fprintf(h, "L:%s>%s:%v\n", l.Target, l.ResolvedSlug, l.Dead)
		if dst, ok := bySlug[l.ResolvedSlug]; ok {
			fmt.Fprintf(h, "Lt:%s|%s\n", dst.OutPath, dst.Title)
		}
	}
	for _, rel := range g.Related(t.Slug, 6) {
		if dst, ok := bySlug[rel.Slug]; ok {
			fmt.Fprintf(h, "R:%s|%s|%s\n", dst.OutPath, dst.Title, dst.Summary)
		}
	}
	bl := append([]string(nil), g.Backlinks(t.Slug)...)
	sort.Strings(bl)
	for _, s := range bl {
		if dst, ok := bySlug[s]; ok {
			fmt.Fprintf(h, "B:%s|%s\n", dst.OutPath, dst.Title)
		}
	}
	// The sidebar shows this page's connected component, so its membership
	// (which sections/topics appear) is part of the rendered output.
	comp := g.Component(t.Slug)
	var members []string
	for s, c := range g.Components {
		if c == comp {
			if dst, ok := bySlug[s]; ok {
				members = append(members, dst.OutPath+"|"+dst.Title)
			}
		}
	}
	sort.Strings(members)
	fmt.Fprintf(h, "C:%d:%s\n", comp, strings.Join(members, ";"))
	return hex.EncodeToString(h.Sum(nil))
}

func layoutHash(layout map[string]string) string {
	keys := make([]string, 0, len(layout))
	for k := range layout {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s:%s\n", k, layout[k])
	}
	return hex.EncodeToString(h.Sum(nil))
}

func hashStrings(parts ...string) string {
	h := sha256.New()
	for _, p := range parts {
		h.Write([]byte(p))
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))
}

// --- listing pages ---

type listingPage struct {
	title, slug, out, intro string
	cards                   []render.Card
	fingerprint             string
}

func buildListings(rd *render.Renderer, topics []*model.Topic, bySlug map[string]*model.Topic, cfg Config) []listingPage {
	var pages []listingPage

	// Group topics by their parent directory for section indexes.
	bySection := map[string][]*model.Topic{}
	hasIndex := map[string]bool{}
	var rootIndex bool
	for _, t := range topics {
		if t.Slug == "index" {
			rootIndex = true
		}
		dir := path.Dir(t.Slug)
		if dir == "." {
			dir = ""
		}
		if strings.HasSuffix(t.OutPath, "/index.html") && t.Slug != "index" {
			hasIndex[t.Slug] = true
			continue
		}
		bySection[dir] = append(bySection[dir], t)
	}

	card := func(t *model.Topic) render.Card {
		return render.Card{Title: t.Title, URL: t.OutPath, Summary: t.Summary}
	}

	// Root index, only if the author didn't write one.
	if !rootIndex {
		var cards []render.Card
		var fp strings.Builder
		for _, t := range sortedTopics(bySection[""]) {
			cards = append(cards, card(t))
			fp.WriteString(t.OutPath + t.Title + t.Summary + ";")
		}
		// also surface top-level sections as cards
		for sec := range bySection {
			if sec == "" || strings.Contains(sec, "/") {
				continue
			}
			if idx, ok := bySlug[sec]; ok {
				cards = append(cards, card(idx))
			}
		}
		pages = append(pages, listingPage{
			title: cfg.SiteTitle, slug: "index", out: "index.html",
			intro: "Browse the knowledge base.", cards: cards, fingerprint: fp.String(),
		})
	}

	// Section indexes for directories lacking an index.md.
	for dir, members := range bySection {
		if dir == "" || hasIndex[dir] {
			continue
		}
		var cards []render.Card
		var fp strings.Builder
		for _, t := range sortedTopics(members) {
			cards = append(cards, card(t))
			fp.WriteString(t.OutPath + t.Title + t.Summary + ";")
		}
		pages = append(pages, listingPage{
			title: humanizeSeg(path.Base(dir)), slug: dir, out: dir + "/index.html",
			intro: "", cards: cards, fingerprint: fp.String(),
		})
	}

	// Label pages: one worklist page per label (e.g. labels/open-problem.html),
	// plus an index. Labels are workflow markers, so the listing is the reader-
	// facing form of the `starbase labels` worklist.
	labelTopics := map[string][]*model.Topic{}
	for _, t := range topics {
		for _, l := range t.Labels {
			labelTopics[l] = append(labelTopics[l], t)
		}
	}
	if len(labelTopics) > 0 {
		var idxCards []render.Card
		var idxFP strings.Builder
		labels := make([]string, 0, len(labelTopics))
		for l := range labelTopics {
			labels = append(labels, l)
		}
		sort.Strings(labels)
		for _, l := range labels {
			slug := parse.Slugify(l)
			var cards []render.Card
			var fp strings.Builder
			for _, t := range sortedTopics(labelTopics[l]) {
				cards = append(cards, card(t))
				fp.WriteString(t.OutPath + t.Title + ";")
			}
			pages = append(pages, listingPage{
				title: "⚑ " + l, slug: "labels/" + slug, out: "labels/" + slug + ".html",
				intro: fmt.Sprintf("Topics labeled %q.", l), cards: cards, fingerprint: fp.String(),
			})
			idxCards = append(idxCards, render.Card{Title: "⚑ " + l, URL: slug + ".html",
				Summary: fmt.Sprintf("%d topics", len(labelTopics[l]))})
			idxFP.WriteString(l + fmt.Sprint(len(labelTopics[l])) + ";")
		}
		pages = append(pages, listingPage{
			title: "Labels", slug: "labels", out: "labels/index.html",
			intro: "All labels.", cards: idxCards, fingerprint: idxFP.String(),
		})
	}

	// Tag pages.
	tagTopics := map[string][]*model.Topic{}
	for _, t := range topics {
		for _, tag := range t.Tags {
			tagTopics[tag] = append(tagTopics[tag], t)
		}
	}
	if len(tagTopics) > 0 {
		var idxCards []render.Card
		var idxFP strings.Builder
		tags := make([]string, 0, len(tagTopics))
		for tag := range tagTopics {
			tags = append(tags, tag)
		}
		sort.Strings(tags)
		for _, tag := range tags {
			slug := parse.Slugify(tag)
			out := "tags/" + slug + ".html"
			var cards []render.Card
			var fp strings.Builder
			for _, t := range sortedTopics(tagTopics[tag]) {
				cards = append(cards, card(t))
				fp.WriteString(t.OutPath + t.Title + ";")
			}
			pages = append(pages, listingPage{
				title: "#" + tag, slug: "tags/" + slug, out: out,
				intro: fmt.Sprintf("Topics tagged %q.", tag), cards: cards, fingerprint: fp.String(),
			})
			idxCards = append(idxCards, render.Card{Title: "#" + tag, URL: slug + ".html",
				Summary: fmt.Sprintf("%d topics", len(tagTopics[tag]))})
			idxFP.WriteString(tag + fmt.Sprint(len(tagTopics[tag])) + ";")
		}
		pages = append(pages, listingPage{
			title: "Tags", slug: "tags", out: "tags/index.html",
			intro: "All tags.", cards: idxCards, fingerprint: idxFP.String(),
		})
	}
	return pages
}

func sortedTopics(ts []*model.Topic) []*model.Topic {
	out := append([]*model.Topic(nil), ts...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Weight != out[j].Weight {
			return out[i].Weight < out[j].Weight
		}
		return out[i].Title < out[j].Title
	})
	return out
}

// --- engine / layout / assets loading ---

func loadEngine(cfg Config) (*tmpl.Engine, []model.Diagnostic, error) {
	eng := tmpl.New()
	if err := eng.LoadFS(assets.FS, "templates"); err != nil {
		return nil, nil, fmt.Errorf("loading built-in templates: %w", err)
	}
	var diags []model.Diagnostic
	projDir := filepath.Join(cfg.ContentDir, "templates")
	if isDir(projDir) {
		if err := eng.LoadFS(os.DirFS(projDir), "."); err != nil {
			diags = append(diags, model.Diagnostic{
				Severity: model.SevError, File: "templates", Message: err.Error()})
		}
	}
	return eng, diags, nil
}

func loadLayout(cfg Config) (map[string]string, error) {
	layout, err := render.LoadLayout(assets.FS, "layout")
	if err != nil {
		return nil, err
	}
	projDir := filepath.Join(cfg.ContentDir, "layout")
	if isDir(projDir) {
		over, err := render.LoadLayout(os.DirFS(projDir), ".")
		if err != nil {
			return nil, err
		}
		for k, v := range over {
			layout[k] = v
		}
	}
	return layout, nil
}

// staticVersion hashes all static assets (built-in + project theme overrides)
// into a short token appended to asset URLs, so browsers fetch fresh CSS/JS
// whenever they change and cache aggressively when they don't.
func staticVersion(cfg Config) string {
	h := sha256.New()
	files := map[string][]byte{}
	fs.WalkDir(assets.FS, "static", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, _ := assets.FS.ReadFile(p)
		files[strings.TrimPrefix(p, "static/")] = b
		return nil
	})
	themeDir := filepath.Join(cfg.ContentDir, "theme")
	if isDir(themeDir) {
		if ents, err := os.ReadDir(themeDir); err == nil {
			for _, e := range ents {
				if !e.IsDir() {
					b, _ := os.ReadFile(filepath.Join(themeDir, e.Name()))
					files[e.Name()] = b
				}
			}
		}
	}
	names := make([]string, 0, len(files))
	for n := range files {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		fmt.Fprintf(h, "%s:%x\n", n, sha256.Sum256(files[n]))
	}
	// Math source (CDN vs vendored, and KaTeX version) affects page output too.
	fmt.Fprintf(h, "katex:%s:vendor=%v\n", vendor.KaTeXVersion, cfg.Vendor)
	return hex.EncodeToString(h.Sum(nil))[:10]
}

// writeSearchIndex emits search.json: a flat list of every topic (title, url,
// section, summary) so the sidebar's search box can find anything, across all
// connected components.
func writeSearchIndex(cfg Config, topics []*model.Topic, bySlug map[string]*model.Topic) error {
	type entry struct {
		T string `json:"t"`
		U string `json:"u"`
		S string `json:"s"`
		D string `json:"d"`
	}
	list := make([]entry, 0, len(topics))
	for _, t := range topics {
		if t.Slug == graph.HubSlug {
			continue
		}
		sec := ""
		if dir := path.Dir(t.Slug); dir != "." {
			if idx, ok := bySlug[dir]; ok {
				sec = idx.Title
			} else {
				sec = humanizeSeg(path.Base(dir))
			}
		} else if strings.HasSuffix(t.OutPath, "/index.html") {
			sec = t.Title // a section landing page
		}
		list = append(list, entry{T: t.Title, U: t.OutPath, S: sec, D: t.Summary})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].T < list[j].T })
	b, err := json.Marshal(list)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(cfg.OutDir, "search.json"), b, 0o644)
}

// writeStatic writes one asset under <out>/static/<name>, skipping it when the
// cache shows the same content is already on disk.
func writeStatic(cfg Config, c *cache.Cache, name string, content []byte) error {
	sum := sha256.Sum256(content)
	h := hex.EncodeToString(sum[:])
	dst := filepath.Join(cfg.OutDir, "static", filepath.FromSlash(name))
	if c.AssetFresh(name, h) && fileExists(dst) {
		return nil
	}
	if err := writeFile(dst, content); err != nil {
		return err
	}
	c.PutAsset(name, h)
	return nil
}

func copyStatic(cfg Config, c *cache.Cache) error {
	err := fs.WalkDir(assets.FS, "static", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, e := assets.FS.ReadFile(p)
		if e != nil {
			return e
		}
		return writeStatic(cfg, c, strings.TrimPrefix(p, "static/"), b)
	})
	if err != nil {
		return err
	}
	// project theme overrides
	themeDir := filepath.Join(cfg.ContentDir, "theme")
	if isDir(themeDir) {
		ents, _ := os.ReadDir(themeDir)
		for _, e := range ents {
			if e.IsDir() {
				continue
			}
			b, err := os.ReadFile(filepath.Join(themeDir, e.Name()))
			if err != nil {
				return err
			}
			if err := writeStatic(cfg, c, e.Name(), b); err != nil {
				return err
			}
		}
	}
	return nil
}

// vendorAssets downloads (or loads from cache) the third-party front-end assets
// and writes them under <out>/static. Returns a warning diagnostic on failure
// rather than aborting the build.
func vendorAssets(cfg Config, c *cache.Cache) *model.Diagnostic {
	files, err := vendor.EnsureKaTeX(cfg.Offline)
	if err != nil {
		return &model.Diagnostic{
			Severity: model.SevWarn, File: "(vendor)",
			Message: fmt.Sprintf("could not vendor KaTeX; math will not render: %v", err),
		}
	}
	for _, f := range files {
		if err := writeStatic(cfg, c, f.RelPath, f.Content); err != nil {
			return &model.Diagnostic{
				Severity: model.SevWarn, File: "(vendor)",
				Message: fmt.Sprintf("writing vendored asset %s: %v", f.RelPath, err),
			}
		}
	}
	return nil
}

// --- file walking + io ---

func collectMarkdown(cfg Config) ([]string, error) {
	var files []string
	root := cfg.ContentDir
	absOut, _ := filepath.Abs(cfg.OutDir)
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if p != root && (strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") || ignoredDirs[name]) {
				return filepath.SkipDir
			}
			if abs, _ := filepath.Abs(p); abs == absOut {
				return filepath.SkipDir
			}
			return nil
		}
		name := d.Name()
		if strings.HasSuffix(name, ".md") && !strings.HasPrefix(name, ".") &&
			(!strings.HasPrefix(name, "_") || name == "_index.md") {
			rel, _ := filepath.Rel(root, p)
			files = append(files, filepath.ToSlash(rel))
		}
		return nil
	})
	sort.Strings(files)
	return files, err
}

func writeFile(p string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	return os.WriteFile(p, content, 0o644)
}

func fileExists(p string) bool { _, err := os.Stat(p); return err == nil }
func isDir(p string) bool      { fi, err := os.Stat(p); return err == nil && fi.IsDir() }

func dedupeDiags(diags []model.Diagnostic) []model.Diagnostic {
	seen := map[string]bool{}
	var out []model.Diagnostic
	for _, d := range diags {
		k := fmt.Sprintf("%d|%s|%d|%s", d.Severity, d.File, d.Line, d.Message)
		if seen[k] {
			continue
		}
		seen[k] = true
		out = append(out, d)
	}
	return out
}

func humanizeSeg(s string) string {
	s = strings.ReplaceAll(strings.ReplaceAll(s, "-", " "), "_", " ")
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
