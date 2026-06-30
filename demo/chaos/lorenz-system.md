---
title: Lorenz System
aliases: [lorenz attractor, lorenz equations]
tags: [chaos]
summary: Three coupled equations, distilled from a weather model, whose solutions trace the iconic butterfly-shaped strange attractor.
weight: 30
---

# Lorenz System

The **Lorenz system** is three simple-looking differential equations that launched modern chaos theory. In 1963 the meteorologist Edward Lorenz stripped a model of atmospheric convection down to three variables and found something nobody expected: a deterministic flow that never repeats, never settles, and never escapes — instead winding forever around a [[Strange Attractor]] shaped like a butterfly's wings.

## The equations

$$\dot x = \sigma(y - x), \qquad \dot y = x(\rho - z) - y, \qquad \dot z = xy - \beta z.$$

Here $x$, $y$, and $z$ are abstract measures of convection intensity and temperature variation. The three constants are the knobs:

- $\sigma$ — the **Prandtl number** (fluid viscosity vs. thermal diffusivity), classically $\sigma = 10$.
- $\rho$ — the **Rayleigh number** (the strength of the driving temperature difference), classically $\rho = 28$.
- $\beta$ — a geometric aspect ratio, classically $\beta = 8/3 \approx 2.667$.

At the classic values the system is chaotic. Lorenz discovered this almost by accident: he restarted a run from numbers rounded to three decimals and the new trajectory diverged completely from the old — the first clear sighting of [[Sensitive Dependence on Initial Conditions|sensitive dependence]].

{{< note kind="note" title="Why three is the magic number" >}}
A continuous [[Flows and Maps|flow]] needs **at least three dimensions** to be chaotic. In two dimensions trajectories cannot cross (determinism) yet must stay bounded, so the Poincaré–Bendixson theorem traps them onto [[Fixed Point|fixed points]] or [[Limit Cycle|limit cycles]] — no room for chaos. The third dimension is exactly what lets a trajectory weave past itself without ever intersecting. The Lorenz system is the minimal, most famous example.
{{< /note >}}

## Flying the attractor

Drive the simulation below. Notice that the trajectory orbits one wing for a while, then flips to the other — but *when* it flips is effectively unpredictable. The set of all these orbits is the [[Lorenz System|Lorenz attractor]], a [[Strange Attractor]] of [[Fractal|fractal]] dimension $\approx 2.06$: thinner than a volume, thicker than a surface.

{{< sim name="lorenz" sigma="10" rho="28" beta="2.667" height="460" caption="The Lorenz attractor at the classic parameters σ=10, ρ=28, β=8/3. Two wings, infinitely many orbits, never a repeat." >}}

## Structure inside the chaos

The Lorenz flow has three [[Fixed Point|fixed points]]: the origin, plus a symmetric pair $C^\pm$ at the centers of the two wings. At $\rho = 28$ all three are unstable, so trajectories are perpetually repelled — yet bounded, so they can never leave. The resolution is to circle the unstable points forever, switching wings, tracing the butterfly.

The volume of any blob of initial conditions *shrinks* at a constant rate, since

$$\nabla\!\cdot\mathbf{f} = \frac{\partial \dot x}{\partial x}+\frac{\partial \dot y}{\partial y}+\frac{\partial \dot z}{\partial z} = -(\sigma + 1 + \beta) < 0.$$

A shrinking volume that nonetheless stretches in one direction is the engine of chaos: contract overall, stretch along one axis, fold the result back on itself. That stretch-and-fold is what gives the [[Attractor]] its layered, [[Fractal|fractal]] cross-section.

{{< quiz question="The Lorenz flow contracts phase-space volume everywhere, yet two nearby trajectories still diverge. How is that possible?" options="Volume contraction is an illusion of the numerics|It stretches along one direction while contracting more strongly in the others|The trajectories aren't really nearby|It only contracts near the fixed points" answer="2" explain="Net volume shrinks because contraction in two directions outweighs the stretching in the third. That single stretching direction is the positive Lyapunov exponent responsible for sensitive dependence." >}}

## See also

- [[Strange Attractor]]
- [[Attractor]]
- [[Sensitive Dependence on Initial Conditions]]


