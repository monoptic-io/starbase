---
title: P versus NP
aliases: [p vs np, p np, p versus np]
tags: [computation]
summary: The open question of whether every problem whose solution is easy to check is also easy to solve.
weight: 80
---

# P versus NP

**P versus NP** is the deepest open problem in computer science and one of the seven Clay Millennium Prize Problems — solve it and collect a million dollars. Stripped to a sentence, it asks:

> **If a solution can be *checked* quickly, can it always be *found* quickly?**

The two letters name two [[Complexity Class]]es. **P** is the problems we can *solve* in polynomial time. **NP** is the problems whose answers we can *verify* in polynomial time. Every solvable-fast problem is obviously verifiable-fast, so $\mathrm{P} \subseteq \mathrm{NP}$. The trillion-dollar question is whether they are actually **equal**.

## Finding versus checking

The gap between the two feels enormous in everyday life, and that intuition is the whole problem.

{{< columns count="2" >}}
**Easy to check.** Handed a completed jigsaw, a solved Sudoku, or a valid travel route under budget, you can confirm it is correct almost instantly — a quick scan. These are NP.

**Hard to find?** Producing that solution from scratch seems to demand searching through an astronomical number of possibilities. *Seems* — but no one has proved it must.
{{< /columns >}}

If **P = NP**, that apparent difficulty is an illusion: every problem whose solution we can recognize, we could also *construct* efficiently. If **P ≠ NP** (which nearly all researchers believe), then some problems are intrinsically harder to solve than to verify — searching really is fundamentally costlier than checking.

{{< note kind="key" title="Why almost everyone bets P ≠ NP" >}}
Thousands of [[NP-Completeness|NP-complete]] problems have been attacked for over fifty years by brilliant people, and **none** has yielded a polynomial algorithm. If P = NP, a single such breakthrough would instantly collapse *all* of them. The continued failure is strong circumstantial evidence — but evidence is not proof, and the question remains formally open.
{{< /note >}}

## What hangs on the answer

This is not an abstract curiosity. A constructive proof that **P = NP** (with a practical algorithm) would reshape the world — for better and worse:

- **Cryptography would collapse.** [[RSA]], [[Diffie–Hellman]], and essentially all of [[Public-Key Cryptography]] rest on problems being hard to solve but easy to check. P = NP would, in principle, unlock them all.
- **Optimization would become trivial.** Protein folding, chip layout, logistics, theorem-proving, perfect machine-learning fits — vast searches would turn easy.
- **Mathematics itself** would shift: finding proofs would become as easy as checking them.

Conversely, a proof that **P ≠ NP** would *confirm* that the hardness our security relies on is real and permanent — and would require entirely new mathematical tools we do not yet possess. Either resolution is a landmark.

{{< note kind="warning" title="A subtle trap" >}}
P = NP would not mean *instant* solutions to everything. An $O(n^{50})$ algorithm is polynomial yet hopelessly slow. And the proof could be **non-constructive** — establishing that a fast algorithm exists without revealing it. The headline consequences assume not just P = NP but a *practical* algorithm, which is a stronger thing.
{{< /note >}}

## How the question is anchored

P versus NP is not chased across infinitely many problems at once. The [[NP-Completeness|NP-complete]] problems are the load-bearing pillars: a single one of them in P would drag *all* of NP into P, proving P = NP in one stroke. That is why effort concentrates on a handful of canonical hard problems — satisfiability, [[Graph]] coloring, the clique problem — knowing they all rise or fall together.

{{< quiz question="What would follow immediately if a polynomial-time algorithm were found for one NP-complete problem?" options="Only that one problem becomes easy|P would equal NP, and every NP problem would be solvable in polynomial time|The halting problem would become decidable|Nothing changes" answer="2" explain="Every problem in NP reduces in polynomial time to any NP-complete problem. A fast algorithm for one therefore gives a fast algorithm for all of NP, proving P = NP — which is exactly why NP-complete problems are the focus of the question." >}}

## See also

- [[NP-Completeness]]
- [[Complexity Class]]
- [[Big-O Notation]]
