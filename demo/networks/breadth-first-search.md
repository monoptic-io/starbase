---
title: Breadth-First Search
aliases: [BFS, breadth first search]
tags: [networks, algorithms, traversal]
summary: Explore a graph in concentric layers from a source — the simplest way to find shortest paths in an unweighted graph.
weight: 20
---

# Breadth-First Search

**Breadth-first search** (BFS) explores a graph the way ripples spread on a pond. Starting from a source node, it first visits every immediate neighbor, then every node two steps away, then three — fanning outward one *layer* at a time. Because it never moves to distance $d+1$ until it has exhausted everything at distance $d$, the moment BFS first reaches a node it has reached it by the fewest possible edges.

That single property makes BFS the canonical tool for **shortest paths in an unweighted [[Graph]]**: the layer in which a node appears *is* its distance from the source.

## The algorithm: a queue and a wavefront

BFS is built on a **queue** — a first-in, first-out line. You enqueue the source, then repeatedly take the front node, mark its unvisited neighbors as discovered, record their distance, and enqueue them. Because nodes leave the queue in the order they entered, the search expands strictly by distance.

{{< note kind="key" title="BFS in five lines" >}}
1. Mark the source visited, distance $0$, and put it in the queue.
2. Pop the front node $u$.
3. For each neighbor $v$ of $u$ not yet visited: mark it visited, set $\text{dist}(v) = \text{dist}(u) + 1$, enqueue $v$.
4. Repeat from step 2 until the queue is empty.

Every node and edge is touched once, so BFS runs in $O(V + E)$ time.
{{< /note >}}

## Watch the wavefront spread

On the grid below, BFS starts from a single cell and floods outward. Each cell is colored by the layer in which it is discovered — its graph distance from the source. The boundary between explored and unexplored is the **frontier**, and it is exactly the queue's contents at that instant. Notice how the colored bands form a diamond: on a 4-neighbor grid, graph distance is *Manhattan* distance, not straight-line distance.

{{< sketch height="380" caption="BFS flooding outward from a source cell. Color encodes distance (the BFS layer). The advancing edge is the frontier — the queue made visible. It re-runs from a new random source each pass." >}}
if (frame === 0 || !state.dist) {
  state.cell = 16;
  state.cols = Math.floor(W / state.cell);
  state.rows = Math.floor(H / state.cell);
  const C = state.cols, R = state.rows;
  state.dist = new Array(C * R).fill(-1);
  state.blocked = new Array(C * R).fill(false);
  // scatter a few walls so the wavefront has to bend around them
  for (let k = 0; k < C * R * 0.16; k++) {
    state.blocked[Math.floor(Math.random() * C * R)] = true;
  }
  const sx = Math.floor(C / 2), sy = Math.floor(R / 2);
  const s = sy * C + sx;
  state.blocked[s] = false;
  state.dist[s] = 0;
  state.queue = [s];
  state.head = 0;
  state.maxd = 0;
  state.done = 0;
}
const C = state.cols, R = state.rows, cell = state.cell;
// expand a handful of frontier nodes per frame for a smooth reveal
for (let step = 0; step < 6 && state.head < state.queue.length; step++) {
  const u = state.queue[state.head++];
  const ux = u % C, uy = (u - ux) / C, ud = state.dist[u];
  const nb = [ [ux + 1, uy], [ux - 1, uy], [ux, uy + 1], [ux, uy - 1] ];
  for (let n = 0; n < 4; n++) {
    const nx = nb[n][0], ny = nb[n][1];
    if (nx < 0 || ny < 0 || nx >= C || ny >= R) continue;
    const v = ny * C + nx;
    if (state.blocked[v] || state.dist[v] >= 0) continue;
    state.dist[v] = ud + 1;
    state.maxd = Math.max(state.maxd, ud + 1);
    state.queue.push(v);
  }
}
if (state.head >= state.queue.length) {
  state.done = (state.done || 0) + 1;
  if (state.done > 70) state.dist = null;
}
const cs = getComputedStyle(document.documentElement);
const border = cs.getPropertyValue('--border').trim() || 'rgba(255,255,255,0.12)';
ctx.fillStyle = '#0c1018';
ctx.fillRect(0, 0, W, H);
for (let y = 0; y < R; y++) {
  for (let x = 0; x < C; x++) {
    const idx = y * C + x;
    if (state.blocked[idx]) {
      ctx.fillStyle = border;
      ctx.fillRect(x * cell, y * cell, cell - 1, cell - 1);
    } else if (state.dist[idx] >= 0) {
      const t = state.maxd ? state.dist[idx] / state.maxd : 0;
      const hue = 200 + t * 130;
      ctx.fillStyle = 'hsl(' + hue + ',80%,58%)';
      ctx.fillRect(x * cell, y * cell, cell - 1, cell - 1);
    }
  }
}
{{< /sketch >}}

## When breadth wins

Because BFS settles each node by distance, it is the right choice whenever every edge counts the same: the fewest hops between two people in a social graph, the minimum moves in a puzzle, the nearest exit in a maze. Its sibling [[Depth-First Search]] instead dives as deep as it can before backing up — better for detecting cycles and ordering dependencies, but with no shortest-path guarantee. And the moment edges carry *unequal* weights, layer-by-layer breadth is no longer enough: you need the weighted generalization, [[Dijkstra's Algorithm]], which is BFS with a priority queue in place of a plain one.

## See also

- [[Depth-First Search]]
- [[Dijkstra's Algorithm]]
- [[Graph]]
