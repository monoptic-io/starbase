// Package registry resolves wiki-link names to real topics at compile time.
//
// A link like [[Damped Pendulum]] is matched against every topic's title,
// aliases, slug, and filename (all normalized). Unresolved links are reported
// as warnings so an agent can see exactly which subjects still need writing.
package registry

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/ryannedolan/sitegen/internal/model"
	"github.com/ryannedolan/sitegen/internal/parse"
)

type Registry struct {
	BySlug map[string]*model.Topic
	byName map[string]*model.Topic // normalized name -> topic
}

// New indexes topics for resolution. Duplicate names produce warnings; the
// first topic registered for a name wins.
func New(topics []*model.Topic) (*Registry, []model.Diagnostic) {
	r := &Registry{
		BySlug: make(map[string]*model.Topic, len(topics)),
		byName: make(map[string]*model.Topic, len(topics)*2),
	}
	var diags []model.Diagnostic

	// Stable order so collision reporting is deterministic.
	sorted := append([]*model.Topic(nil), topics...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].SourcePath < sorted[j].SourcePath })

	for _, t := range sorted {
		r.BySlug[t.Slug] = t
	}
	for _, t := range sorted {
		for _, name := range names(t) {
			key := normalize(name)
			if key == "" {
				continue
			}
			if prev, ok := r.byName[key]; ok && prev != t {
				diags = append(diags, model.Diagnostic{
					Severity: model.SevWarn, File: t.SourcePath, Line: 1,
					Message: fmt.Sprintf("name %q also defined by %s; links will resolve to %s",
						name, prev.SourcePath, prev.SourcePath),
				})
				continue
			}
			r.byName[key] = t
		}
	}
	return r, diags
}

// names returns every string a link might use to reach this topic.
func names(t *model.Topic) []string {
	out := []string{t.Title, t.Slug, path.Base(t.Slug)}
	out = append(out, t.Aliases...)
	return out
}

// Resolve finds the topic a link target refers to.
func (r *Registry) Resolve(target string) (*model.Topic, bool) {
	if t, ok := r.byName[normalize(target)]; ok {
		return t, true
	}
	// Allow explicit slug/path references and #fragment suffixes.
	clean := strings.TrimSpace(target)
	if i := strings.IndexByte(clean, '#'); i >= 0 {
		clean = strings.TrimSpace(clean[:i])
	}
	if t, ok := r.byName[normalize(clean)]; ok {
		return t, true
	}
	if t, ok := r.BySlug[slugPath(clean)]; ok {
		return t, true
	}
	return nil, false
}

// ResolveLinks fills in resolution state on every topic's links and returns a
// warning for each dead link.
func (r *Registry) ResolveLinks(topics []*model.Topic) []model.Diagnostic {
	var diags []model.Diagnostic
	for _, t := range topics {
		for i := range t.Links {
			l := &t.Links[i]
			if dst, ok := r.Resolve(l.Target); ok {
				l.ResolvedSlug = dst.Slug
				l.Dead = false
			} else {
				l.Dead = true
				diags = append(diags, model.Diagnostic{
					Severity: model.SevWarn, File: t.SourcePath, Line: l.Line,
					Message: fmt.Sprintf("dead link: no topic named %q", l.Target),
				})
			}
		}
	}
	return diags
}

func normalize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.Join(strings.Fields(s), " ")
	return strings.TrimRight(s, ".,;:!?")
}

func slugPath(p string) string {
	parts := strings.Split(p, "/")
	for i, seg := range parts {
		parts[i] = parse.Slugify(seg)
	}
	return strings.Join(parts, "/")
}
