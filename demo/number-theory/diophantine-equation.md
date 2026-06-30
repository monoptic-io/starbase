---
title: Diophantine Equation
aliases: [diophantine equations]
tags: [number-theory]
summary: An equation whose solutions must be whole numbers — solvable in the linear case exactly when the gcd of the coefficients divides the constant.
weight: 90
---

# Diophantine Equation

A **Diophantine equation** is a polynomial equation for which we demand **integer** solutions — no fractions, no decimals, just whole numbers. Named for Diophantus of Alexandria (c. 250 AD), they are among the oldest and hardest objects in mathematics. The constraint "integers only" changes everything: an equation with infinitely many real solutions may have none at all in integers, or exactly one, or a whole lattice of them.

## The linear case is completely solved

The simplest family is the **linear Diophantine equation** in two unknowns:

$$ax + by = c.$$

Geometrically this is a straight line; we are asking which of its points land on the integer grid. The answer is a clean theorem resting entirely on the [[Greatest Common Divisor]]:

{{< note kind="key" title="Solvability criterion" >}}
$ax + by = c$ has integer solutions **if and only if** $\gcd(a, b)$ divides $c$.

The reason: every combination $ax + by$ is automatically a multiple of $g = \gcd(a,b)$, so it can only ever equal multiples of $g$. And by **Bézout's identity** ([[Greatest Common Divisor]]) it can reach $g$ itself — hence every multiple of $g$. The reachable values are *exactly* the multiples of the gcd.
{{< /note >}}

So $6x + 9y = 21$ is solvable because $\gcd(6,9) = 3$ divides $21$; but $6x + 9y = 20$ has **no** integer solution, since the left side is always a multiple of $3$ and $20$ is not.

## The lattice of solutions

When solutions exist, there are infinitely many, evenly spaced along the line. If $(x_0, y_0)$ is one solution, the rest are

$$x = x_0 + \frac{b}{g}t, \qquad y = y_0 - \frac{a}{g}t, \qquad t \in \mathbb{Z}.$$

The integer points sit at regular intervals — a one-dimensional lattice marching along the line. The sketch shows which lines catch grid points and which slip between them entirely.

{{< sketch height="380" caption="The line ax + by = c over the integer grid. Green dots are integer solutions, evenly spaced along the line. Drag left↔right to change c: when gcd(a,b) divides c the line threads through lattice points; otherwise it weaves between them with none." >}}
if (frame === 0) {
  state.a = 2; state.b = 3;
  state.c = 12;
}
ctx.clearRect(0, 0, W, H);
const css = (v, f) => (getComputedStyle(document.documentElement).getPropertyValue(v).trim() || f);
const accent = css('--accent', '#5b9cff');
const good = css('--good', '#4ec98f');
const warn = css('--warn', '#e0a458');
const text = css('--text', '#e6e9ef');
const faint = css('--text-faint', 'rgba(230,233,239,0.45)');
const border = css('--border', 'rgba(127,127,127,0.3)');

const a = state.a, b = state.b;
if (mouse.x >= 0 && mouse.x <= W) state.c = -6 + Math.round(mouse.x / W * 24);
const c = state.c;
const gcd = (m, n) => { m = Math.abs(m); n = Math.abs(n); while (n) { [m, n] = [n, m % n]; } return m; };
const g = gcd(a, b);
const solvable = (c % g === 0);

// grid: x,y in [-6,6]
const R = 6;
const ox = W / 2, oy = H / 2;
const u = Math.min(W, H - 30) / (2 * R + 2);
const sx = (gx) => ox + gx * u;
const sy = (gy) => oy - gy * u;
// axes + dots
ctx.strokeStyle = border; ctx.lineWidth = 1;
ctx.beginPath(); ctx.moveTo(sx(-R), sy(0)); ctx.lineTo(sx(R), sy(0));
ctx.moveTo(sx(0), sy(-R)); ctx.lineTo(sx(0), sy(R)); ctx.stroke();
ctx.fillStyle = faint;
for (let gx = -R; gx <= R; gx++) for (let gy = -R; gy <= R; gy++) {
  ctx.beginPath(); ctx.arc(sx(gx), sy(gy), 1.6, 0, 7); ctx.fill();
}
// the line ax+by=c  ->  y = (c - a x)/b
ctx.strokeStyle = solvable ? good : warn; ctx.lineWidth = 2;
ctx.beginPath();
let started = false;
for (let px = -R; px <= R; px += 0.05) {
  const py = (c - a * px) / b;
  if (py < -R - 1 || py > R + 1) { started = false; continue; }
  if (!started) { ctx.moveTo(sx(px), sy(py)); started = true; }
  else ctx.lineTo(sx(px), sy(py));
}
ctx.stroke();
// integer solutions: x from -R..R where (c - a x) divisible by b
ctx.fillStyle = good;
for (let gx = -R; gx <= R; gx++) {
  if ((c - a * gx) % b === 0) {
    const gy = (c - a * gx) / b;
    if (gy >= -R && gy <= R) {
      ctx.beginPath(); ctx.arc(sx(gx), sy(gy), 5.5, 0, 7); ctx.fill();
    }
  }
}
// label
ctx.fillStyle = text; ctx.font = '13px monospace'; ctx.textAlign = 'center';
ctx.fillText(a + 'x + ' + b + 'y = ' + c + '   (gcd ' + g + (solvable ? ' | c : solvable' : ' ∤ c : no solutions') + ')', W / 2, H - 9);
{{< /sketch >}}

## Beyond linear: a different universe

Raise the degree and the difficulty explodes. **Pythagorean triples** solve $x^2 + y^2 = z^2$ and form an infinite, beautifully structured family ($3,4,5$; $5,12,13$; …). But change the exponent and everything breaks:

{{< note kind="note" title="Fermat's Last Theorem" >}}
The equation $x^n + y^n = z^n$ has **no** positive-integer solutions for any $n > 2$. Fermat scribbled this claim in a margin around 1637, claiming a proof he had no room to write. It resisted every mathematician for **358 years** until Andrew Wiles proved it in 1995 — using machinery (elliptic curves, modular forms) light-years beyond anything Fermat could have known.
{{< /note >}}

There is no general algorithm that can decide whether an *arbitrary* Diophantine equation has solutions — this is the negative answer to **Hilbert's tenth problem**. The linear case we solved completely is the calm shore of an ocean that quickly becomes unfathomable.

{{< chart type="bar" data="3-4-5:5, 5-12-13:13, 8-15-17:17, 7-24-25:25, 20-21-29:29" title="Pythagorean triples: hypotenuse of x² + y² = z²" ylabel="z" caption="Integer solutions of x² + y² = z² form an infinite structured family — yet the moment the exponent passes 2, Fermat's Last Theorem says the integer solutions vanish entirely." >}}

{{< quiz question="Does the equation 4x + 6y = 9 have any integer solutions?" options="Yes, exactly one|Yes, infinitely many|No|Only if x = y" answer="3" explain="gcd(4, 6) = 2, which does not divide 9. Since 4x + 6y is always even, it can never equal the odd number 9 — there are no integer solutions." >}}

## See also

- [[Greatest Common Divisor]]
- [[Fundamental Theorem of Arithmetic]]
- [[Continued Fraction]]
