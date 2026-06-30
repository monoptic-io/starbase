// Package model defines the core data types shared across starbase.
//
// A Topic is the unit of content: one markdown file. Topics are identified by
// name (their title plus any aliases) and resolved to real paths at compile
// time, which is what lets articles link to each other like a wiki.
package model

// Topic is a single markdown document and everything we parsed out of it.
type Topic struct {
	// Identity / location.
	SourcePath string // path relative to content root, e.g. "mechanics/pendulum.md"
	Slug       string // url-ish path without extension, e.g. "mechanics/pendulum"
	OutPath    string // output file, e.g. "mechanics/pendulum.html"
	Title      string // display title
	Aliases    []string
	Tags       []string
	Summary    string // short description (frontmatter or first paragraph)
	Draft      bool
	Weight     int // sort hint within a section

	// Parsed content.
	ContentHash string      // hash of the raw file bytes
	Body        string      // markdown body with frontmatter stripped
	Headings    []Heading   // for table of contents + anchors
	Links       []Link      // outbound wiki links
	Shortcodes  []Shortcode // template invocations
	WordCount   int
}

// Heading is a section header used to build a table of contents.
type Heading struct {
	Level int    // 1..6
	Text  string // rendered text
	Slug  string // anchor id, unique within the document
	Line  int
}

// Link is an outbound wiki link: [[Target]] or [[Target|Display]].
type Link struct {
	Target  string // topic name as written by the author
	Display string // optional display text; falls back to Target
	Line    int
	Start   int // byte offset of the link in Body
	End     int

	// Resolution, filled in after the registry is built.
	ResolvedSlug string // slug of the matched topic, "" if dead
	Dead         bool
}

// Shortcode is a template invocation embedded in markdown:
//
//	{{< chart type="line" height="320" >}}
//	{{< sim name="pendulum" >}} ... inner ... {{< /sim >}}
type Shortcode struct {
	Name  string
	Args  map[string]string
	Inner string // block body, empty for self-closing form
	Line  int
	Start int // byte offset of the invocation in Body
	End   int
	Raw   string // exact source text, used for placeholder substitution
}

// Diagnostic is a warning or error surfaced during build/check.
type Diagnostic struct {
	Severity Severity
	File     string
	Line     int
	Message  string
}

type Severity int

const (
	SevWarn Severity = iota
	SevError
)

func (s Severity) String() string {
	if s == SevError {
		return "ERROR"
	}
	return "WARN"
}
