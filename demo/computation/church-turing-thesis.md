---
title: Church–Turing Thesis
aliases: [church turing thesis, church-turing thesis]
tags: [computation]
summary: The claim that every reasonable model of computation can compute exactly the same functions as a Turing machine.
weight: 30
---

# Church–Turing Thesis

The **Church–Turing thesis** is the bedrock assumption of computer science: *any function that can be computed by a mechanical procedure at all can be computed by a [[Turing Machine|Turing machine]].* Whatever "effectively calculable" means — pencil and paper, an abacus, a laptop, a quantum chip, a galaxy-sized supercomputer — it computes **no more** than that simple tape-and-head device.

It is called a *thesis*, not a theorem, for a deliberate reason. "Computable by a mechanical procedure" is an informal, intuitive notion; "computable by a Turing machine" is a precise mathematical one. The thesis asserts these two coincide — and because one side is informal, it can never be *proved*, only supported by overwhelming evidence.

## The convergence that convinced everyone

In the 1930s, several people chased the idea of *effective computation* from completely different directions — and every road led to the same destination.

{{< columns count="2" >}}
**Alonzo Church** built the **λ-calculus**, a system of pure functions and substitution, with no machine in sight.

**Alan Turing** imagined his **machine** — tape, head, states — a mechanical metaphor.

**Kurt Gödel & Herbrand** defined the **general recursive functions** through equations.

**Emil Post** described **rewriting systems** of string-manipulation rules.
{{< /columns >}}

These models look nothing alike. Yet each was proved to compute **exactly the same** set of functions. When utterly independent attempts to capture "computable" all collapse into one class, that class is no accident — it is the real thing.

{{< note kind="key" title="Why it matters" >}}
The thesis is what lets us say "no algorithm exists" and mean it. To prove a problem like the [[Halting Problem]] is unsolvable, it is enough to show *no Turing machine* solves it — because by the thesis, that rules out every possible algorithm in any language on any hardware, present or future.
{{< /note >}}

## What it does and does not say

The thesis is about *what* can be computed, never *how fast*. A modern CPU and a Turing machine solve the same problems, but the CPU may be astronomically quicker — that gap is the subject of [[Complexity Class]]es and [[Big-O Notation]], not computability.

It also survives every "more powerful" machine ever proposed. Add randomness, parallelism, multiple tapes, or quantum superposition, and the set of computable functions does **not** grow — only the *efficiency* changes. Each such model can be simulated by a plain Turing machine, so each computes the same things.

{{< note kind="note" title="Hypothetical exceptions" >}}
The thesis could in principle be broken by a *hypercomputer* — a device exploiting infinite precision, true continuous physics, or a trip past a black hole's horizon to do infinitely many steps in finite time. No such machine is known to be physically realizable, and the **physical** Church–Turing thesis conjectures none ever will be.
{{< /note >}}

{{< quiz question="Why is the Church–Turing thesis a thesis rather than a proven theorem?" options="No one has tried to prove it|It relates a precise formal model to an informal, intuitive notion of computation|It has known counterexamples|It is only true for small inputs" answer="2" explain="One side — 'effectively calculable by any mechanical procedure' — is informal. You cannot formally prove an equivalence with an informal notion, only accumulate evidence. And the evidence is overwhelming: every formal model devised computes exactly the Turing-computable functions." >}}

## See also

- [[Turing Machine]]
- [[Halting Problem]]
- [[Complexity Class]]
