---
title: Cryptography
aliases: [cryptography, crypto]
tags: [cryptography]
summary: The art and mathematics of keeping secrets and proving identity — from a Roman shift cipher to the public-key algebra that secures the modern web.
weight: 70
---

# Cryptography

**Cryptography** is the science of communicating in the presence of adversaries — keeping a message *secret* from everyone but its intended reader, and proving that it really came from who it claims to. For most of history it was a craft of clever tricks; in the twentieth century it became a branch of mathematics with **theorems** about what is possible and what is forever out of reach.

Two questions run through everything here:

- **Confidentiality** — can I scramble a message so only the right person can read it?
- **Authenticity** — can I prove *I* wrote it, and that nobody changed a single bit?

The journey runs from ciphers a child can break to algorithms whose security rests on problems no computer can solve in the lifetime of the universe.

{{< note kind="key" title="Kerckhoffs's principle" >}}
A cryptosystem should be secure **even if everything about it is public except the key**. Never trust a secret algorithm — trust a public one with a secret key. Every system on this page (except the broken historical ones) is designed to be analyzed in the open. Security lives in the key, not in the obscurity of the method.
{{< /note >}}

## Two eras: symmetric and public-key

For thousands of years all cryptography was **symmetric**: sender and receiver shared one secret key, used to both lock and unlock. The [[Caesar Cipher]], the [[Substitution Cipher]], and the unbreakable [[One-Time Pad]] all live here. The catch is brutal — before you can talk secretly, you must *already* have shared a secret, somehow.

In the 1970s came the revolution. **[[Public-Key Cryptography]]** split the key in two: a public lock anyone can snap shut, and a private key only you can open. Suddenly strangers could exchange secrets over a tapped wire. [[Diffie–Hellman]] lets two parties conjure a shared key in plain sight; [[RSA]] turns the difficulty of factoring large numbers into a padlock. All of it is built on the quiet algebra of [[Modular Arithmetic]].

## The threads back to information and probability

Cryptography is not an island — it borrows its deepest ideas from the rest of this field.

- **[[Entropy]] and perfect secrecy.** Shannon proved that a cipher is *perfectly secret* only when the key carries at least as much entropy as the message. That single inequality explains why the [[One-Time Pad]] is unbreakable and why every shorter key is, in principle, attackable.
- **[[Hamming Distance]] and the avalanche.** A good [[Hash Function]] must turn a one-bit change in its input into a change in roughly *half* the output bits — a large Hamming distance from a tiny nudge. Avalanche is what makes a fingerprint trustworthy.
- **[[Modular Arithmetic]] as the engine.** Wrap-around arithmetic on a finite ring of integers is the playground where one-way functions live: easy to compute forward, ruinously hard to reverse.

## The pages in this section

- [[Caesar Cipher]] — shift every letter by a fixed amount; 25 keys, broken in seconds.
- [[Substitution Cipher]] — any letter-to-letter map; an astronomical keyspace undone by frequency analysis.
- [[One-Time Pad]] — XOR with a truly random key; provably perfect secrecy, hopelessly impractical.
- [[Modular Arithmetic]] — clock arithmetic, the algebra beneath all of public-key crypto.
- [[Diffie–Hellman]] — agree on a shared secret over a public channel.
- [[RSA]] — encryption and signatures from the hardness of factoring.
- [[Public-Key Cryptography]] — the big idea: a public lock, a private key.
- [[Hash Function]] — a one-way fingerprint with the avalanche property.
- [[Digital Signature]] — sign with your private key, verify with your public one.

## See also

- [[Entropy]]
- [[One-Time Pad]]
- [[Public-Key Cryptography]]
