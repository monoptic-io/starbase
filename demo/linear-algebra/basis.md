---
title: Basis
aliases: [change of basis, coordinate system]
tags: [linear-algebra]
summary: A set of reference directions that turns a vector into a list of numbers — change the directions and the numbers change, though the vector does not.
weight: 60
---

# Basis

A **basis** is a minimal set of [[Vector|vectors]] whose linear combinations reach every point in a space, with none of them redundant. In the plane, any two non-parallel vectors form a basis. Choosing a basis is what lets you replace an abstract arrow with a concrete list of **coordinates**: the numbers tell you how much of each basis vector to add up. The crucial subtlety is that *the coordinates are not the vector* — they are the vector **as seen from a particular set of axes**.

## Same vector, different numbers

The standard basis $\hat\imath = (1,0)$, $\hat\jmath = (0,1)$ feels canonical only by habit. The very same arrow has different coordinates in a tilted basis. Below, one fixed vector (white) is read off in two coordinate systems — the standard grid and a rotated one. The arrow never moves; only the numbers describing it do.

{{< sketch height="380" caption="One fixed vector (white) seen in two bases: the upright grid (blue) and a tilted grid (amber). Drag to spin the tilted basis — the arrow stays put, but its coordinates in the amber frame change continuously." >}}
if (frame === 0) { state.phi = 0.6; }
const css = getComputedStyle(document.documentElement);
const accent = (css.getPropertyValue('--accent') || '#5b9cff').trim();
const accent2 = (css.getPropertyValue('--accent-2') || '#f59e0b').trim();
const faint = (css.getPropertyValue('--text-faint') || '#7a8390').trim();
const text = (css.getPropertyValue('--text') || '#e6e6e6').trim();

const ox = W / 2, oy = H / 2;
if (mouse.down) { state.phi = Math.atan2(-(mouse.y - oy), mouse.x - ox); }
const phi = state.phi;
const u = Math.min(W, H) / 7;

// the fixed vector in world coords
const v = { x: 1.6, y: 1.1 };

ctx.clearRect(0, 0, W, H);

// standard grid (faint)
ctx.strokeStyle = accent; ctx.globalAlpha = 0.18; ctx.lineWidth = 1;
for (let i = -4; i <= 4; i++) {
  ctx.beginPath(); ctx.moveTo(ox + i * u, 0); ctx.lineTo(ox + i * u, H); ctx.stroke();
  ctx.beginPath(); ctx.moveTo(0, oy + i * u); ctx.lineTo(W, oy + i * u); ctx.stroke();
}
// tilted basis vectors
const e1 = { x: Math.cos(phi), y: Math.sin(phi) };
const e2 = { x: -Math.sin(phi), y: Math.cos(phi) };
ctx.strokeStyle = accent2; ctx.globalAlpha = 0.18;
for (let i = -4; i <= 4; i++) {
  const ax = ox + (i * e1.x) * u, ay = oy - (i * e1.y) * u;
  ctx.beginPath();
  ctx.moveTo(ax - e2.x * u * 5, ay + e2.y * u * 5);
  ctx.lineTo(ax + e2.x * u * 5, ay - e2.y * u * 5); ctx.stroke();
  const bx = ox + (i * e2.x) * u, by = oy - (i * e2.y) * u;
  ctx.beginPath();
  ctx.moveTo(bx - e1.x * u * 5, by + e1.y * u * 5);
  ctx.lineTo(bx + e1.x * u * 5, by - e1.y * u * 5); ctx.stroke();
}
ctx.globalAlpha = 1;

function arrow(vx, vy, col, lw) {
  ctx.strokeStyle = col; ctx.fillStyle = col; ctx.lineWidth = lw;
  ctx.beginPath(); ctx.moveTo(ox, oy); ctx.lineTo(ox + vx, oy + vy); ctx.stroke();
  const ang = Math.atan2(vy, vx);
  ctx.beginPath();
  ctx.moveTo(ox + vx, oy + vy);
  ctx.lineTo(ox + vx - 11 * Math.cos(ang - 0.4), oy + vy - 11 * Math.sin(ang - 0.4));
  ctx.lineTo(ox + vx - 11 * Math.cos(ang + 0.4), oy + vy - 11 * Math.sin(ang + 0.4));
  ctx.closePath(); ctx.fill();
}
// basis arrows
arrow(e1.x * u, -e1.y * u, accent2, 2);
arrow(e2.x * u, -e2.y * u, accent2, 2);
arrow(u, 0, accent, 2);
arrow(0, -u, accent, 2);
// the fixed vector
arrow(v.x * u, -v.y * u, '#ffffff', 3);

// coordinates in each basis
const c2_1 = v.x * e1.x + v.y * e1.y;  // dot with e1
const c2_2 = v.x * e2.x + v.y * e2.y;  // dot with e2
ctx.font = '13px ui-monospace, monospace';
ctx.fillStyle = accent;
ctx.fillText('standard: (' + v.x.toFixed(2) + ', ' + v.y.toFixed(2) + ')', 14, 22);
ctx.fillStyle = accent2;
ctx.fillText('tilted:   (' + c2_1.toFixed(2) + ', ' + c2_2.toFixed(2) + ')', 14, 42);
{{< /sketch >}}

## Change of basis

Switching bases is itself a [[Linear Transformation]]. If the new basis vectors are the columns of a matrix $P$, then $P$ converts new coordinates into old, and $P^{-1}$ converts old into new. When the basis is **orthonormal** — mutually perpendicular unit vectors — the conversion is delightfully cheap: each coordinate is just the [[Dot Product]] of the vector with the corresponding axis, no matrix inverse required.

## Picking the right axes

{{< note kind="key" title="A good basis makes a problem obvious" >}}
Much of applied mathematics is the art of choosing a basis in which a hard problem becomes easy. A [[Fourier Series]] expresses a signal in a basis of **sinusoids**, turning differentiation into multiplication. The [[Eigenvalues and Eigenvectors|eigenvector]] basis of a matrix turns its action into simple stretching along each axis. Same vector, same physics — but the right coordinates make the structure leap out.
{{< /note >}}

## See also

- [[Dot Product]]
- [[Eigenvalues and Eigenvectors]]
- [[Fourier Series]]
