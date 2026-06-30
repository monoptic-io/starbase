---
title: Digital Signature
aliases: [digital signature, signing]
tags: [cryptography]
summary: Sign a message's hash with your private key so anyone can verify, with your public key, that you wrote it and nobody altered it.
weight: 90
---

# Digital Signature

A **digital signature** is the cryptographic counterpart of a handwritten signature — but far stronger, because it is mathematically bound to the *exact* document it signs. It proves three things at once:

- **Authenticity** — the message really came from the claimed sender.
- **Integrity** — not one bit has been altered since signing.
- **Non-repudiation** — the signer cannot later deny having signed it.

It is [[Public-Key Cryptography]] run in reverse. Encryption uses the *public* key to lock and the *private* key to unlock. Signing swaps the roles: the **private key signs**, and the matching **public key verifies**. Only the key's owner can produce a signature, but *anyone* can check it.

## How signing works

You never sign the whole message directly — that would be slow and unwieldy. Instead you sign its **[[Hash Function|hash]]**, a short fingerprint of the document.

{{< note kind="key" title="Sign the hash, verify the hash" >}}
**Signing** (private key $d$):
1. Compute the digest $h = H(m)$ with a [[Hash Function]].
2. Transform $h$ with your private key to produce the signature $\sigma$. In [[RSA]] terms, $\sigma = h^{d} \bmod n$.
3. Send $(m, \sigma)$.

**Verifying** (public key $e$):
1. Recompute the digest $h = H(m)$ from the received message.
2. Undo the signature with the signer's public key, $h' = \sigma^{e} \bmod n$.
3. Accept **iff** $h' = h$. A match could only have been produced by the holder of $d$.
{{< /note >}}

Because the signature is tied to $H(m)$, the [[Hash Function|avalanche effect]] does the heavy lifting: change a single character of $m$ and its hash changes completely, so the recomputed $h$ no longer matches $h'$ and verification fails. The signature guards every bit of the document.

## Why the hash matters

Hashing first is not just an optimization — it is what makes signatures safe and practical.

- **Speed.** Public-key operations are expensive; signing 32 bytes of digest instead of a gigabyte file is enormously faster.
- **Integrity coverage.** A good [[Hash Function]] makes it infeasible to find a *different* message with the same digest, so an attacker cannot lift your signature onto forged content.
- **Identity binding.** Verification ties the message to a specific public key; certificate authorities then vouch that the key belongs to a real person or server.

This is the missing piece for [[Diffie–Hellman]], which agrees a shared key but cannot by itself confirm *who* is on the other end. Signing the exchange defeats the man-in-the-middle and is exactly how TLS authenticates the websites you visit — the padlock in your browser is a verified digital signature.

{{< quiz question="To verify a digital signature, which key does the recipient use?" options="The signer's private key|The signer's public key|Their own private key|A fresh random key" answer="2" explain="The signer signs with their private key; anyone verifies with the matching public key. Only the private-key holder could have produced a signature that the public key validates." >}}

## See also

- [[Public-Key Cryptography]]
- [[Hash Function]]
- [[RSA]]
