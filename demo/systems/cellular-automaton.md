---
title: Cellular Automaton
aliases: [cellular automata, cellular automaton, elementary cellular automaton]
tags: [systems, emergence, discrete]
summary: A grid of cells, each updated by a simple rule that looks only at its neighbors — discrete dynamics that generate astonishing complexity.
weight: 40
---

# Cellular Automaton

A **cellular automaton** is a [[Dynamical System]] stripped to its barest discrete form: a grid of cells, each in one of a few states, all updated *simultaneously* by a single rule that looks only at a cell's immediate neighbors. Time ticks in steps, space is a lattice, and states are discrete — yet from this austere setup pours pattern, structure, and even universal computation. Cellular automata are the cleanest demonstration of the section's thesis: **local rules, global complexity**.

Because the update is a map applied over and over, cellular automata live squarely in the world of [[Flows and Maps]] — they are maps whose state space is an entire grid.

## Elementary cellular automata

The simplest family lives on a one-dimensional row of cells, each either **on** or **off**. A cell's next state depends on just three cells: itself and its two neighbors. There are $2^3 = 8$ possible neighborhoods, and a rule assigns an output to each — so there are exactly $2^8 = 256$ such rules, numbered 0–255 by Wolfram's convention (read the eight outputs as a binary number).

Most are dull. A few are extraordinary. **Rule 30** is the star: from a single on-cell it generates a stream of structure so disordered that it has been used as a random-number generator, yet it is produced by a rule you can state in one line.

{{< note kind="key" title="Rule 30 in one line" >}}
Let $l, c, r$ be the left, center, and right neighbors (each 0 or 1). The new center cell is

$$c' = l \oplus (c \lor r),$$

where $\oplus$ is XOR and $\lor$ is OR. Start from a single 1 and stack each new row beneath the last — the triangle below is what you get.
{{< /note >}}

## Watch Rule 30 build itself

Each row is computed from the one above it and drawn beneath, top to bottom. The left edge marches in a clean periodic stripe while the right dissolves into apparent randomness — order and chaos from the *same* three-cell rule.

{{< sketch height="380" caption="Rule 30 grown from a single seed cell. Rows fill top-to-bottom; when the screen fills it reseeds." >}}
if (frame === 0 || !state.rows) {
  state.cols = 121;
  state.cell = W / state.cols;
  state.maxRows = Math.floor(H / state.cell);
  // seed: single 1 in the middle
  let row = new Array(state.cols).fill(0);
  row[Math.floor(state.cols / 2)] = 1;
  state.rows = [row];
  state.shown = 1;
  state.acc = 0;
}
// advance roughly 1.5 rows per frame for a steady reveal
state.acc += 1.5;
while (state.acc >= 1 && state.rows.length < state.maxRows) {
  state.acc -= 1;
  const prev = state.rows[state.rows.length - 1];
  const n = state.cols;
  const next = new Array(n).fill(0);
  for (let i = 0; i < n; i++) {
    const l = prev[(i - 1 + n) % n];
    const c = prev[i];
    const r = prev[(i + 1) % n];
    next[i] = l ^ (c | r);
  }
  state.rows.push(next);
}
if (state.rows.length >= state.maxRows) state.acc = 0;
// draw
ctx.fillStyle = "#0f1020";
ctx.fillRect(0, 0, W, H);
const s = state.cell;
for (let y = 0; y < state.rows.length; y++) {
  const row = state.rows[y];
  // hue drifts down the triangle for a bit of life
  const hue = 190 + (y / state.maxRows) * 80;
  ctx.fillStyle = "hsl(" + hue + ",75%,62%)";
  for (let x = 0; x < state.cols; x++) {
    if (row[x]) ctx.fillRect(x * s, y * s, s + 0.5, s + 0.5);
  }
}
// reseed once the canvas is full and we've paused a moment
if (state.rows.length >= state.maxRows) {
  state.pause = (state.pause || 0) + 1;
  if (state.pause > 90) { state.rows = null; state.pause = 0; }
}
{{< /sketch >}}

## Classes of behavior

Survey all 256 rules and their long-term behavior sorts into four rough classes: settling to a uniform state, settling to stable or periodic stripes, producing chaotic noise, and — rarest and most interesting — generating *localized structures* that move and interact. That fourth class is where computation hides. The two-dimensional automaton [[Conway's Game of Life]] belongs to it, and its drifting "gliders" are the proof that a grid of bits can carry information across space.

## See also

- [[Conway's Game of Life]]
- [[Flows and Maps]]
- [[Complex Systems]]
