---
title: Superposition Principle
aliases: [superposition]
tags: [waves]
summary: When two waves meet, the medium's displacement is simply the sum of what each wave would do alone — the linearity that underlies all of wave physics.
weight: 30
---

# Superposition Principle

The **superposition principle** states that when several waves overlap in the same medium, the total displacement at every point and instant is the **sum** of the displacements each wave would produce on its own:

$$u_{\text{total}}(x,t)=u_1(x,t)+u_2(x,t)+\cdots$$

It sounds almost too obvious to name, but it is the hinge on which the entire subject turns. Superposition holds because the [[Wave Equation]] is **linear**: if $u_1$ and $u_2$ are each solutions, so is any sum $u_1+u_2$. Every richer phenomenon — [[Interference]], [[Beats]], [[Standing Wave|standing waves]], [[Wave Packet|wave packets]] — is just superposition wearing a different costume.

## Waves pass through each other

The most striking consequence: two waves can occupy the same place at the same time and **emerge completely unchanged**. While they overlap they add — sometimes reinforcing, sometimes cancelling — but each carries its own "memory" and continues as if the other had never been there. Two conversations cross a room and neither is garbled. This is utterly unlike particles, which collide and scatter.

{{< sketch height="320" caption="A crest (accent) and an inverted trough (violet) launched toward each other. As they overlap the string briefly flattens — the displacements cancel — yet each pulse re-emerges intact and keeps going. The green curve is the actual string: the sum." >}}
function cv(n){var v=getComputedStyle(document.documentElement).getPropertyValue(n).trim();return v||'#5b9cff';}
ctx.clearRect(0,0,W,H);
var mid = H/2;
var c = W*0.16;
var span = W*1.7;
var sig = W*0.05;
var A = H*0.30;
function gauss(x,c0){var d=x-c0; return Math.exp(-(d*d)/(2*sig*sig));}

var cR = -W*0.35 + ((t*c) % span);          // positive crest, moving right
var cL =  W*1.35 - ((t*c) % span);          // inverted trough, moving left

// midline
ctx.strokeStyle = cv('--border'); ctx.lineWidth=1;
ctx.beginPath(); ctx.moveTo(0,mid); ctx.lineTo(W,mid); ctx.stroke();

// individual pulses (faint)
ctx.globalAlpha = 0.5; ctx.lineWidth=1.2;
ctx.strokeStyle = cv('--accent');
ctx.beginPath();
for (var x=0;x<=W;x+=2){ var y=mid-A*gauss(x,cR); if(x===0)ctx.moveTo(x,y);else ctx.lineTo(x,y);} ctx.stroke();
ctx.strokeStyle = cv('--accent-2');
ctx.beginPath();
for (var x2=0;x2<=W;x2+=2){ var y2=mid+A*gauss(x2,cL); if(x2===0)ctx.moveTo(x2,y2);else ctx.lineTo(x2,y2);} ctx.stroke();
ctx.globalAlpha = 1;

// the sum = the real string
ctx.strokeStyle = cv('--good'); ctx.lineWidth=2.4;
ctx.beginPath();
for (var x3=0;x3<=W;x3+=2){ var y3=mid - A*(gauss(x3,cR) - gauss(x3,cL)); if(x3===0)ctx.moveTo(x3,y3);else ctx.lineTo(x3,y3);} ctx.stroke();

ctx.fillStyle = cv('--text-dim'); ctx.font="13px sans-serif";
ctx.fillText("crest →", 14, 22);
ctx.fillText("← trough", W-90, 22);
{{< /sketch >}}

## The deeper payoff

Superposition is what makes the [[Fourier Series]] possible. If you can build *any* periodic shape by adding sinusoids — and the wave equation lets each sinusoid evolve independently — then you can solve for an arbitrary wave by tracking its simple sinusoidal pieces and adding the results back up. Decompose, evolve, recombine. That strategy, powered entirely by linearity, runs through the whole of [[Fourier Analysis]].

{{< note kind="note" title="When superposition fails" >}}
Real media are only *approximately* linear. Crank the amplitude high enough — a shock wave, a tsunami in shallow water, light in certain crystals — and the medium's response stops being proportional. Superposition breaks, harmonics get generated, and the clean independence of waves is lost. The linear world is a beautiful, usually-excellent approximation.
{{< /note >}}

## See also

- [[Interference]]
- [[Fourier Series]]
- [[Wave Equation]]
