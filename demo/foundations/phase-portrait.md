---
title: Phase Portrait
aliases: [phase plane]
tags: [foundations]
summary: The qualitative picture of all trajectories in state space, revealing a system's behavior at a glance.
weight: 40
---

# Phase Portrait

A **phase portrait** is the complete picture of a system's trajectories drawn in its [[State Space]]. Instead of plotting one variable against time, you plot the state variables against *each other* and sketch the family of curves the flow produces. The result is a kind of weather map of the dynamics: from it you can read where motion comes to rest, where it circulates, and which way nearby trajectories drift — all without solving a single equation.

This qualitative, geometric reading of behavior is the heart of the modern theory of [[Dynamical System|dynamical systems]], an approach that goes back to Poincaré.

## The anatomy of a portrait

A phase portrait is organized around a few landmark features:

- **Fixed points** — places where the flow stalls, $f(x^*) = 0$. Each is classified by how trajectories approach or flee it: nodes, saddles, spirals, centers. See [[Fixed Point]].
- **Closed orbits** — loops that the motion repeats forever. An isolated, attracting loop is a [[Limit Cycle]].
- **Separatrices** — special trajectories that divide the plane into regions of different fate, typically emanating from a saddle.

Together these skeletal features determine the behavior of *every* trajectory, because ordinary trajectories must thread between them without crossing.

## Read the field

Here is the phase plane of a damped oscillator, $\dot x = y,\; \dot y = -x - 0.3y$. Notice the inward spiral: the friction term bleeds energy, so every trajectory winds down to the stable spiral fixed point at the origin.

{{< sim name="vectorfield" fx="y" fy="-x-0.3*y" caption="Phase portrait of a damped oscillator. The origin is a stable spiral; all trajectories spiral in." >}}

Compare it mentally with the frictionless case $\dot y = -x$, whose portrait is nested closed circles, and with $\dot y = -x + 0.3y$, an *unstable* spiral that throws trajectories outward. Same template, three very different fates — the difference is one coefficient, and the crossover between them is a [[Bifurcation]].

{{< note kind="tip" title="Nullclines: a shortcut to the skeleton" >}}
The curves where $\dot x = 0$ and where $\dot y = 0$ are the **nullclines**. The flow is purely vertical on one and purely horizontal on the other, and fixed points sit exactly where they intersect. Sketching nullclines first is the fastest way to rough out a phase portrait by hand.
{{< /note >}}

## Why it matters

A phase portrait is robust: small changes to the equations usually nudge the picture without reorganizing it. The exceptions — the parameter values where the portrait *does* reorganize — are precisely the bifurcations, and tracking them is how we understand the onset of oscillation, [[Chaos]], and pattern formation.

## See also

- [[Fixed Point]]
- [[Limit Cycle]]
- [[State Space]]
