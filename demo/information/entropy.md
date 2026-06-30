---
title: Entropy
aliases: [shannon entropy, information entropy]
tags: [information]
summary: The average surprise of a random source, measured in bits — and the hard floor on how far it can be compressed.
weight: 10
---

# Entropy

**Entropy** is the average amount of [[Shannon Information|surprise]] you get per symbol from a random source — the expected number of bits needed to name an outcome you do not yet know. For a [[Probability Distribution]] over symbols with probabilities $p_1, p_2, \dots, p_n$, Shannon defined it as

$$H = -\sum_{i} p_i \log_2 p_i \quad \text{bits.}$$

Each term weights a symbol's surprise $-\log_2 p_i$ by how often it occurs. A source that always emits the same symbol has $H = 0$ — no surprise, nothing to learn. A source whose outcomes are spread out has high entropy, because every symbol genuinely tells you something new.

{{< note kind="key" title="Maximized by uniformity" >}}
Among all distributions over $n$ symbols, entropy is **largest when every symbol is equally likely**, giving $H = \log_2 n$ bits. Any imbalance — any predictability — lowers it. Entropy is therefore a precise measure of *how uniform*, or equivalently *how unpredictable*, a distribution is. A fair coin carries 1 bit; a biased coin always carries less.
{{< /note >}}

## The binary entropy curve

For a single biased coin with $P(\text{heads}) = p$, the entropy is

$$H(p) = -p\log_2 p - (1-p)\log_2(1-p).$$

It is zero at the certain extremes $p = 0$ and $p = 1$, and peaks at exactly **1 bit** when $p = \tfrac12$. The curve's gentle dome is one of the most-drawn shapes in the field: a fair coin is maximally informative, and you lose information *slowly* as the coin tilts.

{{< plot fn="-x*Math.log2(x)-(1-x)*Math.log2(1-x)" xmin="0.001" xmax="0.999" ymin="0" ymax="1.05" title="Binary entropy H(p)" caption="The surprise per flip of a biased coin. Maximum 1 bit at p = 0.5; zero when the outcome is certain." >}}

## Feel the average surprise

Drag across the panel below to **pour probability** into four bins. Watch the entropy bar respond: pile everything into one bin and it collapses toward 0; spread it evenly and it climbs to its maximum of $\log_2 4 = 2$ bits.

{{< sketch height="360" caption="Move the mouse left↔right to reshape a 4-symbol distribution; the bar shows its entropy in bits. Flat = maximal, spiky = minimal." >}}
if (frame === 0) {
  state.n = 4;
  state.bias = 0.5; // 0 = uniform, 1 = fully concentrated on bin 0
}
ctx.clearRect(0, 0, W, H);
const css = (v, f) => (getComputedStyle(document.documentElement).getPropertyValue(v).trim() || f);
const accent = css('--accent', '#5b9cff');
const accent2 = css('--accent-2', '#b07bff');
const good = css('--good', '#4ec98f');
const warn = css('--warn', '#e0a458');
const text = css('--text', '#e6e9ef');
const faint = css('--text-faint', 'rgba(230,233,239,0.5)');

// mouse drives a skew parameter in [0,1]
if (mouse.x >= 0 && mouse.x <= W) state.bias = Math.min(1, Math.max(0, mouse.x / W));

// build a distribution: geometric-ish decay controlled by bias
const k = state.bias * 6; // 0..6 decay rate
let p = [];
let Z = 0;
for (let i = 0; i < state.n; i++) { const w = Math.exp(-k * i); p.push(w); Z += w; }
for (let i = 0; i < state.n; i++) p[i] /= Z;

// entropy
let Hh = 0;
for (let i = 0; i < state.n; i++) if (p[i] > 1e-12) Hh += -p[i] * Math.log2(p[i]);
const Hmax = Math.log2(state.n);

// layout
const padL = 40, padR = 30, padT = 28, padB = 46;
const plotW = W - padL - padR;
const barsTop = padT, barsBot = H - padB - 70;
const barsH = barsBot - barsTop;
const bw = plotW / state.n * 0.62;
const gap = plotW / state.n;

// distribution bars
ctx.font = '12px sans-serif';
ctx.fillStyle = faint;
ctx.fillText('distribution (drag mouse)', padL, padT - 10);
for (let i = 0; i < state.n; i++) {
  const x = padL + gap * i + (gap - bw) / 2;
  const h = p[i] * barsH;
  ctx.fillStyle = accent;
  ctx.globalAlpha = 0.85;
  ctx.fillRect(x, barsBot - h, bw, h);
  ctx.globalAlpha = 1;
  ctx.fillStyle = text;
  ctx.fillText(p[i].toFixed(2), x + bw / 2 - 12, barsBot + 16);
  ctx.fillStyle = faint;
  ctx.fillText('s' + (i + 1), x + bw / 2 - 6, barsBot + 32);
}

// entropy meter
const meterTop = barsBot + 50;
const meterH = 22;
ctx.fillStyle = faint;
ctx.fillText('entropy', padL, meterTop - 6);
ctx.strokeStyle = css('--border', 'rgba(255,255,255,0.2)');
ctx.lineWidth = 1;
ctx.strokeRect(padL, meterTop, plotW, meterH);
const frac = Hh / Hmax;
const grad = ctx.createLinearGradient(padL, 0, padL + plotW, 0);
grad.addColorStop(0, good); grad.addColorStop(0.5, accent); grad.addColorStop(1, accent2);
ctx.fillStyle = grad;
ctx.fillRect(padL, meterTop, plotW * frac, meterH);
// max tick
ctx.strokeStyle = warn;
ctx.beginPath(); ctx.moveTo(padL + plotW, meterTop - 4); ctx.lineTo(padL + plotW, meterTop + meterH + 4); ctx.stroke();
ctx.fillStyle = text;
ctx.font = '13px sans-serif';
ctx.fillText('H = ' + Hh.toFixed(3) + ' bits   (max ' + Hmax.toFixed(0) + ')', padL, meterTop + meterH + 22);
{{< /sketch >}}

## Why it is the compression floor

Entropy is not just a tidy formula — it is a *limit*. Shannon's source-coding theorem says no lossless scheme can encode a source in fewer than $H$ bits per symbol on average. Predictable sources (low $H$) compress well; truly random ones (high $H$) cannot be squeezed at all. Every method in [[Data Compression]] is a race to reach this floor, and [[Huffman Coding]] comes within one bit of it.

The same quantity reappears in dynamics. In a [[Chaos|chaotic]] system the [[Lyapunov Exponent]] measures the rate at which stretching reveals ever-finer detail of the initial state — an entropy *production rate*, the speed at which the system manufactures new information.

{{< quiz question="Four symbols are equally likely. What is the entropy of the source?" options="0 bits|1 bit|2 bits|4 bits" answer="3" explain="With n equally likely symbols, H = log2(n). Here log2(4) = 2 bits, the maximum possible for four symbols." >}}

## See also

- [[Shannon Information]]
- [[Data Compression]]
- [[Lyapunov Exponent]]
- [[Probability Distribution]]
