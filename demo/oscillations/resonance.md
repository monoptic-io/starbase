---
title: Resonance
tags: [oscillations]
summary: Drive an oscillator near its natural frequency and the response explodes — the principle behind tuning radios, shattering glass, and toppling bridges.
weight: 25
---

# Resonance

**Resonance** is what happens when you push an oscillator at just the right rhythm. As the driving frequency $\omega$ of a [[Driven Oscillator]] sweeps toward its natural frequency $\omega_0$, the steady-state amplitude rises to a sharp peak. Each push arrives perfectly in phase to add a little more energy than damping removes, and small forces accumulate into large motions.

## The response curve

Plotting amplitude against driving frequency (here in units of $\omega_0$, with light damping) gives the characteristic resonance peak — a **Lorentzian** lineshape:

{{< plot fn="1/Math.sqrt((1-x*x)*(1-x*x)+0.04*x*x)" xmin="0" xmax="3" title="Resonance response curve" caption="Amplitude vs. driving frequency ω/ω₀. The towering peak near 1 is resonance; lighter damping makes it taller and sharper." >}}

Two features matter. The **height** of the peak is set by damping — less damping means a taller, more dangerous spike (in the frictionless limit it diverges). The **width** measures how selective the resonance is: a narrow peak responds only to a tiny band of frequencies, which is exactly what lets a radio pick one station out of many.

## Resonance in the wild

- **Tuning** — an RLC circuit resonates at one frequency; turning the dial retunes $\omega_0$ to select a station.
- **Music** — the body of a violin and the column of air in an organ pipe resonate to amplify particular notes.
- **Destruction** — soldiers break step on bridges because a marching cadence near a structural $\omega_0$ can pump the span to failure. The same physics shatters a wine glass at its ringing pitch.

{{< quiz question="You drive an oscillator exactly at its natural frequency and then reduce the damping. What happens to the steady-state amplitude at the peak?" options="It decreases|It stays the same — only the frequency matters|It increases, and the peak also gets narrower|It drops to zero" answer="3" explain="Less damping removes less energy per cycle, so the resonant amplitude grows (diverging in the frictionless limit). The peak also sharpens — high-Q systems resonate strongly but only over a very narrow frequency band." >}}

{{< note kind="warning" title="Resonance cuts both ways" >}}
The same tall, narrow peak that makes a sensor exquisitely sensitive makes a bridge or a building catastrophically vulnerable. Engineers deliberately add [[Damped Oscillator|damping]] to flatten the peak — trading sensitivity for safety.
{{< /note >}}

## See also

- [[Driven Oscillator]]
- [[Simple Harmonic Oscillator]]
- [[Coupled Oscillators]]
