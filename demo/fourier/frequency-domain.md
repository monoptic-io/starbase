---
title: Frequency Domain
aliases: [time domain]
tags: [fourier]
summary: The second of the two equivalent views of a signal — describing it by which frequencies it contains rather than by how it varies in time.
weight: 30
---

# Frequency Domain

Every signal can be looked at in two ways. In the **time domain** you plot amplitude against the clock: the raw wiggle as it happens. In the **frequency domain** you plot how much of each pure frequency the signal holds. The [[Fourier Transform]] carries you from the first to the second and back, and — crucially — *no information is lost in the trip*. They are two photographs of the same statue taken from perpendicular angles.

## One signal, two pictures

Why bother with a second view? Because facts that are hidden in one become obvious in the other. A noisy recording looks like chaos in time, but in frequency the steady hum of a power line is a single sharp spike you can simply delete. A chord sounds like one complex pressure wave in time, but in frequency it is plainly three notes. The frequency domain is where structure that is *smeared across all of time* collapses into a few clean peaks.

Below, your mouse sets a frequency. The top trace is that pure tone in the time domain; the bottom mark is the very same tone in the frequency domain — a single spike whose *position* is the frequency. One wiggly curve, one lone dot: identical content.

{{< sketch height="340" caption="Move the mouse left↔right to change frequency. Top: the tone in time. Bottom: the same tone as a single spike in frequency." >}}
if (frame === 0) {}
const css = getComputedStyle(document.documentElement);
const accent = css.getPropertyValue('--accent').trim() || '#5b9cff';
const accent2 = css.getPropertyValue('--accent-2').trim() || '#ff7eb6';
const faint = css.getPropertyValue('--text-faint').trim() || '#7a8190';
const border = css.getPropertyValue('--border').trim() || '#2a2f3a';
const frac = mouse.x > 0 ? mouse.x / W : 0.4;
const freq = 1 + frac * 7;
const topY = H * 0.30, amp = H * 0.17;
ctx.strokeStyle = border; ctx.globalAlpha = 0.6; ctx.lineWidth = 1;
ctx.beginPath(); ctx.moveTo(0, topY); ctx.lineTo(W, topY); ctx.stroke();
ctx.globalAlpha = 1;
ctx.strokeStyle = accent; ctx.lineWidth = 2;
ctx.beginPath();
for (let px = 0; px <= W; px++) {
  const xv = px / W * Math.PI * 4;
  const yv = topY - amp * Math.sin(freq * xv + t);
  if (px === 0) ctx.moveTo(px, yv); else ctx.lineTo(px, yv);
}
ctx.stroke();
const axisY = H * 0.82;
ctx.strokeStyle = border; ctx.globalAlpha = 0.7; ctx.lineWidth = 1;
ctx.beginPath(); ctx.moveTo(0, axisY); ctx.lineTo(W, axisY); ctx.stroke();
ctx.globalAlpha = 1;
const fx = (freq - 1) / 7 * (W - 20) + 10;
ctx.strokeStyle = accent2; ctx.lineWidth = 3;
ctx.beginPath(); ctx.moveTo(fx, axisY); ctx.lineTo(fx, axisY - H * 0.24); ctx.stroke();
ctx.fillStyle = accent2;
ctx.beginPath(); ctx.arc(fx, axisY - H * 0.24, 4, 0, 7); ctx.fill();
ctx.fillStyle = faint; ctx.font = '12px sans-serif';
ctx.fillText('time domain', 8, topY - amp - 10);
ctx.fillText('frequency domain', 8, axisY + 18);
ctx.fillText('low f', 6, axisY - 4);
ctx.fillText('high f', W - 42, axisY - 4);
{{< /sketch >}}

## Reading the two domains

A few translations are worth memorizing, because they recur everywhere:

- A **slow** variation in time sits at the **left** (low frequency); a **fast** one sits at the **right** (high frequency).
- A **sharp edge or spike** in time spreads across **many** frequencies — see the [[Gibbs Phenomenon]] and the [[Uncertainty Principle]].
- A **pure repeating** signal collapses to a **few sharp peaks** — its [[Spectrum]].

{{< note kind="tip" title="Same data, better basis" >}}
Moving to the frequency domain is just choosing a smarter coordinate system. The signal hasn't changed; you have rotated your axes onto the directions — the sinusoids — along which the signal is simplest to describe.
{{< /note >}}

## See also

- [[Spectrum]]
- [[Fourier Transform]]
- [[Harmonics]]
