---
title: Number Theory
aliases: [number theory]
tags: [number-theory]
summary: The study of the integers — primes, divisibility, and congruence — whose oldest puzzles became the foundation of modern cryptography.
weight: 75
---

# Number Theory

**Number theory** is the study of the whole numbers $1, 2, 3, \dots$ and the surprisingly deep structure hiding inside them. It asks questions a child can pose — *which numbers are prime? when does one number divide another? what are the integer solutions of this equation?* — and answers them with some of the most beautiful and difficult mathematics ever written. For most of its history it was prized precisely because it was *useless*: pure thought, untainted by application. Then, in the 1970s, the difficulty of factoring large integers became the lock on every secure connection on Earth.

{{< note kind="key" title="The integers are not as simple as they look" >}}
Addition and multiplication of whole numbers are taught in primary school, yet they generate questions that remain unsolved after millennia. Are there infinitely many twin primes? Is every even number a sum of two primes? Number theory is the art of finding hard problems in the simplest possible setting — and occasionally discovering they secure the modern world. The thread runs straight from [[Prime Number|primes]] to [[Modular Arithmetic]] to [[Cryptography]].
{{< /note >}}

## Two pillars: primes and congruence

Everything in this section grows from two ideas. The first is the **[[Prime Number|prime]]** — a number divisible only by $1$ and itself. Primes are the *atoms of multiplication*: the [[Fundamental Theorem of Arithmetic]] says every integer is built from them in exactly one way. The [[Sieve of Eratosthenes]] hunts them down; the [[Ulam Spiral]] reveals their eerie, half-ordered scatter.

The second pillar is **[[Modular Arithmetic|congruence]]** — arithmetic that wraps around a fixed modulus, like a clock. It turns the infinite integers into a finite ring where powers cycle and inverses sometimes exist. [[Euler's Totient Function]] counts the units of that ring, [[Fermat's Little Theorem]] governs its exponents, and together they make [[RSA]] possible.

## The road to cryptography

The deepest modern application of number theory is **[[Cryptography]]**. The asymmetry at its heart is purely number-theoretic: multiplying two primes is instant, but factoring their product is believed to be intractable. [[Euler's Totient Function]] sizes the key space, [[Fermat's Little Theorem]] guarantees that decryption undoes encryption, and [[Modular Arithmetic]] is the arena where it all plays out.

{{< chart type="bar" data="primes:25, twin pairs:8, perfect:1" title="Curiosities among the first 100 integers" ylabel="count" caption="Of the numbers 1–100: 25 are prime, 8 are members of twin-prime pairs below 100, and exactly 1 (namely 28) is a perfect number after 6. Simple counts, deep questions." >}}

## The pages in this section

- [[Prime Number]] — the indivisible atoms of multiplication, and a proof there are infinitely many.
- [[Sieve of Eratosthenes]] — an ancient elimination algorithm that finds every prime up to $N$.
- [[Fundamental Theorem of Arithmetic]] — every integer factors into primes in exactly one way.
- [[Greatest Common Divisor]] — the largest shared divisor, found by Euclid's algorithm.
- [[Euler's Totient Function]] — how many integers below $n$ are coprime to it.
- [[Fermat's Little Theorem]] — $a^{p}\equiv a\pmod p$, the seed of primality testing.
- [[Ulam Spiral]] — coil the integers, mark the primes, and watch diagonals appear.
- [[Continued Fraction]] — nested fractions that give the best rational approximations.
- [[Diophantine Equation]] — equations demanding whole-number answers.

## See also

- [[Cryptography]]
- [[Modular Arithmetic]]
- [[RSA]]
