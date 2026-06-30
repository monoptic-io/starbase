package evidence

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseOutputSingle(t *testing.T) {
	got, err := parseOutput([]byte(`{"value":"4"}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[""].Value != "4" {
		t.Fatalf("single result should key on \"\": %+v", got)
	}
}

func TestParseOutputMap(t *testing.T) {
	got, err := parseOutput([]byte(`{"a":{"value":"1"},"b":{"table":[["x"],["2"]]}}`))
	if err != nil {
		t.Fatal(err)
	}
	if got["a"].Value != "1" || len(got["b"].Table) != 2 {
		t.Fatalf("map results not parsed: %+v", got)
	}
}

func TestUnitKeyReflectsSourceAndDeps(t *testing.T) {
	dir := t.TempDir()
	data := filepath.Join(dir, "data.csv")
	write(t, data, "a,b\n1,2\n")
	write(t, filepath.Join(dir, "main.go"), "//starbase:deps data.csv\npackage main\nfunc main(){}\n")

	k1, err := unitKey(dir)
	if err != nil {
		t.Fatal(err)
	}
	if k2, _ := unitKey(dir); k2 != k1 {
		t.Fatal("key must be stable when nothing changes")
	}
	// changing the declared data dependency must change the key
	write(t, data, "a,b\n9,9\n")
	if k3, _ := unitKey(dir); k3 == k1 {
		t.Fatal("key must change when a declared data dep changes")
	}
}

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
