---
title: Simulated Annealing
aliases: [annealing]
tags: [optimization]
summary: Escape local minima by sometimes stepping uphill — with a probability set by a temperature that slowly cools from bold exploration to careful descent.
weight: 50
---

# Simulated Annealing

[[Gradient Descent|Gradient descent]] only ever goes downhill, so it gets trapped in the first valley it falls into. **Simulated annealing** borrows a trick from metallurgy to escape: heat the system so it can jump around freely, then cool it slowly so it settles into a deep, low-energy configuration. Concretely, at each step you propose a random move. If it lowers the loss (energy $E$), you always take it. If it *raises* the loss by $\Delta E$, you take it anyway — but only with probability

$$P(\text{accept uphill}) = e^{-\Delta E / T}.$$

When the **temperature** $T$ is high, even big uphill jumps are likely, so the walker roams the whole [[Loss Landscape|landscape]], hopping out of any trap. As $T$ falls, uphill moves become rarer and rarer, and the walker is squeezed into the deepest basin it has found. This accept-with-probability rule makes the search a [[Markov Chain]] (the Metropolis algorithm), and choosing moves at random makes it a [[Monte Carlo Method|Monte Carlo method]].

## Cooling into the global minimum

Below, a walker explores a rugged one-dimensional landscape. At high temperature it leaps boldly, even uphill, sampling far-apart valleys. As the **temperature bar** on the right drains, its jumps shrink and uphill moves all but stop, until it settles — usually into the deepest valley, not just the nearest. The best point found so far is marked.

**Click anywhere to reheat** and watch a fresh anneal from a random start.

{{< sketch height="380" caption="Simulated annealing on a bumpy 1-D energy landscape. High temperature = big, frequent uphill-tolerant jumps (exploration); as it cools, the walker funnels into a deep minimum. Click to re-anneal from a random point." >}}
if (frame === 0 || !state.energy) {
  state.umin = -3; state.umax = 3;
  state.energy = function (u) { return 0.3 * u * u + Math.sin(3 * u) + 0.6 * Math.cos(5.3 * u); };
  // Sample range for vertical scaling.
  let lo = Infinity, hi = -Infinity;
  for (let i = 0; i <= 600; i++) {
    const u = state.umin + (i / 600) * (state.umax - state.umin);
    const e = state.energy(u);
    if (e < lo) lo = e; if (e > hi) hi = e;
  }
  state.Emin = lo; state.Emax = hi;
  state.start = function () {
    state.u = state.umin + Math.random() * (state.umax - state.umin);
    state.T = 3.0;
    state.best = state.u; state.bestE = state.energy(state.u);
    state.flash = 0;
    state.done = false;
    state.cool = 0;
  };
  state.start();
}

const cs = getComputedStyle(document.documentElement);
const accent = cs.getPropertyValue('--accent').trim() || '#5b9cff';
const accent2 = cs.getPropertyValue('--accent-2').trim() || '#ff9e64';
const good = cs.getPropertyValue('--good').trim() || '#9ece6a';
const warn = cs.getPropertyValue('--warn').trim() || '#e0af68';
const text = cs.getPropertyValue('--text').trim() || '#e8ecf5';
const faint = cs.getPropertyValue('--text-faint').trim() || 'rgba(232,236,245,0.45)';

ctx.clearRect(0, 0, W, H);
ctx.fillStyle = '#0b0e16';
ctx.fillRect(0, 0, W, H);

const yTop = 44, yBot = H - 34, plotW = W - 70;
const toPX = function (u) { return 10 + (u - state.umin) / (state.umax - state.umin) * (plotW - 20); };
const toPY = function (e) {
  const t = (e - state.Emin) / (state.Emax - state.Emin);
  return yTop + (1 - t) * (yBot - yTop);
};

// Draw the landscape curve.
ctx.strokeStyle = accent; ctx.lineWidth = 2.5;
ctx.beginPath();
for (let i = 0; i <= 240; i++) {
  const u = state.umin + (i / 240) * (state.umax - state.umin);
  const px = toPX(u), py = toPY(state.energy(u));
  if (i === 0) ctx.moveTo(px, py); else ctx.lineTo(px, py);
}
ctx.stroke();

if (mouse.clicked) state.start();

// Annealing steps (a few proposals per frame).
if (!state.done) {
  for (let s = 0; s < 2; s++) {
    const step = 0.04 + state.T * 0.35;
    const uNew = state.u + (Math.random() * 2 - 1) * step;
    const uc = Math.max(state.umin, Math.min(state.umax, uNew));
    const dE = state.energy(uc) - state.energy(state.u);
    if (dE < 0 || Math.random() < Math.exp(-dE / Math.max(state.T, 1e-4))) {
      if (dE > 0) state.flash = 1;
      state.u = uc;
      const e = state.energy(uc);
      if (e < state.bestE) { state.bestE = e; state.best = uc; }
    }
    state.T *= 0.992;
    state.cool++;
  }
  if (state.T < 0.02) { state.done = true; state.pause = 90; }
}
if (state.done) { if (state.pause > 0) state.pause--; else state.start(); }

// Best-so-far marker.
ctx.fillStyle = good;
ctx.beginPath(); ctx.arc(toPX(state.best), toPY(state.bestE), 5, 0, 2 * Math.PI); ctx.fill();
ctx.strokeStyle = good; ctx.lineWidth = 1;
ctx.beginPath(); ctx.moveTo(toPX(state.best), toPY(state.bestE)); ctx.lineTo(toPX(state.best), yBot + 6); ctx.stroke();

// The walker.
const wx = toPX(state.u), wy = toPY(state.energy(state.u));
ctx.fillStyle = state.flash > 0 ? warn : accent2;
ctx.shadowColor = ctx.fillStyle; ctx.shadowBlur = 14;
ctx.beginPath(); ctx.arc(wx, wy, 7, 0, 2 * Math.PI); ctx.fill();
ctx.shadowBlur = 0;
if (state.flash > 0) state.flash -= 0.08;

// Temperature bar.
const barX = W - 40, barY = 30, barH = H - 70;
ctx.strokeStyle = faint; ctx.lineWidth = 1;
ctx.strokeRect(barX, barY, 16, barH);
const frac = Math.max(0, Math.min(1, state.T / 3.0));
ctx.fillStyle = frac > 0.4 ? warn : good;
ctx.fillRect(barX, barY + (1 - frac) * barH, 16, frac * barH);
ctx.fillStyle = text; ctx.font = '11px sans-serif';
ctx.fillText('T', barX + 3, barY - 8);

// HUD.
ctx.fillStyle = text; ctx.font = '13px sans-serif';
ctx.fillText('T = ' + state.T.toFixed(3) + (state.done ? '  (cooled)' : ''), 12, 22);
ctx.fillStyle = good;
ctx.fillText('best E = ' + state.bestE.toFixed(3), 12, H - 12);
{{< /sketch >}}

{{< note kind="key" title="Why slow cooling matters" >}}
Cool too fast and the walker freezes wherever it happens to be — exactly the local-minimum trap you were trying to avoid. Cool slowly enough and, in the idealized limit, annealing is *guaranteed* to find the global minimum. The schedule $T_k$ is the entire art: it trades runtime for the quality of the answer.
{{< /note >}}

## See also

- [[Monte Carlo Method]]
- [[Markov Chain]]
- [[Loss Landscape]]
