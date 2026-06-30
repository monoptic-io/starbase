---
title: Vector
aliases: [vectors]
tags: [linear-algebra]
summary: An arrow with magnitude and direction — the basic object that linear algebra adds, scales, and transforms.
weight: 10
---

# Vector

A **vector** is a quantity with both **magnitude** and **direction** — an arrow. Velocity, force, and displacement are vectors; temperature and mass, which have only size, are not. Once you fix a [[Basis|coordinate frame]], an arrow in the plane is captured by two numbers, its **components**: $\vec v = (v_x, v_y)$, the horizontal and vertical distances from tail to tip. In three dimensions you add a third, and the same arithmetic carries on into any number of dimensions, where the geometric picture fades but the rules do not.

## Components and length

The components are the shadows the arrow casts on the axes. Its magnitude — the length of the arrow — comes straight from the Pythagorean theorem:

{{< eq number="1" >}}
\lVert \vec v \rVert = \sqrt{v_x^2 + v_y^2}.
{{< /eq >}}

Drag the tip below to see the components and length update live. The two grey legs are $v_x$ and $v_y$; the arrow is their resultant.

{{< sketch height="360" caption="Drag anywhere to move the vector's tip. The grey legs are its components; the readout shows length and angle." >}}
if (frame === 0) {
  state.tip = { x: W * 0.34, y: -H * 0.26 };
}
const css = getComputedStyle(document.documentElement);
const accent = (css.getPropertyValue('--accent') || '#5b9cff').trim();
const accent2 = (css.getPropertyValue('--accent-2') || '#f59e0b').trim();
const faint = (css.getPropertyValue('--text-faint') || '#8892a0').trim();
const text = (css.getPropertyValue('--text') || '#e6e6e6').trim();

const ox = W / 2, oy = H / 2;
if (mouse.down) {
  state.tip.x = mouse.x - ox;
  state.tip.y = mouse.y - oy;
}
const vx = state.tip.x, vy = state.tip.y;

ctx.clearRect(0, 0, W, H);
// axes
ctx.strokeStyle = faint; ctx.globalAlpha = 0.5; ctx.lineWidth = 1;
ctx.beginPath(); ctx.moveTo(0, oy); ctx.lineTo(W, oy);
ctx.moveTo(ox, 0); ctx.lineTo(ox, H); ctx.stroke();
ctx.globalAlpha = 1;

const tx = ox + vx, ty = oy + vy;
// component legs
ctx.strokeStyle = faint; ctx.setLineDash([5, 5]); ctx.lineWidth = 1.5;
ctx.beginPath(); ctx.moveTo(ox, oy); ctx.lineTo(tx, oy); ctx.lineTo(tx, ty); ctx.stroke();
ctx.setLineDash([]);

// the vector with arrowhead
ctx.strokeStyle = accent; ctx.fillStyle = accent; ctx.lineWidth = 3;
ctx.beginPath(); ctx.moveTo(ox, oy); ctx.lineTo(tx, ty); ctx.stroke();
const ang = Math.atan2(vy, vx);
ctx.beginPath();
ctx.moveTo(tx, ty);
ctx.lineTo(tx - 12 * Math.cos(ang - 0.4), ty - 12 * Math.sin(ang - 0.4));
ctx.lineTo(tx - 12 * Math.cos(ang + 0.4), ty - 12 * Math.sin(ang + 0.4));
ctx.closePath(); ctx.fill();

// readout (note screen-y points down, so report -vy)
const len = Math.sqrt(vx * vx + vy * vy);
const deg = (-ang * 180 / Math.PI).toFixed(0);
const ux = (vx / 40).toFixed(2), uy = (-vy / 40).toFixed(2);
ctx.fillStyle = accent2;
ctx.font = '14px system-ui, sans-serif';
ctx.fillText('v = (' + ux + ', ' + uy + ')', 14, 22);
ctx.fillStyle = text;
ctx.fillText('|v| = ' + (len / 40).toFixed(2) + '   angle = ' + deg + '°', 14, 42);
{{< /sketch >}}

## Adding and scaling

Two operations make vectors a *vector space*, and almost everything else is built from them. To **add** vectors, place them tip-to-tail and draw the arrow from the first tail to the last tip — equivalently, add componentwise: $(a_x + b_x,\ a_y + b_y)$. To **scale** a vector, multiply every component by a number $c$, stretching its length by $|c|$ and flipping its direction if $c < 0$.

{{< note kind="note" title="Linear combinations" >}}
Combining the two — $c_1\vec v_1 + c_2\vec v_2$ — gives a **linear combination**, the single most important phrase in the subject. The set of all linear combinations of some vectors is their *span*, and a minimal spanning set is a [[Basis]]. A [[Linear Transformation]] is precisely a map that respects linear combinations.
{{< /note >}}

## See also

- [[Dot Product]]
- [[Basis]]
- [[Linear Transformation]]
