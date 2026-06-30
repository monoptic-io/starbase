---
title: NP-Completeness
aliases: [np-complete, np complete, reduction, np-completeness]
tags: [computation]
summary: The hardest problems in NP, all interreducible — solve one efficiently and you solve them all.
weight: 90
---

# NP-Completeness

An **NP-complete** problem is one of the *hardest* problems in [[Complexity Class|NP]], in a precise sense: it is in NP, and **every** other problem in NP can be transformed into it by an efficient (polynomial-time) **reduction**. These problems are the keystones of [[P versus NP]]. Find a fast algorithm for *one* of them and you have, automatically, a fast algorithm for *all* of NP — proving P = NP. Prove that *one* of them needs exponential time and you have proved P ≠ NP. They stand or fall together.

There are thousands of them, scattered across logic, scheduling, [[Graph|graph theory]], number theory, and games — yet underneath they are, computationally, **the same problem wearing different costumes**.

## Reduction: the universal translator

The magic ingredient is the **polynomial-time reduction**. To reduce problem $A$ to problem $B$ is to write an efficient translator that turns any instance of $A$ into an instance of $B$ with the *same yes/no answer*. If such a translator exists, then "$B$ is easy" implies "$A$ is easy" — just translate and solve.

{{< note kind="key" title="What makes a problem NP-complete" >}}
A problem $B$ is NP-complete when **(1)** it is in NP — solutions are checkable in polynomial time — **and (2)** every problem in NP reduces to it. Cook and Levin proved in 1971 that boolean satisfiability (**SAT**) satisfies both. Every NP-complete problem found since has been pinned down by reducing a *known* one to it — a chain that all traces back to SAT.
{{< /note >}}

## A graph you can watch: 3-coloring

A classic NP-complete problem: can the vertices of a [[Graph]] be painted with just **three** colors so that no edge joins two same-colored vertices? Checking a proposed coloring is trivial — glance at every edge. *Finding* one means searching a space of $3^V$ colorings. The sketch below hunts for a valid 3-coloring of a wheel graph by trial: red edges are conflicts (endpoints clash); it keeps trying until every edge is satisfied, then locks in green.

{{< sketch height="380" caption="Searching for a 3-coloring of a wheel graph. Each pass recolors the vertices; edges with matching endpoints flash red as conflicts. When a conflict-free coloring is stumbled upon, the whole graph confirms green — easy to verify, hard to find." >}}
if (frame === 0) {
  state.n = 7;
  state.edges = [[0,1],[1,2],[2,3],[3,4],[4,5],[5,0],[6,0],[6,1],[6,2],[6,3],[6,4],[6,5]];
  state.colors = new Array(state.n).fill(0);
  state.solved = false;
  state.tNext = 0.4;
}
const cs = getComputedStyle(document.documentElement);
const accent = cs.getPropertyValue('--accent').trim() || '#5b9cff';
const good = cs.getPropertyValue('--good').trim() || '#46d39a';
const warn = cs.getPropertyValue('--warn').trim() || '#f2b347';
const text = cs.getPropertyValue('--text').trim() || '#e8edf5';
const border = cs.getPropertyValue('--border').trim() || 'rgba(255,255,255,0.18)';
const palette = [accent, good, warn];
const conflict = '#ff5470';

// try a random coloring each tick until one works, then hold
if (t > state.tNext) {
  if (state.solved) {
    state.solved = false;
    state.tNext = t + 0.4;
  } else {
    for (let i = 0; i < state.n; i++) state.colors[i] = Math.floor(Math.random() * 3);
    let bad = 0;
    for (let e = 0; e < state.edges.length; e++) {
      const a = state.edges[e][0], b = state.edges[e][1];
      if (state.colors[a] === state.colors[b]) bad++;
    }
    state.solved = bad === 0;
    state.tNext = t + (state.solved ? 2.4 : 0.4);
  }
}

ctx.clearRect(0, 0, W, H);
const cx = W / 2, cy = H * 0.54;
const R = Math.min(W, H) * 0.32;
const nodeR = Math.min(22, W * 0.06);
const pos = [];
for (let i = 0; i < 6; i++) {
  const a = -Math.PI / 2 + i * Math.PI / 3;
  pos.push({ x: cx + R * Math.cos(a), y: cy + R * Math.sin(a) });
}
pos.push({ x: cx, y: cy });   // hub

// edges
ctx.lineWidth = 2.5;
for (let e = 0; e < state.edges.length; e++) {
  const a = state.edges[e][0], b = state.edges[e][1];
  const clash = state.colors[a] === state.colors[b];
  ctx.strokeStyle = clash ? conflict : border;
  ctx.lineWidth = clash ? 3.5 : 2;
  ctx.beginPath();
  ctx.moveTo(pos[a].x, pos[a].y);
  ctx.lineTo(pos[b].x, pos[b].y);
  ctx.stroke();
}

// vertices
for (let i = 0; i < state.n; i++) {
  if (state.solved) {
    ctx.beginPath(); ctx.arc(pos[i].x, pos[i].y, nodeR + 5, 0, Math.PI * 2);
    ctx.fillStyle = good; ctx.globalAlpha = 0.18; ctx.fill(); ctx.globalAlpha = 1;
  }
  ctx.beginPath(); ctx.arc(pos[i].x, pos[i].y, nodeR, 0, Math.PI * 2);
  ctx.fillStyle = palette[state.colors[i]]; ctx.fill();
  ctx.lineWidth = 2; ctx.strokeStyle = state.solved ? good : '#0c1018'; ctx.stroke();
}

ctx.textAlign = 'center';
ctx.font = '14px ui-monospace, monospace';
ctx.fillStyle = state.solved ? good : text;
ctx.fillText(state.solved ? 'valid 3-coloring found — no edge conflicts'
                          : 'searching... red edges are conflicts', W / 2, H - 14);
{{< /sketch >}}

Other faces of the same beast include the **clique problem** (is there a set of $k$ mutually connected vertices?), the **traveling salesman** decision problem, and **SAT** itself. A polynomial reduction connects every one of them to all the others.

{{< note kind="note" title="NP-hard vs. NP-complete" >}}
A problem is **NP-hard** if everything in NP reduces to it — it is *at least* as hard as all of NP — but it need not itself be in NP (it might not even be decidable). An **NP-complete** problem is NP-hard *and* a member of NP. Optimization versions ("find the *largest* clique") are typically NP-hard; their yes/no decision versions ("is there a clique of size $k$?") are NP-complete.
{{< /note >}}

## So what do we do about them?

Since (almost certainly) no fast exact algorithm exists, real systems sidestep the wall: **approximation** algorithms that get provably close, **heuristics** like [[Simulated Annealing]] and [[Genetic Algorithm]]s that find good-enough answers, or solvers tuned for the easy instances that arise in practice. NP-completeness does not mean "give up" — it means "stop looking for a *perfect, fast, general* method, because by [[P versus NP]] it probably cannot exist."

{{< quiz question="A problem is shown to be NP-complete. What does this tell you?" options="It cannot be solved by any computer|It is in NP and every NP problem reduces to it|It can be solved in polynomial time|Its solutions cannot be verified" answer="2" explain="NP-complete = in NP (verifiable in polynomial time) AND NP-hard (every NP problem reduces to it). It is among the hardest problems in NP; a polynomial algorithm for it would solve all of NP, settling P vs NP." >}}

## See also

- [[P versus NP]]
- [[Complexity Class]]
- [[Graph]]
