---
title: Greatest Common Divisor
aliases: [gcd, euclidean algorithm]
tags: [number-theory]
summary: The largest integer dividing two numbers, found in a handful of steps by Euclid's algorithm of repeated remainders.
weight: 40
---

# Greatest Common Divisor

The **greatest common divisor** of two integers $a$ and $b$, written $\gcd(a,b)$, is the largest number that divides both. For $48$ and $36$ it is $12$; for $17$ and $5$ it is $1$ (they share no factor but the trivial one, and are called **coprime**). The gcd is the meeting point of two numbers' factor structure — and there is a beautiful, ancient algorithm to find it without ever factoring either number.

## The Euclidean algorithm

Euclid's insight is that the gcd doesn't change if you replace the larger number by its remainder against the smaller:

$$\gcd(a, b) = \gcd(b,\; a \bmod b), \qquad \gcd(a, 0) = a.$$

Repeat until the remainder hits $0$; the last nonzero value is the gcd. For $\gcd(48, 36)$:

$$48 = 1\cdot 36 + 12, \quad 36 = 3\cdot 12 + 0 \;\Rightarrow\; \gcd = 12.$$

It is breathtakingly fast — the number of steps is at worst proportional to the *number of digits*, never the size of the numbers themselves.

## The geometry: tiling a rectangle with squares

There is a picture hiding in the algorithm. Lay out an $a \times b$ rectangle and repeatedly cut off the **largest square** that fits. Each cut removes a $b \times b$ square and leaves a smaller rectangle — exactly the step $a \to a \bmod b$. When the leftover is a perfect square, its side *is* $\gcd(a,b)$: the largest square that tiles the whole rectangle evenly.

{{< sketch height="380" caption="Carving an a × b rectangle into the largest squares that fit — the geometric Euclidean algorithm. The final, smallest square's side is gcd(a, b). Click to try a new pair of dimensions." >}}
if (frame === 0) {
  state.pairs = [[48, 36], [55, 34], [60, 24], [40, 25], [63, 36], [50, 35]];
  state.pi = 0;
  state.squares = null;
  state.reveal = 0;     // how many squares are shown
  state.acc = 0;
}
ctx.clearRect(0, 0, W, H);
const css = (v, f) => (getComputedStyle(document.documentElement).getPropertyValue(v).trim() || f);
const accent = css('--accent', '#5b9cff');
const accent2 = css('--accent-2', '#b07bff');
const good = css('--good', '#4ec98f');
const warn = css('--warn', '#e0a458');
const text = css('--text', '#e6e9ef');
const faint = css('--text-faint', 'rgba(230,233,239,0.5)');

const compute = (a, b) => {
  // produce the list of squares (in original a×b coordinates)
  const sq = [];
  let ox = 0, oy = 0, rw = a, rh = b;
  while (rw > 0 && rh > 0) {
    const s = Math.min(rw, rh);
    if (rw >= rh) {
      const n = Math.floor(rw / rh);
      for (let i = 0; i < n; i++) sq.push({ x: ox + i * s, y: oy, s });
      ox += n * s; rw -= n * s;
    } else {
      const n = Math.floor(rh / rw);
      for (let i = 0; i < n; i++) sq.push({ x: ox, y: oy + i * s, s });
      oy += n * s; rh -= n * s;
    }
  }
  return sq;
};

if (state.squares === null || mouse.clicked) {
  if (mouse.clicked) state.pi = (state.pi + 1) % state.pairs.length;
  const [a, b] = state.pairs[state.pi];
  state.a = a; state.b = b;
  state.squares = compute(a, b);
  state.reveal = 0; state.acc = 0;
}

// reveal squares over time
state.acc += dt;
if (state.acc > 0.45 && state.reveal < state.squares.length) {
  state.acc = 0; state.reveal++;
}

// scale a×b into the canvas
const pad = 30;
const sc = Math.min((W - 2 * pad) / state.a, (H - 2 * pad - 14) / state.b);
const ox = (W - state.a * sc) / 2;
const oy = (H - 14 - state.b * sc) / 2;
// outer rectangle
ctx.strokeStyle = faint; ctx.lineWidth = 1;
ctx.strokeRect(ox, oy, state.a * sc, state.b * sc);
// squares
ctx.textAlign = 'center'; ctx.textBaseline = 'middle';
const gcd = state.squares.length ? state.squares[state.squares.length - 1].s : 1;
for (let i = 0; i < state.reveal; i++) {
  const q = state.squares[i];
  const isGcd = q.s === gcd;
  ctx.fillStyle = isGcd ? warn : (i % 2 ? accent : accent2);
  ctx.globalAlpha = isGcd ? 0.55 : 0.22;
  ctx.fillRect(ox + q.x * sc, oy + q.y * sc, q.s * sc, q.s * sc);
  ctx.globalAlpha = 1;
  ctx.strokeStyle = isGcd ? warn : accent; ctx.lineWidth = isGcd ? 2 : 1;
  ctx.strokeRect(ox + q.x * sc, oy + q.y * sc, q.s * sc, q.s * sc);
  if (q.s * sc > 18) {
    ctx.fillStyle = text; ctx.font = Math.min(13, q.s * sc * 0.4) + 'px monospace';
    ctx.fillText(q.s, ox + (q.x + q.s / 2) * sc, oy + (q.y + q.s / 2) * sc);
  }
}
// status
ctx.font = '13px monospace'; ctx.fillStyle = text; ctx.textAlign = 'center';
ctx.fillText('gcd(' + state.a + ', ' + state.b + ') = ' + gcd, W / 2, H - 9);
{{< /sketch >}}

The number of squares of each size corresponds exactly to the quotients in Euclid's divisions — the rectangle *is* the algorithm drawn out in tile.

## Bézout and the extended algorithm

Run Euclid backwards and something extra falls out: the gcd can always be written as an integer combination of the originals.

{{< note kind="key" title="Bézout's identity" >}}
For any integers $a, b$ there exist integers $x, y$ with
$$a x + b y = \gcd(a, b).$$
For $48$ and $36$: $48(-2) + 36(3) = -96 + 108 = 12$. The **extended Euclidean algorithm** tracks these coefficients as it goes. They are what compute the modular inverses behind [[RSA]] keys — finding $d$ with $ed \equiv 1$ is exactly solving $ed + \varphi y = 1$.
{{< /note >}}

Bézout's identity is also the gateway to the [[Diophantine Equation]]: the equation $ax + by = c$ has integer solutions **if and only if** $\gcd(a,b)$ divides $c$. The gcd is the finest "grid spacing" the combination $ax+by$ can land on.

{{< note kind="tip" title="Coprime = building block of the totient" >}}
When $\gcd(a, b) = 1$ the numbers are **coprime**. Counting how many integers below $n$ are coprime to it is precisely [[Euler's Totient Function]] — and the link to [[Modular Arithmetic]] is direct: $a$ has a multiplicative inverse mod $n$ exactly when $\gcd(a, n) = 1$.
{{< /note >}}

{{< quiz question="The equation 6x + 4y = 9 has how many integer solutions (x, y)?" options="Exactly one|Infinitely many|None|Exactly two" answer="3" explain="gcd(6,4) = 2, and 2 does not divide 9. Since 6x + 4y is always even, it can never equal the odd number 9 — no solutions." >}}

## See also

- [[Modular Arithmetic]]
- [[Diophantine Equation]]
- [[Euler's Totient Function]]
