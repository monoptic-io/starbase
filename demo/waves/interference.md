---
title: Interference
aliases: [wave interference]
tags: [waves]
summary: Two overlapping waves reinforce where their crests align and cancel where crest meets trough, painting a fixed pattern of bright and dark fringes.
weight: 50
---

# Interference

**Interference** is the [[Superposition Principle]] made visible. When two waves overlap, their displacements add — but whether they *reinforce* or *cancel* depends on whether they arrive in step. Align crest with crest and you get **constructive interference**, doubling the amplitude. Align crest with trough and you get **destructive interference**, flat nothing. From two perfectly ordinary wave sources emerges a stable, sculpted pattern of loud and quiet, bright and dark.

## Path difference sets the fringes

Consider two sources oscillating in phase, a distance apart, and a point that sits a distance $r_1$ from one and $r_2$ from the other. The two waves arrive having traveled different distances, so they are out of step by the **path difference** $\Delta r = r_2 - r_1$. The rule is entirely about how many wavelengths fit in that gap:

$$\Delta r = m\lambda \;\Rightarrow\; \textbf{constructive},\qquad \Delta r = \left(m+\tfrac12\right)\lambda \;\Rightarrow\; \textbf{destructive},$$

for integer $m$. The set of points with a *constant* path difference is a hyperbola, which is why the bright fringes fan out in smooth curves — visible directly in the ripple tank below.

{{< sketch height="420" caption="A two-source ripple tank. Each source emits identical circular waves; the field is their sum. Bright violet bands are crests reinforcing (constructive), accent bands are troughs, and the gray channels between are nodal lines where the waves cancel. Move the mouse left/right to change the source spacing and watch the fringes fan." >}}
function cv(n){var v=getComputedStyle(document.documentElement).getPropertyValue(n).trim();return v||'#5b9cff';}
function hx(c){c=(c||'#5b9cff').replace('#','');if(c.length===3)c=c[0]+c[0]+c[1]+c[1]+c[2]+c[2];return [parseInt(c.slice(0,2),16),parseInt(c.slice(2,4),16),parseInt(c.slice(4,6),16)];}
if (frame===0){
  state.cell = 5;
  state.A = hx(cv('--accent'));    // troughs
  state.B = hx(cv('--accent-2'));  // crests
  state.D = hx(cv('--bg-elev2')); if(!state.D) state.D=[20,24,32];
  state.sep = 0.18;
}
if (mouse.x>0 && mouse.x<W){ state.sep = 0.07 + (mouse.x / W) * 0.30; }

var cell = state.cell;
var k = 0.085;        // wavenumber (px^-1)
var w = 4.0;          // angular frequency
var ph = w * t;
var s1x = W*(0.5 - state.sep), s1y = H*0.5;
var s2x = W*(0.5 + state.sep), s2y = H*0.5;
var A=state.A, B=state.B, D=state.D;

for (var py=0; py<H; py+=cell){
  for (var px=0; px<W; px+=cell){
    var dx1=px-s1x, dy1=py-s1y, d1=Math.sqrt(dx1*dx1+dy1*dy1);
    var dx2=px-s2x, dy2=py-s2y, d2=Math.sqrt(dx2*dx2+dy2*dy2);
    // 1/sqrt(r) amplitude falloff keeps near-source from saturating
    var a1 = 1/Math.sqrt(1+d1*0.02);
    var a2 = 1/Math.sqrt(1+d2*0.02);
    var v = a1*Math.cos(k*d1 - ph) + a2*Math.cos(k*d2 - ph);  // range ~[-2,2]
    var s = v/2;            // [-1,1]
    var r,g,b;
    if (s>=0){ var u=s; r=D[0]+(B[0]-D[0])*u; g=D[1]+(B[1]-D[1])*u; b=D[2]+(B[2]-D[2])*u; }
    else { var u2=-s; r=D[0]+(A[0]-D[0])*u2; g=D[1]+(A[1]-D[1])*u2; b=D[2]+(A[2]-D[2])*u2; }
    ctx.fillStyle = 'rgb('+(r|0)+','+(g|0)+','+(b|0)+')';
    ctx.fillRect(px, py, cell, cell);
  }
}
// mark the two sources
ctx.fillStyle = cv('--good');
ctx.beginPath(); ctx.arc(s1x,s1y,5,0,Math.PI*2); ctx.fill();
ctx.beginPath(); ctx.arc(s2x,s2y,5,0,Math.PI*2); ctx.fill();
ctx.fillStyle = cv('--text');
ctx.font = "13px sans-serif";
ctx.fillText("two coherent sources", 12, 20);
{{< /sketch >}}

## Coherence

The pattern only stays still if the two sources keep a *fixed* phase relationship — they must be **coherent**. Two independent light bulbs never interfere visibly because their phases jitter randomly billions of times a second, smearing the fringes into uniform gray. This is exactly why interference experiments split *one* source in two (Young's double slit, an interferometer) rather than using two separate ones.

Interference in *space* gives fringes; interference in *time* between two slightly different frequencies gives [[Beats]]. Both are the same addition rule, read along different axes.

## See also

- [[Superposition Principle]]
- [[Beats]]
- [[Wave]]
