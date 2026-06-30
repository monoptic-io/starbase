---
title: Central Limit Theorem
aliases: [CLT]
tags: [probability]
summary: The sum or average of many independent random quantities tends to a normal distribution, no matter what shape the individual quantities have.
weight: 30
---

# Central Limit Theorem

The **central limit theorem** (CLT) is the reason the bell curve rules the world. It makes a claim that sounds too strong to be true: *add up many independent random quantities — whatever their individual shapes — and their sum is approximately **normal**.* The pieces can be lopsided, flat, spiky, or weird; pile up enough of them and the lumps wash out into the same smooth, symmetric bell. The shape of the parts does not survive; only the bell does.

This is why heights, measurement errors, test scores, and the means of almost any survey come out Gaussian. Each is a sum of countless small independent influences, and the CLT erases the fingerprints of the individual causes.

## The statement

Let $X_1, X_2, \ldots, X_n$ be independent draws from *any* distribution with mean $\mu$ and finite variance $\sigma^2$. Form their average $\bar X_n$. Then as $n$ grows, the standardized average

{{< eq number="1" >}}
\frac{\bar X_n - \mu}{\sigma / \sqrt{n}} \;\xrightarrow{\;n\to\infty\;}\; \mathcal{N}(0,1),
{{< /eq >}}

a standard normal distribution. Two facts hide inside this. First, the average **converges** to $\mu$ — that is the [[Law of Large Numbers]]. Second, the *fluctuations* around $\mu$ shrink like $1/\sqrt{n}$ and, rescaled, take the universal bell shape. The CLT is the law of large numbers' more refined companion: it describes not just where the average lands, but the exact statistics of how it gets there.

{{< note kind="key" title="Why a bell, and why always the same one?" >}}
Convolving (adding) distributions smooths them. Each sum averages out the bumps of the last, and the normal distribution is the unique *fixed shape* of that smoothing — the one form that, added to itself, stays itself (just wider). Sums flow toward it like water finding its level.
{{< /note >}}

## Watch a bell emerge

Below, the machine repeatedly draws $k$ independent random numbers (each uniform on $0$ to $1$), averages them, and drops the result into a histogram. With $k=1$ the histogram is **flat** — a single uniform draw is equally likely anywhere. Nudge $k$ up and the histogram sharpens into a bell; the more numbers you average, the narrower and more Gaussian it gets. **Move your mouse left and right to set $k$** (from $1$ to about $12$) and watch a flat slab fold itself into a bell.

{{< sketch height="380" caption="The CLT in action. Each sample is the average of k independent uniform random numbers; the histogram counts millions of such averages. Mouse x sets k: at k=1 it is flat, and it sharpens to a bell as k rises. A reference normal curve is overlaid." >}}
if (frame === 0 || !state.hist) {
  state.bins = 90;
  state.hist = new Array(state.bins).fill(0);
  state.total = 0;
  state.k = 1;
}
const cs = getComputedStyle(document.documentElement);
const accent = cs.getPropertyValue('--accent').trim() || '#5b9cff';
const accent2 = cs.getPropertyValue('--accent-2').trim() || '#ff9e64';
const good = cs.getPropertyValue('--good').trim() || '#9ece6a';
const faint = cs.getPropertyValue('--text-faint').trim() || 'rgba(255,255,255,0.4)';
// mouse sets k (1..12)
let k = 1 + Math.floor((mouse.x / W) * 12);
if (!isFinite(k) || k < 1) k = 1; if (k > 12) k = 12;
if (k !== state.k) { state.k = k; state.hist = new Array(state.bins).fill(0); state.total = 0; }
// add a batch of samples each frame
const batch = 400;
for (let s = 0; s < batch; s++) {
  let sum = 0;
  for (let j = 0; j < state.k; j++) sum += Math.random();
  const avg = sum / state.k; // in [0,1], mean 0.5
  let b = Math.floor(avg * state.bins);
  if (b < 0) b = 0; if (b >= state.bins) b = state.bins - 1;
  state.hist[b]++;
  state.total++;
}
ctx.clearRect(0, 0, W, H);
ctx.fillStyle = '#0c1018';
ctx.fillRect(0, 0, W, H);
const padB = 34, padT = 16;
const plotH = H - padB - padT;
const binW = W / state.bins;
let hmax = 1; for (let b = 0; b < state.bins; b++) hmax = Math.max(hmax, state.hist[b]);
// histogram bars
ctx.fillStyle = accent;
for (let b = 0; b < state.bins; b++) {
  const h = (state.hist[b] / hmax) * plotH;
  ctx.fillRect(b * binW + 0.5, H - padB - h, binW - 1, h);
}
// overlay theoretical normal: mean 0.5, var = (1/12)/k
const variance = (1 / 12) / state.k;
const sd = Math.sqrt(variance);
ctx.strokeStyle = accent2;
ctx.lineWidth = 2;
ctx.beginPath();
let peak = 1 / (sd * Math.sqrt(2 * Math.PI));
for (let px = 0; px <= W; px += 2) {
  const xv = px / W; // 0..1
  const z = (xv - 0.5) / sd;
  const dens = Math.exp(-0.5 * z * z) / (sd * Math.sqrt(2 * Math.PI));
  const y = H - padB - (dens / peak) * plotH;
  if (px === 0) ctx.moveTo(px, y); else ctx.lineTo(px, y);
}
ctx.stroke();
// labels
ctx.fillStyle = good;
ctx.font = 'bold 14px sans-serif';
ctx.fillText('k = ' + state.k + (state.k === 1 ? '  (uniform — flat)' : '  (average of ' + state.k + ')'), 12, padT + 6);
ctx.fillStyle = faint;
ctx.font = '12px sans-serif';
ctx.fillText('← move mouse to change k →   samples: ' + state.total.toLocaleString(), 12, H - 12);
{{< /sketch >}}

## It is not magic — there are rules

The CLT needs **independence** (or near enough) and a **finite variance**. Quantities with wild, heavy tails — where rare giant values dominate — can break it, converging instead to other "stable" laws. And the parts must each be small relative to the sum; one overwhelming term keeps its own shape. Within those bounds, though, the theorem is astonishingly forgiving, which is exactly why the normal distribution shows up so relentlessly. The [[Galton Board]] makes the same convergence physical, and the [[Probability Distribution]] page collects the cast of characters.

## See also

- [[Galton Board]]
- [[Law of Large Numbers]]
- [[Probability Distribution]]
