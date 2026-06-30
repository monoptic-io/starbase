---
title: Neural Network
aliases: [neural network, deep learning, multilayer perceptron]
tags: [optimization]
summary: Stack layers of simple neurons with nonlinear activations and they learn the curved decision boundaries a single perceptron never could.
weight: 90
---

# Neural Network

A single [[Perceptron]] can only carve space with a straight line. Stack many of them in **layers**, feed each layer's outputs into the next, and you get a **neural network** — a function flexible enough to bend boundaries into any shape and approximate essentially any mapping from inputs to outputs. This is the engine of modern *deep learning*, and a many-layered network is still, at heart, a **multilayer perceptron**.

Each layer does two things. First a [[Linear Transformation|linear]] step — a [[Matrix Multiplication|matrix multiply]] by the layer's weights plus a bias, $\mathbf{z} = W\mathbf{x} + \mathbf{b}$. Then a **nonlinear activation** applied elementwise, $\mathbf{a} = \sigma(\mathbf{z})$. That nonlinearity is the whole trick: without it, stacking linear layers would just collapse into one big linear map (and one straight boundary again). With it, depth lets the network compose simple bends into arbitrarily intricate curved boundaries.

## Signals flowing forward

Computing a network's output is **forward propagation**: activations flow left to right, layer by layer, each neuron summing its weighted inputs and squashing the result. Below, an input layer feeds two hidden layers and an output; pulses trace the connections, blue for excitatory (positive) weights and orange for inhibitory (negative) ones, and neurons brighten as the wave of activation reaches them.

{{< sketch height="380" caption="A small multilayer perceptron computing forward. Activation flows input → hidden → hidden → output; pulses ride each connection (blue = positive weight, orange = negative), and nodes glow as the signal arrives. This forward pass is what backpropagation later differentiates." >}}
if (frame === 0 || !state.layers) {
  state.sizes = [3, 5, 4, 2];
  state.layers = [];
  for (let l = 0; l < state.sizes.length; l++) {
    const nodes = [];
    for (let i = 0; i < state.sizes[l]; i++) nodes.push({});
    state.layers.push(nodes);
  }
  // Random edge weights (sign sets color).
  state.edges = [];
  for (let l = 0; l < state.sizes.length - 1; l++) {
    for (let i = 0; i < state.sizes[l]; i++) {
      for (let j = 0; j < state.sizes[l + 1]; j++) {
        state.edges.push({ l: l, i: i, j: j, w: Math.random() * 2 - 1, off: Math.random() });
      }
    }
  }
  state.speed = 0.55;
}

const cs = getComputedStyle(document.documentElement);
const accent = cs.getPropertyValue('--accent').trim() || '#5b9cff';
const accent2 = cs.getPropertyValue('--accent-2').trim() || '#ff9e64';
const text = cs.getPropertyValue('--text').trim() || '#e8ecf5';
const faint = cs.getPropertyValue('--text-faint').trim() || 'rgba(232,236,245,0.4)';
const border = cs.getPropertyValue('--border').trim() || 'rgba(255,255,255,0.18)';

ctx.clearRect(0, 0, W, H);
ctx.fillStyle = '#0b0e16';
ctx.fillRect(0, 0, W, H);

const nL = state.sizes.length;
const marginX = 70, marginY = 50;
const nodePos = function (l, i) {
  const x = marginX + (W - 2 * marginX) * (l / (nL - 1));
  const cnt = state.sizes[l];
  const y = (cnt === 1) ? H / 2 : marginY + (H - 2 * marginY) * (i / (cnt - 1));
  return [x, y];
};

// The "activation front" sweeps across layers, looping.
const front = (state.speed * t) % (nL - 1 + 1.2);

// Edges + traveling pulses.
for (let e = 0; e < state.edges.length; e++) {
  const ed = state.edges[e];
  const a = nodePos(ed.l, ed.i), b = nodePos(ed.l + 1, ed.j);
  const col = ed.w >= 0 ? accent : accent2;
  ctx.strokeStyle = 'rgba(255,255,255,0.06)';
  ctx.lineWidth = 1;
  ctx.beginPath(); ctx.moveTo(a[0], a[1]); ctx.lineTo(b[0], b[1]); ctx.stroke();
  // Pulse active when the front is crossing this layer gap.
  if (front >= ed.l && front < ed.l + 1) {
    const u = front - ed.l;
    const px = a[0] + (b[0] - a[0]) * u, py = a[1] + (b[1] - a[1]) * u;
    ctx.fillStyle = col;
    ctx.globalAlpha = 0.85;
    ctx.beginPath(); ctx.arc(px, py, 2.5 + 1.5 * Math.abs(ed.w), 0, 2 * Math.PI); ctx.fill();
    ctx.globalAlpha = 1;
  }
}

// Nodes (glow as the front arrives at their layer).
for (let l = 0; l < nL; l++) {
  for (let i = 0; i < state.sizes[l]; i++) {
    const p = nodePos(l, i);
    const near = Math.max(0, 1 - Math.abs(front - l) * 1.4);
    ctx.fillStyle = l === 0 ? accent : (l === nL - 1 ? accent2 : '#16203a');
    ctx.strokeStyle = border; ctx.lineWidth = 1.5;
    ctx.shadowColor = l === nL - 1 ? accent2 : accent;
    ctx.shadowBlur = 18 * near;
    ctx.beginPath(); ctx.arc(p[0], p[1], 13, 0, 2 * Math.PI); ctx.fill(); ctx.stroke();
    ctx.shadowBlur = 0;
  }
}

// Layer labels.
ctx.fillStyle = faint; ctx.font = '12px sans-serif'; ctx.textAlign = 'center';
const labels = ['input', 'hidden', 'hidden', 'output'];
for (let l = 0; l < nL; l++) ctx.fillText(labels[l] || ('layer ' + l), nodePos(l, 0)[0], H - 16);
ctx.textAlign = 'left';
{{< /sketch >}}

## The nonlinearity is the magic

The activation function bends straight lines. A common choice is the hyperbolic tangent, which smoothly saturates large inputs toward $\pm 1$:

{{< plot fn="Math.tanh(x);;Math.max(0,x)" xmin="-4" xmax="4" ymin="-1.5" ymax="3" title="Two activation functions: tanh (S-curve) and ReLU (a hinge)" caption="Squashing nonlinearities like these let stacked layers represent curved, even disconnected, decision regions. Remove them and the whole network collapses back to a single linear map." >}}

Train the weights — usually by [[Gradient Descent]] guided by [[Backpropagation]] — and a network with even one hidden layer becomes a *universal approximator*: given enough neurons it can represent any continuous function to any accuracy. Depth makes that representation efficient, building rich features as a hierarchy of simpler ones.

{{< quiz question="Why does a neural network need nonlinear activation functions between its layers?" options="To make training faster|Without them, stacked linear layers collapse into a single linear map|To save memory|They prevent the weights from changing" answer="2" explain="A composition of linear maps is itself linear. The nonlinearity between layers is what lets depth build curved, expressive boundaries instead of one straight one — exactly what a single perceptron cannot do." >}}

## See also

- [[Perceptron]]
- [[Backpropagation]]
- [[Matrix Multiplication]]
