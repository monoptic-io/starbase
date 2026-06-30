---
name: starbase-authoring
description: Write content for a starbase knowledge base — markdown topics, wiki links, tables of contents, and embedded interactive templates. Use whenever creating or editing .md files that will be built with the `starbase` tool.
---

# Authoring starbase knowledge bases

`starbase` turns a tree of markdown files into a highly interactive but fully
static, wiki-like website. It is built **for agents**: it reports every problem
(dead links, missing template arguments, broken templates) precisely so you can
iterate a subject to completion.

## The golden loop

Author in tight cycles. After *every* batch of edits:

```
starbase check <dir>      # fast: dead links + bad template calls, no rendering
```

Fix everything `check` reports, then occasionally:

```
starbase build <dir> -o _site -title "My KB"   # full incremental render
```

A page is **done** when `starbase check` reports zero errors and zero warnings
for it. Treat warnings as a to-do list, not noise.

## One file = one topic

Each `.md` file is a *topic*, identified by its title (and any aliases). Put
files in folders of any depth; folders become sections in the sidebar.

```markdown
---
title: Damped Pendulum
aliases: [damping, pendulum damping]
tags: [mechanics, oscillations]
summary: A pendulum losing energy to friction, decaying toward rest.
weight: 20            # optional sort order within its section
---

# Damped Pendulum

A damped pendulum loses energy each swing...
```

- `title` — display name and primary link target. If omitted, the first `# H1`
  or the filename is used. **Always set it.**
- `aliases` — alternate names that links can resolve to. Add the plural,
  the abbreviation, common synonyms.
- `summary` — one sentence; shown under the title and in "Related topics"
  cards. If omitted, the first paragraph is used.
- `tags` — generate tag index pages automatically.
- A file named `index.md` (or `_index.md`) is the landing page for its folder.

## Wiki links — the heart of exploration

Link to another topic **by name**, not by path:

```markdown
A [[Pendulum]] is the canonical example. See also [[Entropy|the entropy page]].
```

- `[[Name]]` resolves to the topic whose title/alias/filename matches `Name`
  (case- and spacing-insensitive).
- `[[Name|shown text]]` controls the visible text.
- `[[Name#section]]` links to a heading anchor within that topic.
- A link that resolves to nothing is a **dead link**: it renders in red and
  `check` warns `dead link: no topic named "..."`. This is your signal that a
  topic still needs to be written — either write it, or fix the name.

Link generously. Relatedness and backlinks ("Referenced by") are computed from
this link graph, so well-linked topics surface more connections automatically.
Long articles get an automatic table of contents from their `##`/`###` headings.

## Embedded interactive templates (shortcodes)

Inject charts, simulations, math, and custom visualizations by *invoking a
template*:

```markdown
{{< chart type="line" data="0,1,4,9,16,25" title="Quadratic growth" >}}

{{< sim name="doublependulum" caption="Sensitive dependence on initial conditions." >}}
```

- Run `starbase templates` to list every template and its arguments.
- A **missing required argument is a hard ERROR** (e.g. `{{< chart >}}` with no
  `data`). Supply all required args.
- Block form wraps inner markdown:

  ```markdown
  {{< note kind="tip" title="Remember" >}}
  Energy is conserved only in the **undamped** [[Pendulum]].
  {{< /note >}}
  ```

For building visualizations and simulations, see the **interactive-content**
skill. For building out a whole subject area, see **flesh-out-subject**.

## Math

Write inline math with `$...$` and display math with `$$...$$` directly in
markdown — it is protected from markdown mangling and rendered with KaTeX:

```markdown
The period is $T = 2\pi\sqrt{L/g}$.

$$\ddot\theta = -\frac{g}{L}\sin\theta$$
```

## Style that reads well

- Open with a one- or two-sentence definition a newcomer can grasp.
- Use `##` sections; keep them scannable. Long topics earn a TOC automatically.
- Prefer one strong interactive element over many weak ones.
- End by linking onward: "See also [[X]], [[Y]]." Exploration is the point.
