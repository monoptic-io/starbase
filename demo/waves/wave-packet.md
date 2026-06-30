---
title: Wave Packet
aliases: [wave packets, group velocity]
tags: [waves]
summary: A localized burst of oscillation — a carrier wave wrapped in a traveling envelope — whose envelope and crests move at different speeds in a dispersive medium.
weight: 90
---

# Wave Packet

A **wave packet** is a localized lump of wave: a fast oscillating **carrier** wrapped inside a smooth, slowly varying **envelope**. Unlike an endless sine, a packet is concentrated in a region of space, which is what makes it physical — real signals, pulses of light, and quantum particles are all wave packets. It is the bridge between a wave (spread everywhere) and a particle (somewhere).

You build one by adding many sinusoids of nearby wavenumbers via the [[Superposition Principle]]. Away from the center they interfere destructively and cancel; near the center they reinforce. A Gaussian envelope is the cleanest example:

$$u(x,t)=\underbrace{e^{-\frac{(x-v_g t)^2}{2\sigma^2}}}_{\text{envelope, moves at }v_g}\;\underbrace{\cos\!\big(k_0 x-\omega_0 t\big)}_{\text{carrier, crests move at }v_p}.$$

## Two speeds at once

The defining feature of a packet in a [[Dispersion|dispersive]] medium is that **the envelope and the carrier move at different speeds**. The envelope — the bundle of energy and information — travels at the **group velocity** $v_g=d\omega/dk$. The individual crests inside travel at the **phase velocity** $v_p=\omega/k$. Watch closely and you can see crests born at the trailing edge of the packet, sweeping forward through the envelope, and dissolving at the leading edge.

{{< sketch height="360" caption="A Gaussian wave packet. The faint outline is the envelope, gliding at the group velocity; the orange oscillation inside is the carrier, whose crests slide through the envelope at a different (here faster) phase velocity. Watch a crest appear at the back, march forward, and vanish at the front." >}}
function cv(n){var v=getComputedStyle(document.documentElement).getPropertyValue(n).trim();return v||'#5b9cff';}
ctx.clearRect(0,0,W,H);
var mid = H/2;
var A = H*0.34;
var sig = W*0.10;
var vg = W*0.10;          // group velocity (envelope)
var vp = W*0.24;          // phase velocity (crests) — faster, so crests slide forward
var k0 = 0.14;            // carrier wavenumber (px^-1)
var span = W*1.5;
var cx = -W*0.25 + ((t*vg) % span);   // envelope center

function envelope(x){ var d=x-cx; return Math.exp(-(d*d)/(2*sig*sig)); }

// midline
ctx.strokeStyle = cv('--border'); ctx.lineWidth=1;
ctx.beginPath(); ctx.moveTo(0,mid); ctx.lineTo(W,mid); ctx.stroke();

// envelope outline (top + bottom), faint
ctx.strokeStyle = cv('--accent-2'); ctx.globalAlpha=0.55; ctx.lineWidth=1.5;
ctx.beginPath();
for (var x=0;x<=W;x+=2){ var y=mid-A*envelope(x); if(x===0)ctx.moveTo(x,y);else ctx.lineTo(x,y);} ctx.stroke();
ctx.beginPath();
for (var xb=0;xb<=W;xb+=2){ var yb=mid+A*envelope(xb); if(xb===0)ctx.moveTo(xb,yb);else ctx.lineTo(xb,yb);} ctx.stroke();
ctx.globalAlpha=1;

// the packet = envelope * carrier
ctx.strokeStyle = cv('--accent'); ctx.lineWidth=2.2;
ctx.beginPath();
for (var x2=0;x2<=W;x2+=1.5){
  var carrier = Math.cos(k0*x2 - (k0*vp)*t);
  var y2 = mid - A*envelope(x2)*carrier;
  if (x2===0) ctx.moveTo(x2,y2); else ctx.lineTo(x2,y2);
}
ctx.stroke();

// markers: envelope center (group) and a tracked crest
ctx.fillStyle = cv('--good');
ctx.beginPath(); ctx.arc(cx, mid+A*envelope(cx)+10, 4, 0, Math.PI*2); ctx.fill();
ctx.fillStyle = cv('--text-dim'); ctx.font="13px sans-serif";
ctx.fillText("envelope → v_group", 12, 22);
ctx.fillText("crests → v_phase", 12, 40);
{{< /sketch >}}

## Localization costs bandwidth

A packet narrow in space needs a *wide* spread of wavenumbers to build it, and a packet narrow in wavenumber must be *broad* in space. This reciprocal trade — $\Delta x\,\Delta k \gtrsim \tfrac12$ — is a pure property of the [[Fourier Transform]], and in quantum mechanics, where $p=\hbar k$, it becomes the [[Uncertainty Principle]]. You cannot have a wave that is both perfectly located and perfectly monochromatic.

In a non-dispersive medium the packet would glide along rigidly; with dispersion the constituent frequencies drift apart and the packet **spreads** as it travels — the same broadening that limits how fast pulses can be sent down an optical fiber.

## See also

- [[Dispersion]]
- [[Uncertainty Principle]]
- [[Fourier Transform]]
