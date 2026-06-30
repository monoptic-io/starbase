---
title: Caesar Cipher
aliases: [caesar cipher, shift cipher]
tags: [cryptography]
summary: Shift every letter of the alphabet by a fixed amount — the oldest cipher, and one of the easiest to break.
weight: 10
---

# Caesar Cipher

The **Caesar cipher** encrypts a message by sliding every letter a fixed number of places down the alphabet. With a shift of $k = 3$, $A$ becomes $D$, $B$ becomes $E$, and $Z$ wraps around to $C$. Julius Caesar reportedly used exactly this to shield his military dispatches.

Numbering the letters $A=0, B=1, \dots, Z=25$, encryption and decryption are one line of [[Modular Arithmetic]]:

$$c = (m + k) \bmod 26, \qquad m = (c - k) \bmod 26.$$

The wrap-around — the $\bmod 26$ — is the whole trick. The alphabet is a ring of 26 positions, and the key just rotates it.

## A cipher wheel you can turn

The classic tool for a shift cipher is a **cipher wheel**: two concentric alphabets, the inner one rotated against the outer by the key. Line up the rings and you can read off the substitution at a glance.

{{< sketch height="400" caption="Drag left↔right to set the shift k. The wheel rotates, the sample message re-encrypts live, and the frequency bars show how the telltale English peaks just slide along — handing the key to any code-breaker." >}}
if (frame === 0) {
  state.shift = 3;
  state.sample = "THEQUICKBROWNFOX";
  // approximate English letter frequencies A..Z (percent)
  state.eng = [8.2,1.5,2.8,4.3,12.7,2.2,2.0,6.1,7.0,0.15,0.77,4.0,2.4,6.7,7.5,1.9,0.095,6.0,6.3,9.1,2.8,0.98,2.4,0.15,2.0,0.074];
}
ctx.clearRect(0, 0, W, H);
const css = (v, f) => (getComputedStyle(document.documentElement).getPropertyValue(v).trim() || f);
const accent = css('--accent', '#5b9cff');
const accent2 = css('--accent-2', '#b07bff');
const good = css('--good', '#4ec98f');
const warn = css('--warn', '#e0a458');
const text = css('--text', '#e6e9ef');
const faint = css('--text-faint', 'rgba(230,233,239,0.5)');

// mouse sets the shift 0..25
if (mouse.x >= 0 && mouse.x <= W) state.shift = Math.max(0, Math.min(25, Math.floor(mouse.x / W * 26)));
const k = state.shift;
const A = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';

// ---- cipher wheel (left) ----
const cx = W * 0.26, cy = H * 0.46, R1 = Math.min(W * 0.22, H * 0.40), R2 = R1 * 0.74;
const step = Math.PI * 2 / 26;
ctx.lineWidth = 1;
ctx.strokeStyle = faint;
ctx.beginPath(); ctx.arc(cx, cy, R1 + 14, 0, 7); ctx.stroke();
ctx.beginPath(); ctx.arc(cx, cy, R2 - 14, 0, 7); ctx.stroke();
ctx.textAlign = 'center'; ctx.textBaseline = 'middle';
for (let j = 0; j < 26; j++) {
  const ang = -Math.PI / 2 + j * step;
  // outer ring = ciphertext alphabet (fixed)
  const ox = cx + Math.cos(ang) * R1, oy = cy + Math.sin(ang) * R1;
  ctx.font = '12px sans-serif';
  ctx.fillStyle = text;
  ctx.fillText(A[j], ox, oy);
  // inner ring = plaintext, rotated so plaintext p sits under cipher (p+k)
  const pIdx = (j - k + 26) % 26;
  const ix = cx + Math.cos(ang) * R2, iy = cy + Math.sin(ang) * R2;
  ctx.fillStyle = accent;
  ctx.fillText(A[pIdx], ix, iy);
}
// highlight the A -> (A+k) spoke
const ang0 = -Math.PI / 2 + ((0 + k) % 26) * step;
ctx.strokeStyle = warn; ctx.lineWidth = 2;
ctx.beginPath();
ctx.moveTo(cx + Math.cos(ang0) * (R2 - 10), cy + Math.sin(ang0) * (R2 - 10));
ctx.lineTo(cx + Math.cos(ang0) * (R1 + 10), cy + Math.sin(ang0) * (R1 + 10));
ctx.stroke();
ctx.fillStyle = warn; ctx.font = 'bold 13px sans-serif';
ctx.fillText('k = ' + k, cx, cy);

// ---- sample encryption (top right) ----
const rx = W * 0.52;
ctx.textAlign = 'left';
ctx.font = '13px monospace';
let line = '';
for (const ch of state.sample) line += ch;
let cipher = '';
for (const ch of state.sample) {
  const idx = ch.charCodeAt(0) - 65;
  cipher += A[(idx + k) % 26];
}
ctx.fillStyle = faint; ctx.fillText('plain  ', rx, 30);
ctx.fillStyle = text;  ctx.fillText(line, rx + 52, 30);
ctx.fillStyle = faint; ctx.fillText('cipher ', rx, 50);
ctx.fillStyle = accent2; ctx.fillText(cipher, rx + 52, 50);

// ---- shifted frequency histogram (bottom right) ----
const hx = rx, hy = H - 30, hw = W - rx - 24, hbh = H * 0.5;
ctx.strokeStyle = faint; ctx.lineWidth = 1;
ctx.beginPath(); ctx.moveTo(hx, hy); ctx.lineTo(hx + hw, hy); ctx.stroke();
ctx.font = '10px sans-serif'; ctx.textAlign = 'center';
const bw = hw / 26;
let peak = 0; for (let i = 0; i < 26; i++) peak = Math.max(peak, state.eng[i]);
for (let i = 0; i < 26; i++) {
  // ciphertext letter i carries the frequency of plaintext (i - k)
  const f = state.eng[(i - k + 26) % 26];
  const bh = f / peak * hbh;
  const isE = ((i - k + 26) % 26) === 4; // the English 'E' peak, now shifted
  ctx.fillStyle = isE ? warn : accent;
  ctx.globalAlpha = isE ? 1 : 0.7;
  ctx.fillRect(hx + i * bw + 1, hy - bh, bw - 2, bh);
  ctx.globalAlpha = 1;
  ctx.fillStyle = faint;
  ctx.fillText(A[i], hx + i * bw + bw / 2, hy + 9);
}
ctx.fillStyle = warn; ctx.textAlign = 'left'; ctx.font = '11px sans-serif';
ctx.fillText("the 'E' spike points straight at k", hx, hy - hbh - 6);
{{< /sketch >}}

Watch the orange bar — it marks where English's dominant letter **E** has landed. Its position *is* the key. That is the fatal flaw.

## Why it falls in seconds

A Caesar cipher has exactly **25 useful keys** (a shift of 0 does nothing). An attacker simply tries all of them — a *brute-force* search so small you can do it by hand. Even faster: because the cipher only *slides* the alphabet, it leaves the shape of English untouched. The most common ciphertext letter is almost always a shifted **E**, and that one observation usually reveals $k$ outright.

In the language of [[Entropy]], the key carries only $\log_2 25 \approx 4.6$ bits of uncertainty — a rounding error against the information in a real message. There simply isn't enough secret to hide behind. The natural next step is to scramble the alphabet arbitrarily instead of merely rotating it: the [[Substitution Cipher]].

{{< quiz question="Ciphertext encrypted with a Caesar cipher has 'W' as its most common letter. What is the most likely shift k?" options="k = 5|k = 18|k = 23|impossible to tell" answer="2" explain="The commonest English letter is E (index 4). If E maps to W (index 22), then k = 22 − 4 = 18. The frequency peak hands you the key." >}}

## See also

- [[Substitution Cipher]]
- [[Modular Arithmetic]]
- [[Entropy]]
