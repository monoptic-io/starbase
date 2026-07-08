package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// InitResult summarizes a scaffold.
type InitResult struct {
	Created []string // repo-relative paths written
	Skipped []string // paths left untouched because they already existed
	Skills  EmitResult
	Dir     string
}

// Init scaffolds a content-at-root starbase KB in dir: a couple of starter
// topics, a Pages workflow, a CLAUDE.md operating manual, a .gitignore, and the
// skills under .claude/skills/. Existing files are never overwritten (so it is
// safe to run in a partially set-up repo); skills are emitted edit-safely. The
// result builds green under `starbase check` immediately.
func Init(dir, title, version string, force bool) (InitResult, error) {
	res := InitResult{Dir: dir}
	if strings.TrimSpace(title) == "" {
		title = "Knowledge Base"
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return res, err
	}

	files := map[string]string{
		"index.md":                    fmt.Sprintf(indexTmpl, title),
		"getting-started.md":          gettingStartedTmpl,
		".gitignore":                  gitignoreTmpl,
		"CLAUDE.md":                   fmt.Sprintf(claudeTmpl, title),
		".github/workflows/pages.yml": fmt.Sprintf(pagesTmpl, title),
	}
	rels := make([]string, 0, len(files))
	for rel := range files {
		rels = append(rels, rel)
	}
	sort.Strings(rels)
	for _, rel := range rels {
		target := filepath.Join(dir, filepath.FromSlash(rel))
		if _, err := os.Stat(target); err == nil {
			res.Skipped = append(res.Skipped, rel)
			continue
		}
		if err := writeFile(target, []byte(files[rel])); err != nil {
			return res, err
		}
		res.Created = append(res.Created, rel)
	}

	skills, err := EmitSkills(filepath.Join(dir, SkillsDir), version, force)
	if err != nil {
		return res, err
	}
	res.Skills = skills
	sort.Strings(res.Created)
	sort.Strings(res.Skipped)
	return res, nil
}

const indexTmpl = `---
title: Home
summary: %[1]s — a starbase knowledge base.
---

# %[1]s

Welcome. This is a [[Getting Started|starbase]] knowledge base — a tree of
markdown topics that builds into an interactive, fully static site.

Start with [[Getting Started]].
`

const gettingStartedTmpl = `---
title: Getting Started
aliases: [getting started, start here]
tags: [meta]
summary: How this knowledge base is organized and how to add to it.
---

# Getting Started

Every ` + "`.md`" + ` file is a **topic**, identified by its title. Link between
topics by name with ` + "`[[Wiki Links]]`" + ` — a link to a page that does not exist
yet is a *dead link*, your worklist for what to write next.

Author in a tight loop:

` + "```sh" + `
starbase check .      # dead links + bad template calls — drive this to zero
starbase build . -o _site -title "Your Title"
` + "```" + `

See the guides in ` + "`.claude/skills/`" + ` for the full authoring workflow:
writing topics, embedding interactive charts and simulations, and backing
numbers with re-runnable evidence.

## Add your first topic

Create a new ` + "`.md`" + ` file, give it a ` + "`title`" + `, and link to it from
[[Home]]. Run ` + "`starbase check .`" + ` and fix whatever it reports.
`

const gitignoreTmpl = `# starbase build output
_site/
`

const claudeTmpl = `# %[1]s — a starbase knowledge base

This repo is a **starbase** knowledge base: a tree of markdown topics that
builds into an interactive, fully static website.

## The loop

Author in tight cycles. After every batch of edits:

    starbase check .        # dead links + bad template calls (fast, no render)

Drive every error and warning to zero. A dead link means a topic still needs
writing; an unsupported claim means a number still needs evidence. Then:

    starbase verify .       # re-run evidence checks; every claim must still match
    starbase build . -o _site -title "%[1]s"

## Where the manual lives

Full authoring guides are in ` + "`.claude/skills/`" + ` and are picked up
automatically. They cover writing interlinked topics, embedding charts /
simulations / custom JavaScript widgets, and backing load-bearing numbers with
re-runnable evidence (` + "`{{< val >}}`" + ` / ` + "`{{< data >}}`" + ` / ` + "`{{< claim >}}`" + `).

Run ` + "`starbase templates .`" + ` to list every shortcode and its arguments.

## Upgrades

The skills are version-locked to the ` + "`starbase`" + ` binary. After upgrading it,
` + "`starbase check`" + ` will note if the emitted skills are stale; run
` + "`starbase skills`" + ` to refresh them (your local edits are preserved).
`

const pagesTmpl = `name: Deploy to GitHub Pages

# Post-merge: build the knowledge base and publish it to GitHub Pages.
on:
  push:
    branches: [main]
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: pages
  cancel-in-progress: true

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: stable

      # Install the starbase CLI. Pin @vX.Y.Z for reproducible builds.
      - name: Install starbase
        run: go install github.com/monoptic-io/starbase/cmd/starbase@latest

      - name: Build site
        run: starbase build . -o _site -title %[1]q -strict

      - run: touch _site/.nojekyll
      - uses: actions/configure-pages@v5
      - uses: actions/upload-pages-artifact@v3
        with:
          path: _site

  deploy:
    needs: build
    runs-on: ubuntu-latest
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    steps:
      - id: deployment
        uses: actions/deploy-pages@v4
`
