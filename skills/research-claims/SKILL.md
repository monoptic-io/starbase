---
name: research-claims
description: Back factual statements in a starbase knowledge base with evidence — the computation you ran to support them — using the claim shortcode. Use when building a research/analysis KB where numbers and assertions must be traceable, not asserted.
---

# Evidence-backed claims

In a research KB, a number is only as good as the computation behind it. The
`{{< claim >}}` shortcode lets you state a fact and attach the **implementation
you ran in your sandbox** to derive it, so a reader can dig down to the query and
its result. starbase runs nothing — *you* run the code; the shortcode surfaces it.

## The discipline

**Never hand-type a number you did not compute.** When you assert a quantity:

1. Run the query/script in your sandbox against a real source.
2. Capture both the result and the exact code.
3. Wrap the statement in a `claim` with that evidence attached.

## Syntax

````markdown
{{< claim value="4" source="data/sales.csv" asof="2026-06-30" >}}
The **Midwest** division is served by **4** distinct sales regions.

```sql
SELECT count(*) AS regions FROM sales WHERE division = 'Midwest';
```

```result
regions
4
```
{{< /claim >}}
````

- The **prose** (the statement, with the value in it) renders inline with an
  *evidence* badge and a *How we know this* disclosure.
- The first non-result fenced block is the **implementation** (any language —
  `sql`, `python`, `bash`, …), shown verbatim in the drawer.
- A fence tagged `result` / `csv` / `tsv` / `table` is the **captured result**,
  rendered as a table (or a single value).
- Args: `value` (the asserted value — checked against the result), `source`
  (the dataset/notebook/file), `asof` (when you ran it).

## Verifiable claims (don't make the reader trust you)

Pasting a result is *attested* evidence — the reader has to trust you transcribed
it faithfully. To make a number **un-fakeable**, bind the claim to a `check` and
write a program that recomputes it. `starbase verify` re-runs the program at build
time and **fails if the article and the computation disagree** — so the build, not
the author, is the source of truth.

Each check is a Go `main` package — a *unit* — under `evidence/`. It can do
anything (pure Go, or shell out to DuckDB, a SQL driver, an API) and prints its
result(s) as JSON: a single result named after its directory,

```json
{ "value": "4" }
```

or a map of named results,

```json
{ "midwest-regions": { "value": "4" },
  "revenue-by-division": { "table": [["division","total"],["Midwest","11400000"]] } }
```

Bind the claim with `check="midwest-regions"`. `verify` compares the claim's
embedded `value`/result to the computed one (numbers are compared numerically, so
`11,400,000` matches `11400000`).

### Incremental, like `go test`

starbase caches each unit's result keyed by a hash of its Go sources **plus the
data files it declares**:

```go
//starbase:deps ../../data/sales.csv
```

A unit is re-run only when its code or a declared dep changes. So:

- editing an unrelated page never re-runs anything;
- editing one check's code re-runs only that check;
- changing a data file re-runs every unit that declares it.

**Put each expensive check in its own package directory** so a minutes-long
simulation isn't re-run when you touch something else. Always declare data deps,
or the local cache can go stale. The cache is a local convenience: CI starts
cold and re-runs everything, so CI is always authoritative.

```sh
starbase verify <dir>          # re-runs only changed units; diffs every checked claim
starbase verify <dir> -force   # ignore the cache and re-run all units
```

Wire `verify` into CI. A claim with a `check` that re-executes and matches is
**verified**; one with only a pasted result is **attested**; one with neither is
**unsupported**. Aim to make load-bearing numbers verified.

## The validation loop (this is the point)

`starbase check` treats a claim with **no implementation and no source** as an
`unsupported claim` warning — exactly like a dead wiki link. That is a
coordination signal: one agent can assert something and mark it a claim, and
another agent (or a later pass) picks up the warning, runs the analysis, and
attaches the query — or discovers the number is wrong and corrects it. If `value`
disagrees with a scalar result, `check` warns about the drift.

So the loop mirrors topic-writing:

- **dead link** → a topic that needs writing.
- **unsupported claim** → a statement that needs evidence.

Drive both to zero. A finished research KB has no dead links *and* no unsupported
claims — every assertion traces to something a reader can inspect and, in
principle, re-run.

## Tips

- Keep claims to load-bearing facts; not every sentence needs one.
- Prefer the smallest query that proves the point; show it in full.
- Watch for `$` in prose — a pair of dollar signs is read as inline math. Write
  `\$11.4M` or `USD 11.4M` for currency.
- For a chart of real data, compute the series in a claim and feed the same
  numbers to `{{< chart >}}` so the figure inherits the provenance.
