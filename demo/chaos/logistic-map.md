---
title: Logistic Map
aliases: [logistic equation, logistic map]
tags: [chaos]
summary: A one-line population model, x_{n+1}=r x_n(1-x_n), whose period-doubling cascade is the canonical road from order into chaos.
weight: 50
---

# Logistic Map

The **logistic map** is the cleanest proof that chaos needs no complexity. It is a single line of arithmetic, iterated:

$$x_{n+1} = r\,x_n(1 - x_n).$$

Think of $x_n$ as a population, scaled to lie between 0 (extinct) and 1 (carrying capacity), measured generation by generation. The factor $r x_n$ is growth; the factor $(1-x_n)$ is the brake of overcrowding. One parameter, $r$, controls the whole story — and as you raise it from 0 to 4, this innocent [[Flows and Maps|map]] marches through fixed points, oscillations, a cascade of [[Bifurcation|bifurcations]], and finally full-blown [[Chaos|chaos]].

## The cast of behaviors

For small $r$ the population just dies out. Past $r = 1$ it settles to a steady [[Fixed Point]]. At $r = 3$ that fixed point loses [[Stability|stability]] and the population starts alternating between two values — a period-2 cycle. Raise $r$ further and the period doubles to 4, then 8, then 16, the windows of stability shrinking geometrically, until near $r \approx 3.5699$ the period becomes infinite and the dynamics turn chaotic.

{{< note kind="key" title="Feigenbaum's universal number" >}}
The successive period-doubling thresholds $r_n$ pile up geometrically, and the ratio of the gaps converges to a *universal* constant:

$$\delta = \lim_{n\to\infty}\frac{r_{n}-r_{n-1}}{r_{n+1}-r_{n}} = 4.6692016\ldots$$

Mitchell Feigenbaum found that the **same** number governs the period-doubling route in countless unrelated systems — dripping faucets, heart cells, electronic circuits. The logistic map is the simplest place to meet this deep universality; see [[Feigenbaum Constant]] for where these thresholds fall and how fast they pile up.
{{< /note >}}

## The bifurcation diagram

Here is the showpiece. For every value of $r$ across the width of the canvas, we iterate the map, throw away the transient, and plot the values it *settles onto* down the height. Read it left to right and you watch order dissolve into chaos: one branch splits into two, two into four, four into eight — the **period-doubling cascade** — then shatters into a fractal mist, threaded by surprising white **windows** of restored periodicity (the widest is the period-3 window near $r \approx 3.83$).

**Move your mouse across the diagram** to pick an $r$ and see the live orbit it produces.

{{< sketch height="460" caption="Bifurcation diagram of x→r x(1-x) for r from 2.5 to 4.0. Computed once into an offscreen buffer; hover to probe the attracting orbit at any r. Density (darker) shows where the iterates spend their time." >}}
if (frame === 0) {
  state.rmin = 2.5;
  state.rmax = 4.0;
  // Render the diagram once into an offscreen buffer, then just blit it.
  const buf = document.createElement('canvas');
  buf.width = Math.max(1, Math.floor(W));
  buf.height = Math.max(1, Math.floor(H));
  const b = buf.getContext('2d');
  b.fillStyle = '#0b0e16';
  b.fillRect(0, 0, buf.width, buf.height);
  const transient = 250;
  const samples = 220;
  for (let px = 0; px < buf.width; px++) {
    const r = state.rmin + (px / (buf.width - 1)) * (state.rmax - state.rmin);
    let x = 0.5;
    for (let i = 0; i < transient; i++) x = r * x * (1 - x);
    // color shifts blue->cyan->magenta with r, for depth
    const hue = 210 + 110 * ((r - state.rmin) / (state.rmax - state.rmin));
    b.fillStyle = 'hsla(' + hue + ', 85%, 65%, 0.16)';
    for (let i = 0; i < samples; i++) {
      x = r * x * (1 - x);
      const y = buf.height - 1 - x * (buf.height - 1);
      b.fillRect(px, y, 1, 1);
    }
  }
  state.buf = buf;
}

ctx.clearRect(0, 0, W, H);
ctx.drawImage(state.buf, 0, 0, W, H);

// axes labels
ctx.fillStyle = 'rgba(255,255,255,0.55)';
ctx.font = '11px sans-serif';
ctx.fillText('r = ' + state.rmin.toFixed(1), 4, H - 6);
ctx.fillText('r = ' + state.rmax.toFixed(1), W - 56, H - 6);
ctx.fillText('x = 1', 4, 14);
ctx.fillText('x = 0', 4, H - 18);

// interactive probe
if (mouse.x >= 0 && mouse.x <= W && mouse.y >= 0 && mouse.y <= H) {
  const r = state.rmin + (mouse.x / W) * (state.rmax - state.rmin);
  // vertical cursor
  ctx.strokeStyle = 'rgba(255,255,255,0.5)';
  ctx.lineWidth = 1;
  ctx.beginPath(); ctx.moveTo(mouse.x, 0); ctx.lineTo(mouse.x, H); ctx.stroke();
  // iterate live and mark the settled orbit
  let x = 0.5;
  for (let i = 0; i < 300; i++) x = r * x * (1 - x);
  ctx.fillStyle = '#ffd166';
  for (let i = 0; i < 200; i++) {
    x = r * x * (1 - x);
    const y = H - 1 - x * (H - 1);
    ctx.beginPath(); ctx.arc(mouse.x, y, 1.6, 0, 2 * Math.PI); ctx.fill();
  }
  ctx.fillStyle = 'rgba(255,255,255,0.9)';
  ctx.font = '13px sans-serif';
  ctx.fillText('r = ' + r.toFixed(4), Math.min(mouse.x + 8, W - 90), 20);
}
{{< /sketch >}}

## Reading the chaos quantitatively

Whether a given $r$ is periodic or chaotic is decided by the [[Lyapunov Exponent]] of the map,

$$\lambda(r) = \lim_{N\to\infty}\frac{1}{N}\sum_{n=0}^{N-1}\ln\big|\,r(1 - 2x_n)\,\big|.$$

Where $\lambda < 0$ neighboring states converge and the orbit is a stable cycle; where $\lambda > 0$ they diverge and the orbit is chaotic. The plot below sketches this growth-rate factor $f'(x)=r(1-2x)$ — sensitivity per step lives in its magnitude.

{{< plot fn="2.5*(1-2*x);;4*(1-2*x)" xmin="0" xmax="1" title="Per-step stretch factor f'(x)=r(1-2x) at r=2.5 (stable) and r=4 (chaotic)" caption="Where |f'| exceeds 1, a single step stretches small errors — the seed of a positive Lyapunov exponent." >}}

{{< quiz question="In the bifurcation diagram, what is happening at the values of r where a single curve splits into two?" options="The population goes extinct|A period-doubling bifurcation: the stable cycle's period doubles|Numerical round-off error|The map stops being deterministic" answer="2" explain="Each fork is a period-doubling bifurcation. The previously stable cycle loses stability and is replaced by one of twice the period; the accumulation of these doublings is the route to chaos." >}}

## See also

- [[Bifurcation]]
- [[Flows and Maps]]
- [[Lyapunov Exponent]]
