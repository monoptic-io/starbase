---
title: Public-Key Cryptography
aliases: [public key cryptography, asymmetric cryptography, trapdoor function]
tags: [cryptography]
summary: Split the key in two — a public lock anyone can close and a private key only you can open — using one-way functions with a hidden trapdoor.
weight: 70
---

# Public-Key Cryptography

For all of history, encryption meant a **single shared secret**: the same key locked and unlocked, so the two parties had to meet or trust a courier first. **Public-key cryptography** — also called **asymmetric cryptography** — shattered that assumption. Each person holds a **key pair**:

- a **public key**, published to the world, that anyone can use to *encrypt to you* or *verify your signature*;
- a **private key**, kept secret, that only you can use to *decrypt* or *sign*.

The metaphor is an open padlock. You hand out copies of an unlocked padlock (your public key); anyone can snap one shut around a message and send it back. Only you hold the physical key that opens it. Closing a padlock takes no secret; opening it does.

{{< note kind="key" title="The breakthrough" >}}
Public-key cryptography means **two strangers can communicate secretly with no prior shared secret**, over a channel an adversary fully observes. That single idea — published by Diffie and Hellman in 1976, then realized by [[RSA]] in 1977 — underpins essentially all secure communication today: HTTPS, messaging apps, software updates, cryptocurrencies.
{{< /note >}}

## Trapdoor one-way functions

The whole edifice stands on a special kind of function: a **trapdoor one-way function**.

- **One-way** — easy to compute forward, infeasible to invert. Like smashing a plate: trivial to do, hopeless to reverse from the pieces.
- **Trapdoor** — *unless* you hold a secret, in which case inversion becomes easy again.

The public key describes the one-way function; the private key *is* the trapdoor. Two concrete trapdoors power the real world, both living in [[Modular Arithmetic]]:

- **Factoring** — multiplying primes $p \cdot q$ is easy; factoring $n$ back apart is hard. This is [[RSA]]'s trapdoor.
- **Discrete logarithm** — computing $g^a \bmod p$ is easy; recovering $a$ is hard. This powers [[Diffie–Hellman]] and elliptic-curve schemes.

{{< sketch height="300" caption="The same operation, two directions. Forward (locking) is cheap for everyone — the bar fills instantly. Backward (unlocking) is a cliff for an attacker without the trapdoor, but instant for the private-key holder. Move the mouse to scale the key size and watch the gap explode." >}}
if (frame === 0) { state.size = 0.5; }
ctx.clearRect(0, 0, W, H);
const css = (v, f) => (getComputedStyle(document.documentElement).getPropertyValue(v).trim() || f);
const accent = css('--accent', '#5b9cff');
const good = css('--good', '#4ec98f');
const warn = css('--warn', '#e0a458');
const text = css('--text', '#e6e9ef');
const faint = css('--text-faint', 'rgba(230,233,239,0.5)');

if (mouse.x >= 0 && mouse.x <= W) state.size = Math.max(0.05, Math.min(1, mouse.x / W));
const s = state.size;

const barX = W * 0.30, barW = W * 0.55, barH = 26;
const forward = 0.04 + 0.12 * s;          // forward effort: small, linear
const backward = Math.min(1, Math.pow(2, s * 9) / 360); // reverse: blows up

const drawBar = (y, frac, col, label) => {
  ctx.fillStyle = 'rgba(127,127,127,0.18)';
  ctx.fillRect(barX, y, barW, barH);
  ctx.fillStyle = col;
  ctx.fillRect(barX, y, barW * frac, barH);
  ctx.fillStyle = text; ctx.textAlign = 'right'; ctx.font = '12px sans-serif';
  ctx.fillText(label, barX - 12, y + barH / 2 + 4);
};
ctx.textBaseline = 'middle';
drawBar(H * 0.30, forward, good, 'lock (forward)');
drawBar(H * 0.55, backward, warn, 'break (reverse)');

ctx.textAlign = 'left'; ctx.fillStyle = faint; ctx.font = '12px sans-serif';
ctx.fillText('with the private key, unlocking is instant too — that is the trapdoor', barX, H * 0.80);
ctx.fillText('key size →  ' + Math.round(s * 100) + '%', barX, H * 0.12);
{{< /sketch >}}

## What the two keys buy you

Asymmetry gives two distinct powers, depending on which key acts first:

- **Encrypt with the public key → decrypt with the private key.** Anyone can send *you* a secret; only you can read it. This is confidentiality.
- **Sign with the private key → verify with the public key.** Only *you* can produce the value; anyone can check it. This is a [[Digital Signature]], giving authenticity and integrity.

In practice asymmetric crypto is slow, so it rarely encrypts bulk data directly. Instead it does the *hard part* — agreeing on or delivering a short symmetric key (often via [[Diffie–Hellman]]) — and a fast symmetric cipher handles the rest. Public-key crypto is the handshake; symmetric crypto is the conversation.

## See also

- [[RSA]]
- [[Diffie–Hellman]]
- [[Digital Signature]]
