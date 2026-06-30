---
title: Systems, Signals & Computation
summary: An interactive field guide to how things change, oscillate, propagate, and compute — from a swinging pendulum to chaos, waves, networks, and learning machines.
tags: [overview]
---

# Systems, Signals & Computation

An interactive field guide to the mathematics of change and pattern. It starts
with **dynamical systems** — anything whose state evolves by a fixed rule, like a
swinging [[Pendulum]] — and follows the threads outward: into the **oscillations**
and **[[Waves]]** that motion makes, the **[[Fourier Analysis]]** that decomposes
any signal, the **[[Chaos]]** that simple rules can hide, and the **networks**,
**randomness**, **information**, and **optimization** that turn dynamics into
computation. Almost every page has a live simulation you can poke.

{{< sim name="lorenz" caption="The Lorenz attractor — a deterministic system that never repeats. Watch it trace out its butterfly." >}}

## Ways in

{{< columns count="2" >}}
**[[Foundations]]** — the language of the field: [[State Space]],
[[Phase Portrait]], [[Fixed Point]], [[Stability]], [[Attractor]], and the
[[Lyapunov Exponent]] that measures chaos.

**[[Oscillations]]** — things that repeat: the [[Simple Harmonic Oscillator]],
the [[Pendulum]], [[Resonance]], and self-sustaining [[Limit Cycle]]s.

**[[Waves]]** — oscillation set loose in space: the [[Wave Equation]],
[[Standing Wave]]s, [[Interference]], [[Beats]], and the [[Doppler Effect]].

**[[Fourier Analysis]]** — every signal as a sum of sinusoids: [[Fourier Series]],
the [[Fourier Transform]], [[Harmonics]], and the [[Uncertainty Principle]].

**[[Chaos]]** — deterministic yet unpredictable: the [[Lorenz System]], the
[[Double Pendulum]], the [[Logistic Map]], and [[Strange Attractor]]s.

**[[Complex Systems]]** — many simple parts, emergent wholes: the
[[N-Body Problem]], [[Conway's Game of Life]], [[Predator–Prey Dynamics]], and
[[Reaction–Diffusion]] patterns.

**[[Linear Algebra]]** — the math under all of it: [[Vector]]s, [[Matrix|matrices]]
as [[Linear Transformation]]s, and the [[Eigenvalues and Eigenvectors]] that decide
[[Stability]] and rank the web.

**[[Graph Theory & Networks]]** — nodes and edges everywhere: [[Breadth-First Search]],
[[Dijkstra's Algorithm]], [[PageRank]], and [[Small-World Network]]s.
{{< /columns >}}

{{< note kind="tip" title="How to explore" >}}
Follow the **links** between topics — every page suggests *Related topics* and
shows what *References* it. Drag, pause, and reset the simulations.
There is no single path; wander.
{{< /note >}}

## A taste of the whole field in one picture

Below, the same rule runs from many slightly different starting points. They
track together, then diverge — the essence of [[Sensitive Dependence on Initial Conditions]].

{{< sketch height="300" caption="Many pendulums, released a hair apart, fall out of step." >}}
if (frame === 0) {
  state.ps = [];
  for (let i = 0; i < 24; i++) state.ps.push({ th: 2.2 + i * 0.002, w: 0 });
}
const g = 9.81, L = 1.6, ox = W / 2, oy = 18, len = Math.min(H * 0.8, W * 0.35);
ctx.lineWidth = 1.5;
for (let i = 0; i < state.ps.length; i++) {
  const p = state.ps[i];
  const a = -g / L * Math.sin(p.th);
  p.w += a * dt; p.th += p.w * dt;
  const x = ox + len * Math.sin(p.th), y = oy + len * Math.cos(p.th);
  const hue = 210 + i * 4;
  ctx.strokeStyle = 'hsl(' + hue + ',70%,65%)';
  ctx.globalAlpha = 0.7;
  ctx.beginPath(); ctx.moveTo(ox, oy); ctx.lineTo(x, y); ctx.stroke();
  ctx.globalAlpha = 1;
  ctx.fillStyle = 'hsl(' + hue + ',70%,65%)';
  ctx.beginPath(); ctx.arc(x, y, 4, 0, 7); ctx.fill();
}
{{< /sketch >}}

Start with the [[Foundations]], or jump straight into the [[Lorenz System]].
