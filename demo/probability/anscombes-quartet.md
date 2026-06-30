---
title: Anscombe's Quartet
aliases: [Anscombe quartet, Anscombe's four datasets]
tags: [probability]
summary: Four small datasets engineered to share almost every summary statistic — mean, variance, correlation, and regression line — yet look completely different when plotted. A warning that numbers without a picture can lie.
weight: 67
---

# Anscombe's Quartet

In 1973 the statistician Francis Anscombe published four tiny datasets — eleven $(x, y)$ points each — built to make a single, stubborn point: **summary statistics can be identical while the data they summarize is nothing alike.** Run the standard battery of descriptive numbers over the four sets and they come out the same to two decimals — same mean of $x$, same mean of $y$, same variances, same correlation, same least-squares line. Plot them and the illusion shatters: one is a clean linear trend, one a smooth curve, one a perfect line dragged off course by a single outlier, and one a vertical stack of points fixed by a lone far-flung observation.

The quartet is the canonical argument for *looking* at data before trusting a model fit to it. A [[Probability Distribution]] is more than its first two moments, and a correlation coefficient is a single number squeezed out of a whole cloud of points — squeezing discards exactly the structure that tells you whether the number means anything.

## The four datasets

The datasets live in `data/anscombe-1.csv` through `data/anscombe-4.csv`, each a column of $x$ and a column of $y$. Sets I, II, and III share the same eleven $x$ values; set IV holds $x$ constant at 8 for ten points and places the eleventh far out at 19. By construction the differences are invisible to the usual statistics:

{{< claim check="anscombe-summary" source="data/anscombe-1.csv" asof="2026-06-30" >}}
All four datasets agree, to two decimals, on **mean $x = 9.00$**, **mean $y = 7.50$**, **variance of $x = 11.00$**, **variance of $y \approx 4.12$**, **correlation $\approx 0.82$**, and the least-squares fit **$y = 3.00 + 0.50\,x$** — yet, as the plots show, they could hardly be more different.

```python
# evidence/anscombe-summary/{inputs: data/anscombe-{1,2,3,4}.csv (4th -> set4.csv), run}
# Per-dataset mean x, mean y, sample variance x & y, Pearson correlation,
# and the least-squares slope & intercept. The columns barely move across rows.
files = ["anscombe-1.csv", "anscombe-2.csv", "anscombe-3.csv", "set4.csv"]
for i, path in enumerate(files, start=1):
    mx, my, vx, vy, r, slope, intercept = stats(*load(path))
    print("%3d  %6.2f  %6.2f  %5.2f  %5.2f  %5.3f  %5.2f  %9.2f"
          % (i, mx, my, vx, vy, r, slope, intercept))
```

```result
set  mean_x  mean_y  var_x  var_y   corr  slope  intercept
  1    9.00    7.50  11.00   4.13  0.816   0.50       3.00
  2    9.00    7.50  11.00   4.13  0.816   0.50       3.00
  3    9.00    7.50  11.00   4.12  0.816   0.50       3.00
  4    9.00    7.50  11.00   4.12  0.817   0.50       3.00
```
{{< /claim >}}

The faint wobble in the last column — the correlation reads `0.816` for the first three sets and `0.817` for the fourth — is real but immaterial: rounded to the two decimals a reader actually reports, every dataset carries the identical correlation.

{{< note kind="key" title="What each dataset really is" >}}
- **Set I** — a genuine, noisy linear relationship; the regression line is honest.
- **Set II** — a smooth concave curve; a *straight* line is simply the wrong model.
- **Set III** — a tight linear run with one outlier that single-handedly bends the fitted slope.
- **Set IV** — $x$ is constant except for one point; that lone point defines the entire line.
{{< /note >}}

## The one number everyone quotes

The headline statistic of the quartet is its shared correlation. It is the same — to the precision anyone reports it — across all four datasets, which is exactly why it is such a treacherous summary on its own:

{{< claim value="0.82" check="anscombe-correlation" source="data/anscombe-1.csv" asof="2026-06-30" >}}
Every dataset in the quartet has a Pearson correlation of **0.82** between $x$ and $y$ (each is $0.816\ldots$, the fourth $0.8165$) — yet only Set I is actually a noisy straight line. Identical correlation, four different stories.

```python
# evidence/anscombe-correlation/{inputs: data/anscombe-{1,2,3,4}.csv, run}
# Assert the four correlations agree to two decimals, then print the shared value.
vals = [corr(*load("anscombe-%d.csv" % i)) for i in (1, 2, 3, 4)]
rounded = {"%.2f" % v for v in vals}
assert len(rounded) == 1, "datasets disagree on correlation: %s" % rounded
print(rounded.pop())
```
{{< /claim >}}

A correlation near $0.82$ sounds like a strong, trustworthy linear relationship — and for Set I it is. For Set II it hides a curve; for Set III it has been inflated by a single outlier; for Set IV it is essentially meaningless, manufactured by one point. The number is the same in every case.

## Why it matters

Anscombe's quartet is a compact rebuttal to "the data is summarized by its statistics." The same caution runs through the rest of this section. The [[Central Limit Theorem]] promises that *averages* converge to a bell curve, and the [[Law of Large Numbers]] promises that sample means approach the truth — but both speak about aggregates, and an aggregate can be blind to shape, skew, and outliers in exactly the way the quartet exploits. A least-squares line is itself a projection — the [[Linear Algebra]] of minimizing squared residuals — and like any projection it throws information away; the quartet is four different pre-images of the same low-dimensional shadow.

The modern descendant of Anscombe's idea is the "Datasaurus" dozen, where a cartoon dinosaur and a starburst and a bullseye all share the same means, variances, and correlation. The lesson is unchanged from 1973: **compute the statistics, but plot the data.**

## See also

- [[Probability Distribution]]
- [[Central Limit Theorem]]
- [[Law of Large Numbers]]
- [[Linear Algebra]]
- [[Benford's Law]]
