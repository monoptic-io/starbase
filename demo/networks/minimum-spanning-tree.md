---
title: Minimum Spanning Tree
aliases: [spanning tree, MST, Kruskal, Prim]
tags: [networks, algorithms, weighted]
summary: The cheapest set of edges that connects every node of a weighted graph without forming a loop.
weight: 50
---

# Minimum Spanning Tree

Suppose you must wire up every house in a town — lay cable, build roads, run pipes — and each possible link has a cost. You want everything connected, and you want to spend as little as possible. The answer is a **minimum spanning tree** (MST): a subset of edges of a connected, weighted [[Graph]] that touches every node, contains no cycle, and has the smallest possible total weight.

It is a *tree* because connecting $n$ nodes without any redundant loop takes exactly $n-1$ edges — any more would create a cycle (waste), any fewer would leave something stranded. Among all the spanning trees a graph admits, the MST is the cheapest.

## Two greedy algorithms

Remarkably, you can build the optimal tree by always grabbing the cheapest edge that does not ruin the structure. Two classic algorithms formalize this:

- **Kruskal's algorithm** sorts *all* edges by weight and adds them cheapest-first, skipping any edge that would close a cycle. It grows a forest of fragments that gradually merge into one tree. (Detecting "would this make a cycle?" efficiently is what the union–find data structure is for.)
- **Prim's algorithm** grows a *single* tree from a starting node, repeatedly adding the cheapest edge that links the tree to a node not yet in it. This is strikingly close to [[Dijkstra's Algorithm]] — same priority-queue frontier — but Prim ranks frontier edges by their *own* weight, while Dijkstra ranks nodes by *total distance from the source*.

{{< note kind="key" title="The cut property" >}}
Both algorithms work because of one fact. Split the nodes into any two groups; the **cheapest edge crossing the divide must belong to some MST**. Every greedy step above is just an application of this rule — which is why grabbing the locally cheapest safe edge never sabotages the global optimum.
{{< /note >}}

## MST is not shortest paths

It is tempting to confuse the two weighted classics, but they optimize different things. [[Dijkstra's Algorithm]] minimizes the distance *from one source to each node individually*. An MST minimizes the *total* edge weight needed to hold the whole network together — and the path between two nodes *within* the MST is generally **not** their shortest path in the original graph. One answers "how do I get there fastest?"; the other answers "how do I connect everyone cheapest?"

## Watch Prim grow the tree

Below, scattered points are cities and the cost of a link is the straight-line distance between them. Prim's algorithm starts from one city and, frame by frame, reaches out along the *cheapest* edge to a city not yet connected (the bright candidate). The tree that results is the least-total-length way to wire them all together.

{{< sketch height="380" caption="Prim's algorithm building a minimum spanning tree over random cities. The orange edge is the cheapest link from the tree to an unconnected city, about to be locked in. Re-runs on a fresh scatter when complete." >}}
if (frame === 0 || !state.pts) {
  const N = 26;
  state.pts = [];
  for (let i = 0; i < N; i++) state.pts.push({ x: 24 + Math.random() * (W - 48), y: 24 + Math.random() * (H - 48) });
  state.inTree = new Array(N).fill(false);
  state.inTree[0] = true;
  state.tree = [];          // committed edges [a,b]
  state.count = 1;
  state.cand = null;
  state.done = 0;
  state.tick = 0;
}
const pts = state.pts, N = pts.length;
state.tick++;
// add one edge every ~12 frames so the growth is watchable
if (state.count < N && state.tick % 12 === 0) {
  let best = null, bd = Infinity;
  for (let a = 0; a < N; a++) {
    if (!state.inTree[a]) continue;
    for (let b = 0; b < N; b++) {
      if (state.inTree[b]) continue;
      const dx = pts[a].x - pts[b].x, dy = pts[a].y - pts[b].y, d = dx * dx + dy * dy;
      if (d < bd) { bd = d; best = [a, b]; }
    }
  }
  if (best) { state.tree.push(best); state.inTree[best[1]] = true; state.count++; state.cand = best; }
} else if (state.count >= N) {
  state.cand = null;
  state.done++;
  if (state.done > 90) state.pts = null;
}
const cs = getComputedStyle(document.documentElement);
const accent = cs.getPropertyValue('--accent').trim() || '#5b9cff';
const accent2 = cs.getPropertyValue('--accent-2').trim() || '#ff9e64';
const faint = cs.getPropertyValue('--text-faint').trim() || 'rgba(255,255,255,0.35)';
ctx.fillStyle = '#0c1018';
ctx.fillRect(0, 0, W, H);
// committed tree edges
ctx.strokeStyle = accent;
ctx.lineWidth = 1.8;
for (let e = 0; e < state.tree.length; e++) {
  const a = pts[state.tree[e][0]], b = pts[state.tree[e][1]];
  ctx.beginPath(); ctx.moveTo(a.x, a.y); ctx.lineTo(b.x, b.y); ctx.stroke();
}
// the candidate edge just added
if (state.cand) {
  const a = pts[state.cand[0]], b = pts[state.cand[1]];
  ctx.strokeStyle = accent2; ctx.lineWidth = 2.6;
  ctx.beginPath(); ctx.moveTo(a.x, a.y); ctx.lineTo(b.x, b.y); ctx.stroke();
}
// nodes
for (let i = 0; i < N; i++) {
  ctx.fillStyle = state.inTree[i] ? accent : faint;
  ctx.beginPath(); ctx.arc(pts[i].x, pts[i].y, state.inTree[i] ? 5 : 3.5, 0, 7); ctx.fill();
}
{{< /sketch >}}

{{< note kind="tip" title="Where MSTs show up" >}}
Designing utility, telephone, and transport networks; clustering data (cut the most expensive MST edges and the tree falls into natural groups); approximating the traveling-salesman tour; and even image segmentation. Whenever the goal is *cheap global connectivity*, an MST is lurking.
{{< /note >}}

## See also

- [[Graph]]
- [[Dijkstra's Algorithm]]
- [[Depth-First Search]]
