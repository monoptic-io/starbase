---
title: Fourier Analysis
tags: [fourier]
summary: Any signal is a sum of pure sinusoids — and reading it as such turns the tangled language of time into the clean language of frequency.
weight: 28
---

# Fourier Analysis

**Fourier analysis** rests on one astonishing claim: *any* signal — a plucked string, a heartbeat, a stock price, the brightness of a distant star — can be rebuilt by adding together pure sinusoids of different frequencies, amplitudes, and phases. Nothing about the original signal need look wavy. The sharp corner of a square pulse, the spike of a drumbeat, the slow drift of a tide: all of them are secretly choirs of sine waves singing in superposition.

That claim has a profound consequence. Every signal has *two* equally valid descriptions. One is the familiar view in **time**: amplitude as the clock ticks. The other is the view in **frequency**: how much of each pure tone the signal contains. These are not two different signals — they are the same object seen from two directions, and Fourier analysis is the machinery that rotates between them.

## The time ↔ frequency duality

In the time view you ask *when*. In the frequency view you ask *how fast and how strong*. A problem that is a knot in one view is often a single clean tug in the other — which is exactly why this duality is one of the most useful ideas in all of science and engineering.

{{< note kind="key" title="The whole section in one sentence" >}}
A signal in time and its [[Spectrum]] in frequency carry identical information; the [[Fourier Transform]] is the dictionary that translates between them, and almost everything else here is a consequence of that dictionary.
{{< /note >}}

## How the ideas build

- [[Fourier Series]] — a *periodic* signal is an exact sum of [[Harmonics]] of one fundamental.
- [[Fourier Transform]] — drop periodicity and the sum becomes an integral over all frequencies.
- [[Frequency Domain]] — the second of the two views, and how to read it.
- [[Spectrum]] — the magnitude of that view: where the energy lives.
- [[Gibbs Phenomenon]] — why truncated reconstructions ring at sharp edges.
- [[Discrete Fourier Transform]] — what happens when we sample, and the limits that imposes.
- [[Convolution]] — why multiplying in frequency means smearing in time, the heart of filtering.
- [[Uncertainty Principle]] — you cannot be sharp in both views at once.

## Where it connects

Fourier analysis grew up alongside the study of vibration, so it threads directly back into the rest of this knowledge base. The pure sinusoid is nothing but the motion of a [[Simple Harmonic Oscillator]]; a periodic signal is a superposition of such motions, exactly as the [[Superposition Principle]] for [[Waves]] would predict. The frequencies a system *prefers* to ring at are the subject of [[Resonance]], and when many oscillators interact you get the rich spectra of [[Coupled Oscillators]]. To go the other way — from these tools back to the physics of travelling disturbances — start with [[Oscillations]] and [[Waves]].

## See also

- [[Fourier Series]]
- [[Oscillations]]
- [[Waves]]
