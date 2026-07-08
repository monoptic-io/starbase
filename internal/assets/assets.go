// Package assets holds starbase's built-in templates, page layout, and theme,
// embedded into the binary so the tool works with zero external files.
//
// Resolution order lets a project override anything: a project's own
// templates/ shadows the built-in standard templates, and a project's
// theme/ files shadow the built-in static assets.
package assets

import "embed"

//go:embed templates layout static
var FS embed.FS

// Skills holds the agent-facing authoring guides, embedded so `starbase init`
// and `starbase skills` can plant them into a KB repo (in .claude/skills/,
// where Claude Code discovers them) without any external files. They are
// version-locked to the binary: emitting always writes the manual that matches
// this build.
//
//go:embed skills
var Skills embed.FS
