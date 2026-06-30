---
title: Strange Attractor
aliases: [strange attractor, chaotic attractor]
tags: [chaos]
summary: An attractor with fractal structure, onto which chaotic motion settles while never repeating itself.
weight: 60
---

# Strange Attractor

A **strange attractor** is the geometric home of chaos. Like any [[Attractor]], it is a set that nearby trajectories are drawn toward and then stay on. What makes it *strange* is twofold: the motion on it is [[Chaos|chaotic]] — showing [[Sensitive Dependence on Initial Conditions|sensitive dependence]] — and the set itself is a [[Fractal]], with intricate self-similar structure and a non-integer dimension.

That combination sounds paradoxical. An attractor pulls trajectories *together*; chaos pushes nearby trajectories *apart*. A strange attractor does both at once, and the resolution is a perpetual act of **stretching and folding**.

## Stretch, fold, repeat

Picture a blob of initial conditions on the attractor. The flow stretches it along one direction — that stretching *is* the positive [[Lyapunov Exponent]], pulling neighbors apart and destroying predictability. But the attractor is bounded, so the stretched blob cannot grow forever; instead it gets folded back over itself, like a baker kneading dough. Stretch, fold, stretch, fold, endlessly.

{{< note kind="key" title="Why strange attractors must be fractal" >}}
Repeated stretching-and-folding lays down layer upon layer upon layer, never quite touching, at every scale. Zoom in on a cross-section and you find the same banded structure again — a Cantor-set-like infinity of leaves. That is exactly what gives a strange attractor its **fractal dimension**: it occupies more than a surface but less than a volume. The [[Lorenz System|Lorenz attractor]] has dimension $\approx 2.06$.
{{< /note >}}

## The canonical example

The [[Lorenz System]] traces the most famous strange attractor of all — the butterfly. Run it again here and watch the two properties coexist: every trajectory is sucked onto the same delicate winged surface (attraction), yet two trajectories started side by side soon orbit it out of step (chaos).

{{< sim name="lorenz" sigma="10" rho="28" beta="2.667" height="460" caption="The Lorenz attractor — the archetype of a strange attractor. Trajectories converge onto it, then diverge along it." >}}

The Lorenz butterfly is far from alone. The Rössler attractor folds a single twisted band; the Hénon map produces a strange attractor for a two-dimensional [[Flows and Maps|map]]; the [[Double Pendulum]] and driven oscillators trace their own. All share the stretch-and-fold signature.

## Strange versus ordinary attractors

It helps to line them up. A stable [[Fixed Point]] is a zero-dimensional attractor — motion stops. A [[Limit Cycle]] is a one-dimensional attractor — motion repeats forever. A strange attractor is the next rung: bounded motion that *never* repeats, riding a fractal set.

{{< quiz question="What distinguishes a strange attractor from a limit cycle?" options="A strange attractor is unstable|On a strange attractor motion never repeats and the set is fractal, whereas a limit cycle is a single repeating loop|A limit cycle is higher-dimensional|There is no real difference" answer="2" explain="Both are attractors that bound long-term motion, but a limit cycle is a simple closed loop traversed periodically, while a strange attractor is a fractal set on which the motion is chaotic and aperiodic." >}}

## See also

- [[Fractal]]
- [[Lorenz System]]
- [[Attractor]]
