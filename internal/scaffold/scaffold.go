// Package scaffold bootstraps and maintains a starbase knowledge-base repo. It
// emits the agent-facing skills (embedded in the binary) into a repo's
// .claude/skills/ directory, scaffolds a ready-to-build repo with `Init`, and
// detects when a repo's emitted skills have drifted from the running binary.
//
// The emitted skills are treated as *managed files*: each is stamped in a
// manifest with the hash of what the tool last wrote, so a re-emit can update
// pristine files silently while preserving (and flagging) any the user edited.
package scaffold

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strings"

	"github.com/monoptic-io/starbase/internal/assets"
)

// SkillsDir is the repo-relative location Claude Code discovers skills from.
const SkillsDir = ".claude/skills"

// manifestName is the per-repo stamp recording the version and the canonical
// hash of every skill file the tool last wrote there.
const manifestName = ".starbase-version"

// Version reports the running binary's version. For a `go install …@vX.Y.Z`
// build it is that module version; for a local build it is "(devel)" with the
// short VCS revision appended when available.
func Version() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return "(devel)"
	}
	if v := bi.Main.Version; v != "" && v != "(devel)" {
		return v
	}
	rev := ""
	for _, s := range bi.Settings {
		if s.Key == "vcs.revision" {
			rev = s.Value
		}
	}
	if len(rev) >= 7 {
		return "(devel)+" + rev[:7]
	}
	return "(devel)"
}

// isReleased reports whether v is a real released version (not a local build),
// so drift is only asserted between two comparable versions.
func isReleased(v string) bool { return v != "" && !strings.HasPrefix(v, "(devel)") }

// skillFile is one embedded skill document.
type skillFile struct {
	rel  string // path relative to the skills dir, e.g. "research-claims/SKILL.md"
	data []byte
	hash string
}

// skillFiles reads every embedded skill document.
func skillFiles() ([]skillFile, error) {
	var out []skillFile
	err := fs.WalkDir(assets.Skills, "skills", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, err := assets.Skills.ReadFile(p)
		if err != nil {
			return err
		}
		out = append(out, skillFile{rel: strings.TrimPrefix(p, "skills/"), data: b, hash: hashBytes(b)})
		return nil
	})
	sort.Slice(out, func(i, j int) bool { return out[i].rel < out[j].rel })
	return out, err
}

func hashBytes(b []byte) string {
	s := sha256.Sum256(b)
	return hex.EncodeToString(s[:])
}

func hashFile(p string) (string, bool) {
	b, err := os.ReadFile(p)
	if err != nil {
		return "", false
	}
	return hashBytes(b), true
}

// EmitResult summarizes what an EmitSkills call did.
type EmitResult struct {
	Written   []string // newly created
	Updated   []string // pristine files refreshed to the new version
	Skipped   []string // locally modified — left as-is (re-run with force to overwrite)
	Removed   []string // pristine files for skills no longer in this version
	Unchanged int
	Dir       string // absolute skills dir written to
}

// Changed reports whether the emit did anything the user should notice.
func (r EmitResult) Changed() bool {
	return len(r.Written)+len(r.Updated)+len(r.Removed)+len(r.Skipped) > 0
}

// EmitSkills writes the embedded skills into destDir (typically <repo>/.claude/skills).
// It updates files that are pristine (untouched since the tool last wrote them),
// skips files the user has edited unless force is set, removes pristine files for
// skills that no longer exist, and rewrites the manifest stamped with version.
func EmitSkills(destDir, version string, force bool) (EmitResult, error) {
	res := EmitResult{Dir: destDir}
	files, err := skillFiles()
	if err != nil {
		return res, err
	}
	recorded := readManifest(filepath.Join(destDir, manifestName)) // path -> hash we last wrote

	newManifest := map[string]string{}
	present := map[string]bool{}
	for _, f := range files {
		present[f.rel] = true
		target := filepath.Join(destDir, filepath.FromSlash(f.rel))
		dh, exists := hashFile(target)
		rh := recorded[f.rel]

		switch {
		case !exists:
			if err := writeFile(target, f.data); err != nil {
				return res, err
			}
			res.Written = append(res.Written, f.rel)
			newManifest[f.rel] = f.hash
		case dh == f.hash:
			res.Unchanged++
			newManifest[f.rel] = f.hash
		case rh != "" && dh == rh:
			// Pristine: matches what we last wrote, so it's safe to update.
			if err := writeFile(target, f.data); err != nil {
				return res, err
			}
			res.Updated = append(res.Updated, f.rel)
			newManifest[f.rel] = f.hash
		case force:
			if err := writeFile(target, f.data); err != nil {
				return res, err
			}
			res.Updated = append(res.Updated, f.rel)
			newManifest[f.rel] = f.hash
		default:
			// Locally modified (or authored by hand): leave it, keep flagging it.
			res.Skipped = append(res.Skipped, f.rel)
			newManifest[f.rel] = rh // preserve prior record (may be empty)
		}
	}

	// Remove files for skills that vanished from this version, but only if the
	// on-disk copy is pristine (matches what we last wrote) — never clobber edits.
	for rel, rh := range recorded {
		if present[rel] {
			continue
		}
		target := filepath.Join(destDir, filepath.FromSlash(rel))
		if dh, exists := hashFile(target); exists && dh == rh {
			if err := os.Remove(target); err == nil {
				res.Removed = append(res.Removed, rel)
			}
		}
	}

	if err := writeManifest(filepath.Join(destDir, manifestName), version, newManifest); err != nil {
		return res, err
	}
	sort.Strings(res.Written)
	sort.Strings(res.Updated)
	sort.Strings(res.Skipped)
	sort.Strings(res.Removed)
	return res, nil
}

func writeFile(target string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	return os.WriteFile(target, data, 0o644)
}

// readManifest parses a skills manifest into a path->hash map. A missing or
// unreadable manifest yields an empty map (every file then looks user-authored).
func readManifest(path string) map[string]string {
	out := map[string]string{}
	b, err := os.ReadFile(path)
	if err != nil {
		return out
	}
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "version ") {
			continue
		}
		if f := strings.Fields(line); len(f) == 2 {
			out[f[1]] = f[0]
		}
	}
	return out
}

// manifestVersion returns the version a manifest was stamped with, or "".
func manifestVersion(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(b), "\n") {
		if v, ok := strings.CutPrefix(strings.TrimSpace(line), "version "); ok {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func writeManifest(path, version string, hashes map[string]string) error {
	var b strings.Builder
	b.WriteString("# starbase skills manifest — managed by `starbase skills`.\n")
	b.WriteString("# Regenerate after upgrading the binary; do not edit by hand.\n")
	fmt.Fprintf(&b, "version %s\n", version)
	rels := make([]string, 0, len(hashes))
	for rel := range hashes {
		rels = append(rels, rel)
	}
	sort.Strings(rels)
	for _, rel := range rels {
		if hashes[rel] != "" {
			fmt.Fprintf(&b, "%s  %s\n", hashes[rel], rel)
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

// Drift describes a mismatch between a repo's emitted skills and the binary.
type Drift struct {
	Stale   bool
	Emitted string // version stamped in the repo's manifest
	Current string // running binary version
	Dir     string // skills dir that carries the manifest
}

// DetectDrift searches upward from startDir for a .claude/skills/.starbase-version
// manifest and reports whether it was emitted by a different (released) version
// than current. It never flags drift for local ("(devel)") builds, where the
// comparison would be meaningless.
func DetectDrift(startDir, current string) Drift {
	if !isReleased(current) {
		return Drift{}
	}
	dir := startDir
	for i := 0; i < 24; i++ {
		manifest := filepath.Join(dir, SkillsDir, manifestName)
		if emitted := manifestVersion(manifest); emitted != "" {
			return Drift{
				Stale:   isReleased(emitted) && emitted != current,
				Emitted: emitted,
				Current: current,
				Dir:     filepath.Join(dir, SkillsDir),
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return Drift{}
}
