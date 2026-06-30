---
title: Attractor
aliases: [attractors, basin of attraction]
tags: [foundations]
summary: The set a system settles onto after transients die away — a point, a loop, or something stranger.
weight: 70
---

# Attractor

An **attractor** is a set of states that a [[Dynamical System]] tends toward as
time goes on, and *stays* near once it arrives. Release the system from many
different starting conditions and, after the initial transient fades, the
trajectories all end up doing the same thing — orbiting the same shape in
[[State Space]]. That shape is the attractor, and the collection of starting
points that flow to it is its **basin of attraction**.

## A hierarchy of shapes

Attractors come in a small zoo of types, in order of increasing richness:

- A **fixed point** — the system runs down to rest, like a [[Damped Oscillator]]
  coming to a stop. The attractor is a single point.
- A **limit cycle** — the system settles into a steady repeating oscillation.
  The attractor is a closed loop; see [[Limit Cycle]].
- A **torus** — two incommensurate frequencies combine, and the trajectory winds
  forever over the surface of a doughnut without closing.
- A **strange attractor** — the trajectory is bounded and never repeats, folded
  into a [[Fractal]]. This is the signature of [[Chaos]]; see
  [[Strange Attractor]].

{{< sim name="lorenz" caption="The Lorenz attractor: a bounded, never-repeating strange attractor. Every trajectory is sucked onto this butterfly, yet none ever crosses itself." >}}

## What makes a set an attractor

Three conditions: it is **invariant** (start on it, stay on it), it **attracts**
an open neighborhood of nearby states, and it is **minimal** — no smaller piece
of it does the same job. The attracting property is why attractors are what you
actually observe in nature: unstable behaviors are washed out, and the system is
left riding its attractor.

{{< note kind="key" title="Dissipation is the price of admission" >}}
Only **dissipative** systems — those that lose volume in state space, like
anything with friction — have attractors. A frictionless, energy-conserving
system (a perfect [[Pendulum]]) instead has nested orbits and no attractor: it
remembers its initial energy forever.
{{< /note >}}

## Basins and competition

A system can have several attractors at once, each with its own basin. The
boundaries between basins decide the system's fate from a given start, and those
boundaries can themselves be fractal — making the long-term outcome practically
unpredictable even when each attractor is simple. A change in a parameter can
make an attractor appear, vanish, or swap stability with another: a
[[Bifurcation]].

## See also

- [[Strange Attractor]]
- [[Limit Cycle]]
- [[Lorenz System]]
