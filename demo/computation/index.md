---
title: Computability & Complexity
aliases: [computability and complexity, theory of computation, computation]
tags: [computation]
summary: What a machine can compute in principle, and how much time and space it takes in practice — from the Turing machine to the million-dollar P versus NP question.
weight: 80
---

# Computability & Complexity

Every other section of this field guide *runs a rule forward*: a [[Pendulum]] swings, a [[Wave]] propagates, a [[Markov Chain]] wanders. This section asks the prior question — **what does it even mean to compute a rule, and what are the limits?** It splits cleanly into two halves that grew up together.

- **Computability** — the study of what can be computed *at all*, given unlimited time and memory. Some problems have no algorithm, ever. This is where the [[Turing Machine]], the [[Halting Problem]], and [[Decidability]] live.
- **Complexity** — among the problems we *can* solve, how expensive are they? Measured in time and space as the input grows, this is the world of [[Big-O Notation]], [[Complexity Class]]es, and the famous [[P versus NP]] divide.

{{< note kind="key" title="The two great questions" >}}
**Can it be done?** (computability) and **can it be done *fast*?** (complexity). The first was settled in the 1930s with a resounding *not always*. The second is still wide open — and a correct answer to [[P versus NP]] is worth a literal million dollars.
{{< /note >}}

## The big arc

In 1936, before any electronic computer existed, Alan Turing imagined an absurdly simple machine — a tape, a head, a handful of rules — and proved two things at once: that this toy captures **everything** mechanically computable (the [[Church–Turing Thesis]]), and that even it cannot solve every problem (the [[Halting Problem]]). Computation has a hard ceiling, and it was found before the hardware.

Once you accept that machines are limited, the next question is *cost*. A problem may be solvable yet hopeless: a method that takes $2^n$ steps is useless at $n = 100$. Sorting these problems by their appetite for time and memory gives the [[Complexity Class]] hierarchy, whose central mystery — is checking an answer really easier than finding one? — is [[P versus NP]], anchored by the [[NP-Completeness|NP-complete]] problems that stand or fall together.

{{< sketch height="300" caption="A landscape of problems. Green = solvable and fast (P). Amber = solvable but possibly slow (NP and beyond). Grey = no algorithm exists at all (undecidable). The boundaries are the whole story." >}}
if (frame === 0) {
  state.cells = [[ ]];
  state.cells = null;
}
const cs = getComputedStyle(document.documentElement);
const good = cs.getPropertyValue('--good').trim() || '#46d39a';
const warn = cs.getPropertyValue('--warn').trim() || '#f2b347';
const faint = cs.getPropertyValue('--text-faint').trim() || 'rgba(255,255,255,0.35)';
const text = cs.getPropertyValue('--text').trim() || '#e8edf5';
ctx.clearRect(0, 0, W, H);
const cx = W / 2, cy = H * 0.52;
const rings = [
  { r: Math.min(W, H) * 0.16, col: good, label: 'P  — fast' },
  { r: Math.min(W, H) * 0.28, col: warn, label: 'NP / PSPACE / EXP' },
  { r: Math.min(W, H) * 0.40, col: faint, label: 'undecidable' }
];
for (let i = rings.length - 1; i >= 0; i--) {
  const ring = rings[i];
  const pulse = 1 + 0.015 * Math.sin(t * 1.2 + i);
  ctx.beginPath();
  ctx.arc(cx, cy, ring.r * pulse, 0, Math.PI * 2);
  ctx.fillStyle = ring.col;
  ctx.globalAlpha = 0.12;
  ctx.fill();
  ctx.globalAlpha = 0.9;
  ctx.lineWidth = 2;
  ctx.strokeStyle = ring.col;
  ctx.stroke();
}
ctx.globalAlpha = 1;
ctx.textAlign = 'center';
ctx.font = '13px ui-monospace, monospace';
ctx.fillStyle = good;
ctx.fillText('P', cx, cy + 4);
ctx.fillStyle = warn;
ctx.fillText('NP · PSPACE · EXP', cx, cy - rings[0].r - 14);
ctx.fillStyle = faint;
ctx.fillText('undecidable — no algorithm at all', cx, cy - rings[1].r - 14);
{{< /sketch >}}

## The pages in this section

{{< columns count="2" >}}
**Computability**
- [[Turing Machine]] — tape, head, rules: the definition of *computable*.
- [[Finite Automaton]] — a memoryless cousin that recognizes regular patterns.
- [[Church–Turing Thesis]] — why every reasonable model computes the same things.
- [[Halting Problem]] — the first proven-impossible problem.
- [[Decidability]] — the frontier between solvable and forever out of reach.

**Complexity**
- [[Big-O Notation]] — the language of *how fast it grows*.
- [[Complexity Class]] — P, NP, PSPACE, EXP — sorting problems by resources.
- [[P versus NP]] — is finding as easy as checking?
- [[NP-Completeness]] — the hardest problems in NP, joined at the hip.
{{< /columns >}}

## Threads to the rest of the field

Computation is not sealed off. The shortest-path engine of [[Dijkstra's Algorithm]] is a complexity story about a [[Graph]]; the [[NP-Completeness|NP-complete]] problems are graph-coloring and clique puzzles. And the deepest link runs to [[Information Theory]]: **[[Kolmogorov Complexity]]** measures the size of the *shortest program* that prints a string — a computational definition of information that turns out to be [[Decidability|undecidable]]. The endlessly intricate patterns of a [[Cellular Automaton]] like [[Conway's Game of Life]] are Turing-complete too: simple local rules, full computational power.

Start with the [[Turing Machine]], or jump to the open problem that defines the field: [[P versus NP]].
