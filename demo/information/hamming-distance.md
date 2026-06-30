---
title: Hamming Distance
aliases: [hamming]
tags: [information]
summary: The number of positions in which two equal-length strings differ — the metric behind error correction.
weight: 80
---

# Hamming Distance

The **Hamming distance** between two strings of equal length is simply the number of positions where they disagree. For two binary words it is the count of bits you would have to flip to turn one into the other:

$$d_H(x, y) = \#\{\, i : x_i \neq y_i \,\} = \sum_i (x_i \oplus y_i).$$

For example, `1011101` and `1001001` differ in positions 3, 5, and 6, so their Hamming distance is **3**. It is a genuine *metric* — non-negative, symmetric, zero only for identical strings, and obeying the triangle inequality — which is what lets us reason about codes geometrically, as points scattered in a space of strings.

## Why it sets correcting power

A code is just a chosen set of valid codewords. Its **minimum distance** $d_{\min}$ — the smallest Hamming distance between any two codewords — determines everything about its robustness. Picture each codeword sitting at the center of a ball of nearby strings; if the balls don't overlap, a received word that falls inside one gets decoded to its center.

{{< note kind="key" title="The detection / correction rule" >}}
A code with minimum distance $d_{\min}$ can:

- **detect** up to $d_{\min} - 1$ bit errors, and
- **correct** up to $\left\lfloor \dfrac{d_{\min} - 1}{2} \right\rfloor$ bit errors.

So distance 3 (like [[Error-Correcting Code|Hamming(7,4)]]) corrects 1 error; distance 5 corrects 2. To pack in more correcting power, spread the codewords farther apart — which costs redundancy.
{{< /note >}}

## Distance as a cube

For 3-bit strings, every word is a corner of a cube and Hamming distance is the number of edges of the shortest path between corners. Hover a corner to highlight its neighbors at each distance — adjacent corners differ by 1 bit, face-diagonals by 2, the body-diagonal by 3.

{{< sketch height="360" caption="The 3-bit cube. Each corner is a binary word; Hamming distance is the edge-count between corners. Move the mouse near a corner to color all others by their distance from it (green = far, blue = near)." >}}
if (frame === 0) {
  // 8 corners of the unit cube, with bit labels
  state.pts = [];
  for (let i = 0; i < 8; i++) {
    const bx = (i & 1), by = (i >> 1) & 1, bz = (i >> 2) & 1;
    state.pts.push({b: [bx, by, bz], label: '' + bz + by + bx});
  }
}
ctx.clearRect(0, 0, W, H);
const css = (v, f) => (getComputedStyle(document.documentElement).getPropertyValue(v).trim() || f);
const accent = css('--accent', '#5b9cff');
const good = css('--good', '#4ec98f');
const warn = css('--warn', '#e0a458');
const text = css('--text', '#e6e9ef');
const faint = css('--text-faint', 'rgba(230,233,239,0.55)');
const border = css('--border', 'rgba(255,255,255,0.25)');

// project cube to 2D with a gentle fixed rotation
const cx = W / 2, cy = H / 2, S = Math.min(W, H) * 0.30;
const ang = 0.62, ca = Math.cos(ang), sa = Math.sin(ang);
const proj = (b) => {
  // center coords in [-1,1]
  let X = (b[0] - 0.5) * 2, Y = (b[1] - 0.5) * 2, Z = (b[2] - 0.5) * 2;
  // rotate around Y then tilt
  const px = X * ca + Z * sa;
  const pz = -X * sa + Z * ca;
  const py = Y * 0.82 + pz * 0.32;
  return {x: cx + px * S, y: cy + py * S};
};
const P = state.pts.map(p => proj(p.b));

// find nearest corner to the mouse
let sel = -1, best = 1e9;
for (let i = 0; i < 8; i++) {
  const dx = P[i].x - mouse.x, dy = P[i].y - mouse.y;
  const dd = dx * dx + dy * dy;
  if (dd < best) { best = dd; sel = i; }
}
if (best > (S * 0.9) * (S * 0.9)) sel = -1; // only if reasonably close

const ham = (a, c) => (a[0]^c[0]) + (a[1]^c[1]) + (a[2]^c[2]);
const distColor = (d) => d === 1 ? accent : d === 2 ? warn : good;

// edges: connect corners at Hamming distance 1
ctx.lineWidth = 1.4;
for (let i = 0; i < 8; i++) for (let j = i + 1; j < 8; j++) {
  if (ham(state.pts[i].b, state.pts[j].b) === 1) {
    const hot = sel >= 0 && (i === sel || j === sel);
    ctx.strokeStyle = hot ? accent : border;
    ctx.globalAlpha = hot ? 1 : 0.6;
    ctx.beginPath(); ctx.moveTo(P[i].x, P[i].y); ctx.lineTo(P[j].x, P[j].y); ctx.stroke();
  }
}
ctx.globalAlpha = 1;

// corners
ctx.font = '12px monospace';
ctx.textAlign = 'center';
ctx.textBaseline = 'middle';
for (let i = 0; i < 8; i++) {
  let col = accent, r = 9;
  if (sel >= 0) {
    if (i === sel) { col = text; r = 12; }
    else col = distColor(ham(state.pts[sel].b, state.pts[i].b));
  }
  ctx.beginPath(); ctx.arc(P[i].x, P[i].y, r, 0, 2 * Math.PI);
  ctx.fillStyle = col; ctx.fill();
  ctx.fillStyle = '#0b0e16';
  ctx.fillText(state.pts[i].label, P[i].x, P[i].y + 0.5);
}
ctx.textAlign = 'left';
ctx.textBaseline = 'alphabetic';

ctx.fillStyle = faint;
ctx.font = '13px sans-serif';
if (sel >= 0) ctx.fillText('from ' + state.pts[sel].label + ':  blue = 1 flip · orange = 2 · green = 3', 12, H - 14);
else ctx.fillText('hover a corner to color the rest by Hamming distance', 12, H - 14);
{{< /sketch >}}

## Beyond bits

Hamming distance is not only for codes. It measures how many single-character edits separate two DNA strands, quantifies how different two hash fingerprints are, and serves as a similarity metric in machine learning whenever data is categorical. Anywhere you compare strings position by position, it is the natural ruler — and it is the foundation on which every [[Error-Correcting Code]] is designed.

{{< quiz question="Two codewords have Hamming distance 5. How many bit errors can a code with this as its minimum distance correct?" options="1|2|4|5" answer="2" explain="Correcting power is floor((d−1)/2) = floor(4/2) = 2. Distance 5 also detects up to 4 errors." >}}

## See also

- [[Error-Correcting Code]]
- [[Channel Capacity]]
- [[Data Compression]]
