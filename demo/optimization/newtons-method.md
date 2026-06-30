---
title: Newton's Method
aliases: [newton's method, newton-raphson]
tags: [optimization]
summary: Use curvature, not just slope — fit a tangent (or a parabola) and jump straight to where it predicts the answer, converging breathtakingly fast.
weight: 60
---

# Newton's Method

[[Gradient Descent|Gradient descent]] uses only the *slope* and inches along with a hand-tuned step size. **Newton's method** is greedier and smarter: it also uses the **curvature** — the second derivative — to estimate exactly how far to jump. For finding a root of $f$ (a point where $f(x)=0$), you draw the tangent line at your current guess and slide to where that line crosses zero:

$$x_{n+1} = x_n - \frac{f(x_n)}{f'(x_n)}.$$

For *optimization* you want a minimum, where the derivative of the loss is zero — so you simply apply the same idea to $f = L'$, dividing the gradient by the curvature $L''$ instead of by a guessed learning rate:

$$x_{n+1} = x_n - \frac{L'(x_n)}{L''(x_n)}.$$

The payoff is **quadratic convergence**: once you are near the answer, the number of correct digits roughly *doubles* every step. The catch is that you must compute (and, in many dimensions, invert) the second-derivative information — the Hessian — which can be expensive or unstable far from the solution.

## Tangent lines homing in

Watch Newton's method find a root of $f(x) = x^3 - 2x - 5$. From the current guess it rides the **tangent line** down to the axis; that crossing becomes the next guess. After just a few steps the tangent is essentially sitting on the root.

**Click to choose a new starting point** and watch it converge again.

{{< sketch height="380" caption="Newton's method for f(x)=x³−2x−5. Each step follows the tangent at the current guess to where it hits y=0, giving the next guess. Convergence is so fast the steps pile up on the root. Click to restart from a new x₀." >}}
if (frame === 0 || state.x === undefined) {
  state.xmin = -1; state.xmax = 3.6; state.ymin = -10; state.ymax = 14;
  state.f = function (x) { return x * x * x - 2 * x - 5; };
  state.fp = function (x) { return 3 * x * x - 2; };
  state.restart = function (x0) {
    state.x = x0; state.steps = [x0]; state.timer = 0; state.done = false; state.pause = 0;
  };
  state.restart(3.4);
}

const cs = getComputedStyle(document.documentElement);
const accent = cs.getPropertyValue('--accent').trim() || '#5b9cff';
const accent2 = cs.getPropertyValue('--accent-2').trim() || '#ff9e64';
const good = cs.getPropertyValue('--good').trim() || '#9ece6a';
const text = cs.getPropertyValue('--text').trim() || '#e8ecf5';
const faint = cs.getPropertyValue('--text-faint').trim() || 'rgba(232,236,245,0.4)';
const border = cs.getPropertyValue('--border').trim() || 'rgba(255,255,255,0.18)';

ctx.clearRect(0, 0, W, H);
ctx.fillStyle = '#0b0e16';
ctx.fillRect(0, 0, W, H);

const padL = 36, padR = 16, padT = 16, padB = 28;
const toPX = function (x) { return padL + (x - state.xmin) / (state.xmax - state.xmin) * (W - padL - padR); };
const toPY = function (y) { return padT + (1 - (y - state.ymin) / (state.ymax - state.ymin)) * (H - padT - padB); };

// Axes.
ctx.strokeStyle = border; ctx.lineWidth = 1;
ctx.beginPath(); ctx.moveTo(toPX(state.xmin), toPY(0)); ctx.lineTo(toPX(state.xmax), toPY(0)); ctx.stroke();
ctx.beginPath(); ctx.moveTo(toPX(0), toPY(state.ymin)); ctx.lineTo(toPX(0), toPY(state.ymax)); ctx.stroke();
ctx.fillStyle = faint; ctx.font = '11px sans-serif';
ctx.fillText('y = 0', toPX(state.xmax) - 38, toPY(0) - 6);

// The curve.
ctx.strokeStyle = accent; ctx.lineWidth = 2.5;
ctx.beginPath();
for (let i = 0; i <= 240; i++) {
  const x = state.xmin + (i / 240) * (state.xmax - state.xmin);
  const px = toPX(x), py = toPY(state.f(x));
  if (i === 0) ctx.moveTo(px, py); else ctx.lineTo(px, py);
}
ctx.stroke();

// Set a new start on click.
if (mouse.clicked) {
  const x0 = state.xmin + ((mouse.x - padL) / (W - padL - padR)) * (state.xmax - state.xmin);
  state.restart(Math.max(state.xmin + 0.05, Math.min(state.xmax - 0.05, x0)));
}

// Advance one Newton step periodically.
state.timer++;
if (!state.done && state.timer % 45 === 0) {
  const x = state.x, d = state.fp(x);
  if (Math.abs(d) < 1e-6) { state.done = true; }
  else {
    const xn = x - state.f(x) / d;
    if (!isFinite(xn) || xn < state.xmin - 2 || xn > state.xmax + 2) { state.done = true; }
    else {
      state.x = xn; state.steps.push(xn);
      if (Math.abs(state.f(xn)) < 1e-4 || state.steps.length > 8) state.done = true;
    }
  }
  if (state.done) state.pause = 80;
}
if (state.done) { if (state.pause > 0) state.pause--; else state.restart(3.4); }

// Draw tangent steps: from (x_k, f(x_k)) down to (x_{k+1}, 0).
for (let k = 0; k < state.steps.length; k++) {
  const x = state.steps[k];
  const fx = state.f(x);
  const recent = (k >= state.steps.length - 2);
  // vertical guide from axis up to the curve
  ctx.strokeStyle = recent ? accent2 : faint;
  ctx.lineWidth = recent ? 2 : 1;
  ctx.setLineDash([4, 4]);
  ctx.beginPath(); ctx.moveTo(toPX(x), toPY(0)); ctx.lineTo(toPX(x), toPY(fx)); ctx.stroke();
  ctx.setLineDash([]);
  // tangent line to next root estimate
  const d = state.fp(x);
  if (Math.abs(d) > 1e-6) {
    const xn = x - fx / d;
    ctx.strokeStyle = recent ? accent2 : 'rgba(255,158,100,0.35)';
    ctx.lineWidth = recent ? 2 : 1;
    ctx.beginPath(); ctx.moveTo(toPX(x), toPY(fx)); ctx.lineTo(toPX(xn), toPY(0)); ctx.stroke();
  }
  // point on the curve
  ctx.fillStyle = good;
  ctx.beginPath(); ctx.arc(toPX(x), toPY(fx), recent ? 5 : 3, 0, 2 * Math.PI); ctx.fill();
}

// HUD.
ctx.fillStyle = text; ctx.font = '13px sans-serif';
ctx.fillText('x' + (state.steps.length - 1) + ' = ' + state.x.toFixed(6), 12, H - 10);
ctx.fillText('f(x) = ' + state.f(state.x).toExponential(2), 12, 26);
{{< /sketch >}}

{{< note kind="warning" title="Fast, but not foolproof" >}}
Newton's method is spectacular near the answer and treacherous far from it. If the curvature $f'(x)$ is near zero the step explodes; for a non-[[Convexity|convex]] loss a negative curvature can send it *uphill* toward a maximum. Practical optimizers blend Newton-like curvature with the reliability of gradient descent (quasi-Newton methods like BFGS), getting fast convergence without the fragility.
{{< /note >}}

## See also

- [[Gradient Descent]]
- [[Convexity]]
- [[Momentum]]
