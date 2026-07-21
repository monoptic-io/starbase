---
name: labels
description: Mark starbase topics with workflow labels (open-problem, stub, needs-evidence, speculative, …) and surface them as worklists and site pages. Use when curating a subject's frontier — recording open questions, flagging thin or shaky pages, or asking "what should be worked on next" in a knowledge base.
---

# Labels: surfacing the frontier of a subject

A **label** is a workflow marker on a topic — orthogonal to `tags`. Tags say
what a page is *about* (reader-facing subject terms); labels say what *state*
the page is in or what kind of work it represents. starbase surfaces every
label three ways, so labeling a page is enough to put it on a worklist:

- **`starbase labels <dir>`** prints one tab-separated `label  file  title`
  line per labeled topic — the machine-readable worklist. Filter with
  `-label open-problem`.
- **`labels/<label>.html`** — the build emits a listing page per label (plus a
  `labels/` index), so readers can browse e.g. every open problem in the KB.
- Each labeled page shows a **⚑ label chip** in its header, linking to that
  listing.

Declare labels in frontmatter:

```yaml
---
title: Feigenbaum universality beyond unimodal maps
labels: [open-problem]
tags: [chaos, renormalization]
summary: Does the period-doubling constant survive weaker smoothness assumptions?
---
```

Any label name works — the mechanism is generic. Use a small, consistent
vocabulary per KB; these conventions travel well:

- **`open-problem`** — an unresolved question in the *subject itself* (see below).
- **`stub`** — a page that exists so links resolve but has no real content yet.
- **`needs-evidence`** — prose asserts numbers that should become checked claims
  (see the **research-claims** skill).
- **`speculative`** — content that goes beyond established results; flags where
  a reader (or a later agent pass) should apply skepticism.
- **`low-hanging-fruit`** — a known-tractable extension: an analysis nobody has
  run, a simulation worth adding, an obvious next computation.

## The open-problem convention

Open problems get **their own pages**, not footnotes. For each genuinely
unresolved question in the subject:

1. Create a topic for it (a short page is fine): state the question in the
   first sentence, summarize what *is* known, link the topics it touches with
   `[[...]]`, and note what an answer would change.
2. Label it `labels: [open-problem]` and set a one-line `summary` — the
   listing page reads like a research agenda.
3. Link to it from the established pages it grows out of ("whether this
   extends to X is open — see [[...]]"), so readers hit the frontier from the
   middle of the subject, not only from the index.
4. When a problem is *resolved*, don't delete the page: drop the label,
   rewrite it around the resolution, and keep the backlinks.

An `open-problem` page pairs naturally with `low-hanging-fruit` pages: the
former records questions nobody can answer yet, the latter records work anyone
could do next. Together `starbase labels` becomes the subject's frontier map.

## The worklist loop

Labels complete the diagnostics triad — each is a signal one agent leaves and
another picks up:

- **dead link** → a topic that needs writing (`starbase check`)
- **unsupported claim** → a statement that needs evidence (`starbase check`)
- **label** → deliberate, named state that a page is *in* (`starbase labels`)

Dead links and unsupported claims should be driven to zero; labels are
different — `open-problem` pages may live forever, and that's the point. But
`stub` and `needs-evidence` are debts: when asked to improve a KB, start from

```sh
starbase labels <dir>                     # the full worklist
starbase labels <dir> -label stub        # just the debts of one kind
```

pick an entry, do the work, and remove the label in the same edit.
