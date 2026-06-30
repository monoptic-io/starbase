---
name: flesh-out-subject
description: Build out an entire subject area as a connected starbase knowledge base — planning the topic map, writing interlinked pages, and iterating with `starbase check` until there are no dead links. Use when asked to create or expand a knowledge base on some domain.
---

# Fleshing out a subject area

Goal: a *connected* web of topics where a reader can start anywhere and keep
exploring, every link resolves, and key ideas are illustrated interactively.
Use this together with the **starbase-authoring** and **interactive-content**
skills.

## 1. Map before you write

Sketch the topic graph first. For a subject, identify:

- **Sections** (folders): the major sub-areas. e.g. for dynamical systems:
  `foundations/`, `oscillations/`, `chaos/`, `systems/`.
- **Pillar topics**: 1–2 central pages per section that many others link to.
- **Leaf topics**: specific concepts, examples, phenomena.
- **Cross-links**: which topics naturally reference which. Aim for a dense web,
  not a tree — co-citation drives the "Related topics" panel.

Write a folder `index.md` for each section introducing it and linking to its
members.

## 2. Write in dependency order, link forward freely

Write pillar topics first, but **link to topics that don't exist yet** as you
go — `[[Lyapunov Exponent]]` even before that page exists. Dead links are not
errors to avoid; they are your **worklist**. After a pass:

```
starbase check <dir>
```

Every `dead link: no topic named "X"` is a page still to write (or a name to
fix with an alias). Keep going until the warnings are gone.

## 3. Make every page pull its weight

A good topic page:

- defines the idea in the first sentence;
- has 2–4 `##` sections (enough for an auto TOC on longer ones);
- contains **at least one interactive element** where the concept is dynamic
  (a `sim`, `plot`, `chart`, or custom `sketch`);
- ends with "See also" links to 2–3 sibling topics.

Set `summary`, `tags`, and `aliases` in frontmatter so related-topic cards and
tag pages read well.

## 4. Converge to zero diagnostics

Iterate the loop until clean:

```
starbase check <dir>          # 0 errors, 0 warnings is the bar
starbase build <dir> -o _site -title "<Subject>"
```

Open the built site and click around as a reader would. Ask:

- Can I reach every section from the home page in one or two clicks?
- Does each page offer somewhere interesting to go next?
- Do the "Related topics" and "Referenced by" panels look sensible? If a key
  topic has no backlinks, add links to it from where it's relevant.
- Are the interactive pieces actually illuminating, or decorative? Improve or
  cut the decorative ones.

## 5. Coverage checklist

Before declaring a subject done, confirm:

- [ ] `starbase check` is completely clean.
- [ ] Every section folder has an `index.md`.
- [ ] No orphan pages (every topic has at least one inbound link).
- [ ] Pillar topics are richly linked and have strong interactive demos.
- [ ] Tags are used consistently (they generate browsable index pages).
- [ ] The home page invites exploration and frames the subject.

A subject is finished when a curious reader could fall into it and not hit a
dead end.
