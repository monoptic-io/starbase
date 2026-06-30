---
title: Genetic Algorithm
aliases: [genetic algorithm, evolutionary algorithm]
tags: [optimization]
summary: Optimize with no gradient at all — keep a population of candidate solutions, let the fittest reproduce with mutation and crossover, and watch quality climb.
weight: 70
---

# Genetic Algorithm

Some landscapes are too rugged, discrete, or mysterious to differentiate — there is no gradient to follow. A **genetic algorithm** sidesteps that entirely by imitating natural selection. Instead of one point creeping downhill, you maintain a whole **population** of candidate solutions and breed them:

1. **Evaluate** each individual's *fitness* (how good a solution it is).
2. **Select** parents, favoring the fitter ones.
3. **Crossover** — combine two parents into a child, mixing their "genes."
4. **Mutate** — randomly tweak a few genes to inject new variety.
5. Repeat. The next generation is, on average, a little fitter than the last.

There is no derivative anywhere — only *variation* (mutation and crossover) and *selection* (survival of the fittest). Over generations, that loop reliably climbs toward good solutions, and it cheerfully handles problems gradient descent cannot touch.

## Evolving toward a target phrase

The classic demonstration: evolve random gibberish into a target sentence. Each individual is a string; its fitness is simply how many letters it gets right. No string is ever *designed* — selection and mutation alone drive the population from noise to the exact phrase. Watch the **best fitness climb** generation by generation; matched letters turn green.

**Click to start a fresh population** of random strings.

{{< sketch height="400" caption="A genetic algorithm evolving random strings toward a target phrase. Fitness = number of correct letters; selection favors the fittest, crossover mixes parents, mutation adds variety. The curve tracks best fitness rising to 100%. Click to reseed." >}}
if (frame === 0 || !state.pop) {
  state.target = "GENETIC ALGORITHM";
  state.alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ ";
  state.N = 180;
  state.mut = 0.045;
  const L = state.target.length;
  state.rand = function () {
    let s = "";
    for (let i = 0; i < L; i++) s += state.alphabet[Math.floor(Math.random() * state.alphabet.length)];
    return s;
  };
  state.fit = function (s) {
    let f = 0;
    for (let i = 0; i < L; i++) if (s[i] === state.target[i]) f++;
    return f;
  };
  state.seed = function () {
    state.pop = [];
    for (let i = 0; i < state.N; i++) state.pop.push(state.rand());
    state.gen = 0;
    state.history = [];
    state.best = state.pop[0];
    state.bestFit = state.fit(state.best);
    state.timer = 0;
    state.solved = false;
    state.pause = 0;
  };
  state.seed();
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

if (mouse.clicked) state.seed();

const L = state.target.length;

// One generation every few frames.
state.timer++;
if (!state.solved && state.timer % 5 === 0) {
  // Evaluate + track best.
  let best = state.pop[0], bf = state.fit(best);
  for (let i = 1; i < state.N; i++) {
    const f = state.fit(state.pop[i]);
    if (f > bf) { bf = f; best = state.pop[i]; }
  }
  state.best = best; state.bestFit = bf;
  state.history.push(bf / L);
  if (state.history.length > 240) state.history.shift();
  if (bf === L) { state.solved = true; state.pause = 110; }
  else {
    const tournament = function () {
      let a = state.pop[Math.floor(Math.random() * state.N)];
      for (let k = 0; k < 3; k++) {
        const b = state.pop[Math.floor(Math.random() * state.N)];
        if (state.fit(b) > state.fit(a)) a = b;
      }
      return a;
    };
    const next = [best]; // elitism
    while (next.length < state.N) {
      const p1 = tournament(), p2 = tournament();
      const cut = Math.floor(Math.random() * L);
      let child = "";
      for (let i = 0; i < L; i++) {
        let c = (i < cut) ? p1[i] : p2[i];
        if (Math.random() < state.mut) c = state.alphabet[Math.floor(Math.random() * state.alphabet.length)];
        child += c;
      }
      next.push(child);
    }
    state.pop = next;
    state.gen++;
  }
}
if (state.solved) { if (state.pause > 0) state.pause--; else state.seed(); }

// Draw the best individual, highlighting correct letters.
ctx.font = '22px monospace';
const chW = ctx.measureText('M').width;
const total = L * chW;
let sx = (W - total) / 2;
const sy = 64;
for (let i = 0; i < L; i++) {
  const ch = state.best[i];
  ctx.fillStyle = (ch === state.target[i]) ? good : faint;
  ctx.fillText(ch === ' ' ? '·' : ch, sx + i * chW, sy);
}
ctx.font = '12px sans-serif';
ctx.fillStyle = faint;
ctx.fillText('target:  ' + state.target, sx, 28);

// Fitness curve.
const gx = 50, gy = 110, gw = W - 70, gh = H - gy - 36;
ctx.strokeStyle = border; ctx.lineWidth = 1;
ctx.strokeRect(gx, gy, gw, gh);
ctx.fillStyle = faint; ctx.font = '11px sans-serif';
ctx.fillText('100%', 12, gy + 10);
ctx.fillText('0%', 20, gy + gh);
ctx.fillText('generations →', gx + 6, gy + gh + 18);
ctx.strokeStyle = accent2; ctx.lineWidth = 2.5;
ctx.beginPath();
for (let i = 0; i < state.history.length; i++) {
  const px = gx + (i / Math.max(1, 239)) * gw;
  const py = gy + (1 - state.history[i]) * gh;
  if (i === 0) ctx.moveTo(px, py); else ctx.lineTo(px, py);
}
ctx.stroke();

// HUD.
ctx.fillStyle = text; ctx.font = '13px sans-serif';
ctx.fillText('generation ' + state.gen, gx, gy - 8);
ctx.fillStyle = good;
ctx.fillText('best fitness ' + state.bestFit + ' / ' + L, gx + 130, gy - 8);
{{< /sketch >}}

{{< note kind="key" title="Exploration vs exploitation" >}}
Mutation is **exploration** (try new things); selection is **exploitation** (keep what works). Too much mutation and the population never settles; too little and it stagnates in a local optimum — the same tension [[Simulated Annealing]] manages with temperature. Crossover adds a twist unique to populations: it can splice a good "head" from one solution onto a good "tail" from another.
{{< /note >}}

## See also

- [[Predator–Prey Dynamics]]
- [[Monte Carlo Method]]
- [[Simulated Annealing]]
