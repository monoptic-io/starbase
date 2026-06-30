# starbase

A static knowledge-base generator **for agents**. It turns a tree of markdown
files into a highly interactive but fully static, wiki-like website: articles
link to each other *by name* (resolved to real paths at build time), long pages
get tables of contents, related topics and backlinks are computed from the link
graph the way a search engine would, and pages can embed charts, plots,
simulations, and arbitrary interactive visualizations through a template system.

Every problem it can detect — a dead link, a missing template argument, a broken
template — is reported precisely, so an agent can iterate a subject to completion.

The sidebar adapts to the link graph: each page shows only its own **connected
component** — its reachable "world" — so disjoint subjects stay in separate,
focused navigations that merge automatically the moment a single link bridges
them. A build-emitted index powers a **search box** that spans every component,
so you can always jump anywhere.

**▶ See the live demo: https://monoptic-io.github.io/starbase/** — an interactive
field guide to systems, signals & computation, built and deployed from this repo
by GitHub Actions.

## Install / build

```sh
go build -o starbase ./cmd/starbase
```

## Commands

```sh
starbase check <dir>                 # fast validation: dead links + bad template calls
starbase build <dir> -o _site \      # full incremental render
        -title "My KB"
starbase templates [dir]             # list embedded templates and their arguments
```

`check` parses, resolves every wiki link, and validates every template
invocation **without rendering** — ideal for tight authoring loops. `build`
renders the site incrementally: a page is only re-rendered when its own content,
any template it uses, or any topic it links to / is related to has changed.

Exit code is non-zero if there are errors (or, with `-strict`, warnings).

## Authoring model

- **One markdown file = one topic**, in folders of any depth (folders become
  sidebar sections). Frontmatter sets `title`, `aliases`, `tags`, `summary`,
  `weight`, `draft`.
- **Wiki links**: `[[Topic Name]]`, `[[Name|display]]`, `[[Name#section]]`.
  Resolved by title/alias/filename. Unresolved links are warnings (your
  worklist) and render in red.
- **Math**: inline `$...$` and display `$$...$$`, protected from markdown
  mangling and rendered with KaTeX (linked from a CDN by default, or vendored
  locally with `--vendor`; see below).
- **Shortcodes** invoke templates with validated arguments:

  ```markdown
  {{< sim name="lorenz" >}}
  {{< plot fn="Math.sin(x)" title="sine" >}}
  {{< chart type="line" data="0,1,4,9,16" >}}
  {{< note kind="tip" title="Remember" >}} inner **markdown** {{< /note >}}
  ```

  A missing **required** argument is a hard error. Run `starbase templates` to
  see every template and its parameters. Custom interactive visualizations are
  authored with `{{< sketch >}} …JavaScript… {{< /sketch >}}`.

### Built-in templates

`chart`, `plot`, `sim` (pendulum · doublependulum · lorenz · nbody · life ·
vectorfield · wave · interference · wavepacket), `sketch`, `note`, `quiz`, `eq`,
`figure`, `columns`, `embed`.

### Project overrides

A content directory may contain `templates/` (custom or overriding shortcode
templates), `layout/` (override the page layout), and `theme/` (override
`theme.css` / `app.js`). Built-ins are embedded in the binary; project files
shadow them by name.

## Evidence-backed claims

For research/analysis KBs, a `{{< claim >}}` shortcode ties a statement to the
computation that produced it — the query/script the authoring agent ran in its
sandbox, plus the captured result — rendered inline with a *How we know this*
disclosure so a reader can dig down to how a number was derived. starbase
executes nothing; it surfaces what the agent computed.

`check` flags a claim with no implementation or source as an **unsupported
claim** — the same coordination signal as a dead link: one agent asserts, the
warning tells the swarm to go find the evidence (or correct the value). A finished
research KB has no dead links *and* no unsupported claims. See the
`research-claims` skill.

## Third-party assets & offline builds

starbase ships no third-party front-end code in its repository. By default a
build links external assets (currently just KaTeX) from a CDN — ideal for a
public site like the demo.

For an air-gapped or intranet deployment, build with `--vendor`: starbase
downloads the assets on demand (verifying them by checksum), caches them under
your user cache directory, and bundles a local copy into the site so it works
with **no external requests**. Add `--offline` to require the cache and never
touch the network.

```sh
starbase build site -o _site                 # links KaTeX from a CDN
starbase build site -o _site --vendor        # downloads + bundles KaTeX locally
starbase build site -o _site --vendor --offline   # cache only, no network
```

## Continuous integration

`.github/workflows/` contains two workflows:

- **ci.yml** (pull requests): builds, vets, tests, and runs
  `starbase check demo -strict` — validating the demo's links and template
  calls *without* rendering it.
- **pages.yml** (push to `main`): renders the demo and publishes it to GitHub
  Pages.

## How it works

```
parse      frontmatter + wiki links + shortcodes (code-fence aware)
registry   resolve link names → topics (titles, aliases, slugs)
graph      backlinks, PageRank authority, related = direct links
           + co-citation + bibliographic coupling, connected components
render     goldmark → HTML, heading anchors + TOC, shortcode expansion,
           math, page layout, component-scoped collapsible sidebar
build      per-page fingerprints drive incremental rendering;
           emits a search index spanning all components
```

The `internal/` packages are small and single-purpose
(`model`, `parse`, `registry`, `graph`, `tmpl`, `claim`, `render`, `cache`,
`vendor`, `build`).

## Skills

The `skills/` directory contains authoring guides for agents:
`starbase-authoring`, `interactive-content`, `flesh-out-subject`, and `research-claims`.

## Demo

`demo/` is a 15-section interactive field guide — dynamical-systems foundations,
oscillations, waves, Fourier analysis, chaos, complex systems, linear algebra,
graph theory & networks, probability, information theory, optimization & learning,
cryptography, number theory, and computability & complexity — plus a deliberately
disjoint **music theory** section that demonstrates the reachability-scoped
sidebar. Authored by sub-agents to exercise the tool.
It is published live at **https://monoptic-io.github.io/starbase/**. Build it
locally with:

```sh
starbase build demo -o demo/_site -title "Systems, Signals & Computation"
```
