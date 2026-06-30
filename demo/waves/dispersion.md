---
title: Dispersion
aliases: [wave dispersion]
tags: [waves]
summary: When wave speed depends on frequency, a pulse made of many frequencies spreads and reshapes as it travels, and phase and group velocities part ways.
weight: 80
---

# Dispersion

**Dispersion** is what happens when different frequencies travel at different speeds. In the ideal [[Wave Equation]] every frequency moves at the same $c$, so a pulse keeps its shape forever. Real media rarely cooperate: deep-water waves, light in glass, and electrons in a crystal all carry high and low frequencies at *different* speeds. A pulse is a [[Superposition Principle|superposition]] of many frequencies, so if they march at different rates the pulse **spreads out and distorts** as it goes. The rainbow from a prism is dispersion: glass bends blue more than red because it is *slower* for blue.

## The dispersion relation

All the physics lives in one function — the **dispersion relation** $\omega(k)$, tying angular frequency to wavenumber. A non-dispersive medium has the straight line $\omega = ck$; anything curved means dispersion. From it come two distinct velocities:

$$v_{\text{phase}}=\frac{\omega}{k}\qquad\text{(speed of an individual crest)},$$

$$v_{\text{group}}=\frac{d\omega}{dk}\qquad\text{(speed of the overall envelope)}.$$

The **phase velocity** is how fast a single crest moves; the **group velocity** is how fast the *packet* of energy and information moves. When $\omega=ck$ they are equal and nothing disperses. When $\omega(k)$ curves, they differ — sometimes dramatically, with crests sliding through the envelope faster than the envelope itself advances.

{{< plot fn="x ;; x - 0.18*x*x*x" xmin="0" xmax="3" ymin="0" ymax="3.2" height="300" title="Dispersion relation ω(k): straight = non-dispersive, curved = dispersive" caption="The straight line ω = c·k (no dispersion) versus a curved relation. On the curve the slope dω/dk (group velocity) differs from the chord ω/k (phase velocity), so a pulse spreads." >}}

## Group velocity carries the message

A crucial subtlety: the **group velocity** is what carries energy and signal, not the phase velocity. It is entirely possible — in fact common — for individual crests to race ahead while the bundle they belong to lags behind, the crests appearing at the back of the packet, sweeping forward, and vanishing at the front. The object that makes this concrete is the [[Wave Packet]]: a localized envelope (moving at $v_{\text{group}}$) filled with oscillations (moving at $v_{\text{phase}}$).

{{< note kind="tip" title="Why a prism makes a rainbow" >}}
Decompose white light into its frequencies with a [[Fourier Transform]] and send each through glass. Because the refractive index — and hence the speed — depends on frequency, each color refracts by a different angle. Dispersion is the mechanism; the spectrum is the result.
{{< /note >}}

## See also

- [[Wave Packet]]
- [[Fourier Transform]]
- [[Wave Equation]]
