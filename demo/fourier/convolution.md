---
title: Convolution
aliases: [convolution theorem]
tags: [fourier]
summary: The sliding-overlap operation that blends two signals — and which, magically, becomes simple multiplication in the frequency domain.
weight: 80
---

# Convolution

**Convolution** is the operation that slides one signal across another, multiplying and accumulating the overlap at every shift. For two signals $f$ and $g$:

$$(f * g)(\tau) = \int_{-\infty}^{\infty} f(u)\,g(\tau - u)\,du.$$

Intuitively, $g$ is a little **kernel** — a stencil — that you drag along $f$; at each position the output is the area of their overlap. Smear a signal with a blob-shaped kernel and you blur it; smear it with a kernel that has positive and negative parts and you sharpen or differentiate it. Every linear filter — every blur, echo, equalizer, and edge detector — is a convolution.

## The kernel sliding across the signal

Below, a bump-shaped kernel (the moving highlight) slides left to right across a two-peak signal (top). At each position the running output (bottom) records how much the kernel currently overlaps the signal. Where the kernel sits on a peak, the output rises; between peaks, it dips.

{{< sketch height="340" caption="Top: a fixed signal with a sliding kernel. Bottom: the convolution — the accumulated overlap at each position." >}}
if (frame === 0) {}
const css = getComputedStyle(document.documentElement);
const accent = css.getPropertyValue('--accent').trim() || '#5b9cff';
const accent2 = css.getPropertyValue('--accent-2').trim() || '#ff7eb6';
const good = css.getPropertyValue('--good').trim() || '#6ad08f';
const faint = css.getPropertyValue('--text-faint').trim() || '#7a8190';
const border = css.getPropertyValue('--border').trim() || '#2a2f3a';
const sig = (u) => Math.exp(-Math.pow((u - 0.32) * 7, 2)) + 0.75 * Math.exp(-Math.pow((u - 0.66) * 9, 2));
const Kw = 0.05;
const topY = H * 0.32, sigAmp = H * 0.22;
const baseY = H * 0.88, outAmp = H * 0.26;
ctx.strokeStyle = border; ctx.globalAlpha = 0.5; ctx.lineWidth = 1;
ctx.beginPath(); ctx.moveTo(0, topY); ctx.lineTo(W, topY); ctx.stroke();
ctx.beginPath(); ctx.moveTo(0, baseY); ctx.lineTo(W, baseY); ctx.stroke();
ctx.globalAlpha = 1;
ctx.strokeStyle = accent; ctx.lineWidth = 2;
ctx.beginPath();
for (let px = 0; px <= W; px++) {
  const u = px / W;
  const y = topY - sigAmp * sig(u);
  if (px === 0) ctx.moveTo(px, y); else ctx.lineTo(px, y);
}
ctx.stroke();
const kc = (t * 0.13) % 1;
ctx.fillStyle = good; ctx.globalAlpha = 0.22;
ctx.beginPath();
for (let px = 0; px <= W; px++) {
  const u = px / W;
  const k = Math.exp(-Math.pow((u - kc) / Kw, 2));
  const y = topY - sigAmp * k;
  if (px === 0) ctx.moveTo(px, y); else ctx.lineTo(px, y);
}
ctx.lineTo(W, topY); ctx.lineTo(0, topY); ctx.closePath(); ctx.fill();
ctx.globalAlpha = 1;
let maxOut = 0;
const outs = [];
const N = 90;
for (let i = 0; i <= N; i++) {
  const c = i / N;
  let a = 0;
  for (let u = 0; u <= 1; u += 0.02) a += sig(u) * Math.exp(-Math.pow((u - c) / Kw, 2)) * 0.02;
  outs.push(a);
  if (a > maxOut) maxOut = a;
}
ctx.strokeStyle = accent2; ctx.lineWidth = 2;
ctx.beginPath();
for (let i = 0; i <= N; i++) {
  const px = i / N * W;
  const y = baseY - outAmp * (outs[i] / maxOut);
  if (i === 0) ctx.moveTo(px, y); else ctx.lineTo(px, y);
}
ctx.stroke();
const ci = Math.round(kc * N);
const cyOut = baseY - outAmp * (outs[ci] / maxOut);
ctx.fillStyle = accent2;
ctx.beginPath(); ctx.arc(kc * W, cyOut, 4.5, 0, 7); ctx.fill();
ctx.strokeStyle = good; ctx.globalAlpha = 0.5; ctx.setLineDash([4, 4]);
ctx.beginPath(); ctx.moveTo(kc * W, topY); ctx.lineTo(kc * W, baseY); ctx.stroke();
ctx.setLineDash([]); ctx.globalAlpha = 1;
ctx.fillStyle = faint; ctx.font = '12px sans-serif';
ctx.fillText('signal · kernel', 8, topY - sigAmp - 8);
ctx.fillText('convolution output', 8, baseY + 14);
{{< /sketch >}}

## The convolution theorem

Convolution looks expensive — a sliding integral at every point — yet the [[Fourier Transform]] collapses it into something trivial. **Convolution in time is multiplication in frequency**, and vice versa:

$$\widehat{f * g} = \hat f \cdot \hat g.$$

This is the engine of practical signal processing. To filter a signal, you do not laboriously slide a kernel across it; you transform to the [[Frequency Domain]], *multiply* by the filter's response, and transform back. A blur becomes "attenuate the high frequencies"; an echo becomes "boost a comb of frequencies." The hard operation in one domain is the easy one in the other.

{{< note kind="key" title="Filtering is multiplication in disguise" >}}
Because of this theorem, designing a filter and choosing which frequencies to keep are the *same act*. The kernel you slide in time and the curve you multiply in frequency are a [[Fourier Transform]] pair.
{{< /note >}}

## See also

- [[Fourier Transform]]
- [[Frequency Domain]]
- [[Gibbs Phenomenon]]
