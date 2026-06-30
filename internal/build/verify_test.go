package build

import (
	"testing"

	"github.com/monoptic-io/starbase/internal/claim"
	"github.com/monoptic-io/starbase/internal/evidence"
)

func TestValueEqual(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{"11,400,000", "11400000", true},
		{"4", "4.0", true},
		{"$11.4", "11.4", true},
		{"Midwest", "midwest", true},
		{"4", "5", false},
		{"Midwest", "West", false},
	}
	for _, c := range cases {
		if got := valueEqual(c.a, c.b); got != c.want {
			t.Errorf("valueEqual(%q,%q)=%v want %v", c.a, c.b, got, c.want)
		}
	}
}

func TestTableEqualNumberFormatting(t *testing.T) {
	a := [][]string{{"division", "total"}, {"Midwest", "11400000"}}
	b := [][]string{{"division", "total"}, {"Midwest", "11,400,000"}}
	if !tableEqual(a, b) {
		t.Error("tables differing only in number formatting should match")
	}
	b[1][1] = "9000000"
	if tableEqual(a, b) {
		t.Error("tables with different values must not match")
	}
}

func TestCompareClaim(t *testing.T) {
	// matching scalar
	if msg := compareClaim(claim.Info{Check: "x", Value: "4"}, evidence.Check{Value: "4"}); msg != "" {
		t.Errorf("expected match, got %q", msg)
	}
	// a fabricated value is caught
	if msg := compareClaim(claim.Info{Check: "x", Value: "5"}, evidence.Check{Value: "4"}); msg == "" {
		t.Error("expected mismatch for fabricated value")
	}
	// a propagated error from the check
	if msg := compareClaim(
		claim.Info{Check: "x", Result: "division,total\nMidwest,11400000", ResultFmt: "csv"},
		evidence.Check{Table: [][]string{{"division", "total"}, {"Midwest", "9000000"}}}); msg == "" {
		t.Error("expected mismatch for wrong table")
	}
}
