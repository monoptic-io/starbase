---
title: Small-World Network
aliases: [small world, six degrees, Watts-Strogatz]
tags: [networks, structure]
summary: Networks that are richly clustered yet only a few steps across — the structure behind "six degrees of separation."
weight: 80
---

# Small-World Network

Your friends mostly know each other — society is full of tight little cliques. And yet any two people on Earth are joined by a chain of only a handful of acquaintances. Those two facts seem to pull in opposite directions: dense local clustering *should* trap you in your neighborhood, far from everyone else. A **small-world network** reconciles them. It is a [[Graph]] that is simultaneously **highly clustered** (your neighbors are neighbors) and has a **short average path length** (everyone is just a few hops away). This is the network behind "six degrees of separation."

## Clustering and path length

Two numbers capture the tension:

- **Clustering coefficient** — the chance that two of a node's neighbors are themselves connected. High clustering means dense local triangles, the signature of communities.
- **Average path length** — the typical number of edges on the shortest route between two random nodes. Low path length means the whole network is shallow.

A plain ring lattice — each node wired only to its nearest neighbors — is *highly clustered* but has a *long* path length: to cross it you must shuffle step by step around the rim. A purely random graph is the opposite: short paths, but almost no clustering. The small-world regime lives in between, and the surprise is how *little* it takes to get there.

## The Watts–Strogatz rewiring

In 1998 Watts and Strogatz showed the recipe. Start from the clustered ring lattice. Now, with probability $p$, **rewire** each edge to a random target — replacing a local link with a long-range *shortcut*. The astonishing result: just a *few* shortcuts collapse the average path length almost to that of a random graph, while the clustering stays high. A handful of long jumps stitch the distant arcs of the ring together, and the world shrinks.

{{< note kind="key" title="A few shortcuts change everything" >}}
At $p = 0$ you have a clustered but vast ring; at $p = 1$ a short but clusterless random graph. The small-world magic happens at *tiny* $p$: path length plummets long before clustering does. A network can be almost entirely local and still be globally shallow — all it takes is a sprinkling of shortcuts.
{{< /note >}}

## Rewire the ring yourself

Below is a ring lattice. **Move the mouse left and right to set the rewiring probability $p$** (left edge = ordered ring, right edge = heavily rewired). As $p$ rises, local edges are replaced by chords leaping across the circle, and the measured average path length — shown at the bottom — drops sharply while the ring's clustered backbone survives.

{{< sketch height="420" caption="A Watts-Strogatz ring. Move the mouse left/right to set the rewiring probability p. A few random shortcuts (chords across the circle) collapse the average path length while clustering persists." >}}
if (frame === 0 || !state.base) {
  state.N = 30;
  state.K = 4; // neighbors each side total -> K/2 each side
  // base ring adjacency: store undirected edges as pairs (i, (i+s)%N)
  state.base = [];
  for (let i = 0; i < state.N; i++) {
    for (let s = 1; s <= state.K / 2; s++) state.base.push([i, (i + s) % state.N]);
  }
  state.lastBucket = -1;
  state.edges = state.base.map(e => [e[0], e[1]]);
  state.apl = 0;
}
const N = state.N;
// p from mouse x (clamped); default mild if mouse outside
let p = 0.0;
if (mouse.x > 0 && mouse.x < W) p = Math.max(0, Math.min(1, mouse.x / W));
// only rebuild when p changes enough (bucketize) to avoid per-frame churn
const bucket = Math.round(p * 20);
if (bucket !== state.lastBucket) {
  state.lastBucket = bucket;
  const pp = bucket / 20;
  // rewire from a fresh copy of the base ring
  const edges = [];
  const adj = Array.from({ length: N }, () => new Set());
  for (let e = 0; e < state.base.length; e++) {
    let a = state.base[e][0], b = state.base[e][1];
    if (Math.random() < pp) {
      // rewire endpoint b to a random non-self, non-duplicate target
      let tries = 0, nb = b;
      do { nb = Math.floor(Math.random() * N); tries++; } while ((nb === a || adj[a].has(nb)) && tries < 20);
      b = nb;
    }
    if (a !== b && !adj[a].has(b)) { adj[a].add(b); adj[b].add(a); edges.push([a, b]); }
  }
  state.edges = edges;
  state.adj = adj.map(s => [...s]);
  // average shortest path length via BFS from every node
  let total = 0, pairs = 0;
  for (let s = 0; s < N; s++) {
    const dist = new Array(N).fill(-1); dist[s] = 0;
    const q = [s]; let h = 0;
    while (h < q.length) {
      const u = q[h++];
      const nbrs = state.adj[u];
      for (let k = 0; k < nbrs.length; k++) {
        const v = nbrs[k];
        if (dist[v] < 0) { dist[v] = dist[u] + 1; q.push(v); }
      }
    }
    for (let v = 0; v < N; v++) if (v !== s && dist[v] > 0) { total += dist[v]; pairs++; }
  }
  state.apl = pairs ? total / pairs : 0;
  state.p = pp;
}
const cs = getComputedStyle(document.documentElement);
const accent = cs.getPropertyValue('--accent').trim() || '#5b9cff';
const accent2 = cs.getPropertyValue('--accent-2').trim() || '#ff9e64';
const border = cs.getPropertyValue('--border').trim() || 'rgba(255,255,255,0.18)';
const text = cs.getPropertyValue('--text').trim() || '#e8ecf4';
ctx.fillStyle = '#0c1018';
ctx.fillRect(0, 0, W, H);
const cx = W / 2, cy = H / 2 - 6, rad = Math.min(W, H) * 0.36;
const pos = [];
for (let i = 0; i < N; i++) {
  const a = (i / N) * Math.PI * 2 - Math.PI / 2;
  pos.push({ x: cx + Math.cos(a) * rad, y: cy + Math.sin(a) * rad });
}
for (let e = 0; e < state.edges.length; e++) {
  const a = state.edges[e][0], b = state.edges[e][1];
  // local ring edges dim; long shortcuts bright
  const span = Math.min((b - a + N) % N, (a - b + N) % N);
  const isShort = span > state.K / 2;
  ctx.strokeStyle = isShort ? accent2 : border;
  ctx.lineWidth = isShort ? 1.6 : 1.0;
  ctx.globalAlpha = isShort ? 0.9 : 0.6;
  ctx.beginPath(); ctx.moveTo(pos[a].x, pos[a].y); ctx.lineTo(pos[b].x, pos[b].y); ctx.stroke();
}
ctx.globalAlpha = 1;
ctx.fillStyle = accent;
for (let i = 0; i < N; i++) { ctx.beginPath(); ctx.arc(pos[i].x, pos[i].y, 4.5, 0, 7); ctx.fill(); }
ctx.fillStyle = text;
ctx.font = '13px sans-serif';
ctx.fillText('p = ' + (state.p || 0).toFixed(2) + '   avg path length = ' + (state.apl || 0).toFixed(2), 12, H - 12);
{{< /sketch >}}

Most real networks — neurons, power grids, the web, and the [[Complex Systems]] of society itself — turn out to be small worlds. The structure is not a curiosity; it is what lets information, disease, and influence travel a vast network in just a few steps. Its frequent companion, the [[Scale-Free Network]], explains where the long-range shortcuts and their hubs come from in the first place.

## See also

- [[Scale-Free Network]]
- [[Complex Systems]]
- [[Graph]]
