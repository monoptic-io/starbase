package build

import (
	"testing"

	"github.com/monoptic-io/starbase/internal/claim"
	"github.com/monoptic-io/starbase/internal/evidence"
)

func TestNormalizeText(t *testing.T) {
	cases := [][2]string{
		{"4\n", "4"},
		{"  4  ", "4"},
		{"a,b \nc,d\t\n", "a,b\nc,d"},
		{"\r\nx\r\n", "x"},
	}
	for _, c := range cases {
		if got := normalizeText(c[0]); got != c[1] {
			t.Errorf("normalizeText(%q)=%q want %q", c[0], got, c[1])
		}
	}
}

func TestCompareClaim(t *testing.T) {
	// scalar: value matches stdout (trimmed)
	if msg := compareClaim(claim.Info{Check: "x", Value: "4"}, evidence.Check{Output: "4\n"}); msg != "" {
		t.Errorf("expected match, got %q", msg)
	}
	// a fabricated value is caught — exact text, no numeric fuzzing
	if msg := compareClaim(claim.Info{Check: "x", Value: "5"}, evidence.Check{Output: "4\n"}); msg == "" {
		t.Error("expected mismatch for fabricated value")
	}
	// exact-match: differing number formatting now mismatches (by design)
	if msg := compareClaim(claim.Info{Check: "x", Value: "11.4M"}, evidence.Check{Output: "11400000"}); msg == "" {
		t.Error("expected mismatch: text contract does not normalize numbers")
	}
	// a pasted result block is compared verbatim (trailing whitespace ignored)
	want := "division,total\nMidwest,11400000"
	if msg := compareClaim(
		claim.Info{Check: "x", Result: want, ResultFmt: "csv"},
		evidence.Check{Output: want + "\n"}); msg != "" {
		t.Errorf("expected table match, got %q", msg)
	}
	if msg := compareClaim(
		claim.Info{Check: "x", Result: want, ResultFmt: "csv"},
		evidence.Check{Output: "division,total\nMidwest,9000000"}); msg == "" {
		t.Error("expected mismatch for wrong table")
	}
}
