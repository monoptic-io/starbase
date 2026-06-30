---
title: Backpropagation
aliases: [backprop, backpropagation]
tags: [optimization]
summary: The chain rule run in reverse — compute the loss's gradient with respect to every weight in one efficient backward sweep through the network.
weight: 100
---

# Backpropagation

A [[Neural Network]] is trained by [[Gradient Descent]], which needs the gradient of the loss with respect to *every* weight — potentially billions of them. Computing each one separately would be hopeless. **Backpropagation** is the algorithm that gets them all in a single backward pass, and it is nothing more than the **chain rule** of calculus applied with ruthless efficiency.

The idea: a network is a long composition of functions, $L = \ell(f_n(\cdots f_2(f_1(\mathbf{x}))\cdots))$. The chain rule says the derivative of a composition is the *product* of the local derivatives along the way:

{{< eq number="1" >}}
\frac{\partial L}{\partial \theta} = \frac{\partial L}{\partial a_n}\cdot\frac{\partial a_n}{\partial a_{n-1}}\cdots\frac{\partial a_{k}}{\partial \theta}.
{{< /eq >}}

Backprop computes this right-to-left. After a **forward pass** stores every neuron's activation, error is injected at the output and flows *backward*: at each layer the incoming gradient is multiplied by that layer's local derivative, handed to the layer before it, and used to update that layer's weights. One forward sweep, one backward sweep, and every gradient is known.

## Error flowing backward

Below, a signal first runs **forward** through a chain of operations (faint pulses, left to right) to produce a loss. Then a gradient is injected at the loss and travels **backward**, and at each edge it is multiplied by that step's local derivative — the running product is the chain rule accumulating. Watch the gradient value change as it propagates back to the input.

{{< sketch height="340" caption="The two passes of backpropagation along a chain. Forward (left→right) computes activations; backward (right→left) multiplies local derivatives via the chain rule so the gradient at every node is known after a single reverse sweep." >}}
if (frame === 0 || !state.deriv) {
  state.nodes = ['x', 'f₁', 'f₂', 'f₃', 'L'];
  state.deriv = [0.8, -1.3, 0.6, 1.7]; // local derivative on each edge
  state.phase = 'forward';
  state.p = 0;
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

const N = state.nodes.length;
const marginX = 56;
const y = H / 2 - 10;
const nx = function (k) { return marginX + (W - 2 * marginX) * (k / (N - 1)); };

// Advance the animation.
state.p += dt * 0.7;
if (state.p >= N - 1) {
  state.p = 0;
  state.phase = (state.phase === 'forward') ? 'backward' : 'forward';
}
const forward = state.phase === 'forward';
// position of the active pulse measured from the left
const pos = forward ? state.p : (N - 1 - state.p);

// Edges.
for (let k = 0; k < N - 1; k++) {
  ctx.strokeStyle = border; ctx.lineWidth = 2;
  ctx.beginPath(); ctx.moveTo(nx(k) + 18, y); ctx.lineTo(nx(k + 1) - 18, y); ctx.stroke();
  // local derivative label on each edge
  ctx.fillStyle = faint; ctx.font = '11px sans-serif'; ctx.textAlign = 'center';
  ctx.fillText("∂=" + state.deriv[k].toFixed(1), (nx(k) + nx(k + 1)) / 2, y - 26);
}

// Moving pulse.
const segf = Math.floor(state.p);
const u = state.p - segf;
let pulseX;
if (forward) pulseX = nx(segf) + (nx(segf + 1) - nx(segf)) * u;
else { const kk = N - 1 - segf; pulseX = nx(kk) - (nx(kk) - nx(kk - 1)) * u; }
ctx.fillStyle = forward ? accent : accent2;
ctx.shadowColor = ctx.fillStyle; ctx.shadowBlur = 16;
ctx.beginPath(); ctx.arc(pulseX, y, 7, 0, 2 * Math.PI); ctx.fill();
ctx.shadowBlur = 0;

// Accumulated gradient during the backward pass.
let grad = 1;
const reached = Math.floor(pos + 0.5);
for (let k = N - 2; k >= reached; k--) grad *= state.deriv[k];

// Nodes.
for (let k = 0; k < N; k++) {
  const active = Math.abs(pos - k) < 0.5;
  ctx.fillStyle = (k === 0) ? accent : (k === N - 1 ? accent2 : '#16203a');
  ctx.strokeStyle = active ? (forward ? accent : accent2) : border;
  ctx.lineWidth = active ? 3 : 1.5;
  ctx.beginPath(); ctx.arc(nx(k), y, 18, 0, 2 * Math.PI); ctx.fill(); ctx.stroke();
  ctx.fillStyle = text; ctx.font = '14px sans-serif'; ctx.textAlign = 'center';
  ctx.fillText(state.nodes[k], nx(k), y + 5);
}
ctx.textAlign = 'left';

// HUD.
ctx.fillStyle = forward ? accent : accent2;
ctx.font = '14px sans-serif';
ctx.fillText(forward ? 'forward pass  →   computing activations'
                     : '←  backward pass   ∂L/∂node = ' + grad.toFixed(3), 16, 28);
{{< /sketch >}}

{{< note kind="key" title="Why backprop made deep learning possible" >}}
The miracle is the *cost*. A naive gradient would re-traverse the network once per parameter. Backpropagation reuses the shared intermediate derivatives, computing **all** gradients in time proportional to a single forward pass. This is reverse-mode automatic differentiation, and it is the reason training networks with millions of weights is feasible at all.
{{< /note >}}

## See also

- [[Neural Network]]
- [[Gradient Descent]]
- [[Perceptron]]
