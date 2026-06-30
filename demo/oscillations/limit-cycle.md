---
title: Limit Cycle
tags: [oscillations]
summary: An isolated closed orbit that a nonlinear system generates and maintains entirely on its own — the mathematics of self-sustained rhythm.
weight: 26
---

# Limit Cycle

A **limit cycle** is an isolated closed loop in [[State Space|state space]] that nearby trajectories spiral *onto*. Unlike the [[Simple Harmonic Oscillator]]'s family of nested orbits — where the amplitude is set by how hard you started it — a limit cycle has *one* preferred amplitude that the system returns to no matter where it begins. It is an [[Attractor]], but a one-dimensional, oscillating one.

This is the mathematics of **self-sustained oscillation**: a heartbeat, a firing neuron, a chirping circuit, a beating laser. None of these needs a periodic [[Driven Oscillator|external drive]]; the rhythm is intrinsic, born from a balance between an internal energy source and dissipation.

## The van der Pol oscillator

The classic example is the **van der Pol oscillator**, originally a model of a vacuum-tube circuit:

$$\ddot x - \mu(1-x^2)\dot x + x = 0.$$

The trick is in the damping term $-\mu(1-x^2)\dot x$. When the amplitude is *small* ($x^2 < 1$) the term is *negative* damping — it pumps energy in and the oscillation grows. When the amplitude is *large* ($x^2 > 1$) it becomes ordinary positive damping, bleeding energy away. The system is squeezed from both sides onto a single stable loop.

{{< sketch height="380" caption="A trajectory (orange) spiraling outward from near the center and inward from far away — both converging onto the same stable limit cycle (cyan). Move your mouse to relaunch from a new starting point." >}}
if (frame === 0) {
  state.pts = [];
  state.x = 0.05; state.y = 0.0;
  state.mu = 1.2;
}
// relaunch from mouse position when pressed or moved into canvas
if (mouse.down) {
  state.x = (mouse.x - W/2) / (W*0.16);
  state.y = (mouse.y - H/2) / (H*0.16);
  state.pts = [];
}
// integrate van der Pol with small steps
var steps = 6;
for (var i = 0; i < steps; i++) {
  var h = 0.02;
  var dx = state.y;
  var dy = state.mu*(1 - state.x*state.x)*state.y - state.x;
  state.x += dx*h;
  state.y += dy*h;
}
state.pts.push({x: state.x, y: state.y});
if (state.pts.length > 900) state.pts.shift();

// map state coords to screen
function sx(x){ return W/2 + x*W*0.16; }
function sy(y){ return H/2 + y*H*0.16; }

ctx.clearRect(0,0,W,H);
// faint axes
ctx.strokeStyle = "rgba(255,255,255,0.12)";
ctx.beginPath(); ctx.moveTo(0,H/2); ctx.lineTo(W,H/2);
ctx.moveTo(W/2,0); ctx.lineTo(W/2,H); ctx.stroke();

// draw the attracting limit cycle by integrating a long-settled orbit once
if (!state.cycle) {
  state.cycle = [];
  var cx = 2.0, cy = 0.0;
  for (var k = 0; k < 4000; k++) {
    var ddx = cy;
    var ddy = state.mu*(1 - cx*cx)*cy - cx;
    cx += ddx*0.01; cy += ddy*0.01;
    if (k > 2000) state.cycle.push({x: cx, y: cy});
  }
}
ctx.strokeStyle = "rgba(80,220,235,0.9)";
ctx.lineWidth = 2;
ctx.beginPath();
for (var j = 0; j < state.cycle.length; j++) {
  var p = state.cycle[j];
  if (j===0) ctx.moveTo(sx(p.x), sy(p.y)); else ctx.lineTo(sx(p.x), sy(p.y));
}
ctx.stroke();

// draw the live trajectory
ctx.strokeStyle = "rgba(255,160,60,0.85)";
ctx.lineWidth = 1.5;
ctx.beginPath();
for (var m = 0; m < state.pts.length; m++) {
  var q = state.pts[m];
  if (m===0) ctx.moveTo(sx(q.x), sy(q.y)); else ctx.lineTo(sx(q.x), sy(q.y));
}
ctx.stroke();

// moving head
var head = state.pts[state.pts.length-1];
if (head) {
  ctx.fillStyle = "#ffd060";
  ctx.beginPath(); ctx.arc(sx(head.x), sy(head.y), 4, 0, Math.PI*2); ctx.fill();
}
ctx.fillStyle = "rgba(255,255,255,0.6)";
ctx.font = "13px sans-serif";
ctx.fillText("van der Pol  μ=" + state.mu.toFixed(1) + "  — click/drag to relaunch", 12, 20);
{{< /sketch >}}

## Why limit cycles need nonlinearity

A *linear* system can only spiral in (damped) or out (unstable) — it can never settle onto a fixed-amplitude loop. Limit cycles are inescapably nonlinear: it takes a term like $\mu(1-x^2)$ that *changes sign with amplitude* to trap the motion. They typically appear when a [[Fixed Point]] loses [[Stability]] through a **Hopf** [[Bifurcation]], spinning a stable equilibrium off into a growing oscillation as a parameter crosses a threshold.

## See also

- [[Phase Portrait]]
- [[Bifurcation]]
- [[Attractor]]
