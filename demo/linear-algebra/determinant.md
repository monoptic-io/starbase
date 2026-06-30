---
title: Determinant
tags: [linear-algebra]
summary: A single number that says how much a transformation scales area or volume — and whether it flips orientation.
weight: 50
---

# Determinant

The **determinant** of a square [[Matrix]] is one number that summarizes the most basic effect of its [[Linear Transformation]]: **how much it scales area** (in 2D) or **volume** (in 3D). The unit square has area $1$; after the transformation it becomes a parallelogram, and its new area *is* the determinant. A sign rides along too — a negative determinant means the transformation turned space inside-out, a mirror flip.

## Area of the parallelogram

For a $2\times2$ matrix the formula is short:

{{< eq number="1" >}}
\det\begin{bmatrix} a & b \\ c & d \end{bmatrix} = ad - bc.
{{< /eq >}}

The two columns $(a,c)$ and $(b,d)$ are the edges of the image parallelogram, and $ad - bc$ is its signed area. The sketch sends the unit square through a slowly changing matrix; watch the readout track the parallelogram's area exactly, and watch it pass through **zero** at the instant the shape collapses to a line.

{{< sketch height="380" caption="The unit square (faint) mapped to a parallelogram (filled). Its signed area equals the determinant — which crosses zero, and goes negative, as the shape collapses and flips." >}}
if (frame === 0) { state.t = 0; }
const css = getComputedStyle(document.documentElement);
const accent = (css.getPropertyValue('--accent') || '#5b9cff').trim();
const faint = (css.getPropertyValue('--text-faint') || '#7a8390').trim();
const good = (css.getPropertyValue('--good') || '#34d399').trim();
const warn = (css.getPropertyValue('--warn') || '#f87171').trim();
const text = (css.getPropertyValue('--text') || '#e6e6e6').trim();

state.t += dt * 0.6;
// animated matrix columns
const a = 1.3, b = 0.9 * Math.sin(state.t);
const c = 0.5 * Math.sin(state.t * 0.7), d = 1.1;
const det = a * d - b * c;

const unit = Math.min(W, H) / 5;
const ox = W * 0.42, oy = H / 2;
function S(p, q) { return { x: ox + (p * a + q * b) * unit, y: oy - (p * c + q * d) * unit }; }

ctx.clearRect(0, 0, W, H);
// axes
ctx.strokeStyle = faint; ctx.globalAlpha = 0.3; ctx.lineWidth = 1;
ctx.beginPath(); ctx.moveTo(0, oy); ctx.lineTo(W, oy);
ctx.moveTo(ox, 0); ctx.lineTo(ox, H); ctx.stroke();

// original unit square (faint outline)
ctx.globalAlpha = 0.4; ctx.setLineDash([4, 4]);
ctx.strokeStyle = faint;
ctx.strokeRect(ox, oy - unit, unit, unit);
ctx.setLineDash([]); ctx.globalAlpha = 1;

// transformed parallelogram
const q00 = S(0, 0), q10 = S(1, 0), q11 = S(1, 1), q01 = S(0, 1);
ctx.fillStyle = (det < 0) ? warn : accent;
ctx.globalAlpha = 0.22;
ctx.beginPath();
ctx.moveTo(q00.x, q00.y); ctx.lineTo(q10.x, q10.y);
ctx.lineTo(q11.x, q11.y); ctx.lineTo(q01.x, q01.y); ctx.closePath(); ctx.fill();
ctx.globalAlpha = 1;
ctx.strokeStyle = (det < 0) ? warn : accent; ctx.lineWidth = 2;
ctx.stroke();

// edge vectors
ctx.lineWidth = 2.5;
ctx.strokeStyle = accent;
ctx.beginPath(); ctx.moveTo(q00.x, q00.y); ctx.lineTo(q10.x, q10.y); ctx.stroke();
ctx.strokeStyle = good;
ctx.beginPath(); ctx.moveTo(q00.x, q00.y); ctx.lineTo(q01.x, q01.y); ctx.stroke();

// readout
ctx.font = '14px ui-monospace, monospace';
ctx.fillStyle = text;
ctx.fillText('area scale = ' + Math.abs(det).toFixed(2), 14, 24);
ctx.fillStyle = (det < 0) ? warn : good;
ctx.fillText('det = ' + det.toFixed(2) + (det < 0 ? '  (flipped)' : ''), 14, 44);
{{< /sketch >}}

## When the determinant is zero

{{< note kind="key" title="Collapse and irreversibility" >}}
A determinant of **zero** means the transformation squashes area to nothing: the two columns have become parallel, and the whole plane is crushed onto a line (or a point). Information is lost and the map cannot be undone — a matrix is invertible **iff** its determinant is nonzero. A zero determinant also signals that the matrix has a zero [[Eigenvalues and Eigenvectors|eigenvalue]], with the crushed direction as its eigenvector.
{{< /note >}}

## Composition multiplies

Determinants behave beautifully under [[Matrix Multiplication]]: $\det(AB) = \det(A)\det(B)$. Two transformations that each double area combine to quadruple it — scale factors simply multiply, which is exactly what areas should do when you apply one map after another.

## See also

- [[Linear Transformation]]
- [[Eigenvalues and Eigenvectors]]
- [[Matrix Multiplication]]
