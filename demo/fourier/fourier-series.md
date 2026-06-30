---
title: Fourier Series
aliases: [harmonic synthesis]
tags: [fourier]
summary: A periodic signal is an exact sum of harmonics — sinusoids at integer multiples of one fundamental frequency.
weight: 10
---

# Fourier Series

A **Fourier series** expresses any reasonable periodic signal as a sum of sinusoids whose frequencies are integer multiples — the [[Harmonics]] — of a single **fundamental**. If a signal $f(t)$ repeats with period $T$, then

$$f(t) = a_0 + \sum_{n=1}^{\infty}\Big[a_n\cos(n\omega_0 t) + b_n\sin(n\omega_0 t)\Big],\qquad \omega_0 = \frac{2\pi}{T}.$$

Each coefficient measures *how much* of that harmonic the signal contains, recovered by the orthogonality of sines and cosines:

$$a_n = \frac{2}{T}\int_0^T f(t)\cos(n\omega_0 t)\,dt,\qquad b_n = \frac{2}{T}\int_0^T f(t)\sin(n\omega_0 t)\,dt.$$

The remarkable part is that the building blocks are nothing exotic: every term is the displacement of a [[Simple Harmonic Oscillator]], and the series is their [[Superposition Principle|superposition]].

## Building a square wave

The clearest demonstration is the square wave, which turns out to be a sum over the **odd** harmonics only, each weighted by $1/n$:

$$f(t) = \frac{4}{\pi}\left(\sin t + \frac{1}{3}\sin 3t + \frac{1}{5}\sin 5t + \cdots\right).$$

Add a few terms and the flat plateaus and steep edges of a square wave emerge out of pure curves:

{{< plot fn="(4/Math.PI)*Math.sin(x) ;; (4/Math.PI)*(Math.sin(x)+Math.sin(3*x)/3) ;; (4/Math.PI)*(Math.sin(x)+Math.sin(3*x)/3+Math.sin(5*x)/5+Math.sin(7*x)/7+Math.sin(9*x)/9)" title="Square wave from 1, 2, and 5 odd harmonics" caption="More harmonics → flatter tops and sharper edges. The wiggles never fully vanish — see Gibbs Phenomenon." >}}

## Epicycles: synthesis made visible

There is a beautiful mechanical picture of the same sum. Attach a circle to the rim of a bigger circle, and another to that one, and so on — one circle per harmonic, each spinning at its own integer rate with radius $\tfrac{4}{\pi}\cdot\tfrac1n$. The tip of the final arm traces the signal. Below, the chained circles for the odd harmonics turn, and the height of their tip is drawn out to the right: a square wave, written by spinning wheels.

{{< sketch height="360" caption="Chained epicycles for the odd harmonics 1, 1/3, 1/5, … Their tip traces a square wave to the right." >}}
if (frame === 0) { state.path = []; }
const css = getComputedStyle(document.documentElement);
const accent = css.getPropertyValue('--accent').trim() || '#5b9cff';
const accent2 = css.getPropertyValue('--accent-2').trim() || '#ff7eb6';
const faint = css.getPropertyValue('--text-faint').trim() || '#7a8190';
const cy = H * 0.5;
const cx = W * 0.30;
const scale = Math.min(H * 0.34, W * 0.18);
let x = cx, y = cy;
ctx.lineWidth = 1;
for (let i = 0; i < 8; i++) {
  const n = 2 * i + 1;
  const px = x, py = y;
  const radius = scale * 4 / (Math.PI * n);
  ctx.globalAlpha = 0.35;
  ctx.strokeStyle = faint;
  ctx.beginPath(); ctx.arc(px, py, radius, 0, Math.PI * 2); ctx.stroke();
  const ang = n * t;
  x += radius * Math.cos(ang);
  y += radius * Math.sin(ang);
  ctx.globalAlpha = 0.9;
  ctx.strokeStyle = accent;
  ctx.beginPath(); ctx.moveTo(px, py); ctx.lineTo(x, y); ctx.stroke();
}
ctx.globalAlpha = 1;
state.path.unshift(y);
const maxLen = Math.floor(W * 0.6);
if (state.path.length > maxLen) state.path.pop();
const waveX = W * 0.55;
ctx.strokeStyle = faint;
ctx.globalAlpha = 0.5;
ctx.beginPath(); ctx.moveTo(x, y); ctx.lineTo(waveX, state.path[0]); ctx.stroke();
ctx.globalAlpha = 1;
ctx.strokeStyle = accent2;
ctx.lineWidth = 2;
ctx.beginPath();
for (let i = 0; i < state.path.length; i++) {
  const xx = waveX + i;
  const yy = state.path[i];
  if (i === 0) ctx.moveTo(xx, yy); else ctx.lineTo(xx, yy);
}
ctx.stroke();
ctx.fillStyle = accent;
ctx.beginPath(); ctx.arc(x, y, 3.5, 0, 7); ctx.fill();
{{< /sketch >}}

Watch what each wheel does: the big one sets the basic back-and-forth, and each smaller, faster wheel sharpens the corners. Stop adding wheels and the corners stay slightly rounded and rippling — the unavoidable [[Gibbs Phenomenon]].

## See also

- [[Harmonics]]
- [[Gibbs Phenomenon]]
- [[Simple Harmonic Oscillator]]
