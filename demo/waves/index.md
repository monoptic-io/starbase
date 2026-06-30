---
title: Waves
aliases: [wave physics]
tags: [waves]
summary: From a single swinging oscillator to ripples that travel, interfere, and carry frequency through space — the bridge between Oscillations and Fourier Analysis.
weight: 25
---

# Waves

Set one oscillator going and it traces a rhythm in *time*. Now line up a chain of them, each nudging its neighbor, and that rhythm starts to *travel*: a disturbance in one place reappears, an instant later, in the next. That propagating disturbance is a **wave**. This section follows a single arc — **oscillation → wave → interference → frequency** — that turns the local, in-place motion of the [[Oscillations]] chapter into something that moves, spreads, overlaps, and ultimately decomposes into a spectrum.

A wave is what a [[Simple Harmonic Oscillator]] becomes when you give it room to spread. Every point of the medium oscillates just as a mass on a spring does; the wave is the *pattern those oscillations make across space*. That is the thread connecting this chapter backward to [[Oscillations]] and forward to [[Fourier Analysis]], where we learn that *any* wave is a sum of pure sinusoidal oscillations.

## The arc of this section

- [[Wave]] — anatomy of a moving disturbance: wavelength, frequency, amplitude, and the master relation $v=f\lambda$.
- [[Wave Equation]] — the single partial differential equation $\partial_{tt}u=c^2\partial_{xx}u$ that every non-dispersive wave obeys, and d'Alembert's left- and right-movers.
- [[Superposition Principle]] — the linearity that lets waves pass through one another and add up cleanly. Everything else is a consequence of it.
- [[Standing Wave]] — what happens when a wave is trapped: nodes, antinodes, and the harmonic ladder $f_n=nv/2L$.
- [[Interference]] — two sources, one pattern of bright and dark fringes set by path difference.
- [[Beats]] — interference in *time*: two near frequencies producing a slow throb.
- [[Doppler Effect]] — how motion of the source or observer shifts the frequency you receive.
- [[Dispersion]] — when the speed depends on frequency, and phase and group velocities part ways.
- [[Wave Packet]] — a localized lump of wave, the object that carries both energy and information.

## Where it connects

{{< note kind="key" title="One idea, two chapters" >}}
A wave is oscillation unrolled across space, so it inherits everything from [[Oscillations]] — restoring forces, natural frequencies, [[Resonance]] — and adds propagation. And because every wave is built from pure sinusoids, this chapter hands you straight to [[Fourier Analysis]]: [[Standing Wave|standing waves]] are the [[Harmonics]], a [[Wave Packet]] is a [[Fourier Transform|Fourier]] superposition, and [[Beats]] are a two-tone [[Spectrum]] you can hear.
{{< /note >}}

## See also

- [[Oscillations]]
- [[Coupled Oscillators]]
- [[Fourier Analysis]]
