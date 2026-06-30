---
title: Spectrum
aliases: [power spectrum, frequency spectrum]
tags: [fourier]
summary: The magnitude of a signal's frequency content — the map of where its energy lives, read as a landscape of peaks.
weight: 50
---

# Spectrum

The **spectrum** of a signal is the magnitude of its [[Fourier Transform]]: for each frequency, *how much* of it is present, with the phase set aside. Where the full transform $X(f)$ is a complex number, the spectrum keeps only $|X(f)|$ (or its square, the **power spectrum**). It answers the most common question we ask of a signal — *what frequencies is it made of, and how strong are they?* — as a single readable picture.

Reading a signal as a spectrum means reading it as a landscape of **peaks**. A tall peak at some frequency says "there is a strong, steady oscillation right here." Broad, low hills say "noise, energy spread everywhere." The whole craft of spectral analysis is learning to read that terrain.

## A line spectrum

A periodic signal has a **discrete** spectrum: energy lives only at the fundamental and its [[Harmonics]], so the picture is a row of sharp lines. Here is the spectrum of a sawtooth-like tone, whose harmonics fall off as $1/n$:

{{< chart type="bar" data="1:1, 2:0.5, 3:0.33, 4:0.25, 5:0.2, 6:0.167, 7:0.143, 8:0.125" title="A 1/n line spectrum (harmonic number : amplitude)" >}}

Each bar marks one harmonic; the steady decay is the spectral signature of a signal with sharp corners. Smooth signals decay fast; jagged ones decay slowly.

## What the spectrum tells you

- **Pitch** is the lowest strong peak — the fundamental.
- **Timbre** is the *pattern* of the remaining peaks (see [[Harmonics]]).
- **Noise** is a raised, featureless floor between the peaks.
- **A pure tone** is a single spike; a **click** is a flat spectrum touching every frequency at once.

{{< note kind="key" title="Magnitude keeps the 'what', drops the 'where'" >}}
By discarding phase, the spectrum tells you *which* frequencies are present but not *how they line up* in time. Two very different-looking signals can share a spectrum. To reconstruct the original signal exactly you need the full [[Fourier Transform]], phase included.
{{< /note >}}

## See also

- [[Fourier Transform]]
- [[Harmonics]]
- [[Frequency Domain]]
