---
title: Probability & Random Processes
aliases: [probability, random processes]
tags: [probability]
summary: How randomness, far from being formless, organizes itself into sharp and predictable structure — bell curves, steady states, and laws of large numbers.
weight: 55
---

# Probability & Random Processes

Flip one coin and you cannot say what happens. Flip ten thousand and you can say, almost to the percent, how many land heads. This is the quiet miracle at the center of probability: **individual randomness, piled up, becomes collective certainty.** The single step is unpredictable; the long run is lawful. This section is a tour of that paradox — randomness with structure.

The structure shows up everywhere once you look. The aimless stagger of a [[Random Walk]] spreads at a precise rate, $\sqrt{n}$. Sums of wildly different random quantities all relax into the same bell-shaped curve — the [[Central Limit Theorem]]. A memoryless hop between states settles into a fixed long-run distribution — a [[Markov Chain]]. Even our *beliefs* update by a rule, [[Bayes' Theorem]]. Chance is not the absence of law; it is a different kind of law.

{{< note kind="key" title="Two threads run through this section" >}}
- **Randomness builds geometry.** A [[Markov Chain]]'s stationary distribution is an eigenvector — the very same machinery powers [[PageRank]], where a random web-surfer's wandering ranks the whole internet.
- **Discrete randomness has a continuous limit.** Shrink the steps of a [[Random Walk]] and it becomes [[Brownian Motion]], the jiggling thread that ties probability to diffusion, to [[Reaction–Diffusion]] patterns, and to the unpredictability of [[Chaos]].
{{< /note >}}

## The pages in this section

- [[Random Walk]] — a sum of random steps; the prototype of all diffusion.
- [[Probability Distribution]] — how likelihood is spread over the possible outcomes.
- [[Central Limit Theorem]] — why the bell curve is everywhere.
- [[Law of Large Numbers]] — why averages converge to the truth.
- [[Markov Chain]] — memoryless dynamics and their stationary states.
- [[Brownian Motion]] — the continuous limit of a random walk.
- [[Monte Carlo Method]] — computing hard answers by rolling dice.
- [[Bayes' Theorem]] — turning evidence into updated belief.
- [[Galton Board]] — the central limit theorem made of falling beads.
- [[Benford's Law]] — why real-world numbers begin with the digit 1 about 30% of the time.

## How it connects

Probability is the hinge between the deterministic systems of the earlier sections and the *information* of the next. A [[Markov Chain]] is a [[Dynamical System]] whose rule is a transition matrix; the entropy of a [[Probability Distribution]] is the bridge to [[Information Theory]] and [[Entropy]]. Randomness is not a detour from structure — it is one more way the same mathematics speaks.

## See also

- [[Markov Chain]]
- [[Brownian Motion]]
- [[Information Theory]]
