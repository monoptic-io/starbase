---
title: Prime Number
aliases: [primes, prime]
tags: [number-theory]
summary: An integer greater than 1 divisible only by 1 and itself — the indivisible building block from which every other integer is made.
weight: 10
---

# Prime Number

A **prime number** is a whole number greater than $1$ whose only divisors are $1$ and itself: $2, 3, 5, 7, 11, 13, \dots$ Everything else — $4 = 2\times2$, $6 = 2\times3$, $12 = 2\times2\times3$ — is **composite**, built by multiplying smaller pieces. The primes are the pieces that cannot be broken further. They are, quite literally, the *atoms of multiplication*, and the [[Fundamental Theorem of Arithmetic]] makes that metaphor a theorem.

{{< note kind="note" title="Why 1 is not prime" >}}
The number $1$ is deliberately excluded. If it counted as prime, factorizations would no longer be unique — $6 = 2\times3 = 1\times2\times3 = 1\times1\times2\times3$ and so on forever. Banishing $1$ keeps the [[Fundamental Theorem of Arithmetic]] clean: every integer has *exactly one* prime factorization.
{{< /note >}}

## The sieve made visible

The fastest way to *see* primality is to eliminate everything that isn't prime. Lay the integers out in a grid and cross out every multiple of $2$, then every multiple of $3$, then $5$, and so on. Whatever survives is prime. This is the [[Sieve of Eratosthenes]], running live below.

{{< sketch height="420" caption="The Sieve of Eratosthenes on a 10-wide grid of 1–120. Each pass picks the next surviving prime (ringed) and strikes out its multiples. Survivors glow; struck numbers fade. Click to restart the sieve." >}}
if (frame === 0) {
  state.N = 120;
  state.cols = 10;
  state.struck = [];           // struck[i] = true if eliminated
  for (let i = 0; i <= state.N; i++) state.struck.push(i < 2);
  state.p = 2;                  // current prime being processed
  state.m = 0;                  // current multiple index (m*p)
  state.acc = 0;
  state.done = false;
}
ctx.clearRect(0, 0, W, H);
const css = (v, f) => (getComputedStyle(document.documentElement).getPropertyValue(v).trim() || f);
const accent = css('--accent', '#5b9cff');
const good = css('--good', '#4ec98f');
const warn = css('--warn', '#e0a458');
const text = css('--text', '#e6e9ef');
const faint = css('--text-faint', 'rgba(230,233,239,0.45)');
const border = css('--border', 'rgba(127,127,127,0.3)');

if (mouse.clicked) {            // restart the sieve
  state.struck = [];
  for (let i = 0; i <= state.N; i++) state.struck.push(i < 2);
  state.p = 2; state.m = 0; state.acc = 0; state.done = false;
}

// advance the sieve on a timer
state.acc += dt;
if (!state.done && state.acc > 0.10) {
  state.acc = 0;
  if (state.p * state.p > state.N) {
    state.done = true;
  } else {
    state.m++;
    const v = state.p * state.m;
    if (v <= state.N) {
      if (v !== state.p) state.struck[v] = true;
    } else {
      // next prime
      state.p++;
      while (state.p <= state.N && state.struck[state.p]) state.p++;
      state.m = 0;
    }
  }
}

// grid geometry
const rows = Math.ceil((state.N) / state.cols);
const pad = Math.min(W, H) * 0.04;
const gw = (W - pad * 2) / state.cols;
const gh = (H - pad * 2 - 22) / rows;
const cell = Math.min(gw, gh);
const ox = (W - cell * state.cols) / 2;
const oy = pad;
ctx.textAlign = 'center'; ctx.textBaseline = 'middle';
ctx.font = Math.floor(cell * 0.42) + 'px monospace';
for (let v = 1; v <= state.N; v++) {
  const idx = v - 1;
  const cxp = ox + (idx % state.cols) * cell + cell / 2;
  const cyp = oy + Math.floor(idx / state.cols) * cell + cell / 2;
  const struck = state.struck[v];
  const isCurrent = v === state.p && !state.done;
  // cell background
  if (!struck) {
    ctx.fillStyle = good; ctx.globalAlpha = 0.16;
    ctx.fillRect(cxp - cell / 2 + 1, cyp - cell / 2 + 1, cell - 2, cell - 2);
    ctx.globalAlpha = 1;
  }
  if (isCurrent) {
    ctx.strokeStyle = warn; ctx.lineWidth = 2;
    ctx.beginPath(); ctx.arc(cxp, cyp, cell * 0.42, 0, 7); ctx.stroke();
  }
  ctx.fillStyle = struck ? faint : (v < 2 ? faint : good);
  ctx.fillText(v, cxp, cyp);
  if (struck && v >= 2) {
    ctx.strokeStyle = border; ctx.lineWidth = 1;
    ctx.beginPath(); ctx.moveTo(cxp - cell * 0.3, cyp + cell * 0.3);
    ctx.lineTo(cxp + cell * 0.3, cyp - cell * 0.3); ctx.stroke();
  }
}
// status line
ctx.textAlign = 'center'; ctx.font = '12px monospace';
ctx.fillStyle = text;
let primes = 0;
for (let v = 2; v <= state.N; v++) if (!state.struck[v]) primes++;
const msg = state.done ? (primes + ' primes ≤ ' + state.N) : ('striking multiples of ' + state.p);
ctx.fillText(msg, W / 2, H - 11);
{{< /sketch >}}

## Infinitely many — Euclid's proof

There is no largest prime. Euclid proved it around 300 BC with one of the most elegant arguments in all mathematics, a **proof by contradiction**.

Suppose the primes were a finite list $p_1, p_2, \dots, p_k$. Form the number

$$Q = p_1 p_2 \cdots p_k + 1.$$

Now $Q$ leaves remainder $1$ when divided by *every* prime on the list, so none of them divides it. But every integer above $1$ has at least one prime factor (by the [[Fundamental Theorem of Arithmetic]]). That factor is a prime not on our list — contradiction. The list can never be complete, so the primes go on forever.

{{< note kind="tip" title="Q itself need not be prime" >}}
A common misreading is that $Q = p_1\cdots p_k + 1$ is always prime. It isn't — for example $2\cdot3\cdot5\cdot7\cdot11\cdot13 + 1 = 30031 = 59 \times 509$. The proof only needs that $Q$ has *some* prime factor outside the list, which is all the contradiction requires.
{{< /note >}}

## Thinning out, but never stopping

Primes become rarer as you climb, but only slowly. The **prime counting function** $\pi(x)$ — how many primes are $\le x$ — grows like $x / \ln x$, the celebrated **Prime Number Theorem**. The density near $x$ is roughly $1/\ln x$, so primes thin out logarithmically: never absent, never regular.

{{< chart type="line" data="10:4, 50:15, 100:25, 500:95, 1000:168, 5000:669" title="π(x): primes up to x" xlabel="x" ylabel="number of primes" caption="The count of primes climbs without bound but ever more slowly — the signature of the Prime Number Theorem, π(x) ≈ x / ln x." >}}

This unpredictable-yet-structured distribution is exactly what makes primes useful in [[RSA]]: it is easy to *find* large primes (they are common enough), but hard to *factor* their products.

{{< quiz question="Why does Euclid's argument show there is no largest prime?" options="Because p₁⋯pₖ + 1 is always itself prime|Because p₁⋯pₖ + 1 has a prime factor that cannot be on any finite list|Because primes get denser as numbers grow|Because every even number is composite" answer="2" explain="Q = p₁⋯pₖ + 1 is coprime to every listed prime, yet must have a prime factor — necessarily a new one. So no finite list can contain them all." >}}

## See also

- [[Sieve of Eratosthenes]]
- [[Fundamental Theorem of Arithmetic]]
- [[RSA]]
