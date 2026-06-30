---
title: Eigenvalues and Eigenvectors
aliases: [eigenvalue, eigenvector, eigenvectors, eigen]
tags: [linear-algebra]
summary: The special directions a transformation merely stretches — and the factors by which it stretches them — the hidden skeleton of a matrix.
weight: 70
---

# Eigenvalues and Eigenvectors

For most vectors, a [[Linear Transformation]] both rotates and stretches. But almost every matrix has a few special directions it leaves *pointing the same way*, only scaled. A nonzero vector $\vec v$ with

$$A\vec v = \lambda\,\vec v$$

is an **eigenvector** of $A$, and the scalar $\lambda$ is its **eigenvalue** — the factor by which that direction is stretched ($|\lambda| > 1$), shrunk ($|\lambda| < 1$), or flipped ($\lambda < 0$). Eigenvectors are the axes along which the transformation is purely a scaling, and finding them strips a complicated matrix down to its skeleton.

## Finding them

The defining equation rearranges to $(A - \lambda I)\vec v = \vec 0$, which has a nonzero solution only when the matrix $A - \lambda I$ collapses space — that is, when its [[Determinant]] vanishes:

{{< eq number="1" >}}
\det(A - \lambda I) = 0.
{{< /eq >}}

This **characteristic equation** is a polynomial in $\lambda$; its roots are the eigenvalues, and back-substituting each one yields its eigenvector direction. For a symmetric matrix like $A=\begin{bmatrix}2&1\\1&2\end{bmatrix}$, the eigenvalues are $3$ and $1$, with perpendicular eigenvectors $(1,1)$ and $(1,-1)$.

## Watch the dominant direction win

Apply a matrix over and over and a remarkable thing happens: almost any starting vector swings toward the eigenvector with the **largest** eigenvalue, because that direction grows fastest and drowns the others out. Below, a ring of vectors is repeatedly hit with $A=\begin{bmatrix}2&1\\1&2\end{bmatrix}$ (and renormalized each step so they stay on screen). Watch them sweep onto the dominant eigenvector $(1,1)$, the green line.

{{< sketch height="420" caption="A ring of directions, repeatedly transformed by A and renormalized. They collapse onto the dominant eigenvector (green); the other eigenvector (faint) is repulsive. Click to release a fresh ring." >}}
if (frame === 0) {
  state.A = [[2, 1], [1, 2]];
  state.vecs = [];
  state.step = 0; state.timer = 0; state.wasDown = false;
  state.reset = true;
}
const css = getComputedStyle(document.documentElement);
const accent = (css.getPropertyValue('--accent') || '#5b9cff').trim();
const good = (css.getPropertyValue('--good') || '#34d399').trim();
const faint = (css.getPropertyValue('--text-faint') || '#7a8390').trim();
const text = (css.getPropertyValue('--text') || '#e6e6e6').trim();

function initRing() {
  state.vecs = [];
  const M = 36;
  for (let i = 0; i < M; i++) {
    const a = (i / M) * Math.PI * 2;
    state.vecs.push({ x: Math.cos(a), y: Math.sin(a) });
  }
  state.step = 0;
}
if (state.reset) { initRing(); state.reset = false; }

// click to relaunch
const justPressed = mouse.down && !state.wasDown;
if (justPressed) initRing();
state.wasDown = mouse.down;

const A = state.A;
// advance one transform step on a timer; after converging, hold, then relaunch
const STEPS = 11, HOLD = 5;  // HOLD extra ticks of pause on the converged ring
state.timer += dt;
if (state.timer > 0.45) {
  state.timer -= 0.45;
  if (state.step < STEPS) {
    for (let i = 0; i < state.vecs.length; i++) {
      const v = state.vecs[i];
      const nx = A[0][0] * v.x + A[0][1] * v.y;
      const ny = A[1][0] * v.x + A[1][1] * v.y;
      const L = Math.hypot(nx, ny) || 1;
      v.x = nx / L; v.y = ny / L;
    }
    state.step++;
  } else if (state.step < STEPS + HOLD) {
    state.step++;                 // pause on the converged ring
  } else {
    initRing();                   // relaunch a fresh ring
  }
}

const ox = W / 2, oy = H / 2;
const R = Math.min(W, H) * 0.40;

ctx.clearRect(0, 0, W, H);
// unit circle
ctx.strokeStyle = faint; ctx.globalAlpha = 0.3; ctx.lineWidth = 1;
ctx.beginPath(); ctx.arc(ox, oy, R, 0, 7); ctx.stroke();
ctx.globalAlpha = 1;

// eigenvector reference lines: dominant (1,1), other (1,-1)
function evLine(dx, dy, col, w) {
  const L = Math.hypot(dx, dy);
  const ux = dx / L, uy = dy / L;
  ctx.strokeStyle = col; ctx.lineWidth = w;
  ctx.beginPath();
  ctx.moveTo(ox - ux * R, oy + uy * R);
  ctx.lineTo(ox + ux * R, oy - uy * R);
  ctx.stroke();
}
ctx.globalAlpha = 0.5; evLine(1, -1, faint, 1.5);   // subdominant, screen-y flipped
ctx.globalAlpha = 0.9; evLine(1, 1, good, 2.5);      // dominant
ctx.globalAlpha = 1;

// the transforming ring of directions
for (let i = 0; i < state.vecs.length; i++) {
  const v = state.vecs[i];
  const tx = ox + v.x * R, ty = oy - v.y * R;
  // alignment with dominant eigenvector (1,1)/sqrt2
  const al = Math.abs((v.x + v.y) / Math.SQRT2);
  ctx.globalAlpha = 0.35 + 0.5 * al;
  ctx.strokeStyle = accent; ctx.lineWidth = 1.5;
  ctx.beginPath(); ctx.moveTo(ox, oy); ctx.lineTo(tx, ty); ctx.stroke();
  ctx.fillStyle = al > 0.985 ? good : accent;
  ctx.beginPath(); ctx.arc(tx, ty, 3, 0, 7); ctx.fill();
}
ctx.globalAlpha = 1;

ctx.font = '13px ui-monospace, monospace';
ctx.fillStyle = good;
ctx.fillText('dominant eigenvector (1,1), λ = 3', 14, 24);
ctx.fillStyle = text;
ctx.fillText('iteration ' + Math.min(state.step, 11), 14, 44);
{{< /sketch >}}

This iterate-and-normalize loop is the **power method**, and it is not a toy: it is precisely how [[PageRank]] finds the dominant eigenvector of a web's link [[Matrix]] to score every page at once.

## Why eigenvalues run the world

{{< note kind="key" title="The thread through the field guide" >}}
Eigenvalues decide the long-term fate of linear processes, so they surface everywhere a system is linearized:

- **[[Stability]]** — the eigenvalues of the linearization at a [[Fixed Point]] decide whether perturbations grow or decay.
- **[[Coupled Oscillators]]** — the *normal modes* are eigenvectors of the coupling matrix, each ringing at its own frequency.
- **[[PageRank]]** — importance scores are the dominant eigenvector of the link matrix.

Whenever a system's behavior is "the same thing, scaled, repeated," eigenvalues are the right language.
{{< /note >}}

## See also

- [[Linear Transformation]]
- [[Singular Value Decomposition]]
- [[Stability]]
