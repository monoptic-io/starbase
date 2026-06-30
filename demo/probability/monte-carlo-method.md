---
title: Monte Carlo Method
aliases: [monte carlo]
tags: [probability]
summary: Estimate a hard quantity by generating many random samples and averaging the results.
weight: 70
---

# Monte Carlo Method

The **Monte Carlo method** turns a hard calculation into a game of chance. When a quantity is too tangled to compute directly — a high-dimensional integral, the value of a board position, the behavior of a messy physical system — you instead *sample it at random* many times and average. The [[Law of Large Numbers]] guarantees the average converges to the true answer, and the [[Central Limit Theorem]] tells you how fast. Named after the Monaco casino, it trades exact analysis for honest dice.

The method was born on the Manhattan Project, where physicists needed to track neutrons diffusing through material — a problem with no clean formula but an easy *random* description. Today it estimates everything from financial risk to the volume of complicated shapes to the moves a game-playing AI should consider.

## The core idea

To find some quantity, design a random experiment whose **average outcome equals the quantity you want**, then run it over and over:

{{< eq number="1" >}}
\text{answer} \;\approx\; \frac{1}{N}\sum_{i=1}^{N} f(\text{sample}_i).
{{< /eq >}}

The error shrinks like $1/\sqrt{N}$ — slow but utterly general. It does not care how many dimensions the problem has, which is exactly why Monte Carlo dominates where grid-based methods drown. To halve the error you must *quadruple* the samples, the same $\sqrt{N}$ trade-off that haunts the [[Random Walk]].

{{< note kind="tip" title="Randomness as a feature, not a bug" >}}
It feels backwards to inject randomness into a deterministic question. But random sampling probes a huge space *evenly and cheaply* without ever building the whole thing. The noise is the price; the dimension-independence is the prize.
{{< /note >}}

## Estimating π with darts

Here is the cleanest Monte Carlo there is. Throw darts uniformly at a square, and draw the largest circle that fits inside it. The fraction of darts landing **inside the circle** equals the ratio of areas:

$$ \frac{\text{circle}}{\text{square}} = \frac{\pi r^2}{(2r)^2} = \frac{\pi}{4}. $$

So $\pi \approx 4 \times (\text{fraction inside})$. Below, darts rain down at random; hits inside the circle are colored, misses are dimmed, and the live estimate of $\pi$ ticks toward $3.14159\ldots$ as the count climbs. Watch it converge — and watch it converge *slowly*, jittering by less and less.

{{< sketch height="400" caption="Estimating π by throwing random darts at a square. The fraction landing inside the inscribed circle approaches π/4, so 4× that fraction estimates π. The live estimate sharpens like 1/√N. Auto-resets after many darts." >}}
if (frame === 0 || !state.pts) {
  state.pts = [];     // recent darts to draw
  state.inside = 0;
  state.total = 0;
  state.maxN = 6000;
}
const cs = getComputedStyle(document.documentElement);
const accent = cs.getPropertyValue('--accent').trim() || '#5b9cff';
const accent2 = cs.getPropertyValue('--accent-2').trim() || '#ff9e64';
const good = cs.getPropertyValue('--good').trim() || '#9ece6a';
const faint = cs.getPropertyValue('--text-faint').trim() || 'rgba(255,255,255,0.4)';
const border = cs.getPropertyValue('--border').trim() || 'rgba(255,255,255,0.18)';
ctx.clearRect(0, 0, W, H);
ctx.fillStyle = '#0c1018';
ctx.fillRect(0, 0, W, H);
// square region centered, side = S
const S = Math.min(W, H) - 60;
const ox = (W - S) / 2, oy = 36;
const r = S / 2, ccx = ox + r, ccy = oy + r;
// throw a batch of darts
const batch = 40;
for (let s = 0; s < batch && state.total < state.maxN; s++) {
  const x = Math.random() * 2 - 1;   // -1..1
  const y = Math.random() * 2 - 1;
  const hit = (x * x + y * y) <= 1;
  if (hit) state.inside++;
  state.total++;
  state.pts.push([ccx + x * r, ccy + y * r, hit]);
  if (state.pts.length > 3000) state.pts.shift();
}
// square + circle
ctx.strokeStyle = border;
ctx.lineWidth = 1.5;
ctx.strokeRect(ox, oy, S, S);
ctx.beginPath(); ctx.arc(ccx, ccy, r, 0, 7); ctx.stroke();
// darts
for (let i = 0; i < state.pts.length; i++) {
  const p = state.pts[i];
  ctx.fillStyle = p[2] ? accent : faint;
  ctx.globalAlpha = p[2] ? 0.7 : 0.4;
  ctx.fillRect(p[0], p[1], 2, 2);
}
ctx.globalAlpha = 1;
// estimate readout
const est = state.total ? 4 * state.inside / state.total : 0;
ctx.fillStyle = good;
ctx.font = 'bold 16px sans-serif';
ctx.fillText('π ≈ ' + est.toFixed(5), 12, 24);
ctx.fillStyle = faint;
ctx.font = '12px sans-serif';
ctx.fillText('darts: ' + state.total + '    inside: ' + state.inside + '    (true π = 3.14159)', 12, H - 12);
// error bar
ctx.fillStyle = accent2;
ctx.fillText('error: ' + Math.abs(est - Math.PI).toFixed(4), W - 130, 24);
if (state.total >= state.maxN) {
  state.hold = (state.hold || 0) + 1;
  if (state.hold > 130) { state.pts = []; state.inside = 0; state.total = 0; state.hold = 0; }
}
{{< /sketch >}}

## Where it shines

Monte Carlo wins whenever the problem is high-dimensional or has no closed form: pricing exotic financial derivatives, simulating particle physics, rendering realistic light in computer graphics, and powering the tree search in modern game engines. It is the computational embodiment of the [[Law of Large Numbers]] — and a reminder that a well-aimed random guess, repeated enough, beats a clever formula you cannot find.

## See also

- [[Law of Large Numbers]]
- [[Random Walk]]
- [[Central Limit Theorem]]
