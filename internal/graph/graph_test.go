package graph

import (
	"testing"

	"github.com/ryannedolan/sitegen/internal/model"
)

func topic(slug string, links ...string) *model.Topic {
	t := &model.Topic{Slug: slug}
	for _, l := range links {
		t.Links = append(t.Links, model.Link{ResolvedSlug: l})
	}
	return t
}

func TestBacklinksAndDirectRelated(t *testing.T) {
	// a -> b, c -> b : b is linked by a and c.
	topics := []*model.Topic{topic("a", "b"), topic("b"), topic("c", "b")}
	g := Build(topics)

	if bl := g.Backlinks("b"); len(bl) != 2 {
		t.Fatalf("b backlinks = %v, want 2", bl)
	}
	// a and c share an outbound target (b): bibliographic coupling makes them related.
	rel := g.Related("a", 5)
	found := false
	for _, r := range rel {
		if r.Slug == "c" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected a related to c via coupling, got %+v", rel)
	}
}

func TestPageRankFavorsHub(t *testing.T) {
	// everyone links to hub.
	topics := []*model.Topic{
		topic("hub"), topic("x", "hub"), topic("y", "hub"), topic("z", "hub"),
	}
	g := Build(topics)
	for _, s := range []string{"x", "y", "z"} {
		if g.Rank["hub"] <= g.Rank[s] {
			t.Errorf("hub rank %.4f should exceed %s rank %.4f", g.Rank["hub"], s, g.Rank[s])
		}
	}
}
