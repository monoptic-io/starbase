---
title: Double Pendulum
aliases: [double pendulum]
tags: [chaos]
summary: Two rods hinged end to end — the simplest everyday mechanical system that exhibits genuine chaos.
weight: 40
---

# Double Pendulum

Hang one [[Pendulum]] from the bottom of another and you have built the **double pendulum** — perhaps the simplest contraption in all of physics whose motion is genuinely [[Chaos|chaotic]]. A single pendulum swings predictably; add a second joint and the system erupts into tumbling, flipping, never-repeating motion that no one can forecast more than a few seconds ahead.

What makes it remarkable is that there is no trick, no hidden randomness. The equations are exact Newtonian mechanics, fully deterministic. Yet the double pendulum is the textbook demonstration of [[Sensitive Dependence on Initial Conditions]]: release two of them from *almost* the same angle and within seconds they are doing completely different things.

## Try to predict it

Give the simulation a push and watch. For large swings the inner and outer arms exchange energy in a way that quickly becomes impossible to anticipate. The trajectory traced by the outer tip is a tangled, space-filling scribble.

{{< sim name="doublependulum" height="460" caption="A double pendulum. Drag to set it swinging, or let it fall from the horizontal — small differences in release lead to wildly different motion." >}}

{{< note kind="tip" title="Energy decides the character" >}}
The double pendulum is not chaotic at *all* energies. Nudged gently from hanging-straight-down, it behaves like a pair of [[Coupled Oscillators]] — quasi-periodic, even pretty. Crank up the energy so the arms can swing over the top and the motion becomes fully chaotic. Energy is the knob that drives this system through a sequence of [[Bifurcation|bifurcations]] from order into chaos, beautifully exposed by a [[Poincaré Section]].
{{< /note >}}

## Why it goes chaotic

The configuration is fixed by two angles, $\theta_1$ and $\theta_2$, but because the system is second-order its full [[State Space]] is four-dimensional: $(\theta_1, \theta_2, \dot\theta_1, \dot\theta_2)$. That is more than enough room for chaos — recall a continuous flow needs only three dimensions to escape the Poincaré–Bendixson trap.

The Lagrangian equations of motion are nastily nonlinear, coupling the two arms through $\sin$ and $\cos$ of the angle difference:

$$ (m_1+m_2)\ell_1\ddot\theta_1 + m_2\ell_2\ddot\theta_2\cos(\theta_1-\theta_2) + m_2\ell_2\dot\theta_2^{\,2}\sin(\theta_1-\theta_2) + (m_1+m_2)g\sin\theta_1 = 0. $$

There is no closed-form solution. The nonlinear coupling is precisely the stretch-and-fold mechanism that gives the system a positive [[Lyapunov Exponent]] — the quantitative stamp of chaos.

## Conserved, but not predictable

If there is no friction, the double pendulum conserves total energy exactly — its trajectory is forever pinned to a constant-energy surface in state space. Conservation and chaos coexist happily: energy says *where* the motion may go, chaos says you cannot predict *when* it gets there. Slicing that energy surface with a [[Poincaré Section]] turns the bewildering flow into a cleaner two-dimensional map you can actually read.

{{< quiz question="A frictionless double pendulum conserves energy exactly. Does that make its motion predictable?" options="Yes — conservation laws always make motion predictable|No — energy constrains where it can go but not the chaotic timing of its motion|Yes — but only for small swings|No — because energy is not actually conserved" answer="2" explain="A conservation law confines the trajectory to a constant-energy surface, yet on that surface the motion can still be chaotic, with sensitive dependence making long-term prediction impossible." >}}

## See also

- [[Pendulum]]
- [[Sensitive Dependence on Initial Conditions]]
- [[Poincaré Section]]
