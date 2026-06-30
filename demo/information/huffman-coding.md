---
title: Huffman Coding
aliases: [huffman, prefix code]
tags: [information]
summary: A greedy algorithm that builds the optimal prefix code — frequent symbols get short codewords, rare ones long.
weight: 40
---

# Huffman Coding

**Huffman coding** is the classic algorithm for building an *optimal* symbol code: given how often each symbol appears, it produces the variable-length binary code with the shortest possible average length. Frequent symbols (like the letter *e* in English) get **short** codewords; rare symbols (*z*, *q*) get **long** ones. David Huffman invented it in 1952 as a term-paper assignment — and it has been quietly compressing the world's data ever since, inside ZIP, JPEG, MP3, and more.

The codes it produces are **prefix codes**: no codeword is a prefix of another. That property lets a decoder read a bitstream left to right and split it into symbols unambiguously, with no separators.

## The greedy rule

The algorithm is disarmingly simple and provably optimal:

1. Start with one **leaf node** per symbol, each holding its frequency.
2. Repeatedly take the **two least-frequent** nodes and merge them under a new parent whose frequency is their sum.
3. Stop when a single node — the **root** — remains.
4. Read codewords off the tree: label each left branch `0` and each right branch `1`, then spell out the path from root to each leaf.

The least-frequent symbols get merged first, so they end up **deepest** in the tree and earn the **longest** codewords. That is exactly what we want.

## Watch the tree build itself

The sketch below runs the algorithm live. Leaves sit at the bottom with their frequencies. Step by step, the two lowest-frequency nodes are pulled together and joined under a new parent (the running merge is highlighted). When the root is reached, the finished codewords appear — short for the common symbols, long for the rare ones.

{{< sketch height="440" caption="Huffman construction on six symbols. Each step merges the two least-frequent nodes; the resulting prefix codewords are read off as left=0 / right=1. The animation loops." >}}
if (frame === 0) {
  // --- symbols and frequencies ---
  const syms = [['A',40],['B',20],['C',15],['D',10],['E',8],['F',7]];
  let nodes = [];
  let idc = 0;
  for (const s of syms) nodes.push({id: idc++, sym: s[0], f: s[1], l: null, r: null, born: -1});
  state.all = nodes.slice();
  // --- build Huffman tree, recording merge order ---
  let work = nodes.slice();
  let step = 0;
  while (work.length > 1) {
    // pick two smallest by frequency (ties: lower id first) for determinism
    work.sort((a,b) => a.f - b.f || a.id - b.id);
    const a = work.shift();
    const b = work.shift();
    const p = {id: idc++, sym: null, f: a.f + b.f, l: a, r: b, born: step++};
    state.all.push(p);
    work.push(p);
  }
  state.root = work[0];
  state.steps = step; // number of merges

  // --- assign codewords ---
  state.codes = {};
  (function walk(n, code) {
    if (!n) return;
    if (n.sym !== null) { state.codes[n.sym] = code || '0'; return; }
    walk(n.l, code + '0');
    walk(n.r, code + '1');
  })(state.root, '');

  // --- layout: leaves by in-order index, depth from root ---
  let order = 0, maxDepth = 0;
  (function place(n, depth) {
    if (!n) return;
    n.depth = depth;
    if (depth > maxDepth) maxDepth = depth;
    if (n.sym !== null) { n.ox = order++; return; }
    place(n.l, depth + 1);
    place(n.r, depth + 1);
    n.ox = (n.l.ox + n.r.ox) / 2;
  })(state.root, 0);
  state.leafCount = order;
  state.maxDepth = maxDepth;
}

ctx.clearRect(0, 0, W, H);
const css = (v, f) => (getComputedStyle(document.documentElement).getPropertyValue(v).trim() || f);
const accent = css('--accent', '#5b9cff');
const accent2 = css('--accent-2', '#b07bff');
const good = css('--good', '#4ec98f');
const warn = css('--warn', '#e0a458');
const text = css('--text', '#e6e9ef');
const faint = css('--text-faint', 'rgba(230,233,239,0.55)');
const border = css('--border', 'rgba(255,255,255,0.2)');

// timeline: reveal one merge per beat, then hold, then loop
const beat = 0.9;
const hold = 3.0;
const total = state.steps * beat + hold;
const tt = t % total;
const shown = Math.min(state.steps, Math.floor(tt / beat) + 1);
const done = tt >= state.steps * beat;

// layout transform
const padX = 46, padTop = 36, padBot = 56;
const plotW = W - padX * 2;
const plotH = H - padTop - padBot;
const X = (ox) => padX + (state.leafCount <= 1 ? 0.5 : ox / (state.leafCount - 1)) * plotW;
const Y = (d) => padTop + (state.maxDepth === 0 ? 0 : d / state.maxDepth) * plotH;

// which nodes/edges are visible
const visible = (n) => n.born < shown; // leaves have born -1 (always), merges by step

// draw edges first
for (const n of state.all) {
  if (n.sym !== null || !visible(n)) continue;
  const isCur = (n.born === shown - 1) && !done;
  for (const c of [n.l, n.r]) {
    ctx.strokeStyle = isCur ? warn : border;
    ctx.lineWidth = isCur ? 2.4 : 1.4;
    ctx.globalAlpha = isCur ? 1 : 0.8;
    ctx.beginPath();
    ctx.moveTo(X(n.ox), Y(n.depth));
    ctx.lineTo(X(c.ox), Y(c.depth));
    ctx.stroke();
    // bit label
    const mx = (X(n.ox) + X(c.ox)) / 2, my = (Y(n.depth) + Y(c.depth)) / 2;
    ctx.globalAlpha = 1;
    ctx.fillStyle = faint;
    ctx.font = '11px sans-serif';
    ctx.fillText(c === n.l ? '0' : '1', mx + 4, my);
  }
}
ctx.globalAlpha = 1;

// draw nodes
ctx.font = '12px sans-serif';
ctx.textAlign = 'center';
ctx.textBaseline = 'middle';
for (const n of state.all) {
  if (n.sym === null && !visible(n)) continue;
  const x = X(n.ox), y = Y(n.depth);
  const leaf = n.sym !== null;
  const isCur = !leaf && n.born === shown - 1 && !done;
  const r = leaf ? 16 : 13;
  ctx.beginPath(); ctx.arc(x, y, r, 0, 2 * Math.PI);
  ctx.fillStyle = leaf ? accent : (isCur ? warn : 'rgba(120,130,160,0.55)');
  ctx.globalAlpha = leaf ? 0.95 : 0.9;
  ctx.fill();
  ctx.globalAlpha = 1;
  ctx.lineWidth = isCur ? 2.4 : 1;
  ctx.strokeStyle = isCur ? warn : border;
  ctx.stroke();
  ctx.fillStyle = leaf ? '#0b0e16' : text;
  if (leaf) {
    ctx.fillText(n.sym, x, y - 4);
    ctx.font = '10px sans-serif';
    ctx.fillText(n.f + '', x, y + 7);
    ctx.font = '12px sans-serif';
  } else {
    ctx.fillText(n.f + '', x, y);
  }
}
ctx.textAlign = 'left';
ctx.textBaseline = 'alphabetic';

// status line
ctx.fillStyle = faint;
ctx.font = '12px sans-serif';
if (!done) ctx.fillText('merging the two least-frequent nodes…  step ' + shown + ' / ' + state.steps, padX, 20);
else ctx.fillText('tree complete — read codewords as root→leaf paths (0 = left, 1 = right)', padX, 20);

// codeword table once done
if (done) {
  const keys = Object.keys(state.codes).sort();
  let cx = padX, cy = H - 30;
  ctx.font = '13px monospace';
  for (const kkey of keys) {
    const label = kkey + ':' + state.codes[kkey];
    ctx.fillStyle = good;
    ctx.fillText(label, cx, cy);
    cx += ctx.measureText(label).width + 18;
  }
}
{{< /sketch >}}

## Why "greedy" lands on the optimum

It feels too easy that always-merge-the-two-smallest could be optimal — but it is, and the proof rests on one observation. In *any* optimal prefix code, the two least-frequent symbols can be taken to sit at the deepest level as siblings. Merging them first is therefore never a mistake; you can keep doing it, and induction carries the optimality all the way to the root. The result comes within **less than one bit per symbol** of the [[Entropy]] floor — and arithmetic coding closes even that small gap.

{{< note kind="tip" title="Where Huffman lives" >}}
Pure Huffman is still the back end of **DEFLATE** (ZIP, gzip, PNG) and the entropy stage of **JPEG**. It shines when symbol frequencies are stable and known. When they drift, adaptive and arithmetic coders take over — but the core idea, *cheap codes for common symbols*, is Huffman's.
{{< /note >}}

{{< quiz question="In a Huffman tree, which symbol receives the longest codeword?" options="The most frequent one|The least frequent one|The alphabetically first one|They are all the same length" answer="2" explain="Least-frequent symbols are merged earliest, so they sit deepest in the tree and get the longest root-to-leaf paths — exactly the codewords we can most afford to make long." >}}

## See also

- [[Data Compression]]
- [[Entropy]]
- [[Shannon Information]]
