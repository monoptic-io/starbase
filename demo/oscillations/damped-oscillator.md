---
title: Damped Oscillator
aliases: [damping]
tags: [oscillations]
summary: Add friction to a harmonic oscillator and its energy bleeds away — gently ringing down, just barely settling, or sluggishly creeping home.
weight: 23
---

# Damped Oscillator

Real oscillators lose energy. Add a velocity-proportional drag — air resistance, internal friction, electrical resistance — to the [[Simple Harmonic Oscillator]] and you get the **damped oscillator**:

$$m\ddot x + c\dot x + kx = 0.$$

The new middle term $c\dot x$ always opposes the motion, draining energy on every pass. How the system relaxes depends entirely on how the damping $c$ compares to the stiffness and mass, summarized by the dimensionless **damping ratio** $\zeta = c/(2\sqrt{mk})$.

## Three regimes

{{< columns count="3" >}}
**Underdamped ($\zeta < 1$)** — the system still oscillates, but inside a shrinking envelope $e^{-\zeta\omega_0 t}$. It rings down over many cycles. This is a plucked string or a struck bell.

**Critically damped ($\zeta = 1$)** — the fastest possible return to rest *without* overshooting. Car suspensions and door closers are tuned near here.

**Overdamped ($\zeta > 1$)** — so much drag that the system creeps back to equilibrium slowly, never crossing it. Think of a spoon settling in honey.
{{< /columns >}}

The plot below shows the underdamped case: an oscillation $e^{-0.3t}\cos(3t)$ caged between its two exponential envelopes.

{{< plot fn="Math.exp(-0.3*x)*Math.cos(3*x) ;; Math.exp(-0.3*x) ;; -Math.exp(-0.3*x)" title="Underdamped decay" caption="The oscillation (middle curve) is bounded by the decaying envelope ±e^(−0.3t). Each swing is a little smaller than the last." >}}

## Watch it ring down

Damping turns the conservative pendulum's perpetual loop into an inward spiral in the [[Phase Portrait]] — every orbit shrinks toward the stable [[Fixed Point]] at the bottom. Give the simulation below a healthy dose of friction and watch the swings die away:

{{< sim name="pendulum" length="1.5" angle="2.4" damping="0.4" caption="A pendulum with damping 0.4. Energy leaks out each swing until it hangs still at the bottom." >}}

{{< note kind="tip" title="Quality factor" >}}
Engineers summarize damping with the **Q factor**, roughly the number of oscillations before the amplitude decays appreciably. A wine glass has Q in the thousands (it rings); a car suspension has Q near $1/2$ (it just settles). Low Q is the price you pay to avoid [[Resonance]] disasters when you start *driving* the system.
{{< /note >}}

## See also

- [[Driven Oscillator]]
- [[Simple Harmonic Oscillator]]
- [[Stability]]
