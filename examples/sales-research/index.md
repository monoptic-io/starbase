---
title: Sales Research
summary: A tiny evidence-backed knowledge base — every number is re-computed at build time from a Go check.
---
# Sales Research

A worked example of **verifiable claims**. Each highlighted statement in the
[[Findings]] carries a check — a Go program in `evidence/` that recomputes the
value from the source data. `starbase verify` re-runs it and fails the build if
the article and the data disagree, so a number can't be faked.

Open *How we know this* on any claim to see the query and result; the binding
`check=` re-runs the real computation.
