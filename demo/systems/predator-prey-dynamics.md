---
title: Predator–Prey Dynamics
aliases: [lotka-volterra, lotka volterra, predator prey, predator-prey dynamics]
tags: [systems, emergence, oscillations]
summary: Two coupled species — foxes and rabbits — whose populations rise and fall in eternal pursuit, tracing closed orbits in the phase plane.
weight: 60
---

# Predator–Prey Dynamics

**Predator–prey dynamics** describe two species locked in a feedback loop: prey feed the predators, predators thin the prey, and the thinned prey let the predators starve — which lets the prey rebound. The result is not a march to equilibrium but a perpetual **oscillation**, populations chasing each other up and down forever. The simplest model, the **Lotka–Volterra equations**, turns this story into two coupled differential equations and is a cornerstone example of a many-body [[Dynamical System]] with emergent rhythm.

## The Lotka–Volterra equations

Let $x$ be the prey population and $y$ the predators. Then

$$\dot x = \alpha x - \beta x y, \qquad \dot y = -\gamma y + \delta x y.$$

Read each term as a story. Prey grow on their own at rate $\alpha$ (plenty of grass) but are eaten at a rate proportional to *encounters*, $\beta x y$. Predators die off at rate $\gamma$ without food, but every encounter $\delta x y$ feeds their growth. The nonlinear $xy$ coupling — the rate at which the two species meet — is the entire engine of the dynamics.

There is a single coexistence [[Fixed Point]] where both rates balance,

$$x^* = \frac{\gamma}{\delta}, \qquad y^* = \frac{\alpha}{\beta},$$

but the populations almost never sit there. Instead they circle it.

## Closed orbits in the phase plane

Plot predators against prey and each cycle traces a **closed loop**: the system returns exactly to where it began and repeats. The phase plane fills with a family of nested loops, one for each starting amplitude — a [[Phase Portrait]] of pure circulation around a *center*. The orbits are closed but not isolated, so they are neutral cycles rather than a true [[Limit Cycle]]: nudge the system and it simply rides a different loop instead of returning to the same one.

In the simulation below, the **left panel** is the phase plane (prey horizontal, predators vertical) and the **right panel** is the time series. Watch the predator peak chase the prey peak — predators always crest a quarter-cycle *after* their food does.

{{< sketch height="340" caption="Lotka–Volterra dynamics. Left: the closed orbit in the prey–predator phase plane. Right: oscillating populations, predators lagging prey." >}}
if (frame === 0 || !state.init) {
  state.init = true;
  state.a = 0.9; state.b = 0.5; state.c = 1.0; state.d = 0.4;
  state.xmax = 5.5; state.ymax = 4.0;
  state.x = 1.0; state.y = 1.0;       // start off the equilibrium
  state.trail = [];                    // phase-plane points
  state.ts = [];                       // time series {x,y}
  state.dt = 0.012;
}
const f = (x, y) => [state.a*x - state.b*x*y, -state.c*y + state.d*x*y];
// RK4 integration, a few steps per frame, to keep the orbit closed
const steps = 6, h = state.dt;
for (let s = 0; s < steps; s++) {
  let x = state.x, y = state.y;
  const k1 = f(x, y);
  const k2 = f(x + 0.5*h*k1[0], y + 0.5*h*k1[1]);
  const k3 = f(x + 0.5*h*k2[0], y + 0.5*h*k2[1]);
  const k4 = f(x + h*k3[0], y + h*k3[1]);
  state.x = x + (h/6)*(k1[0] + 2*k2[0] + 2*k3[0] + k4[0]);
  state.y = y + (h/6)*(k1[1] + 2*k2[1] + 2*k3[1] + k4[1]);
}
state.trail.push([state.x, state.y]);
if (state.trail.length > 1400) state.trail.shift();
state.ts.push([state.x, state.y]);
if (state.ts.length > 480) state.ts.shift();

// background
ctx.fillStyle = "#0f1020";
ctx.fillRect(0, 0, W, H);
const pad = 34;
const split = W * 0.46;

// ---- phase plane (left) ----
const px = v => pad + (v / state.xmax) * (split - 2*pad);
const py = v => H - pad - (v / state.ymax) * (H - 2*pad);
ctx.strokeStyle = "#33365a"; ctx.lineWidth = 1;
ctx.strokeRect(pad, pad, split - 2*pad, H - 2*pad);
// equilibrium
const ex = state.c/state.d, ey = state.a/state.b;
ctx.fillStyle = "#8a8fb5";
ctx.beginPath(); ctx.arc(px(ex), py(ey), 3, 0, 7); ctx.fill();
// orbit trail
ctx.strokeStyle = "#5ad6c0"; ctx.lineWidth = 1.5;
ctx.beginPath();
for (let i = 0; i < state.trail.length; i++) {
  const [tx, ty] = state.trail[i];
  if (i === 0) ctx.moveTo(px(tx), py(ty)); else ctx.lineTo(px(tx), py(ty));
}
ctx.stroke();
// current point
ctx.fillStyle = "#ffd166";
ctx.beginPath(); ctx.arc(px(state.x), py(state.y), 4, 0, 7); ctx.fill();
ctx.fillStyle = "#9aa0c8"; ctx.font = "11px sans-serif";
ctx.fillText("prey →", split/2 - 14, H - 12);

// ---- time series (right) ----
const tx0 = split + pad, tw = W - split - pad - 8;
const ty = v => H - pad - (v / Math.max(state.xmax, state.ymax)) * (H - 2*pad);
ctx.strokeStyle = "#33365a";
ctx.strokeRect(tx0, pad, tw, H - 2*pad);
const drawSeries = (idx, color) => {
  ctx.strokeStyle = color; ctx.lineWidth = 1.5;
  ctx.beginPath();
  const n = state.ts.length;
  for (let i = 0; i < n; i++) {
    const sx = tx0 + (i / 480) * tw;
    if (i === 0) ctx.moveTo(sx, ty(state.ts[i][idx]));
    else ctx.lineTo(sx, ty(state.ts[i][idx]));
  }
  ctx.stroke();
};
drawSeries(0, "#5ad6c0");  // prey
drawSeries(1, "#ff6b8a");  // predators
ctx.fillStyle = "#5ad6c0"; ctx.fillText("prey", tx0 + 6, pad + 14);
ctx.fillStyle = "#ff6b8a"; ctx.fillText("predators", tx0 + 6, pad + 28);
{{< /sketch >}}

## Beyond the textbook model

The bare Lotka–Volterra model is famously fragile — its neutral cycles are an artifact of its simplicity. Add realism (prey that saturate their environment, predators that get full) and the closed loops typically collapse onto a single isolated [[Limit Cycle]], a self-correcting oscillation the system returns to after any disturbance. Push the parameters further and the cycle can lose stability in a [[Bifurcation]], the gateway to richer behavior. Despite its flaws, the model remains the archetype for coupled oscillation in ecology, epidemiology, and chemistry.

## See also

- [[Limit Cycle]]
- [[Phase Portrait]]
- [[Coupled Oscillators]]
