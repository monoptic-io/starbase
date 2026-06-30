---
title: Wave
aliases: [traveling wave, wave motion]
tags: [waves]
summary: A disturbance that travels through a medium while each point of the medium merely oscillates in place — carrying energy without carrying matter.
weight: 10
---

# Wave

A **wave** is a self-propagating disturbance: a pattern that moves through a medium while the medium itself stays put. Drop a pebble in a pond and ripples race outward, but a floating leaf only bobs up and down — it does not travel with the ripple. That is the defining paradox of wave motion. **Energy and shape propagate; matter oscillates in place.** Each particle is just a [[Simple Harmonic Oscillator]], and the wave is the choreography linking one oscillator to the next.

## Anatomy of a wave

A pure traveling sine wave is described by

$$u(x,t)=A\sin\!\left(kx-\omega t\right),$$

and three numbers fix everything about it:

- **Amplitude** $A$ — the peak displacement, setting how much energy the wave carries (energy $\propto A^2$).
- **Wavelength** $\lambda$ — the distance between successive crests; the spatial period. The **wavenumber** is $k=2\pi/\lambda$.
- **Frequency** $f$ — how many crests pass a fixed point per second; the temporal period is $T=1/f$ and the **angular frequency** is $\omega=2\pi f$.

These are tied together by the single most important relation in the subject, linking how often the medium wiggles to how fast the pattern moves:

{{< eq number="1" >}}v = f\lambda = \frac{\omega}{k}{{< /eq >}}

The speed $v$ is usually a fixed property of the *medium* (tension and density for a string, stiffness for sound). So if you raise the frequency, the wavelength must shrink to keep the product constant.

{{< sketch height="320" caption="A traveling sine wave moving left to right. The orange dot is a single particle of the medium — it only moves up and down, tracing simple harmonic motion as the pattern slides past. Move the mouse left/right to change frequency, up/down to change amplitude." >}}
function cv(n){var v=getComputedStyle(document.documentElement).getPropertyValue(n).trim();return v||'#5b9cff';}
if (frame===0){ state.k=4; state.amp=0.7; }
if (mouse.down || (mouse.x>0 && mouse.x<W)){
  state.k = 2 + (mouse.x / W) * 8;
  state.amp = 0.25 + (1 - mouse.y / H) * 0.65;
}
ctx.clearRect(0,0,W,H);
var mid = H/2;
var A = state.amp * H * 0.32;
var k = state.k * Math.PI / W;   // spatial frequency across canvas
var w = 2.2;                     // temporal angular frequency
var phase = w * t;

// midline
ctx.strokeStyle = cv('--border');
ctx.lineWidth = 1;
ctx.beginPath(); ctx.moveTo(0,mid); ctx.lineTo(W,mid); ctx.stroke();

// the wave
ctx.strokeStyle = cv('--accent');
ctx.lineWidth = 2;
ctx.beginPath();
for (var x=0; x<=W; x+=2){
  var y = mid - A*Math.sin(k*x - phase);
  if (x===0) ctx.moveTo(x,y); else ctx.lineTo(x,y);
}
ctx.stroke();

// one tracked particle near the left third
var px = W*0.32;
var py = mid - A*Math.sin(k*px - phase);
// its vertical guide line
ctx.strokeStyle = cv('--text-faint');
ctx.setLineDash([3,4]);
ctx.beginPath(); ctx.moveTo(px,mid); ctx.lineTo(px,py); ctx.stroke();
ctx.setLineDash([]);
ctx.fillStyle = cv('--warn');
ctx.beginPath(); ctx.arc(px,py,5,0,Math.PI*2); ctx.fill();

// wavelength indicator (one period)
var lambda = (2*Math.PI)/k;
ctx.strokeStyle = cv('--accent-2');
ctx.lineWidth = 1.5;
ctx.beginPath(); ctx.moveTo(W*0.1, H-18); ctx.lineTo(W*0.1+lambda, H-18); ctx.stroke();
ctx.fillStyle = cv('--text-dim');
ctx.font = "13px sans-serif";
ctx.fillText("λ", W*0.1 + lambda/2 - 4, H-24);
{{< /sketch >}}

## Transverse and longitudinal

Waves come in two flavors, distinguished by *the direction the medium moves relative to the direction the wave travels*:

- In a **transverse** wave the particles oscillate *perpendicular* to the direction of travel — a wave on a string, light, the ripples above. Crests and troughs.
- In a **longitudinal** wave the particles oscillate *along* the direction of travel — sound in air, a compression pulse down a slinky. Compressions and rarefactions instead of crests and troughs.

Both obey the same mathematics; only the geometry of the displacement differs. Either way, the next step is the equation that governs how the shape propagates: the [[Wave Equation]].

## See also

- [[Wave Equation]]
- [[Simple Harmonic Oscillator]]
- [[Oscillations]]
