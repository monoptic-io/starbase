---
title: Lyapunov Exponent
aliases: [lyapunov exponents, largest lyapunov exponent]
tags: [foundations]
summary: The average rate at which nearby trajectories pull apart — the quantitative fingerprint of chaos.
weight: 90
---

# Lyapunov Exponent

The **Lyapunov exponent** $\lambda$ measures how fast two trajectories that
start infinitesimally close together separate as time passes. If their initial
gap is $\delta_0$, then on average it grows (or shrinks) like

{{< eq number="1" >}}
|\delta(t)| \approx |\delta_0|\,e^{\lambda t}, \qquad
\lambda = \lim_{t\to\infty}\frac{1}{t}\ln\frac{|\delta(t)|}{|\delta_0|}.
{{< /eq >}}

The **sign** of $\lambda$ tells you everything about the system's predictability:

- $\lambda < 0$ — neighbors converge; the system forgets its initial condition
  and settles onto a [[Fixed Point]] or [[Limit Cycle]].
- $\lambda = 0$ — neighbors drift apart only linearly; marginal, as on a torus.
- $\lambda > 0$ — neighbors diverge exponentially. This is
  [[Sensitive Dependence on Initial Conditions]], the defining mark of
  [[Chaos]].

## Exponential divergence, seen

The gap between two chaotic trajectories doesn't creep — it explodes. The plot
shows the same initial error of $10^{-4}$ growing under a positive exponent. On
a log scale it would be a straight line of slope $\lambda$; on this linear scale
it looks like nothing, then suddenly everything.

{{< plot fn="0.0001*Math.exp(0.9*x)" xmin="0" xmax="12" ymin="0" ymax="6" title="A tiny error amplified by a positive Lyapunov exponent" caption="An undetectable initial difference becomes order-1 in finite time." >}}

## The predictability horizon

A positive exponent imposes a hard limit on forecasting. If you know the state
to precision $\delta_0$ and can tolerate error up to $\Delta$, your predictions
stay useful only until the **Lyapunov time** $t_h \approx \frac{1}{\lambda}\ln\frac{\Delta}{\delta_0}$.
Because the dependence on your initial precision is *logarithmic*, buying a
thousand times better measurements only buys you a little more forecast time —
the reason weather prediction hits a wall a couple of weeks out.

{{< note kind="key" title="A spectrum, not a number" >}}
An $n$-dimensional system has $n$ Lyapunov exponents, one per direction of
stretching or squeezing. The **largest** decides chaos. Their **sum** is the
rate at which a blob of states changes volume — negative for any dissipative
system with an [[Attractor]]. Positive largest exponent **and** shrinking volume
together force the folded geometry of a [[Strange Attractor]].
{{< /note >}}

## Many trajectories, one fate

Watch a fan of nearby starts under a chaotic flow. They march together
convincingly — and then, at no particular moment, they don't.

{{< sketch height="280" caption="Twelve trajectories of a chaotic map, started within a pixel of each other, losing all correlation." >}}
if (frame === 0) {
  state.xs = [];
  for (let i = 0; i < 12; i++) state.xs.push(0.5 + i * 1e-4);
  state.hist = state.xs.map(() => []);
  state.k = 0;
}
const r = 3.95;
if (frame % 2 === 0 && state.k < W) {
  for (let i = 0; i < state.xs.length; i++) {
    state.xs[i] = r * state.xs[i] * (1 - state.xs[i]);
    state.hist[i].push(state.xs[i]);
  }
  state.k++;
}
const ac = getComputedStyle(document.documentElement).getPropertyValue('--accent').trim() || '#5b9cff';
for (let i = 0; i < state.hist.length; i++) {
  ctx.strokeStyle = 'hsl(' + (210 + i * 10) + ',70%,64%)';
  ctx.globalAlpha = 0.8; ctx.lineWidth = 1.2;
  ctx.beginPath();
  for (let j = 0; j < state.hist[i].length; j++) {
    const px = j, py = H - state.hist[i][j] * H;
    j ? ctx.lineTo(px, py) : ctx.moveTo(px, py);
  }
  ctx.stroke();
}
ctx.globalAlpha = 1;
{{< /sketch >}}

## See also

- [[Sensitive Dependence on Initial Conditions]]
- [[Chaos]]
- [[Logistic Map]]
