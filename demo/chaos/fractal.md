---
title: Fractal
aliases: [fractal, fractals, self-similarity]
tags: [chaos]
summary: A shape that repeats its structure at every scale, with a fractional dimension that measures how intricately it fills space.
weight: 70
---

# Fractal

A **fractal** is a shape that looks the same — or statistically similar — no matter how far you zoom in. Coastlines, ferns, lungs, lightning, and the cross-section of every [[Strange Attractor]] share this property of **self-similarity**: the whole is echoed in its parts, down to arbitrarily small scales. Fractals are the natural geometry of [[Chaos|chaos]], because the stretch-and-fold dynamics that generate chaos lay down structure at every scale at once.

## Self-similarity and infinite detail

Zoom into a smooth curve and it eventually looks like a straight line — it is *locally simple*. Zoom into a fractal and the complexity never resolves; new detail keeps appearing. The classic constructions make this concrete: the Koch curve adds a bump to every segment forever, becoming infinitely long while enclosing a finite area; the Cantor set repeatedly deletes middle thirds, leaving infinitely many points of total length zero.

The **Julia set** below is the showpiece. Each point of the plane is colored by how quickly the iteration $z \mapsto z^2 + c$ runs away to infinity; the black region is the set of points that never escape. Its boundary is an infinitely detailed fractal whose shape is dictated entirely by the constant $c$.

**Click and drag** on the image to change $c$ and watch the fractal morph — from connected filaments to exploding dust — in real time.

{{< sketch height="460" caption="Filled Julia set of z → z² + c. Drag anywhere to steer c by mouse position; release to keep it. The black interior never escapes; the glowing boundary is fractal at every scale." >}}
if (frame === 0) {
  // small offscreen buffer for speed; scaled up to fill the canvas
  state.bw = 280;
  state.bh = Math.max(1, Math.round(280 * H / W));
  state.buf = document.createElement('canvas');
  state.buf.width = state.bw;
  state.buf.height = state.bh;
  state.bctx = state.buf.getContext('2d');
  state.img = state.bctx.createImageData(state.bw, state.bh);
  state.c = { re: -0.8, im: 0.156 };
  state.lastc = { re: 999, im: 999 };
  state.maxIter = 90;
}

function renderJulia(cre, cim) {
  const bw = state.bw, bh = state.bh, maxIter = state.maxIter;
  const data = state.img.data;
  const spanX = 3.2, spanY = spanX * bh / bw; // view [-1.6,1.6] x scaled
  let p = 0;
  for (let py = 0; py < bh; py++) {
    const y0 = (py / (bh - 1) - 0.5) * spanY;
    for (let px = 0; px < bw; px++) {
      const x0 = (px / (bw - 1) - 0.5) * spanX;
      let zx = x0, zy = y0, n = 0;
      while (n < maxIter) {
        const x2 = zx * zx, y2 = zy * zy;
        if (x2 + y2 > 16) break;
        const xt = x2 - y2 + cre;
        zy = 2 * zx * zy + cim;
        zx = xt;
        n++;
      }
      if (n >= maxIter) {
        data[p] = 6; data[p + 1] = 8; data[p + 2] = 16; data[p + 3] = 255;
      } else {
        // smooth escape count
        const mag = Math.sqrt(zx * zx + zy * zy);
        const mu = (n + 1 - Math.log(Math.log(mag)) / Math.log(2)) / maxIter;
        const t = Math.max(0, Math.min(1, mu));
        // warm/cool polynomial palette
        const r = 9 * (1 - t) * t * t * t * 255;
        const g = 15 * (1 - t) * (1 - t) * t * t * 255;
        const b = 8.5 * (1 - t) * (1 - t) * (1 - t) * t * 255;
        data[p] = r; data[p + 1] = g; data[p + 2] = b; data[p + 3] = 255;
      }
      p += 4;
    }
  }
  state.bctx.putImageData(state.img, 0, 0);
}

// steer c with the mouse while dragging
if (mouse.down && mouse.x >= 0 && mouse.x <= W && mouse.y >= 0 && mouse.y <= H) {
  state.c.re = (mouse.x / W - 0.5) * 2.0;   // re in [-1, 1]
  state.c.im = (mouse.y / H - 0.5) * 1.6;   // im in [-0.8, 0.8]
}
if (state.c.re !== state.lastc.re || state.c.im !== state.lastc.im) {
  renderJulia(state.c.re, state.c.im);
  state.lastc.re = state.c.re;
  state.lastc.im = state.c.im;
}

// blit (nearest-neighbour scale-up keeps it crisp and cheap)
ctx.imageSmoothingEnabled = true;
ctx.clearRect(0, 0, W, H);
ctx.drawImage(state.buf, 0, 0, W, H);

ctx.fillStyle = 'rgba(255,255,255,0.9)';
ctx.font = '13px sans-serif';
ctx.fillText('c = ' + state.c.re.toFixed(3) + ' + ' + state.c.im.toFixed(3) + ' i', 10, 20);
ctx.fillStyle = 'rgba(255,255,255,0.55)';
ctx.fillText('drag to change c', 10, H - 10);
{{< /sketch >}}

## Fractal dimension

How do you measure something rougher than a curve but thinner than a region? With a **dimension that need not be a whole number**. The box-counting dimension asks how the number $N(\varepsilon)$ of little boxes of size $\varepsilon$ needed to cover the set grows as the boxes shrink:

{{< eq >}}D = \lim_{\varepsilon \to 0}\frac{\ln N(\varepsilon)}{\ln(1/\varepsilon)}{{< /eq >}}

For a smooth line you get $D = 1$, for a filled square $D = 2$. The Koch curve gives $D = \ln 4/\ln 3 \approx 1.26$ — genuinely between a line and an area. The cross-section of the [[Lorenz System|Lorenz attractor]] sits at $D \approx 2.06$. A fractional dimension is precisely the quantitative statement that detail persists across scales.

{{< note kind="note" title="Strange attractors are fractals" >}}
This is the bridge back to dynamics. The repeated stretching and folding that powers a [[Strange Attractor]] stacks infinitely many layers into a bounded region — a fractal in [[State Space]]. Its (usually non-integer) dimension is one of the cleanest fingerprints distinguishing chaotic motion from mere noise.
{{< /note >}}

{{< quiz question="The Koch curve has box-counting dimension about 1.26. What does a non-integer dimension between 1 and 2 tell you?" options="The measurement was done incorrectly|The set is more space-filling than a smooth line but does not fill an area — detail persists at every scale|The curve is two-dimensional|The curve has finite length" answer="2" explain="A dimension strictly between 1 and 2 quantifies a shape rougher and more space-filling than a 1D curve yet thinner than a 2D region, the hallmark of self-similar fractal detail at all scales." >}}

## See also

- [[Strange Attractor]]
- [[Lorenz System]]
- [[State Space]]
