---
title: Sieve of Eratosthenes
aliases: [sieve]
tags: [number-theory]
summary: An ancient algorithm that finds every prime up to N by repeatedly striking out the multiples of each prime in turn.
weight: 20
---

# Sieve of Eratosthenes

The **Sieve of Eratosthenes** is a method, more than two thousand years old, for finding *all* the [[Prime Number|primes]] up to some limit $N$ — not by testing each number, but by **elimination**. Write down the integers from $2$ to $N$. Circle the first, $2$: it is prime. Now cross out every larger multiple of $2$. Move to the next surviving number, $3$: it is prime, so cross out its multiples. Repeat. Whatever is never crossed out is exactly the set of primes.

It is named for Eratosthenes of Cyrene, the polymath who also measured the circumference of the Earth. The sieve's charm is that it never performs a single division: primality falls out of pure marking.

## Watch it run

Below, the integers $2 \dots N$ are laid in a row. A scanning bar selects each new prime, then sweeps across striking out its multiples. The work shrinks pass by pass — by the time the current prime $p$ satisfies $p^2 > N$, every composite is already gone.

{{< sketch height="380" caption="The sieve in motion over 2…200. The ringed cell is the current prime; the moving marker strikes its multiples. Once p² > N the sieve halts — survivors are all prime. Click anywhere to run it again." >}}
if (frame === 0) {
  state.N = 200;
  state.alive = [];            // alive[i] = not yet struck
  for (let i = 0; i <= state.N; i++) state.alive.push(i >= 2);
  state.p = 2;
  state.cur = 4;               // current multiple being struck (p*2 to start)
  state.acc = 0;
  state.done = false;
  state.flash = 0;             // x of the marker for a brief highlight
}
ctx.clearRect(0, 0, W, H);
const css = (v, f) => (getComputedStyle(document.documentElement).getPropertyValue(v).trim() || f);
const accent = css('--accent', '#5b9cff');
const good = css('--good', '#4ec98f');
const warn = css('--warn', '#e0a458');
const text = css('--text', '#e6e9ef');
const faint = css('--text-faint', 'rgba(230,233,239,0.4)');

if (mouse.clicked) {
  state.alive = [];
  for (let i = 0; i <= state.N; i++) state.alive.push(i >= 2);
  state.p = 2; state.cur = 4; state.acc = 0; state.done = false;
}

// layout: a near-square grid sized to fill the canvas
const cols = Math.ceil(Math.sqrt(state.N * W / H));
const rows = Math.ceil(state.N / cols);
const pad = 8;
const cell = Math.min((W - 2 * pad) / cols, (H - 2 * pad - 20) / rows);
const ox = (W - cell * cols) / 2;
const oy = pad;
const cellOf = (v) => {
  const i = v - 1;
  return { x: ox + (i % cols) * cell, y: oy + Math.floor(i / cols) * cell };
};

// advance: strike one multiple per tick
state.acc += dt;
if (!state.done && state.acc > 0.045) {
  state.acc = 0;
  if (state.p * state.p > state.N) {
    state.done = true;
  } else if (state.cur > state.N) {
    state.p++;
    while (state.p <= state.N && !state.alive[state.p]) state.p++;
    state.cur = state.p * 2;
  } else {
    state.alive[state.cur] = false;
    state.flash = state.cur;
    state.cur += state.p;
  }
}

// draw cells
for (let v = 2; v <= state.N; v++) {
  const c = cellOf(v);
  const alive = state.alive[v];
  if (alive) {
    ctx.fillStyle = good; ctx.globalAlpha = (v === state.p && !state.done) ? 0.5 : 0.22;
    ctx.fillRect(c.x + 1, c.y + 1, cell - 2, cell - 2);
    ctx.globalAlpha = 1;
  } else {
    ctx.fillStyle = faint; ctx.globalAlpha = 0.10;
    ctx.fillRect(c.x + 1, c.y + 1, cell - 2, cell - 2);
    ctx.globalAlpha = 1;
  }
}
// ring the current prime
if (!state.done) {
  const c = cellOf(state.p);
  ctx.strokeStyle = warn; ctx.lineWidth = 2;
  ctx.strokeRect(c.x + 1.5, c.y + 1.5, cell - 3, cell - 3);
  // highlight the cell just struck
  if (state.flash >= 2 && state.flash <= state.N && state.cur <= state.N + state.p) {
    const f = cellOf(state.flash);
    ctx.strokeStyle = accent; ctx.lineWidth = 1.5;
    ctx.strokeRect(f.x + 1.5, f.y + 1.5, cell - 3, cell - 3);
  }
}
// status
ctx.textAlign = 'center'; ctx.textBaseline = 'middle'; ctx.font = '12px monospace';
ctx.fillStyle = text;
let primes = 0;
for (let v = 2; v <= state.N; v++) if (state.alive[v]) primes++;
ctx.fillText(state.done ? ('done — ' + primes + ' primes ≤ ' + state.N)
  : ('prime ' + state.p + ': striking multiples'), W / 2, H - 10);
{{< /sketch >}}

## Why it is fast

The sieve's genius is in *not* repeating work. Two optimizations make it efficient.

- **Start at $p^2$.** When you reach prime $p$, every multiple $2p, 3p, \dots, (p-1)p$ was already struck by a smaller prime factor. So the first *new* multiple to cross out is $p^2$.
- **Stop at $\sqrt{N}$.** Once $p^2 > N$, any remaining composite would need a prime factor larger than $\sqrt N$ paired with one smaller — but the smaller factor already eliminated it. Everything left is prime.

Counting the operations gives a running time of about $N \ln\ln N$ — almost linear in $N$. The doubly-logarithmic factor grows so slowly it is nearly a constant.

$$\sum_{p \le N}\frac{N}{p} \;\approx\; N\ln\ln N.$$

{{< chart type="bar" data="2:100, 3:67, 5:40, 7:29, 11:18" title="Cross-outs per prime, N = 200" xlabel="prime p" ylabel="multiples struck" caption="Each prime strikes roughly N/p numbers, so the work falls off sharply: 2 does the heavy lifting, and by 11 only a handful of cells remain to cross out." >}}

{{< note kind="key" title="Elimination beats testing" >}}
To list the primes up to $N$, the sieve is far faster than checking each number for primality one at a time. It trades arithmetic for memory: keep an array of $N$ flags and just walk through it striking multiples. This same idea — pay with space to save time — recurs throughout the [[Fundamental Theorem of Arithmetic|theory of factorization]].
{{< /note >}}

{{< quiz question="When the sieve reaches a new prime p, why can it begin crossing out at p² instead of 2p?" options="Because 2p is always prime|Because every multiple of p below p² already has a smaller prime factor and was struck|Because p² is the largest multiple of p|Because the sieve skips even numbers" answer="2" explain="Any multiple kp with k < p has a prime factor of k that is smaller than p, so that earlier prime already struck it. The first uniquely-new multiple is p·p." >}}

## See also

- [[Prime Number]]
- [[Fundamental Theorem of Arithmetic]]
- [[Number Theory]]
