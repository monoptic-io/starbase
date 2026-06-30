---
title: Centrality
aliases: [degree centrality, betweenness, importance]
tags: [networks, metrics]
summary: A family of measures that score which nodes matter most in a network — by connections, by bridging, or by closeness.
weight: 70
---

# Centrality

In any [[Graph]], some nodes matter more than others — but *matter how?* **Centrality** is not one number but a family of them, each formalizing a different intuition about importance. The most popular person, the airport whose closure strands the most travelers, and the gossip best placed to spread news fastest are all "central," yet they can be entirely different nodes. Choosing a centrality measure means choosing what kind of importance you mean.

## Three classic answers

- **Degree centrality** — *how many connections?* Simply count a node's edges. The person with the most friends, the website with the most links. It is instant to compute and often a fine first cut, but it is purely local: it cannot tell a hub wired into the global network from one buried in an isolated clique.
- **Betweenness centrality** — *how often are you on the path between others?* For every pair of nodes, find the shortest path between them; a node's betweenness is the share of those paths that run *through* it. High-betweenness nodes are **bridges** — remove one and distant regions of the network fall apart. They control flow even with few connections of their own.
- **Closeness centrality** — *how near are you to everyone?* The inverse of a node's average shortest-path distance to all others. High closeness means short reach to the whole network — ideal for spreading something fast, computed with the help of [[Breadth-First Search]] from each node.

{{< note kind="key" title="Same node, different verdicts" >}}
A single edge connecting two large clusters has *low degree* (only two links) but *enormous betweenness* (every cross-cluster path uses it). Degree, betweenness, and closeness can crown three different kings of the same graph — which is exactly why a network analyst keeps all of them on hand.
{{< /note >}}

## A bridge node lights up

The graph below is two tight clusters joined by a single **bridge** node. By *degree* the bridge is unremarkable. By *betweenness* it is the most important node in the network: pulse-by-pulse, almost every path that crosses from one side to the other must pass through it. Watch the traffic — sized to flow through each node — and the bridge glows brightest.

{{< sketch height="380" caption="Two clusters joined by one bridge node. Node brightness and size track betweenness — the fraction of cross-network paths passing through. The unassuming bridge dominates. (Layout is fixed; the pulse animates flow.)" >}}
if (frame === 0 || !state.nodes) {
  state.nodes = [];
  const Lc = { x: W * 0.27, y: H * 0.5 }, Rc = { x: W * 0.73, y: H * 0.5 };
  // left cluster 0..5, right cluster 6..11, bridge node 12
  for (let i = 0; i < 6; i++) {
    const a = (i / 6) * Math.PI * 2;
    state.nodes.push({ x: Lc.x + Math.cos(a) * H * 0.26, y: Lc.y + Math.sin(a) * H * 0.26, side: 0 });
  }
  for (let i = 0; i < 6; i++) {
    const a = (i / 6) * Math.PI * 2;
    state.nodes.push({ x: Rc.x + Math.cos(a) * H * 0.26, y: Rc.y + Math.sin(a) * H * 0.26, side: 1 });
  }
  state.nodes.push({ x: W * 0.5, y: H * 0.5, side: 2 }); // bridge = 12
  state.edges = [];
  // dense-ish within each cluster
  for (let i = 0; i < 6; i++) for (let j = i + 1; j < 6; j++) if (Math.random() < 0.6) state.edges.push([i, j]);
  for (let i = 6; i < 12; i++) for (let j = i + 1; j < 12; j++) if (Math.random() < 0.6) state.edges.push([i, j]);
  // attach bridge to one node of each cluster, and clusters only meet through it
  state.edges.push([0, 12]); state.edges.push([6, 12]);
  // betweenness-ish weight: bridge huge, others modest (precomputed illustrative values)
  state.bet = new Array(13).fill(0.15);
  state.bet[12] = 1.0; state.bet[0] = 0.45; state.bet[6] = 0.45;
}
const cs = getComputedStyle(document.documentElement);
const accent = cs.getPropertyValue('--accent').trim() || '#5b9cff';
const accent2 = cs.getPropertyValue('--accent-2').trim() || '#ff9e64';
const border = cs.getPropertyValue('--border').trim() || 'rgba(255,255,255,0.18)';
ctx.fillStyle = '#0c1018';
ctx.fillRect(0, 0, W, H);
ctx.strokeStyle = border;
ctx.lineWidth = 1.2;
for (let e = 0; e < state.edges.length; e++) {
  const a = state.nodes[state.edges[e][0]], b = state.nodes[state.edges[e][1]];
  ctx.beginPath(); ctx.moveTo(a.x, a.y); ctx.lineTo(b.x, b.y); ctx.stroke();
}
// a pulse traveling across the bridge to suggest cross-cluster flow
const phase = (t * 0.6) % 1;
const A = state.nodes[0], M = state.nodes[12], B = state.nodes[6];
let px, py;
if (phase < 0.5) { const u = phase / 0.5; px = A.x + (M.x - A.x) * u; py = A.y + (M.y - A.y) * u; }
else { const u = (phase - 0.5) / 0.5; px = M.x + (B.x - M.x) * u; py = M.y + (B.y - M.y) * u; }
ctx.fillStyle = accent2;
ctx.beginPath(); ctx.arc(px, py, 4, 0, 7); ctx.fill();
// nodes sized/colored by betweenness
for (let i = 0; i < state.nodes.length; i++) {
  const b = state.bet[i];
  const rad = 5 + b * 16;
  const hue = 210 + b * 30;
  ctx.fillStyle = i === 12 ? accent2 : 'hsl(' + hue + ',75%,' + (45 + b * 20) + '%)';
  ctx.beginPath(); ctx.arc(state.nodes[i].x, state.nodes[i].y, rad, 0, 7); ctx.fill();
}
{{< /sketch >}}

## A spectrum of importance

[[PageRank]] is itself a centrality measure — a recursive, eigenvector-flavored cousin of degree centrality in which a connection's worth depends on the importance of *whoever it comes from*. That is the deep idea uniting the family: importance can be counted (degree), positional (betweenness, closeness), or self-referential (PageRank). Which one to trust depends entirely on what flows through your network and what you are trying to protect, spread, or rank.

## See also

- [[PageRank]]
- [[Graph]]
- [[Breadth-First Search]]
