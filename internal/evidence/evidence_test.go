package evidence

import (
	"os"
	"path/filepath"
	"testing"
)

func TestKebab(t *testing.T) {
	cases := map[string]string{"MidwestRegions": "midwest-regions", "RSA": "r-s-a", "Foo": "foo"}
	for in, want := range cases {
		if got := kebab(in); got != want {
			t.Errorf("kebab(%q)=%q want %q", in, got, want)
		}
	}
}

func TestDiscoverFuncs(t *testing.T) {
	dir := t.TempDir()
	gf := filepath.Join(dir, "x.go")
	write(t, gf, `package x
func MidwestRegions() (int, error) { return 4, nil }
func RevenueTable() [][]string { return nil }
func Helper() {}
func priv() int { return 0 }
func WithArg(n int) int { return n }
type T struct{}
func (T) Method() int { return 0 }
`)
	funcs := discoverFuncs([]string{gf})
	got := map[string]int{}
	for _, f := range funcs {
		got[f.Name] = f.Results
	}
	if len(got) != 2 || got["midwest-regions"] != 2 || got["revenue-table"] != 1 {
		t.Fatalf("expected two checks (midwest-regions:2, revenue-table:1), got %+v", got)
	}
}

func TestUnitKeyReflectsSourceAndDeps(t *testing.T) {
	dir := t.TempDir()
	data := filepath.Join(dir, "data.csv")
	write(t, data, "a,b\n1,2\n")
	gf := filepath.Join(dir, "x.go")
	write(t, gf, "//starbase:deps data.csv\npackage x\nfunc F() int { return 1 }\n")
	p := &pkg{Dir: dir, GoFiles: []string{gf}}

	k1 := unitKey(p, dir)
	if unitKey(p, dir) != k1 {
		t.Fatal("key must be stable when nothing changes")
	}
	write(t, data, "a,b\n9,9\n")
	if unitKey(p, dir) == k1 {
		t.Fatal("key must change when a declared data dep changes")
	}
}

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
