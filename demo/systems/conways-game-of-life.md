---
title: Conway's Game of Life
aliases: [game of life, life, conway's game of life, conway game of life]
tags: [systems, emergence, discrete]
summary: Two rules on a grid of cells produce gliders, oscillators, and self-replicating structure — the most famous cellular automaton.
weight: 50
---

# Conway's Game of Life

**Conway's Game of Life** is a two-dimensional [[Cellular Automaton]] invented by John Conway in 1970, and it is the canonical demonstration that life-like complexity can spring from almost nothing. The board is an infinite grid of cells, each **alive** or **dead**. There are no players and no moves after the first: you set an initial pattern, and from then on the whole grid updates in lockstep according to two short rules. What follows — blinking oscillators, gliding spaceships, glider guns, even working logic gates — is pure emergence.

## The rules: B3/S23

Each cell looks at its eight surrounding neighbors and obeys:

{{< note kind="key" title="Birth on 3, Survival on 2 or 3" >}}
- A **dead** cell with exactly **3** live neighbors becomes alive (**birth**).
- A **live** cell with **2 or 3** live neighbors stays alive (**survival**).
- Every other cell dies or stays dead — from loneliness (fewer than 2) or overcrowding (more than 3).

In shorthand: **B3/S23**. Two numbers, and the rest is consequence.
{{< /note >}}

The balance is delicate. Too generous a rule and the grid floods to a uniform on-state; too strict and everything dies. B3/S23 sits right at the knife's edge where structures can persist *and* move *and* interact — the same "edge of chaos" that makes a [[Cellular Automaton]] interesting.

## A zoo of emergent objects

Out of two rules comes a whole taxonomy:

- **Still lifes** — patterns that never change (the block, the beehive). These are [[Fixed Point|fixed points]] of the grid map.
- **Oscillators** — patterns that cycle through a repeating sequence (the blinker, the pulsar). Each is a discrete [[Limit Cycle]], a closed loop the dynamics settle onto — an [[Attractor]] in the space of grid states.
- **Spaceships** — patterns that translate across the grid as they oscillate. The humble **glider** is the smallest, ferrying a packet of "aliveness" diagonally forever.

Because gliders carry information and can be made to collide in controlled ways, the Game of Life is **Turing-complete**: in principle you can build any computer inside it. Universal computation, from B3/S23.

{{< sim name="life" cell="10" caption="Conway's Game of Life. Look for blinkers (period-2 oscillators) and gliders drifting diagonally across the board." >}}

## Why it matters

The Game of Life reframed a deep question: complexity does not require a complex cause. A simulated universe with two rules generates open-ended structure, self-organization, and computation — a vivid, visual argument that the elaborate order we see in nature can emerge from simple, local laws iterated many times over. It is the beating heart of the study of [[Complex Systems]].

## See also

- [[Cellular Automaton]]
- [[Attractor]]
- [[Complex Systems]]
