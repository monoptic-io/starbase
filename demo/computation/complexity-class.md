---
title: Complexity Class
aliases: [complexity classes, complexity class]
tags: [computation]
summary: A grouping of problems by the computational resources — time or memory — needed to solve them.
weight: 70
---

# Complexity Class

A **complexity class** is a bucket that holds all problems solvable within some budget of resources — a bound on **time** (number of steps) or **space** (amount of memory), measured as a function of input size with [[Big-O Notation]]. Sorting problems into these buckets is how theorists map the landscape of difficulty: not "can this be solved?" (that is [[Decidability]]) but "how *expensive* is it, fundamentally?"

The classes nest inside one another like Russian dolls, from the comfortably efficient out to the astronomically expensive.

## The headline classes

{{< note kind="key" title="P and NP, the two that matter most" >}}
**P** (polynomial time) — problems *solvable* in $O(n^k)$ time for some fixed $k$. This is the practical definition of **tractable**: sorting, shortest paths, primality testing. The runtime grows politely.

**NP** (nondeterministic polynomial time) — problems whose *proposed solutions can be checked* in polynomial time, even if finding one seems to need a search. Given a filled-in Sudoku you can verify it in seconds; finding it may be far harder.
{{< /note >}}

Every problem in **P** is also in **NP** (if you can *solve* it fast, you can *check* a solution fast — just re-solve). Whether the reverse holds — whether easy-to-check always means easy-to-solve — is the [[P versus NP]] question, the most famous open problem in the field.

Beyond NP the buckets keep growing:

- **PSPACE** — solvable using a *polynomial amount of memory*, with no limit on time. Roomy enough to contain all of NP. Many two-player games (generalized chess, Go) live here.
- **EXP** (EXPTIME) — solvable in exponential time $O(2^{n^k})$. Provably larger than P: some problems *require* exponential time, full stop.

What we know for certain is a chain of inclusions:

$$\mathrm{P} \subseteq \mathrm{NP} \subseteq \mathrm{PSPACE} \subseteq \mathrm{EXP}.$$

## The frustrating state of knowledge

Here is the scandal: we know $\mathrm{P} \subsetneq \mathrm{EXP}$ — the two *endpoints* of that chain are genuinely different. So at least one of the $\subseteq$ links must be a strict $\subsetneq$. Yet we cannot prove a single one of the intermediate inclusions is strict. It is entirely possible (though wildly unlikely) that $\mathrm{P} = \mathrm{NP} = \mathrm{PSPACE}$, with the gap hiding entirely between PSPACE and EXP.

{{< chart type="bar" data="2,8,20,40" labels="P,NP,PSPACE,EXP" title="Nested classes — schematic 'reach' of each (illustrative)" ylabel="problems reachable" caption="Each class contains all the ones before it. The boundaries drawn here look crisp, but proving any single inclusion is strict — most urgently P vs NP — is beyond current mathematics." >}}

{{< note kind="note" title="Why polynomial = 'tractable'?" >}}
It is a useful fiction. An $O(n^{100})$ algorithm is polynomial yet useless, and $O(2^{0.0001 n})$ is exponential yet fine for a while. But polynomials are closed under composition and addition — chain two efficient steps and you stay efficient — which makes P a robust, machine-independent notion of "feasible". In practice, problems in P almost always have *small*-exponent algorithms.
{{< /note >}}

## Why classes are robust

A crucial fact makes this whole taxonomy meaningful: the classes barely care about the machine. By the [[Church–Turing Thesis|Church–Turing principle]], reasonable models simulate each other with only polynomial overhead, so **P is the same class** whether you use a Turing machine, a laptop, or a [[Finite Automaton|register machine]]. That stability is what lets P and NP name real, hardware-independent facts about problems — and what makes [[NP-Completeness]] such a powerful tool.

{{< quiz question="What property defines the class NP?" options="Problems that cannot be solved|Problems solvable in polynomial time|Problems whose candidate solutions can be verified in polynomial time|Problems needing exponential memory" answer="3" explain="NP is the class of problems for which a proposed solution (a 'certificate') can be checked in polynomial time. Solving may be hard; verifying is easy. P ⊆ NP because anything you can solve quickly you can also verify quickly." >}}

## See also

- [[P versus NP]]
- [[Big-O Notation]]
- [[NP-Completeness]]
