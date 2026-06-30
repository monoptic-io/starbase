---
title: Poincaré Section
aliases: [poincare map, poincaré map, poincare section, surface of section]
tags: [chaos]
summary: A slice through a continuous flow that turns it into a discrete map, exposing the hidden structure of complicated trajectories.
weight: 80
---

# Poincaré Section

A **Poincaré section** is a brilliantly simple trick for taming a complicated [[Flows and Maps|flow]]. Instead of following a trajectory's every wiggle through its full [[State Space]], you stretch a surface across the flow and record only the points where the trajectory *punctures* it. A continuous, hard-to-read curve in three or more dimensions becomes a tidy scatter of dots in two — a [[Flows and Maps|map]] you can actually analyze.

Henri Poincaré invented the idea while wrestling with the [[Three-Body Problem]]. It remains the single most useful instrument for diagnosing [[Chaos|chaos]] in continuous systems.

## From a flow to a map

Pick a surface $\Sigma$ that trajectories cross transversally. Each time the trajectory passes through $\Sigma$ in a chosen direction, mark the crossing point. The rule "from one crossing, find the next" defines the **Poincaré map** (or *return map*)

$$P : \Sigma \to \Sigma, \qquad x_{n+1} = P(x_n),$$

which advances the system one full loop at a time. The map inherits the dynamics but drops a dimension — and dropping a dimension is exactly what makes the structure legible.

{{< sketch height="400" caption="A trajectory spiraling through 3-D state space, punctured by a section plane. Each crossing drops a dot onto the plane — the Poincaré map. Drag to rotate the view." >}}
if (frame === 0) {
  state.yaw = -0.6;
  state.dragX = 0;
  state.pts = [];      // recorded crossings (in plane coords u,v)
  state.theta = 0;
  state.lastSign = 0;
}
// rotate view by dragging
if (mouse.down) { state.yaw = -0.6 + (mouse.x / W - 0.5) * 3.0; }

const cx = W * 0.5, cy = H * 0.5, scale = Math.min(W, H) * 0.30;
const cosY = Math.cos(state.yaw), sinY = Math.sin(state.yaw);
// project a 3-D point (camera tilted a bit on x)
function project(x, y, z) {
  // rotate about vertical (y) axis
  const xr = x * cosY - z * sinY;
  const zr = x * sinY + z * cosY;
  // slight tilt down
  const tilt = 0.5;
  const yr = y * Math.cos(tilt) - zr * Math.sin(tilt);
  return { X: cx + xr * scale, Y: cy - yr * scale };
}

ctx.clearRect(0, 0, W, H);

// --- draw the section plane (the plane z = 0), as a grid quad ---
ctx.strokeStyle = 'rgba(120,170,255,0.35)';
ctx.lineWidth = 1;
const lim = 1.3;
for (let i = -3; i <= 3; i++) {
  const f = (i / 3) * lim;
  let a = project(f, -lim, 0), b = project(f, lim, 0);
  ctx.beginPath(); ctx.moveTo(a.X, a.Y); ctx.lineTo(b.X, b.Y); ctx.stroke();
  a = project(-lim, f, 0); b = project(lim, f, 0);
  ctx.beginPath(); ctx.moveTo(a.X, a.Y); ctx.lineTo(b.X, b.Y); ctx.stroke();
}

// --- advance the trajectory: a spiral on a slowly precessing torus ---
// param: x = (R + r cos a) cos th ... we simplify to a 3-D Lissajous-like spiral
function traj(th) {
  const R = 0.9;
  const a = th * 5.0;            // fast winding
  const r = 0.45;
  const x = (R + r * Math.cos(a)) * Math.cos(th);
  const y = (R + r * Math.cos(a)) * Math.sin(th);
  const z = r * Math.sin(a);     // oscillates through the plane z=0
  return { x, y, z };
}

// draw a length of trajectory as a fading tail
const steps = 240, span = 1.4;
let prev = null;
for (let i = 0; i <= steps; i++) {
  const th = state.theta - span + (i / steps) * span;
  const p = traj(th);
  const pr = project(p.x, p.y, p.z);
  if (prev) {
    const alpha = 0.15 + 0.55 * (i / steps);
    ctx.strokeStyle = 'rgba(255,209,102,' + alpha + ')';
    ctx.lineWidth = 1.8;
    ctx.beginPath(); ctx.moveTo(prev.X, prev.Y); ctx.lineTo(pr.X, pr.Y); ctx.stroke();
  }
  prev = pr;
}

// head of trajectory
const head = traj(state.theta);
const hp = project(head.x, head.y, head.z);
ctx.fillStyle = '#ff6b6b';
ctx.beginPath(); ctx.arc(hp.X, hp.Y, 4, 0, 2 * Math.PI); ctx.fill();

// detect upward crossings of the plane z = 0 -> record (x,y) as plane coords
const sign = head.z >= 0 ? 1 : -1;
if (state.lastSign < 0 && sign >= 0) {
  state.pts.push({ u: head.x, v: head.y });
  if (state.pts.length > 60) state.pts.shift();
}
state.lastSign = sign;

// draw recorded crossing dots lying in the plane
for (const q of state.pts) {
  const pr = project(q.u, q.v, 0);
  ctx.fillStyle = '#5b9cff';
  ctx.beginPath(); ctx.arc(pr.X, pr.Y, 3, 0, 2 * Math.PI); ctx.fill();
}

ctx.fillStyle = 'rgba(255,255,255,0.7)';
ctx.font = '12px sans-serif';
ctx.fillText('blue dots = Poincaré map (crossings of the plane)', 10, H - 12);

state.theta += dt * 1.2;
{{< /sketch >}}

## Reading chaos off the section

The payoff is diagnostic. The *pattern* of dots tells you the character of the motion at a glance:

- **A single dot** — the trajectory closes after one loop: a periodic orbit, a [[Limit Cycle]].
- **A handful of dots** — a period-$n$ cycle, the kind produced by a [[Bifurcation|period-doubling]] cascade.
- **A smooth closed curve** — quasi-periodic motion on a torus (two incommensurate frequencies).
- **A fuzzy fractal scatter** — chaos. The dots fill out the cross-section of a [[Strange Attractor]].

This is how a tumbling [[Double Pendulum]] is analyzed: take its four-dimensional flow, fix the energy, slice with a section, and the resulting dots reveal islands of order floating in a sea of chaotic spray.

{{< note kind="tip" title="A section also speeds up computation" >}}
Because the Poincaré map skips the smooth motion between crossings, it compresses long, expensive integrations into a few essential numbers. Many results about [[Stability]] of periodic orbits — and the [[Lyapunov Exponent]] of a flow — are most easily computed on the section rather than on the full flow.
{{< /note >}}

{{< quiz question="On a Poincaré section, what does a fuzzy, fractal-looking scatter of points indicate?" options="A numerical bug|Chaotic motion on a strange attractor|A stable fixed point|Perfectly periodic motion" answer="2" explain="A periodic orbit leaves finitely many dots and a quasi-periodic one a smooth curve; a fractal scatter that never closes up is the signature of chaos, tracing the attractor's cross-section." >}}

## See also

- [[Flows and Maps]]
- [[Double Pendulum]]
- [[Strange Attractor]]
