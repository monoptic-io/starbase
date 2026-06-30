// Package graph turns the resolved link structure into the signals that make
// the site explorable: backlinks, an authority score (PageRank), and a
// "related topics" ranking computed the way a search engine would — from
// direct links, co-citation, and bibliographic coupling.
package graph

import (
	"sort"

	"github.com/monoptic-io/starbase/internal/model"
)

type Graph struct {
	Out     map[string][]string  // slug -> outbound slugs (resolved, deduped, no self)
	In      map[string][]string  // slug -> inbound slugs (backlinks)
	Rank    map[string]float64   // PageRank authority in [0,1]
	related map[string][]Related // cached top-N related per slug
}

type Related struct {
	Slug  string
	Score float64
}

// Build constructs the graph from topics whose links are already resolved.
func Build(topics []*model.Topic) *Graph {
	g := &Graph{
		Out:  map[string][]string{},
		In:   map[string][]string{},
		Rank: map[string]float64{},
	}
	outSet := map[string]map[string]bool{}
	for _, t := range topics {
		g.Out[t.Slug] = nil
		outSet[t.Slug] = map[string]bool{}
	}
	for _, t := range topics {
		for _, l := range t.Links {
			if l.Dead || l.ResolvedSlug == "" || l.ResolvedSlug == t.Slug {
				continue
			}
			if outSet[t.Slug][l.ResolvedSlug] {
				continue
			}
			outSet[t.Slug][l.ResolvedSlug] = true
			g.Out[t.Slug] = append(g.Out[t.Slug], l.ResolvedSlug)
			g.In[l.ResolvedSlug] = append(g.In[l.ResolvedSlug], t.Slug)
		}
	}
	g.pageRank()
	g.computeRelated()
	return g
}

// Backlinks returns the slugs that link to the given topic.
func (g *Graph) Backlinks(slug string) []string { return g.In[slug] }

// Related returns up to n related topics, best first.
func (g *Graph) Related(slug string, n int) []Related {
	r := g.related[slug]
	if len(r) > n {
		r = r[:n]
	}
	return r
}

// pageRank runs the standard iterative algorithm with damping 0.85.
func (g *Graph) pageRank() {
	const (
		damping = 0.85
		iters   = 40
	)
	n := float64(len(g.Out))
	if n == 0 {
		return
	}
	rank := make(map[string]float64, len(g.Out))
	for s := range g.Out {
		rank[s] = 1 / n
	}
	for i := 0; i < iters; i++ {
		next := make(map[string]float64, len(rank))
		var dangling float64
		for s, outs := range g.Out {
			if len(outs) == 0 {
				dangling += rank[s]
			}
		}
		base := (1-damping)/n + damping*dangling/n
		for s := range rank {
			next[s] = base
		}
		for s, outs := range g.Out {
			if len(outs) == 0 {
				continue
			}
			share := damping * rank[s] / float64(len(outs))
			for _, dst := range outs {
				next[dst] += share
			}
		}
		rank = next
	}
	g.Rank = rank
}

func (g *Graph) computeRelated() {
	g.related = make(map[string][]Related, len(g.Out))
	inSet := toSets(g.In)
	outSet := toSets(g.Out)

	for a := range g.Out {
		scores := map[string]float64{}

		// Direct links in either direction are the strongest signal.
		for _, b := range g.Out[a] {
			scores[b] += 3
		}
		for _, b := range g.In[a] {
			scores[b] += 2
		}
		// Bibliographic coupling: a and b cite the same topics.
		for x := range outSet[a] {
			for _, b := range g.In[x] { // others that also link to x
				if b != a {
					scores[b] += 1
				}
			}
		}
		// Co-citation: a and b are cited together by the same topics.
		for x := range inSet[a] {
			for _, b := range g.Out[x] { // others also linked from x
				if b != a {
					scores[b] += 1
				}
			}
		}

		rel := make([]Related, 0, len(scores))
		for b, s := range scores {
			// Blend in a little of b's authority to break ties toward hubs.
			rel = append(rel, Related{Slug: b, Score: s + g.Rank[b]})
		}
		sort.Slice(rel, func(i, j int) bool {
			if rel[i].Score != rel[j].Score {
				return rel[i].Score > rel[j].Score
			}
			return rel[i].Slug < rel[j].Slug
		})
		g.related[a] = rel
	}
}

func toSets(m map[string][]string) map[string]map[string]bool {
	out := make(map[string]map[string]bool, len(m))
	for k, vs := range m {
		s := make(map[string]bool, len(vs))
		for _, v := range vs {
			s[v] = true
		}
		out[k] = s
	}
	return out
}
