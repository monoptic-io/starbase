---
title: Gibbs Phenomenon
aliases: [ringing, overshoot]
tags: [fourier]
summary: The persistent ~9% overshoot that a truncated Fourier sum makes near a jump discontinuity — and that never goes away no matter how many terms you add.
weight: 60
---

# Gibbs Phenomenon

Build a square wave from its [[Fourier Series]] and stop after finitely many [[Harmonics]], and something stubborn happens at every sharp edge: the partial sum **overshoots** the true value, ringing above and below the jump before it settles. Adding more terms squeezes the ripples into a narrower and narrower band around the discontinuity — but it does **not** shrink the height of the overshoot. The peak holds fast at about **9%** of the jump, forever. This is the **Gibbs phenomenon**.

It is a genuinely counterintuitive result. We expect "more terms = better approximation," and in most senses that is true: the energy of the error goes to zero. Yet right at the edge, the maximum error refuses to die.

## Watch the overshoot refuse to die

Move your mouse left to right to set the number of harmonics in the partial sum of a square wave. The dashed line marks the true level; notice the little horn of overshoot that hugs the jump. As you add terms, the horn slides *closer* to the edge and gets *thinner* — but its height barely budges.

{{< sketch height="340" caption="Move the mouse to add harmonics. The overshoot near each jump narrows but stays ~9% tall — the Gibbs phenomenon." >}}
if (frame === 0) {}
const css = getComputedStyle(document.documentElement);
const accent = css.getPropertyValue('--accent').trim() || '#5b9cff';
const warn = css.getPropertyValue('--warn').trim() || '#ffb454';
const faint = css.getPropertyValue('--text-faint').trim() || '#7a8190';
const border = css.getPropertyValue('--border').trim() || '#2a2f3a';
const frac = mouse.x > 0 ? mouse.x / W : 0.25;
const terms = Math.max(1, Math.round(1 + frac * 39));
const cy = H * 0.5, amp = H * 0.34;
ctx.strokeStyle = border; ctx.globalAlpha = 0.6; ctx.lineWidth = 1;
ctx.beginPath(); ctx.moveTo(0, cy); ctx.lineTo(W, cy); ctx.stroke();
ctx.globalAlpha = 1;
ctx.strokeStyle = faint; ctx.globalAlpha = 0.55; ctx.lineWidth = 1.5;
ctx.setLineDash([5, 5]);
ctx.beginPath();
for (let px = 0; px <= W; px++) {
  const xv = px / W * Math.PI * 2;
  const sq = (xv < Math.PI) ? 1 : -1;
  const y = cy - amp * sq;
  if (px === 0) ctx.moveTo(px, y); else ctx.lineTo(px, y);
}
ctx.stroke();
ctx.setLineDash([]);
ctx.globalAlpha = 0.4;
ctx.strokeStyle = warn;
ctx.beginPath(); ctx.moveTo(0, cy - amp * 1.09); ctx.lineTo(W, cy - amp * 1.09); ctx.stroke();
ctx.globalAlpha = 1;
ctx.strokeStyle = accent; ctx.lineWidth = 2;
ctx.beginPath();
for (let px = 0; px <= W; px++) {
  const xv = px / W * Math.PI * 2;
  let s = 0;
  for (let k = 1; k <= terms; k++) {
    const n = 2 * k - 1;
    s += Math.sin(n * xv) / n;
  }
  s *= 4 / Math.PI;
  const y = cy - amp * s;
  if (px === 0) ctx.moveTo(px, y); else ctx.lineTo(px, y);
}
ctx.stroke();
ctx.fillStyle = warn; ctx.font = '12px sans-serif';
ctx.fillText('≈ 9% overshoot', 10, cy - amp * 1.09 - 6);
ctx.fillStyle = faint;
ctx.fillText(terms + ' harmonic' + (terms === 1 ? '' : 's'), 10, H - 12);
{{< /sketch >}}

## Why it happens

The partial sum is the true signal smeared by a fixed averaging window (the *Dirichlet kernel*). That window has wiggling side-lobes whose *integrated* area shrinks with more terms but whose *peak* relative size is a constant set by the integral $\tfrac{2}{\pi}\int_0^{\pi}\tfrac{\sin u}{u}\,du \approx 1.0895$ — hence the famous ~8.95% overshoot above the jump.

{{< note kind="warning" title="It is not a bug to be fixed by 'more terms'" >}}
Throwing more harmonics at a discontinuity will never remove the overshoot; it only moves it closer to the edge. To genuinely tame the ringing you must *taper* the harmonics' weights — a technique called **windowing** — trading a little edge sharpness for a smooth result. The same ringing shows up whenever a sharp feature is reconstructed from limited bandwidth, from audio to JPEG images.
{{< /note >}}

## See also

- [[Fourier Series]]
- [[Harmonics]]
- [[Convolution]]
