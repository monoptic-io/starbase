---
title: Standing Wave
aliases: [standing waves, normal mode, stationary wave]
tags: [waves]
summary: Trap a wave between two fixed ends and it stops traveling — freezing into a pattern of motionless nodes and violently swinging antinodes.
weight: 40
---

# Standing Wave

A **standing wave** is what a traveling wave becomes when it is **confined**. Tie a string down at both ends, send a wave along it, and the reflections from the two ends superpose with the incoming wave. For most frequencies the result is a mess, but at special **resonant** frequencies the pattern locks into place: certain points, the **nodes**, never move at all, while points halfway between them, the **antinodes**, swing with maximum amplitude. The wave no longer travels — it *stands*.

Mathematically it is two identical [[Wave|traveling waves]] going opposite directions, summed via the [[Superposition Principle]]:

$$\underbrace{A\sin(kx-\omega t)}_{\text{right}}+\underbrace{A\sin(kx+\omega t)}_{\text{left}}=\underbrace{2A\sin(kx)}_{\text{fixed shape}}\;\underbrace{\cos(\omega t)}_{\text{shared wobble}}.$$

Notice that $x$ and $t$ have **separated**. The spatial shape $\sin(kx)$ is frozen; only its overall height breathes in and out as $\cos(\omega t)$. Every point oscillates in lockstep — that is precisely a **normal mode**, the same object you met for [[Coupled Oscillators]], now with infinitely many beads.

## The harmonic ladder

Fixed ends force the string to hold a node at each end, so only whole numbers of half-wavelengths fit: $L=n\lambda/2$. That quantizes the allowed frequencies into a ladder:

{{< eq number="1" >}}f_n=\frac{nv}{2L},\qquad n=1,2,3,\dots{{< /eq >}}

The lowest, $f_1$, is the **fundamental**; the rest are its [[Harmonics]] — exactly the integer-multiple frequencies that a [[Fourier Series]] uses to build any periodic shape. This is why a plucked string sounds like a definite musical pitch: it can only ring at $f_1,2f_1,3f_1,\dots$

{{< sketch height="340" caption="A string vibrating in n loops between two fixed ends. Move the mouse left to right to climb the harmonic ladder n = 1, 2, 3, … The black dots mark the motionless nodes." >}}
function cv(n){var v=getComputedStyle(document.documentElement).getPropertyValue(n).trim();return v||'#5b9cff';}
if (frame===0){ state.n=2; }
if (mouse.x>0 && mouse.x<W){ state.n = 1 + Math.floor((mouse.x / W) * 6); }
var n = Math.max(1, Math.min(7, state.n));
ctx.clearRect(0,0,W,H);
var mid = H/2;
var x0 = W*0.08, x1 = W*0.92, L = x1 - x0;
var A = H*0.30;
var env = Math.cos(2.6 * t);   // the cos(ωt) breathing

// equilibrium line
ctx.strokeStyle = cv('--border'); ctx.lineWidth=1;
ctx.beginPath(); ctx.moveTo(x0,mid); ctx.lineTo(x1,mid); ctx.stroke();

// faint outline of the full envelope (max excursion)
ctx.strokeStyle = cv('--text-faint'); ctx.lineWidth=1; ctx.globalAlpha=0.5;
ctx.beginPath();
for (var xa=x0; xa<=x1; xa+=2){ var ya=mid - A*Math.sin(n*Math.PI*(xa-x0)/L); if(xa===x0)ctx.moveTo(xa,ya);else ctx.lineTo(xa,ya);} ctx.stroke();
ctx.beginPath();
for (var xb=x0; xb<=x1; xb+=2){ var yb=mid + A*Math.sin(n*Math.PI*(xb-x0)/L); if(xb===x0)ctx.moveTo(xb,yb);else ctx.lineTo(xb,yb);} ctx.stroke();
ctx.globalAlpha=1;

// the live string
ctx.strokeStyle = cv('--accent'); ctx.lineWidth=2.6;
ctx.beginPath();
for (var x=x0; x<=x1; x+=2){
  var y = mid - A*env*Math.sin(n*Math.PI*(x-x0)/L);
  if (x===x0) ctx.moveTo(x,y); else ctx.lineTo(x,y);
}
ctx.stroke();

// nodes (including endpoints)
ctx.fillStyle = cv('--text');
for (var k=0; k<=n; k++){
  var nx = x0 + (k/n)*L;
  ctx.beginPath(); ctx.arc(nx, mid, 4, 0, Math.PI*2); ctx.fill();
}
// endpoints clamps
ctx.fillStyle = cv('--warn');
ctx.fillRect(x0-5, mid-12, 4, 24);
ctx.fillRect(x1+1, mid-12, 4, 24);

ctx.fillStyle = cv('--text-dim'); ctx.font="14px sans-serif";
ctx.fillText("n = " + n + "   (f" + n + " = " + n + "·v/2L)", x0, H-14);
{{< /sketch >}}

## Resonance and instruments

You can only *excite* a standing wave by driving the string near one of its $f_n$ — drive off-resonance and the reflections fight the input. This is [[Resonance]] in spatial form: the boundary conditions pick out the natural frequencies, and energy pours in only when you match one. Strings, organ pipes, drumheads, and the cavity of a laser all work this way; the shape of the boundary chooses which [[Harmonics]] are allowed to ring.

## See also

- [[Harmonics]]
- [[Resonance]]
- [[Coupled Oscillators]]
