---
title: Random Walk
aliases: [random walks, drunkard's walk]
tags: [probability]
summary: A path built by adding up independent random steps; its typical distance from the start grows like the square root of the number of steps.
weight: 10
---

# Random Walk

A **random walk** is the simplest random process there is: start at a point, take a step in a random direction, then another, then another. Each step forgets the past — it does not matter how you got here, only that you take one more independent jump. The classic picture is the **drunkard's walk**: a stumbling figure who flips a coin at every lamppost to decide left or right. Where do they end up?

The surprising answer is that the walk has a precise statistical shape even though every single step is unpredictable. On average the drunkard goes nowhere — left and right cancel. But their *typical distance* from the start grows steadily, and at a very specific rate.

## Spread grows like the square root

After $n$ steps of size $1$, the position $X_n = s_1 + s_2 + \cdots + s_n$ is a sum of independent $\pm 1$ steps. Its mean is zero, but its variance adds up: each step contributes $1$, so $\operatorname{Var}(X_n) = n$. The typical distance is the standard deviation:

{{< eq number="1" >}}
\sqrt{\operatorname{Var}(X_n)} \;=\; \sqrt{n}.
{{< /eq >}}

This $\sqrt{n}$ law is the signature of diffusion. To get twice as far from home you need *four* times as many steps. It is why ink spreads slowly through still water, why a smell takes minutes to cross a room, and why a stock price wanders the way it does. The walk explores, but it explores reluctantly.

{{< note kind="note" title="Sum of steps, shape of a bell" >}}
Because $X_n$ is a sum of many independent steps, the [[Central Limit Theorem]] guarantees that its distribution approaches a **normal** (Gaussian) curve, with width $\sqrt{n}$. That is exactly the bell building up at the bottom of the sketch below — diffusion and the bell curve are the same fact seen two ways.
{{< /note >}}

## A crowd of walkers spreading out

Below, hundreds of walkers all start at the center line and step left or right at random each frame. Two things happen at once. The **cloud spreads** — its edges creep outward like $\sqrt{n}$, slowing as it widens. And the **histogram at the bottom** fills in: it counts how many walkers sit at each horizontal position, and it grows into a smooth bell curve, sharp in the middle, thin at the tails.

{{< sketch height="380" caption="Hundreds of independent random walkers spreading from the center. Top: each dot is one walker stepping left/right. Bottom: a live histogram of their positions building into a bell curve of width √n. Auto-resets." >}}
if (frame === 0 || !state.walkers) {
  state.N = 600;
  state.walkers = new Array(state.N).fill(0).map(() => 0);
  state.step = 0;
  state.maxSteps = 240;
}
const cs = getComputedStyle(document.documentElement);
const accent = cs.getPropertyValue('--accent').trim() || '#5b9cff';
const accent2 = cs.getPropertyValue('--accent-2').trim() || '#ff9e64';
const faint = cs.getPropertyValue('--text-faint').trim() || 'rgba(255,255,255,0.4)';
const border = cs.getPropertyValue('--border').trim() || 'rgba(255,255,255,0.18)';
ctx.clearRect(0, 0, W, H);
ctx.fillStyle = '#0c1018';
ctx.fillRect(0, 0, W, H);
const cx = W / 2;
const cloudH = H * 0.52;
const stepPx = Math.max(1, W / 240) * 0.9;
// advance: one random step per walker per frame
if (state.step < state.maxSteps) {
  for (let i = 0; i < state.N; i++) {
    state.walkers[i] += (Math.random() < 0.5 ? -1 : 1);
  }
  state.step++;
}
// center guide line
ctx.strokeStyle = border;
ctx.lineWidth = 1;
ctx.beginPath(); ctx.moveTo(cx, 0); ctx.lineTo(cx, cloudH); ctx.stroke();
// draw walker cloud: y by index band, x by position
ctx.fillStyle = accent;
ctx.globalAlpha = 0.5;
for (let i = 0; i < state.N; i++) {
  const x = cx + state.walkers[i] * stepPx;
  const y = 8 + (i / state.N) * (cloudH - 16);
  if (x > 0 && x < W) { ctx.fillRect(x, y, 2, 2); }
}
ctx.globalAlpha = 1;
// histogram at bottom
const bins = 81;
const hist = new Array(bins).fill(0);
const binW = W / bins;
for (let i = 0; i < state.N; i++) {
  const x = cx + state.walkers[i] * stepPx;
  let b = Math.floor(x / binW);
  if (b < 0) b = 0; if (b >= bins) b = bins - 1;
  hist[b]++;
}
let hmax = 1; for (let b = 0; b < bins; b++) hmax = Math.max(hmax, hist[b]);
const histTop = cloudH + 14;
const histH = H - histTop - 8;
ctx.fillStyle = accent2;
for (let b = 0; b < bins; b++) {
  const h = (hist[b] / hmax) * histH;
  ctx.fillRect(b * binW + 0.5, H - 8 - h, binW - 1, h);
}
// labels
ctx.fillStyle = faint;
ctx.font = '12px sans-serif';
ctx.fillText('step ' + state.step + '   width ≈ √n = ' + Math.sqrt(state.step).toFixed(1), 10, H - histH - 6);
if (state.step >= state.maxSteps) {
  state.hold = (state.hold || 0) + 1;
  if (state.hold > 90) { state.walkers = null; state.hold = 0; }
}
{{< /sketch >}}

## Walks in more dimensions

Nothing forces the walk to a line. A walker on a plane or in space follows the same $\sqrt{n}$ spreading, and in the continuous limit — infinitely many infinitely small steps — the path becomes [[Brownian Motion]], the mathematically exact model of a diffusing particle. A random walk is also the simplest interesting [[Markov Chain]]: the next position depends only on the current one, never on the route taken to reach it.

## See also

- [[Brownian Motion]]
- [[Central Limit Theorem]]
- [[Markov Chain]]
