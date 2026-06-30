---
title: PageRank
aliases: [page rank, random surfer]
tags: [networks, algorithms, linear-algebra]
summary: A node's importance is the long-run chance a random surfer following links lands on it — the dominant eigenvector of the link matrix.
weight: 60
---

# PageRank

How do you rank the importance of a page in a web of millions, using nothing but the links between them? **PageRank**, the algorithm that launched Google, gives a beautifully circular answer: *a page is important if important pages link to it.* Importance is defined in terms of itself — and the way out of that circle is one of the loveliest applications of [[Linear Algebra]] to a real-world [[Graph]].

## The random surfer

Picture someone browsing at random. They sit on a page, pick one of its outgoing links uniformly at random, click it, and repeat — forever. **A page's PageRank is the fraction of time this random surfer spends on it.** Pages with many incoming links get visited often; a link *from* a frequently-visited page is worth more than a link from an obscure one, because the surfer arrives there more often to begin with. Importance flows along edges and pools where the structure concentrates it.

One fix keeps the wandering well-behaved. A surfer can get stuck on a page with no outgoing links, or trapped in a small cluster. So with a small probability $d \approx 0.15$ — the **damping** — the surfer *teleports* to a random page instead of following a link. This keeps every page reachable and the long-run fractions unique.

{{< note kind="key" title="Importance is a fixed point" >}}
Let $r_i$ be the rank of page $i$, and let page $j$ have $L_j$ outgoing links. Each step, page $j$ hands an equal share of its rank to every page it points to:

$$ r_i = \frac{d}{N} + (1-d) \sum_{j \to i} \frac{r_j}{L_j}. $$

The ranks that satisfy this equation are *self-consistent* — feed them through one more step and they come back unchanged. That stationary distribution **is** PageRank.
{{< /note >}}

## It is an eigenvector

Stack the ranks into a vector $\mathbf{r}$ and the link structure into a matrix $M$ (with the teleport term folded in). The fixed-point equation becomes simply

$$ \mathbf{r} = M\,\mathbf{r}. $$

That is the definition of an **eigenvector** with eigenvalue $1$. PageRank is nothing but the **dominant eigenvector** of the web's link matrix — exactly the [[Eigenvalues and Eigenvectors]] story from the linear-algebra section, made enormous. And you compute it the simplest way imaginable: start with any guess, multiply by $M$ over and over, and watch it converge. This **power iteration** is the random surfer's wandering written as repeated [[Matrix Multiplication]].

{{< note kind="note" title="This site runs on it" >}}
The *related topics* this knowledge base suggests are ranked with a PageRank-style computation over the very graph of wiki-links you are clicking. The page you are reading is a node; every `[[link]]` is an edge; importance flows along them just as described here.
{{< /note >}}

## Watch the rank settle

Each node below starts with an equal share of rank. Every step, a node pushes its rank evenly out along its edges (with a dash of teleport mixed back in). A node's *size* tracks its current rank. Watch the mass slosh around and then settle — the well-connected nodes swell, the peripheral ones shrink, and after a dozen iterations the picture stops changing. That frozen distribution is the dominant eigenvector.

{{< sketch height="400" caption="Power iteration computing PageRank. Node area is proportional to current rank; rank flows along directed edges each step and converges to the stationary distribution. Re-runs on a new graph after it settles." >}}
if (frame === 0 || !state.nodes) {
  const N = 12;
  state.nodes = [];
  for (let i = 0; i < N; i++) {
    const ang = (i / N) * Math.PI * 2;
    state.nodes.push({
      x: W / 2 + Math.cos(ang) * Math.min(W, H) * 0.34,
      y: H / 2 + Math.sin(ang) * Math.min(W, H) * 0.34
    });
  }
  // random directed links: each node points to 1-3 others
  state.out = [];
  for (let i = 0; i < N; i++) {
    const k = 1 + Math.floor(Math.random() * 3), set = new Set();
    while (set.size < k) { const t = Math.floor(Math.random() * N); if (t !== i) set.add(t); }
    state.out.push([...set]);
  }
  state.r = new Array(N).fill(1 / N);
  state.tick = 0;
  state.iters = 0;
  state.done = 0;
}
const N = state.nodes.length, d = 0.15;
state.tick++;
// one power-iteration step every 14 frames
if (state.tick % 14 === 0 && state.iters < 40) {
  const nr = new Array(N).fill(d / N);
  for (let j = 0; j < N; j++) {
    const outs = state.out[j];
    if (outs.length === 0) { for (let i = 0; i < N; i++) nr[i] += (1 - d) * state.r[j] / N; continue; }
    for (let t = 0; t < outs.length; t++) nr[outs[t]] += (1 - d) * state.r[j] / outs.length;
  }
  state.r = nr;
  state.iters++;
  if (state.iters >= 40) { state.done = 1; }
}
if (state.done) { state.doneCt = (state.doneCt || 0) + 1; if (state.doneCt > 100) { state.nodes = null; state.doneCt = 0; state.done = 0; } }
const cs = getComputedStyle(document.documentElement);
const accent = cs.getPropertyValue('--accent').trim() || '#5b9cff';
const accent2 = cs.getPropertyValue('--accent-2').trim() || '#ff9e64';
const border = cs.getPropertyValue('--border').trim() || 'rgba(255,255,255,0.18)';
ctx.fillStyle = '#0c1018';
ctx.fillRect(0, 0, W, H);
// directed edges with arrowheads
ctx.strokeStyle = border;
ctx.fillStyle = border;
ctx.lineWidth = 1.1;
for (let j = 0; j < N; j++) {
  for (let t = 0; t < state.out[j].length; t++) {
    const a = state.nodes[j], b = state.nodes[state.out[j][t]];
    const dx = b.x - a.x, dy = b.y - a.y, L = Math.sqrt(dx * dx + dy * dy);
    const ux = dx / L, uy = dy / L;
    const ex = b.x - ux * 16, ey = b.y - uy * 16;
    ctx.beginPath(); ctx.moveTo(a.x + ux * 14, a.y + uy * 14); ctx.lineTo(ex, ey); ctx.stroke();
    // arrowhead
    ctx.beginPath();
    ctx.moveTo(ex, ey);
    ctx.lineTo(ex - ux * 6 - uy * 4, ey - uy * 6 + ux * 4);
    ctx.lineTo(ex - ux * 6 + uy * 4, ey - uy * 6 - ux * 4);
    ctx.fill();
  }
}
// nodes sized by rank
let rmax = 0; for (let i = 0; i < N; i++) rmax = Math.max(rmax, state.r[i]);
for (let i = 0; i < N; i++) {
  const rad = 6 + (state.r[i] / rmax) * 22;
  ctx.fillStyle = i === 0 ? accent : accent;
  ctx.globalAlpha = 0.85;
  ctx.beginPath(); ctx.arc(state.nodes[i].x, state.nodes[i].y, rad, 0, 7); ctx.fill();
  ctx.globalAlpha = 1;
}
ctx.fillStyle = accent2;
ctx.font = '12px sans-serif';
ctx.fillText('iteration ' + state.iters, 10, H - 10);
{{< /sketch >}}

## Beyond the web

PageRank long ago escaped search engines. It scores influence in social networks, importance of species in food webs, key papers in citation graphs, and central roads in transport maps — anywhere "important things point to important things" makes sense. It is one face of a larger family of [[Centrality]] measures, distinguished by its recursive, eigenvector-based definition of importance.

## See also

- [[Eigenvalues and Eigenvectors]]
- [[Centrality]]
- [[Graph]]
