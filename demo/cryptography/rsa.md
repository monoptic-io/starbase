---
title: RSA
aliases: [rsa, rsa encryption]
tags: [cryptography]
summary: Public-key encryption and signatures built on the difficulty of factoring a product of two large primes.
weight: 60
---

# RSA

**RSA** — named for Rivest, Shamir, and Adleman, who published it in 1977 — was the first practical scheme to deliver the full promise of [[Public-Key Cryptography]]: a **public key** anyone can use to encrypt to you, and a **private key** only you can use to decrypt. Its security rests on a single asymmetry of [[Modular Arithmetic]]: multiplying two large primes is easy, but **factoring** their product back apart is, as far as anyone knows, impossibly slow.

## Building the keys

1. Pick two large primes $p$ and $q$ and multiply them: $n = p q$. This $n$ is public; its factors are not.
2. Compute $\varphi(n) = (p-1)(q-1)$ — secret.
3. Choose a public exponent $e$ coprime to $\varphi(n)$ (commonly $65537$).
4. Find the private exponent $d$ as the modular inverse of $e$: $\;ed \equiv 1 \pmod{\varphi(n)}$.

Your **public key** is the pair $(n, e)$; your **private key** is $d$. Encryption and decryption are each one modular exponentiation:

$$c = m^{e} \bmod n, \qquad m = c^{d} \bmod n.$$

The reason decryption undoes encryption is Fermat/Euler's theorem: $m^{ed} = m^{1 + k\varphi(n)} \equiv m \pmod{n}$.

## A tiny worked example

Real RSA uses primes hundreds of digits long. To *see* the gears turn, take absurdly small ones.

{{< note kind="note" title="RSA with p = 3, q = 11" >}}
- $n = 3 \times 11 = 33$, and $\varphi(n) = 2 \times 10 = 20$.
- Choose $e = 7$ (coprime to 20). The private $d$ solves $7d \equiv 1 \pmod{20}$, giving $d = 3$ (since $7 \times 3 = 21 \equiv 1$).
- **Public key** $(33, 7)$; **private key** $3$.

Encrypt the message $m = 4$:
$$c = 4^{7} \bmod 33 = 16384 \bmod 33 = 16.$$
Decrypt $c = 16$:
$$m = 16^{3} \bmod 33 = 4096 \bmod 33 = 4. \checkmark$$
The plaintext comes back exactly. Anyone with $(33,7)$ can encrypt; only the holder of $d = 3$ can reverse it.
{{< /note >}}

The catch for an attacker: to find $d$ they need $\varphi(n) = (p-1)(q-1)$, which means knowing $p$ and $q$ — i.e. **factoring $n$**. For $n = 33$ that is trivial. For a 2048-bit $n$ it is the wall.

## Why factoring is the whole game

Multiplying is a near-instant operation even on huge numbers; factoring the result is believed to take time that explodes with the size of $n$. The chart sketches that gulf in difficulty.

{{< chart type="bar" data="1,2,4,30,3000" labels="multiply,trial 64-bit,trial 256-bit,factor 512-bit,factor 1024-bit" title="Relative effort: forward vs. reverse (illustrative log-ish scale)" ylabel="hardness" caption="Going forward — multiplying p·q — is cheap at any size. Going backward — factoring n — becomes astronomically expensive as the key grows. That one-way gap is RSA's security." >}}

{{< note kind="warning" title="Textbook RSA is not safe RSA" >}}
The clean $m^e \bmod n$ shown here is **textbook RSA** and is insecure in practice: it is deterministic (identical messages give identical ciphertext) and malleable. Real deployments add randomized **padding** (OAEP) and never encrypt raw long messages directly — instead RSA wraps a short symmetric key. The math is the same; the engineering around it is what makes it secure.
{{< /note >}}

Because the public and private exponents are interchangeable, running the *private* key first and the *public* key to check produces not secrecy but a **[[Digital Signature]]** — proof that only the key's owner could have created the value.

{{< quiz question="In RSA, what secret must an attacker learn to compute the private exponent d from the public key (n, e)?" options="The message m|The factorization of n into p and q|The value of e|Nothing — d is public" answer="2" explain="d is the inverse of e modulo φ(n) = (p−1)(q−1). Computing φ(n) requires the prime factors of n, so breaking RSA reduces to factoring." >}}

## See also

- [[Modular Arithmetic]]
- [[Public-Key Cryptography]]
- [[Digital Signature]]
