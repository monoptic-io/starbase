---
title: Markov Chain
aliases: [markov chains, markov process, transition matrix]
tags: [probability]
summary: A random process that hops between states where the next state depends only on the current one, forgetting all earlier history.
weight: 50
---

# Markov Chain

A **Markov chain** is a random process with no memory. It lives in a set of **states** and hops between them at each step, but the rule for where it goes next depends *only on where it is now* — never on the path that brought it there. This **memorylessness** (the *Markov property*) is a drastic simplification, and yet it is rich enough to model weather, board games, queues, shuffles, language, and the random web-surfer behind [[PageRank]].

All the dynamics live in one object: the **transition matrix** $P$, whose entry $P_{ij}$ is the probability of jumping from state $i$ to state $j$. Each row lists where you might go from one state and sums to $1$ — you must land somewhere.

## From a state to a distribution

Instead of tracking a single wandering token, track a *probability distribution* over the states — a vector $\pi$ saying how much "mass" sits on each. One step of the chain advances it by a [[Matrix]] multiply:

{{< eq number="1" >}}
\pi_{t+1} = \pi_t \, P .
{{< /eq >}}

Apply $P$ again and again and, for most chains, the distribution stops changing. It reaches a **stationary distribution** $\pi^\star$ that satisfies

{{< eq number="2" >}}
\pi^\star = \pi^\star P .
{{< /eq >}}

{{< note kind="key" title="The stationary distribution is an eigenvector" >}}
Equation (2) says $\pi^\star$ is a **left eigenvector of $P$ with eigenvalue $1$** — exactly the [[Eigenvalues and Eigenvectors]] idea. Repeatedly multiplying by $P$ is *power iteration*, and it converges to that dominant eigenvector. This is why [[PageRank]] — the stationary distribution of a random surfer's Markov chain — *is* an eigenvector computation in disguise.
{{< /note >}}

## Watch the mass converge

Below are a few states joined by weighted arrows (the transition probabilities). Probability mass — drawn as the **area of each node** — starts piled on one state and then flows along the arrows every step, splitting according to the weights. No matter where it starts, it sloshes around and then **freezes into the stationary distribution**: the unique mix that flows into itself unchanged. The bars on the right track the same distribution converging.

{{< sketch height="400" caption="A 5-state Markov chain converging to its stationary distribution. Node area is the current probability mass; it flows along weighted arrows each step (power iteration) and settles onto the eigenvector with eigenvalue 1. Bars at right show the same distribution. Auto-resets from a fresh random pile." >}}
if (frame === 0 || !state.nodes) {
  const N = 5;
  state.N = N;
  state.nodes = [];
  for (let i = 0; i < N; i++) {
    const ang = -Math.PI / 2 + (i / N) * Math.PI * 2;
    state.nodes.push({
      x: W * 0.36 + Math.cos(ang) * Math.min(W, H) * 0.30,
      y: H * 0.5 + Math.sin(ang) * Math.min(W, H) * 0.30
    });
  }
  // random row-stochastic transition matrix
  state.P = [];
  for (let i = 0; i < N; i++) {
    const row = [];
    let s = 0;
    for (let j = 0; j < N; j++) { const v = Math.random() * (i === j ? 0.6 : 1) + 0.05; row.push(v); s += v; }
    for (let j = 0; j < N; j++) row[j] /= s;
    state.P.push(row);
  }
  // start mass all on one node
  state.pi = new Array(N).fill(0);
  state.pi[Math.floor(Math.random() * N)] = 1;
  state.tick = 0;
  state.iters = 0;
}
const N = state.N;
const cs = getComputedStyle(document.documentElement);
const accent = cs.getPropertyValue('--accent').trim() || '#5b9cff';
const accent2 = cs.getPropertyValue('--accent-2').trim() || '#ff9e64';
const good = cs.getPropertyValue('--good').trim() || '#9ece6a';
const faint = cs.getPropertyValue('--text-faint').trim() || 'rgba(255,255,255,0.4)';
const border = cs.getPropertyValue('--border').trim() || 'rgba(255,255,255,0.18)';
ctx.clearRect(0, 0, W, H);
ctx.fillStyle = '#0c1018';
ctx.fillRect(0, 0, W, H);
state.tick++;
// one step every 16 frames
if (state.tick % 16 === 0 && state.iters < 60) {
  const np = new Array(N).fill(0);
  for (let i = 0; i < N; i++) for (let j = 0; j < N; j++) np[j] += state.pi[i] * state.P[i][j];
  state.pi = np;
  state.iters++;
}
// edges with arrowheads, width ~ transition prob
for (let i = 0; i < N; i++) {
  for (let j = 0; j < N; j++) {
    if (i === j || state.P[i][j] < 0.06) continue;
    const a = state.nodes[i], b = state.nodes[j];
    const dx = b.x - a.x, dy = b.y - a.y, L = Math.sqrt(dx * dx + dy * dy);
    const ux = dx / L, uy = dy / L;
    // offset perpendicular so i->j and j->i don't overlap
    const ox = -uy * 6, oy = ux * 6;
    const sx = a.x + ux * 20 + ox, sy = a.y + uy * 20 + oy;
    const ex = b.x - ux * 20 + ox, ey = b.y - uy * 20 + oy;
    ctx.strokeStyle = border;
    ctx.lineWidth = 0.5 + state.P[i][j] * 6;
    ctx.beginPath(); ctx.moveTo(sx, sy); ctx.lineTo(ex, ey); ctx.stroke();
    ctx.fillStyle = border;
    ctx.beginPath();
    ctx.moveTo(ex, ey);
    ctx.lineTo(ex - ux * 8 - uy * 5, ey - uy * 8 + ux * 5);
    ctx.lineTo(ex - ux * 8 + uy * 5, ey - uy * 8 - ux * 5);
    ctx.fill();
  }
}
// nodes sized by mass
for (let i = 0; i < N; i++) {
  const rad = 8 + Math.sqrt(state.pi[i]) * 34;
  ctx.fillStyle = accent;
  ctx.globalAlpha = 0.85;
  ctx.beginPath(); ctx.arc(state.nodes[i].x, state.nodes[i].y, rad, 0, 7); ctx.fill();
  ctx.globalAlpha = 1;
  ctx.fillStyle = good;
  ctx.font = '11px sans-serif';
  ctx.fillText((state.pi[i] * 100).toFixed(0) + '%', state.nodes[i].x - 10, state.nodes[i].y + 4);
}
// stationary bars at right
const bx = W * 0.78, bw = (W * 0.20) / N, bTop = H * 0.18, bH = H * 0.62;
ctx.strokeStyle = border;
ctx.beginPath(); ctx.moveTo(bx - 6, bTop + bH); ctx.lineTo(bx + N * bw + 4, bTop + bH); ctx.stroke();
for (let i = 0; i < N; i++) {
  const h = state.pi[i] * bH;
  ctx.fillStyle = accent2;
  ctx.fillRect(bx + i * bw + 2, bTop + bH - h, bw - 4, h);
}
ctx.fillStyle = faint;
ctx.font = '12px sans-serif';
ctx.fillText('step ' + state.iters + (state.iters >= 60 ? '  (stationary)' : ''), 10, H - 12);
ctx.fillText('distribution', bx, bTop - 6);
if (state.iters >= 60) {
  state.hold = (state.hold || 0) + 1;
  if (state.hold > 120) { state.nodes = null; state.hold = 0; }
}
{{< /sketch >}}

## When does a chain settle?

Not every chain has a single stationary destiny. A chain that can eventually get from any state to any other (**irreducible**) and does not lock into a rigid cycle (**aperiodic**) is guaranteed to converge to one unique $\pi^\star$ from any start. Chains with isolated traps or perfect periodicity can fail this — they may oscillate forever or depend on where they began. When the conditions hold, the long-run behavior is beautifully indifferent to initial conditions, the opposite of the [[Chaos]] story.

A Markov chain is also just a [[Dynamical System]] whose state is a probability vector and whose rule is the matrix $P$ — randomness and linear algebra, fused.

## See also

- [[PageRank]]
- [[Eigenvalues and Eigenvectors]]
- [[Random Walk]]
