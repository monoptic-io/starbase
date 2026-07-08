package scaffold

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEmitSkillsFreshThenIdempotent(t *testing.T) {
	dir := t.TempDir()
	res, err := EmitSkills(dir, "v1.0.0", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Written) == 0 {
		t.Fatal("first emit should write skill files")
	}
	if _, err := os.Stat(filepath.Join(dir, manifestName)); err != nil {
		t.Fatalf("manifest not written: %v", err)
	}
	// Second emit with no changes: everything unchanged, nothing written.
	res2, err := EmitSkills(dir, "v1.0.0", false)
	if err != nil {
		t.Fatal(err)
	}
	if res2.Changed() || res2.Unchanged != len(res.Written) {
		t.Fatalf("re-emit should be a no-op, got %+v", res2)
	}
}

func TestEmitSkillsPreservesEditsUnlessForced(t *testing.T) {
	dir := t.TempDir()
	if _, err := EmitSkills(dir, "v1.0.0", false); err != nil {
		t.Fatal(err)
	}
	// Pick any emitted file and edit it.
	files, _ := skillFiles()
	edited := filepath.Join(dir, filepath.FromSlash(files[0].rel))
	if err := os.WriteFile(edited, []byte("HAND EDIT"), 0o644); err != nil {
		t.Fatal(err)
	}

	res, err := EmitSkills(dir, "v1.1.0", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != files[0].rel {
		t.Fatalf("edited file should be skipped, got %+v", res)
	}
	if b, _ := os.ReadFile(edited); string(b) != "HAND EDIT" {
		t.Fatal("edit was clobbered without -force")
	}
	// With force, it is overwritten back to canonical content.
	res, err = EmitSkills(dir, "v1.1.0", true)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Updated) != 1 {
		t.Fatalf("force should update the edited file, got %+v", res)
	}
	if b, _ := os.ReadFile(edited); string(b) != string(files[0].data) {
		t.Fatal("force did not restore canonical content")
	}
}

func TestEmitSkillsRemovesPristineOrphans(t *testing.T) {
	dir := t.TempDir()
	if _, err := EmitSkills(dir, "v1.0.0", false); err != nil {
		t.Fatal(err)
	}
	// Simulate a skill that existed in a prior version: a file + a manifest entry
	// whose hash matches it (pristine). A fresh emit should remove it.
	orphan := filepath.Join(dir, "old-skill", "SKILL.md")
	if err := writeFile(orphan, []byte("gone in this version")); err != nil {
		t.Fatal(err)
	}
	m := readManifest(filepath.Join(dir, manifestName))
	m["old-skill/SKILL.md"] = hashBytes([]byte("gone in this version"))
	if err := writeManifest(filepath.Join(dir, manifestName), "v1.0.0", m); err != nil {
		t.Fatal(err)
	}

	res, err := EmitSkills(dir, "v1.1.0", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Removed) != 1 || res.Removed[0] != "old-skill/SKILL.md" {
		t.Fatalf("pristine orphan should be removed, got %+v", res)
	}
	if _, err := os.Stat(orphan); !os.IsNotExist(err) {
		t.Fatal("orphan file still present")
	}
}

func TestInitLayout(t *testing.T) {
	dir := t.TempDir()
	res, err := Init(dir, "Coral Reefs", "v1.0.0", false)
	if err != nil {
		t.Fatal(err)
	}
	// Content lives in a subdir; repo meta lives at the root, separate from it.
	mustExist(t, dir, ContentSubdir+"/index.md")
	mustExist(t, dir, ContentSubdir+"/getting-started.md")
	mustExist(t, dir, "CLAUDE.md")
	mustExist(t, dir, ".github/workflows/pages.yml")
	mustExist(t, dir, SkillsDir+"/research-claims/SKILL.md")
	// No markdown meta files inside the content dir — nothing to deny-list.
	entries, _ := os.ReadDir(filepath.Join(dir, ContentSubdir))
	for _, e := range entries {
		if e.Name() == "CLAUDE.md" || e.Name() == "README.md" {
			t.Fatalf("meta file %s leaked into content dir", e.Name())
		}
	}
	if len(res.Skills.Written) == 0 {
		t.Fatal("skills should be emitted")
	}
}

func mustExist(t *testing.T, dir, rel string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(rel))); err != nil {
		t.Fatalf("expected %s: %v", rel, err)
	}
}

func TestDetectDrift(t *testing.T) {
	repo := t.TempDir()
	skills := filepath.Join(repo, SkillsDir)
	if _, err := EmitSkills(skills, "v1.0.0", false); err != nil {
		t.Fatal(err)
	}

	// Different released version → stale. Search works from a nested content dir.
	nested := filepath.Join(repo, "content", "deep")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	if d := DetectDrift(nested, "v1.2.0"); !d.Stale || d.Emitted != "v1.0.0" {
		t.Fatalf("expected stale drift v1.0.0 vs v1.2.0, got %+v", d)
	}
	// Same version → not stale.
	if d := DetectDrift(repo, "v1.0.0"); d.Stale {
		t.Fatalf("same version should not be stale, got %+v", d)
	}
	// Local build → never stale (comparison is meaningless).
	if d := DetectDrift(repo, "(devel)"); d.Stale {
		t.Fatalf("devel build should never report drift, got %+v", d)
	}
	// No manifest anywhere → not stale.
	if d := DetectDrift(t.TempDir(), "v9.9.9"); d.Stale {
		t.Fatalf("missing manifest should not be stale, got %+v", d)
	}
}
