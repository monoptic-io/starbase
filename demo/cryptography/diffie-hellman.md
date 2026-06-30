---
title: Diffie–Hellman
aliases: [diffie hellman, dh, key exchange]
tags: [cryptography]
summary: Two strangers agree on a shared secret over a fully public channel — the paint-mixing trick that launched public-key cryptography.
weight: 50
---

# Diffie–Hellman

**Diffie–Hellman key exchange** solves a problem that sounds impossible: two people who have never met, shouting across a room full of eavesdroppers, end up holding the *same secret number* — and no listener can figure out what it is. Published in 1976 by Whitfield Diffie and Martin Hellman, it cracked open the field of [[Public-Key Cryptography]].

It does not encrypt a message. It does something more fundamental: it lets Alice and Bob *agree on a shared key* over an open wire, which they can then feed into a fast symmetric cipher or a [[One-Time Pad]].

## The paint-mixing intuition

Forget numbers for a moment. Imagine mixing paint:

1. Alice and Bob publicly agree on a **common color** — say yellow. Everyone, including the eavesdropper, sees it.
2. Each secretly picks a **private color** and mixes it into the yellow. Alice makes one blend, Bob another. They **swap these public blends** in the open.
3. Each stirs their *own* private color into the *other's* blend.

Both now hold yellow + Alice's secret + Bob's secret — the **same** mixture. The eavesdropper saw the two public blends but cannot un-mix paint to recover the private colors. Separating mixed paint is the "hard problem."

{{< sketch height="400" caption="Alice and Bob each mix the shared base with a private tint, exchange the blends (the traveling dots), then mix in their own secret again. Both arrive at the identical shared color at center — the eavesdropper only ever saw the two public blends. Drag to change the private tints." >}}
if (frame === 0) {
  state.base = 55;   // hue of public base
  state.aSec = 200;  // Alice private hue
  state.bSec = 320;  // Bob private hue
}
ctx.clearRect(0, 0, W, H);
const css = (v, f) => (getComputedStyle(document.documentElement).getPropertyValue(v).trim() || f);
const accent = css('--accent', '#5b9cff');
const faint = css('--text-faint', 'rgba(230,233,239,0.5)');
const text = css('--text', '#e6e9ef');

// mouse: x -> Alice secret hue, y -> Bob secret hue
if (mouse.x >= 0 && mouse.x <= W) state.aSec = Math.floor(mouse.x / W * 360);
if (mouse.y >= 0 && mouse.y <= H) state.bSec = Math.floor(mouse.y / H * 360);

// hsl(hue) -> rgb
const h2rgb = (h) => {
  const s = 0.6, l = 0.55;
  const c = (1 - Math.abs(2 * l - 1)) * s, hp = h / 60, x = c * (1 - Math.abs(hp % 2 - 1));
  let r = 0, g = 0, b = 0;
  if (hp < 1) [r, g, b] = [c, x, 0]; else if (hp < 2) [r, g, b] = [x, c, 0];
  else if (hp < 3) [r, g, b] = [0, c, x]; else if (hp < 4) [r, g, b] = [0, x, c];
  else if (hp < 5) [r, g, b] = [x, 0, c]; else [r, g, b] = [c, 0, x];
  const m = l - c / 2;
  return [Math.round((r + m) * 255), Math.round((g + m) * 255), Math.round((b + m) * 255)];
};
const mix = (cols) => {
  let r = 0, g = 0, b = 0; for (const c of cols) { r += c[0]; g += c[1]; b += c[2]; }
  const k = cols.length; return [Math.round(r / k), Math.round(g / k), Math.round(b / k)];
};
const str = (c) => 'rgb(' + c[0] + ',' + c[1] + ',' + c[2] + ')';
const blob = (x, y, r, col, label) => {
  ctx.fillStyle = str(col);
  ctx.beginPath(); ctx.arc(x, y, r, 0, 7); ctx.fill();
  ctx.fillStyle = faint; ctx.textAlign = 'center'; ctx.font = '11px sans-serif';
  if (label) ctx.fillText(label, x, y + r + 14);
};

const base = h2rgb(state.base);
const aS = h2rgb(state.aSec), bS = h2rgb(state.bSec);
const aPub = mix([base, aS]);   // Alice's public blend
const bPub = mix([base, bS]);   // Bob's public blend
const shared = mix([base, aS, bS]); // = mix(bPub,aS) = mix(aPub,bS)

const yTop = H * 0.14, yMid = H * 0.5, yBot = H * 0.84;
const xA = W * 0.16, xB = W * 0.84, xC = W * 0.5;

// labels
ctx.fillStyle = text; ctx.font = 'bold 13px sans-serif'; ctx.textAlign = 'center';
ctx.fillText('ALICE', xA, 16); ctx.fillText('BOB', xB, 16);

// public base (top center, on the open wire)
blob(xC, yTop, 16, base, 'public base');

// private secrets
blob(xA, yTop, 14, aS, 'a (secret)');
blob(xB, yTop, 14, bS, 'b (secret)');

// public blends, with a traveling dot showing the exchange
blob(xA, yMid, 18, aPub, 'A = base+a (public)');
blob(xB, yMid, 18, bPub, 'B = base+b (public)');
const phase = (t % 2) / 2; // 0..1
const lerp = (p, q, u) => [p[0] + (q[0] - p[0]) * u, p[1] + (q[1] - p[1]) * u, p[2] + (q[2] - p[2]) * u];
let dpA = lerp([xA, yMid], [xB, yMid], phase);
let dpB = lerp([xB, yMid], [xA, yMid], phase);
ctx.fillStyle = str(aPub); ctx.beginPath(); ctx.arc(dpA[0], yMid - 26, 7, 0, 7); ctx.fill();
ctx.fillStyle = str(bPub); ctx.beginPath(); ctx.arc(dpB[0], yMid + 26, 7, 0, 7); ctx.fill();

// shared secret (bottom) — identical on both sides
blob(xA, yBot, 20, shared, '');
blob(xB, yBot, 20, shared, '');
ctx.fillStyle = text; ctx.font = 'bold 12px sans-serif';
ctx.fillText('same shared secret', xC, yBot + 4);
{{< /sketch >}}

## The math behind the paint

Replace paint with [[Modular Arithmetic]] and the trick becomes precise. Everyone agrees on a large prime $p$ and a base $g$.

{{< eq number="1" >}}
A = g^{a} \bmod p, \qquad B = g^{b} \bmod p
{{< /eq >}}

Alice keeps her secret exponent $a$ private and sends $A$; Bob keeps $b$ private and sends $B$. Now each raises what they received to their *own* secret:

$$\text{Alice computes } B^{a} = g^{ba} \bmod p, \qquad \text{Bob computes } A^{b} = g^{ab} \bmod p.$$

Since $g^{ba} = g^{ab}$, **they hold the same value** — the shared secret $s = g^{ab} \bmod p$. The exponents commute; that is the whole engine.

{{< note kind="key" title="Why the eavesdropper is stuck" >}}
The wire carries $p$, $g$, $A = g^a$, and $B = g^b$ — all public. To find the shared secret, an attacker must recover $a$ from $g^a \bmod p$ (or $b$ from $g^b$). That is the **discrete logarithm problem**, and for a well-chosen large prime no one knows how to do it efficiently. Mixing forward is cheap; un-mixing is the wall security rests on.
{{< /note >}}

## What it is and isn't

Diffie–Hellman gives you **confidentiality of a freshly agreed key**, but on its own it does *not* tell Alice who is really on the other end — a "man in the middle" could exchange keys with each side separately. Pinning down identity needs a [[Digital Signature]] or a trusted [[Public-Key Cryptography|public key]]. In practice DH establishes the session key, and signatures vouch for who you are talking to.

{{< quiz question="An eavesdropper records p, g, A = g^a mod p and B = g^b mod p. Why can't they compute the shared secret g^(ab) mod p?" options="The values are encrypted|They would need to solve the discrete logarithm to recover a or b|g^(ab) is never actually transmitted, so it doesn't exist|They can, Diffie–Hellman is insecure" answer="2" explain="Computing the shared secret from public values requires recovering a private exponent — the discrete logarithm problem — which is believed intractable for large p." >}}

## See also

- [[Modular Arithmetic]]
- [[Public-Key Cryptography]]
- [[RSA]]
