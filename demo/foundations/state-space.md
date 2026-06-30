---
title: State Space
aliases: [phase space]
tags: [foundations]
summary: The space of all possible states of a system, in which each point is a complete snapshot and motion traces a trajectory.
weight: 30
---

# State Space

The **state space** (or **phase space**) is the set of all states a system can occupy. Each *point* is a full snapshot — every variable you need to determine the future, gathered into a single coordinate. A pendulum's state, for instance, needs both its angle and its angular velocity, so its state space is two-dimensional; a single number would not be enough to say where it goes next.

Once you adopt this viewpoint, the evolution rule of a [[Dynamical System]] becomes a recipe for *moving points around* the state space. The path a point follows is a **trajectory** (or *orbit*), and the collection of all trajectories is the [[Phase Portrait]].

## Points, dimensions, and trajectories

For a flow $\dot x = f(x)$ on $\mathbb{R}^n$, the function $f$ assigns a velocity vector to every point — a **vector field**. A trajectory is a curve that is everywhere tangent to those arrows, threading the field like a streamline in a fluid.

Two facts give state space its rigid beauty:

- **Trajectories never cross.** Because the rule is deterministic, only one trajectory passes through each point. If two crossed, that intersection would have two futures.
- **Dimension counts.** A continuous flow needs at least three dimensions before it can behave chaotically; in one or two dimensions the no-crossing rule pens trajectories in. This is why the [[Lorenz System]] lives in 3D.

{{< note kind="note" title="Why velocity belongs in the state" >}}
For mechanical systems, position alone is not a state: the same position can be moving left or right. Pairing each position with its velocity (or momentum) makes the future unique. That doubling is why phase space is so often even-dimensional.
{{< /note >}}

## Watch the flow

Below is the vector field for the simple oscillator $\dot x = y,\; \dot y = -x$. Drop a point anywhere and it circles the origin forever — a closed trajectory that is the state-space signature of steady oscillation. Each little arrow is the velocity the rule assigns at that location.

{{< sim name="vectorfield" fx="y" fy="-x" caption="State space of x' = y, y' = -x. Trajectories are nested circles; the origin is a center." >}}

Add a touch of friction, $\dot y = -x - 0.3y$, and those circles become inward spirals — every trajectory winds down to the resting [[Fixed Point]] at the origin, the phase-space picture of a [[Damped Oscillator]].

## Reading behavior off the geometry

Because the geometry *is* the dynamics, qualitative questions become visual ones. Closed loops mean periodic motion (see [[Limit Cycle]]); points where all arrows vanish are equilibria; regions that all nearby trajectories drain into are [[Attractor|attractors]]. Learning to read these shapes is the whole content of the [[Phase Portrait]].

## See also

- [[Phase Portrait]]
- [[Dynamical System]]
- [[Attractor]]
