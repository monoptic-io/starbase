---
title: Bayes' Theorem
aliases: [bayes, bayesian, conditional probability]
tags: [probability]
summary: A rule for updating the probability of a hypothesis as new evidence arrives.
weight: 80
---

# Bayes' Theorem

**Bayes' theorem** is the mathematics of learning from evidence. It takes a belief you held *before* seeing some data — the **prior** — and tells you exactly how to revise it *after* — the **posterior**. It is the bridge between "what I thought was likely" and "what I should think now," and it underlies spam filters, medical diagnosis, machine learning, and the scientific method's habit of updating theories against observation.

The whole theorem follows from a simple fact about **conditional probability**: the chance of $A$ *and* $B$ can be read two ways, $P(A)P(B\mid A) = P(B)P(A\mid B)$. Rearranged, that is Bayes' rule.

## The rule

For a hypothesis $H$ and evidence $E$:

{{< eq number="1" >}}
P(H \mid E) = \frac{P(E \mid H)\,P(H)}{P(E)}.
{{< /eq >}}

Read it piece by piece. $P(H)$ is the **prior** — how plausible the hypothesis was beforehand. $P(E\mid H)$ is the **likelihood** — how well the hypothesis predicts the evidence. $P(E)$ is the total chance of seeing that evidence at all, the normalizer that keeps things summing to $1$. Multiply prior by likelihood, divide by the evidence's overall rate, and out comes the updated belief $P(H\mid E)$.

{{< note kind="key" title="The prior never disappears" >}}
A common mistake is to read $P(E\mid H)$ — "the test fires when the disease is present" — as if it were $P(H\mid E)$ — "the disease is present when the test fires." They are wildly different when the hypothesis is *rare*. Bayes' theorem is precisely the correction: a strong likelihood multiplied by a tiny prior can still leave you fairly unconvinced.
{{< /note >}}

## The base-rate surprise

Nothing makes Bayes' theorem click like the classic medical-test puzzle. The numbers feel impossible until you trace where everyone goes.

{{< columns count="2" >}}
**The setup**

A disease affects **1 in 1,000** people. A test is **99% accurate**: it correctly flags 99% of sick people, and correctly clears 99% of healthy people. You test positive.

*How likely are you to actually have the disease?* Most people answer "about 99%." The truth is closer to **9%**.

---

**Why so low?**

Imagine **100,000 people**.

- **100** are sick. The test catches $99$ of them. → **99 true positives**
- **99,900** are healthy. The test wrongly flags $1\%$ of them: $999$. → **999 false positives**

Of the $99 + 999 = 1{,}098$ positive tests, only $99$ are truly sick:

$$ P(\text{sick}\mid +) = \frac{99}{1098} \approx 9\%. $$

The disease is so **rare** that the small flood of false positives from the huge healthy majority swamps the few true ones.
{{< /columns >}}

The lesson generalizes far past medicine: when a hypothesis starts out unlikely, even strong evidence may only nudge it. The base rate — the prior — is not a detail to wave away; it is half the calculation.

{{< quiz question="In the example, what makes the posterior so much lower than the test's 99% accuracy?" options="The test is poorly designed|The disease's low base rate makes false positives outnumber true positives|99% accuracy actually means 9% accuracy|The math only works for diseases" answer="2" explain="Because only 1 in 1,000 people are sick, the 1% false-positive rate applied to the huge healthy population produces far more false positives than there are true positives. The rare prior dominates." >}}

## Updating, again and again

Bayes' theorem is not a one-shot calculation — it is a *loop*. Today's posterior becomes tomorrow's prior, ready to be updated by the next piece of evidence. Run enough evidence through it and beliefs converge on the truth, the same way a [[Probability Distribution]] of guesses sharpens with data. This iterative tightening is the engine of Bayesian inference and a close cousin of how the [[Law of Large Numbers]] grinds noisy samples into certainty.

## See also

- [[Probability Distribution]]
- [[Law of Large Numbers]]
- [[Information Theory]]
