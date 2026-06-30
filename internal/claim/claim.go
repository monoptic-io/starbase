// Package claim implements evidence-backed claims: a statement whose value is
// produced by a computation the authoring agent already ran in its sandbox. The
// claim carries the result and the implementation (a query/script, shown
// verbatim) so a reader can dig down to how a number was derived.
//
// starbase never executes anything — the agent does. starbase's job is to
// surface the computation and, like a dead wiki link, to flag a claim that has
// no evidence attached so the agent swarm converges on supporting it.
package claim

import (
	"fmt"
	"strings"

	"github.com/monoptic-io/starbase/internal/model"
)

// Info is a parsed claim: the prose statement plus its attached evidence.
type Info struct {
	Prose     string // the claim text (markdown, fences removed)
	Value     string // the asserted value (should equal the computed result)
	Lang      string // implementation language (for highlighting)
	Code      string // the implementation: a query/script, verbatim
	HasImpl   bool
	Result    string // the captured result (raw text: scalar, CSV, or table)
	ResultFmt string // "result" | "csv" | "tsv" | "table" | ...
	Source    string // dataset / notebook / file the computation ran against
	AsOf      string // when the computation was run
}

// Parse pulls a claim apart. Inside the shortcode body, the first non-result
// fenced code block is the implementation; a fence tagged result/csv/tsv/table
// is the captured result; everything else is the prose statement.
func Parse(sc model.Shortcode) Info {
	in := Info{
		Value:  strings.TrimSpace(sc.Args["value"]),
		Source: strings.TrimSpace(sc.Args["source"]),
		AsOf:   strings.TrimSpace(sc.Args["asof"]),
	}
	prose, fences := splitFences(sc.Inner)
	for _, f := range fences {
		key := strings.ToLower(firstWord(f.info))
		if isResultTag(key) {
			if in.Result == "" {
				in.Result, in.ResultFmt = strings.TrimRight(f.body, "\n"), key
			}
			continue
		}
		if !in.HasImpl {
			in.Code, in.HasImpl = strings.TrimRight(f.body, "\n"), true
			in.Lang = firstNonEmpty(key, strings.TrimSpace(sc.Args["lang"]))
		}
	}
	in.Prose = strings.TrimSpace(prose)
	return in
}

// Validate surfaces the research-loop diagnostics for a claim.
func Validate(in Info, file string, line int) []model.Diagnostic {
	var diags []model.Diagnostic
	if !in.HasImpl && in.Source == "" {
		diags = append(diags, model.Diagnostic{
			Severity: model.SevWarn, File: file, Line: line,
			Message: "unsupported claim: attach an implementation (a code/query block) or a source",
		})
	}
	if in.Value != "" && in.Result != "" {
		if s := Scalar(in.Result); s != "" && !strings.EqualFold(s, in.Value) {
			diags = append(diags, model.Diagnostic{
				Severity: model.SevWarn, File: file, Line: line,
				Message: fmt.Sprintf("claim value %q does not match its result %q", in.Value, s),
			})
		}
	}
	return diags
}

// Scalar extracts a single value from a result when it looks like one (a lone
// value, or a one-column/one-row table), else "".
func Scalar(result string) string {
	lines := nonEmptyLines(result)
	if len(lines) == 0 {
		return ""
	}
	cells := splitCells(lines[len(lines)-1])
	if len(lines) <= 2 && len(cells) == 1 {
		return strings.TrimSpace(cells[len(cells)-1])
	}
	return ""
}

// Rows turns a result block into a header + rows for table rendering. The
// second return is true when the first row looks like a header.
func Rows(result, format string) [][]string {
	var rows [][]string
	for _, ln := range nonEmptyLines(result) {
		rows = append(rows, splitCells(ln))
	}
	return rows
}

// --- helpers ---

type fence struct{ info, body string }

func splitFences(s string) (string, []fence) {
	lines := strings.Split(s, "\n")
	var prose []string
	var fences []fence
	for i := 0; i < len(lines); {
		t := strings.TrimSpace(lines[i])
		if strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~") {
			marker := t[:3]
			info := strings.TrimSpace(t[3:])
			var body []string
			i++
			for i < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[i]), marker) {
				body = append(body, lines[i])
				i++
			}
			i++ // skip the closing fence
			fences = append(fences, fence{info: info, body: strings.Join(body, "\n")})
			continue
		}
		prose = append(prose, lines[i])
		i++
	}
	return strings.Join(prose, "\n"), fences
}

func isResultTag(s string) bool {
	switch s {
	case "result", "results", "output", "csv", "tsv", "table", "rows":
		return true
	}
	return false
}

func splitCells(line string) []string {
	line = strings.TrimSpace(strings.Trim(line, "|"))
	var parts []string
	switch {
	case strings.Contains(line, "|"):
		parts = strings.Split(line, "|")
	case strings.Contains(line, "\t"):
		parts = strings.Split(line, "\t")
	case strings.Contains(line, ","):
		parts = strings.Split(line, ",")
	default:
		parts = strings.Fields(line)
	}
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func nonEmptyLines(s string) []string {
	var out []string
	for _, ln := range strings.Split(s, "\n") {
		t := strings.TrimSpace(ln)
		// markdown table separator row, e.g. |---|---|
		if t == "" || strings.Trim(t, "-|: ") == "" {
			continue
		}
		out = append(out, ln)
	}
	return out
}

func firstWord(s string) string {
	f := strings.Fields(s)
	if len(f) == 0 {
		return ""
	}
	return f[0]
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
