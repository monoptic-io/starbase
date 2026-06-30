---
title: Rhythm
aliases: [meter, beat, tempo, time signature]
tags: [music]
summary: The organization of music in time — beats grouped into measures by a time signature, ticking past at a tempo set in beats per minute.
weight: 60
---

# Rhythm

If [[Pitch]] is the *vertical* dimension of music, **rhythm** is the *horizontal* one: the organization of sound in **time**. Even with no pitch at all — a drum, a clap, a footstep — rhythm alone can carry a groove. It rests on three nested ideas: the **beat**, the **measure**, and the **tempo**.

## Beat, meter, and tempo

- The **beat** is the steady pulse you tap your foot to — the heartbeat of the music.
- **Meter** groups beats into repeating bundles called **measures** (or *bars*). A **time signature** like 4/4 or 3/4 records the grouping: the top number is *how many beats per measure*, the bottom *which note value gets the beat*. The first beat of each measure, the **downbeat**, feels strongest.
- **Tempo** is how fast the beats go, in *beats per minute* (BPM). At 120 BPM each beat lasts exactly half a second:

$$\text{beat duration} = \frac{60}{\text{BPM}}\ \text{seconds}.$$

The same melody at 60 BPM feels solemn; at 160 BPM it races. Same pitches, same meter — only the clock changed.

## A live metronome

The metronome below swings at a fixed tempo while a row of beat-lights counts out the measure. **Click** to cycle the time signature (4/4 → 3/4 → 6/8) and watch the grouping change — the bright **downbeat** lands on beat one every measure.

{{< sketch height="320" caption="A swinging metronome and beat grid. Click to cycle the time signature; the accented downbeat (beat 1) marks the start of every measure." >}}
if (frame === 0) {
  const cs = getComputedStyle(document.documentElement);
  state.accent    = cs.getPropertyValue('--accent').trim()     || '#5b9cff';
  state.accent2   = cs.getPropertyValue('--accent-2').trim()   || '#ff6b6b';
  state.good      = cs.getPropertyValue('--good').trim()       || '#46c98b';
  state.border    = cs.getPropertyValue('--border').trim()     || 'rgba(255,255,255,0.18)';
  state.textFaint = cs.getPropertyValue('--text-faint').trim() || 'rgba(255,255,255,0.5)';
  state.text      = cs.getPropertyValue('--text').trim()       || '#e8e8ef';
  state.sigs = [[4,4],[3,4],[6,8]];
  state.pick = 0;
  state.beats = 0;       // elapsed beats (continuous)
  state.lastBeat = -1;   // which beat last triggered a tick
  state.wasDown = false;
}
const accent = state.accent, accent2 = state.accent2, good = state.good,
      border = state.border, textFaint = state.textFaint, text = state.text;

ctx.clearRect(0,0,W,H);
const clicked = mouse.clicked || (mouse.down && !state.wasDown);
state.wasDown = mouse.down;
if (clicked) state.pick = (state.pick + 1) % state.sigs.length;

const BPM = 100;
const beatsPerBar = state.sigs[state.pick][0];
state.beats += dt * (BPM / 60);
const beatPos = state.beats % beatsPerBar;          // 0..beatsPerBar
const curBeat = Math.floor(beatPos);

// tick exactly once when the metronome lands on a new beat
if (frame>0 && curBeat !== state.lastBeat) {
  if (window.sgTone) {
    if (curBeat === 0) sgTone(1500, 0.05, 'square');  // accented downbeat
    else               sgTone(900, 0.04);             // other beats
  }
  state.lastBeat = curBeat;
}

// --- metronome: a pendulum swinging once per beat ---
const px = W*0.5, py = H*0.78, armLen = Math.min(H*0.5, W*0.4);
const swing = 0.6 * Math.sin(beatPos % 1 * Math.PI*2 - Math.PI/2); // -.6..+.6 rad
const ang = swing;
const bobx = px + armLen * Math.sin(ang), boby = py - armLen * Math.cos(ang);
ctx.strokeStyle = border; ctx.lineWidth = 2;
ctx.beginPath(); ctx.moveTo(px, py); ctx.lineTo(bobx, boby); ctx.stroke();
// pivot
ctx.fillStyle = textFaint; ctx.beginPath(); ctx.arc(px, py, 5, 0, 7); ctx.fill();
// bob, pulsing brighter near each beat tick
const tick = 1 - Math.abs((beatPos % 1) - 0.5) * 2;   // 0..1, peak mid-swing
ctx.fillStyle = accent;
ctx.beginPath(); ctx.arc(bobx, boby, 8 + 4*tick, 0, 7); ctx.fill();

// --- beat-light row ---
const n = beatsPerBar, gap = W*0.04;
const rw = Math.min(W*0.10, (W*0.7 - gap*(n-1)) / n);
const totalW = n*rw + (n-1)*gap, sx = (W - totalW)/2, sy = H*0.16;
ctx.textAlign = 'center'; ctx.textBaseline = 'middle';
for (let i=0;i<n;i++){
  const x = sx + i*(rw+gap);
  const on = (i === curBeat);
  const down = (i === 0);
  ctx.fillStyle = on ? (down ? accent2 : good) : 'rgba(255,255,255,0.06)';
  ctx.fillRect(x, sy, rw, rw*0.7);
  ctx.strokeStyle = border; ctx.strokeRect(x, sy, rw, rw*0.7);
  ctx.fillStyle = on ? '#0c0f16' : textFaint; ctx.font = '13px sans-serif';
  ctx.fillText(String(i+1), x+rw/2, sy + rw*0.35);
}

ctx.textAlign = 'left'; ctx.textBaseline = 'alphabetic';
ctx.fillStyle = text; ctx.font = '600 18px sans-serif';
ctx.fillText(state.sigs[state.pick][0] + '/' + state.sigs[state.pick][1], W*0.05, H*0.10);
ctx.fillStyle = textFaint; ctx.font = '13px sans-serif';
ctx.fillText(BPM + ' BPM   (click to change time signature)', W*0.05 + 50, H*0.10);
{{< /sketch >}}

A waltz in 3/4 and a march in 4/4 can share the very same notes; what makes one *swing* and the other *stride* is purely how their beats are grouped.

## See also

- [[Pitch]]
