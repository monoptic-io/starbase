---
title: Oscillations
tags: [oscillations]
summary: How systems swing, settle, and sustain rhythm — from a single bob on a spring to whole populations of synchronizing oscillators.
weight: 20
---

# Oscillations

An **oscillation** is motion that returns. Pull a mass on a spring aside and let go; nudge a pendulum; strike a tuning fork — each repeats a pattern in time rather than running off to infinity or grinding to a halt. Oscillation is the simplest non-trivial thing a [[Dynamical System]] can do, and it is everywhere: in clocks and bridges, hearts and lasers, planetary librations and AC circuits. This section builds the idea up one layer at a time.

We begin with the frictionless ideal, the [[Simple Harmonic Oscillator]], whose sinusoidal motion is the template for everything that follows. The [[Pendulum]] shows how that template is only an approximation — a small-angle shadow of a richer nonlinear system. Adding friction gives the [[Damped Oscillator]]; pushing back against that friction with an external force gives the [[Driven Oscillator]], and tuning the drive to the system's natural rhythm produces [[Resonance]]. Finally we leave the linear world entirely: a [[Limit Cycle]] is an oscillation a system *generates on its own*, and [[Coupled Oscillators]] reveal how separate rhythms lock together into collective motion.

## The arc of this section

- [[Simple Harmonic Oscillator]] — the linear restoring force and its sine-wave solution.
- [[Pendulum]] — small-angle harmonic motion versus the full nonlinear swing.
- [[Damped Oscillator]] — friction, decay, and the under/critical/over-damped trichotomy.
- [[Driven Oscillator]] — external forcing, transients, and steady state.
- [[Resonance]] — the dramatic amplitude peak when drive meets natural frequency.
- [[Limit Cycle]] — self-sustained oscillation with no external clock.
- [[Coupled Oscillators]] — normal modes, beating, and synchronization.

## Why it matters

{{< note kind="key" title="The unifying picture" >}}
Every oscillator traces a closed (or spiraling) loop in its [[State Space]]. Reading that loop in a [[Phase Portrait]] tells you more than any single time series: a center means perpetual oscillation, an inward spiral means damping, and an isolated closed orbit means a [[Limit Cycle]]. Keep the geometry in mind as you read.
{{< /note >}}

## See also

- [[Foundations]]
- [[Phase Portrait]]
- [[Chaos]]
