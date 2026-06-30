---
title: Probability Distribution
aliases: [distribution, probability density]
tags: [probability]
summary: A complete description of how likelihood is spread across the possible outcomes of a random quantity.
weight: 20
---

# Probability Distribution

A **probability distribution** is the full ledger of a random quantity: for every possible outcome, how likely it is. It answers not just "what can happen?" but "*how often*?" — and that complete accounting is what lets us reason precisely about something whose individual results are uncertain. Roll a die and you cannot predict the face; name its distribution — each of $1$ through $6$ with probability $\tfrac{1}{6}$ — and you have said everything there is to say.

The total likelihood always sums (or integrates) to $1$: *something* must happen. Beyond that single constraint, distributions come in endlessly many shapes, each encoding a different kind of randomness.

## Discrete and continuous

When outcomes are separate and countable — coin flips, dice, the number of emails in an hour — the distribution is **discrete**, and we list a probability for each value. When outcomes form a continuum — a height, a waiting time, a measurement error — we instead use a **probability density**: probability per unit length, where the *area* under the curve over an interval gives the chance of landing in it.

{{< note kind="note" title="Density is not probability" >}}
For a continuous quantity the chance of any *exact* value is zero — there are infinitely many. Only ranges carry probability, given by the area under the density. A tall density just means outcomes cluster there, not that any single point is "likely."
{{< /note >}}

## Three distributions worth knowing

**Uniform** — every outcome equally likely; the flat distribution of a fair die or an honest spinner. Maximum ignorance, maximum [[Entropy]].

**Binomial** — count the heads in $n$ independent flips. Its bars rise to a peak at the expected number and fall away symmetrically (for a fair coin).

**Normal** — the bell curve, set by a mean (where it centers) and a standard deviation (how wide). It is the universal limit shape that the [[Central Limit Theorem]] explains, and it appears whenever many small independent effects add up.

{{< chart type="bar" data="0.5 2 4.4 9.2 15.3 20.1 21.0 15.3 9.2 4.4 1.0 0.5" title="Binomial(10, ½): the chance of getting k heads in 10 fair flips" xlabel="number of heads k" ylabel="probability (%)" >}}

The binomial above already looks bell-shaped — and it should. A binomial is a sum of $n$ independent coin flips, so as $n$ grows it converges to a normal curve, a first glimpse of the central limit theorem at work.

## Summarizing a distribution

We rarely carry the whole distribution around. Two numbers usually suffice: the **mean** (the balance point, the expected value) and the **variance** (how far outcomes typically stray from it). The [[Law of Large Numbers]] says a sample average homes in on the mean; the [[Central Limit Theorem]] says the *errors* in that average are themselves normally distributed, with a width set by the variance. The distribution is the object; mean and variance are its shadow.

## See also

- [[Central Limit Theorem]]
- [[Entropy]]
- [[Random Walk]]
