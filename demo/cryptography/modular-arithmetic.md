---
title: Modular Arithmetic
aliases: [modular arithmetic, clock arithmetic, congruence, mod]
tags: [cryptography]
summary: Arithmetic that wraps around a fixed modulus — the finite-ring algebra that makes public-key cryptography possible.
weight: 40
---

# Modular Arithmetic

**Modular arithmetic** is arithmetic that *wraps around*. On a 12-hour clock, 4 hours after 10 o'clock is not 14 — it is **2**, because we wrap at 12. Written formally, $10 + 4 \equiv 2 \pmod{12}$. We say two numbers are **congruent mod $n$** when they leave the same remainder on division by $n$:

$$a \equiv b \pmod{n} \quad\Longleftrightarrow\quad n \text{ divides } (a - b).$$

Fix a modulus $n$ and the integers collapse into just $n$ residues, $\{0, 1, \dots, n-1\}$ — a finite ring. Addition, subtraction, and multiplication all still work; results simply fold back into that range. This little finite world is the stage on which [[RSA]] and [[Diffie–Hellman]] perform.

## Powers that cycle

The operation that powers public-key cryptography is **modular exponentiation**: computing $a^k \bmod n$. Because there are only finitely many residues, the powers $a^1, a^2, a^3, \dots$ can never run off to infinity — they must eventually **repeat**, tracing a cycle through the ring. Watch one form.

{{< sketch height="400" caption="Points 0…n−1 sit around a ring (n = 17). Drag left↔right to choose the base a; the orange hop walks a, a², a³, … mod n. The trail closes into a loop — the order of a. Some bases visit every nonzero point; others get stuck in a short cycle." >}}
if (frame === 0) {
  state.n = 17;
  state.a = 3;
  state.val = 1;     // current a^k mod n, starts at a^0 = 1
  state.k = 0;
  state.edges = [];  // recent hops {from,to,age}
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

const n = state.n;
// base from mouse (2..n-1)
if (mouse.x >= 0 && mouse.x <= W) {
  const na = 2 + Math.floor(mouse.x / W * (n - 2));
  const clamped = Math.max(2, Math.min(n - 1, na));
  if (clamped !== state.a) { state.a = clamped; state.val = 1; state.k = 0; state.edges = []; }
}

// advance the exponent on a timer
state.acc += dt;
if (state.acc > 0.55) {
  state.acc = 0;
  const next = (state.val * state.a) % n;
  state.edges.push({ from: state.val, to: next, age: 0 });
  state.val = next;
  state.k++;
  if (state.edges.length > 40) state.edges.shift();
}
for (const e of state.edges) e.age += dt;

// ring geometry
const cx = W * 0.40, cy = H * 0.52, R = Math.min(W * 0.32, H * 0.40);
const pos = (i) => {
  const ang = -Math.PI / 2 + i / n * Math.PI * 2;
  return { x: cx + Math.cos(ang) * R, y: cy + Math.sin(ang) * R };
};
// edges (hops)
ctx.lineCap = 'round';
for (const e of state.edges) {
  const p = pos(e.from), q = pos(e.to);
  const a = Math.max(0.08, 1 - e.age / 12);
  ctx.strokeStyle = accent2; ctx.globalAlpha = a; ctx.lineWidth = 1.6;
  ctx.beginPath(); ctx.moveTo(p.x, p.y); ctx.lineTo(q.x, q.y); ctx.stroke();
}
ctx.globalAlpha = 1;
// nodes
ctx.textAlign = 'center'; ctx.textBaseline = 'middle'; ctx.font = '11px sans-serif';
for (let i = 0; i < n; i++) {
  const p = pos(i);
  const here = i === state.val;
  ctx.fillStyle = here ? warn : 'rgba(127,127,127,0.25)';
  ctx.beginPath(); ctx.arc(p.x, p.y, here ? 9 : 5, 0, 7); ctx.fill();
  ctx.fillStyle = here ? '#000' : faint;
  ctx.fillText(i, p.x, p.y);
}
// readout panel (right)
const rx = W * 0.74;
ctx.textAlign = 'left'; ctx.fillStyle = text; ctx.font = '13px monospace';
ctx.fillText('n = ' + n, rx, H * 0.30);
ctx.fillStyle = accent;
ctx.fillText('a = ' + state.a, rx, H * 0.38);
ctx.fillStyle = warn;
ctx.fillText('a^' + state.k + ' mod n', rx, H * 0.50);
ctx.fillText('   = ' + state.val, rx, H * 0.58);
ctx.fillStyle = faint; ctx.font = '11px sans-serif';
ctx.fillText('drag to change a', rx, H * 0.70);
{{< /sketch >}}

When $n$ is prime, some bases are **generators**: their powers visit *every* nonzero residue before returning home. That property is exactly what [[Diffie–Hellman]] needs.

## The two trapdoors crypto leans on

Modular arithmetic matters to cryptography because it hides operations that are **easy one way, brutally hard to reverse** inside its finite ring.

- **Modular exponentiation is fast.** Even for thousand-digit numbers, $a^k \bmod n$ takes only a few dozen squarings (repeatedly square and multiply). Going *forward* is cheap.
- **The inverse is hard.** Recovering $k$ from $a^k \bmod n$ — the **discrete logarithm** — has no known fast method for well-chosen $n$. That asymmetry is the trapdoor under [[Diffie–Hellman]].
- **Factoring is hard.** Multiplying two large primes $p \cdot q = n$ is instant; splitting $n$ back into $p$ and $q$ is believed intractable. That is the trapdoor under [[RSA]].

{{< note kind="tip" title="Fermat's little theorem" >}}
For a prime $p$ and any $a$ not divisible by $p$, $a^{p-1} \equiv 1 \pmod{p}$. This is the gear that makes [[RSA]] decryption undo encryption exactly — choosing the private exponent $d$ as the modular inverse of $e$ guarantees $m^{ed} \equiv m$. The whole scheme is Fermat's theorem wearing a key.
{{< /note >}}

{{< quiz question="On a 12-hour clock, what is 7 × 5 mod 12?" options="35|11|1|5" answer="2" explain="7 × 5 = 35, and 35 = 2·12 + 11, so 35 ≡ 11 (mod 12). Multiplication wraps around just like addition." >}}

## See also

- [[RSA]]
- [[Diffie–Hellman]]
- [[Caesar Cipher]]
