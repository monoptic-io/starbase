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

	// value-laundering: an honest result block must not excuse a fabricated value
	out := "85 of 292 figures (29.1%) lead with digit 1\n"
	if msg := compareClaim(
		claim.Info{Check: "x", Value: "35.0%", Result: out},
		evidence.Check{Output: out}); msg == "" {
		t.Error("a fabricated value must be caught even when the result block is honest")
	}
	// a value present as a whole word verifies, tolerant of surrounding punctuation
	if msg := compareClaim(
		claim.Info{Check: "x", Value: "29.1%", Result: out},
		evidence.Check{Output: out}); msg != "" {
		t.Errorf("value present as a token should verify, got %q", msg)
	}
	// '=' vs space tolerated: "maxdev 1.7%" matches "maxdev=1.7%"
	if !valueInText("maxdev 1.7%", "N=292 chi2=3.22 maxdev=1.7%") {
		t.Error("value should match across '=' / whitespace differences")
	}
	// whole-word: "4" must not match inside "42"
	if valueInText("4", "42") {
		t.Error("value must match whole words, not substrings")
	}
	// a check-bound claim with neither a result block nor a value is verified by
	// injection — there is no author-typed number to disagree with the computation
	if msg := compareClaim(claim.Info{Check: "x"}, evidence.Check{Output: "anything"}); msg != "" {
		t.Errorf("a check-bound claim relying on injection should verify, got %q", msg)
	}
}
