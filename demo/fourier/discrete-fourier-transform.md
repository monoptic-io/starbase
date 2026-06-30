---
title: Discrete Fourier Transform
aliases: [DFT, FFT, fast fourier transform, sampling]
tags: [fourier]
summary: The Fourier transform of sampled data — and the Nyquist limit and aliasing that sampling forces upon us.
weight: 70
---

# Discrete Fourier Transform

Computers cannot integrate a signal over all of continuous time; they hold a finite list of **samples**. The **discrete Fourier transform** (DFT) is the version of the [[Fourier Transform]] built for that reality. Given $N$ samples $x_0,\dots,x_{N-1}$ it returns $N$ frequency coefficients:

$$X_k = \sum_{n=0}^{N-1} x_n\,e^{-2\pi i kn/N},\qquad k = 0,\dots,N-1.$$

Computed naively this costs $N^2$ operations; the celebrated **fast Fourier transform (FFT)** reorganizes the same arithmetic into $N\log N$, and that single speedup is what put Fourier analysis inside every phone, radio, and oscilloscope.

## The price of sampling: Nyquist and aliasing

Sampling is not free. If you measure a signal only every $\Delta t$ seconds, you can faithfully represent frequencies only up to the **Nyquist limit** $f_{\text{Nyq}} = \tfrac{1}{2\Delta t}$ — half the sampling rate. Anything faster cannot be told apart from something slower: a high frequency, sampled too coarsely, **masquerades as a low one**. This impostor is called an **alias**.

Below, the faint curve is a genuine high-frequency sine. The dots are samples taken too sparsely, and the bold line connects them — a slow, lazy wave that was never really there. Your eye, like the DFT, is fooled.

{{< sketch height="320" caption="Faint: a true high-frequency sine. Dots: coarse samples. Bold: the low-frequency 'alias' they imply — a frequency that does not exist in the signal." >}}
if (frame === 0) {}
const css = getComputedStyle(document.documentElement);
const accent = css.getPropertyValue('--accent').trim() || '#5b9cff';
const accent2 = css.getPropertyValue('--accent-2').trim() || '#ff7eb6';
const faint = css.getPropertyValue('--text-faint').trim() || '#7a8190';
const border = css.getPropertyValue('--border').trim() || '#2a2f3a';
const cy = H * 0.5, amp = H * 0.32;
const fHigh = 11;
const drift = t * 0.5;
ctx.strokeStyle = border; ctx.globalAlpha = 0.5; ctx.lineWidth = 1;
ctx.beginPath(); ctx.moveTo(0, cy); ctx.lineTo(W, cy); ctx.stroke();
ctx.globalAlpha = 0.6;
ctx.strokeStyle = faint; ctx.lineWidth = 1.5;
ctx.beginPath();
for (let px = 0; px <= W; px++) {
  const ph = px / W * Math.PI * 2 * fHigh + drift;
  const y = cy - amp * Math.sin(ph);
  if (px === 0) ctx.moveTo(px, y); else ctx.lineTo(px, y);
}
ctx.stroke();
ctx.globalAlpha = 1;
const Ns = 12;
const pts = [];
for (let i = 0; i <= Ns; i++) {
  const px = i / Ns * W;
  const ph = px / W * Math.PI * 2 * fHigh + drift;
  pts.push({ x: px, y: cy - amp * Math.sin(ph) });
}
ctx.strokeStyle = accent2; ctx.lineWidth = 2.5;
ctx.beginPath();
pts.forEach((p, i) => i === 0 ? ctx.moveTo(p.x, p.y) : ctx.lineTo(p.x, p.y));
ctx.stroke();
ctx.fillStyle = accent;
pts.forEach(p => { ctx.beginPath(); ctx.arc(p.x, p.y, 4, 0, 7); ctx.fill(); });
ctx.fillStyle = faint; ctx.font = '12px sans-serif';
ctx.fillText('true signal (fast) vs. aliased reconstruction (slow)', 8, H - 10);
{{< /sketch >}}

## Living with the limit

The cure for aliasing is to remove offending frequencies *before* sampling with an **anti-aliasing filter**, or to sample fast enough that everything real sits below Nyquist. The same trap explains why wagon wheels spin backward on film and why a poorly chosen audio sample rate adds phantom tones. Once the DFT has the samples, its output is read as a [[Spectrum]] — but only the part below Nyquist can be trusted.

{{< note kind="warning" title="The Nyquist–Shannon rule" >}}
To capture a signal containing frequencies up to $f_{\max}$, you must sample at more than $2f_{\max}$. Sample slower and the lost high frequencies do not simply vanish — they fold down and contaminate the low frequencies you *can* see.
{{< /note >}}

## See also

- [[Fourier Transform]]
- [[Spectrum]]
- [[Frequency Domain]]
