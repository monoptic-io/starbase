package parse

import (
	"strings"
	"testing"

	"github.com/monoptic-io/starbase/internal/model"
)

func TestWikiLinks(t *testing.T) {
	body := "See [[Gravity]] and [[Entropy|the entropy page]].\n`not [[a link]]` here.\n"
	links := ScanLinks(body)
	if len(links) != 2 {
		t.Fatalf("want 2 links, got %d: %+v", len(links), links)
	}
	if links[0].Target != "Gravity" || links[0].Display != "" {
		t.Errorf("link0 = %+v", links[0])
	}
	if links[1].Target != "Entropy" || links[1].Display != "the entropy page" {
		t.Errorf("link1 = %+v", links[1])
	}
}

func TestShortcodePairingWithInterleavedSelfClose(t *testing.T) {
	// A self-closing tag precedes a paired one whose close shares no name with
	// the earlier opener — the old single-regex approach mis-paired these.
	body := `{{< sim name="pendulum" >}}
{{< chart data="1,2,3" >}}
{{< note kind="tip" title="Hi" >}}
inner [[Gravity]] text
{{< /note >}}`
	codes := ScanShortcodes(body)
	if len(codes) != 3 {
		t.Fatalf("want 3 shortcodes, got %d: %+v", len(codes), names(codes))
	}
	byName := map[string]string{}
	for _, c := range codes {
		byName[c.Name] = c.Inner
	}
	if _, ok := byName["sim"]; !ok {
		t.Errorf("missing sim; got %v", names(codes))
	}
	note, ok := byName["note"]
	if !ok {
		t.Fatalf("missing note; got %v", names(codes))
	}
	if want := "inner [[Gravity]] text"; !contains(note, want) {
		t.Errorf("note inner = %q, want to contain %q", note, want)
	}
}

func TestShortcodeArgs(t *testing.T) {
	codes := ScanShortcodes(`{{< chart type="line" height=320 label='a b' >}}`)
	if len(codes) != 1 {
		t.Fatalf("want 1, got %d", len(codes))
	}
	a := codes[0].Args
	if a["type"] != "line" || a["height"] != "320" || a["label"] != "a b" {
		t.Errorf("args = %+v", a)
	}
}

func TestMaskCodeFence(t *testing.T) {
	body := "text\n```\n{{< sim name=\"x\" >}}\n[[link]]\n```\nreal [[Gravity]]\n"
	if got := len(ScanShortcodes(body)); got != 0 {
		t.Errorf("shortcodes in fenced code should be ignored, got %d", got)
	}
	if got := len(ScanLinks(body)); got != 1 {
		t.Errorf("only the link outside the fence should count, got %d", got)
	}
}

func TestFrontmatterAndSlug(t *testing.T) {
	topic, _, err := Bytes("mechanics/Damped Pendulum.md",
		[]byte("---\ntitle: Damped Pendulum\naliases: [damping, decay]\ntags: [a, b]\n---\n# Heading\nBody.\n"))
	if err != nil {
		t.Fatal(err)
	}
	if topic.Title != "Damped Pendulum" {
		t.Errorf("title = %q", topic.Title)
	}
	if topic.Slug != "mechanics/damped-pendulum" {
		t.Errorf("slug = %q", topic.Slug)
	}
	if len(topic.Aliases) != 2 || topic.Aliases[0] != "damping" {
		t.Errorf("aliases = %v", topic.Aliases)
	}
}

func names(cs []model.Shortcode) []string {
	var out []string
	for _, c := range cs {
		out = append(out, c.Name)
	}
	return out
}
func contains(s, sub string) bool { return strings.Contains(s, sub) }
