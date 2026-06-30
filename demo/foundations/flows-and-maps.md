---
title: Flows and Maps
aliases: [flows, maps, iterated map, continuous and discrete dynamics]
tags: [foundations]
summary: The two ways time can pass — continuously, as a flow of differential equations, or in discrete steps, as an iterated map.
weight: 100
---

# Flows and Maps

Dynamical systems come in two flavors, set by how time advances. In a **flow**,
time is continuous and the rule is a differential equation, $\dot x = f(x)$; the
state slides smoothly along a trajectory. In a **map**, time ticks in discrete
steps and the rule is an iteration, $x_{n+1} = F(x_n)$; the state hops from one
value to the next. Both describe a [[Dynamical System]] evolving in
[[State Space]] — they are two dialects of the same language.

## Flows: continuous time

A flow integrates a velocity field. Drop a point anywhere and it follows the
arrows of the [[Phase Portrait]], tracing a smooth curve that never crosses
itself. Planetary orbits, the swing of a [[Pendulum]], the convection of the
[[Lorenz System]] — all are flows.

{{< sim name="vectorfield" fx="y" fy="-x - 0.25*y" caption="The velocity field of a damped oscillator, dx/dt = y, dy/dt = -x - 0.25y. A flow is just 'follow the arrows', forever." >}}

## Maps: discrete time

A map applies its rule over and over. Models that are naturally
step-by-step — a population censused once a year, the [[Logistic Map]], the
return of an orbit — live here. Maps are easy to iterate and can be wildly
complex in one dimension, where a flow would need at least three.

{{< note kind="tip" title="A cobweb is a map made visible" >}}
To iterate a map graphically, bounce between its curve $y=F(x)$ and the diagonal
$y=x$: up to the curve to apply the rule, across to the diagonal to feed the
output back as the next input. The staircase you trace is a **cobweb plot**, and
its shape — spiraling in, settling, or filling an interval — reveals the
dynamics at a glance.
{{< /note >}}

## The bridge: Poincaré sections

The two pictures are deeply connected. Take a flow and record the state only each
time it pierces a chosen surface, and you have reduced it to a map — a
[[Poincaré Section]]. This trades a continuous trajectory for a sequence of
points, dropping the dimension by one while preserving the essential dynamics. It
is the standard way to find the [[Limit Cycle]]s and [[Strange Attractor]]s
hidden inside a flow.

{{< sketch height="300" caption="A cobweb plot of the logistic map: bounce between the parabola and the diagonal to iterate. Drag to move the starting point." >}}
const m = 36, w = W - 2 * m, h = H - 2 * m, r = 3.2;
const ac = getComputedStyle(document.documentElement).getPropertyValue('--accent').trim() || '#5b9cff';
const dim = getComputedStyle(document.documentElement).getPropertyValue('--text-faint').trim() || '#888';
const X = v => m + v * w, Y = v => H - m - v * h;
ctx.strokeStyle = dim; ctx.globalAlpha = 0.6; ctx.lineWidth = 1;
ctx.strokeRect(m, m, w, h);
ctx.beginPath(); ctx.moveTo(X(0), Y(0)); ctx.lineTo(X(1), Y(1)); ctx.stroke(); // diagonal
ctx.globalAlpha = 1; ctx.strokeStyle = dim; ctx.lineWidth = 2; ctx.beginPath();
for (let i = 0; i <= 100; i++) { const x = i / 100, y = r * x * (1 - x); i ? ctx.lineTo(X(x), Y(y)) : ctx.moveTo(X(x), Y(y)); }
ctx.stroke();
let x0 = mouse.down ? Math.max(0.01, Math.min(0.99, (mouse.x - m) / w)) : 0.08;
ctx.strokeStyle = ac; ctx.lineWidth = 1.5; ctx.beginPath(); ctx.moveTo(X(x0), Y(0));
let x = x0;
for (let i = 0; i < 60; i++) { const y = r * x * (1 - x); ctx.lineTo(X(x), Y(y)); ctx.lineTo(X(y), Y(y)); x = y; }
ctx.stroke();
{{< /sketch >}}

## See also

- [[Poincaré Section]]
- [[Logistic Map]]
- [[Dynamical System]]
