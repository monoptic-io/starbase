---
title: Stability
aliases: [stable, unstable, linear stability]
tags: [foundations]
summary: Whether small disturbances to an equilibrium grow or decay — the difference between a valley and a hilltop.
weight: 60
---

# Stability

**Stability** asks a simple question: if you nudge a system away from an equilibrium, does it return or run away? A marble in a bowl rolls back — the bottom is **stable**. A marble balanced on a dome rolls off — the top is **unstable**. The same distinction governs every [[Fixed Point]] of a [[Dynamical System]], and it decides which equilibria you ever actually observe, since unstable ones are destroyed by the slightest perturbation.

## The Lyapunov picture

Aleksandr Lyapunov gave the idea its precise form. A fixed point $x^*$ is **stable** if trajectories that start nearby *stay* nearby, and **asymptotically stable** if they additionally *return* to it as $t \to \infty$. The intuition is energy: if you can find a quantity that only ever decreases along trajectories and bottoms out at $x^*$ — a *Lyapunov function*, the mathematical bowl — then the state is forced to slide downhill into equilibrium.

## Linearization and eigenvalues

For most purposes you can read stability straight off the linearization. Near $x^*$, a displacement $u$ obeys $\dot u \approx A u$ with $A = Df(x^*)$, and perturbations evolve like $e^{\lambda t}$ where $\lambda$ are the eigenvalues of $A$:

{{< eq number="1" >}}
\operatorname{Re}(\lambda) < 0 \ \text{for all } \lambda \;\Longrightarrow\; \text{stable}, \qquad \operatorname{Re}(\lambda) > 0 \ \text{for some } \lambda \;\Longrightarrow\; \text{unstable}.
{{< /eq >}}

The sign of the real part is everything: negative means the perturbation decays, positive means it explodes. The imaginary part only sets whether the approach is a smooth slide (real eigenvalues, a node) or a ringing spiral (complex eigenvalues).

## Decay versus growth

The plot contrasts two perturbations: one to a stable equilibrium, $e^{-0.7x}$, which shrinks away, and one to an unstable equilibrium, $e^{0.7x}$, which blows up. Same starting size, opposite fates — set entirely by the sign of the exponent.

{{< plot fn="Math.exp(-0.7*x);;Math.exp(0.7*x)" xmin="0" xmax="5" ymin="0" ymax="6" title="Stable (decaying) vs. unstable (growing) perturbation" caption="Lower curve: Re(lambda) < 0. Upper curve: Re(lambda) > 0." >}}

{{< note kind="warning" title="The borderline cases" >}}
When the largest real part is exactly zero — a center, or a zero eigenvalue — linearization is inconclusive. The nonlinear terms decide, and these *marginal* cases are precisely where a [[Bifurcation]] is about to happen.
{{< /note >}}

## Beyond equilibria

Stability is not only about fixed points. A whole trajectory can be stable or unstable, and the *rate* at which neighbors converge or diverge along it is measured by the [[Lyapunov Exponent]]. When that rate is positive on an [[Attractor]], you have the stretching that defines [[Chaos]].

## See also

- [[Lyapunov Exponent]]
- [[Bifurcation]]
- [[Fixed Point]]
