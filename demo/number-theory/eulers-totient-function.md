---
title: Euler's Totient Function
aliases: [totient, phi function, euler phi]
tags: [number-theory]
summary: φ(n) counts the integers from 1 to n that share no factor with n — the size of the group of units that powers RSA.
weight: 50
---

# Euler's Totient Function

**Euler's totient function**, written $\varphi(n)$, counts how many integers from $1$ to $n$ are **coprime** to $n$ — that is, share no common factor with it beyond $1$. For $n = 12$, the integers coprime to it are $1, 5, 7, 11$, so $\varphi(12) = 4$. These are exactly the numbers that have a multiplicative inverse in [[Modular Arithmetic|arithmetic mod $n$]], which makes $\varphi$ the quiet engine inside [[RSA]].

## Counting the coprime residues

Below, the integers $1 \dots n$ are arranged around a ring. Those coprime to $n$ — sharing only a [[Greatest Common Divisor|gcd]] of $1$ — light up; the rest dim. Their count is $\varphi(n)$.

{{< sketch height="380" caption="Residues 1…n around a ring. Highlighted points are coprime to n (gcd = 1); dim points share a factor. The lit count is φ(n). Drag left↔right to change n and watch φ jump around — high for primes, low for highly composite numbers." >}}
if (frame === 0) {
  state.n = 12;
}
ctx.clearRect(0, 0, W, H);
const css = (v, f) => (getComputedStyle(document.documentElement).getPropertyValue(v).trim() || f);
const accent = css('--accent', '#5b9cff');
const good = css('--good', '#4ec98f');
const warn = css('--warn', '#e0a458');
const text = css('--text', '#e6e9ef');
const faint = css('--text-faint', 'rgba(230,233,239,0.4)');

// n from mouse x, range 3..36
if (mouse.x >= 0 && mouse.x <= W) {
  state.n = 3 + Math.floor(mouse.x / W * 33);
  state.n = Math.max(3, Math.min(36, state.n));
}
const n = state.n;
const gcd = (a, b) => { while (b) { [a, b] = [b, a % b]; } return a; };

const cx = W * 0.40, cy = H * 0.50, R = Math.min(W * 0.32, H * 0.40);
let phi = 0;
ctx.textAlign = 'center'; ctx.textBaseline = 'middle';
ctx.font = Math.max(8, Math.min(12, 320 / n)) + 'px monospace';
for (let i = 1; i <= n; i++) {
  const ang = -Math.PI / 2 + (i - 1) / n * Math.PI * 2;
  const x = cx + Math.cos(ang) * R, y = cy + Math.sin(ang) * R;
  const co = gcd(i, n) === 1;
  if (co) phi++;
  ctx.beginPath(); ctx.arc(x, y, co ? 8 : 5, 0, 7);
  ctx.fillStyle = co ? good : 'rgba(127,127,127,0.25)'; ctx.fill();
  ctx.fillStyle = co ? '#0b0d12' : faint;
  if (n <= 30) ctx.fillText(i, x, y);
}
// readout
const rx = W * 0.74;
ctx.textAlign = 'left'; ctx.font = '14px monospace';
ctx.fillStyle = text; ctx.fillText('n = ' + n, rx, H * 0.40);
ctx.fillStyle = good; ctx.fillText('φ(' + n + ') = ' + phi, rx, H * 0.50);
ctx.fillStyle = faint; ctx.font = '11px sans-serif';
ctx.fillText('drag to change n', rx, H * 0.62);
{{< /sketch >}}

## A product formula from the primes

You never have to count one by one. The totient is **multiplicative** — $\varphi(ab) = \varphi(a)\varphi(b)$ whenever $a$ and $b$ are coprime — and for a prime power $\varphi(p^k) = p^k - p^{k-1}$. Combining these via the [[Fundamental Theorem of Arithmetic]] gives a clean product over the distinct primes dividing $n$:

{{< eq number="1" >}}
\varphi(n) \;=\; n\prod_{p \mid n}\left(1 - \frac{1}{p}\right)
{{< /eq >}}

For $n = 12 = 2^2\cdot 3$: $\varphi(12) = 12\left(1-\tfrac12\right)\left(1-\tfrac13\right) = 12\cdot\tfrac12\cdot\tfrac23 = 4$ — matching the four points we lit up. The two special cases worth memorizing:

- **For a prime $p$:** $\varphi(p) = p - 1$, since every smaller number is coprime to a prime.
- **For two primes:** $\varphi(pq) = (p-1)(q-1)$ — the exact quantity [[RSA]] needs.

{{< chart type="bar" data="9:6, 10:4, 11:10, 12:4, 13:12, 14:6, 15:8, 16:8" title="φ(n) for n = 9…16" xlabel="n" ylabel="φ(n)" caption="The totient spikes at primes (φ(11)=10, φ(13)=12 — one less than n) and dips for numbers rich in small factors (φ(12)=4). It is jagged, not smooth — a direct readout of n's factor structure." >}}

## Euler's theorem and RSA

The reason $\varphi$ rules cryptography is **Euler's theorem**, the generalization of [[Fermat's Little Theorem]] to any modulus:

$$a^{\varphi(n)} \equiv 1 \pmod{n} \qquad\text{whenever } \gcd(a, n) = 1.$$

This says the coprime residues, under multiplication, form a group of size $\varphi(n)$, so raising to that power returns to the identity. [[RSA]] chooses its public and private exponents $e, d$ so that $ed \equiv 1 \pmod{\varphi(n)}$; Euler's theorem then guarantees $m^{ed} \equiv m$, so decryption perfectly undoes encryption.

{{< note kind="key" title="Why factoring breaks RSA" >}}
To compute the private key you need $\varphi(n) = (p-1)(q-1)$. That requires knowing $p$ and $q$ — i.e. **factoring** $n$. Anyone who could factor large numbers could compute $\varphi$ and read every RSA message. The security of the modern web rests on $\varphi$ being easy to compute *if you know the primes* and apparently impossible *if you don't*.
{{< /note >}}

{{< quiz question="If p and q are distinct primes, what is φ(pq)?" options="pq − 1|p + q|(p − 1)(q − 1)|pq − p − q" answer="3" explain="φ is multiplicative on coprime parts and φ(p) = p−1, so φ(pq) = (p−1)(q−1). This is exactly the modulus RSA uses to derive its private exponent." >}}

## See also

- [[RSA]]
- [[Fermat's Little Theorem]]
- [[Prime Number]]
