---
title: Coupled Oscillators
tags: [oscillations]
summary: Connect two oscillators and they stop acting alone — exchanging energy through beats, splitting into normal modes, and ultimately synchronizing.
weight: 27
---

# Coupled Oscillators

Couple two oscillators — link two pendulums with a spring, two pacemaker cells through ion channels, two clocks bolted to the same beam — and they cease to be independent. Energy flows back and forth across the coupling, and the pair develops collective behaviors that neither half has alone: **normal modes**, **beating**, and **synchronization**.

## Normal modes

For two identical masses joined by springs, the motion looks complicated, but it is a superposition of two simple patterns that each oscillate at a single clean frequency:

- The **in-phase mode** — both masses swing together. The coupling spring never stretches, so this mode oscillates at the bare natural frequency $\omega_0$.
- The **anti-phase mode** — the masses swing oppositely. The coupling spring works hard, stiffening the restoring force and raising the frequency above $\omega_0$.

Any motion whatsoever is a mix of these two **normal modes**. They are the [[Phase Portrait|geometry]]'s preferred axes — the coordinates in which the coupled system falls apart into two independent [[Simple Harmonic Oscillator|simple oscillators]].

## Beats: energy sloshing back and forth

Start one pendulum swinging and leave its partner at rest. Because that initial state is an *equal mix* of the two normal modes — which run at slightly different frequencies — the modes drift in and out of step. The result is **beating**: the energy migrates entirely from one pendulum to the other and back, a slow throb at the difference frequency riding on the fast oscillation.

{{< sketch height="320" caption="A beating waveform: two nearby frequencies superposed. The fast oscillation is enveloped by a slow throb as the components fall in and out of phase — exactly how energy sloshes between two coupled pendulums. Drag horizontally to detune them." >}}
if (frame === 0) { state.detune = 0.12; }
if (mouse.down) {
  state.detune = 0.02 + (mouse.x / W) * 0.30;
}
ctx.clearRect(0,0,W,H);
var f1 = 1.0, f2 = 1.0 + state.detune;
var mid = H/2;

// midline
ctx.strokeStyle = "rgba(255,255,255,0.12)";
ctx.beginPath(); ctx.moveTo(0,mid); ctx.lineTo(W,mid); ctx.stroke();

// scroll phase with time
var ph = t * 2.2;
function wave(x){
  var u = (x / W) * 24;   // spatial extent
  return Math.cos(f1*u + ph) + Math.cos(f2*u + ph);
}
// envelope = 2*cos(half difference)
function env(x){
  var u = (x / W) * 24;
  return 2*Math.abs(Math.cos((f2-f1)/2 * u));
}
var amp = H*0.22;

// draw envelope (dim)
ctx.strokeStyle = "rgba(120,200,255,0.35)";
ctx.lineWidth = 1;
ctx.beginPath();
for (var x=0; x<=W; x+=2){ var y = mid - env(x)*amp/2; if(x===0)ctx.moveTo(x,y);else ctx.lineTo(x,y);} ctx.stroke();
ctx.beginPath();
for (var x2=0; x2<=W; x2+=2){ var y2 = mid + env(x2)*amp/2; if(x2===0)ctx.moveTo(x2,y2);else ctx.lineTo(x2,y2);} ctx.stroke();

// draw the beating wave
ctx.strokeStyle = "rgba(255,170,70,0.95)";
ctx.lineWidth = 2;
ctx.beginPath();
for (var x3=0; x3<=W; x3+=1){ var y3 = mid - wave(x3)*amp/2; if(x3===0)ctx.moveTo(x3,y3);else ctx.lineTo(x3,y3);} ctx.stroke();

ctx.fillStyle = "rgba(255,255,255,0.6)";
ctx.font = "13px sans-serif";
ctx.fillText("beat / detune = " + state.detune.toFixed(2) + "  — drag to retune", 12, 20);
{{< /sketch >}}

The closer the two frequencies, the slower and deeper the beats — drag the sketch to detune them and watch the throb stretch out.

## Synchronization

Add even weak coupling between *self-sustaining* oscillators (each a [[Limit Cycle]]) and something remarkable happens: they can lock to a common rhythm. This is **synchronization**, and it is why Huygens' two pendulum clocks on a shared beam drifted into anti-phase, why fireflies flash in unison, and why heart-cells beat as one. Pushed to large populations, the same coupling — when it overpowers the spread of natural frequencies — triggers a sudden collective onset of order, a phase transition akin to a [[Bifurcation]]. The continuous-space cousin of these locking interactions drives pattern formation in [[Reaction–Diffusion]] systems.

{{< note kind="key" title="From two to a trillion" >}}
The leap from a pair of coupled pendulums to synchronized fireflies, pacemaker cells, and power-grid generators is mostly a leap in *number*. The same ingredient — oscillators nudging each other through a shared connection — scales from beats to collective rhythm.
{{< /note >}}

## See also

- [[Resonance]]
- [[Limit Cycle]]
- [[Reaction–Diffusion]]
