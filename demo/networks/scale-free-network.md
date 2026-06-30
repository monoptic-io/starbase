---
title: Scale-Free Network
aliases: [scale free, power law, preferential attachment, Barabasi-Albert]
tags: [networks, structure]
summary: Networks dominated by a few enormous hubs, grown by "the rich get richer" and obeying a power-law degree distribution.
weight: 90
---

# Scale-Free Network

Some networks are democratic — most nodes have roughly the same number of links, and a node with ten times the average is freakishly rare. Many real networks are nothing like that. The web has a handful of pages with millions of inbound links and billions with a few. Air travel has a few mega-hub airports and thousands of small fields. These are **scale-free networks**: their connectivity is wildly uneven, ruled by a small number of enormous **hubs**.

The name comes from their **degree distribution**. The fraction of nodes with degree $k$ falls off as a *power law*,

$$ P(k) \sim k^{-\gamma}, \qquad \gamma \approx 2\text{–}3. $$

A power law has no characteristic "scale" — no typical degree around which everything clusters — which is exactly what lets a few nodes tower over the rest. Contrast a bell-shaped distribution, where a node a hundred times the average is essentially impossible.

## Hubs are grown, not designed

Where do the hubs come from? The **Barabási–Albert** model gives a startlingly simple answer with two ingredients:

- **Growth** — the network is not fixed; new nodes arrive over time.
- **Preferential attachment** — a newcomer is more likely to link to a node that *already* has many links. The well-connected get more connections. The rich get richer.

That feedback loop is all it takes. An early node that happens to gain a few extra links becomes a slightly more attractive target, draws still more links, and snowballs into a hub. The power law emerges on its own, with no designer and no central plan — pure [[Complex Systems]] behavior.

{{< note kind="key" title="Robust yet fragile" >}}
Scale-free structure has a famous double edge. Knock out *random* nodes and the network barely notices — almost all nodes are small, so you almost always hit a leaf. But strike at the *hubs* deliberately and it shatters fast. The same hub-dominated shape that makes the internet resilient to random failure makes it vulnerable to targeted attack — and makes super-spreaders pivotal in an epidemic.
{{< /note >}}

## Grow a network and watch hubs emerge

Below, the network grows one node at a time. Each arrival attaches to existing nodes chosen *in proportion to how many links they already have*. Watch a few lucky early nodes pull ahead and bloom into hubs while the rest stay small. Node size and color track degree — the [[Centrality]] of each node made visible.

{{< sketch height="420" caption="Barabasi-Albert growth: each new node attaches preferentially to already-popular nodes. A few hubs emerge and dominate. Node size and color encode degree. Re-runs after it fills." >}}
if (frame === 0 || !state.nodes) {
  state.nodes = [];
  state.edges = [];
  state.targets = [];   // attachment list: each node appears once per link (for preferential sampling)
  state.maxN = 70;
  state.m = 2;          // links each new node makes
  // seed with a small connected core
  for (let i = 0; i < 3; i++) state.nodes.push({ x: W / 2 + (Math.random() - 0.5) * 60, y: H / 2 + (Math.random() - 0.5) * 60, vx: 0, vy: 0, deg: 0 });
  const seed = [ [0, 1], [1, 2], [2, 0] ];
  for (let e = 0; e < seed.length; e++) {
    const a = seed[e][0], b = seed[e][1];
    state.edges.push([a, b]); state.nodes[a].deg++; state.nodes[b].deg++;
    state.targets.push(a, b);
  }
  state.tick = 0;
  state.done = 0;
}
const nodes = state.nodes, edges = state.edges;
state.tick++;
// add a node every ~10 frames
if (nodes.length < state.maxN && state.tick % 10 === 0) {
  const nx = W / 2 + (Math.random() - 0.5) * W * 0.5;
  const ny = H / 2 + (Math.random() - 0.5) * H * 0.5;
  const idx = nodes.length;
  nodes.push({ x: nx, y: ny, vx: 0, vy: 0, deg: 0 });
  const chosen = new Set();
  let guard = 0;
  while (chosen.size < state.m && guard < 50) {
    guard++;
    const cand = state.targets[Math.floor(Math.random() * state.targets.length)];
    if (cand !== idx) chosen.add(cand);
  }
  chosen.forEach(c => {
    edges.push([idx, c]);
    nodes[idx].deg++; nodes[c].deg++;
    state.targets.push(idx, c);
  });
} else if (nodes.length >= state.maxN) {
  state.done++;
  if (state.done > 120) state.nodes = null;
}
// light force layout so structure stays legible
const N = nodes.length;
for (let i = 0; i < N; i++) {
  for (let j = i + 1; j < N; j++) {
    let dx = nodes[i].x - nodes[j].x, dy = nodes[i].y - nodes[j].y;
    let d2 = dx * dx + dy * dy + 0.01, d = Math.sqrt(d2);
    const f = 900 / d2;
    nodes[i].vx += dx / d * f; nodes[i].vy += dy / d * f;
    nodes[j].vx -= dx / d * f; nodes[j].vy -= dy / d * f;
  }
}
const L = 46;
for (let e = 0; e < edges.length; e++) {
  const a = edges[e][0], b = edges[e][1];
  let dx = nodes[b].x - nodes[a].x, dy = nodes[b].y - nodes[a].y;
  let d = Math.sqrt(dx * dx + dy * dy) + 0.01;
  const f = (d - L) * 0.02;
  nodes[a].vx += dx / d * f; nodes[a].vy += dy / d * f;
  nodes[b].vx -= dx / d * f; nodes[b].vy -= dy / d * f;
}
for (let i = 0; i < N; i++) {
  nodes[i].vx += (W / 2 - nodes[i].x) * 0.006;
  nodes[i].vy += (H / 2 - nodes[i].y) * 0.006;
  nodes[i].vx *= 0.8; nodes[i].vy *= 0.8;
  const vmax = 14;
  nodes[i].vx = Math.max(-vmax, Math.min(vmax, nodes[i].vx));
  nodes[i].vy = Math.max(-vmax, Math.min(vmax, nodes[i].vy));
  nodes[i].x += nodes[i].vx; nodes[i].y += nodes[i].vy;
  nodes[i].x = Math.max(10, Math.min(W - 10, nodes[i].x));
  nodes[i].y = Math.max(10, Math.min(H - 10, nodes[i].y));
}
const cs = getComputedStyle(document.documentElement);
const accent = cs.getPropertyValue('--accent').trim() || '#5b9cff';
const accent2 = cs.getPropertyValue('--accent-2').trim() || '#ff9e64';
const border = cs.getPropertyValue('--border').trim() || 'rgba(255,255,255,0.16)';
ctx.fillStyle = '#0c1018';
ctx.fillRect(0, 0, W, H);
ctx.strokeStyle = border;
ctx.lineWidth = 1.0;
for (let e = 0; e < edges.length; e++) {
  const a = nodes[edges[e][0]], b = nodes[edges[e][1]];
  ctx.beginPath(); ctx.moveTo(a.x, a.y); ctx.lineTo(b.x, b.y); ctx.stroke();
}
let dmax = 1; for (let i = 0; i < N; i++) dmax = Math.max(dmax, nodes[i].deg);
for (let i = 0; i < N; i++) {
  const t = nodes[i].deg / dmax;
  const rad = 3 + t * 16;
  ctx.fillStyle = t > 0.5 ? accent2 : accent;
  ctx.globalAlpha = 0.9;
  ctx.beginPath(); ctx.arc(nodes[i].x, nodes[i].y, rad, 0, 7); ctx.fill();
  ctx.globalAlpha = 1;
}
{{< /sketch >}}

## The telltale straight line

Plot a power-law degree distribution on log–log axes and it becomes a straight line with slope $-\gamma$ — the fingerprint analysts look for. Most nodes sit at low degree (the tall left bars); a long, thin **tail** of hubs stretches far to the right, nodes a bell curve would forbid.

{{< chart type="bar" data="1:48, 2:22, 3:12, 4:7, 5:4, 6:3, 7:2, 8:1, 12:1, 20:1" title="A scale-free degree distribution: many small nodes, a long tail of hubs (degree : count)" xlabel="degree k" ylabel="number of nodes" >}}

Scale-free and [[Small-World Network]] structure usually travel together: preferential-attachment hubs double as the long-range shortcuts that make a network shallow, so real systems are typically *both* hub-dominated and small-world at once.

## See also

- [[Small-World Network]]
- [[Centrality]]
- [[Complex Systems]]
