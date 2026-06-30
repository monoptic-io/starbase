---
title: Turing Machine
aliases: [turing machine]
tags: [computation]
summary: An idealized machine of tape, head, and rules that defines exactly what it means for something to be computable.
weight: 10
---

# Turing Machine

A **Turing machine** is the simplest device powerful enough to compute anything a computer can. Alan Turing dreamed it up in 1936 — years before working hardware — not to *build* but to *define*: it pins down, with mathematical precision, what the word **computable** means. Strip a computer down to its absolute essentials and this is what remains.

The whole machine is just four parts:

- an **infinite tape** divided into cells, each holding a symbol (say `0`, `1`, or blank);
- a **head** that sits over one cell, able to read and rewrite it;
- a **state** drawn from a small finite set (the machine's entire "memory" beyond the tape);
- a **transition table** of rules of the form *"in state $q$ reading symbol $s$: write $s'$, move left or right, go to state $q'$."*

That is the entire universe of the machine. Astonishingly, it is enough.

## Watch one run

The machine below is a **binary counter**. Its head sits at the rightmost bit and runs the schoolbook increment rule: a `1` becomes `0` and the head carries left; the first `0` it meets becomes `1` and the machine halts. Then it counts again. Every sweep of the head to the left *is* a carry propagating — the same logic in your calculator, laid bare.

{{< sketch height="300" caption="A Turing machine incrementing a binary number. State INC scans left flipping 1→0 on each carry; the first 0 flips to 1 and the machine reaches DONE. The head pointer turns green at the halt." >}}
if (frame === 0) {
  state.N = 11;
  state.tape = new Array(state.N).fill(0);
  state.head = state.N - 1;
  state.mode = 'INC';
  state.tNext = 0.6;
}
const cs = getComputedStyle(document.documentElement);
const accent = cs.getPropertyValue('--accent').trim() || '#5b9cff';
const accent2 = cs.getPropertyValue('--accent-2').trim() || '#c08bff';
const good = cs.getPropertyValue('--good').trim() || '#46d39a';
const text = cs.getPropertyValue('--text').trim() || '#e8edf5';
const border = cs.getPropertyValue('--border').trim() || 'rgba(255,255,255,0.15)';
const faint = cs.getPropertyValue('--text-faint').trim() || 'rgba(255,255,255,0.4)';

// step the machine on a timer
if (t > state.tNext) {
  state.tNext = t + 0.5;
  if (state.mode === 'INC') {
    if (state.tape[state.head] === 0) {
      state.tape[state.head] = 1;
      state.mode = 'DONE';
    } else {
      state.tape[state.head] = 0;
      state.head -= 1;
      if (state.head < 0) {            // overflowed past the top bit: roll over
        state.tape = new Array(state.N).fill(0);
        state.head = state.N - 1;
      }
    }
  } else {
    state.head = state.N - 1;          // rewind and count again
    state.mode = 'INC';
  }
}

ctx.clearRect(0, 0, W, H);
const N = state.N;
const cw = Math.min(46, (W - 40) / N);
const x0 = (W - cw * N) / 2;
const cy = H * 0.52;
const half = cw / 2;

ctx.lineWidth = 2;
ctx.textAlign = 'center';
ctx.textBaseline = 'middle';
ctx.font = (cw * 0.5) + 'px ui-monospace, monospace';
for (let i = 0; i < N; i++) {
  const x = x0 + i * cw;
  const isHead = i === state.head;
  ctx.fillStyle = state.tape[i] ? (isHead ? accent : 'rgba(91,156,255,0.18)') : 'rgba(255,255,255,0.03)';
  ctx.fillRect(x, cy - half, cw - 2, cw - 2);
  ctx.strokeStyle = isHead ? accent2 : border;
  ctx.strokeRect(x, cy - half, cw - 2, cw - 2);
  ctx.fillStyle = state.tape[i] ? text : faint;
  ctx.fillText(state.tape[i], x + half - 1, cy);
}

// head pointer above the active cell
const hx = x0 + state.head * cw + half - 1;
const hy = cy - half - 6;
ctx.fillStyle = state.mode === 'DONE' ? good : accent2;
ctx.beginPath();
ctx.moveTo(hx, hy);
ctx.lineTo(hx - 7, hy - 12);
ctx.lineTo(hx + 7, hy - 12);
ctx.closePath();
ctx.fill();

let val = 0;
for (let i = 0; i < N; i++) val = val * 2 + state.tape[i];
ctx.font = '14px ui-monospace, monospace';
ctx.fillStyle = state.mode === 'DONE' ? good : accent;
ctx.fillText('state: ' + state.mode, W / 2, cy + cw + 4);
ctx.fillStyle = faint;
ctx.fillText('tape value = ' + val, W / 2, cy + cw + 26);
{{< /sketch >}}

Two states (`INC`, `DONE`) and a one-line rule produce unbounded counting. Give it more states and a richer table and the *same machinery* can sort, search, factor, or simulate any other computer.

## Why this toy is the whole point

The Turing machine matters precisely because it is so weak-looking. If even this can be made to compute a function, the function is **computable**; if no Turing machine can, no machine can — that universality is the content of the [[Church–Turing Thesis]]. A single **universal** Turing machine can read another machine's table off its tape and imitate it, which is exactly what a stored-program computer does: code is just data on the tape.

{{< note kind="note" title="Tape vs. no tape" >}}
Take the tape away and forbid the head from writing, and you get a [[Finite Automaton]] — far weaker, able only to recognize *regular* patterns. The rewritable, unbounded tape is the entire difference between a pattern-matcher and a universal computer. Memory is power.
{{< /note >}}

This same definition is what makes impossibility provable. Because *computable* now has an exact meaning, one can show certain problems have **no** machine at all — most famously the [[Halting Problem]]: no Turing machine can decide whether an arbitrary machine eventually halts. And the local-rule universality on display here is not unique to tapes; a [[Cellular Automaton]] such as [[Conway's Game of Life]] is Turing-complete from nothing but cells flipping on and off.

{{< quiz question="What single feature gives a Turing machine its power over a finite automaton?" options="More states|A faster clock|An unbounded, rewritable tape for memory|Random access to any cell" answer="3" explain="A finite automaton has only its current state to remember the past. The Turing machine's infinite, rewritable tape is unbounded external memory — that is the entire jump from regular languages to universal computation." >}}

## See also

- [[Finite Automaton]]
- [[Church–Turing Thesis]]
- [[Halting Problem]]
- [[Cellular Automaton]]
