---
title: Continued Fraction
aliases: [continued fractions]
tags: [number-theory]
summary: A representation of a real number as a nested stack of fractions, yielding the best possible rational approximations.
weight: 80
---

# Continued Fraction

A **continued fraction** writes a number as an integer plus a fraction whose denominator is itself an integer plus a fraction, nested down as far as you like:

$$x = a_0 + \cfrac{1}{a_1 + \cfrac{1}{a_2 + \cfrac{1}{a_3 + \cdots}}}.$$

The list of integers $[a_0; a_1, a_2, a_3, \dots]$ *is* the number. This unusual notation turns out to give the **best rational approximations** any number can have — better, in a precise sense, than the decimal system ever could.

## How to build one

The algorithm is pure [[Greatest Common Divisor|Euclidean algorithm]] in disguise. Take the whole-number part, subtract it off, **flip** the remainder, and repeat:

$$\tfrac{45}{16} = 2 + \tfrac{13}{16} = 2 + \cfrac{1}{16/13} = 2 + \cfrac{1}{1 + \cfrac{3}{13}} = 2 + \cfrac{1}{1 + \cfrac{1}{4 + \cfrac{1}{3}}}.$$

So $\frac{45}{16} = [2; 1, 4, 3]$. The integers $2, 1, 4, 3$ are exactly the **quotients** Euclid's algorithm produces for $\gcd(45, 16)$ — the continued fraction and the gcd are the same computation read two ways. For a rational number the list always **terminates**; for an irrational it runs forever.

## Convergents: the best approximations

Truncating the list early gives a **convergent** — a rational that approximates $x$. The remarkable theorem is that each convergent is the *best possible* rational approximation for its size of denominator: no fraction with a smaller denominator gets closer.

{{< sketch height="360" caption="Convergents of π = [3; 7, 15, 1, 292, …] closing in on the true value. Each bar is one convergent p/q; its height is the error |p/q − π| on a log scale. Watch 22/7, then 333/106, then 355/113 — each a leap closer. Click to cycle through π, the golden ratio, and √2." >}}
if (frame === 0) {
  state.sets = [
    { name: 'π', cf: [3, 7, 15, 1, 292, 1, 1, 1, 2], val: Math.PI },
    { name: 'φ (golden ratio)', cf: [1, 1, 1, 1, 1, 1, 1, 1, 1, 1], val: (1 + Math.sqrt(5)) / 2 },
    { name: '√2', cf: [1, 2, 2, 2, 2, 2, 2, 2, 2], val: Math.SQRT2 }
  ];
  state.si = 0;
}
ctx.clearRect(0, 0, W, H);
const css = (v, f) => (getComputedStyle(document.documentElement).getPropertyValue(v).trim() || f);
const accent = css('--accent', '#5b9cff');
const accent2 = css('--accent-2', '#b07bff');
const good = css('--good', '#4ec98f');
const warn = css('--warn', '#e0a458');
const text = css('--text', '#e6e9ef');
const faint = css('--text-faint', 'rgba(230,233,239,0.5)');

if (mouse.clicked) state.si = (state.si + 1) % state.sets.length;
const S = state.sets[state.si];

// compute convergents p/q from the cf list
const conv = [];
let pm1 = 1, pm2 = 0, qm1 = 0, qm2 = 1;
for (let i = 0; i < S.cf.length; i++) {
  const a = S.cf[i];
  const p = a * pm1 + pm2, q = a * qm1 + qm2;
  conv.push({ p, q, err: Math.abs(p / q - S.val) });
  pm2 = pm1; pm1 = p; qm2 = qm1; qm1 = q;
}

const pad = 30;
const n = conv.length;
const bw = (W - 2 * pad) / n * 0.7;
const gap = (W - 2 * pad) / n;
// log-error scale
const logs = conv.map(c => Math.log10(Math.max(c.err, 1e-12)));
const lo = Math.min(...logs), hi = Math.max(...logs);
const baseY = H - 54, topY = 40;
ctx.textAlign = 'center'; ctx.textBaseline = 'bottom';
for (let i = 0; i < n; i++) {
  const x = pad + gap * i + gap / 2;
  const frac = (logs[i] - hi) / (lo - hi || 1);    // 0 (worst) .. 1 (best)
  const h = Math.max(3, frac * (baseY - topY));
  ctx.fillStyle = i === n - 1 ? warn : accent; ctx.globalAlpha = 0.85;
  ctx.fillRect(x - bw / 2, baseY - h, bw, h); ctx.globalAlpha = 1;
  ctx.fillStyle = text; ctx.font = '11px monospace';
  ctx.fillText(conv[i].p + '/' + conv[i].q, x, baseY - h - 3);
  ctx.fillStyle = faint; ctx.font = '9px monospace';
  ctx.save(); ctx.translate(x, baseY + 4); ctx.textBaseline = 'top';
  ctx.fillText('q=' + conv[i].q, 0, 0); ctx.restore();
}
ctx.fillStyle = text; ctx.font = '13px monospace'; ctx.textAlign = 'left';
ctx.textBaseline = 'top';
ctx.fillText('convergents of ' + S.name + '  (taller = closer)', pad, 12);
{{< /sketch >}}

For $\pi$, the convergents are $3,\ \frac{22}{7},\ \frac{333}{106},\ \frac{355}{113}, \dots$ The third, $\frac{355}{113}$, matches $\pi$ to six decimal places using a three-digit denominator — an approximation so good it was known in 5th-century China.

## The golden ratio: the worst-approximable number

The continued fraction also tells you which numbers are **hardest** to approximate by rationals. A large entry $a_i$ means the previous convergent was already excellent (you barely needed the next term). So the number that is *least* approximable is the one with the *smallest possible* entries everywhere — all $1$s:

$$\varphi = 1 + \cfrac{1}{1 + \cfrac{1}{1 + \cfrac{1}{1 + \cdots}}} = [1; 1, 1, 1, \dots] = \frac{1+\sqrt 5}{2}.$$

This is the **golden ratio**. Its all-ones expansion makes it the most stubbornly irrational number there is — a fact that explains why sunflower seeds and pinecones space their spirals by the golden angle: it is the packing least prone to settling into a repeating, gappy rational pattern.

{{< note kind="tip" title="Convergents alternate around the target" >}}
Successive convergents straddle the true value: one below, the next above, each closer than the last. And consecutive convergents $\frac{p_{k-1}}{q_{k-1}}, \frac{p_k}{q_k}$ always satisfy $p_k q_{k-1} - p_{k-1} q_k = \pm 1$ — a Bézout relation ([[Greatest Common Divisor]]) that guarantees each convergent is already in lowest terms.
{{< /note >}}

{{< quiz question="Why is the golden ratio φ = [1; 1, 1, 1, …] called the 'most irrational' number?" options="Because it is the largest irrational|Because its all-1s continued fraction makes its rational approximations converge as slowly as possible|Because it cannot be written as a continued fraction|Because it is transcendental" answer="2" explain="Large continued-fraction terms signal a sharp jump in approximation quality. The smallest possible terms — all 1s — give the slowest convergence, so no number resists rational approximation more than φ." >}}

## See also

- [[Greatest Common Divisor]]
- [[Diophantine Equation]]
- [[Number Theory]]
