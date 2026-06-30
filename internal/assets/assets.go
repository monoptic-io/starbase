// Package assets holds sitegen's built-in templates, page layout, and theme,
// embedded into the binary so the tool works with zero external files.
//
// Resolution order lets a project override anything: a project's own
// templates/ shadows the built-in standard templates, and a project's
// theme/ files shadow the built-in static assets.
package assets

import "embed"

//go:embed templates layout static
var FS embed.FS
