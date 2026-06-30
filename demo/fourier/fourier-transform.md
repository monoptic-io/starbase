---
title: Fourier Transform
tags: [fourier]
summary: The generalization of the Fourier series to non-periodic signals — an integral that resolves any signal into a continuous band of frequencies.
weight: 20
---

# Fourier Transform

A [[Fourier Series]] only works for *periodic* signals, where the allowed frequencies form a discrete ladder of [[Harmonics]]. But most real signals — a single clap, a transient, a pulse that happens once and never again — are not periodic. The **Fourier transform** handles them by letting the period stretch to infinity, at which point the discrete ladder of harmonics fills in to a *continuum*. The sum becomes an integral:

$$X(f) = \int_{-\infty}^{\infty} x(t)\,e^{-2\pi i f t}\,dt.$$

Here $x(t)$ is the signal in time and $X(f)$ is its representation in the [[Frequency Domain]]: a complex number for every frequency $f$, whose magnitude says *how much* of that frequency is present and whose phase says *where* it sits. The factor $e^{-2\pi i f t} = \cos(2\pi f t) - i\sin(2\pi f t)$ is a pure rotation — the transform correlates the signal against a sinusoid of every possible frequency at once.

## A perfectly reversible change of view

Nothing is lost. The **inverse transform** rebuilds the signal exactly, reassembling it from its frequency content:

$$x(t) = \int_{-\infty}^{\infty} X(f)\,e^{+2\pi i f t}\,df.$$

So $x(t)$ and $X(f)$ are two encodings of one object. This is the [[Frequency Domain|time ↔ frequency]] duality made precise.

{{< note kind="note" title="Why the complex exponential?" >}}
Sines and cosines are awkward to manipulate; the single complex exponential $e^{-2\pi i ft}$ packages both into one term and turns calculus into algebra. Differentiation in time becomes multiplication by $2\pi i f$ in frequency — which is why the transform makes differential equations so much easier.
{{< /note >}}

## The properties that make it powerful

Three structural facts do most of the heavy lifting in applications:

- **Linearity** — the transform of a sum is the sum of the transforms, the [[Superposition Principle]] in another guise.
- **The [[Convolution|convolution theorem]]** — convolving two signals in time is the same as *multiplying* their transforms in frequency. Filtering becomes a single multiplication.
- **Duality and scaling** — squeezing a signal in time stretches its [[Spectrum]] in frequency, the seed of the [[Uncertainty Principle]].

## See also

- [[Frequency Domain]]
- [[Convolution]]
- [[Uncertainty Principle]]
- [[Spectrum]]
