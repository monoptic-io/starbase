---
title: Bifurcation
aliases: [bifurcations, bifurcation diagram]
tags: [foundations]
summary: A qualitative change in a system's behavior as a parameter crosses a critical value — equilibria born, destroyed, or swapping stability.
weight: 80
---

# Bifurcation

A **bifurcation** is a sudden, qualitative change in the long-term behavior of a
[[Dynamical System]] as a control parameter is varied smoothly. For most values
of the parameter, nudging it a little changes the dynamics only a little. But at
isolated **critical values**, the structure of the [[Phase Portrait]] reorganizes:
a [[Fixed Point]] appears or disappears, an equilibrium loses [[Stability]], or a
steady state gives way to an oscillation. The parameter moved continuously; the
behavior jumped.

## The common local bifurcations

Three patterns account for most one-parameter bifurcations:

- **Saddle-node** — two fixed points (one stable, one unstable) collide and
  annihilate. Below the threshold there are two equilibria; above it, none. This
  is how steady states are created and destroyed.
- **Pitchfork** — a single stable equilibrium loses stability and gives birth to
  two new stable ones, symmetric about it. A ruler pressed end-on stays straight
  until it buckles left *or* right.
- **Hopf** — a fixed point loses stability and throws off a [[Limit Cycle]]: the
  system stops resting and starts oscillating. This is the birth of rhythm.

{{< note kind="key" title="The normal form" >}}
The pitchfork is captured by $\dot x = rx - x^3$. For $r<0$ the only equilibrium
is $x^*=0$ (stable). As $r$ passes zero it goes unstable and two new branches
$x^* = \pm\sqrt{r}$ peel away. The diagram below is exactly this $\pm\sqrt{r}$.
{{< /note >}}

## A bifurcation diagram

Plot the locations of the equilibria against the parameter $r$ and you get a
**bifurcation diagram** — a map of the system's repertoire. Here is the
supercritical pitchfork: one branch splits cleanly into two at $r=0$.

{{< sketch height="320" caption="Pitchfork bifurcation: the stable equilibrium at zero splits into two as r crosses zero. Solid = stable, dashed = unstable." >}}
const ox = 60, oy = H / 2, sx = (W - 80) / 2, sy = H * 0.42;
const ac = getComputedStyle(document.documentElement).getPropertyValue('--accent').trim() || '#5b9cff';
const dim = getComputedStyle(document.documentElement).getPropertyValue('--text-faint').trim() || '#888';
ctx.strokeStyle = dim; ctx.globalAlpha = 0.5; ctx.lineWidth = 1;
ctx.beginPath(); ctx.moveTo(ox, 10); ctx.lineTo(ox, H - 16); ctx.moveTo(ox, oy); ctx.lineTo(W - 16, oy); ctx.stroke();
ctx.globalAlpha = 1; ctx.fillStyle = dim; ctx.font = '12px sans-serif';
ctx.fillText('x*', ox - 26, 16); ctx.fillText('r', W - 22, oy - 6);
// unstable branch x*=0 for r>0 (dashed)
ctx.strokeStyle = dim; ctx.setLineDash([4, 4]); ctx.lineWidth = 1.5;
ctx.beginPath(); ctx.moveTo(ox, oy); ctx.lineTo(W - 16, oy); ctx.stroke();
ctx.setLineDash([]);
// stable branch x*=0 for r<0, then +/- sqrt(r)
ctx.strokeStyle = ac; ctx.lineWidth = 2.5;
ctx.beginPath(); ctx.moveTo(ox - 0, oy); ctx.lineTo(ox, oy); ctx.stroke();
for (let s of [1, -1]) {
  ctx.beginPath();
  for (let i = 0; i <= 100; i++) {
    const r = i / 100;                 // 0..1
    const x = s * Math.sqrt(r);
    const px = ox + r * sx, py = oy - x * sy;
    i ? ctx.lineTo(px, py) : ctx.moveTo(px, py);
  }
  ctx.stroke();
}
// flat stable branch for r<0
ctx.beginPath(); ctx.moveTo(ox - sx, oy); ctx.lineTo(ox, oy); ctx.stroke();
{{< /sketch >}}

## Why it matters

Bifurcations are where systems *change character*: lasers switch on, fluids
begin to convect, populations start to cycle, a heartbeat turns arrhythmic. A
cascade of period-doubling bifurcations is one of the standard roads to
[[Chaos]] — the [[Logistic Map]] travels exactly this route.

## See also

- [[Logistic Map]]
- [[Stability]]
- [[Limit Cycle]]
