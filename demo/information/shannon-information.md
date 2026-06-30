---
title: Shannon Information
aliases: [self-information, surprise, bit, shannon]
tags: [information]
summary: The information content of a single outcome — its surprise, equal to minus the log of its probability.
weight: 20
---

# Shannon Information

The **Shannon information** (or *self-information*, or *surprise*) of an outcome is the amount you learn when it happens. Shannon set it equal to

$$I(x) = -\log_2 p(x) \quad \text{bits,}$$

where $p(x)$ is the probability of that outcome. The rule is simple and inevitable: **the rarer the event, the more information it carries.** Learning that a tossed fair coin came up heads ($p = \tfrac12$) gives you exactly 1 bit. Learning the outcome of one fair six-sided die ($p = \tfrac16$) gives $\log_2 6 \approx 2.585$ bits. Learning something certain ($p = 1$) gives 0 bits — you already knew it.

{{< note kind="note" title="Why a logarithm?" >}}
Three demands pin down the formula. Surprise must (1) **decrease** as probability rises, (2) be **zero** for a certain event, and (3) **add up** for independent events — finding out two unrelated facts should give you the sum of their informations. Only $-\log p$ satisfies all three, because $\log(p_1 p_2) = \log p_1 + \log p_2$ turns the *product* of independent probabilities into a *sum* of bits.
{{< /note >}}

## The shape of surprise

Plotted against probability, surprise sweeps from a towering spike near the impossible ($p \to 0$, information $\to \infty$) down to nothing at the certain ($p = 1$). Common events are cheap; rare ones are expensive. This is exactly why a good code spends *short* codewords on frequent symbols and saves the long ones for rarities.

{{< plot fn="-Math.log2(x)" xmin="0.01" xmax="1" ymin="0" ymax="6.7" title="Self-information I(p) = -log₂ p" caption="Surprise in bits versus the probability of the event. Rare events (left) carry many bits; certainties (right) carry none." >}}

## From one outcome to the average

[[Entropy]] is nothing more than the *average* self-information of a source — each outcome's surprise $-\log_2 p_i$ weighted by how often it occurs:

$$H = \sum_i p_i\,I(x_i) = -\sum_i p_i \log_2 p_i.$$

So self-information is the atom and entropy is the expectation built from it. A single **bit** is the natural unit throughout: the information in one equally likely yes/no answer. Using base-$e$ instead of base-2 logarithms measures the same thing in *nats*; the choice of base is just a choice of unit.

{{< quiz question="An event has probability 1/8. What is its self-information?" options="1 bit|3 bits|8 bits|0 bits" answer="2" explain="I = -log2(1/8) = log2(8) = 3 bits. Each halving of probability adds exactly one bit of surprise." >}}

## See also

- [[Entropy]]
- [[Data Compression]]
- [[Probability Distribution]]
