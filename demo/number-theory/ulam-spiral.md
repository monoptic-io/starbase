---
title: Ulam Spiral
aliases: [prime spiral]
tags: [number-theory]
summary: Writing the integers in a square spiral and marking the primes reveals startling diagonal alignments where order seems to leak out of randomness.
weight: 70
---

# Ulam Spiral

The **Ulam spiral** is one of the most famous accidental discoveries in mathematics. In 1963, bored during a lecture, Stanisław Ulam began writing the integers in a square spiral and idly circling the [[Prime Number|primes]]. He expected a random scatter. Instead the primes lined up along **diagonals** — long, conspicuous streaks that have no business appearing in something as irregular as the primes. The pattern is real, partial, and still not fully explained.

## Coil the integers, mark the primes

Start with $1$ at the center and spiral outward: $2$ to its right, then up, left, down, around and around. Color a cell whenever its number is prime. The diagonals emerge on their own.

{{< sketch height="440" caption="The Ulam spiral. Integers coil out from the center; primes are lit. Notice the diagonal streaks — primes clustering along lines like 4k²+something. The spiral grows outward as it fills; click to restart and rebuild it." >}}
if (frame === 0) {
  state.max = 0;          // how many numbers placed so far
  state.acc = 0;
  state.isPrime = (m) => {
    if (m < 2) return false;
    if (m % 2 === 0) return m === 2;
    for (let d = 3; d * d <= m; d += 2) if (m % d === 0) return false;
    return true;
  };
}
ctx.clearRect(0, 0, W, H);
const css = (v, f) => (getComputedStyle(document.documentElement).getPropertyValue(v).trim() || f);
const accent = css('--accent', '#5b9cff');
const good = css('--good', '#4ec98f');
const warn = css('--warn', '#e0a458');
const text = css('--text', '#e6e9ef');
const faint = css('--text-faint', 'rgba(230,233,239,0.4)');

if (mouse.clicked) state.max = 0;

// grid size (odd) chosen to fit the canvas
const margin = 24;
const side = Math.min(W, H - 16) - 2 * margin;
const G = 41;                       // 41×41 spiral → numbers 1..1681
const cell = side / G;
const ox = (W - cell * G) / 2;
const oy = (H - 16 - cell * G) / 2 + 4;
const total = G * G;

// grow the fill over time
state.acc += dt;
if (state.acc > 0.012 && state.max < total) {
  state.acc = 0;
  state.max = Math.min(total, state.max + Math.ceil(total / 240));
}

// spiral coordinate of value v (1-based): returns grid (gx, gy) centered
// walk the spiral; cache positions once
if (!state.pos || state.pos.length < total + 1) {
  state.pos = new Array(total + 1);
  let x = 0, y = 0, v = 1;
  state.pos[1] = { x, y };
  let step = 1, dir = 0; // dir: 0=right,1=up,2=left,3=down
  const dx = [1, 0, -1, 0], dy = [0, -1, 1, 0];
  while (v < total) {
    for (let leg = 0; leg < 2 && v < total; leg++) {
      for (let s = 0; s < step && v < total; s++) {
        x += dx[dir]; y += dy[dir]; v++;
        state.pos[v] = { x, y };
      }
      dir = (dir + 1) % 4;
    }
    step++;
  }
}

const c0 = (G - 1) / 2;            // center grid index
for (let v = 1; v <= state.max; v++) {
  if (!state.isPrime(v)) continue;
  const p = state.pos[v];
  const gx = c0 + p.x, gy = c0 + p.y;
  if (gx < 0 || gx >= G || gy < 0 || gy >= G) continue;
  const px = ox + gx * cell, py = oy + gy * cell;
  ctx.fillStyle = good; ctx.globalAlpha = 0.9;
  ctx.fillRect(px + 0.5, py + 0.5, cell - 1, cell - 1);
}
ctx.globalAlpha = 1;
// faint frame + status
ctx.strokeStyle = faint; ctx.lineWidth = 1;
ctx.strokeRect(ox, oy, cell * G, cell * G);
ctx.fillStyle = text; ctx.textAlign = 'center'; ctx.font = '12px monospace';
ctx.fillText('primes up to ' + state.max + ' / ' + total, W / 2, H - 6);
{{< /sketch >}}

## Why the diagonals?

The diagonals are not an illusion — they trace **quadratic polynomials**. Moving along a diagonal of the spiral steps you through values of an expression like $4k^2 + bk + c$. Some of these quadratics are unusually rich in primes; the most celebrated is **Euler's polynomial**

$$k^2 + k + 41,$$

which is prime for every $k$ from $0$ to $39$ — forty primes in a row. Each prime-dense quadratic paints one bright diagonal across the spiral.

{{< chart type="bar" data="0:41, 1:43, 2:47, 3:53, 4:61, 5:71, 6:83, 7:97" title="Euler's polynomial k² + k + 41 (all prime)" xlabel="k" ylabel="value" caption="Every value shown is prime. This single quadratic produces 40 consecutive primes (k = 0…39), and along the Ulam spiral it lights up as one of the most striking diagonals." >}}

{{< note kind="note" title="Order from the absence of small factors" >}}
A diagonal stays prime-rich when its quadratic rarely produces multiples of small primes. Euler's $k^2+k+41$ is, for instance, *never* divisible by any prime below $41$ — so it slips past the [[Sieve of Eratosthenes]]' early cuts far more often than a random number would. The diagonals are a visible echo of which quadratics dodge small factors. They are real structure, yet they do not let us *predict* primes — the spiral teases order without delivering a formula.
{{< /note >}}

## Same primes, different canvas

The effect is not unique to the square spiral. Plot the primes on a **polar** spiral, or with a hexagonal coiling, and different alignments appear or dissolve — proof that the diagonals are partly an artifact of *how* we arrange the integers, and partly a genuine fact about prime-dense quadratics. The honest summary: the primes are neither random nor regular, and the Ulam spiral is the clearest picture we have of that in-between.

{{< quiz question="What mathematical objects do the diagonal streaks in the Ulam spiral correspond to?" options="Arithmetic sequences like 2k|Prime-rich quadratic polynomials like k²+k+41|Powers of 2|Fibonacci numbers" answer="2" explain="Diagonals of the spiral step through values of quadratics a·k²+b·k+c. Some quadratics, like Euler's k²+k+41, avoid small prime factors and so stay unusually prime-dense — lighting up as diagonals." >}}

## See also

- [[Prime Number]]
- [[Sieve of Eratosthenes]]
- [[Number Theory]]
