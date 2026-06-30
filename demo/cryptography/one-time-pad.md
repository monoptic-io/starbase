---
title: One-Time Pad
aliases: [one time pad, vernam cipher]
tags: [cryptography]
summary: XOR a message with a truly random key as long as itself — the only cipher proven unbreakable, and the only one almost nobody can use.
weight: 30
---

# One-Time Pad

The **one-time pad** is the unicorn of cryptography: a cipher *mathematically proven* to be unbreakable. The recipe is almost embarrassingly simple. Take a key that is **truly random**, **as long as the message**, and **never reused**, then combine it with the message bit by bit using XOR (exclusive-or, addition [[Modular Arithmetic|mod 2]]):

$$c_i = m_i \oplus k_i, \qquad m_i = c_i \oplus k_i.$$

XOR is its own inverse, so the receiver — who has the same pad — just XORs again to peel the key back off. Invented in 1882 and patented by Gilbert Vernam in 1919, it is also called the **Vernam cipher**.

## Why it is perfectly secret

Here is the magic. For *any* ciphertext, every possible plaintext of the same length is **equally likely**, because for each candidate plaintext there is exactly one key that would have produced that ciphertext. The ciphertext therefore reveals *nothing*: an attacker who intercepts `1011` cannot tell whether you sent `YES` or `NO` or random noise — all are consistent with some key.

{{< note kind="key" title="Shannon's perfect secrecy" >}}
In 1949 Claude Shannon proved this is **perfect secrecy**: $P(m \mid c) = P(m)$ — seeing the ciphertext does not change the probability of any message. He also proved its price: perfect secrecy *requires* the key to have at least as much [[Entropy]] as the message, $H(K) \ge H(M)$. The one-time pad meets that bound exactly, which is why no cipher can be more secure — and why none can be secure with a shorter key.
{{< /note >}}

## Watch the key drown the message

The key is a fresh fair coin for every bit — a binary [[Random Walk]] with no pattern to grab onto. XORing a message with it produces ciphertext that is *itself* indistinguishable from random. Below, the structured message bits vanish under a random key, and XORing the key back in restores them perfectly.

{{< sketch height="320" caption="Top: a structured message. Middle: a fresh random key (re-roll by clicking). Bottom: their XOR — statistically pure noise that gives the message no cover to leak from." >}}
if (frame === 0) {
  state.N = 48;
  state.msg = [];
  for (let i = 0; i < state.N; i++) state.msg.push((i % 6 < 3) ? 1 : 0); // a blocky pattern
  state.key = [];
  for (let i = 0; i < state.N; i++) state.key.push(Math.random() < 0.5 ? 1 : 0);
  state.roll = 0;
}
ctx.clearRect(0, 0, W, H);
const css = (v, f) => (getComputedStyle(document.documentElement).getPropertyValue(v).trim() || f);
const accent = css('--accent', '#5b9cff');
const accent2 = css('--accent-2', '#b07bff');
const good = css('--good', '#4ec98f');
const warn = css('--warn', '#e0a458');
const text = css('--text', '#e6e9ef');
const faint = css('--text-faint', 'rgba(230,233,239,0.5)');

if (mouse.clicked) {
  for (let i = 0; i < state.N; i++) state.key[i] = Math.random() < 0.5 ? 1 : 0;
  state.roll++;
}

const cols = state.N;
const cw = (W - 80) / cols;
const ch = Math.min(38, H * 0.16);
const x0 = 70;
const rows = [
  { y: H * 0.10, bits: state.msg, col: accent,  label: 'message' },
  { y: H * 0.42, bits: state.key, col: warn,    label: 'key (random)' },
  { y: H * 0.74, bits: null,      col: good,    label: 'cipher = m⊕k' }
];
ctx.textAlign = 'right'; ctx.textBaseline = 'middle'; ctx.font = '11px sans-serif';
for (const r of rows) {
  ctx.fillStyle = faint;
  ctx.fillText(r.label, x0 - 8, r.y + ch / 2);
  for (let i = 0; i < cols; i++) {
    let b;
    if (r.bits) b = r.bits[i];
    else b = state.msg[i] ^ state.key[i];
    ctx.fillStyle = b ? r.col : 'rgba(127,127,127,0.18)';
    ctx.fillRect(x0 + i * cw + 0.5, r.y, cw - 1, ch);
  }
}
ctx.textAlign = 'left'; ctx.fillStyle = faint; ctx.font = '11px sans-serif';
ctx.fillText('click to draw a fresh key — the cipher reshuffles into noise every time', x0, H - 8);
{{< /sketch >}}

## Why almost nobody uses it

If it is perfect, why is the world not built on it? Because the very thing that makes it secure makes it impractical: **the key must be as long as everything you will ever send**, truly random, shared in advance through some already-secure channel, and *never reused*. To send a gigabyte secretly you must first secretly share a gigabyte of key — which raises the same problem you started with.

{{< note kind="warning" title="Reuse is fatal" >}}
Use a pad twice and the secrecy collapses. If $c_1 = m_1 \oplus k$ and $c_2 = m_2 \oplus k$, then $c_1 \oplus c_2 = m_1 \oplus m_2$ — the key cancels and the two messages bleed into each other. This exact mistake has broken real systems (famously the Soviet VENONA traffic). "One-time" is not a suggestion.
{{< /note >}}

The one-time pad is reserved for the highest-stakes links — diplomatic and intelligence channels — where couriers can hand-deliver key material. For everyone else, the practical path is to stop sharing huge secret keys altogether and let mathematics carry the secret instead. That is the promise of [[Diffie–Hellman]] and [[Public-Key Cryptography]].

## See also

- [[Entropy]]
- [[Random Walk]]
- [[Modular Arithmetic]]
