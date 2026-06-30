---
title: Singular Value Decomposition
aliases: [SVD]
tags: [linear-algebra]
summary: Every matrix factors into a rotation, a pure scaling, and another rotation — the most useful decomposition in applied linear algebra.
weight: 80
---

# Singular Value Decomposition

The **singular value decomposition** says something almost too clean to believe: *every* [[Matrix]] $A$, of any shape, can be written as

{{< eq number="1" >}}
A = U\,\Sigma\,V^{\top},
{{< /eq >}}

where $U$ and $V$ are **rotations** (orthogonal matrices) and $\Sigma$ is a **diagonal scaling**. In words, any [[Linear Transformation]] — however much it seems to shear and skew — is really just *rotate, stretch along perpendicular axes, rotate again*. The stretch factors on $\Sigma$'s diagonal are the **singular values** $\sigma_1 \ge \sigma_2 \ge \cdots \ge 0$, and they measure how much the map amplifies space along each of its principal directions.

## Rotate, scale, rotate

Geometrically, $V^{\top}$ first turns the input so that the transformation's natural axes line up with the coordinate axes; $\Sigma$ stretches each of those axes by its singular value; and $U$ rotates the result into its final orientation. The image of the unit circle under any matrix is an **ellipse**, and the SVD hands you that ellipse directly: its semi-axis lengths are the singular values, and their directions are the columns of $U$.

{{< sketch height="360" caption="The unit circle (faint) mapped by a matrix into an ellipse. The ellipse's perpendicular semi-axes are the singular values σ₁ ≥ σ₂ — the rotate-scale-rotate skeleton of the transformation." >}}
if (frame === 0) { state.t = 0; }
const css = getComputedStyle(document.documentElement);
const accent = (css.getPropertyValue('--accent') || '#5b9cff').trim();
const accent2 = (css.getPropertyValue('--accent-2') || '#f59e0b').trim();
const faint = (css.getPropertyValue('--text-faint') || '#7a8390').trim();
const good = (css.getPropertyValue('--good') || '#34d399').trim();
const text = (css.getPropertyValue('--text') || '#e6e6e6').trim();

state.t += dt * 0.5;
// animated matrix
const a = 1.4, b = 0.8 * Math.sin(state.t);
const c = 0.6 * Math.sin(state.t * 0.6), d = 1.0;

const ox = W / 2, oy = H / 2;
const u = Math.min(W, H) / 4.6;

ctx.clearRect(0, 0, W, H);
ctx.strokeStyle = faint; ctx.globalAlpha = 0.3; ctx.lineWidth = 1;
ctx.beginPath(); ctx.moveTo(0, oy); ctx.lineTo(W, oy);
ctx.moveTo(ox, 0); ctx.lineTo(ox, H); ctx.stroke();

// original unit circle
ctx.setLineDash([4, 4]); ctx.strokeStyle = faint; ctx.globalAlpha = 0.6;
ctx.beginPath(); ctx.arc(ox, oy, u, 0, 7); ctx.stroke();
ctx.setLineDash([]); ctx.globalAlpha = 1;

// mapped ellipse
ctx.strokeStyle = accent; ctx.lineWidth = 2.5; ctx.beginPath();
for (let i = 0; i <= 64; i++) {
  const th = i / 64 * Math.PI * 2;
  const x = Math.cos(th), y = Math.sin(th);
  const px = a * x + b * y, py = c * x + d * y;
  const sx = ox + px * u, sy = oy - py * u;
  if (i === 0) ctx.moveTo(sx, sy); else ctx.lineTo(sx, sy);
}
ctx.closePath(); ctx.stroke();

// singular values via SVD of [[a,b],[c,d]] (2x2 closed form)
const E = (a + d) / 2, F = (a - d) / 2, G = (c + b) / 2, Hh = (c - b) / 2;
const Q = Math.hypot(E, Hh), Rr = Math.hypot(F, G);
const s1 = Q + Rr, s2 = Math.abs(Q - Rr);
const a1 = Math.atan2(Hh, E), a2 = Math.atan2(G, F);
const theta = (a2 - a1) / 2;        // angle of left singular vectors (U)
// draw the two principal semi-axes (U directions, scaled by singular values)
function axis(ang, len, col) {
  const dx = Math.cos(ang) * len * u, dy = Math.sin(ang) * len * u;
  ctx.strokeStyle = col; ctx.lineWidth = 3;
  ctx.beginPath(); ctx.moveTo(ox, oy); ctx.lineTo(ox + dx, oy - dy); ctx.stroke();
}
axis(theta, s1, good);
axis(theta + Math.PI / 2, s2, accent2);

ctx.font = '13px ui-monospace, monospace';
ctx.fillStyle = good;
ctx.fillText('σ₁ = ' + s1.toFixed(2), 14, 24);
ctx.fillStyle = accent2;
ctx.fillText('σ₂ = ' + s2.toFixed(2), 14, 44);
{{< /sketch >}}

## Low-rank approximation

{{< note kind="key" title="Keep the big singular values, drop the rest" >}}
Because the singular values are sorted, the first few capture most of what a matrix does. Keep only the largest $k$ — zeroing the rest — and you get the **best possible rank-$k$ approximation** of $A$. This is the heart of data compression and dimensionality reduction: an image, a recommendation table, or a term–document matrix is often *almost* low-rank, so a handful of singular values reconstruct it faithfully while storing a fraction of the numbers.
{{< /note >}}

## Relation to eigenvalues

The SVD is the well-behaved cousin of [[Eigenvalues and Eigenvectors|eigendecomposition]]. Eigenvectors need not be perpendicular and exist cleanly only for square matrices, but singular values are always real, non-negative, and come with orthogonal axes — for *any* matrix at all. In fact the singular values of $A$ are the square roots of the eigenvalues of $A^{\top}A$, tying the two ideas together.

## See also

- [[Eigenvalues and Eigenvectors]]
- [[Matrix]]
- [[Linear Transformation]]
