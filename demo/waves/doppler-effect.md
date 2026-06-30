---
title: Doppler Effect
aliases: [doppler shift]
tags: [waves]
summary: Motion of the source or observer compresses the wavefronts ahead and stretches them behind, raising the frequency you receive as it approaches and lowering it as it recedes.
weight: 70
---

# Doppler Effect

The **Doppler effect** is the change in observed frequency caused by relative motion between a [[Wave]] source and an observer. It is the rising-then-falling wail of a passing siren, the pitch drop of a race car, and — for light — the redshift that tells us distant galaxies are flying away. The wave's speed through the medium does not change; what changes is how the *wavefronts pile up*.

## Why the fronts bunch

Picture a source emitting a crest every period $T$. If the source is **moving toward you**, it creeps a little closer between successive emissions, so each new crest leaves from slightly nearer than the last. The crests crowd together ahead of the source — shorter wavelength, **higher frequency**. Behind the source they spread apart — longer wavelength, **lower frequency**. The wavefronts themselves stay perfect circles; it is only their *centers* that march forward.

For a source moving at speed $v_s$ through a medium where waves travel at $v$, the frequency an observer receives is

$$f_{\text{obs}}=f_{\text{src}}\,\frac{v}{v \mp v_s},$$

with the minus sign (higher pitch) when the source approaches and the plus sign (lower pitch) when it recedes.

{{< sketch height="420" caption="A source (green) glides rightward, emitting a circular wavefront at a steady beat. The fronts bunch ahead of it and spread behind — the Doppler effect made geometric. Move the mouse left/right to set the source speed; push it past the wave speed and the fronts pile into a Mach cone (a sonic boom)." >}}
function cv(n){var v=getComputedStyle(document.documentElement).getPropertyValue(n).trim();return v||'#5b9cff';}
if (frame===0){ state.waves=[]; state.sx=W*0.06; state.acc=0; state.vfrac=0.45; }
if (mouse.x>0 && mouse.x<W){ state.vfrac = (mouse.x / W) * 1.4; }  // 0..1.4 of wave speed

var c = W*0.16;                 // wave speed (px/s)
var vsrc = state.vfrac * c;     // source speed
var period = 0.45;              // emission interval (s)
var sy = H*0.5;

// advance source; reset at right edge
state.sx += vsrc * dt;
if (state.sx > W*0.94){ state.sx = W*0.06; state.waves = []; }

// emit wavefronts on the beat
state.acc += dt;
while (state.acc >= period){
  state.acc -= period;
  state.waves.push({x: state.sx, y: sy, t0: t});
}
// cull old fronts that have left the canvas
if (state.waves.length > 60) state.waves.shift();

ctx.clearRect(0,0,W,H);

// draw expanding fronts, fading with age
for (var i=0; i<state.waves.length; i++){
  var wv = state.waves[i];
  var r = c * (t - wv.t0);
  if (r <= 1) continue;
  var age = (t - wv.t0);
  var alpha = Math.max(0, 0.85 - age*0.20);
  ctx.strokeStyle = cv('--accent');
  ctx.globalAlpha = alpha;
  ctx.lineWidth = 1.6;
  ctx.beginPath(); ctx.arc(wv.x, wv.y, r, 0, Math.PI*2); ctx.stroke();
}
ctx.globalAlpha = 1;

// the source
ctx.fillStyle = cv('--good');
ctx.beginPath(); ctx.arc(state.sx, sy, 6, 0, Math.PI*2); ctx.fill();
// little velocity arrow
ctx.strokeStyle = cv('--good'); ctx.lineWidth=2;
ctx.beginPath(); ctx.moveTo(state.sx, sy); ctx.lineTo(state.sx+22, sy); ctx.stroke();

ctx.fillStyle = cv('--text-dim'); ctx.font="13px sans-serif";
var label = "source speed = " + state.vfrac.toFixed(2) + " × wave speed";
if (state.vfrac > 1.0) label += "  (supersonic — shock cone!)";
ctx.fillText(label, 12, 22);
ctx.fillStyle = cv('--text-faint');
ctx.fillText("ahead: bunched, higher f", 12, H-30);
ctx.fillText("behind: stretched, lower f", 12, H-12);
{{< /sketch >}}

## The supersonic limit

Push the source past the wave speed and the geometry changes character. The source now outruns its own wavefronts, which stack up along a cone — the **Mach cone**. For sound that cone is the shock wave heard as a sonic boom; for a charged particle in a medium it is the optical analogue, Cherenkov radiation. The Doppler formula above diverges at $v_s=v$ precisely because the fronts ahead collapse onto a single surface.

## See also

- [[Wave]]
- [[Interference]]
- [[Oscillations]]
