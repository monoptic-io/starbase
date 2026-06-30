---
title: Hash Function
aliases: [hash function, cryptographic hash, sha, message digest]
tags: [cryptography]
summary: A one-way fingerprint that turns any input into a fixed-size digest, where flipping one input bit scrambles roughly half the output.
weight: 80
---

# Hash Function

A **cryptographic hash function** crushes an input of *any* size — a word, a file, a movie — into a short, fixed-length **digest** (a "fingerprint"), in a way that is easy to compute but impossible to reverse or fool. SHA-256, for instance, maps anything to exactly 256 bits. It is the workhorse of modern cryptography: it backs [[Digital Signature]]s, password storage, blockchains, and integrity checks everywhere.

A good hash $h = H(m)$ must satisfy three properties:

- **Pre-image resistance** — given a digest $h$, you cannot find any $m$ with $H(m) = h$ (one-way).
- **Second pre-image resistance** — given $m$, you cannot find a *different* $m'$ with the same digest.
- **Collision resistance** — you cannot find *any* two inputs that hash to the same value.

## The avalanche effect

The property that makes a fingerprint trustworthy is the **avalanche effect**: change the input by a single bit and about **half** of the output bits flip, with no visible pattern. The digest of `hello` and the digest of `hellp` should look completely unrelated. Measured with [[Hamming Distance]] — the number of differing bits — a one-bit input nudge should land you roughly halfway across the entire output space.

{{< sketch height="380" caption="Two inputs differing by a single character, run through a small illustrative hash. The output bit-grids share almost nothing — click to re-roll the inputs and watch the avalanche hold near 50% every time. (This toy hash is for demonstration, not real security.)" >}}
if (frame === 0) { state.nonce = 1734; }
ctx.clearRect(0, 0, W, H);
const css = (v, f) => (getComputedStyle(document.documentElement).getPropertyValue(v).trim() || f);
const accent = css('--accent', '#5b9cff');
const accent2 = css('--accent-2', '#b07bff');
const good = css('--good', '#4ec98f');
const warn = css('--warn', '#e0a458');
const text = css('--text', '#e6e9ef');
const faint = css('--text-faint', 'rgba(230,233,239,0.5)');

if (mouse.clicked) state.nonce = (state.nonce * 1103515245 + 12345) >>> 8;

// a small illustrative avalanche hash -> 32 bits (NOT cryptographically secure)
const hash32 = (str) => {
  let h = 2166136261 >>> 0;
  for (let i = 0; i < str.length; i++) {
    h ^= str.charCodeAt(i);
    h = Math.imul(h, 16777619) >>> 0;
  }
  h ^= h >>> 16; h = Math.imul(h, 2246822507) >>> 0;
  h ^= h >>> 13; h = Math.imul(h, 3266489909) >>> 0;
  h ^= h >>> 16;
  return h >>> 0;
};
const inA = 'msg-' + state.nonce;
const inB = 'msg-' + (state.nonce + 1); // differs by one character
const hA = hash32(inA), hB = hash32(inB);

// bit extraction
const bits = (x) => { const a = []; for (let i = 31; i >= 0; i--) a.push((x >>> i) & 1); return a; };
const bA = bits(hA), bB = bits(hB);
let diff = 0; for (let i = 0; i < 32; i++) if (bA[i] !== bB[i]) diff++;

ctx.textAlign = 'left'; ctx.textBaseline = 'middle';
ctx.font = '13px monospace';
ctx.fillStyle = faint; ctx.fillText('input A:', 20, H * 0.10);
ctx.fillStyle = accent; ctx.fillText('"' + inA + '"', 90, H * 0.10);
ctx.fillStyle = faint; ctx.fillText('input B:', 20, H * 0.18);
ctx.fillStyle = accent2; ctx.fillText('"' + inB + '"', 90, H * 0.18);
ctx.fillStyle = warn; ctx.font = '11px sans-serif';
ctx.fillText('(one character apart)', 230, H * 0.18);

// draw two bit grids (4 rows x 8 cols)
const cols = 8, rowsN = 4;
const cell = Math.min((W - 60) / cols, 26);
const gx = (W - cols * cell) / 2;
const drawGrid = (bitsArr, otherArr, gy, label, col) => {
  ctx.fillStyle = faint; ctx.font = '12px sans-serif'; ctx.textAlign = 'left';
  ctx.fillText(label, gx, gy - 12);
  for (let i = 0; i < 32; i++) {
    const r = Math.floor(i / cols), c = i % cols;
    const x = gx + c * cell, y = gy + r * cell;
    const on = bitsArr[i] === 1;
    const differs = bitsArr[i] !== otherArr[i];
    ctx.fillStyle = on ? col : 'rgba(127,127,127,0.15)';
    ctx.fillRect(x + 1, y + 1, cell - 2, cell - 2);
    if (differs) { ctx.strokeStyle = warn; ctx.lineWidth = 2; ctx.strokeRect(x + 1, y + 1, cell - 2, cell - 2); }
  }
};
const g1 = H * 0.34, g2 = H * 0.34 + rowsN * cell + 36;
drawGrid(bA, bB, g1, 'hash(A) — 32 bits', accent);
drawGrid(bB, bA, g2, 'hash(B) — 32 bits', accent2);

// avalanche readout
ctx.textAlign = 'center'; ctx.fillStyle = good; ctx.font = 'bold 13px sans-serif';
ctx.fillText(diff + ' of 32 bits differ  (' + Math.round(diff / 32 * 100) + '%)  —  orange boxes = flipped', W / 2, H - 14);
{{< /sketch >}}

Notice the digest stays near 50% flipped no matter which inputs you roll. That near-total decorrelation, viewed through [[Entropy]], is what makes a digest behave like a random draw and lets it stand in for the message itself.

## Where hashes earn their keep

The avalanche and one-wayness combine into a few indispensable jobs:

- **Integrity.** Publish a file's digest; a single altered byte avalanches the hash, exposing tampering instantly.
- **Signatures.** Public-key signing is slow, so you sign the *short* hash of a message rather than the whole thing — the heart of a [[Digital Signature]].
- **Passwords.** Store $H(\text{password})$, never the password. A breach leaks digests, and one-wayness keeps the originals hidden.
- **Commitments and blockchains.** Hash-link records so any change to history breaks every fingerprint downstream.

{{< note kind="warning" title="A demo hash, not a real one" >}}
The 32-bit toy above shows the *flavor* of avalanche, but real hashes like SHA-256 produce 256 bits and resist attacks the toy does not. Never invent your own hash for production — small or homemade hashes have collisions an adversary can find. The principle is universal; the security is in the vetted construction.
{{< /note >}}

{{< quiz question="A good cryptographic hash has the avalanche property. What does flipping a single input bit do to the output?" options="Flips exactly one output bit|Changes nothing — the hash is stable|Flips about half the output bits, unpredictably|Makes the output one bit longer" answer="3" explain="Avalanche means a one-bit input change cascades to roughly 50% of output bits with no discernible pattern, so similar inputs yield unrelated digests." >}}

## See also

- [[Hamming Distance]]
- [[Digital Signature]]
- [[Entropy]]
