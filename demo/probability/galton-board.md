---
title: Galton Board
aliases: [bean machine, quincunx]
tags: [probability]
summary: A board of pegs that turns a stream of falling balls into a bell-shaped pile — the central limit theorem made physical.
weight: 90
---

# Galton Board

A **Galton board** — also called a *bean machine* or *quincunx* — is the central limit theorem you can hold in your hands. Balls drop from a single point at the top and tumble through a triangular lattice of pegs. At each peg a ball bounces left or right with equal chance, then falls to the next row and bounces again. After many rows the balls collect in bins along the bottom — and they always pile up into the same shape: a **bell curve**.

Invented by Sir Francis Galton in the 1870s, it is a mechanical proof that order emerges from accumulated randomness. No ball's path is predictable, yet the *pile* is utterly reliable.

## Each ball is a random walk

Follow one ball. At each of the $n$ rows it steps left ($-1$) or right ($+1$) with probability $\tfrac{1}{2}$ — exactly a [[Random Walk]]. Its final bin is the sum of those $n$ independent $\pm 1$ choices, so the number of balls landing in each bin follows a **binomial distribution**. And a binomial with many trials is, by the [[Central Limit Theorem]], approximately **normal**:

{{< eq number="1" >}}
\text{bin position} = \sum_{i=1}^{n} s_i, \qquad s_i = \pm 1 \text{ each with prob } \tfrac{1}{2}.
{{< /eq >}}

The center bins fill fastest because there are *many* left/right combinations that cancel out to land there, but only *one* path — all-right — that reaches the far edge. Counting paths is counting the binomial coefficients of Pascal's triangle, and their smooth envelope is the bell.

{{< note kind="key" title="Why the bell is inevitable" >}}
Every ball is a sum of $n$ independent equal coin-flips. Sum enough independent things and the CLT forces a normal distribution — no matter that each flip is a crude two-way choice. The Galton board is that theorem rendered in wood and bouncing beads. Widen it (more rows) and the pile only hugs the bell more tightly.
{{< /note >}}

## Watch the bell build

Below, balls fall one after another, bouncing left or right at each peg and settling into bins at the bottom. Any single ball's route is a coin-flip cascade — but watch the bins. They fill into the unmistakable bell curve, taller in the middle, thinner at the edges, growing smoother with every ball.

{{< sketch height="420" caption="A Galton board: balls bounce left/right at each peg (a random walk) and stack into bins. The pile grows into a binomial/normal bell curve — the central limit theorem made physical. Resets when the bins fill." >}}
if (frame === 0 || !state.balls) {
  state.rows = 10;
  state.bins = state.rows + 1;
  state.counts = new Array(state.bins).fill(0);
  state.balls = [];
  state.spawnTimer = 0;
  state.totalDropped = 0;
  state.maxTotal = 500;
}
const cs = getComputedStyle(document.documentElement);
const accent = cs.getPropertyValue('--accent').trim() || '#5b9cff';
const accent2 = cs.getPropertyValue('--accent-2').trim() || '#ff9e64';
const good = cs.getPropertyValue('--good').trim() || '#9ece6a';
const faint = cs.getPropertyValue('--text-faint').trim() || 'rgba(255,255,255,0.4)';
const border = cs.getPropertyValue('--border').trim() || 'rgba(255,255,255,0.18)';
ctx.clearRect(0, 0, W, H);
ctx.fillStyle = '#0c1018';
ctx.fillRect(0, 0, W, H);
const rows = state.rows;
const topY = 30;
const pegGapY = (H * 0.52) / rows;
const pegGapX = Math.min(W / (rows + 3), pegGapY * 1.15);
const cx = W / 2;
const binTop = topY + rows * pegGapY + 14;
const binAreaH = H - binTop - 10;
// draw pegs
ctx.fillStyle = border;
for (let r = 0; r < rows; r++) {
  for (let c = 0; c <= r; c++) {
    const px = cx + (c - r / 2) * pegGapX;
    const py = topY + r * pegGapY;
    ctx.beginPath(); ctx.arc(px, py, 2.2, 0, 7); ctx.fill();
  }
}
// spawn a ball periodically
state.spawnTimer++;
if (state.spawnTimer >= 7 && state.totalDropped < state.maxTotal && state.balls.length < 30) {
  state.spawnTimer = 0;
  state.totalDropped++;
  state.balls.push({ x: cx, y: topY - 14, row: -1, offset: 0, vy: 0, done: false });
}
// bin geometry
const binW = pegGapX;
const leftEdge = cx - (state.bins / 2) * binW;
function binX(i) { return leftEdge + (i + 0.5) * binW; }
// update + draw balls
ctx.fillStyle = accent;
for (let b = state.balls.length - 1; b >= 0; b--) {
  const ball = state.balls[b];
  ball.y += 2.4;
  if (ball.row < rows - 1 && ball.y >= topY + (ball.row + 1) * pegGapY) {
    ball.row++;
    ball.offset += (Math.random() < 0.5 ? -0.5 : 0.5);
    ball.x = cx + ball.offset * pegGapX;
  } else if (ball.row >= rows - 1 && !ball.settled) {
    // settle into a bin
    if (ball.y >= binTop) {
      const idx = Math.round(ball.offset + rows / 2);
      const i = Math.max(0, Math.min(state.bins - 1, idx));
      state.counts[i]++;
      state.balls.splice(b, 1);
      continue;
    }
  }
  // smoothly track x toward target
  const targetX = cx + ball.offset * pegGapX;
  ball.x += (targetX - ball.x) * 0.3;
  ctx.beginPath(); ctx.arc(ball.x, ball.y, 3, 0, 7); ctx.fill();
}
// draw bins / histogram
let cmax = 1; for (let i = 0; i < state.bins; i++) cmax = Math.max(cmax, state.counts[i]);
const colH = binAreaH;
for (let i = 0; i < state.bins; i++) {
  const h = (state.counts[i] / cmax) * colH;
  ctx.fillStyle = accent2;
  ctx.fillRect(binX(i) - binW / 2 + 1, H - 10 - h, binW - 2, h);
}
// bin dividers
ctx.strokeStyle = border;
ctx.lineWidth = 1;
for (let i = 0; i <= state.bins; i++) {
  const x = leftEdge + i * binW;
  ctx.beginPath(); ctx.moveTo(x, binTop); ctx.lineTo(x, H - 10); ctx.stroke();
}
// overlay ideal bell (binomial mean rows/2, var rows/4)
const mean = state.bins / 2 - 0.5;
const sd = Math.sqrt(rows) / 2;
ctx.strokeStyle = good;
ctx.lineWidth = 2;
ctx.beginPath();
let peak = 1 / (sd * Math.sqrt(2 * Math.PI));
for (let i = 0; i < state.bins; i += 0.1) {
  const z = (i - mean) / sd;
  const dens = Math.exp(-0.5 * z * z) / (sd * Math.sqrt(2 * Math.PI));
  const x = binX(i);
  const y = H - 10 - (dens / peak) * colH;
  if (i === 0) ctx.moveTo(x, y); else ctx.lineTo(x, y);
}
ctx.stroke();
// label
ctx.fillStyle = faint;
ctx.font = '12px sans-serif';
ctx.fillText('balls dropped: ' + state.totalDropped, 12, H - 14);
if (state.totalDropped >= state.maxTotal && state.balls.length === 0) {
  state.hold = (state.hold || 0) + 1;
  if (state.hold > 130) {
    state.counts = new Array(state.bins).fill(0);
    state.totalDropped = 0; state.hold = 0;
  }
}
{{< /sketch >}}

## A machine you can read two ways

The Galton board is a rare object that is simultaneously a physical experiment and a mathematical proof. As physics, it is gravity and elastic bounces. As mathematics, it is the binomial distribution converging to the normal — the [[Central Limit Theorem]] and the [[Probability Distribution]] of a [[Random Walk]] stacked into a literal pile. Galton built it to argue that complex, bell-shaped variation in nature need not have a complex cause: many tiny independent nudges suffice.

## See also

- [[Central Limit Theorem]]
- [[Probability Distribution]]
- [[Random Walk]]
