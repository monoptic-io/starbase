package evidence

import (
	"os"
	"path/filepath"
	"testing"
)

func write(t *testing.T, path, content string, mode os.FileMode) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), mode); err != nil {
		t.Fatal(err)
	}
}

func TestUnitKeyReflectsScriptAndInputs(t *testing.T) {
	dir := t.TempDir()
	data := filepath.Join(dir, "data.csv")
	write(t, data, "a,b\n1,2\n", 0o644)
	run := filepath.Join(dir, "run")
	write(t, run, "#!/bin/sh\n# starbase:inputs data.csv\ncat data.csv\n", 0o755)

	k1 := unitKey(run, dir)
	if k1 == "" || unitKey(run, dir) != k1 {
		t.Fatal("key must be non-empty and stable when nothing changes")
	}
	write(t, data, "a,b\n9,9\n", 0o644)
	if unitKey(run, dir) == k1 {
		t.Fatal("key must change when a declared input changes")
	}
	k2 := unitKey(run, dir)
	write(t, run, "#!/bin/sh\n# starbase:inputs data.csv\nhead -1 data.csv\n", 0o755)
	if unitKey(run, dir) == k2 {
		t.Fatal("key must change when the run script changes")
	}
}

func TestRunExecutesCachesAndReruns(t *testing.T) {
	kb := t.TempDir()
	ev := filepath.Join(kb, "evidence", "greeting")
	if err := os.MkdirAll(ev, 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(kb, "name.txt"), "world\n", 0o644)
	write(t, filepath.Join(ev, "run"),
		"#!/bin/sh\n# starbase:inputs name.txt\nprintf 'hello %s' \"$(cat name.txt)\"\n", 0o755)

	r, present, err := Run(kb, false)
	if err != nil || !present {
		t.Fatalf("Run: err=%v present=%v", err, present)
	}
	if got := r.Checks["greeting"]; got.Output != "hello world" || got.Err != "" {
		t.Fatalf("output=%q err=%q", got.Output, got.Err)
	}
	if len(r.Units) != 1 || r.Units[0].Cached {
		t.Fatalf("first run should execute, got %+v", r.Units)
	}

	// unchanged -> cached
	r, _, _ = Run(kb, false)
	if !r.Units[0].Cached {
		t.Fatal("second run should be cached")
	}

	// changed input -> re-run with new output
	write(t, filepath.Join(kb, "name.txt"), "moon\n", 0o644)
	r, _, _ = Run(kb, false)
	if r.Units[0].Cached || r.Checks["greeting"].Output != "hello moon" {
		t.Fatalf("changed input should re-run, got cached=%v out=%q", r.Units[0].Cached, r.Checks["greeting"].Output)
	}
}

func TestNonZeroExitIsError(t *testing.T) {
	kb := t.TempDir()
	ev := filepath.Join(kb, "evidence", "boom")
	if err := os.MkdirAll(ev, 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(ev, "run"), "#!/bin/sh\necho nope >&2\nexit 3\n", 0o755)

	r, _, _ := Run(kb, true)
	if got := r.Checks["boom"]; got.Err == "" {
		t.Fatalf("non-zero exit should surface an error, got %+v", got)
	}
}
