---
title: Pendulum
tags: [oscillations]
summary: A swinging bob that looks like a simple harmonic oscillator only when it barely moves — and reveals its nonlinear soul at large angles.
weight: 22
---

# Pendulum

A **pendulum** is a bob of mass $m$ swinging on a rod of length $L$ under gravity. Resolving the gravitational pull along the arc gives an equation of motion for the angle $\theta$ from vertical:

$$\ddot\theta = -\frac{g}{L}\sin\theta.$$

That single $\sin\theta$ is the whole story. It makes the pendulum *nonlinear*, and it is the difference between a textbook idealization and the real, slightly unruly thing on the wall.

## The small-angle approximation

For small swings, $\sin\theta \approx \theta$, and the equation collapses to that of a [[Simple Harmonic Oscillator]]:

$$\ddot\theta \approx -\frac{g}{L}\theta, \qquad T \approx 2\pi\sqrt{\frac{L}{g}}.$$

In this regime the period is independent of amplitude — Galileo's famous observation. But the approximation quietly fails as the swing grows: a pendulum released from near-horizontal takes noticeably *longer* per swing than the formula predicts, because $\sin\theta < \theta$ weakens the restoring pull at large angles.

{{< sim name="pendulum" length="1.5" angle="2.4" caption="Released from 2.4 rad — far outside the small-angle regime. Watch it linger near the top of each swing." >}}

## The full nonlinear swing

Across its whole range the pendulum has a far richer [[Phase Portrait]] than the SHO's tidy ellipses:

- For low energy it **oscillates**, tracing closed loops around the stable hanging [[Fixed Point]] at $\theta = 0$.
- The upright position $\theta = \pi$ is an *unstable* fixed point — a saddle.
- Give it enough energy and it stops swinging back and forth and instead **rotates** over the top, circulating forever.

The boundary between swinging and spinning is a special trajectory called the *separatrix*. Push the idea further — pin the pivot to a motor, or chain two together into a [[Double Pendulum]] — and the pendulum becomes a gateway to [[Chaos]].

{{< quiz question="Why does a wide-swinging pendulum have a longer period than the small-angle formula T = 2π√(L/g) predicts?" options="Air resistance slows it down|Because sin θ < θ, so the restoring torque is weaker than the linear approximation at large angles|The rod stretches|Gravity is weaker at the top of the swing" answer="2" explain="The true restoring term is sin θ, which falls below the linear θ used in the SHO approximation. Weaker restoring torque means slower return and a longer period — an amplitude dependence the linear model can't capture." >}}

## See also

- [[Simple Harmonic Oscillator]]
- [[Double Pendulum]]
- [[Phase Portrait]]
