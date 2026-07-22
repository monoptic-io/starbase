package build

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/monoptic-io/starbase/internal/model"
	"github.com/monoptic-io/starbase/internal/parse"
)

func write(t *testing.T, root, rel, content string) {
	t.Helper()
	p := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestCopyAssetsMirrorsAndSkips(t *testing.T) {
	content := t.TempDir()
	out := t.TempDir()
	write(t, content, "index.md", "# Home\n")
	write(t, content, "data/experiment.log", "step 0 loss 4.2\n")
	write(t, content, "evidence/x/run", "#!/bin/sh\necho 4\n")
	write(t, content, "templates/kpi.html", "<div></div>")
	write(t, content, ".hidden/secret.txt", "no")
	write(t, content, "data/_scratch.txt", "no")

	n, diags := copyAssets(Config{ContentDir: content, OutDir: out})
	if len(diags) != 0 {
		t.Fatalf("diags: %+v", diags)
	}
	if n != 1 {
		t.Errorf("copied %d files, want 1", n)
	}
	if _, err := os.Stat(filepath.Join(out, "data", "experiment.log")); err != nil {
		t.Errorf("data/experiment.log not published: %v", err)
	}
	for _, absent := range []string{"index.md", "evidence/x/run", "templates/kpi.html", ".hidden/secret.txt", "data/_scratch.txt"} {
		if _, err := os.Stat(filepath.Join(out, filepath.FromSlash(absent))); err == nil {
			t.Errorf("%s should not be published", absent)
		}
	}

	// Second run is a no-op (mtime+size cache).
	if n, _ := copyAssets(Config{ContentDir: content, OutDir: out}); n != 0 {
		t.Errorf("recopied %d files, want 0", n)
	}
}

func TestValidateAssetRefs(t *testing.T) {
	content := t.TempDir()
	write(t, content, "data/run.log", "ok\n")
	write(t, content, "notes/topic.md", "")
	write(t, content, "evidence/x/run", "#!/bin/sh\n")
	body := "see [log](data/run.log), [missing](data/nope.log), " +
		"[esc](../outside.txt), [ev](evidence/x/run), [src](notes/topic.md), " +
		"[web](https://example.com/x.log), [frag](#section)\n"
	topic := &model.Topic{SourcePath: "index.md", Body: body}

	diags := validateAssetRefs(Config{ContentDir: content}, topic)
	if len(diags) != 4 {
		t.Fatalf("want 4 warnings (missing, outside, evidence, .md), got %d: %+v", len(diags), diags)
	}
	for _, d := range diags {
		if d.Severity != model.SevWarn {
			t.Errorf("severity = %v", d.Severity)
		}
	}
}

func TestScanAssetRefsMasksCode(t *testing.T) {
	refs := parse.ScanAssetRefs("`[x](a/b.txt)` and [y](c/d.txt)\n```\n[z](e/f.txt)\n```\n")
	if len(refs) != 1 || refs[0].Target != "c/d.txt" {
		t.Fatalf("got %+v", refs)
	}
}
