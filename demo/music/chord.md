---
title: Chord
aliases: [chords, triad, triads]
tags: [music]
summary: Three or more pitches sounded together, most often built by stacking thirds into a triad — and whether that triad sounds major (bright) or minor (dark).
weight: 40
---

# Chord

A **chord** is three or more pitches sounded at once. The workhorse of Western music is the simplest chord of all: the **triad**, three notes built by stacking two [[Interval|thirds]] on top of a root. Take a root, add the note a third above it, then a third above *that*, and you have spanned a perfect fifth from bottom to top — a stable, full-bodied sound.

## Major versus minor

A third comes in two sizes, and which one you put on the bottom decides the whole mood of the chord:

- **Major triad** — a *major* third (4 semitones) then a *minor* third (3): intervals `0–4–7`. Bright, open, resolved. C–E–G.
- **Minor triad** — a *minor* third (3) then a *major* third (4): intervals `0–3–7`. Darker, more pensive. C–E♭–G.

Both span the same perfect fifth (7 semitones); only the **middle** note moves, by a single semitone — yet that one semitone is the difference between sunshine and shadow. This is the cell from which all richer [[Harmony]] is grown.

{{< note kind="tip" title="Stacking further" >}}
Keep stacking thirds and the triad grows: add one more and you get a *seventh chord* (four notes), the staple of jazz and blues. Every such chord still traces back to the same recipe — thirds piled on a root.
{{< /note >}}

## Build a triad

Move the mouse to choose the **root**; the keyboard shows the resulting triad — root in blue, third in orange, fifth in green. Use the **Major / minor** button to flip the quality and watch only the middle note slide by a semitone.

{{< sketch height="340" caption="Hover to set the root; the triad's three notes light up (root blue, third orange, fifth green). Toggle Major/minor and watch only the middle note shift one semitone." >}}
if (frame === 0) {
  const cs = getComputedStyle(document.documentElement);
  state.accent    = cs.getPropertyValue('--accent').trim()     || '#5b9cff';
  state.accent2   = cs.getPropertyValue('--accent-2').trim()   || '#ff6b6b';
  state.good      = cs.getPropertyValue('--good').trim()       || '#46c98b';
  state.border    = cs.getPropertyValue('--border').trim()     || 'rgba(255,255,255,0.18)';
  state.textFaint = cs.getPropertyValue('--text-faint').trim() || 'rgba(255,255,255,0.5)';
  state.text      = cs.getPropertyValue('--text').trim()       || '#e8e8ef';
  state.root = 0;
  state.minor = false;
  state.wasDown = false;
}
const accent = state.accent, accent2 = state.accent2, good = state.good,
      border = state.border, textFaint = state.textFaint, text = state.text;
const NAMES  = ['C','C#','D','D#','E','F','F#','G','G#','A','A#','B'];
const WHITES = [0,2,4,5,7,9,11];
const BLACKS = [[1,0],[3,1],[6,3],[8,4],[10,5]];

function kbGeom() {
  const kx = W*0.05, kw = W*0.90, ww = kw/7;
  const ky = H*0.30, kh = H*0.50, bw = ww*0.62, bh = kh*0.62;
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
    ctx.font = '11px sans-serif';
    ctx.fillText(NAMES[s], x+g.ww/2, g.ky+g.kh-8);
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

// quality toggle button
let triadChanged = false;
const bx = W*0.05, by = H*0.05, bw0 = W*0.34, bh0 = H*0.12;
if (clicked && mouse.x>=bx && mouse.x<=bx+bw0 && mouse.y>=by && mouse.y<=by+bh0) {
  state.minor = !state.minor; triadChanged = true;
}
ctx.textBaseline = 'middle'; ctx.textAlign = 'center';
ctx.fillStyle = state.minor ? 'rgba(255,255,255,0.04)' : accent;
ctx.fillRect(bx, by, bw0, bh0);
ctx.strokeStyle = border; ctx.strokeRect(bx, by, bw0, bh0);
ctx.fillStyle = state.minor ? text : '#0c0f16'; ctx.font = '13px sans-serif';
ctx.fillText(state.minor ? 'minor triad' : 'Major triad', bx+bw0/2, by+bh0/2);

// root from hover (over the keyboard area only); a click locks it and sounds the chord
const g = kbGeom();
if (mouse.y > g.ky - 4) {
  const k = keyAt(mouse.x, mouse.y);
  if (k>=0) { state.root = k; if (clicked) triadChanged = true; }
}
const r = state.root;
const third = state.minor ? 3 : 4;
const notes = [r, (r+third)%12, (r+7)%12];
if (triadChanged && window.sgTone) {                 // play the three notes together
  const base = 261.63;
  sgChord([base*Math.pow(2, r/12), base*Math.pow(2, (r+third)/12), base*Math.pow(2, (r+7)/12)], 0.9);
}
const cols  = [accent, accent2, good];
const hl = {};
for (let i=0;i<3;i++) hl[notes[i]] = cols[i];
drawKeyboard(hl);

ctx.textBaseline = 'alphabetic'; ctx.textAlign = 'left';
ctx.fillStyle = text; ctx.font = '600 18px sans-serif';
ctx.fillText(NAMES[r] + (state.minor?' minor':' major') + ':  ' + notes.map(n=>NAMES[n]).join(' – '),
             bx + bw0 + W*0.04, by + bh0/2);
{{< /sketch >}}

## See also

- [[Harmony]]
- [[Interval]]
