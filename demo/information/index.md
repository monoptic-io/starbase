---
title: Information Theory
aliases: [information theory]
tags: [information]
summary: The mathematics of measuring information — quantifying surprise, and the hard limits it sets on compression and communication.
weight: 60
---

# Information Theory

**Information theory** is the science of *how much* — how much a message tells you, how far it can be squeezed, and how fast it can be pushed through a noisy wire without error. Born in a single 1948 paper by Claude Shannon, it replaced the vague notion of "information" with a number measured in **bits**, and then proved exactly what those bits can and cannot do.

The whole edifice rests on one idea: **information is surprise**. A message that you could have predicted tells you nothing; a message you did not expect tells you a lot. Make that intuition quantitative and everything else follows — the [[Entropy]] of a source, the floor on [[Data Compression]], the capacity of a [[Channel Capacity|noisy channel]], the redundancy needed for [[Error-Correcting Code|error correction]].

{{< note kind="key" title="One idea, measured in bits" >}}
The **bit** is the unit of information: the answer to one yes/no question whose two outcomes are equally likely. Everything in this section is an accounting of bits — how few you need to name an outcome ([[Shannon Information]]), how few on average to name a stream of them ([[Entropy]]), and how many extra you must spend to survive noise ([[Error-Correcting Code]]).
{{< /note >}}

## The two limits Shannon proved

Information theory is anchored by two theorems that bracket every communication system ever built.

- **Source coding** (compression). No lossless code can represent a source in fewer bits per symbol than its entropy. [[Entropy]] is the floor; [[Huffman Coding]] and its relatives chase it. This ties the abstract [[Entropy]] of a [[Probability Distribution]] directly to file sizes on a disk.
- **Channel coding** (communication). Every noisy channel has a [[Channel Capacity|capacity]] $C$. Below $C$ you can communicate with *arbitrarily small* error using enough redundancy; above $C$ you cannot. The bridge between input and output is [[Mutual Information]].

## Information meets randomness and chaos

Information theory does not live alone. Its measures are expectations over a [[Probability Distribution]], so it shares a border with [[Probability & Random Processes]]. And its deepest question — what is the shortest possible description of an object? — is [[Kolmogorov Complexity]], which says a string is *random* exactly when it is incompressible. That same idea reaches into [[Chaos]]: a chaotic trajectory generates fresh, incompressible information at a rate set by its [[Lyapunov Exponent]], so determinism and unpredictability shake hands here.

## The pages in this section

- [[Entropy]] — average surprise, the master quantity.
- [[Shannon Information]] — the surprise of a single outcome, $-\log_2 p$ bits.
- [[Data Compression]] — squeezing out redundancy down to the entropy floor.
- [[Huffman Coding]] — provably optimal symbol-by-symbol codes.
- [[Mutual Information]] — how much one variable reveals about another.
- [[Channel Capacity]] — the top speed of a noisy channel.
- [[Error-Correcting Code]] — redundancy that detects and repairs flipped bits.
- [[Hamming Distance]] — counting the bit-flips between two strings.
- [[Kolmogorov Complexity]] — the shortest program that prints a string.

## See also

- [[Probability Distribution]]
- [[Kolmogorov Complexity]]
- [[Chaos]]
