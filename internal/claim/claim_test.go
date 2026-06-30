package claim

import (
	"strings"
	"testing"

	"github.com/monoptic-io/starbase/internal/model"
)

func sc(inner string, args map[string]string) model.Shortcode {
	return model.Shortcode{Name: "claim", Inner: inner, Args: args}
}

func TestParseExtractsImplAndResult(t *testing.T) {
	in := Parse(sc(
		"There are **4** regions.\n```sql\nSELECT count(*) FROM r;\n```\n```result\nn\n4\n```\n",
		map[string]string{"value": "4", "source": "r.csv"}))
	if !in.HasImpl || in.Lang != "sql" {
		t.Fatalf("implementation not parsed: %+v", in)
	}
	if !strings.Contains(in.Code, "count(*)") {
		t.Errorf("code missing query: %q", in.Code)
	}
	if in.Result == "" {
		t.Error("result not parsed")
	}
	if in.Prose == "" || strings.Contains(in.Prose, "```") {
		t.Errorf("prose should be fence-free: %q", in.Prose)
	}
	if in.Value != "4" || in.Source != "r.csv" {
		t.Errorf("args: %+v", in)
	}
}

func TestUnsupportedClaimWarns(t *testing.T) {
	d := Validate(Parse(sc("A bare claim with no evidence.", nil)), "f.md", 1)
	if len(d) != 1 || d[0].Severity != model.SevWarn {
		t.Fatalf("want one unsupported warning, got %+v", d)
	}
}

func TestSupportedClaimIsClean(t *testing.T) {
	d := Validate(Parse(sc("X.\n```sql\nSELECT 1;\n```", nil)), "f.md", 1)
	if len(d) != 0 {
		t.Fatalf("supported claim should be clean, got %+v", d)
	}
}

func TestValueDriftWarns(t *testing.T) {
	d := Validate(Parse(sc(
		"X.\n```sql\nSELECT count(*);\n```\n```result\nn\n5\n```",
		map[string]string{"value": "4"})), "f.md", 1)
	if len(d) != 1 {
		t.Fatalf("value 4 vs result 5 should warn, got %+v", d)
	}
}
