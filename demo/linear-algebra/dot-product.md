---
title: Dot Product
aliases: [inner product, projection]
tags: [linear-algebra]
summary: A single number, built from two vectors, that encodes the angle between them and how much one casts onto the other.
weight: 55
---

# Dot Product

The **dot product** of two [[Vector|vectors]] multiplies them into a single number that measures how much they point the same way. Two equivalent formulas define it:

{{< eq number="1" >}}
\vec a \cdot \vec b = a_x b_x + a_y b_y = \lVert \vec a\rVert\,\lVert \vec b\rVert \cos\theta.
{{< /eq >}}

The left form is pure arithmetic — multiply matching components and add. The right form reveals the geometry: it is the product of the two lengths times the cosine of the angle $\theta$ between them. Setting the two equal lets you *recover an angle from coordinates*, which is the seed of nearly all geometry done with numbers.

## Projection

The dot product is the engine of **projection** — dropping one vector's shadow onto another's direction. The component of $\vec a$ along a unit vector $\hat b$ is just $\vec a \cdot \hat b$, and the projected vector is

$$\operatorname{proj}_{\vec b}\,\vec a = \frac{\vec a \cdot \vec b}{\vec b \cdot \vec b}\,\vec b.$$

Drag the tip of $\vec a$ below and watch its shadow slide along $\vec b$. The projection is longest when the vectors align and shrinks to nothing when they are perpendicular.

{{< sketch height="360" caption="Drag to move vector a (blue). Its projection onto b (amber) is shown as the green shadow; the dashed line is the perpendicular dropped from a's tip." >}}
if (frame === 0) { state.a = { x: 120, y: -70 }; }
const css = getComputedStyle(document.documentElement);
const accent = (css.getPropertyValue('--accent') || '#5b9cff').trim();
const accent2 = (css.getPropertyValue('--accent-2') || '#f59e0b').trim();
const good = (css.getPropertyValue('--good') || '#34d399').trim();
const faint = (css.getPropertyValue('--text-faint') || '#7a8390').trim();
const text = (css.getPropertyValue('--text') || '#e6e6e6').trim();

const ox = W / 2, oy = H / 2;
if (mouse.down) { state.a.x = mouse.x - ox; state.a.y = mouse.y - oy; }
const a = state.a;
const b = { x: W * 0.30, y: 0 };  // b points along +x

ctx.clearRect(0, 0, W, H);
ctx.strokeStyle = faint; ctx.globalAlpha = 0.3; ctx.lineWidth = 1;
ctx.beginPath(); ctx.moveTo(0, oy); ctx.lineTo(W, oy);
ctx.moveTo(ox, 0); ctx.lineTo(ox, H); ctx.stroke();
ctx.globalAlpha = 1;

// projection scalar t along b
const bb = b.x * b.x + b.y * b.y;
const t = (a.x * b.x + a.y * b.y) / bb;
const proj = { x: t * b.x, y: t * b.y };

function arrow(vx, vy, col, lw) {
  ctx.strokeStyle = col; ctx.fillStyle = col; ctx.lineWidth = lw;
  ctx.beginPath(); ctx.moveTo(ox, oy); ctx.lineTo(ox + vx, oy + vy); ctx.stroke();
  const ang = Math.atan2(vy, vx);
  ctx.beginPath();
  ctx.moveTo(ox + vx, oy + vy);
  ctx.lineTo(ox + vx - 12 * Math.cos(ang - 0.4), oy + vy - 12 * Math.sin(ang - 0.4));
  ctx.lineTo(ox + vx - 12 * Math.cos(ang + 0.4), oy + vy - 12 * Math.sin(ang + 0.4));
  ctx.closePath(); ctx.fill();
}

// perpendicular from a's tip to the line of b
ctx.strokeStyle = faint; ctx.setLineDash([5, 5]); ctx.lineWidth = 1.5;
ctx.beginPath(); ctx.moveTo(ox + a.x, oy + a.y); ctx.lineTo(ox + proj.x, oy + proj.y); ctx.stroke();
ctx.setLineDash([]);

arrow(b.x, b.y, accent2, 2.5);
arrow(a.x, a.y, accent, 3);
// projection shadow (thick, on top of b's line)
ctx.strokeStyle = good; ctx.lineWidth = 5;
ctx.beginPath(); ctx.moveTo(ox, oy); ctx.lineTo(ox + proj.x, oy + proj.y); ctx.stroke();

// readout
const la = Math.sqrt(a.x * a.x + a.y * a.y);
const lb = Math.sqrt(bb);
const dot = a.x * b.x + a.y * b.y;
const cos = la > 0 ? dot / (la * lb) : 0;
ctx.font = '13px ui-monospace, monospace';
ctx.fillStyle = text;
ctx.fillText('a · b ' + (dot >= 0 ? '> 0' : '< 0') + '   cos θ = ' + cos.toFixed(2), 14, 22);
ctx.fillStyle = good;
ctx.fillText(Math.abs(cos) < 0.04 ? 'a ⟂ b  (a · b = 0)' : '', 14, 42);
{{< /sketch >}}

## Orthogonality

{{< note kind="key" title="Zero means perpendicular" >}}
When $\vec a \cdot \vec b = 0$ the vectors are **orthogonal** — at right angles — because $\cos 90° = 0$. This is the cleanest test in all of linear algebra, and it is the property that makes a [[Basis]] convenient: in an *orthonormal* basis, each coordinate of a vector is simply its dot product with the corresponding axis. That trick is exactly how the coefficients of a [[Fourier Series]] are computed — by projecting a signal onto each sinusoid.
{{< /note >}}

## See also

- [[Vector]]
- [[Basis]]
- [[Fourier Series]]
