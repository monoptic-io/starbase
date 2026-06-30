---
title: Three-Body Problem
aliases: [three-body, three-body problem, 3-body problem]
tags: [systems, emergence, chaos, gravitation]
summary: Three gravitating masses — the smallest many-body system, and the birthplace of chaos in celestial mechanics.
weight: 30
---

# Three-Body Problem

The **three-body problem** is the [[N-Body Problem]] with $N = 3$: three masses, each attracting the other two, set loose under gravity. It looks like a modest step up from the tidy, solvable two-body case — yet it is one of the most consequential problems in the history of science. There is no general formula for the orbits. For almost all starting conditions the motion is **chaotic**, and the three-body problem is where mathematicians first understood what that means.

## Why it cannot be solved

The two-body problem is *integrable*: enough conserved quantities exist to reduce it to a curve you can write down. Add a third body and the count falls short — the conserved quantities no longer suffice to pin the motion, and the leftover freedom fills with chaotic wandering. When Poincaré studied this in the 1880s he found that trajectories near certain orbits tangle into an infinitely intricate web. He had, without quite naming it, discovered [[Chaos]].

{{< note kind="key" title="Non-integrability in one sentence" >}}
A system is *integrable* when conserved quantities constrain it onto simple surfaces you can solve for; it is *non-integrable* when they do not. The three-body problem is the canonical non-integrable system — the reason "closed-form solution" gives way to "numerical simulation."
{{< /note >}}

## Sensitive dependence

Run two three-body systems whose starting positions differ by a hair, and their orbits track each other for a while — then diverge completely. This is [[Sensitive Dependence on Initial Conditions]], the defining symptom of chaos and the reason long-term prediction is hopeless in practice even though the equations are exactly deterministic. The rate of separation is measured by a positive [[Lyapunov Exponent]].

{{< sim name="nbody" bodies="3" caption="Three gravitating masses. Reload and the dance is never the same twice; tiny differences explode into wholly different orbits." >}}

## Islands of order

Chaos is not the whole story. Hidden among the tangle are exact, eternally repeating **periodic solutions** — special arrangements that thread the needle. Lagrange found configurations where three bodies sit at the corners of an equilateral triangle and rotate rigidly; more recently, the "figure-eight" orbit was discovered, in which three equal masses chase one another along a single looping path. These choreographies are exquisitely balanced and, like a pencil on its tip, mostly unstable — but they show that order and chaos share the same equations.

## See also

- [[N-Body Problem]]
- [[Chaos]]
- [[Sensitive Dependence on Initial Conditions]]
