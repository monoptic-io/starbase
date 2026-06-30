---
title: Uncertainty Principle
aliases: [time-frequency uncertainty, time-bandwidth product]
tags: [fourier]
summary: A signal cannot be sharply localized in both time and frequency at once — the tighter one view, the broader the other.
weight: 90
---

# Uncertainty Principle

There is a hard limit on how well a signal can be pinned down in both [[Frequency Domain|time and frequency]] at the same time. A signal that is **narrow in time** — a brief click — is necessarily **broad in frequency**, splattering energy across a wide band. A signal that is **narrow in frequency** — a pure, sustained tone — must be **broad in time**, droning on and on. You cannot have both sharp. Quantitatively, if $\Delta t$ and $\Delta f$ measure the spread in each domain,

$$\Delta t \,\Delta f \;\gtrsim\; \frac{1}{4\pi}.$$

This is *not* a statement about measurement clumsiness or quantum mechanics — it is a theorem about the [[Fourier Transform]] itself. The quantum Heisenberg principle is this very inequality applied to a particle's wavefunction.

## Squeeze one, the other spreads

The extremal case — the signal that comes *closest* to equality — is the Gaussian bell. Below, your mouse sets the width of a Gaussian pulse in time (left). Its transform, also a Gaussian, is drawn in frequency (right). Make the time pulse skinny and watch the frequency bell fatten; widen the time pulse and the frequencies collapse toward a single tone. Their product stays roughly fixed.

{{< sketch height="320" caption="Move the mouse to set the time-domain width (left). The frequency-domain width (right) moves the opposite way — their product is bounded below." >}}
if (frame === 0) {}
const css = getComputedStyle(document.documentElement);
const accent = css.getPropertyValue('--accent').trim() || '#5b9cff';
const accent2 = css.getPropertyValue('--accent-2').trim() || '#ff7eb6';
const faint = css.getPropertyValue('--text-faint').trim() || '#7a8190';
const border = css.getPropertyValue('--border').trim() || '#2a2f3a';
const frac = mouse.x > 0 ? mouse.x / W : 0.5;
const sigma = 0.05 + frac * 0.30;
const cxL = W * 0.25, cxR = W * 0.75;
const midY = H * 0.62, amp = H * 0.40;
const halfW = W * 0.46;
ctx.strokeStyle = border; ctx.globalAlpha = 0.7; ctx.lineWidth = 1;
ctx.beginPath(); ctx.moveTo(W / 2, H * 0.10); ctx.lineTo(W / 2, H * 0.90); ctx.stroke();
ctx.beginPath(); ctx.moveTo(0, midY); ctx.lineTo(W, midY); ctx.stroke();
ctx.globalAlpha = 1;
const wt = sigma * halfW;
ctx.strokeStyle = accent; ctx.lineWidth = 2;
ctx.beginPath();
for (let px = 0; px <= W / 2; px++) {
  const d = (px - cxL) / wt;
  const y = midY - amp * Math.exp(-d * d);
  if (px === 0) ctx.moveTo(px, y); else ctx.lineTo(px, y);
}
ctx.stroke();
const wf = (W * 0.010) / sigma;
ctx.strokeStyle = accent2; ctx.lineWidth = 2;
ctx.beginPath();
let started = false;
for (let px = Math.floor(W / 2); px <= W; px++) {
  const d = (px - cxR) / wf;
  const y = midY - amp * Math.exp(-d * d);
  if (!started) { ctx.moveTo(px, y); started = true; } else ctx.lineTo(px, y);
}
ctx.stroke();
ctx.fillStyle = faint; ctx.font = '12px sans-serif';
ctx.fillText('time domain', 10, H * 0.16);
ctx.fillText('frequency domain', W / 2 + 10, H * 0.16);
ctx.fillText('Δt · Δf  ≈  constant', W * 0.35, H - 10);
{{< /sketch >}}

## Why narrow demands wide

To synthesize a feature that is sharp and short-lived, you need sinusoids that *cancel everywhere except* in that brief window — and cancelling over all of time while reinforcing in one spot requires a wide range of frequencies beating against each other. Conversely, a single frequency, by definition, oscillates identically forever and cannot be confined. The math just makes this bookkeeping exact.

{{< note kind="note" title="A pervasive trade-off" >}}
The same bound governs the [[Wave Packet]] (localized in space ⇒ a spread of wavenumbers), the resolution of a spectrum analyzer (watch longer to resolve finer frequencies), and — through [[Dispersion]] — why a sharp pulse spreads as it travels. Sharpness is a resource you must spend in one domain to gain in the other.
{{< /note >}}

## See also

- [[Wave Packet]]
- [[Fourier Transform]]
- [[Dispersion]]
