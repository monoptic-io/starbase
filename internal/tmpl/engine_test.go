package tmpl

import (
	"strings"
	"testing"

	"github.com/ryannedolan/sitegen/internal/model"
)

func newEngine(t *testing.T, name, src string) *Engine {
	e := New()
	if err := e.register(name, src); err != nil {
		t.Fatal(err)
	}
	return e
}

func TestRequiredArgIsError(t *testing.T) {
	e := newEngine(t, "chart", `{{/* @param data required */}}{{/* @param color default="#fff" */}}<i>{{.data}}</i>`)

	sc := model.Shortcode{Name: "chart", Args: map[string]string{}}
	d := e.Validate(sc, "f.md")
	if len(d) != 1 || d[0].Severity != model.SevError {
		t.Fatalf("expected one error for missing required arg, got %+v", d)
	}

	sc.Args["data"] = "1,2,3"
	if d := e.Validate(sc, "f.md"); len(d) != 0 {
		t.Fatalf("expected no diagnostics, got %+v", d)
	}
}

func TestUnknownTemplateAndArg(t *testing.T) {
	e := newEngine(t, "note", `{{/* @param kind default="note" */}}x`)

	if d := e.Validate(model.Shortcode{Name: "missing"}, "f.md"); len(d) != 1 || d[0].Severity != model.SevError {
		t.Fatalf("unknown template should error, got %+v", d)
	}
	d := e.Validate(model.Shortcode{Name: "note", Args: map[string]string{"typo": "x"}}, "f.md")
	if len(d) != 1 || d[0].Severity != model.SevWarn {
		t.Fatalf("unknown arg should warn, got %+v", d)
	}
}

func TestOpenParamsSuppressesUnknownArgWarning(t *testing.T) {
	e := newEngine(t, "sim", `{{/* @params open */}}{{/* @param name required */}}x`)
	d := e.Validate(model.Shortcode{Name: "sim", Args: map[string]string{"name": "p", "length": "2"}}, "f.md")
	if len(d) != 0 {
		t.Fatalf("open params should accept extra args, got %+v", d)
	}
}

func TestRenderUsesDefaultsAndID(t *testing.T) {
	e := newEngine(t, "box", `{{/* @param color default="red" */}}<b id="{{.ID}}">{{.color}}/{{.Page.Title}}</b>`)
	out, d := e.Render(model.Shortcode{Name: "box", Args: map[string]string{}}, "", "", PageContext{Title: "T"}, "id7", "f.md")
	if len(d) != 0 {
		t.Fatalf("unexpected diags: %+v", d)
	}
	if s := string(out); !strings.Contains(s, `id="id7"`) || !strings.Contains(s, "red/T") {
		t.Fatalf("render = %q", s)
	}
}
