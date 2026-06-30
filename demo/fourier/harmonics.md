---
title: Harmonics
aliases: [overtones, fundamental frequency, timbre]
tags: [fourier]
summary: The integer multiples of a fundamental frequency whose relative strengths give every sound its characteristic color.
weight: 40
---

# Harmonics

When something vibrates with a definite pitch — a string, an air column, a bell — it rarely vibrates at just one frequency. It rings at a **fundamental frequency** $f_0$ *and* at a whole ladder of **harmonics** (also called **overtones**) at integer multiples $2f_0, 3f_0, 4f_0, \dots$ all at once. The fundamental sets the pitch you name; the harmonics, and how strong each one is, set the **timbre** — the quality that lets you tell instruments apart.

This ladder is exactly the set of frequencies a [[Fourier Series]] is built from. A periodic sound *is* its harmonics, summed.

## Why a violin and a flute sound different

Play the note A at 440 Hz on a violin and on a flute and you hear the same pitch — both have their fundamental at 440 Hz — yet no one confuses them. The difference is the **recipe of harmonic amplitudes**. A flute is nearly a pure tone: a strong fundamental and weak overtones. A bowed violin is rich and buzzy: substantial energy in many higher harmonics. Same pitch, different spectral fingerprint.

{{< chart type="bar" data="1,0.12,0.06,0.03,0.02,0.01,0.008" title="Flute-like spectrum: strong fundamental, faint overtones" >}}

{{< chart type="bar" data="1,0.8,0.65,0.5,0.55,0.32,0.28" title="Violin-like spectrum: rich, slowly decaying overtones" >}}

Read left to right, each bar is the strength of the fundamental, then the 2nd harmonic, the 3rd, and so on. The two instruments draw the same first bar but utterly different silhouettes — and your ear reads that silhouette as "flute" or "violin."

## Where the ladder comes from

Harmonics are not arbitrary; they are forced by geometry. A string clamped at both ends can only support a [[Standing Wave]] whose length fits a whole number of half-wavelengths, and those allowed modes have frequencies $f_0, 2f_0, 3f_0,\dots$ — the harmonic series exactly. Driving such a system near one of these frequencies produces [[Resonance]], which is why instruments speak so readily at their natural pitches.

{{< note kind="note" title="Harmonics vs. overtones" >}}
The two words count differently. The *fundamental* is the 1st harmonic. The *first overtone* is the 2nd harmonic. "Harmonic" counts from the fundamental; "overtone" counts the tones *above* it. Same ladder, offset labels.
{{< /note >}}

## See also

- [[Fourier Series]]
- [[Standing Wave]]
- [[Resonance]]
