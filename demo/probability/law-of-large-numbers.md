---
title: Law of Large Numbers
aliases: [LLN]
tags: [probability]
summary: As you collect more independent samples, their running average converges to the true expected value.
weight: 40
---

# Law of Large Numbers

The **law of large numbers** (LLN) is the promise that makes statistics possible: *average enough independent samples and the average converges to the truth.* Any single measurement is noisy and might mislead you. But the **running mean** — the average of everything you have seen so far — settles down, drifting ever closer to the underlying expected value as the sample count grows. It is the formal version of the intuition that a fair coin flipped many times lands heads "about half" the time.

This is why casinos always win, why polls of more people are sharper, and why a long experiment beats a short one. Randomness does not cancel on any single trial — it cancels in the *aggregate*.

## The statement

For independent samples $X_1, X_2, \ldots$ each with mean $\mu$, the sample average

{{< eq number="1" >}}
\bar X_n = \frac{1}{n}\sum_{i=1}^{n} X_i \;\xrightarrow{\;n\to\infty\;}\; \mu .
{{< /eq >}}

The convergence is real but **slow and shaky**. Early on, a few unlucky draws can swing the average wildly. As $n$ grows the wobble shrinks — at the same $1/\sqrt{n}$ rate that governs a [[Random Walk]] — and the curve flattens onto $\mu$. The [[Central Limit Theorem]] is the fine print here: it says precisely how big the remaining wobble is at each $n$.

{{< note kind="warning" title="No memory, no 'due' " >}}
The LLN does *not* say a coin that landed heads ten times is now "due" for tails — that gambler's fallacy imagines the coin remembers. Past flips have no pull on future ones. The average converges only because new independent draws gradually *outnumber* and dilute the early swing, never because anything corrects it.
{{< /note >}}

## A noisy mean settling onto the truth

Below, a fair process (think a coin paying $+1$ for heads, $-1$ for tails, true mean $0$) is sampled forever. The jagged line is the **running average** after each new draw. Watch it: violent at the left edge where only a handful of samples exist, then steadier and steadier, hugging the dashed line at the true value as the sample count climbs into the thousands.

{{< sketch height="340" caption="The running average of a fair ±1 process converging to its true mean of 0. Early on a few samples swing it wildly; with more samples the wobble shrinks like 1/√n and the line settles onto the dashed truth. Auto-resets." >}}
if (frame === 0 || !state.hist) {
  state.hist = [];     // running-mean history (downsampled to width)
  state.n = 0;
  state.sum = 0;
  state.maxN = 4000;
}
const cs = getComputedStyle(document.documentElement);
const accent = cs.getPropertyValue('--accent').trim() || '#5b9cff';
const accent2 = cs.getPropertyValue('--accent-2').trim() || '#ff9e64';
const faint = cs.getPropertyValue('--text-faint').trim() || 'rgba(255,255,255,0.4)';
const border = cs.getPropertyValue('--border').trim() || 'rgba(255,255,255,0.18)';
ctx.clearRect(0, 0, W, H);
ctx.fillStyle = '#0c1018';
ctx.fillRect(0, 0, W, H);
// add samples this frame
const perFrame = 18;
for (let s = 0; s < perFrame && state.n < state.maxN; s++) {
  state.sum += (Math.random() < 0.5 ? 1 : -1);
  state.n++;
  state.hist.push(state.sum / state.n);
}
const padT = 20, padB = 24, padL = 8, padR = 8;
const plotH = H - padT - padB;
const midY = padT + plotH / 2;
const yScale = plotH / 2 / 1.0; // mean ranges roughly within ±1
// true-value line
ctx.strokeStyle = accent2;
ctx.lineWidth = 1.5;
ctx.setLineDash([6, 5]);
ctx.beginPath(); ctx.moveTo(padL, midY); ctx.lineTo(W - padR, midY); ctx.stroke();
ctx.setLineDash([]);
// shrinking ±1/√n envelope
ctx.strokeStyle = border;
ctx.lineWidth = 1;
for (const sgn of [1, -1]) {
  ctx.beginPath();
  for (let px = 0; px <= W - padL - padR; px += 3) {
    const idx = Math.floor((px / (W - padL - padR)) * (state.hist.length || 1));
    const nn = Math.max(1, idx * (state.maxN / Math.max(1, state.hist.length)));
    const env = sgn * 1.0 / Math.sqrt(nn);
    const y = midY - env * yScale;
    if (px === 0) ctx.moveTo(padL + px, y); else ctx.lineTo(padL + px, y);
  }
  ctx.stroke();
}
// running-mean curve
ctx.strokeStyle = accent;
ctx.lineWidth = 2;
ctx.beginPath();
const L = state.hist.length;
for (let i = 0; i < L; i++) {
  const x = padL + (i / Math.max(1, L - 1)) * (W - padL - padR);
  let y = midY - state.hist[i] * yScale;
  if (y < padT) y = padT; if (y > H - padB) y = H - padB;
  if (i === 0) ctx.moveTo(x, y); else ctx.lineTo(x, y);
}
ctx.stroke();
// labels
ctx.fillStyle = faint;
ctx.font = '12px sans-serif';
ctx.fillText('n = ' + state.n, padL + 4, padT + 4);
ctx.fillText('running average', padL + 4, H - 8);
ctx.fillStyle = accent2;
ctx.fillText('true mean = 0', W - padR - 96, midY - 6);
if (state.n >= state.maxN) {
  state.hold = (state.hold || 0) + 1;
  if (state.hold > 110) { state.hist = []; state.n = 0; state.sum = 0; state.hold = 0; }
}
{{< /sketch >}}

## Why it matters for computing

The LLN is the engine under the [[Monte Carlo Method]]: to estimate a hard quantity, sample it at random many times and average — the LLN guarantees the estimate converges to the right answer, and the [[Central Limit Theorem]] tells you the error bars. Every simulation that "runs it a million times and takes the mean" is cashing in this theorem.

## See also

- [[Central Limit Theorem]]
- [[Monte Carlo Method]]
- [[Probability Distribution]]
