---
title: Brownian Motion
aliases: [wiener process, brownian]
tags: [probability]
summary: The continuous-time limit of a random walk — the erratic, jittery path traced by a particle buffeted by countless tiny random kicks.
weight: 60
---

# Brownian Motion

**Brownian motion** is what a [[Random Walk]] becomes when you take its steps infinitely small and infinitely frequent. The result is a continuous but impossibly jagged path — nowhere smooth, jittering at every scale, yet governed by clean statistical laws. It is named for the botanist Robert Brown, who in 1827 watched pollen grains in water twitch under his microscope without apparent cause. The cause, Einstein later showed, was the grain being kicked from all sides by unseen water molecules — and the math of those kicks is this process.

As the mathematical idealization (the **Wiener process**), Brownian motion is the cornerstone of stochastic calculus, the model for diffusion, and the noisy heartbeat behind everything from stock prices to the spread of heat.

## From discrete steps to a continuous thread

Take a random walk of $n$ steps in a fixed time $T$. Shrink each step and crank up the count, and in the limit the corner-to-corner path smooths into a continuous curve $B(t)$ with three defining properties:

- **Independent increments** — disjoint stretches of the path are independent of one another.
- **Gaussian increments** — the displacement over a time interval $\Delta t$ is normally distributed with variance proportional to $\Delta t$.
- **Continuity** — the path never jumps, yet is differentiable nowhere.

That middle property is the inherited $\sqrt{n}$ law of the random walk: the typical distance traveled grows like $\sqrt{t}$. Spread accumulates with the *square root* of time, the universal signature of diffusion.

{{< note kind="note" title="Why nowhere smooth?" >}}
A smooth path has a well-defined velocity. But Brownian motion's displacement over $\Delta t$ scales like $\sqrt{\Delta t}$, so the "velocity" $\sqrt{\Delta t}/\Delta t = 1/\sqrt{\Delta t}$ blows up as $\Delta t \to 0$. There is no instantaneous speed — only ceaseless, scale-free jitter. Zoom in on any piece and it looks just as rough as the whole, a statistical [[Fractal]].
{{< /note >}}

## Wandering in the plane

Below, several particles undergo Brownian motion in 2D: each takes a small Gaussian-distributed kick every frame, tracing a tangled trail. The paths never quite repeat, never settle, and never smooth out — exactly the restless wandering Brown saw in his pollen.

{{< sketch height="380" caption="Several particles performing 2D Brownian motion. Each takes an independent Gaussian step every frame; the trails are continuous but jagged at every scale. Auto-resets when the trails grow long." >}}
if (frame === 0 || !state.parts) {
  state.K = 5;
  state.parts = [];
  for (let i = 0; i < state.K; i++) {
    state.parts.push({ x: W / 2, y: H / 2, trail: [] });
  }
  state.age = 0;
}
const cs = getComputedStyle(document.documentElement);
const accent = cs.getPropertyValue('--accent').trim() || '#5b9cff';
const accent2 = cs.getPropertyValue('--accent-2').trim() || '#ff9e64';
const good = cs.getPropertyValue('--good').trim() || '#9ece6a';
const palette = [accent, accent2, good, '#bb9af7', '#7dcfff'];
ctx.clearRect(0, 0, W, H);
ctx.fillStyle = '#0c1018';
ctx.fillRect(0, 0, W, H);
const sigma = Math.min(W, H) * 0.012;
state.age++;
// gaussian via Box-Muller, advance + record trail
function gauss() {
  let u = 0, v = 0;
  while (u === 0) u = Math.random();
  while (v === 0) v = Math.random();
  return Math.sqrt(-2 * Math.log(u)) * Math.cos(2 * Math.PI * v);
}
for (let i = 0; i < state.K; i++) {
  const p = state.parts[i];
  p.x += gauss() * sigma;
  p.y += gauss() * sigma;
  // soft reflective bounds
  if (p.x < 4) p.x = 4 + (4 - p.x); if (p.x > W - 4) p.x = (W - 4) - (p.x - (W - 4));
  if (p.y < 4) p.y = 4 + (4 - p.y); if (p.y > H - 4) p.y = (H - 4) - (p.y - (H - 4));
  p.trail.push([p.x, p.y]);
  if (p.trail.length > 600) p.trail.shift();
}
// draw trails
for (let i = 0; i < state.K; i++) {
  const tr = state.parts[i].trail;
  ctx.strokeStyle = palette[i % palette.length];
  ctx.globalAlpha = 0.55;
  ctx.lineWidth = 1.4;
  ctx.beginPath();
  for (let j = 0; j < tr.length; j++) {
    if (j === 0) ctx.moveTo(tr[j][0], tr[j][1]); else ctx.lineTo(tr[j][0], tr[j][1]);
  }
  ctx.stroke();
  ctx.globalAlpha = 1;
  // head
  ctx.fillStyle = palette[i % palette.length];
  ctx.beginPath(); ctx.arc(state.parts[i].x, state.parts[i].y, 3, 0, 7); ctx.fill();
}
if (state.age > 620) {
  state.hold = (state.hold || 0) + 1;
  if (state.hold > 60) { state.parts = null; state.hold = 0; }
}
{{< /sketch >}}

## Diffusion, prices, and patterns

Average over a whole cloud of Brownian particles and the sharp randomness smooths into the **diffusion equation** — the same heat-spreading law that, paired with local reaction rules, drives the spots and stripes of [[Reaction–Diffusion]] systems. Add a steady drift and you get the standard model of a fluctuating stock price. Brownian motion is where probability stops being about discrete dice and becomes the continuous calculus of noise.

## See also

- [[Random Walk]]
- [[Reaction–Diffusion]]
- [[Central Limit Theorem]]
