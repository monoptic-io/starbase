---
title: Circle of Fifths
aliases: [circle of fifths]
tags: [music]
summary: The twelve keys arranged in a ring by ascending fifths, so that neighbors share almost all their notes and sharps or flats accumulate one step at a time.
weight: 70
---

# Circle of Fifths

The **circle of fifths** is the single most useful map in music. Arrange the twelve keys around a clock so that each step *clockwise* rises by a perfect fifth — C, G, D, A, E, … — and a remarkable order emerges: neighbors on the circle are the keys most closely related, sharing all but one of their notes, and the [[Key Signature|sharps and flats]] pile up one at a time as you travel around.

## Reading the ring

Start at C at the top, with no sharps and no flats. Each step **clockwise** (up a fifth) *adds one sharp*: G has one (F♯), D has two, A has three, and so on. Each step **counter-clockwise** (down a fifth, i.e. up a fourth) *adds one flat*: F has one (B♭), B♭ has two, E♭ has three. Walk far enough either way and the sharp and flat spellings meet at the bottom of the circle — the same pitches, written two ways.

Why fifths? Because the perfect fifth is the most consonant [[Interval]] after the octave, keys a fifth apart overlap heavily, so moving between them sounds smooth. That is exactly why the strongest moves in [[Harmony]] are also moves by a fifth.

## Spin the circle

**Click any key** to select it. The selection lights up along with its two **neighbors** — the keys a fifth above and below, its closest relatives. The readout names the key and counts its sharps or flats.

{{< sketch height="420" caption="The twelve major keys in fifths order. Click a key to highlight it and its two nearest neighbors (a fifth on either side); the panel reports its sharps or flats." >}}
if (frame === 0) {
  const cs = getComputedStyle(document.documentElement);
  state.accent    = cs.getPropertyValue('--accent').trim()     || '#5b9cff';
  state.accent2   = cs.getPropertyValue('--accent-2').trim()   || '#ff6b6b';
  state.good      = cs.getPropertyValue('--good').trim()       || '#46c98b';
  state.border    = cs.getPropertyValue('--border').trim()     || 'rgba(255,255,255,0.18)';
  state.textFaint = cs.getPropertyValue('--text-faint').trim() || 'rgba(255,255,255,0.5)';
  state.text      = cs.getPropertyValue('--text').trim()       || '#e8e8ef';
  state.sel = 0;
  state.wasDown = false;
}
const accent = state.accent, accent2 = state.accent2, good = state.good,
      border = state.border, textFaint = state.textFaint, text = state.text;
// clockwise from top: ascending fifths
const KEYS  = ['C','G','D','A','E','B','F#','Db','Ab','Eb','Bb','F'];
const MINOR = ['Am','Em','Bm','F#m','C#m','G#m','D#m','Bbm','Fm','Cm','Gm','Dm'];
// signature: positive = sharps, negative = flats
const SIG   = [0,1,2,3,4,5,6,-5,-4,-3,-2,-1];

ctx.clearRect(0,0,W,H);
const clicked = mouse.clicked || (mouse.down && !state.wasDown);
state.wasDown = mouse.down;

const cx = W*0.5, cy = H*0.52;
const rOut = Math.min(W, H)*0.40, rIn = rOut*0.66;
function pos(i, r){ const a = -Math.PI/2 + i*(Math.PI*2/12); return {x: cx + r*Math.cos(a), y: cy + r*Math.sin(a)}; }

// hit-test outer ring on click
if (clicked) {
  for (let i=0;i<12;i++){ const p = pos(i, rOut);
    if (Math.hypot(mouse.x-p.x, mouse.y-p.y) < rOut*0.16) {
      state.sel = i;
      if (window.sgTone) {                           // sound the selected key's I chord
        const pc = (i*7) % 12, base = 261.63;        // tonic pitch class (ascending fifths)
        sgChord([base*Math.pow(2, pc/12), base*Math.pow(2, (pc+4)/12), base*Math.pow(2, (pc+7)/12)], 0.9);
      }
      break;
    } }
}
const sel = state.sel;
const nbrs = [(sel+11)%12, (sel+1)%12];   // a fifth either side

// connecting ring
ctx.strokeStyle = border; ctx.lineWidth = 1;
ctx.beginPath(); ctx.arc(cx, cy, rOut, 0, 7); ctx.stroke();
ctx.beginPath(); ctx.arc(cx, cy, rIn, 0, 7); ctx.stroke();

ctx.textAlign = 'center'; ctx.textBaseline = 'middle';
for (let i=0;i<12;i++){
  const p = pos(i, rOut), rad = rOut*0.15;
  const isSel = (i===sel), isNbr = nbrs.includes(i);
  ctx.fillStyle = isSel ? accent : (isNbr ? good : 'rgba(255,255,255,0.05)');
  ctx.beginPath(); ctx.arc(p.x, p.y, rad, 0, 7); ctx.fill();
  ctx.strokeStyle = border; ctx.lineWidth = 1; ctx.stroke();
  ctx.fillStyle = (isSel||isNbr) ? '#0c0f16' : text;
  ctx.font = '600 16px sans-serif';
  ctx.fillText(KEYS[i], p.x, p.y);
  // relative minor on inner ring
  const q = pos(i, rIn);
  ctx.fillStyle = isSel ? accent2 : textFaint;
  ctx.font = '12px sans-serif';
  ctx.fillText(MINOR[i], q.x, q.y);
}

// center readout
const s = SIG[sel];
const sigtxt = s===0 ? 'no sharps or flats'
             : s>0 ? (s + (s===1?' sharp':' sharps'))
                   : (-s + (-s===1?' flat':' flats'));
ctx.fillStyle = text; ctx.font = '600 22px sans-serif';
ctx.fillText(KEYS[sel] + ' major', cx, cy-10);
ctx.fillStyle = textFaint; ctx.font = '13px sans-serif';
ctx.fillText(sigtxt, cx, cy+14);
ctx.fillText('rel. minor: ' + MINOR[sel], cx, cy+32);
{{< /sketch >}}

{{< note kind="tip" title="A composer's compass" >}}
Want a chord that sounds related but fresh? Reach for a neighbor on the circle. Want a dramatic, distant shift? Jump to the far side. The circle turns "which keys go together" from guesswork into geography.
{{< /note >}}

## See also

- [[Key Signature]]
- [[Harmony]]
