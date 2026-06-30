---
title: Depth-First Search
aliases: [DFS, depth first search]
tags: [networks, algorithms, traversal]
summary: Plunge as deep as possible along each branch, then backtrack — the traversal behind spanning trees, cycle detection, and topological order.
weight: 30
---

# Depth-First Search

**Depth-first search** (DFS) explores a [[Graph]] like someone wandering a maze with one hand on the wall: follow a corridor as far as it goes, and only when you hit a dead end do you *backtrack* to the last junction with an unexplored passage. Where [[Breadth-First Search]] fans out cautiously in layers, DFS commits — it races down a single branch to its end before considering any alternative.

This deep-first, backtrack-when-stuck discipline is captured by a **stack** (last-in, first-out), or equivalently by *recursion*, which uses the call stack for you.

## The algorithm

To run DFS from a node $u$: mark it visited, then recurse into each unvisited neighbor in turn. When a node has no unvisited neighbors left, the recursion unwinds — that unwinding *is* the backtracking. Like BFS, every vertex and edge is examined once, so DFS runs in $O(V + E)$.

{{< note kind="key" title="DFS, recursively" >}}
```
visit(u):
    mark u visited
    for each neighbor v of u:
        if v not visited:
            visit(v)          # dive deeper before trying the next neighbor
```
The order in which nodes are *finished* (recursion returns) is the secret behind topological sorting and strongly-connected-component algorithms.
{{< /note >}}

## What depth buys you

Diving deep is the wrong tool for shortest paths — the first time DFS reaches a node it may have taken a wildly roundabout route. But the structure DFS leaves behind is exactly what several classic problems need:

- **Spanning trees.** The edges DFS actually follows form a *DFS tree* — a connected, loop-free subgraph touching every reachable node. (BFS produces its own spanning tree, shallower and bushier.)
- **Cycle detection.** If DFS ever encounters an edge leading back to a node still on the current stack — a *back edge* — the graph contains a cycle. This is how you check whether a dependency graph is acyclic.
- **Ordering.** Recording nodes as their recursion *finishes* yields a **topological order** of a directed acyclic graph: build steps, course prerequisites, spreadsheet recalculation.

## A maze is a depth-first tree

Maze generation makes DFS tangible. Carve a grid by walking depth-first from a starting cell, knocking down a wall whenever you step to an unvisited neighbor; when you get stuck, backtrack. The single winding corridor with no loops that results is precisely a DFS spanning tree — every cell reachable, exactly one path between any two.

{{< sketch height="380" caption="DFS carving a maze. The bright cell is the current head plunging into unvisited territory; the dimmer trail is the stack it will backtrack along. The finished maze is a depth-first spanning tree. Re-runs when complete." >}}
if (frame === 0 || !state.grid) {
  state.cell = 22;
  state.cols = Math.floor(W / state.cell);
  state.rows = Math.floor(H / state.cell);
  const C = state.cols, R = state.rows;
  // each cell stores walls [top,right,bottom,left] and visited flag
  state.grid = [];
  for (let i = 0; i < C * R; i++) state.grid.push({ w: [true, true, true, true], v: false });
  state.stack = [0];
  state.grid[0].v = true;
  state.done = 0;
}
const C = state.cols, R = state.rows, cell = state.cell, g = state.grid;
// carve several steps per frame
for (let step = 0; step < 2 && state.stack.length > 0; step++) {
  const cur = state.stack[state.stack.length - 1];
  const cx = cur % C, cy = (cur - cx) / C;
  const nbs = [];
  // neighbor index, wall on current, wall on neighbor
  if (cy > 0 && !g[cur - C].v) nbs.push([cur - C, 0, 2]);
  if (cx < C - 1 && !g[cur + 1].v) nbs.push([cur + 1, 1, 3]);
  if (cy < R - 1 && !g[cur + C].v) nbs.push([cur + C, 2, 0]);
  if (cx > 0 && !g[cur - 1].v) nbs.push([cur - 1, 3, 1]);
  if (nbs.length === 0) { state.stack.pop(); continue; }
  const pick = nbs[Math.floor(Math.random() * nbs.length)];
  g[cur].w[pick[1]] = false;
  g[pick[0]].w[pick[2]] = false;
  g[pick[0]].v = true;
  state.stack.push(pick[0]);
}
if (state.stack.length === 0) {
  state.done++;
  if (state.done > 80) state.grid = null;
}
const cs = getComputedStyle(document.documentElement);
const accent = cs.getPropertyValue('--accent').trim() || '#5b9cff';
const accent2 = cs.getPropertyValue('--accent-2').trim() || '#ff9e64';
ctx.fillStyle = '#0c1018';
ctx.fillRect(0, 0, W, H);
// dim the stack trail
const onStack = new Set(state.stack);
for (let y = 0; y < R; y++) {
  for (let x = 0; x < C; x++) {
    const idx = y * C + x;
    if (onStack.has(idx)) {
      ctx.fillStyle = 'rgba(91,156,255,0.18)';
      ctx.fillRect(x * cell, y * cell, cell, cell);
    }
  }
}
// the head
if (state.stack.length) {
  const h = state.stack[state.stack.length - 1];
  ctx.fillStyle = accent2;
  ctx.fillRect((h % C) * cell + 3, Math.floor(h / C) * cell + 3, cell - 6, cell - 6);
}
// walls
ctx.strokeStyle = accent;
ctx.lineWidth = 1.4;
for (let y = 0; y < R; y++) {
  for (let x = 0; x < C; x++) {
    const w = g[y * C + x].w, px = x * cell, py = y * cell;
    ctx.beginPath();
    if (w[0]) { ctx.moveTo(px, py); ctx.lineTo(px + cell, py); }
    if (w[1]) { ctx.moveTo(px + cell, py); ctx.lineTo(px + cell, py + cell); }
    if (w[2]) { ctx.moveTo(px, py + cell); ctx.lineTo(px + cell, py + cell); }
    if (w[3]) { ctx.moveTo(px, py); ctx.lineTo(px, py + cell); }
    ctx.stroke();
  }
}
{{< /sketch >}}

{{< note kind="tip" title="Two traversals, one graph" >}}
Run BFS and DFS on the *same* graph and you discover the same set of nodes — they differ only in *order* and in the *tree* they leave behind. Choosing between them is choosing what structure you want to fall out of the walk.
{{< /note >}}

## See also

- [[Breadth-First Search]]
- [[Minimum Spanning Tree]]
- [[Graph]]
