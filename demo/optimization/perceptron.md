---
title: Perceptron
aliases: [perceptron, linear classifier]
tags: [optimization]
summary: The simplest learning machine — a weighted sum past a threshold that learns to separate two classes by nudging its weights every time it errs.
weight: 80
---

# Perceptron

The **perceptron** is the atom of machine learning: a single artificial neuron, invented by Frank Rosenblatt in 1958. It takes an input vector $\mathbf{x}$, forms the weighted sum $\mathbf{w}\cdot\mathbf{x} + b$ — a [[Dot Product]] plus a bias — and fires $+1$ if that exceeds zero, $-1$ otherwise:

$$\hat{y} = \operatorname{sign}(\mathbf{w}\cdot\mathbf{x} + b).$$

Geometrically, $\mathbf{w}\cdot\mathbf{x} + b = 0$ is a **line** (a hyperplane in higher dimensions), and the perceptron simply asks which side of that line a point falls on. The weight vector $\mathbf{w}$ is the line's normal direction; learning means rotating and shifting that line until it cleanly separates the two classes. It is the simplest [[Linear Transformation|linear]] decision rule there is.

## The learning rule

The perceptron learns from its mistakes, one at a time. Show it a labeled point $(\mathbf{x}, y)$. If it classifies correctly, do nothing. If it errs, nudge the weights *toward* getting that point right:

$$\mathbf{w} \leftarrow \mathbf{w} + \eta\,y\,\mathbf{x}, \qquad b \leftarrow b + \eta\,y.$$

That is it. This is a form of [[Gradient Descent]] on a hinge-shaped error, and Rosenblatt proved the **perceptron convergence theorem**: if the two classes *can* be separated by a line, this rule is guaranteed to find one in a finite number of updates.

## Watch the line learn

Two clouds of points, two classes. The decision line starts at a random angle and **rotates into place** as the perceptron processes points and corrects its mistakes; misclassified points are ringed until the line sweeps to the right side of them. The arrow shows $\mathbf{w}$, the direction the line considers "positive."

**Click to add a point** (the label alternates each click) and watch the boundary adjust.

{{< sketch height="400" caption="A perceptron learning to separate two point clouds. The line is the decision boundary w·x+b=0; the arrow is the weight vector w. Each misclassified point nudges the weights, rotating the line into a separator. Click to add points and perturb the problem." >}}
if (frame === 0 || !state.pts) {
  state.pts = [];
  const cloud = function (cx, cy, label, n) {
    for (let i = 0; i < n; i++) {
      const a = Math.random() * Math.PI * 2, r = Math.random() * 0.22;
      state.pts.push({ x: cx + Math.cos(a) * r, y: cy + Math.sin(a) * r, label: label });
    }
  };
  cloud(-0.45, -0.3, -1, 22);
  cloud(0.45, 0.32, 1, 22);
  state.w0 = Math.random() * 2 - 1;
  state.w1 = Math.random() * 2 - 1;
  state.b = Math.random() * 0.4 - 0.2;
  state.nextLabel = 1;
  state.idx = 0;
}

const cs = getComputedStyle(document.documentElement);
const accent = cs.getPropertyValue('--accent').trim() || '#5b9cff';
const accent2 = cs.getPropertyValue('--accent-2').trim() || '#ff9e64';
const good = cs.getPropertyValue('--good').trim() || '#9ece6a';
const text = cs.getPropertyValue('--text').trim() || '#e8ecf5';
const faint = cs.getPropertyValue('--text-faint').trim() || 'rgba(232,236,245,0.45)';
const border = cs.getPropertyValue('--border').trim() || 'rgba(255,255,255,0.18)';

ctx.clearRect(0, 0, W, H);
ctx.fillStyle = '#0b0e16';
ctx.fillRect(0, 0, W, H);

const toPX = function (x) { return (x + 1) / 2 * W; };
const toPY = function (y) { return (1 - (y + 1) / 2) * H; };
const fromPX = function (px) { return px / W * 2 - 1; };
const fromPY = function (py) { return -(py / H * 2 - 1); };

// Add a point on click.
if (mouse.clicked) {
  state.pts.push({ x: fromPX(mouse.x), y: fromPY(mouse.y), label: state.nextLabel });
  state.nextLabel = -state.nextLabel;
}

// Learn: process a couple of points per frame.
const eta = 0.04;
for (let s = 0; s < 2 && state.pts.length > 0; s++) {
  const p = state.pts[state.idx % state.pts.length];
  state.idx++;
  const pred = (state.w0 * p.x + state.w1 * p.y + state.b) >= 0 ? 1 : -1;
  if (pred !== p.label) {
    state.w0 += eta * p.label * p.x;
    state.w1 += eta * p.label * p.y;
    state.b += eta * p.label;
  }
}

// Decision line w0*x + w1*y + b = 0.
let A, B;
if (Math.abs(state.w1) > Math.abs(state.w0)) {
  const x1 = -1.3, x2 = 1.3;
  A = [x1, -(state.w0 * x1 + state.b) / state.w1];
  B = [x2, -(state.w0 * x2 + state.b) / state.w1];
} else {
  const y1 = -1.3, y2 = 1.3;
  A = [-(state.w1 * y1 + state.b) / state.w0, y1];
  B = [-(state.w1 * y2 + state.b) / state.w0, y2];
}
ctx.strokeStyle = good; ctx.lineWidth = 3;
ctx.beginPath(); ctx.moveTo(toPX(A[0]), toPY(A[1])); ctx.lineTo(toPX(B[0]), toPY(B[1])); ctx.stroke();

// Weight-vector arrow from the line's nearest point to the origin.
const wmag = Math.hypot(state.w0, state.w1) || 1;
const foot = [-state.b * state.w0 / (wmag * wmag), -state.b * state.w1 / (wmag * wmag)];
const tip = [foot[0] + state.w0 / wmag * 0.32, foot[1] + state.w1 / wmag * 0.32];
ctx.strokeStyle = faint; ctx.lineWidth = 2;
ctx.beginPath(); ctx.moveTo(toPX(foot[0]), toPY(foot[1])); ctx.lineTo(toPX(tip[0]), toPY(tip[1])); ctx.stroke();
const ang = Math.atan2(toPY(tip[1]) - toPY(foot[1]), toPX(tip[0]) - toPX(foot[0]));
ctx.fillStyle = faint;
ctx.beginPath();
ctx.moveTo(toPX(tip[0]), toPY(tip[1]));
ctx.lineTo(toPX(tip[0]) - 9 * Math.cos(ang - 0.4), toPY(tip[1]) - 9 * Math.sin(ang - 0.4));
ctx.lineTo(toPX(tip[0]) - 9 * Math.cos(ang + 0.4), toPY(tip[1]) - 9 * Math.sin(ang + 0.4));
ctx.fill();

// Points (ring the misclassified ones).
let correct = 0;
for (let i = 0; i < state.pts.length; i++) {
  const p = state.pts[i];
  const pred = (state.w0 * p.x + state.w1 * p.y + state.b) >= 0 ? 1 : -1;
  if (pred === p.label) correct++;
  const px = toPX(p.x), py = toPY(p.y);
  ctx.fillStyle = p.label === 1 ? accent2 : accent;
  ctx.beginPath(); ctx.arc(px, py, 6, 0, 2 * Math.PI); ctx.fill();
  if (pred !== p.label) {
    ctx.strokeStyle = '#ff5d6c'; ctx.lineWidth = 2;
    ctx.beginPath(); ctx.arc(px, py, 10, 0, 2 * Math.PI); ctx.stroke();
  }
}

// HUD.
ctx.fillStyle = text; ctx.font = '13px sans-serif';
ctx.fillText('accuracy ' + correct + ' / ' + state.pts.length, 12, 22);
ctx.fillStyle = state.nextLabel === 1 ? accent2 : accent;
ctx.fillText('next click adds a ' + (state.nextLabel === 1 ? 'orange (+1)' : 'blue (−1)') + ' point', 12, H - 12);
{{< /sketch >}}

{{< note kind="warning" title="The wall a single perceptron hits" >}}
A perceptron can only ever draw a *straight* boundary. The famous counterexample is **XOR**: two classes arranged so no single line can separate them. This limitation, publicized by Minsky and Papert in 1969, stalled the field for years — until people stacked perceptrons into a [[Neural Network]] and learned curved boundaries.
{{< /note >}}

## See also

- [[Dot Product]]
- [[Linear Transformation]]
- [[Neural Network]]
