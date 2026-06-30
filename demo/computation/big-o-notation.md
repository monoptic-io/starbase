---
title: Big-O Notation
aliases: [big o, asymptotic complexity, time complexity, big-o notation]
tags: [computation]
summary: A way to describe how an algorithm's cost grows with input size, keeping only the dominant term and ignoring constants.
weight: 60
---

# Big-O Notation

**Big-O notation** answers the question that matters when an algorithm meets real data: *how does its cost grow as the input gets bigger?* It deliberately throws away everything inessential — constant factors, low-order terms, the speed of your particular machine — and keeps only the **shape** of the growth. We say a routine runs in $O(n^2)$ time if, for large enough $n$, its step count is bounded by some constant times $n^2$.

Why ignore constants? Because they are a property of the *hardware and implementation*, not the *algorithm*. A faster CPU shifts a constant; it cannot rescue an algorithm whose work explodes as $2^n$. For large inputs, growth rate dominates everything else.

{{< note kind="key" title="Reading O(...)" >}}
$O(f(n))$ is an **upper bound on growth**: cost grows *no faster than* $f(n)$, up to a constant, for large $n$. Drop constants and lower terms — $3n^2 + 50n + 99$ is simply $O(n^2)$, because once $n$ is large the $n^2$ term buries the rest.
{{< /note >}}

## The growth zoo

A handful of growth rates cover almost everything you will meet. The gulf between them is enormous — and it only widens as $n$ grows.

{{< plot fn="Math.log2(x) ;; x ;; x*Math.log2(x) ;; x*x ;; Math.pow(2,x)" xmin="1" xmax="10" ymin="0" ymax="40" samples="300" height="340" title="Five growth rates on the same axes" caption="Bottom to top: log n (binary search), n (a single scan), n log n (good sorting), n² (nested loops), 2ⁿ (brute force). At x=10 they span from ~3 to over 1000 — and 2ⁿ has only just begun to erupt." >}}

The same picture in words, from gentlest to most ruinous:

- $O(1)$ — **constant.** Array lookup, hash insert. Input size is irrelevant.
- $O(\log n)$ — **logarithmic.** Halve the search space each step: binary search, balanced trees.
- $O(n)$ — **linear.** Look at each item once.
- $O(n \log n)$ — **linearithmic.** The best comparison sorts (merge sort, heap sort).
- $O(n^2)$ — **quadratic.** Every pair: nested loops, naive sorts.
- $O(2^n)$, $O(n!)$ — **exponential / factorial.** Try every subset or ordering. Hopeless past small $n$.

{{< note kind="warning" title="Why exponential is a wall, not a hill" >}}
At a billion steps per second, an $O(2^n)$ algorithm handles $n = 30$ in a second, $n = 50$ in two weeks, and $n = 70$ in **37,000 years**. Adding *one* element doubles the work. This cliff is exactly why the [[NP-Completeness|NP-complete]] problems — for which only exponential algorithms are known — are considered intractable, and why [[P versus NP]] is the question it is.
{{< /note >}}

## Best, worst, and average

Big-O usually describes the **worst case** — the guarantee that cost never exceeds this bound. Sometimes the typical case differs sharply: quicksort is $O(n^2)$ in the worst case but $O(n \log n)$ on average, which is why it is fast in practice. Sibling notations sharpen the picture: $\Omega$ is a lower bound (no faster than), and $\Theta$ pins growth from both sides (exactly this rate).

This vocabulary is what lets us *rank* algorithms independently of hardware — the foundation for classifying problems by difficulty in [[Complexity Class]]es, and for appreciating why an $O((V+E)\log V)$ method like [[Dijkstra's Algorithm]] scales gracefully to enormous graphs.

{{< quiz question="An algorithm runs in 5n² + 1000n + 200 steps. Its Big-O time complexity is..." options="O(5n²)|O(n² + n)|O(n²)|O(1000n)" answer="3" explain="Big-O keeps only the fastest-growing term and drops constant factors. For large n the n² term dominates the linear and constant terms, and the coefficient 5 is discarded — leaving O(n²)." >}}

## See also

- [[Complexity Class]]
- [[Dijkstra's Algorithm]]
- [[P versus NP]]
