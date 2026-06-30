---
title: Halting Problem
aliases: [halting, halting problem]
tags: [computation]
summary: The proof that no program can decide, for every program and input, whether it eventually halts or runs forever.
weight: 40
---

# Halting Problem

The **halting problem** asks for something that sounds eminently reasonable: a single program `halts(P, x)` that, given any program `P` and input `x`, returns `true` if `P` eventually stops on `x` and `false` if it loops forever. Such an oracle would be a debugger's dream — it would catch every infinite loop before you ran it.

In 1936 Turing proved that **no such program can exist**. The halting problem is *undecidable*: not merely hard, but provably beyond the reach of any algorithm, on any computer, ever. It was the first problem shown to be unsolvable, and it remains the template for nearly every impossibility result in computing.

## The diagonal contradiction

The proof is a jewel — a few lines of self-reference that close like a trap. Suppose, for contradiction, that the decider `halts(P, x)` exists. Then we can build this mischievous program:

{{< note kind="note" title="The program that breaks the oracle" >}}
```
function paradox(P):
    if halts(P, P):     # does P halt when fed its own code?
        loop forever
    else:
        stop
```
`paradox` does the **opposite** of whatever its input would do on itself: if `P` halts on `P`, then `paradox` loops; if `P` loops on `P`, then `paradox` stops.
{{< /note >}}

Now ask the fatal question: **does `paradox(paradox)` halt?**

- If it **halts**, then by its own code `halts(paradox, paradox)` returned `true`, which sends it into `loop forever` — so it does *not* halt. Contradiction.
- If it **runs forever**, then `halts(paradox, paradox)` returned `false`, which sends it to `stop` — so it *does* halt. Contradiction.

Every branch contradicts itself. The only false assumption was that `halts` existed in the first place. Therefore it cannot. This is the same **diagonalization** Cantor used to prove the reals are uncountable, turned against computation.

{{< note kind="key" title="The core trick: self-reference" >}}
The contradiction springs entirely from feeding a program *its own description* as input, then negating the answer. Any system powerful enough to (a) simulate programs and (b) negate a yes/no result can be made to contradict itself this way. The price of universal computation is that some questions about it become unanswerable from inside.
{{< /note >}}

## Why it matters everywhere

The halting problem is not an exotic edge case — it is a barrier you hit constantly, because countless practical questions secretly *contain* it.

- **"Does this program ever crash / leak memory / reach this line?"** — undecidable in general, by reduction from halting. This is why no compiler can flag every bug.
- **"Are these two programs equivalent?"** — undecidable.
- **"Will this loop terminate?"** — the halting problem itself.

Each is proven impossible by showing that a solver for it would also solve halting. That ripple outward is the subject of [[Decidability]], where the halting problem is the seed from which a whole hierarchy of unsolvable problems grows.

{{< note kind="warning" title="Undecidable ≠ useless" >}}
"No algorithm works for *every* case" does not mean "no algorithm helps in *any* case". Real tools — type checkers, static analyzers, termination provers — soundly decide many specific programs and answer "don't know" on the rest. Undecidability forbids a *perfect, total* decider, not a useful partial one.
{{< /note >}}

{{< quiz question="What logical technique drives the proof that the halting problem is undecidable?" options="Proof by induction|Diagonalization via self-reference and contradiction|Exhaustive search|The pigeonhole principle" answer="2" explain="You assume a halting decider exists, then build a program that asks the decider about itself and does the opposite. Both answers lead to contradiction — the hallmark of diagonalization, the same method behind Cantor's uncountability proof and Gödel's incompleteness theorems." >}}

## See also

- [[Decidability]]
- [[Turing Machine]]
- [[Church–Turing Thesis]]
