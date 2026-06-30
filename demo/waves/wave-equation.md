---
title: Wave Equation
aliases: [wave PDE]
tags: [waves]
summary: The partial differential equation whose solutions are anything that travels left or right at a fixed speed without changing shape.
weight: 20
---

# Wave Equation

The **wave equation** is the law every non-dispersive wave obeys:

$$\frac{\partial^2 u}{\partial t^2}=c^2\frac{\partial^2 u}{\partial x^2}.$$

Read it as a statement about curvature: the *acceleration* of the medium at a point ($\partial_{tt}u$) is proportional to how sharply the medium is *bent* there ($\partial_{xx}u$). A taut string that is curved upward gets pulled downward by its own tension, and the constant of proportionality $c^2$ sets how fast the resulting motion races along. The speed $c=\sqrt{T/\mu}$ for a string of tension $T$ and linear density $\mu$.

## From masses on springs to a continuum

The wave equation is not fundamental — it *emerges*. Take a row of [[Coupled Oscillators]]: beads of mass $m$ on a string, each tied to its neighbors by springs. Newton's law for bead $n$ is

$$m\,\ddot u_n = \kappa\,(u_{n+1}-2u_n+u_{n-1}).$$

The right-hand side is a discrete second difference — the lattice's version of a second derivative. Now let the beads shrink and crowd together, taking the spacing to zero. The second difference becomes $\partial_{xx}u$, and the discrete chain melts into the smooth wave equation. **A wave is the continuum limit of infinitely many coupled oscillators**, which is exactly why every point still moves like a [[Simple Harmonic Oscillator]].

## d'Alembert: left- and right-movers

The general solution, found by d'Alembert, is breathtakingly simple:

$$u(x,t)=f(x-ct)+g(x+ct).$$

*Any* function $f$ of the combination $x-ct$ is a shape that slides rightward at speed $c$ without distorting; any $g(x+ct)$ slides leftward. The full solution is just a **right-mover plus a left-mover**, set by the initial shape and velocity. Nothing about the profile matters — a smooth bump, a sharp kink, a sine — it simply translates.

{{< sketch height="320" caption="Two pulses launched from opposite ends: a right-mover (accent) and a left-mover (violet). Each keeps its shape exactly as d'Alembert promises. Watch them overlap in the middle and emerge unchanged." >}}
function cv(n){var v=getComputedStyle(document.documentElement).getPropertyValue(n).trim();return v||'#5b9cff';}
ctx.clearRect(0,0,W,H);
var mid = H*0.6;
var c = W*0.16;            // wave speed (px/s)
var span = W*1.6;          // travel loop length
var sig = W*0.045;         // pulse width
var A = H*0.28;

function gauss(x,c0){var d=x-c0; return Math.exp(-(d*d)/(2*sig*sig));}

// right-mover starts left of screen, left-mover starts right of screen
var cR = -W*0.3 + ((t*c) % span);
var cL =  W*1.3 - ((t*c) % span);

// midline
ctx.strokeStyle = cv('--border');
ctx.lineWidth = 1;
ctx.beginPath(); ctx.moveTo(0,mid); ctx.lineTo(W,mid); ctx.stroke();

// individual movers, faint
ctx.lineWidth = 1.2;
ctx.strokeStyle = cv('--accent-soft') || 'rgba(91,156,255,0.3)';
ctx.globalAlpha = 0.45;
ctx.strokeStyle = cv('--accent');
ctx.beginPath();
for (var x=0;x<=W;x+=2){ var y=mid-A*gauss(x,cR); if(x===0)ctx.moveTo(x,y);else ctx.lineTo(x,y);} ctx.stroke();
ctx.strokeStyle = cv('--accent-2');
ctx.beginPath();
for (var x2=0;x2<=W;x2+=2){ var y2=mid-A*gauss(x2,cL); if(x2===0)ctx.moveTo(x2,y2);else ctx.lineTo(x2,y2);} ctx.stroke();
ctx.globalAlpha = 1;

// the actual string = sum of the two
ctx.strokeStyle = cv('--good');
ctx.lineWidth = 2.4;
ctx.beginPath();
for (var x3=0;x3<=W;x3+=2){ var y3=mid-A*(gauss(x3,cR)+gauss(x3,cL)); if(x3===0)ctx.moveTo(x3,y3);else ctx.lineTo(x3,y3);} ctx.stroke();

ctx.fillStyle = cv('--text-dim');
ctx.font = "13px sans-serif";
ctx.fillText("f(x−ct) →", 14, 24);
ctx.fillText("← g(x+ct)", W-104, 24);
{{< /sketch >}}

Where the two pulses overlap they simply add (the [[Superposition Principle]]) and then continue on, perfectly intact. If $c$ instead depended on frequency, the shape would *not* survive — that is the story of [[Dispersion]]. And if the ends are tied down so the movers reflect and recombine forever, you get a [[Standing Wave]].

## See also

- [[Coupled Oscillators]]
- [[Standing Wave]]
- [[Dispersion]]
