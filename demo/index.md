---
title: Dynamical Systems
summary: An explorable field guide to how things change over time — from a swinging pendulum to deterministic chaos and emergent complexity.
tags: [overview]
---

# Dynamical Systems

A **dynamical system** is anything whose state evolves in time according to a
fixed rule. That single idea — a *state* plus an *evolution rule* — is enough to
describe a swinging [[Pendulum]], the weather, planetary orbits, populations of
predators and prey, and the flicker of neurons. This knowledge base is a guided
tour through the field, with live simulations you can poke on almost every page.

{{< sim name="lorenz" caption="The Lorenz attractor — a deterministic system that never repeats. Watch it trace out its butterfly." >}}

## Six ways in

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
