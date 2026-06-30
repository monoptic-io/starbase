---
title: Finite Automaton
aliases: [finite state machine, dfa, fsm, finite automaton]
tags: [computation]
summary: A memoryless machine of states and transitions that recognizes exactly the regular languages.
weight: 20
---

# Finite Automaton

A **finite automaton** (or finite-state machine) is the [[Turing Machine|Turing machine]] with its memory amputated. It has a finite set of **states**, a designated **start** state, some **accepting** states, and a **transition** rule that, on reading each input symbol, jumps from one state to another. There is no tape to write on — the machine's *entire* recollection of everything it has read is the single state it currently sits in.

That sounds crippling, and in a sense it is: a finite automaton can only recognize **regular languages** — patterns like "ends in `.com`", "has an even number of `1`s", or "is a valid phone number". But that modest power is everywhere: every regex engine, lexer, vending machine, traffic light, and network protocol is a finite automaton in disguise.

## Watch one decide

The machine below recognizes binary strings whose value is a **multiple of three**. It needs just three states — one for each possible remainder mod 3 — because the running remainder is the *only* thing it must remember. On each bit `b` it updates $r \leftarrow (2r + b) \bmod 3$. If it ends in state $q_0$ (remainder 0), the string is accepted.

{{< sketch height="380" caption="A 3-state DFA testing whether a binary string is divisible by 3. The active state glows; the firing transition lights up as each bit is read. Green halo = accept (ended in q0), red = reject. A fresh random string runs each pass." >}}
if (frame === 0) {
  state.fresh = function () {
    const len = 5 + Math.floor(Math.random() * 4);
    let s = '';
    for (let i = 0; i < len; i++) s += Math.random() < 0.5 ? '0' : '1';
    state.input = s;
    state.pos = 0;
    state.cur = 0;
    state.prev = 0;
    state.bit = -1;
    state.phase = 'run';
  };
  state.fresh();
  state.tNext = 1.0;
}
const cs = getComputedStyle(document.documentElement);
const accent = cs.getPropertyValue('--accent').trim() || '#5b9cff';
const good = cs.getPropertyValue('--good').trim() || '#46d39a';
const warn = cs.getPropertyValue('--warn').trim() || '#f2b347';
const text = cs.getPropertyValue('--text').trim() || '#e8edf5';
const border = cs.getPropertyValue('--border').trim() || 'rgba(255,255,255,0.18)';
const faint = cs.getPropertyValue('--text-faint').trim() || 'rgba(255,255,255,0.4)';

// advance on a timer
if (t > state.tNext) {
  if (state.phase === 'run') {
    if (state.pos < state.input.length) {
      const b = +state.input[state.pos];
      state.prev = state.cur;
      state.cur = (2 * state.cur + b) % 3;
      state.bit = b;
      state.pos += 1;
      state.tNext = t + 0.85;
    } else {
      state.phase = state.cur === 0 ? 'accept' : 'reject';
      state.tNext = t + 1.8;
    }
  } else {
    state.fresh();
    state.tNext = t + 0.85;
  }
}

ctx.clearRect(0, 0, W, H);
const cx = W / 2, cy = H * 0.6;
const R = Math.min(W, H) * 0.30;
const nodeR = Math.min(30, W * 0.09);
const nodes = [];
for (let i = 0; i < 3; i++) {
  const a = -Math.PI / 2 + i * 2 * Math.PI / 3;
  nodes.push({ x: cx + R * Math.cos(a), y: cy + R * Math.sin(a) });
}

function edge(ai, bi, active, label) {
  const A = nodes[ai], B = nodes[bi];
  let dx = B.x - A.x, dy = B.y - A.y;
  const d = Math.hypot(dx, dy); dx /= d; dy /= d;
  const px = -dy, py = dx, off = 9;
  const sx = A.x + dx * nodeR + px * off, sy = A.y + dy * nodeR + py * off;
  const ex = B.x - dx * nodeR + px * off, ey = B.y - dy * nodeR + py * off;
  ctx.strokeStyle = active ? accent : border;
  ctx.lineWidth = active ? 3 : 1.5;
  ctx.beginPath(); ctx.moveTo(sx, sy); ctx.lineTo(ex, ey); ctx.stroke();
  const ah = 8;
  ctx.fillStyle = active ? accent : border;
  ctx.beginPath();
  ctx.moveTo(ex, ey);
  ctx.lineTo(ex - dx * ah - px * ah * 0.5, ey - dy * ah - py * ah * 0.5);
  ctx.lineTo(ex - dx * ah + px * ah * 0.5, ey - dy * ah + py * ah * 0.5);
  ctx.closePath(); ctx.fill();
  ctx.fillStyle = active ? accent : faint;
  ctx.font = '13px ui-monospace, monospace';
  ctx.textAlign = 'center';
  ctx.fillText(label, (sx + ex) / 2 + px * 11, (sy + ey) / 2 + py * 11);
}

function selfLoop(ai, active, label) {
  const A = nodes[ai];
  const ly = A.y - nodeR - 13;
  ctx.strokeStyle = active ? accent : border;
  ctx.lineWidth = active ? 3 : 1.5;
  ctx.beginPath(); ctx.arc(A.x, ly, 12, 0, Math.PI * 2); ctx.stroke();
  ctx.fillStyle = active ? accent : faint;
  ctx.font = '13px ui-monospace, monospace';
  ctx.textAlign = 'center';
  ctx.fillText(label, A.x, ly - 17);
}

const stepped = state.pos > 0 && state.phase === 'run';
function on(from, to) { return stepped && state.prev === from && state.cur === to; }

// transitions: q0--1-->q1, q1--1-->q0, q1--0-->q2, q2--0-->q1; self: q0--0, q2--1
edge(0, 1, on(0, 1), '1');
edge(1, 0, on(1, 0), '1');
edge(1, 2, on(1, 2), '0');
edge(2, 1, on(2, 1), '0');
selfLoop(0, on(0, 0), '0');
selfLoop(2, on(2, 2), '1');

// nodes
for (let i = 0; i < 3; i++) {
  const A = nodes[i];
  const isCur = i === state.cur;
  let glow = border;
  if (isCur && state.phase === 'accept') glow = good;
  else if (isCur && state.phase === 'reject') glow = warn;
  else if (isCur) glow = accent;
  if (isCur) {
    ctx.beginPath(); ctx.arc(A.x, A.y, nodeR + 6, 0, Math.PI * 2);
    ctx.fillStyle = glow; ctx.globalAlpha = 0.18; ctx.fill(); ctx.globalAlpha = 1;
  }
  ctx.beginPath(); ctx.arc(A.x, A.y, nodeR, 0, Math.PI * 2);
  ctx.fillStyle = '#0c1018'; ctx.fill();
  ctx.lineWidth = 2.5; ctx.strokeStyle = glow; ctx.stroke();
  if (i === 0) {  // accepting state: double ring
    ctx.beginPath(); ctx.arc(A.x, A.y, nodeR - 4, 0, Math.PI * 2);
    ctx.lineWidth = 1.5; ctx.strokeStyle = glow; ctx.stroke();
  }
  ctx.fillStyle = text;
  ctx.font = '15px ui-monospace, monospace';
  ctx.textAlign = 'center'; ctx.textBaseline = 'middle';
  ctx.fillText('q' + i, A.x, A.y);
}

// input string with read cursor
ctx.textBaseline = 'alphabetic';
ctx.font = '20px ui-monospace, monospace';
const chw = 18;
const ix0 = W / 2 - (state.input.length * chw) / 2;
const iy = H * 0.13;
for (let i = 0; i < state.input.length; i++) {
  const read = i < state.pos;
  ctx.fillStyle = (i === state.pos - 1) ? accent : (read ? faint : text);
  ctx.fillText(state.input[i], ix0 + i * chw + chw / 2, iy);
}
ctx.font = '13px ui-monospace, monospace';
ctx.fillStyle = state.phase === 'accept' ? good : state.phase === 'reject' ? warn : faint;
const msg = state.phase === 'accept' ? 'ACCEPT — divisible by 3'
          : state.phase === 'reject' ? 'REJECT — not divisible by 3'
          : 'reading...';
ctx.fillText(msg, W / 2, iy + 26);
{{< /sketch >}}

Three states, six transitions, and the machine flawlessly tests divisibility on strings of any length — using a fixed amount of memory no matter how long the input grows. That bounded memory is the signature of a finite automaton.

## What it can and cannot do

The reach of finite automata is captured by a beautiful equivalence: a language is recognizable by *some* finite automaton **if and only if** it can be described by a regular expression. Deterministic (DFA) and nondeterministic (NFA) variants look different but accept exactly the same languages — every NFA can be flattened into a DFA.

{{< note kind="warning" title="The wall: counting" >}}
A finite automaton **cannot** recognize the language of balanced parentheses, or "$n$ zeros followed by $n$ ones". To check that two counts match, you must remember a number that can grow without bound — and a finite state set cannot. The moment a problem needs unbounded memory, you have left regular territory and need the tape of a [[Turing Machine]].
{{< /note >}}

That single limitation — no unbounded counting — is exactly where the finite automaton ends and the universal computer begins. Everything else in this section is a study of what that extra memory buys you.

## See also

- [[Turing Machine]]
- [[Church–Turing Thesis]]
- [[Complexity Class]]
