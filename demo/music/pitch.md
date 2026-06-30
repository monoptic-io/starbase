---
title: Pitch
aliases: [note, musical note, octave]
tags: [music]
summary: How high or low a tone sounds, set by its frequency, and why the twelve-tone octave repeats all the way up and down the keyboard.
weight: 10
---

# Pitch

**Pitch** is how high or low a tone sounds. Physically it is governed by *frequency* — the number of times per second the air pressure cycles — measured in hertz (Hz). A slow vibration (few cycles per second) sounds low; a fast one sounds high. The single most important fact about pitch is that our ears hear it *logarithmically*: doubling the frequency does not sound "twice as high," it sounds like *the same note, one octave up*.

## The octave and the twelve tones

When one frequency is exactly double another — say 220 Hz and 440 Hz — the two tones blend so completely that we give them the same letter name. That doubling is the **octave**, and it is the closest thing music has to a law of nature.

Western music slices each octave into **twelve equal steps**, the *semitones*. Climb twelve semitones and you have multiplied the frequency by two; you are back where you started, an octave higher. Because every step multiplies the frequency by the same factor, that factor must be the twelfth root of two, $2^{1/12} \approx 1.0595$. This scheme is called **equal temperament**, and it lets the same note land in tune in every key.

Taking the note A above middle C as a reference of 440 Hz, any note $n$ semitones away has frequency

{{< eq number="1" >}}f = 440 \cdot 2^{\,n/12}{{< /eq >}}

so $n = +12$ gives 880 Hz (an octave up) and $n = -12$ gives 220 Hz (an octave down). Notice this is closely related to the way a vibrating string's overtones stack up in whole-number multiples — but pitch as we *name* it is the logarithmic ladder above, not that overtone series.

{{< note kind="note" title="Why twelve?" >}}
Twelve is a happy accident of arithmetic: powers of $2^{1/12}$ land remarkably close to the simple frequency ratios (like 3:2 and 4:3) that the ear finds pleasing. Those near-misses are exactly what makes an [[Interval]] sound consonant.
{{< /note >}}

## Play an octave

The keyboard below is one octave: seven white keys (C D E F G A B) and five black keys (the sharps and flats) wedged between them. **Click a key** to select it — the panel shows its name, its distance $n$ in semitones from A, and the frequency from Equation&nbsp;1.

{{< sketch height="320" caption="A one-octave keyboard. Click any key to hear its place in the twelve-tone system: name, semitone offset n from A4, and frequency f = 440·2^(n/12)." >}}
if (frame === 0) {
  const cs = getComputedStyle(document.documentElement);
  state.accent    = cs.getPropertyValue('--accent').trim()     || '#5b9cff';
  state.accent2   = cs.getPropertyValue('--accent-2').trim()   || '#ff6b6b';
  state.border    = cs.getPropertyValue('--border').trim()     || 'rgba(255,255,255,0.18)';
  state.textFaint = cs.getPropertyValue('--text-faint').trim() || 'rgba(255,255,255,0.5)';
  state.text      = cs.getPropertyValue('--text').trim()       || '#e8e8ef';
  state.sel = 9;          // start on A
  state.wasDown = false;
}
const accent = state.accent, accent2 = state.accent2, border = state.border,
      textFaint = state.textFaint, text = state.text;
const NAMES  = ['C','C#','D','D#','E','F','F#','G','G#','A','A#','B'];
const WHITES = [0,2,4,5,7,9,11];
const BLACKS = [[1,0],[3,1],[6,3],[8,4],[10,5]]; // [semitone, white-to-its-left]

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
if (clicked) {
  const k = keyAt(mouse.x, mouse.y);
  if (k>=0) {
    state.sel = k;
    if (window.sgTone) sgTone(440 * Math.pow(2, (k-9)/12), 0.6);  // sound the pressed note
  }
}
const hover = keyAt(mouse.x, mouse.y);

const hl = {};
if (hover>=0) hl[hover] = accent2;
hl[state.sel] = accent;
drawKeyboard(hl);

// info panel
const s = state.sel, n = s - 9;                 // semitones from A4
const f = 440 * Math.pow(2, n/12);
ctx.textAlign = 'left';
ctx.fillStyle = text; ctx.font = '600 20px sans-serif';
ctx.fillText('Note ' + NAMES[s], W*0.05, H*0.18);
ctx.fillStyle = textFaint; ctx.font = '14px sans-serif';
ctx.fillText('n = ' + (n>=0?'+':'') + n + ' semitones from A    f = ' + f.toFixed(1) + ' Hz',
             W*0.05, H*0.18 + 22);
{{< /sketch >}}

Every key you press is some $n$ in Equation&nbsp;1; the black keys simply fill in the half-steps between the white ones. Twelve presses to the right and the pattern — and the frequency ratio of 2 — repeats exactly.

## See also

- [[Interval]]
- [[Scale]]
