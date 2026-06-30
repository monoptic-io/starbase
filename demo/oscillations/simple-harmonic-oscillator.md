---
title: Simple Harmonic Oscillator
aliases: [SHO, harmonic oscillator]
tags: [oscillations]
summary: A mass pulled back by a force proportional to its displacement — the linear ideal whose sinusoidal motion underlies every other oscillator.
weight: 21
---

# Simple Harmonic Oscillator

The **simple harmonic oscillator** (SHO) is the purest oscillation there is: a system whose restoring force grows in exact proportion to how far it has been displaced. For a mass $m$ on a spring of stiffness $k$, Newton's second law reads

$$m\ddot x = -kx,$$

the minus sign saying the force always points *back* toward equilibrium. Divide by $m$ and define the **natural angular frequency** $\omega_0 = \sqrt{k/m}$ to get the equation in its canonical form, $\ddot x = -\omega_0^2 x$.

## The sinusoidal solution

Any function whose second derivative is its own negative (up to a constant) is a sine or cosine, so the general motion is

$$x(t) = A\cos(\omega_0 t + \varphi),$$

with amplitude $A$ and phase $\varphi$ fixed by the initial position and velocity. The motion is **isochronous**: the period $T = 2\pi/\omega_0 = 2\pi\sqrt{m/k}$ depends on the mass and stiffness but *not* on the amplitude. A gentle wobble and a violent swing take exactly the same time — a fact that made the harmonic oscillator the heart of mechanical clocks.

{{< plot fn="Math.cos(x)" title="x(t) = cos(t): the canonical harmonic motion" caption="One full cycle every 2π. Amplitude sets the height; it never changes the period." >}}

## Energy flows back and forth

The total energy is constant, but it sloshes between two forms — kinetic and potential:

$$E = \underbrace{\tfrac12 m\dot x^2}_{\text{kinetic}} + \underbrace{\tfrac12 k x^2}_{\text{potential}} = \tfrac12 k A^2.$$

At the turning points all the energy is potential; flying through the center it is all kinetic. Because energy never leaves, the trajectory in [[State Space]] is a closed ellipse traced forever — the [[Phase Portrait]] of a center.

{{< note kind="note" title="Why the SHO is everywhere" >}}
Near any smooth potential minimum, the energy looks like a parabola, $V(x) \approx V_0 + \tfrac12 V''(x_0)\,(x-x_0)^2$. So *every* stable equilibrium behaves like a simple harmonic oscillator for small enough disturbances. That is why the SHO is the first model physicists reach for — and why the [[Pendulum]] is "simple harmonic" only when it barely swings.
{{< /note >}}

## See also

- [[Pendulum]]
- [[Damped Oscillator]]
- [[Resonance]]
