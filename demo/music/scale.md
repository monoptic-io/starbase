---
title: Scale
aliases: [musical scale, major scale, minor scale, scales]
tags: [music]
summary: An ordered ladder of pitches within an octave — like the major scale's W-W-H-W-W-W-H pattern — that supplies the palette of notes a piece draws from.
weight: 30
---

# Scale

A **scale** is an ordered set of pitches spanning an octave — the palette of notes a melody is allowed to use. Where an [[Interval]] measures the gap between *two* notes, a scale is a whole *staircase* of intervals, and the exact pattern of step sizes is what gives each scale its flavor.

## Whole steps and half steps

The steps of a scale come in two sizes: a **half step** (one semitone, `H`) and a **whole step** (two semitones, `W`). The familiar **major scale** is built from the pattern

$$\text{W–W–H–W–W–W–H}$$

Starting on C and following that recipe lands exactly on the white keys: C–D–E–F–G–A–B–C. The two half-steps (E→F and B→C) are what make the major scale sound bright and "finished" when it returns home.

Change the pattern and you change the mood. The **natural minor** scale is W–H–W–W–H–W–W; its lowered third and sixth give it that wistful, shadowed quality. Different starting points within the major pattern produce the seven [[Mode|modes]], each with its own color.

{{< note kind="note" title="Scale, key, and notes" >}}
A scale is a *pattern of steps*; a [[Key Signature]] is the *bookkeeping* that records which notes that pattern lands on. "C major scale" and "the key of C major" describe the same seven pitches from two angles.
{{< /note >}}

## Light up a scale

Pick a scale with the buttons; the keyboard lights its notes in order, and the strip beneath shows the W/H step pattern. Notice how every major scale, wherever it starts, keeps the *same* shape — only shifted.

{{< sketch height="340" caption="Click a scale name to light its notes on the keyboard. The pattern strip shows the whole-step (W) and half-step (H) recipe; all major scales share one shape." >}}
if (frame === 0) {
  const cs = getComputedStyle(document.documentElement);
  state.accent    = cs.getPropertyValue('--accent').trim()     || '#5b9cff';
  state.accent2   = cs.getPropertyValue('--accent-2').trim()   || '#ff6b6b';
  state.good      = cs.getPropertyValue('--good').trim()       || '#46c98b';
  state.border    = cs.getPropertyValue('--border').trim()     || 'rgba(255,255,255,0.18)';
  state.textFaint = cs.getPropertyValue('--text-faint').trim() || 'rgba(255,255,255,0.5)';
  state.text      = cs.getPropertyValue('--text').trim()       || '#e8e8ef';
  state.pick = 0;
  state.wasDown = false;
}
const accent = state.accent, accent2 = state.accent2, good = state.good,
      border = state.border, textFaint = state.textFaint, text = state.text;
const NAMES  = ['C','C#','D','D#','E','F','F#','G','G#','A','A#','B'];
const WHITES = [0,2,4,5,7,9,11];
const BLACKS = [[1,0],[3,1],[6,3],[8,4],[10,5]];
const SCALES = [
  { name:'C major',        steps:[2,2,1,2,2,2,1] },
  { name:'A natural minor',steps:[2,1,2,2,1,2,2], root:9 },
  { name:'C pentatonic',   steps:[2,2,3,2,3] }
];
function scaleNotes(sc){ const r = sc.root||0; const out=[r]; let a=r;
  for(let i=0;i<sc.steps.length-1;i++){ a=(a+sc.steps[i])%12; out.push(a);} return out; }

function kbGeom() {
  const kx = W*0.05, kw = W*0.90, ww = kw/7;
  const ky = H*0.30, kh = H*0.46, bw = ww*0.62, bh = kh*0.62;
  return { kx, kw, ww, ky, kh, bw, bh };
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
function arpeggiate(sc) {                            // play the scale's notes in sequence
  if (!window.sgTone) return;
  const base = 261.63;                               // middle C
  let abs = sc.root || 0;
  const fs = [base * Math.pow(2, abs/12)];
  for (let i=0;i<sc.steps.length;i++){ abs += sc.steps[i]; fs.push(base * Math.pow(2, abs/12)); }
  for (let i=0;i<fs.length;i++){
    setTimeout(function(){ if (window.sgTone) sgTone(fs[i], 0.35); }, i*140);
  }
}

ctx.clearRect(0,0,W,H);
const clicked = mouse.clicked || (mouse.down && !state.wasDown);
state.wasDown = mouse.down;

// scale-selector buttons across the top
const bw0 = W*0.30, bh0 = H*0.12, gap = W*0.015, by = H*0.05;
ctx.textBaseline = 'middle'; ctx.textAlign = 'center';
for (let i=0;i<SCALES.length;i++){
  const bx = W*0.05 + i*(bw0+gap);
  if (clicked && mouse.x>=bx && mouse.x<=bx+bw0 && mouse.y>=by && mouse.y<=by+bh0) {
    state.pick = i; arpeggiate(SCALES[i]);           // arpeggiate the chosen scale
  }
  const on = state.pick===i;
  ctx.fillStyle = on ? accent : 'rgba(255,255,255,0.04)';
  ctx.fillRect(bx, by, bw0, bh0);
  ctx.strokeStyle = border; ctx.strokeRect(bx, by, bw0, bh0);
  ctx.fillStyle = on ? '#0c0f16' : text; ctx.font = '13px sans-serif';
  ctx.fillText(SCALES[i].name, bx+bw0/2, by+bh0/2);
}

// click a key on the keyboard to hear a single note
if (clicked) { const k = keyAt(mouse.x, mouse.y); if (k>=0 && window.sgTone) sgTone(261.63 * Math.pow(2, k/12), 0.5); }

const sc = SCALES[state.pick];
const notes = scaleNotes(sc);
const hl = {};
for (let i=0;i<notes.length;i++) hl[notes[i]] = (i===0) ? accent2 : good;
drawKeyboard(hl);

// step pattern strip
const g = kbGeom();
ctx.textBaseline = 'alphabetic'; ctx.textAlign = 'left';
ctx.fillStyle = textFaint; ctx.font = '13px sans-serif';
const pat = sc.steps.map(s => s===1?'H':(s===2?'W':'W+')).join(' – ');
ctx.fillText('steps:  ' + pat, g.kx, g.ky + g.kh + 28);
ctx.fillStyle = text; ctx.font = '600 14px sans-serif';
ctx.fillText('notes:  ' + notes.map(n=>NAMES[n]).join('  '), g.kx, g.ky + g.kh + 50);
{{< /sketch >}}

## See also

- [[Interval]]
- [[Mode]]
- [[Key Signature]]
