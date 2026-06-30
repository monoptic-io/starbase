// Package parse turns a markdown file on disk into a model.Topic: it pulls out
// frontmatter, outbound wiki links, and template (shortcode) invocations.
//
// Extraction is deliberately code-aware: fenced code blocks and inline code
// spans are masked before scanning, so `[[not a link]]` inside a code sample is
// never mistaken for a real wiki link.
package parse

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/monoptic-io/starbase/internal/model"
	yaml "gopkg.in/yaml.v2"
)

var (
	reWikiLink   = regexp.MustCompile(`\[\[([^\]\n|]+)(?:\|([^\]\n]+))?\]\]`)
	reTag        = regexp.MustCompile(`\{\{<\s*(/?)\s*([\w-]+)((?:[^>]|>(?:[^}]|$))*?)(/?)>\}\}`)
	reArg        = regexp.MustCompile(`([\w-]+)\s*=\s*("([^"]*)"|'([^']*)'|(\S+))`)
	reH1         = regexp.MustCompile(`(?m)^#{1}\s+(.+?)\s*$`)
	reInlineCode = regexp.MustCompile("`[^`\n]*`")
)

// File reads and parses a single markdown file. root is the content root and
// rel is the file's path relative to it (slash-separated).
func File(root, rel string) (*model.Topic, []model.Diagnostic, error) {
	abs := path.Join(root, rel)
	raw, err := os.ReadFile(abs)
	if err != nil {
		return nil, nil, err
	}
	return Bytes(rel, raw)
}

// Bytes parses already-read file content. Exposed for testing and caching.
func Bytes(rel string, raw []byte) (*model.Topic, []model.Diagnostic, error) {
	sum := sha256.Sum256(raw)
	t := &model.Topic{
		SourcePath:  rel,
		ContentHash: hex.EncodeToString(sum[:]),
	}
	t.Slug, t.OutPath = pathToSlug(rel)

	body, fm := splitFrontmatter(string(raw))
	t.Body = body
	var diags []model.Diagnostic
	if fm != "" {
		var meta map[string]any
		if err := yaml.Unmarshal([]byte(fm), &meta); err != nil {
			diags = append(diags, model.Diagnostic{
				Severity: model.SevError, File: rel, Line: 1,
				Message: fmt.Sprintf("invalid frontmatter: %v", err),
			})
		} else {
			applyFrontmatter(t, meta)
		}
	}

	masked := maskCode(body)

	// Title falls back to the first H1, then a humanized filename.
	if t.Title == "" {
		if m := reH1.FindStringSubmatch(masked); m != nil {
			// recover the real text from the unmasked body at the same span
			loc := reH1.FindStringSubmatchIndex(masked)
			t.Title = strings.TrimSpace(body[loc[2]:loc[3]])
		}
	}
	if t.Title == "" {
		t.Title = humanize(strings.TrimSuffix(path.Base(rel), ".md"))
	}
	if t.Summary == "" {
		t.Summary = firstParagraph(masked, body)
	}

	t.Links = extractLinks(masked, body, rel)
	t.Shortcodes = extractShortcodes(masked, body)
	t.WordCount = len(strings.Fields(stripMarkup(masked)))

	return t, diags, nil
}

// ScanLinks finds wiki links in an arbitrary markdown string (used by the
// renderer for nested/inner content). Offsets are relative to s.
func ScanLinks(s string) []model.Link {
	return extractLinks(maskCode(s), s, "")
}

// ScanShortcodes finds shortcode invocations in an arbitrary markdown string.
func ScanShortcodes(s string) []model.Shortcode {
	return extractShortcodes(maskCode(s), s)
}

func extractLinks(masked, orig, rel string) []model.Link {
	var links []model.Link
	for _, loc := range reWikiLink.FindAllStringSubmatchIndex(masked, -1) {
		target := strings.TrimSpace(orig[loc[2]:loc[3]])
		display := ""
		if loc[4] >= 0 {
			display = strings.TrimSpace(orig[loc[4]:loc[5]])
		}
		links = append(links, model.Link{
			Target:  target,
			Display: display,
			Line:    lineOf(masked, loc[0]),
			Start:   loc[0],
			End:     loc[1],
		})
	}
	return links
}

// extractShortcodes tokenizes shortcode tags and pairs them with a stack, the
// way Hugo does. An opening tag becomes a block (paired) shortcode if a matching
// close tag appears later; otherwise it is a self-closing invocation. This is
// robust to interleaving that a single regex cannot express (RE2 has no
// backreferences).
func extractShortcodes(masked, orig string) []model.Shortcode {
	type openTag struct {
		name             string
		argStart, argEnd int
		tagStart, tagEnd int
	}
	var stack []openTag
	var codes []model.Shortcode

	mk := func(name string, argS, argE, start, end int, inner string) model.Shortcode {
		return model.Shortcode{
			Name: name, Args: parseArgs(orig[argS:argE]), Inner: inner,
			Line: lineOf(masked, start), Start: start, End: end, Raw: orig[start:end],
		}
	}

	for _, loc := range reTag.FindAllStringSubmatchIndex(masked, -1) {
		isClose := loc[3] > loc[2] // group 1: leading slash
		isSelf := loc[9] > loc[8]  // group 4: trailing slash
		name := orig[loc[4]:loc[5]]
		argS, argE := loc[6], loc[7]
		start, end := loc[0], loc[1]

		switch {
		case isClose:
			idx := -1
			for i := len(stack) - 1; i >= 0; i-- {
				if stack[i].name == name {
					idx = i
					break
				}
			}
			if idx < 0 {
				continue // stray closing tag
			}
			for i := idx + 1; i < len(stack); i++ { // unclosed inner opens
				o := stack[i]
				codes = append(codes, mk(o.name, o.argStart, o.argEnd, o.tagStart, o.tagEnd, ""))
			}
			o := stack[idx]
			codes = append(codes, mk(o.name, o.argStart, o.argEnd, o.tagStart, end, orig[o.tagEnd:start]))
			stack = stack[:idx]
		case isSelf:
			codes = append(codes, mk(name, argS, argE, start, end, ""))
		default:
			stack = append(stack, openTag{name, argS, argE, start, end})
		}
	}
	for _, o := range stack { // opens with no close are self-closing
		codes = append(codes, mk(o.name, o.argStart, o.argEnd, o.tagStart, o.tagEnd, ""))
	}
	sort.Slice(codes, func(i, j int) bool { return codes[i].Start < codes[j].Start })
	return codes
}

func parseArgs(s string) map[string]string {
	args := map[string]string{}
	for _, m := range reArg.FindAllStringSubmatch(s, -1) {
		key := m[1]
		var val string
		switch {
		case m[3] != "" || strings.HasPrefix(strings.TrimSpace(m[2]), `"`):
			val = m[3]
		case m[4] != "" || strings.HasPrefix(strings.TrimSpace(m[2]), `'`):
			val = m[4]
		default:
			val = m[5]
		}
		args[key] = val
	}
	return args
}

// --- frontmatter ---

func splitFrontmatter(s string) (body, fm string) {
	if !strings.HasPrefix(s, "---\n") && !strings.HasPrefix(s, "---\r\n") {
		return s, ""
	}
	rest := s[strings.IndexByte(s, '\n')+1:]
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return s, ""
	}
	fm = rest[:end]
	after := rest[end+4:] // past "\n---"
	if i := strings.IndexByte(after, '\n'); i >= 0 {
		body = after[i+1:]
	}
	return body, fm
}

func applyFrontmatter(t *model.Topic, meta map[string]any) {
	t.Title = str(meta["title"])
	t.Summary = first(str(meta["summary"]), str(meta["description"]))
	t.Aliases = strList(meta["aliases"], meta["alias"])
	t.Tags = strList(meta["tags"], meta["tag"])
	if d, ok := meta["draft"].(bool); ok {
		t.Draft = d
	}
	if w, ok := meta["weight"].(int); ok {
		t.Weight = w
	}
}

// --- code masking (length-preserving so offsets map back to the original) ---

// MaskCode replaces fenced and inline code with same-length blanks so callers
// can scan for markup without matching inside code samples. Byte offsets in the
// result map exactly onto the original string.
func MaskCode(s string) string { return maskCode(s) }

func maskCode(s string) string {
	lines := strings.Split(s, "\n")
	inFence := false
	marker := ""
	for i, line := range lines {
		trimmed := strings.TrimLeft(line, " ")
		switch {
		case inFence:
			ended := strings.HasPrefix(trimmed, marker)
			lines[i] = blank(line)
			if ended {
				inFence = false
			}
		case strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~"):
			inFence = true
			marker = trimmed[:3]
			lines[i] = blank(line)
		default:
			lines[i] = reInlineCode.ReplaceAllStringFunc(line, blank)
		}
	}
	return strings.Join(lines, "\n")
}

func blank(s string) string {
	b := make([]byte, len(s))
	for i := range b {
		if s[i] == '\n' {
			b[i] = '\n'
		} else {
			b[i] = ' '
		}
	}
	return string(b)
}

// --- helpers ---

func lineOf(s string, off int) int {
	return strings.Count(s[:off], "\n") + 1
}

func pathToSlug(rel string) (slug, out string) {
	rel = strings.TrimSuffix(rel, ".md")
	base := path.Base(rel)
	dir := path.Dir(rel)
	if base == "index" || base == "_index" || base == "README" || base == "readme" {
		if dir == "." || dir == "" {
			return "index", "index.html"
		}
		return slugifyPath(dir), slugifyPath(dir) + "/index.html"
	}
	slug = slugifyPath(rel)
	return slug, slug + ".html"
}

func slugifyPath(p string) string {
	parts := strings.Split(p, "/")
	for i, seg := range parts {
		parts[i] = Slugify(seg)
	}
	return strings.Join(parts, "/")
}

var reNonSlug = regexp.MustCompile(`[^a-z0-9]+`)

// Slugify produces a stable, url-safe anchor/path segment.
func Slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = reNonSlug.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

func humanize(s string) string {
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func firstParagraph(masked, orig string) string {
	for _, para := range strings.Split(masked, "\n\n") {
		p := strings.TrimSpace(para)
		if p == "" || strings.HasPrefix(p, "#") || strings.HasPrefix(p, "{{") || strings.HasPrefix(p, "|") {
			continue
		}
		loc := strings.Index(masked, para)
		text := stripMarkup(orig[loc : loc+len(para)])
		text = strings.Join(strings.Fields(text), " ")
		if len(text) > 200 {
			text = text[:197] + "..."
		}
		return text
	}
	return ""
}

var reMarkup = regexp.MustCompile(`[#>*_\x60\[\]]|\]\(|\(http\S+\)`)

func stripMarkup(s string) string {
	s = reWikiLink.ReplaceAllString(s, "$1")
	s = reMarkup.ReplaceAllString(s, " ")
	return s
}

func str(v any) string {
	if s, ok := v.(string); ok {
		return strings.TrimSpace(s)
	}
	return ""
}

func first(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func strList(vs ...any) []string {
	var out []string
	for _, v := range vs {
		switch x := v.(type) {
		case []any:
			for _, e := range x {
				if s := str(e); s != "" {
					out = append(out, s)
				}
			}
		case string:
			for _, part := range strings.Split(x, ",") {
				if s := strings.TrimSpace(part); s != "" {
					out = append(out, s)
				}
			}
		}
	}
	return out
}
