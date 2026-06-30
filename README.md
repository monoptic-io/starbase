# sitegen

A static knowledge-base generator **for agents**. It turns a tree of markdown
files into a highly interactive but fully static, wiki-like website: articles
link to each other *by name* (resolved to real paths at build time), long pages
get tables of contents, related topics and backlinks are computed from the link
graph the way a search engine would, and pages can embed charts, plots,
simulations, and arbitrary interactive visualizations through a template system.

Every problem it can detect — a dead link, a missing template argument, a broken
template — is reported precisely, so an agent can iterate a subject to completion.

## Install / build

```sh
go build -o sitegen ./cmd/sitegen
```

## Commands

```sh
sitegen check <dir>                 # fast validation: dead links + bad template calls
sitegen build <dir> -o _site \      # full incremental render
        -title "My KB"
sitegen templates [dir]             # list embedded templates and their arguments
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
  mangling and rendered with KaTeX. KaTeX (CSS, JS, and woff2 fonts) is
  **vendored and embedded** in the binary, so the generated site is fully
  self-contained and works offline — no CDN.
- **Shortcodes** invoke templates with validated arguments:

  ```markdown
  {{< sim name="lorenz" >}}
  {{< plot fn="Math.sin(x)" title="sine" >}}
  {{< chart type="line" data="0,1,4,9,16" >}}
  {{< note kind="tip" title="Remember" >}} inner **markdown** {{< /note >}}
  ```

  A missing **required** argument is a hard error. Run `sitegen templates` to
  see every template and its parameters. Custom interactive visualizations are
  authored with `{{< sketch >}} …JavaScript… {{< /sketch >}}`.

### Built-in templates

`chart`, `plot`, `sim` (pendulum · doublependulum · lorenz · nbody · life ·
vectorfield), `sketch`, `note`, `quiz`, `eq`, `figure`, `columns`, `embed`.

### Project overrides

A content directory may contain `templates/` (custom or overriding shortcode
templates), `layout/` (override the page layout), and `theme/` (override
`theme.css` / `app.js`). Built-ins are embedded in the binary; project files
shadow them by name.

## How it works

```
parse      frontmatter + wiki links + shortcodes (code-fence aware)
registry   resolve link names → topics (titles, aliases, slugs)
graph      backlinks, PageRank authority, related = direct links
           + co-citation + bibliographic coupling
render     goldmark → HTML, heading anchors + TOC, shortcode expansion,
           math, page layout with related/backlinks panels
build      per-page fingerprints drive incremental rendering
```

The `internal/` packages are small and single-purpose
(`model`, `parse`, `registry`, `graph`, `tmpl`, `render`, `cache`, `build`).

## Skills

The `skills/` directory contains authoring guides for agents:
`sitegen-authoring`, `interactive-content`, and `flesh-out-subject`.

## Demo

`demo/` is a knowledge base on dynamical systems (foundations, oscillations,
chaos, complex systems) authored to exercise the tool. Build it with:

```sh
sitegen build demo -o demo/_site -title "Dynamical Systems"
```
