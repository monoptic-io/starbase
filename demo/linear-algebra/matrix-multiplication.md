---
title: Matrix Multiplication
tags: [linear-algebra]
summary: Multiplying matrices composes their transformations — do one, then the other — and the order matters.
weight: 40
---

# Matrix Multiplication

To multiply two matrices is to **compose their transformations**: $AB$ is the single [[Linear Transformation]] you get by first applying $B$, then applying $A$. The strange-looking arithmetic — each entry of the product is a row of $A$ dotted with a column of $B$ — is exactly what is forced on you once you demand that $(AB)\vec v = A(B\vec v)$ for every [[Vector]] $\vec v$.

## The rule, and why it looks that way

For $2\times2$ matrices,

{{< eq number="1" >}}
\begin{bmatrix} a & b \\ c & d \end{bmatrix}
\begin{bmatrix} e & f \\ g & h \end{bmatrix}
=
\begin{bmatrix} ae+bg & af+bh \\ ce+dg & cf+dh \end{bmatrix}.
{{< /eq >}}

The clean way to see it: the columns of $B$ are where $B$ sends the basis vectors. Pushing each of those columns through $A$ gives the columns of $AB$. So the product's first column is $A$ acting on $B$'s first column — composition, column by column. For the shapes to line up, the number of columns of $A$ must equal the number of rows of $B$.

## Order matters

Composition is generally **non-commutative**: $AB \neq BA$. Rotating a shape and *then* stretching it is not the same as stretching and then rotating. The picture below applies a rotation $R$ and a horizontal shear $S$ to the same letter in both orders; the results disagree.

{{< sketch height="360" caption="The same F under RS (rotate then shear, blue) versus SR (shear then rotate, amber). Different orders, different results — matrix multiplication does not commute." >}}
if (frame === 0) { state.th = 0; }
const css = getComputedStyle(document.documentElement);
const accent = (css.getPropertyValue('--accent') || '#5b9cff').trim();
const accent2 = (css.getPropertyValue('--accent-2') || '#f59e0b').trim();
const faint = (css.getPropertyValue('--text-faint') || '#7a8390').trim();
const text = (css.getPropertyValue('--text') || '#e6e6e6').trim();

state.th += dt * 0.5;
const a = 0.6 + 0.45 * Math.sin(state.th);   // animated rotation angle
const k = 1.1;                                // shear strength
const ct = Math.cos(a), st = Math.sin(a);
// R = rotation, Sh = horizontal shear
function mul(M, N) {
  return [[M[0][0]*N[0][0]+M[0][1]*N[1][0], M[0][0]*N[0][1]+M[0][1]*N[1][1]],
          [M[1][0]*N[0][0]+M[1][1]*N[1][0], M[1][0]*N[0][1]+M[1][1]*N[1][1]]];
}
const R = [[ct, -st], [st, ct]];
const Sh = [[1, k], [0, 1]];
const RS = mul(R, Sh);   // first shear, then rotate
const SR = mul(Sh, R);   // first rotate, then shear

const F = [
  [-0.20, -0.45], [-0.20, 0.45], [0.25, 0.45],
  [0.25, 0.30], [-0.05, 0.30], [-0.05, 0.12],
  [0.18, 0.12], [0.18, -0.03], [-0.05, -0.03], [-0.05, -0.45]
];

function drawF(M, cx, col, label) {
  const u = Math.min(W, H) / 4.2;
  ctx.strokeStyle = faint; ctx.globalAlpha = 0.3; ctx.lineWidth = 1;
  ctx.beginPath(); ctx.moveTo(cx, 20); ctx.lineTo(cx, H - 28);
  ctx.moveTo(cx - u, H / 2 - 8); ctx.lineTo(cx + u, H / 2 - 8); ctx.stroke();
  ctx.globalAlpha = 1;
  ctx.fillStyle = col; ctx.globalAlpha = 0.85;
  ctx.beginPath();
  for (let i = 0; i < F.length; i++) {
    const x = M[0][0]*F[i][0] + M[0][1]*F[i][1];
    const y = M[1][0]*F[i][0] + M[1][1]*F[i][1];
    const sx = cx + x * u, sy = (H / 2 - 8) - y * u;
    if (i === 0) ctx.moveTo(sx, sy); else ctx.lineTo(sx, sy);
  }
  ctx.closePath(); ctx.fill();
  ctx.globalAlpha = 1;
  ctx.fillStyle = text; ctx.font = '14px system-ui, sans-serif';
  ctx.textAlign = 'center';
  ctx.fillText(label, cx, H - 8);
  ctx.textAlign = 'left';
}

ctx.clearRect(0, 0, W, H);
drawF(RS, W * 0.28, accent, 'R S  (shear, then rotate)');
drawF(SR, W * 0.72, accent2, 'S R  (rotate, then shear)');
{{< /sketch >}}

## Associative, not commutative

{{< note kind="note" title="What does carry over" >}}
Multiplication *is* associative — $(AB)C = A(BC)$ — because composing functions is associative; you can regroup a chain of transformations freely. It just isn't commutative, so you may never reorder them. The identity matrix is the do-nothing transformation, and a matrix has an inverse exactly when its [[Determinant]] is nonzero, since only then is the warp reversible.
{{< /note >}}

## See also

- [[Linear Transformation]]
- [[Determinant]]
- [[Matrix]]
