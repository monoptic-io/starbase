---
title: Beats
aliases: [beat frequency]
tags: [waves]
summary: Two tones of nearly equal frequency drift in and out of step, producing a slow throb whose rate is the difference of the two frequencies.
weight: 60
---

# Beats

**Beats** are interference unfolding in *time*. Sound two notes of almost the same pitch — two guitar strings a hair out of tune — and you hear neither a clean chord nor a clash, but a single tone that **pulses**: loud, soft, loud, soft. That slow throb is the two frequencies repeatedly drifting into step (reinforcing) and out of step (cancelling), exactly the [[Interference]] rule applied to overlap in time rather than space.

## The envelope at the average, the throb at the difference

Add two equal-amplitude tones at $f_1$ and $f_2$ and a product-to-sum identity rearranges the sum into a fast carrier riding inside a slow envelope:

$$\cos(2\pi f_1 t)+\cos(2\pi f_2 t)=\underbrace{2\cos\!\big(2\pi \tfrac{f_1-f_2}{2}\,t\big)}_{\text{slow envelope}}\;\underbrace{\cos\!\big(2\pi \tfrac{f_1+f_2}{2}\,t\big)}_{\text{fast carrier}}.$$

The ear hears the carrier at the **average** frequency $(f_1+f_2)/2$, swelling and fading under the envelope. Because loudness peaks *twice* per envelope cycle (at both the positive and negative bulge), the audible **beat frequency** is the full difference:

{{< eq number="1" >}}f_{\text{beat}}=|f_1-f_2|{{< /eq >}}

Tune one string toward the other and the beats slow down; when they vanish entirely, the strings are in unison. Piano tuners do exactly this by ear.

{{< sketch height="360" caption="Top two traces: two pure tones at nearby frequencies. Bottom trace: their sum, with the slow beat envelope drawn faintly. The sum swells and collapses where the tones fall in and out of phase. Drag left/right to detune the second tone — more detuning, faster beats." >}}
function cv(n){var v=getComputedStyle(document.documentElement).getPropertyValue(n).trim();return v||'#5b9cff';}
if (frame===0){ state.detune = 0.10; }
if (mouse.x>0 && mouse.x<W){ state.detune = 0.02 + (mouse.x / W) * 0.28; }
ctx.clearRect(0,0,W,H);
var f1 = 1.0, f2 = 1.0 + state.detune;
var ph = t * 2.4;                  // scroll in time
var rowH = H/3;
var y1 = rowH*0.5, y2 = rowH*1.5, y3 = rowH*2.5;
var amp = rowH*0.32;
var scale = 26;                    // spatial cycles across canvas

function w1(x){ return Math.cos(f1*(x/W)*scale + ph); }
function w2(x){ return Math.cos(f2*(x/W)*scale + ph); }
function envv(x){ return 2*Math.abs(Math.cos(((f2-f1)/2)*(x/W)*scale + 0)); }

// faint baselines
ctx.strokeStyle = cv('--border'); ctx.lineWidth=1;
[y1,y2,y3].forEach(function(yc){ ctx.beginPath(); ctx.moveTo(0,yc); ctx.lineTo(W,yc); ctx.stroke(); });

// tone 1
ctx.strokeStyle = cv('--accent'); ctx.lineWidth=1.8;
ctx.beginPath();
for (var x=0;x<=W;x+=2){ var y=y1-amp*w1(x); if(x===0)ctx.moveTo(x,y);else ctx.lineTo(x,y);} ctx.stroke();
// tone 2
ctx.strokeStyle = cv('--accent-2'); ctx.lineWidth=1.8;
ctx.beginPath();
for (var xb=0;xb<=W;xb+=2){ var yb=y2-amp*w2(xb); if(xb===0)ctx.moveTo(xb,yb);else ctx.lineTo(xb,yb);} ctx.stroke();

// envelope (faint) on the sum
ctx.strokeStyle = cv('--text-faint'); ctx.lineWidth=1; ctx.globalAlpha=0.6;
ctx.beginPath();
for (var xe=0;xe<=W;xe+=2){ var ye=y3-(amp*0.5)*envv(xe); if(xe===0)ctx.moveTo(xe,ye);else ctx.lineTo(xe,ye);} ctx.stroke();
ctx.beginPath();
for (var xf=0;xf<=W;xf+=2){ var yf=y3+(amp*0.5)*envv(xf); if(xf===0)ctx.moveTo(xf,yf);else ctx.lineTo(xf,yf);} ctx.stroke();
ctx.globalAlpha=1;
// the sum
ctx.strokeStyle = cv('--good'); ctx.lineWidth=2.2;
ctx.beginPath();
for (var xs=0;xs<=W;xs+=2){ var ys=y3-(amp*0.5)*(w1(xs)+w2(xs)); if(xs===0)ctx.moveTo(xs,ys);else ctx.lineTo(xs,ys);} ctx.stroke();

ctx.fillStyle = cv('--text-dim'); ctx.font="13px sans-serif";
ctx.fillText("f1", 10, y1-amp-6);
ctx.fillText("f2", 10, y2-amp-6);
ctx.fillText("sum  (beat = |f1−f2|)", 10, y3-amp-6);
{{< /sketch >}}

## The same physics as coupled pendulums

If beats feel familiar, they should: two weakly [[Coupled Oscillators]] do exactly this. Their two normal modes sit at slightly different frequencies, and energy sloshes fully from one pendulum to the other and back at the difference frequency — a beat you can watch instead of hear. Near [[Resonance]], the difference between drive and natural frequency similarly controls how the response builds and ebbs.

## See also

- [[Coupled Oscillators]]
- [[Resonance]]
- [[Interference]]
