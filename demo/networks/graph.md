---
title: Graph
aliases: [graph theory, network, nodes and edges]
tags: [networks, foundations]
summary: A set of nodes joined by edges — the minimal mathematical object for describing relationships of any kind.
weight: 10
---

# Graph

A **graph** is the barest possible model of connection: a set of **nodes** (also called *vertices*) and a set of **edges**, each edge joining a pair of nodes. That is all. A graph carries no geometry, no distance, no embedding in space — only the relation *these two are connected*. Its power comes precisely from that austerity: anything you can phrase as "things and the links between them" is a graph.

Formally we write $G = (V, E)$, where $V$ is the set of vertices and $E \subseteq V \times V$ the set of edges. The number of edges meeting a node is its **degree**. Two nodes joined by an edge are **neighbors**, and a sequence of edges leading from one node to another is a **path**.

## Flavors of graph

The single picture splits into a small family of variants, chosen to match what you are modeling:

- **Undirected vs. directed.** A friendship is mutual, so its edge has no arrow — an *undirected* graph. A hyperlink or a one-way street points one way, giving a *directed* graph (a **digraph**), where an edge from $a$ to $b$ need not imply one from $b$ to $a$.
- **Unweighted vs. weighted.** Sometimes an edge simply exists. Other times it carries a number — a distance, a cost, a capacity, a strength. A *weighted* graph attaches a value to each edge, and most of the interesting optimization problems ([[Dijkstra's Algorithm]], [[Minimum Spanning Tree]]) live here.
- **Sparse vs. dense.** A graph with $n$ nodes can have up to $\binom{n}{2}$ edges. Real networks are usually *sparse* — each node connects to only a handful of others — which is what makes algorithms that run in $O(V+E)$ time so valuable.

{{< note kind="note" title="Two ways to store a graph" >}}
An **adjacency matrix** is an $n \times n$ grid whose $(i,j)$ entry marks the edge $i \to j$; it is simple but wastes space on sparse graphs. An **adjacency list** stores, for each node, just its neighbors — compact, and the natural input for [[Breadth-First Search]] and friends. The matrix view, though, is what links graphs to [[Linear Algebra]] and powers [[PageRank]].
{{< /note >}}

## A graph wants to lay itself out

A graph has no built-in shape, so to *draw* one we invent positions. A **force-directed layout** treats the graph as a physical system: every node is a charged particle that pushes all the others away, and every edge is a spring that pulls its two endpoints together. Release the system and it relaxes into a configuration where clusters spread out, connected nodes sit close, and the structure becomes legible. It is the same balance of repulsion and attraction that organizes a hanging mobile.

Drag any node below — the rest of the network will flex and re-settle around it.

{{< sketch height="420" caption="A small graph finding its own shape. Nodes repel like charges; edges pull like springs; the whole thing relaxes to equilibrium. Click and drag any node to disturb it." >}}
if (frame === 0) {
  const N = 24;
  state.nodes = [];
  for (let i = 0; i < N; i++) {
    state.nodes.push({
      x: W / 2 + (Math.random() - 0.5) * W * 0.5,
      y: H / 2 + (Math.random() - 0.5) * H * 0.5,
      vx: 0, vy: 0
    });
  }
  // a spanning tree guarantees connectivity, plus a few extra edges for loops
  state.edges = [];
  for (let i = 1; i < N; i++) state.edges.push([i, Math.floor(Math.random() * i)]);
  for (let k = 0; k < 7; k++) {
    const a = Math.floor(Math.random() * N), b = Math.floor(Math.random() * N);
    if (a !== b) state.edges.push([a, b]);
  }
  state.drag = -1;
}
const cs = getComputedStyle(document.documentElement);
const accent = cs.getPropertyValue('--accent').trim() || '#5b9cff';
const accent2 = cs.getPropertyValue('--accent-2').trim() || '#ff9e64';
const border = cs.getPropertyValue('--border').trim() || 'rgba(255,255,255,0.18)';
const nodes = state.nodes, edges = state.edges, N = nodes.length;

// pick up the nearest node on press
if (mouse.down && state.drag < 0) {
  let best = -1, bd = 30 * 30;
  for (let i = 0; i < N; i++) {
    const dx = nodes[i].x - mouse.x, dy = nodes[i].y - mouse.y, d = dx * dx + dy * dy;
    if (d < bd) { bd = d; best = i; }
  }
  state.drag = best;
}
if (!mouse.down) state.drag = -1;

// repulsion between every pair
for (let i = 0; i < N; i++) {
  for (let j = i + 1; j < N; j++) {
    let dx = nodes[i].x - nodes[j].x, dy = nodes[i].y - nodes[j].y;
    let d2 = dx * dx + dy * dy + 0.01, d = Math.sqrt(d2);
    const f = 2400 / d2;
    const fx = dx / d * f, fy = dy / d * f;
    nodes[i].vx += fx; nodes[i].vy += fy;
    nodes[j].vx -= fx; nodes[j].vy -= fy;
  }
}
// spring attraction along edges
const L = 78;
for (let e = 0; e < edges.length; e++) {
  const a = edges[e][0], b = edges[e][1];
  let dx = nodes[b].x - nodes[a].x, dy = nodes[b].y - nodes[a].y;
  let d = Math.sqrt(dx * dx + dy * dy) + 0.01;
  const f = (d - L) * 0.02;
  const fx = dx / d * f, fy = dy / d * f;
  nodes[a].vx += fx; nodes[a].vy += fy;
  nodes[b].vx -= fx; nodes[b].vy -= fy;
}
// gentle pull toward center so nothing drifts off-screen
for (let i = 0; i < N; i++) {
  nodes[i].vx += (W / 2 - nodes[i].x) * 0.004;
  nodes[i].vy += (H / 2 - nodes[i].y) * 0.004;
}
// integrate with damping + velocity clamp
for (let i = 0; i < N; i++) {
  if (i === state.drag) { nodes[i].x = mouse.x; nodes[i].y = mouse.y; nodes[i].vx = 0; nodes[i].vy = 0; continue; }
  nodes[i].vx *= 0.82; nodes[i].vy *= 0.82;
  const vmax = 16;
  nodes[i].vx = Math.max(-vmax, Math.min(vmax, nodes[i].vx));
  nodes[i].vy = Math.max(-vmax, Math.min(vmax, nodes[i].vy));
  nodes[i].x += nodes[i].vx; nodes[i].y += nodes[i].vy;
  nodes[i].x = Math.max(12, Math.min(W - 12, nodes[i].x));
  nodes[i].y = Math.max(12, Math.min(H - 12, nodes[i].y));
}
// draw
ctx.clearRect(0, 0, W, H);
ctx.strokeStyle = border;
ctx.lineWidth = 1.2;
for (let e = 0; e < edges.length; e++) {
  const a = nodes[edges[e][0]], b = nodes[edges[e][1]];
  ctx.beginPath(); ctx.moveTo(a.x, a.y); ctx.lineTo(b.x, b.y); ctx.stroke();
}
for (let i = 0; i < N; i++) {
  const r = i === state.drag ? 9 : 6;
  ctx.fillStyle = i === state.drag ? accent2 : accent;
  ctx.beginPath(); ctx.arc(nodes[i].x, nodes[i].y, r, 0, 7); ctx.fill();
}
{{< /sketch >}}

This same relaxation idea is how diagrams of [[Complex Systems]] — protein interactions, social circles, the internet — are drawn so that their hidden community structure becomes visible.

## See also

- [[Breadth-First Search]]
- [[Centrality]]
- [[Complex Systems]]
