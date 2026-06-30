---
title: Dynamical System
aliases: [dynamical systems theory]
tags: [foundations]
summary: A state together with a fixed rule that determines how that state evolves in time.
weight: 20
---

# Dynamical System

A **dynamical system** is a pair of ingredients: a **state** that captures everything you need to know about the system right now, and an **evolution rule** that says how the state changes a moment later. Nothing else is required. Give me the state today and the rule, and the entire future (and often the past) is determined. This is the engine behind orbits, circuits, populations, and pendulums alike.

Formally, the state lives in a [[State Space]] $X$, and the rule generates a trajectory $x(t)$ through that space. The two great families of rules — continuous and discrete — are studied side by side in [[Flows and Maps]].

## Continuous vs. discrete time

In **continuous time**, the rule is a differential equation — a *flow*. The state moves smoothly, and the rule specifies its velocity at every point:

$$\dot x = f(x), \qquad x \in \mathbb{R}^n.$$

A planet's position and momentum, a capacitor's voltage, a chemical concentration — all evolve this way.

In **discrete time**, the rule is a *map* applied over and over, advancing the state in steps:

$$x_{n+1} = F(x_n).$$

Population counts measured generation by generation, or the iterates of the [[Logistic Map]], live here. A flow can always be turned into a map by sampling it at regular intervals or on a [[Poincaré Section]] — a bridge we lean on constantly.

{{< note kind="key" title="The defining property: determinism" >}}
Same state, same rule $\Rightarrow$ same future. A dynamical system carries no memory beyond its current state and no dice beyond its rule. Randomness, when it appears, is layered on top — the bare system is deterministic.
{{< /note >}}

## A first example

Take the linear system $\dot x = -x$. The rule says "always move toward the origin at a rate proportional to your distance." Its solutions decay exponentially, $x(t) = x_0 e^{-t}$, sliding into the [[Fixed Point]] at the origin no matter where they start.

{{< plot fn="Math.exp(-x)" xmin="0" xmax="5" ymin="0" ymax="1.1" title="Trajectory of x' = -x from x0 = 1" caption="Every initial condition relaxes to the equilibrium at 0." >}}

Change one sign to get $\dot x = +x$ and the same equilibrium becomes a launchpad: trajectories flee to infinity. The rule barely changed, yet the behavior inverted — a foretaste of how [[Stability]] and [[Bifurcation]] hinge on small details.

## What we want to know

We rarely solve these systems exactly. Instead we ask qualitative questions that a picture can answer:

- Where are the equilibria, the [[Fixed Point|fixed points]] where $f(x)=0$?
- Are they stable? (See [[Stability]].)
- What does the long-term motion look like — does it land on an [[Attractor]]?
- How sensitive is the motion to its starting point? (See [[Lyapunov Exponent]].)

The right setting for all of these is the geometry of the [[State Space]] and its [[Phase Portrait]].

## See also

- [[State Space]]
- [[Flows and Maps]]
- [[Attractor]]
