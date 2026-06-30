---
title: Linear Transformation
aliases: [linear map]
tags: [linear-algebra]
summary: A way of warping space that keeps grid lines straight, evenly spaced, and pinned to a fixed origin — exactly what a matrix encodes.
weight: 30
---

# Linear Transformation

A **linear transformation** is a function $T$ that takes [[Vector|vectors]] to vectors while preserving the two operations that define a vector space:

$$T(\vec u + \vec v) = T(\vec u) + T(\vec v), \qquad T(c\,\vec v) = c\,T(\vec v).$$

Geometrically these algebraic rules have a vivid consequence: **grid lines stay straight and evenly spaced, and the origin never moves.** Space can rotate, stretch, shear, or flip — but it cannot bend or tear. Every such transformation of the plane is captured by a single $2\times2$ [[Matrix]], whose two columns are simply where the transformation sends the basis vectors $(1,0)$ and $(0,1)$.

## Drag the basis vectors

Below is the identity grid warped by a live matrix. The two arrows are the matrix's **columns** — the landing spots of the $x$- and $y$-axes. Grab either arrowhead and drag: the entire grid, and the shaded unit square with its **F**, deform to follow. Notice what *cannot* happen — lines never curve, parallel lines stay parallel, and the origin stays put. That is linearity, made visible.

{{< sketch height="440" caption="Drag the blue or amber arrowhead to set the matrix columns. The grid, unit square, and F warp accordingly; the origin stays fixed. Watch the determinant flip sign when you turn the square inside-out." >}}
if (frame === 0) {
  state.c1 = { x: 1.0, y: 0.25 };
  state.c2 = { x: -0.35, y: 1.0 };
  state.drag = null;
  state.wasDown = false;
}
const css = getComputedStyle(document.documentElement);
const accent = (css.getPropertyValue('--accent') || '#5b9cff').trim();
const accent2 = (css.getPropertyValue('--accent-2') || '#f59e0b').trim();
const faint = (css.getPropertyValue('--text-faint') || '#7a8390').trim();
const good = (css.getPropertyValue('--good') || '#34d399').trim();
const warn = (css.getPropertyValue('--warn') || '#f87171').trim();
const text = (css.getPropertyValue('--text') || '#e6e6e6').trim();

const c1 = state.c1, c2 = state.c2;
const unit = Math.min(W, H) / 6.2;
const ox = W / 2, oy = H / 2;
const N = 3;
function S(p, q) {
  return { x: ox + (p * c1.x + q * c2.x) * unit, y: oy - (p * c1.y + q * c2.y) * unit };
}

// --- dragging ---
const justPressed = mouse.down && !state.wasDown;
if (justPressed) {
  const h1 = S(1, 0), h2 = S(0, 1);
  const d1 = Math.hypot(mouse.x - h1.x, mouse.y - h1.y);
  const d2 = Math.hypot(mouse.x - h2.x, mouse.y - h2.y);
  state.drag = (d1 < d2) ? 1 : 2;
}
if (!mouse.down) state.drag = null;
if (mouse.down && state.drag) {
  let mx = (mouse.x - ox) / unit, my = -(mouse.y - oy) / unit;
  mx = Math.max(-N, Math.min(N, mx));
  my = Math.max(-N, Math.min(N, my));
  if (state.drag === 1) { c1.x = mx; c1.y = my; } else { c2.x = mx; c2.y = my; }
}
state.wasDown = mouse.down;

ctx.clearRect(0, 0, W, H);

// --- warped grid ---
ctx.lineWidth = 1;
for (let i = -N; i <= N; i++) {
  const main = (i === 0);
  ctx.strokeStyle = faint;
  ctx.globalAlpha = main ? 0.7 : 0.3;
  let a = S(i, -N), b = S(i, N);
  ctx.beginPath(); ctx.moveTo(a.x, a.y); ctx.lineTo(b.x, b.y); ctx.stroke();
  let c = S(-N, i), d = S(N, i);
  ctx.beginPath(); ctx.moveTo(c.x, c.y); ctx.lineTo(d.x, d.y); ctx.stroke();
}
ctx.globalAlpha = 1;

// --- transformed unit square ---
const det = c1.x * c2.y - c2.x * c1.y;
const p00 = S(0, 0), p10 = S(1, 0), p11 = S(1, 1), p01 = S(0, 1);
ctx.fillStyle = (det < 0) ? warn : accent;
ctx.globalAlpha = 0.16;
ctx.beginPath();
ctx.moveTo(p00.x, p00.y); ctx.lineTo(p10.x, p10.y);
ctx.lineTo(p11.x, p11.y); ctx.lineTo(p01.x, p01.y); ctx.closePath();
ctx.fill();
ctx.globalAlpha = 1;

// --- the F, to expose flips and shears ---
const Fpts = [
  [0.30, 0.15], [0.30, 0.85], [0.70, 0.85],
  [0.70, 0.70], [0.45, 0.70], [0.45, 0.55],
  [0.62, 0.55], [0.62, 0.40], [0.45, 0.40], [0.45, 0.15]
];
ctx.fillStyle = (det < 0) ? warn : accent;
ctx.globalAlpha = 0.85;
ctx.beginPath();
for (let i = 0; i < Fpts.length; i++) {
  const s = S(Fpts[i][0], Fpts[i][1]);
  if (i === 0) ctx.moveTo(s.x, s.y); else ctx.lineTo(s.x, s.y);
}
ctx.closePath(); ctx.fill();
ctx.globalAlpha = 1;

// --- basis arrows ---
function arrow(tip, col) {
  ctx.strokeStyle = col; ctx.fillStyle = col; ctx.lineWidth = 3;
  ctx.beginPath(); ctx.moveTo(ox, oy); ctx.lineTo(tip.x, tip.y); ctx.stroke();
  const ang = Math.atan2(tip.y - oy, tip.x - ox);
  ctx.beginPath();
  ctx.moveTo(tip.x, tip.y);
  ctx.lineTo(tip.x - 13 * Math.cos(ang - 0.4), tip.y - 13 * Math.sin(ang - 0.4));
  ctx.lineTo(tip.x - 13 * Math.cos(ang + 0.4), tip.y - 13 * Math.sin(ang + 0.4));
  ctx.closePath(); ctx.fill();
}
const t1 = S(1, 0), t2 = S(0, 1);
arrow(t1, accent);
arrow(t2, accent2);
// draggable handles
ctx.fillStyle = accent;
ctx.beginPath(); ctx.arc(t1.x, t1.y, state.drag === 1 ? 8 : 6, 0, 7); ctx.fill();
ctx.fillStyle = accent2;
ctx.beginPath(); ctx.arc(t2.x, t2.y, state.drag === 2 ? 8 : 6, 0, 7); ctx.fill();

// --- readout ---
ctx.font = '13px ui-monospace, monospace';
ctx.fillStyle = text;
ctx.fillText('A = [ ' + c1.x.toFixed(2) + '  ' + c2.x.toFixed(2) + ' ]', 14, 22);
ctx.fillText('    [ ' + c1.y.toFixed(2) + '  ' + c2.y.toFixed(2) + ' ]', 14, 40);
ctx.fillStyle = (det < 0) ? warn : good;
ctx.fillText('det A = ' + det.toFixed(2), 14, 62);
{{< /sketch >}}

## Why columns are everything

Because $T$ respects linear combinations, knowing what it does to the basis vectors fixes what it does to *all* vectors. If $\vec v = x\,\hat\imath + y\,\hat\jmath$, then $T(\vec v) = x\,T(\hat\imath) + y\,T(\hat\jmath)$ — and $T(\hat\imath), T(\hat\jmath)$ are exactly the columns of the matrix. This is why "apply a matrix" and "apply a linear transformation" mean the same thing, and why composing transformations is just [[Matrix Multiplication]].

## Reading the geometry

{{< note kind="key" title="Two numbers you can see" >}}
The shaded square's signed area is the [[Determinant]]: it reports how much the transformation scales area, and it turns negative (here, red) when the map flips orientation — when the **F** becomes its mirror image. And the special directions that the warp only stretches, never rotates, are the [[Eigenvalues and Eigenvectors|eigenvectors]]. Both are properties of the same matrix you are dragging.
{{< /note >}}

## See also

- [[Matrix]]
- [[Eigenvalues and Eigenvectors]]
- [[Determinant]]
