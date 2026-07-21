package evidence

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

// isolate the user cache dir so the persisted result/input caches don't leak
// between tests or machines.
func isolateCache(t *testing.T) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CACHE_HOME", filepath.Join(home, "cache"))
}

func TestParseInputs(t *testing.T) {
	specs := parseInputs([]byte("# a comment\ndata/sales.csv\n\nhttp://x/y.csv -> y.csv\nhttp://x/z.csv --> renamed.csv\n"))
	want := []inputSpec{
		{Source: "data/sales.csv"},
		{Source: "http://x/y.csv", LocalName: "y.csv"},
		{Source: "http://x/z.csv", LocalName: "renamed.csv"},
	}
	if len(specs) != len(want) {
		t.Fatalf("got %d specs, want %d: %+v", len(specs), len(want), specs)
	}
	for i := range want {
		if specs[i] != want[i] {
			t.Errorf("spec %d = %+v, want %+v", i, specs[i], want[i])
		}
	}
}

func mkCheck(t *testing.T, kb, name, run, inputs string) {
	t.Helper()
	dir := filepath.Join(kb, "evidence", name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(dir, "run"), run, 0o755)
	if inputs != "" {
		write(t, filepath.Join(dir, "inputs"), inputs, 0o644)
	}
}

func TestRunStagesCachesAndReruns(t *testing.T) {
	isolateCache(t)
	kb := t.TempDir()
	write(t, filepath.Join(kb, "name.txt"), "world\n", 0o644)
	// reads the input by its staged basename, not its repo path
	mkCheck(t, kb, "greeting", "#!/bin/sh\nprintf 'hello %s' \"$(cat name.txt)\"\n", "name.txt\n")

	r, present, err := Run(kb, Options{})
	if err != nil || !present {
		t.Fatalf("Run: err=%v present=%v", err, present)
	}
	if got := r.Checks["greeting"]; got.Output != "hello world" || got.Err != "" {
		t.Fatalf("output=%q err=%q", got.Output, got.Err)
	}
	if r.Units[0].Cached {
		t.Fatal("first run should execute")
	}

	r, _, _ = Run(kb, Options{})
	if !r.Units[0].Cached {
		t.Fatal("unchanged inputs should hit the cache")
	}

	write(t, filepath.Join(kb, "name.txt"), "moon\n", 0o644)
	r, _, _ = Run(kb, Options{})
	if r.Units[0].Cached || r.Checks["greeting"].Output != "hello moon" {
		t.Fatalf("changed input should re-run, got cached=%v out=%q", r.Units[0].Cached, r.Checks["greeting"].Output)
	}
}

func TestHTTPProviderFetchesOnceThenCaches(t *testing.T) {
	isolateCache(t)
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hits++
		fmt.Fprint(w, "x,y\n1,2\n")
	}))
	defer srv.Close()

	kb := t.TempDir()
	mkCheck(t, kb, "remote", "#!/bin/sh\ncat data.csv\n", srv.URL+"/data.csv -> data.csv\n")

	r, _, _ := Run(kb, Options{})
	if got := r.Checks["remote"].Output; got != "x,y\n1,2\n" {
		t.Fatalf("output=%q", got)
	}
	if hits != 1 {
		t.Fatalf("expected 1 fetch, got %d", hits)
	}
	// second run: cached bytes reused, no re-fetch (URL treated as immutable)
	Run(kb, Options{})
	if hits != 1 {
		t.Fatalf("expected no re-fetch on cached run, got %d hits", hits)
	}
}

func TestAttestationsServeColdCacheAndTrustMode(t *testing.T) {
	isolateCache(t)
	kb := t.TempDir()
	write(t, filepath.Join(kb, "name.txt"), "world\n", 0o644)
	mkCheck(t, kb, "greeting", "#!/bin/sh\nprintf 'hello %s' \"$(cat name.txt)\"\n", "name.txt\n")

	// First run executes and records an attestation.
	Run(kb, Options{})
	if _, err := os.Stat(filepath.Join(kb, "evidence", "attestations.json")); err != nil {
		t.Fatalf("expected attestations.json after a successful run: %v", err)
	}

	// A cold cache (fresh machine / CI) is served by the attestation, not re-run.
	isolateCache(t)
	r, _, _ := Run(kb, Options{})
	if u := r.Units[0]; !u.Trusted || u.Cached || u.Err != "" {
		t.Fatalf("cold-cache run should trust the attestation, got %+v", u)
	}
	if r.Checks["greeting"].Output != "hello world" {
		t.Fatalf("output=%q", r.Checks["greeting"].Output)
	}

	// Trust mode with a valid attestation succeeds without executing.
	isolateCache(t)
	r, _, _ = Run(kb, Options{Trust: true})
	if u := r.Units[0]; !u.Trusted || u.Err != "" {
		t.Fatalf("trust mode should accept the attestation, got %+v", u)
	}

	// Changing an input goes stale: trust mode fails instead of executing.
	isolateCache(t)
	write(t, filepath.Join(kb, "name.txt"), "moon\n", 0o644)
	r, _, _ = Run(kb, Options{Trust: true})
	if u := r.Units[0]; u.Err == "" {
		t.Fatalf("trust mode with a stale attestation must fail, got %+v", u)
	}

	// A normal run re-executes and refreshes the attestation; trust then passes.
	r, _, _ = Run(kb, Options{})
	if got := r.Checks["greeting"].Output; got != "hello moon" {
		t.Fatalf("output=%q", got)
	}
	isolateCache(t)
	r, _, _ = Run(kb, Options{Trust: true})
	if u := r.Units[0]; !u.Trusted || u.Err != "" {
		t.Fatalf("trust mode after re-verify should pass, got %+v", u)
	}
}

func TestNonZeroExitIsError(t *testing.T) {
	isolateCache(t)
	kb := t.TempDir()
	mkCheck(t, kb, "boom", "#!/bin/sh\necho nope >&2\nexit 3\n", "")
	r, _, _ := Run(kb, Options{Force: true})
	if r.Checks["boom"].Err == "" {
		t.Fatalf("non-zero exit should surface an error, got %+v", r.Checks["boom"])
	}
}

func TestMissingInputIsError(t *testing.T) {
	isolateCache(t)
	kb := t.TempDir()
	mkCheck(t, kb, "needs", "#!/bin/sh\ncat gone.csv\n", "gone.csv\n")
	r, _, _ := Run(kb, Options{Force: true})
	if r.Checks["needs"].Err == "" {
		t.Fatal("a missing input should surface an error before running")
	}
}
