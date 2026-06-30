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

## Look at the data

Here is the whole argument in one picture. Each panel plots one dataset against the **same** least-squares line $y = 3.00 + 0.50\,x$ that every one of them produces. The summary statistics are identical; the shapes could hardly be more different.

{{< sketch height="380" caption="Anscombe's quartet: identical means, variances, correlation, and regression line — four different stories. Set II bends, Set III has a lone high outlier, Set IV is a single point doing all the work." >}}
const sets = [
  [[10,8.04],[8,6.95],[13,7.58],[9,8.81],[11,8.33],[14,9.96],[6,7.24],[4,4.26],[12,10.84],[7,4.82],[5,5.68]],
  [[10,9.14],[8,8.14],[13,8.74],[9,8.77],[11,9.26],[14,8.1],[6,6.13],[4,3.1],[12,9.13],[7,7.26],[5,4.74]],
  [[10,7.46],[8,6.77],[13,12.74],[9,7.11],[11,7.81],[14,8.84],[6,6.08],[4,5.39],[12,8.15],[7,6.42],[5,5.73]],
  [[8,6.58],[8,5.76],[8,7.71],[8,8.84],[8,8.47],[8,7.04],[8,5.25],[19,12.5],[8,5.56],[8,7.91],[8,6.89]]
];
const labels = ['I', 'II', 'III', 'IV'];
const gap = 12, pw = (W - gap) / 2, ph = (H - gap) / 2;
const xlo = 2, xhi = 20, ylo = 2, yhi = 14, m = 26;
for (let i = 0; i < 4; i++) {
  const col = i % 2, row = (i / 2) | 0;
  const x0 = col * (pw + gap), y0 = row * (ph + gap);
  const ax = x0 + m, ab = y0 + ph - 20, plotW = pw - m - 8, plotH = ph - m - 8;
  const sx = v => ax + (v - xlo) / (xhi - xlo) * plotW;
  const sy = v => ab - (v - ylo) / (yhi - ylo) * plotH;
  ctx.strokeStyle = 'rgba(150,160,180,0.35)'; ctx.lineWidth = 1;
  ctx.strokeRect(ax, ab - plotH, plotW, plotH);
  ctx.strokeStyle = '#e0654f'; ctx.lineWidth = 2;
  ctx.beginPath(); ctx.moveTo(sx(xlo), sy(3 + 0.5 * xlo)); ctx.lineTo(sx(xhi), sy(3 + 0.5 * xhi)); ctx.stroke();
  ctx.fillStyle = '#5aa9e6';
  for (const p of sets[i]) { ctx.beginPath(); ctx.arc(sx(p[0]), sy(p[1]), 3.5, 0, 7); ctx.fill(); }
  ctx.fillStyle = '#cfd6e6'; ctx.font = '13px system-ui, sans-serif';
  ctx.fillText('Set ' + labels[i], ax + 5, y0 + 16);
}
{{< /sketch >}}

## Why it matters

Anscombe's quartet is a compact rebuttal to "the data is summarized by its statistics." The same caution runs through the rest of this section. The [[Central Limit Theorem]] promises that *averages* converge to a bell curve, and the [[Law of Large Numbers]] promises that sample means approach the truth — but both speak about aggregates, and an aggregate can be blind to shape, skew, and outliers in exactly the way the quartet exploits. A least-squares line is itself a projection — the [[Linear Algebra]] of minimizing squared residuals — and like any projection it throws information away; the quartet is four different pre-images of the same low-dimensional shadow.

The modern descendant of Anscombe's idea is the "Datasaurus" dozen, where a cartoon dinosaur and a starburst and a bullseye all share the same means, variances, and correlation. The lesson is unchanged from 1973: **compute the statistics, but plot the data.**

## See also

- [[Probability Distribution]]
- [[Central Limit Theorem]]
- [[Law of Large Numbers]]
- [[Linear Algebra]]
- [[Benford's Law]]
