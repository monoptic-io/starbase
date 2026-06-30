---
title: Interval
aliases: [musical interval, semitone, step]
tags: [music]
summary: The distance between two pitches, measured in semitones, and the reason some gaps sound sweet and consonant while others sound tense.
weight: 20
---

# Interval

An **interval** is the distance between two pitches. It is the atom of harmony: once you can hear *how far apart* two notes are, you can hear everything built from them. Intervals are counted in **semitones** — the smallest step on the keyboard — and each size has a traditional name: a gap of 7 semitones is a *perfect fifth*, 4 a *major third*, 12 an *octave*.

## Why some intervals sound sweet

Consonance is not a matter of taste; it tracks arithmetic. Two tones sound *consonant* when their frequencies form a simple whole-number ratio, because then their overtones line up instead of clashing. The octave is the simplest ratio of all, 2:1. The perfect fifth is 3:2; the perfect fourth, 4:3; the major third, roughly 5:4. The more complicated the ratio, the more the tones beat against each other and the more *tension* we hear.

| Interval | Semitones | Approx. ratio | Character |
|---|---|---|---|
| Unison | 0 | 1:1 | identical |
| Major third | 4 | 5:4 | bright, sweet |
| Perfect fourth | 5 | 4:3 | open, stable |
| Tritone | 6 | 45:32 | restless, tense |
| Perfect fifth | 7 | 3:2 | strong, hollow |
| Octave | 12 | 2:1 | same note, higher |

Because [[Pitch]] is logarithmic, intervals *add*: stack a major third (4) on top of a minor third (3) and you get a perfect fifth (7). That simple bit of arithmetic is the seed of every [[Chord]].

## Hear the distance

The keyboard fixes a **root** on C (blue). Move the mouse to pick the **second note** (orange); the panel names the interval and tells you whether it leans consonant or tense. Click to freeze your choice.

{{< sketch height="320" caption="Root C is fixed in blue. Hover to choose the upper note (orange); the readout names the interval and its consonance. Click to lock it in." >}}
if (frame === 0) {
  const cs = getComputedStyle(document.documentElement);
  state.accent    = cs.getPropertyValue('--accent').trim()     || '#5b9cff';
  state.accent2   = cs.getPropertyValue('--accent-2').trim()   || '#ff6b6b';
  state.good      = cs.getPropertyValue('--good').trim()       || '#46c98b';
  state.warn      = cs.getPropertyValue('--warn').trim()       || '#ffd166';
  state.border    = cs.getPropertyValue('--border').trim()     || 'rgba(255,255,255,0.18)';
  state.textFaint = cs.getPropertyValue('--text-faint').trim() || 'rgba(255,255,255,0.5)';
  state.text      = cs.getPropertyValue('--text').trim()       || '#e8e8ef';
  state.sel = 7;          // start on the fifth
  state.wasDown = false;
}
const accent = state.accent, accent2 = state.accent2, good = state.good, warn = state.warn,
      border = state.border, textFaint = state.textFaint, text = state.text;
const NAMES  = ['C','C#','D','D#','E','F','F#','G','G#','A','A#','B'];
const IVAL = ['unison','minor 2nd','major 2nd','minor 3rd','major 3rd','perfect 4th',
              'tritone','perfect 5th','minor 6th','major 6th','minor 7th','major 7th','octave'];
const TENSE = [0,2,1,1,0,0,2,0,1,0,1,2,0]; // 0 consonant, 1 mild, 2 tense
const WHITES = [0,2,4,5,7,9,11];
const BLACKS = [[1,0],[3,1],[6,3],[8,4],[10,5]];

function kbGeom() {
  const kx = W*0.05, kw = W*0.90, ww = kw/7;
  const ky = H*0.30, kh = H*0.55, bw = ww*0.62, bh = kh*0.62;
  return { kx, kw, ww, ky, kh, bw, bh };
}
function keyAt(mx, my) {
  const g = kbGeom();
  for (const bk of BLACKS) {
    const cx = g.kx + (bk[1]+1)*g.ww, x = cx - g.bw/2;
    if (mx>=x && mx<=x+g.bw && my>=g.ky && my<=g.ky+g.bh) return bk[0];
  }
  for (let i=0;i<7;i++){
    const x = g.kx + i*g.ww;
    if (mx>=x && mx<=x+g.ww && my>=g.ky && my<=g.ky+g.kh) return WHITES[i];
  }
  return -1;
}
function drawKeyboard(hl) {
  const g = kbGeom();
  ctx.textAlign = 'center'; ctx.textBaseline = 'alphabetic';
  for (let i=0;i<7;i++){
    const s = WHITES[i], x = g.kx + i*g.ww;
    ctx.fillStyle = hl[s] || '#f6f7fb';
    ctx.fillRect(x+1, g.ky, g.ww-2, g.kh);
    ctx.strokeStyle = border; ctx.lineWidth = 1;
    ctx.strokeRect(x+1, g.ky, g.ww-2, g.kh);
    ctx.fillStyle = hl[s] ? '#0c0f16' : textFaint;
    ctx.font = '12px sans-serif';
    ctx.fillText(NAMES[s], x+g.ww/2, g.ky+g.kh-10);
  }
  for (const bk of BLACKS){
    const s = bk[0], cx = g.kx+(bk[1]+1)*g.ww, x = cx-g.bw/2;
    ctx.fillStyle = hl[s] || '#161a24';
    ctx.fillRect(x, g.ky, g.bw, g.bh);
    ctx.strokeStyle = border; ctx.strokeRect(x, g.ky, g.bw, g.bh);
  }
}

ctx.clearRect(0,0,W,H);
const clicked = mouse.clicked || (mouse.down && !state.wasDown);
state.wasDown = mouse.down;
const hover = keyAt(mouse.x, mouse.y);
if (clicked && hover>=0) {
  state.sel = hover;
  if (window.sgTone) {                              // sound both notes of the interval
    const C = 261.63;                               // root C
    sgChord([C, C * Math.pow(2, hover/12)], 0.8);
  }
}
const second = (hover>=0) ? hover : state.sel;

const hl = {};
hl[second] = accent2;
hl[0] = accent;                                  // root C stays blue
const semis = second;                            // distance above C within the octave
const tcol = TENSE[semis]===0 ? good : (TENSE[semis]===1 ? warn : accent2);
drawKeyboard(hl);

ctx.textAlign = 'left';
ctx.fillStyle = text; ctx.font = '600 19px sans-serif';
ctx.fillText('C  to  ' + NAMES[second] + '  —  ' + IVAL[semis], W*0.05, H*0.16);
ctx.fillStyle = tcol; ctx.font = '14px sans-serif';
const label = TENSE[semis]===0 ? 'consonant' : (TENSE[semis]===1 ? 'mild' : 'tense');
ctx.fillText(semis + ' semitones · ' + label, W*0.05, H*0.16 + 22);
{{< /sketch >}}

{{< quiz question="How many semitones is a perfect fifth?" options="4|5|7|12" answer="3" explain="A perfect fifth spans 7 semitones and approximates the frequency ratio 3:2 — one of the most consonant intervals after the octave." >}}

## See also

- [[Pitch]]
- [[Chord]]
