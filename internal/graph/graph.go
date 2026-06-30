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
	Out        map[string][]string  // slug -> outbound slugs (resolved, deduped, no self)
	In         map[string][]string  // slug -> inbound slugs (backlinks)
	Rank       map[string]float64   // PageRank authority in [0,1]
	related    map[string][]Related // cached top-N related per slug
	Components map[string]int       // slug -> connected-component id (0 = largest)
	compSize   map[int]int
}

// HubSlug is excluded from connectivity: the home page links to every section as
// a directory, which would otherwise fuse all topics into one component.
const HubSlug = "index"

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
	g.computeComponents()
	return g
}

// computeComponents finds weakly-connected components (links treated as
// undirected), ignoring edges incident to the home hub. Component ids are
// assigned by size, largest first, so id 0 is the main cluster.
func (g *Graph) computeComponents() {
	parent := make(map[string]string, len(g.Out))
	var find func(string) string
	find = func(x string) string {
		for parent[x] != x {
			parent[x] = parent[parent[x]]
			x = parent[x]
		}
		return x
	}
	union := func(a, b string) {
		ra, rb := find(a), find(b)
		if ra != rb {
			parent[ra] = rb
		}
	}
	slugs := make([]string, 0, len(g.Out))
	for s := range g.Out {
		parent[s] = s
		slugs = append(slugs, s)
	}
	sort.Strings(slugs)
	for _, s := range slugs {
		if s == HubSlug {
			continue
		}
		for _, d := range g.Out[s] {
			if d == HubSlug {
				continue
			}
			union(s, d)
		}
	}
	// Tally sizes per root, then rank roots by size (ties broken by slug).
	size := map[string]int{}
	for _, s := range slugs {
		size[find(s)]++
	}
	roots := make([]string, 0, len(size))
	for r := range size {
		roots = append(roots, r)
	}
	sort.Slice(roots, func(i, j int) bool {
		if size[roots[i]] != size[roots[j]] {
			return size[roots[i]] > size[roots[j]]
		}
		return roots[i] < roots[j]
	})
	id := map[string]int{}
	for i, r := range roots {
		id[r] = i
	}
	g.Components = make(map[string]int, len(slugs))
	g.compSize = map[int]int{}
	for _, s := range slugs {
		c := id[find(s)]
		g.Components[s] = c
		g.compSize[c]++
	}
}

// Component returns the connected-component id of a slug (-1 if unknown).
func (g *Graph) Component(slug string) int {
	if c, ok := g.Components[slug]; ok {
		return c
	}
	return -1
}

// ComponentCount returns the number of connected components.
func (g *Graph) ComponentCount() int { return len(g.compSize) }

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
