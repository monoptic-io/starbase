---
title: N-Body Problem
aliases: [n-body, n-body problem, gravitational n-body]
tags: [systems, emergence, gravitation]
summary: Predicting the motion of many masses pulling on one another through gravity — the original many-body problem.
weight: 20
---

# N-Body Problem

The **N-body problem** asks a question as old as astronomy: given $N$ masses, each pulling on every other through gravity, how do they move? Each body feels the summed attraction of all the rest, and in turn changes how it pulls on them. The state is the set of all positions and velocities; the rule is Newton's law of gravitation applied to every pair at once. It is the founding example of a many-body [[Dynamical System]], and the place where the dream of perfect prediction first ran into trouble.

The force on body $i$ is the vector sum over all the others:

$$\ddot{\mathbf r}_i = G \sum_{j \neq i} m_j \frac{\mathbf r_j - \mathbf r_i}{\lVert \mathbf r_j - \mathbf r_i \rVert^3}.$$

With $N$ bodies that is a coupled system of $N$ second-order vector equations — every body in every other body's equation. The coupling is what makes it hard, and what makes it interesting.

## Two bodies are easy, three are not

For $N = 2$ the problem is completely solvable: Kepler's ellipses, parabolas, and hyperbolas, written down in closed form. The two-body motion is *integrable* — there are exactly as many conserved quantities (energy, momentum, angular momentum) as degrees of freedom, and they pin the motion onto simple curves.

Add a single body and that good fortune evaporates. The [[Three-Body Problem]] has no general closed-form solution, and for most initial conditions the motion is chaotic. The jump from two to three is one of the great surprises in the history of physics: difficulty does not grow gently with $N$, it arrives all at once.

{{< sim name="nbody" bodies="6" caption="Six gravitating masses. Watch how close encounters fling bodies onto new paths — the hallmark of strong coupling." >}}

## What survives the chaos

Even when individual orbits are unpredictable, the whole system obeys strict bookkeeping. Total momentum, total angular momentum, and total energy are conserved, and the center of mass drifts in a straight line forever. These conservation laws constrain the motion onto a lower-dimensional surface in [[State Space]] without ever taming it completely.

{{< note kind="note" title="Why simulations, not formulas" >}}
Because closed-form solutions are unavailable for $N \geq 3$, practical celestial mechanics is *numerical*: integrate the equations forward in tiny time steps. The cost of the naive force sum grows like $N^2$ per step — every pair must be considered — which is why galaxy-scale simulations lean on clever tree and multipole approximations.
{{< /note >}}

## See also

- [[Three-Body Problem]]
- [[Dynamical System]]
- [[Chaos]]
