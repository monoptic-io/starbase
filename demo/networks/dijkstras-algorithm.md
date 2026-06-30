---
title: Dijkstra's Algorithm
aliases: [shortest path, dijkstra, dijkstras algorithm]
tags: [networks, algorithms, weighted]
summary: Find the cheapest route in a weighted graph by greedily settling the nearest unfinished node, one at a time.
weight: 40
---

# Dijkstra's Algorithm

When the edges of a [[Graph]] carry **weights** — miles of road, minutes of travel, dollars of cost — the shortest path is no longer the one with the fewest hops. A route of three long edges may cost more than five short ones. **Dijkstra's algorithm** solves this *single-source shortest path* problem on graphs with non-negative weights, and it does so with one elegant idea: always finish the closest node you have not yet finished.

It is best understood as [[Breadth-First Search]] upgraded for weights. BFS expands a uniform wavefront because every edge costs $1$; Dijkstra expands a wavefront that bulges where edges are cheap and lags where they are expensive — and to keep the closest frontier node always at hand, it swaps BFS's plain queue for a **priority queue**.

## The greedy frontier

Each node holds a *tentative distance*: the best total cost found to it so far, starting at $\infty$ for all but the source. The algorithm repeatedly:

1. Picks the unfinished node $u$ with the smallest tentative distance and marks it **settled** — its distance is now final.
2. **Relaxes** each edge $u \to v$: if reaching $v$ through $u$ is cheaper than $v$'s current tentative distance, lower it.

{{< note kind="key" title="Why greed is safe here" >}}
When Dijkstra settles a node it claims that node's distance is final. This is only valid because **all weights are non-negative**: no path discovered later can sneak in at a *lower* total, since every additional edge only adds cost. Introduce a negative edge and the guarantee collapses — that case needs the Bellman–Ford algorithm instead.
{{< /note >}}

With a binary-heap priority queue the whole thing runs in $O((V + E)\log V)$.

## The wavefront, weighted

Below, each grid cell carries a random *terrain cost* (darker cells are slow, expensive ground). Dijkstra grows its settled region outward from the center, but instead of a clean diamond like BFS, the frontier *races* through cheap terrain and *crawls* through costly patches. Brightness shows the final shortest-path cost to each settled cell — the contour lines of cost bend around the expensive regions exactly as water finds the path of least resistance.

{{< sketch height="380" caption="Dijkstra settling outward over costly terrain. The frontier rushes through cheap ground and stalls in expensive ground; color shows total cost from the center. Re-runs from a new random landscape." >}}
if (frame === 0 || !state.dist) {
  state.cell = 14;
  state.cols = Math.floor(W / state.cell);
  state.rows = Math.floor(H / state.cell);
  const C = state.cols, R = state.rows, n = C * R;
  // smooth-ish random cost field in [1, 9]
  state.cost = new Array(n);
  for (let i = 0; i < n; i++) state.cost[i] = 1 + Math.floor(Math.random() * 9);
  state.dist = new Array(n).fill(Infinity);
  state.settled = new Array(n).fill(false);
  const s = Math.floor(R / 2) * C + Math.floor(C / 2);
  state.dist[s] = 0;
  // simple array-based priority list: [dist, node]
  state.pq = [ [0, s] ];
  state.maxd = 1;
  state.done = 0;
}
const C = state.cols, R = state.rows, cell = state.cell;
// settle a few nodes per frame
for (let step = 0; step < 12 && state.pq.length > 0; step++) {
  // extract-min (linear scan; the grid is small)
  let bi = 0;
  for (let i = 1; i < state.pq.length; i++) if (state.pq[i][0] < state.pq[bi][0]) bi = i;
  const item = state.pq.splice(bi, 1)[0];
  const u = item[1];
  if (state.settled[u]) { step--; continue; }
  state.settled[u] = true;
  state.maxd = Math.max(state.maxd, state.dist[u]);
  const ux = u % C, uy = (u - ux) / C;
  const nb = [ [ux + 1, uy], [ux - 1, uy], [ux, uy + 1], [ux, uy - 1] ];
  for (let k = 0; k < 4; k++) {
    const nx = nb[k][0], ny = nb[k][1];
    if (nx < 0 || ny < 0 || nx >= C || ny >= R) continue;
    const v = ny * C + nx;
    const nd = state.dist[u] + state.cost[v];
    if (nd < state.dist[v]) { state.dist[v] = nd; state.pq.push([nd, v]); }
  }
}
if (state.pq.length === 0) {
  state.done++;
  if (state.done > 70) state.dist = null;
}
ctx.fillStyle = '#0c1018';
ctx.fillRect(0, 0, W, H);
for (let y = 0; y < R; y++) {
  for (let x = 0; x < C; x++) {
    const idx = y * C + x;
    if (state.settled[idx]) {
      const t = state.dist[idx] / state.maxd;
      const hue = 210 + t * 120;
      ctx.fillStyle = 'hsl(' + hue + ',80%,' + (62 - t * 22) + '%)';
    } else {
      // unsettled: show terrain darkness
      const c = state.cost[idx];
      const g = 14 + (10 - c) * 3;
      ctx.fillStyle = 'rgb(' + g + ',' + (g + 4) + ',' + (g + 12) + ')';
    }
    ctx.fillRect(x * cell, y * cell, cell - 1, cell - 1);
  }
}
{{< /sketch >}}

## Cousins and uses

Dijkstra is the engine inside every routing system — road navigation, network packet routing, game pathfinding. Add a heuristic that estimates remaining distance and it becomes **A\***, which steers the frontier toward a goal instead of growing in all directions. Strip the weights back to a constant and it degenerates exactly into [[Breadth-First Search]]. And though it answers a different question — *cheapest route between two points* — it shares its greedy, frontier-growing spirit with the algorithms behind the [[Minimum Spanning Tree]], which instead seek the cheapest way to connect *everything*.

## See also

- [[Breadth-First Search]]
- [[Minimum Spanning Tree]]
- [[Graph]]
