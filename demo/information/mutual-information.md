---
title: Mutual Information
aliases: [mutual info]
tags: [information]
summary: How many bits knowing one variable tells you about another — the information they share.
weight: 50
---

# Mutual Information

**Mutual information** measures how much learning one random variable reduces your uncertainty about another. If $X$ and $Y$ are two variables, their mutual information $I(X;Y)$ is the number of bits that knowing $X$ saves you when you try to guess $Y$ (and, symmetrically, vice versa). It is the rigorous answer to "how related are these two things?" — sharper than correlation, because it captures *any* kind of dependence, not just linear.

In terms of [[Entropy]] it has a clean shape:

$$I(X;Y) = H(X) + H(Y) - H(X,Y) = H(Y) - H(Y\mid X).$$

The first form says: the shared information is what you'd double-count if you added the two uncertainties separately. The second says: it is how much $Y$'s uncertainty *drops*, from $H(Y)$ down to the leftover $H(Y\mid X)$, once $X$ is revealed.

{{< note kind="key" title="Independence ⇔ zero" >}}
$I(X;Y) \ge 0$ always, and $I(X;Y) = 0$ **exactly when $X$ and $Y$ are independent** — knowing one tells you nothing about the other. At the other extreme, if $Y$ is a deterministic copy of $X$, then $I(X;Y) = H(X)$: they share everything. Mutual information lives between these poles.
{{< /note >}}

## The overlap picture

Think of $H(X)$ and $H(Y)$ as two overlapping circles of uncertainty. Their **union** is the joint entropy $H(X,Y)$; their **intersection** is the mutual information. As two variables become more dependent, the circles slide together and the shared sliver grows. The sketch lets you drag them.

{{< sketch height="320" caption="Two circles of uncertainty. Their overlap is the mutual information I(X;Y); drag the mouse left↔right to vary the dependence from independent (no overlap) to identical (full overlap)." >}}
if (frame === 0) state.dep = 0.5;
ctx.clearRect(0, 0, W, H);
const css = (v, f) => (getComputedStyle(document.documentElement).getPropertyValue(v).trim() || f);
const accent = css('--accent', '#5b9cff');
const accent2 = css('--accent-2', '#b07bff');
const text = css('--text', '#e6e9ef');
const faint = css('--text-faint', 'rgba(230,233,239,0.55)');

if (mouse.x >= 0 && mouse.x <= W) state.dep = Math.min(1, Math.max(0, mouse.x / W));
const dep = state.dep; // 0 independent, 1 identical

const cy = H / 2 - 6;
const R = Math.min(W * 0.2, H * 0.32);
const maxSep = 2 * R;            // fully apart (tangent) when independent
const sep = maxSep * (1 - dep);  // distance between centers
const cxL = W / 2 - sep / 2;
const cxR = W / 2 + sep / 2;

// circles
ctx.globalCompositeOperation = 'lighter';
ctx.beginPath(); ctx.arc(cxL, cy, R, 0, 2 * Math.PI);
ctx.fillStyle = accent; ctx.globalAlpha = 0.42; ctx.fill();
ctx.beginPath(); ctx.arc(cxR, cy, R, 0, 2 * Math.PI);
ctx.fillStyle = accent2; ctx.globalAlpha = 0.42; ctx.fill();
ctx.globalCompositeOperation = 'source-over';
ctx.globalAlpha = 1;

ctx.lineWidth = 1.5;
ctx.strokeStyle = accent;
ctx.beginPath(); ctx.arc(cxL, cy, R, 0, 2 * Math.PI); ctx.stroke();
ctx.strokeStyle = accent2;
ctx.beginPath(); ctx.arc(cxR, cy, R, 0, 2 * Math.PI); ctx.stroke();

// labels
ctx.fillStyle = text;
ctx.font = '14px sans-serif';
ctx.fillText('H(X)', cxL - R * 0.7, cy - R * 0.4);
ctx.fillText('H(Y)', cxR + R * 0.2, cy - R * 0.4);

// lens area as a proxy for I(X;Y)
const d = sep;
let lens = 0;
if (d < 2 * R) {
  const part = Math.acos(d / (2 * R)) * R * R - (d / 4) * Math.sqrt(Math.max(0, 4 * R * R - d * d));
  lens = 2 * part;
}
const full = Math.PI * R * R;
const frac = Math.min(1, lens / full);

ctx.fillStyle = faint;
ctx.font = '13px sans-serif';
ctx.fillText('dependence: ' + (dep * 100).toFixed(0) + '%', 14, 22);
ctx.fillStyle = text;
ctx.fillText('I(X;Y) ∝ overlap = ' + (frac * 100).toFixed(0) + '% of one circle', 14, H - 16);
{{< /sketch >}}

## Where it matters

Mutual information is the quantity Shannon maximizes to define [[Channel Capacity]]: the capacity of a noisy channel is the largest $I(X;Y)$ achievable between its input $X$ and output $Y$, taken over all input distributions. It is also the workhorse of machine learning and statistics — feature selection, clustering quality, and registration of medical images all hinge on measuring how much one signal tells you about another. Because it sees *any* dependence, it catches relationships that a plain correlation coefficient, blind to nonlinearity, would miss.

{{< quiz question="Two variables X and Y are statistically independent. What is I(X;Y)?" options="Maximal|Exactly zero|Equal to H(X)|Negative" answer="2" explain="Mutual information is zero precisely when the variables are independent — knowing one reveals nothing about the other. It can never be negative." >}}

## See also

- [[Entropy]]
- [[Channel Capacity]]
- [[Probability Distribution]]
