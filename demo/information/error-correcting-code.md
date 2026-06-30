---
title: Error-Correcting Code
aliases: [error correction, ecc, hamming code, parity]
tags: [information]
summary: Codes that add structured redundancy so corrupted bits can be detected — and even located and repaired.
weight: 70
---

# Error-Correcting Code

An **error-correcting code** adds carefully structured redundancy to a message so that, when noise flips some of the bits in transit, the receiver can not only *notice* the damage but *repair* it — without ever asking for a retransmission. This is what lets a scratched DVD still play, a deep-space probe whisper across billions of kilometers, and a memory chip shrug off a cosmic ray.

The trick is to use only a sparse set of **valid codewords**, chosen so that every two of them differ in many bit positions (a large [[Hamming Distance]]). A few flipped bits then leave you nearer to the *original* codeword than to any other, so the decoder snaps the corrupted word back to the closest legal one.

## Parity, the smallest example

Add one **parity bit** equal to the XOR of all the data bits, so every valid codeword has an even number of 1s. Flip any single bit and the parity becomes odd — instant detection. But a single parity bit only *detects*; it cannot say *which* bit flipped, so it cannot correct. To locate the error you need several overlapping parity checks.

## Hamming(7,4): locate and fix one flip

Richard Hamming's 1950 code packs **4 data bits** into a **7-bit** codeword using **3 parity bits**, placed at positions 1, 2, and 4. Each parity bit checks a different overlapping group of positions:

- $p_1$ (pos 1) checks positions **1, 3, 5, 7**
- $p_2$ (pos 2) checks positions **2, 3, 6, 7**
- $p_4$ (pos 4) checks positions **4, 5, 6, 7**

The genius is the overlap. When a single bit flips, the *pattern* of which checks fail spells out, in binary, the **exact position** of the error. Compute the three failing/passing checks as a 3-bit number — the **syndrome** — and it literally equals the index of the broken bit. Read it, flip that bit back, done.

**Click any of the seven bits below to corrupt it.** Watch the parity checks light up red and the syndrome pinpoint the flip — then read the corrected codeword the decoder recovers.

{{< sketch height="420" caption="An interactive Hamming(7,4) codeword. Click a bit to flip it; the three parity checks (covering overlapping position groups) fail in a pattern whose binary value is the error's position. Click the same bit again, or any bit, to keep experimenting." >}}
if (frame === 0) {
  // a valid Hamming(7,4) codeword, indices 0..6 = positions 1..7
  state.b = [0,1,1,0,0,1,1];
  state.prev = false;
}
ctx.clearRect(0, 0, W, H);
const css = (v, f) => (getComputedStyle(document.documentElement).getPropertyValue(v).trim() || f);
const accent = css('--accent', '#5b9cff');
const accent2 = css('--accent-2', '#b07bff');
const good = css('--good', '#4ec98f');
const warn = css('--warn', '#e0a458');
const text = css('--text', '#e6e9ef');
const faint = css('--text-faint', 'rgba(230,233,239,0.55)');
const border = css('--border', 'rgba(255,255,255,0.25)');

// geometry for the 7 boxes
const n = 7;
const bw = Math.min(54, W / (n + 2));
const gap = bw * 0.32;
const totalW = n * bw + (n - 1) * gap;
const x0 = (W - totalW) / 2;
const boxY = 54, boxH = bw;
const boxX = (i) => x0 + i * (bw + gap);
const parityPos = {0:true, 1:true, 3:true}; // indices that are parity bits

// hit test + flip on fresh click
if (mouse.down && !state.prev) {
  for (let i = 0; i < n; i++) {
    if (mouse.x >= boxX(i) && mouse.x <= boxX(i) + bw && mouse.y >= boxY && mouse.y <= boxY + boxH) {
      state.b[i] ^= 1;
      break;
    }
  }
}
state.prev = mouse.down;

const b = state.b;
// parity checks (XOR over each covering group), indices 0-based
const c1 = b[0]^b[2]^b[4]^b[6];
const c2 = b[1]^b[2]^b[5]^b[6];
const c4 = b[3]^b[4]^b[5]^b[6];
const syndrome = c1 * 1 + c2 * 2 + c4 * 4; // 1-based position of error, 0 = none
const errIdx = syndrome - 1;

// corrected codeword
const corrected = b.slice();
if (syndrome >= 1 && syndrome <= 7) corrected[errIdx] ^= 1;

// --- header ---
ctx.fillStyle = faint;
ctx.font = '12px sans-serif';
ctx.fillText('received codeword — click a bit to flip it', x0, 32);

// --- draw boxes ---
ctx.textAlign = 'center';
ctx.textBaseline = 'middle';
for (let i = 0; i < n; i++) {
  const x = boxX(i);
  const isErr = (i === errIdx);
  ctx.beginPath();
  ctx.rect(x, boxY, bw, boxH);
  ctx.fillStyle = parityPos[i] ? 'rgba(176,123,255,0.16)' : 'rgba(91,156,255,0.16)';
  ctx.fill();
  ctx.lineWidth = isErr ? 3 : 1.4;
  ctx.strokeStyle = isErr ? warn : border;
  ctx.stroke();
  ctx.fillStyle = b[i] ? text : faint;
  ctx.font = 'bold 20px monospace';
  ctx.fillText(b[i] + '', x + bw / 2, boxY + boxH / 2 + 1);
  // position label
  ctx.fillStyle = faint;
  ctx.font = '10px sans-serif';
  ctx.fillText('pos ' + (i + 1), x + bw / 2, boxY - 8);
  ctx.fillText(parityPos[i] ? 'p' : 'd', x + bw / 2, boxY + boxH + 12);
}

// --- parity check rows ---
ctx.textAlign = 'left';
const groups = [
  {name: 'p₁ checks 1,3,5,7', cover: [0,2,4,6], c: c1},
  {name: 'p₂ checks 2,3,6,7', cover: [1,2,5,6], c: c2},
  {name: 'p₄ checks 4,5,6,7', cover: [3,4,5,6], c: c4},
];
let gy = boxY + boxH + 40;
ctx.font = '13px sans-serif';
for (const g of groups) {
  const ok = g.c === 0;
  ctx.fillStyle = ok ? good : warn;
  ctx.beginPath(); ctx.arc(x0 + 7, gy - 4, 6, 0, 2 * Math.PI); ctx.fill();
  ctx.fillStyle = text;
  ctx.fillText(g.name, x0 + 22, gy);
  ctx.fillStyle = ok ? good : warn;
  ctx.fillText(ok ? 'pass' : 'FAIL', x0 + 200, gy);
  gy += 26;
}

// --- syndrome + status ---
ctx.font = '13px monospace';
ctx.fillStyle = faint;
ctx.fillText('syndrome (c₄c₂c₁) = ' + c4 + '' + c2 + '' + c1 + ' = ' + syndrome, x0, gy + 4);
gy += 28;
ctx.font = '14px sans-serif';
if (syndrome === 0) {
  ctx.fillStyle = good;
  ctx.fillText('✓ All checks pass — no error detected.', x0, gy);
} else {
  ctx.fillStyle = warn;
  ctx.fillText('Error located at position ' + syndrome + ' → decoder flips it back.', x0, gy);
  gy += 24;
  ctx.fillStyle = good;
  ctx.font = '13px monospace';
  ctx.fillText('corrected: ' + corrected.join(' '), x0, gy);
}
ctx.textAlign = 'left';
ctx.textBaseline = 'alphabetic';
{{< /sketch >}}

{{< note kind="warning" title="One error, not two" >}}
Hamming(7,4) has minimum [[Hamming Distance]] 3, so it can **correct any single** bit-flip or **detect any double** flip — but not both at once, and it cannot correct two. Flip two bits in the sketch and the syndrome will confidently point at the *wrong* single position. Correcting power grows with code distance: distance $d$ corrects up to $\lfloor (d-1)/2 \rfloor$ errors.
{{< /note >}}

## The cost and the payoff

Redundancy is not free — Hamming(7,4) spends 7 bits to carry 4, a rate of $4/7$. But the [[Channel Capacity]] theorem promises that as long as you stay below capacity, codes exist that push the error rate as low as you like while keeping the rate high. Modern codes (Reed–Solomon, LDPC, turbo, polar) are the descendants of Hamming's insight, and they operate startlingly close to Shannon's limit.

{{< quiz question="In Hamming(7,4), what does the 3-bit syndrome tell you?" options="How many bits were sent|The binary position of the single flipped bit|The original data, directly|Whether the channel is noisy" answer="2" explain="The pattern of failing parity checks, read as a binary number, equals the index of the erroneous bit — so the decoder knows exactly which bit to flip back." >}}

## See also

- [[Hamming Distance]]
- [[Channel Capacity]]
- [[Data Compression]]
