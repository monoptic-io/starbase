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

A check is a directory under `evidence/` containing an executable `run` and an
`inputs` manifest. The contract is just **stdout and exit code**, like a golden
test — so `run` can be *any* program. Bind it to a claim with
`check="<directory name>"`:

```
evidence/midwest-regions/inputs        # one source per line; # comments allowed
  data/sales.csv                       #   a local file (relative to the KB root)…
  https://example.com/feed.csv         #   …or an http(s) URL
  https://example.com/d.csv -> d.csv   #   …optionally renamed with -> name

evidence/midwest-regions/run           # chmod +x — needs a #! line
  awk -F, 'NR>1 && $2=="Midwest"{n++} END{print n+0}' sales.csv
```

Each input is resolved by a **provider** (file or http) and staged into a fresh
working directory under its basename; `verify` runs `./run` **there**, so it reads
inputs by name (`sales.csv`) and sees nothing else. It compares `run`'s **stdout,
trimmed, against the result the claim embeds** — its result block if it has one,
else its `value`. The comparison is **exact text**, not numeric: `11.4M` does
*not* match `11400000`, so make both sides agree (let the program print what the
prose says, or vice-versa). A non-zero exit — or an input that won't resolve (a
missing file, a 404) — is a check failure, reported at `evidence/<check>/run`.
Because the contract is text, the program can be a shell one-liner over DuckDB, a
Python script, a compiled binary — whatever fits.

Show the reader the *real* command — paste `run` (or its core) as the claim's
implementation, so the drawer matches what the build executes.

### Incremental, like `go test`

starbase caches each check keyed by a hash of its `run` script, its `inputs`
manifest, and the **resolved content of every input**. A check re-runs only when
one of those changes:

- editing an unrelated page never re-runs anything;
- editing one check's `run` re-runs only that check;
- changing a declared input re-runs every check that uses it.

**One expensive check per directory** so a minutes-long run isn't re-run when you
touch something else. A fetched URL is treated as immutable between builds — a
local `verify` reuses the cached bytes rather than re-fetching. That's safe
because the cache is only a local convenience: CI starts cold, re-fetches, and the
**output comparison** catches any source that has drifted. So you don't pin URLs;
you let the build fail when the recomputed output stops matching the article.

```sh
starbase verify <dir>          # re-runs only changed checks; diffs every checked claim
starbase verify <dir> -force   # ignore the cache, re-fetch URLs, re-run everything
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
