---
title: Reaction–Diffusion
aliases: [gray-scott, gray scott, turing pattern, reaction diffusion, reaction-diffusion]
tags: [systems, emergence, pattern-formation]
summary: Two chemicals that react and spread can spontaneously paint stripes, spots, and labyrinths — Turing's mechanism for biological pattern.
weight: 70
---

# Reaction–Diffusion

A **reaction–diffusion** system is a continuous cousin of the [[Cellular Automaton]]: instead of discrete on/off cells, each point on a surface holds the *concentrations* of two or more chemicals that (1) **react** with one another and (2) **diffuse** through space. Locally, the rules are dull chemistry. Globally — across a whole grid of coupled points — they spontaneously organize into spots, stripes, mazes, and pulsing waves. This is the prototypical mechanism of **pattern formation**, and one of the most visually striking examples of emergence in all of science.

## Turing's surprising idea

In 1952 Alan Turing pointed out something counterintuitive. Diffusion, on its own, *smooths* things out — it erases patterns. Yet couple it to the right reaction and it does the opposite: it *creates* them. The trick is two chemicals diffusing at **different speeds**. A slow-spreading "activator" builds local peaks while a fast-spreading "inhibitor" races outward to suppress neighboring peaks. The competition sets a characteristic spacing, and a uniform soup destabilizes into a regular pattern. The transition from "featureless" to "patterned" as a parameter crosses a threshold is a [[Bifurcation]] — a **Turing instability**.

{{< note kind="key" title="The Gray–Scott equations" >}}
Two chemicals $u$ and $v$ react via $u + 2v \to 3v$ (the $uv^2$ term) while $u$ is fed in and $v$ is removed:

$$\dot u = D_u \nabla^2 u - u v^2 + F(1-u),$$
$$\dot v = D_v \nabla^2 v + u v^2 - (F+k)\,v.$$

Here $\nabla^2$ is the **Laplacian** (diffusion), $F$ is the **feed** rate and $k$ the **kill** rate. Tiny changes to $F$ and $k$ switch the system between spots, stripes, and self-replicating blobs.
{{< /note >}}

## A living Gray–Scott grid

The simulation below integrates the Gray–Scott equations on a small grid, a few iterations per frame. It starts from a uniform field seeded with a small disturbance; within seconds the diffusing chemicals carve out growing, dividing, mitosis-like blobs. Nothing in the rule mentions "blob" — the shape is entirely emergent.

{{< sketch height="360" caption="Gray–Scott reaction–diffusion. A seeded patch grows and divides into a coral of self-replicating spots — Turing patterns from two reacting chemicals." >}}
if (frame === 0 || !state.u) {
  const gw = 96, gh = 64;
  state.gw = gw; state.gh = gh;
  state.u = new Float32Array(gw * gh);
  state.v = new Float32Array(gw * gh);
  state.u2 = new Float32Array(gw * gh);
  state.v2 = new Float32Array(gw * gh);
  for (let i = 0; i < gw * gh; i++) { state.u[i] = 1.0; state.v[i] = 0.0; }
  // seed a few patches of v in the middle
  const seed = (cx, cy, r) => {
    for (let y = -r; y <= r; y++) for (let x = -r; x <= r; x++) {
      const gx = cx + x, gy = cy + y;
      if (gx < 0 || gy < 0 || gx >= gw || gy >= gh) continue;
      const idx = gy * gw + gx;
      state.u[idx] = 0.5; state.v[idx] = 0.25 + Math.random() * 0.15;
    }
  };
  seed(gw >> 1, gh >> 1, 4);
  seed((gw >> 1) - 14, (gh >> 1) + 6, 3);
  seed((gw >> 1) + 12, (gh >> 1) - 5, 3);
  // Gray-Scott parameters: "coral / mitosis" regime
  state.Du = 0.16; state.Dv = 0.08; state.F = 0.060; state.k = 0.062;
}
const gw = state.gw, gh = state.gh;
const u = state.u, v = state.v;
let u2 = state.u2, v2 = state.v2;
const Du = state.Du, Dv = state.Dv, F = state.F, kk = state.k;
const ITERS = 8;
for (let it = 0; it < ITERS; it++) {
  for (let y = 0; y < gh; y++) {
    const ym = ((y - 1 + gh) % gh) * gw;
    const yp = ((y + 1) % gh) * gw;
    const y0 = y * gw;
    for (let x = 0; x < gw; x++) {
      const xm = (x - 1 + gw) % gw;
      const xp = (x + 1) % gw;
      const i = y0 + x;
      // Laplacian: orthogonal 0.2, diagonal 0.05, center -1
      const lapU = 0.2*(u[y0+xm] + u[y0+xp] + u[ym+x] + u[yp+x])
                 + 0.05*(u[ym+xm] + u[ym+xp] + u[yp+xm] + u[yp+xp])
                 - u[i];
      const lapV = 0.2*(v[y0+xm] + v[y0+xp] + v[ym+x] + v[yp+x])
                 + 0.05*(v[ym+xm] + v[ym+xp] + v[yp+xm] + v[yp+xp])
                 - v[i];
      const uvv = u[i] * v[i] * v[i];
      let nu = u[i] + (Du*lapU - uvv + F*(1 - u[i]));
      let nv = v[i] + (Dv*lapV + uvv - (F + kk)*v[i]);
      u2[i] = nu < 0 ? 0 : (nu > 1 ? 1 : nu);
      v2[i] = nv < 0 ? 0 : (nv > 1 ? 1 : nv);
    }
  }
  // swap buffers
  for (let i = 0; i < gw*gh; i++) { u[i] = u2[i]; v[i] = v2[i]; }
}
// render
const cw = W / gw, ch = H / gh;
ctx.fillStyle = "#07101a";
ctx.fillRect(0, 0, W, H);
for (let y = 0; y < gh; y++) {
  for (let x = 0; x < gw; x++) {
    const val = v[y * gw + x];
    if (val < 0.05) continue;
    const t = Math.min(1, val * 3.2);
    const hue = 200 - t * 150;            // teal -> magenta as v grows
    const light = 12 + t * 58;
    ctx.fillStyle = "hsl(" + hue + ",85%," + light + "%)";
    ctx.fillRect(x * cw, y * ch, cw + 0.6, ch + 0.6);
  }
}
{{< /sketch >}}

## Patterns in nature

Reaction–diffusion is widely believed to underlie the markings of animals — the spots of a leopard, the stripes of a zebra and angelfish, the spacing of hair follicles and the ridges on a fingertip. The same mathematics describes spreading chemical waves (the Belousov–Zhabotinsky reaction), the branching of corals, and the dunes of a desert. Whenever a short-range "activate" competes with a long-range "inhibit," expect a pattern — a deep link to the rhythms of [[Coupled Oscillators]], where local interaction likewise begets global form.

## See also

- [[Coupled Oscillators]]
- [[Bifurcation]]
- [[Cellular Automaton]]
