---
title: Fermat's Little Theorem
aliases: [fermat little theorem]
tags: [number-theory]
summary: For a prime p, a^p ≡ a (mod p) for every integer a — a congruence that underlies fast primality testing.
weight: 60
---

# Fermat's Little Theorem

**Fermat's little theorem** states that if $p$ is a [[Prime Number|prime]], then for *every* integer $a$,

$$a^{p} \equiv a \pmod{p},$$

and, when $a$ is not a multiple of $p$, the sharper form

$$a^{p-1} \equiv 1 \pmod{p}.$$

It is "little" only by contrast with Fermat's Last Theorem. In practice it is enormous: this single congruence is the foundation of nearly every fast primality test and a load-bearing beam under [[RSA]].

## What it says, concretely

Take $p = 7$ and $a = 3$. Then $3^{6} = 729 = 104\cdot 7 + 1$, so $3^{6}\equiv 1\pmod 7$. Try any base coprime to $7$ and the same thing happens — raising it to the $6$th power lands back on $1$. The theorem promises this is no coincidence: in [[Modular Arithmetic|arithmetic mod a prime]], the powers of any nonzero base cycle through a length that always divides $p - 1$.

{{< plot fn="((Math.round(Math.pow(3,Math.round(x))) % 7) + 7) % 7" xmin="1" xmax="6" ymin="0" ymax="7" samples="6" height="300" title="3^k mod 7 for k = 1…6" caption="The powers 3, 2, 6, 4, 5, 1 cycle through every nonzero residue mod 7 and return to 1 at k = p − 1 = 6 — Fermat's little theorem in a single orbit." >}}

Here $3$ is a **primitive root** mod $7$: its powers visit all six nonzero residues before closing the loop. Not every base does — but every base's cycle length divides $p-1$, which is exactly why $a^{p-1}\equiv 1$.

## Why it is true

There is a one-line combinatorial proof. Consider the nonzero residues $\{1, 2, \dots, p-1\}$ and multiply each by $a$ (with $\gcd(a,p)=1$). Because multiplication by $a$ is invertible mod $p$, this just **permutes** the same set. So the two products are equal:

$$\prod_{k=1}^{p-1}(a k) \equiv \prod_{k=1}^{p-1} k \pmod p.$$

The left side is $a^{p-1}\,(p-1)!$ and the right is $(p-1)!$. Since $(p-1)!$ is coprime to $p$, cancel it to get $a^{p-1}\equiv 1$. It is a special case of **Euler's theorem** ([[Euler's Totient Function]]), which replaces $p-1$ by $\varphi(n)$ for any modulus $n$.

{{< note kind="key" title="Fermat as a primality test" >}}
Flip the theorem into a test. To probe whether $n$ is prime, pick a base $a$ and compute $a^{n-1}\bmod n$. If the result is **not** $1$, then $n$ is *definitely composite* — no further work needed. If it *is* $1$, $n$ is "probably prime." Because $a^{n-1}\bmod n$ takes only a few dozen modular squarings even for thousand-digit $n$, this is blisteringly fast — the basis of the **Fermat** and **Miller–Rabin** tests.
{{< /note >}}

## The catch: Carmichael numbers

The Fermat test has a flaw. A few rare composites masquerade as primes by passing it for *every* base coprime to them. The smallest is $561 = 3\cdot 11\cdot 17$: it is composite, yet $a^{560}\equiv 1 \pmod{561}$ for all $a$ coprime to $561$. These are the **Carmichael numbers**, and they are why robust testing uses the stronger Miller–Rabin refinement rather than Fermat's theorem alone.

{{< chart type="bar" data="561:1, 1105:1, 1729:1, 2465:1, 2821:1" title="The first Carmichael numbers" ylabel="(each is a Fermat liar for all bases)" caption="Rare composites — 561, 1105, 1729, 2465, 2821 — that satisfy aⁿ⁻¹ ≡ 1 for every coprime base, fooling the naive Fermat test. They thin out but never stop." >}}

{{< quiz question="A Fermat primality test computes a^(n−1) mod n and gets a result of 5, not 1. What can you conclude?" options="n is definitely prime|n is definitely composite|n is probably prime|Nothing at all" answer="2" explain="Fermat's theorem guarantees a^(n−1) ≡ 1 for prime n (and coprime a). Any other value is a certificate that n is composite — the test can prove compositeness outright, but never proves primality." >}}

## See also

- [[Modular Arithmetic]]
- [[Euler's Totient Function]]
- [[RSA]]
