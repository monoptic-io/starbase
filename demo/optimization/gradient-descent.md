---
title: Gradient Descent
aliases: [gradient descent, steepest descent]
tags: [optimization]
summary: The fundamental learning algorithm — repeatedly step a little way downhill along the negative gradient of a loss until you reach a valley.
weight: 10
---

# Gradient Descent

**Gradient descent** is the single most important algorithm in machine learning, and it is almost embarrassingly simple. You are standing somewhere on a [[Loss Landscape|loss surface]] and you want to get low. The gradient $\nabla L$ points in the direction of *steepest increase*, so its negative points straight downhill. Take a small step that way, recompute, and repeat:

$$\theta \leftarrow \theta - \eta\,\nabla L(\theta).$$

Here $\theta$ is everything you can tune (the weights of a model), $L$ is the loss you want to minimize, and $\eta$ — the **learning rate** — is how big a step you dare to take. Too small and you crawl; too large and you overshoot the valley and bounce out. The whole craft of training comes down to choosing a good direction and a good step size, over and over.

## Why the negative gradient is "downhill"

The gradient is the vector of partial derivatives, $\nabla L = \big(\partial L/\partial\theta_1,\dots\big)$. For a tiny step $\Delta\theta$, the change in loss is, to first order, the [[Dot Product]] $\Delta L \approx \nabla L \cdot \Delta\theta$. Among all steps of a fixed length, the one that makes $\Delta L$ as *negative* as possible is the one pointing exactly opposite the gradient — that is what "steepest descent" means. Each update is therefore the greediest locally-improving move you can make.

{{< note kind="tip" title="The learning rate is everything" >}}
The step size $\eta$ is the knob you will fiddle with most. Picture a narrow valley: a small $\eta$ inches safely to the bottom; a large $\eta$ ricochets between the walls and may diverge. When descent "explodes" into `NaN`, the learning rate is the first suspect.
{{< /note >}}

## Roll a ball down the loss surface

Below is a [[Loss Landscape|loss landscape]] drawn as **contour lines** — each ring is a level set, and the valleys (two of them here) glow brightest. A ball starts at a point and follows $-\nabla L$ downhill, one step per frame, tracing its path behind it. Watch how it always cuts *perpendicular* to the contours, the steepest way down, and how its starting point decides which basin it falls into.

**Click anywhere to drop a new ball** and see where it rolls.

{{< sketch height="380" caption="Gradient descent on a two-basin loss surface shown as contour rings. The ball follows the negative gradient, always crossing contours at right angles. Click to drop a fresh ball — the starting point decides which minimum it finds." >}}
if (frame === 0) {
  state.xmin = -3; state.xmax = 3; state.ymin = -3; state.ymax = 3;
  state.loss = function (x, y) {
    return 1.3
      - 1.00 * Math.exp(-(((x - 1) * (x - 1)) + ((y - 1) * (y - 1))) / 1.5)
      - 0.80 * Math.exp(-(((x + 1.3) * (x + 1.3)) + ((y + 0.8) * (y + 0.8))) / 0.8)
      + 0.05 * (x * x + y * y);
  };
  state.grad = function (x, y) {
    let gx = 0.1 * x, gy = 0.1 * y;
    const e1 = Math.exp(-(((x - 1) * (x - 1)) + ((y - 1) * (y - 1))) / 1.5);
    gx += 1.00 * e1 * (2 * (x - 1) / 1.5); gy += 1.00 * e1 * (2 * (y - 1) / 1.5);
    const e2 = Math.exp(-(((x + 1.3) * (x + 1.3)) + ((y + 0.8) * (y + 0.8))) / 0.8);
    gx += 0.80 * e2 * (2 * (x + 1.3) / 0.8); gy += 0.80 * e2 * (2 * (y + 0.8) / 0.8);
    return [gx, gy];
  };
  // Render contour heatmap once into an offscreen buffer.
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
  const NB = 14;
  const bands = new Int16Array(IW * IH);
  for (let i = 0; i < IW * IH; i++) bands[i] = Math.floor(((Lg[i] - lo) / (hi - lo)) * NB);
  for (let py = 0; py < IH; py++) {
    for (let px = 0; px < IW; px++) {
      const i = py * IW + px;
      const f = 1 - (Lg[i] - lo) / (hi - lo);
      let r = Math.floor(12 + f * 70), g = Math.floor(18 + f * 100), bl = Math.floor(30 + f * 160);
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
  state.ball = { x: -2.3, y: 2.3, trail: [] };
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

// Drop a new ball on click.
if (mouse.clicked) {
  const nx = state.xmin + (mouse.x / W) * (state.xmax - state.xmin);
  const ny = state.ymin + (mouse.y / H) * (state.ymax - state.ymin);
  state.ball = { x: nx, y: ny, trail: [] };
}

// One gradient-descent step per frame.
const ball = state.ball;
const eta = 0.06;
const g = state.grad(ball.x, ball.y);
const gmag = Math.hypot(g[0], g[1]);
if (gmag > 1e-3) {
  ball.x -= eta * g[0];
  ball.y -= eta * g[1];
  ball.trail.push([ball.x, ball.y]);
  if (ball.trail.length > 600) ball.trail.shift();
}

// Draw the trail.
ctx.strokeStyle = accent;
ctx.lineWidth = 2;
ctx.beginPath();
for (let i = 0; i < ball.trail.length; i++) {
  const px = toPX(ball.trail[i][0]), py = toPY(ball.trail[i][1]);
  if (i === 0) ctx.moveTo(px, py); else ctx.lineTo(px, py);
}
ctx.stroke();

// Draw the ball with a glow.
const bx = toPX(ball.x), by = toPY(ball.y);
ctx.fillStyle = accent2;
ctx.shadowColor = accent2; ctx.shadowBlur = 14;
ctx.beginPath(); ctx.arc(bx, by, 7, 0, 2 * Math.PI); ctx.fill();
ctx.shadowBlur = 0;

// HUD.
ctx.fillStyle = text;
ctx.font = '13px sans-serif';
ctx.fillText('loss = ' + state.loss(ball.x, ball.y).toFixed(3), 10, 20);
ctx.fillStyle = gmag > 1e-3 ? good : accent2;
ctx.fillText(gmag > 1e-3 ? 'descending…' : 'at a minimum', 10, H - 12);
{{< /sketch >}}

## See also

- [[Loss Landscape]]
- [[Momentum]]
- [[Dynamical System]]
