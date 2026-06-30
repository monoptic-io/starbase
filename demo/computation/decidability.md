---
title: Decidability
aliases: [decidable, undecidable, decidability]
tags: [computation]
summary: The line between problems an algorithm can always settle and those no algorithm can settle at all.
weight: 50
---

# Decidability

A yes/no problem is **decidable** if some algorithm answers it correctly for *every* input and always halts. It is **undecidable** if no such algorithm exists. This is the sharpest line in all of computer science: not between easy and hard, but between *possible* and *impossible*. An undecidable problem is not waiting for a cleverer programmer or a faster machine — by the [[Church–Turing Thesis]], it is beyond computation itself.

The flagship undecidable problem is the [[Halting Problem]]: no program can decide whether an arbitrary program halts. From that single seed, undecidability spreads across mathematics and computing.

## How undecidability spreads: reduction

You rarely prove a new problem undecidable from scratch. Instead you **reduce** the halting problem to it: show that *if* you could solve your problem `B`, you could solve halting too. Since halting is unsolvable, so is `B`. One impossibility, leveraged endlessly.

{{< note kind="key" title="Reduction in one line" >}}
To prove **B is undecidable**: assume a decider for `B`, then use it as a subroutine to decide the [[Halting Problem]]. The assumed decider must not exist. Almost every undecidability proof has this shape — the halting problem is the universal donor.
{{< /note >}}

Through reductions, a startling range of natural questions turn out undecidable:

- **Program equivalence** — do two programs compute the same function?
- **Rice's theorem** — *every* non-trivial question about what a program *computes* (as opposed to how it is written) is undecidable. This is sweeping: "does it ever output 7?", "is it virus-free?", "does it compute the identity?" — all undecidable.
- **Hilbert's tenth problem** — does a polynomial equation have integer solutions? Proven undecidable in 1970.
- **The Post correspondence problem** and many tiling and grammar questions.

## A finer map than "solvable / not"

Undecidable problems are not all equally unsolvable. Many sit just out of reach in a structured way.

{{< note kind="note" title="Semi-decidable: half an answer" >}}
The halting problem is **semi-decidable** (recognizable): just *run* the program; if it halts, you will eventually see it and answer "yes". The catch is the "no" case — you can never be sure it won't halt *one step later*, so you may wait forever. A problem is fully **decidable** exactly when both it and its complement are semi-decidable.
{{< /note >}}

This stratifies undecidable problems into an infinite tower (the *arithmetical hierarchy*), each level harder than the last — a landscape of impossibility with its own fine structure.

## The deepest link: information

Decidability reaches into [[Information Theory]] through [[Kolmogorov Complexity]] — the length of the *shortest program* that outputs a given string. It is the ultimate measure of a string's compressibility, and it is **undecidable**: no algorithm can compute it. The proof is a cousin of the [[Halting Problem|halting]] argument (formalized as Berry's paradox — "the smallest number not describable in fewer than twelve words" describes it in eleven). A genuinely random string is simply one no program shorter than itself can produce — and we can never *prove* a specific string is random.

{{< quiz question="A problem is decidable exactly when..." options="It has a 'yes' answer for every input|Both it and its complement are semi-decidable|Some program prints 'yes' when the answer is yes|It can be solved quickly" answer="2" explain="Decidable means an algorithm always halts with the correct yes-or-no answer. Equivalently, both the problem and its complement are semi-decidable: you can confirm 'yes' answers and confirm 'no' answers, so running both searches in parallel always terminates with a verdict." >}}

## See also

- [[Halting Problem]]
- [[Kolmogorov Complexity]]
- [[Church–Turing Thesis]]
