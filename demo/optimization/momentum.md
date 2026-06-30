---
title: Momentum
aliases: [momentum, heavy ball]
tags: [optimization]
summary: Give the optimizer inertia — accumulate a velocity from past gradients so it powers through ravines and coasts over small bumps.
weight: 40
---

# Momentum

Plain [[Gradient Descent|gradient descent]] has no memory: each step depends only on the gradient *right now*. In a long, narrow ravine — steep walls, gentle floor — that is a disaster. The optimizer keeps overreacting to the steep sidewalls, zig-zagging back and forth across the valley while barely creeping along its length. **Momentum** fixes this by giving the optimizer *inertia*, exactly like a heavy ball rolling downhill. Instead of jumping in the instantaneous gradient direction, you accumulate a velocity:

$$v \leftarrow \mu\,v - \eta\,\nabla L(\theta), \qquad \theta \leftarrow \theta + v.$$

The coefficient $\mu$ (typically around $0.9$) is how much of the previous velocity carries over. Oscillating components — the side-to-side bounces — point in opposite directions on alternate steps and **cancel out** in the running average, while the consistent downhill component along the valley floor **adds up**, building speed. The result is faster, smoother convergence and enough built-up momentum to coast straight over tiny bumps and shallow local dips.

## Watch inertia win the ravine

Both balls below start at the same place in a stretched, ravine-shaped bowl (steep top-to-bottom, shallow left-to-right). The **plain gradient-descent ball** zig-zags wildly across the steep direction, wasting almost all its motion. The **momentum ball** lets those zig-zags cancel and glides along the floor to the bottom far sooner.

{{< sketch height="360" caption="Same start, same ravine. Plain gradient descent (left color) bounces between the steep walls; momentum (right color) averages those bounces away and accelerates along the valley floor to the minimum. The run loops." >}}
if (frame === 0 || !state.gd) {
  state.xmin = -5; state.xmax = 5; state.ymin = -2.4; state.ymax = 2.4;
  // Ravine: gentle in x, steep in y.
  state.loss = function (x, y) { return 0.05 * x * x + 1.0 * y * y; };
  // Render elliptical contours once.
  const IW = Math.max(1, Math.floor(W)), IH = Math.max(1, Math.floor(H));
  const buf = document.createElement('canvas');
  buf.width = IW; buf.height = IH;
  const b = buf.getContext('2d');
  const img = b.createImageData(IW, IH);
  const d = img.data;
  const Lg = new Float32Array(IW * IH);
  let lo = Infinity, hi = -Infinity;
  for (let py = 0; py < IH; py++) {
    const y = state.ymin + (py / (IH - 1)) * (state.ymax - state.ymin);
    for (let px = 0; px < IW; px++) {
      const x = state.xmin + (px / (IW - 1)) * (state.xmax - state.xmin);
      const L = state.loss(x, y);
      Lg[py * IW + px] = L;
      if (L < lo) lo = L; if (L > hi) hi = L;
    }
  }
  const NB = 12;
  const bands = new Int16Array(IW * IH);
  for (let i = 0; i < IW * IH; i++) bands[i] = Math.floor(((Lg[i] - lo) / (hi - lo)) * NB);
  for (let py = 0; py < IH; py++) {
    for (let px = 0; px < IW; px++) {
      const i = py * IW + px;
      const f = 1 - (Lg[i] - lo) / (hi - lo);
      let r = Math.floor(12 + f * 60), g = Math.floor(16 + f * 90), bl = Math.floor(28 + f * 150);
      let edge = false;
      if (px > 0 && bands[i] !== bands[i - 1]) edge = true;
      if (py > 0 && bands[i] !== bands[i - IW]) edge = true;
      if (edge) { r += 45; g += 55; bl += 70; }
      const j = i * 4;
      d[j] = Math.min(255, r); d[j + 1] = Math.min(255, g); d[j + 2] = Math.min(255, bl); d[j + 3] = 255;
    }
  }
  b.putImageData(img, 0, 0);
  state.buf = buf;
  state.reset = function () {
    state.gd = { x: -4.5, y: 1.9, trail: [] };
    state.mo = { x: -4.5, y: 1.9, vx: 0, vy: 0, trail: [] };
    state.tick = 0;
    state.pause = 0;
  };
  state.reset();
}

const cs = getComputedStyle(document.documentElement);
const accent = cs.getPropertyValue('--accent').trim() || '#5b9cff';
const accent2 = cs.getPropertyValue('--accent-2').trim() || '#ff9e64';
const good = cs.getPropertyValue('--good').trim() || '#9ece6a';
const text = cs.getPropertyValue('--text').trim() || '#e8ecf5';

ctx.clearRect(0, 0, W, H);
ctx.drawImage(state.buf, 0, 0, W, H);

const toPX = function (x) { return (x - state.xmin) / (state.xmax - state.xmin) * W; };
const toPY = function (y) { return (y - state.ymin) / (state.ymax - state.ymin) * H; };
const grad = function (x, y) { return [0.1 * x, 2.0 * y]; };

// Advance one optimization step every few frames so the motion is watchable.
state.tick++;
if (state.pause > 0) {
  state.pause--;
} else if (state.tick % 4 === 0) {
  // Plain gradient descent (near its stability limit → zig-zag).
  const etaG = 0.82;
  const gg = grad(state.gd.x, state.gd.y);
  state.gd.x -= etaG * gg[0];
  state.gd.y -= etaG * gg[1];
  state.gd.trail.push([state.gd.x, state.gd.y]);
  // Momentum (heavy ball).
  const etaM = 0.4, mu = 0.82;
  const gm = grad(state.mo.x, state.mo.y);
  state.mo.vx = mu * state.mo.vx - etaM * gm[0];
  state.mo.vy = mu * state.mo.vy - etaM * gm[1];
  state.mo.x += state.mo.vx;
  state.mo.y += state.mo.vy;
  state.mo.trail.push([state.mo.x, state.mo.y]);
  const gdDone = Math.hypot(state.gd.x, state.gd.y) < 0.08;
  const moDone = Math.hypot(state.mo.x, state.mo.y) < 0.08 && Math.hypot(state.mo.vx, state.mo.vy) < 0.02;
  if ((gdDone && moDone) || state.gd.trail.length > 220) state.pause = 70;
  if (state.pause > 0 && state.gd.trail.length > 1) {
    // schedule reset after the pause by flagging
    state.pendingReset = true;
  }
}
if (state.pause === 0 && state.pendingReset) { state.pendingReset = false; state.reset(); }

function drawTrail(b, color) {
  ctx.strokeStyle = color; ctx.lineWidth = 2;
  ctx.beginPath();
  for (let i = 0; i < b.trail.length; i++) {
    const px = toPX(b.trail[i][0]), py = toPY(b.trail[i][1]);
    if (i === 0) ctx.moveTo(px, py); else ctx.lineTo(px, py);
  }
  ctx.stroke();
  const cx = toPX(b.x), cy = toPY(b.y);
  ctx.fillStyle = color; ctx.shadowColor = color; ctx.shadowBlur = 12;
  ctx.beginPath(); ctx.arc(cx, cy, 6.5, 0, 2 * Math.PI); ctx.fill();
  ctx.shadowBlur = 0;
}
drawTrail(state.gd, accent);
drawTrail(state.mo, accent2);

// Legend.
ctx.font = '13px sans-serif';
ctx.fillStyle = accent;
ctx.fillText('● plain gradient descent — ' + state.gd.trail.length + ' steps', 12, 22);
ctx.fillStyle = accent2;
ctx.fillText('● momentum — ' + state.mo.trail.length + ' steps', 12, 42);
{{< /sketch >}}

{{< note kind="tip" title="Momentum is a damped oscillator" >}}
The update $v \leftarrow \mu v - \eta\nabla L$ is exactly a [[Damped Oscillator|damped harmonic oscillator]] descending a potential: the gradient is the restoring force and $(1-\mu)$ plays the role of friction. Tune $\mu$ too high and the ball overshoots and rings; too low and you lose the acceleration. The same physics that governs a [[Pendulum]] governs your optimizer.
{{< /note >}}

## See also

- [[Gradient Descent]]
- [[Loss Landscape]]
- [[Newton's Method]]
