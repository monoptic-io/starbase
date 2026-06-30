---
title: Fundamental Theorem of Arithmetic
aliases: [unique factorization, prime factorization]
tags: [number-theory]
summary: Every integer greater than 1 is a product of primes in exactly one way, up to the order of the factors.
weight: 30
---

# Fundamental Theorem of Arithmetic

The **Fundamental Theorem of Arithmetic** says that every integer greater than $1$ can be written as a product of [[Prime Number|primes]], and — crucially — that this factorization is **unique** apart from reordering. The number $60$ is $2^2 \times 3 \times 5$, and there is no other set of primes whose product is $60$. This is what justifies calling primes the *atoms of multiplication*: just as molecules decompose into elements in one definite way, integers decompose into primes in one definite way.

{{< eq number="1" >}}
n \;=\; p_1^{a_1}\,p_2^{a_2}\cdots p_k^{a_k}, \qquad p_1 < p_2 < \cdots < p_k \text{ prime}
{{< /eq >}}

Written in this **canonical form** — primes in increasing order with their exponents — the factorization is literally one of a kind. Two integers are equal if and only if their canonical forms match exponent for exponent.

## Build the tree

Any factorization can be grown as a **factor tree**: split the number into any two factors, then keep splitting the composite ones until only primes remain at the leaves. Start from different splits and the tree looks different — but the multiset of leaves is always the same.

{{< sketch height="380" caption="A factor tree for n. Composite nodes (blue) split into two factors; prime leaves glow green. Click to grow a tree for a new number — however you split, the prime leaves are always the same multiset." >}}
if (frame === 0) {
  state.targets = [60, 72, 84, 90, 120, 126, 96, 100];
  state.ti = 0;
  state.tree = null;
  state.built = false;
  state.t0 = 0;
}
ctx.clearRect(0, 0, W, H);
const css = (v, f) => (getComputedStyle(document.documentElement).getPropertyValue(v).trim() || f);
const accent = css('--accent', '#5b9cff');
const good = css('--good', '#4ec98f');
const warn = css('--warn', '#e0a458');
const text = css('--text', '#e6e9ef');
const border = css('--border', 'rgba(127,127,127,0.35)');
const faint = css('--text-faint', 'rgba(230,233,239,0.45)');

const smallestFactor = (m) => {
  for (let d = 2; d * d <= m; d++) if (m % d === 0) return d;
  return m; // prime
};
const build = (m) => {
  const f = smallestFactor(m);
  if (f === m) return { v: m, prime: true };
  return { v: m, prime: false, a: build(f), b: build(m / f) };
};

if (state.tree === null || mouse.clicked) {
  if (mouse.clicked) state.ti = (state.ti + 1) % state.targets.length;
  state.tree = build(state.targets[state.ti]);
  state.t0 = t;
}

// assign x positions by in-order leaf spread, depth = y
const W2 = W, H2 = H;
const layout = (node, depth, xleft, xright) => {
  node.y = 36 + depth * ((H2 - 70) / 4);
  if (node.prime) { node.x = (xleft + xright) / 2; return; }
  const mid = (xleft + xright) / 2;
  layout(node.a, depth + 1, xleft, mid);
  layout(node.b, depth + 1, mid, xright);
  node.x = (node.a.x + node.b.x) / 2;
};
layout(state.tree, 0, W2 * 0.08, W2 * 0.92);

// draw edges then nodes
ctx.lineWidth = 1.4; ctx.strokeStyle = border;
const drawEdges = (node) => {
  if (node.prime) return;
  for (const c of [node.a, node.b]) {
    ctx.beginPath(); ctx.moveTo(node.x, node.y + 12); ctx.lineTo(c.x, c.y - 12); ctx.stroke();
    drawEdges(c);
  }
};
drawEdges(state.tree);

ctx.textAlign = 'center'; ctx.textBaseline = 'middle'; ctx.font = '13px monospace';
const drawNodes = (node) => {
  const r = 15;
  ctx.beginPath(); ctx.arc(node.x, node.y, r, 0, 7);
  ctx.fillStyle = node.prime ? good : accent; ctx.globalAlpha = node.prime ? 0.9 : 0.85;
  ctx.fill(); ctx.globalAlpha = 1;
  ctx.fillStyle = '#0b0d12'; ctx.fillText(node.v, node.x, node.y);
  if (!node.prime) { drawNodes(node.a); drawNodes(node.b); }
};
drawNodes(state.tree);

// collect leaves into canonical form
const leaves = [];
const gather = (node) => { if (node.prime) leaves.push(node.v); else { gather(node.a); gather(node.b); } };
gather(state.tree);
leaves.sort((a, b) => a - b);
const counts = {};
for (const l of leaves) counts[l] = (counts[l] || 0) + 1;
let s = state.tree.v + ' = ';
s += Object.keys(counts).map(p => counts[p] > 1 ? p + '^' + counts[p] : '' + p).join(' · ');
ctx.font = '14px monospace'; ctx.fillStyle = text;
ctx.fillText(s, W2 / 2, H2 - 14);
ctx.font = '11px sans-serif'; ctx.fillStyle = faint;
ctx.fillText('click for another number', W2 / 2, 14);
{{< /sketch >}}

## Why uniqueness is the hard part

That every number *factors* into primes is easy: keep pulling out a factor until you can't. The deep claim is that the result is **unique**. The proof rests on **Euclid's lemma**:

{{< note kind="key" title="Euclid's lemma" >}}
If a prime $p$ divides a product $ab$, then $p$ divides $a$ or $p$ divides $b$. Primes cannot "hide" across a multiplication — they must land wholly inside one factor. From this, uniqueness follows: if two prime factorizations of $n$ existed, Euclid's lemma forces every prime in one to appear in the other, and matching them off shows the lists are identical.
{{< /note >}}

Uniqueness is not automatic — it can *fail* in other number systems. Among the numbers of the form $a + b\sqrt{-5}$, the integer $6$ factors two genuinely different ways, $6 = 2\cdot 3 = (1+\sqrt{-5})(1-\sqrt{-5})$, and neither factorization refines the other. That such worlds exist is exactly why the ordinary integers' unique factorization is a *theorem* worth naming.

## Reading a number from its primes

The canonical form is a fingerprint that makes many quantities instantly computable. If $n = p_1^{a_1}\cdots p_k^{a_k}$, then the **number of divisors** of $n$ is simply

$$d(n) = (a_1+1)(a_2+1)\cdots(a_k+1),$$

because each prime $p_i$ may appear in a divisor anywhere from $0$ to $a_i$ times. For $n = 60 = 2^2\cdot 3\cdot 5$ that gives $(2{+}1)(1{+}1)(1{+}1) = 12$ divisors. The same fingerprint feeds the [[Greatest Common Divisor]] (take the smaller exponent of each shared prime) and [[Euler's Totient Function]].

{{< quiz question="Using its prime factorization 72 = 2³·3², how many positive divisors does 72 have?" options="6|8|12|18" answer="3" explain="The divisor count is (3+1)(2+1) = 4·3 = 12. Each divisor chooses 0–3 copies of 2 and 0–2 copies of 3." >}}

## See also

- [[Prime Number]]
- [[Greatest Common Divisor]]
- [[Euler's Totient Function]]
