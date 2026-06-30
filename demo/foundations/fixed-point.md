---
title: Fixed Point
aliases: [equilibrium, fixed points]
tags: [foundations]
summary: A state where the dynamics stand still — an equilibrium — classified by how nearby trajectories behave around it.
weight: 50
---

# Fixed Point

A **fixed point** (or **equilibrium**) is a state where the system, left alone, stays put forever. For a flow it is a point $x^*$ where the velocity vanishes,

$$f(x^*) = 0,$$

so there is nowhere to go. For a map it is a point that returns to itself, $F(x^*) = x^*$. Fixed points are the anchors of a [[Phase Portrait]]: every other trajectory is organized by how it approaches, avoids, or orbits them.

Finding them is the easy part — just solve $f(x)=0$. The interesting question is what happens *near* one.

## Classifying a fixed point

Zoom in close enough and a smooth flow looks linear, governed by the **Jacobian matrix** $A = Df(x^*)$. The eigenvalues of $A$ tell you the local geography:

- **Node** — eigenvalues real, same sign. Trajectories dive straight in (stable, both negative) or straight out (unstable, both positive).
- **Saddle** — real eigenvalues of opposite sign. Trajectories approach along one direction and flee along another; saddles are always unstable.
- **Spiral (focus)** — complex eigenvalues. Trajectories wind in or out, rotating as they go — the hallmark of decaying or growing oscillation.
- **Center** — purely imaginary eigenvalues. Trajectories form closed loops, neither approaching nor receding, as in a frictionless oscillator.

Whether each type attracts or repels is the subject of [[Stability]], which reads off the *sign* of the real parts of those eigenvalues.

{{< note kind="key" title="Linearization in one line" >}}
Let $u = x - x^*$ be a small displacement. Then $\dot u \approx A u$, so near a fixed point the messy nonlinear flow is replaced by a linear one whose behavior is dictated entirely by the eigenvalues of $A$.
{{< /note >}}

## See the types side by side

The vector field below is the damped oscillator $\dot x = y,\; \dot y = -x - 0.3y$, whose only fixed point — the origin — is a **stable spiral**. Trajectories rotate inward toward it.

{{< sim name="vectorfield" fx="y" fy="-x-0.3*y" caption="A stable spiral at the origin. Crank the friction past 2 and the spiral straightens into a node." >}}

Swap the rule to $\dot x = y,\; \dot y = x$ and the origin becomes a **saddle**: most trajectories sweep in, bend, and shoot back out along the diagonal escape directions.

## One point, many roles

The *same* equilibrium can change character as a parameter varies — a stable spiral can lose its stability and spawn a [[Limit Cycle]], or two fixed points can collide and annihilate. Those transitions are exactly what [[Bifurcation]] catalogs.

## See also

- [[Stability]]
- [[Phase Portrait]]
- [[Bifurcation]]
