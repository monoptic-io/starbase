---
title: Substitution Cipher
aliases: [substitution cipher, monoalphabetic cipher]
tags: [cryptography]
summary: Replace each letter by any other under a fixed scrambled alphabet — an enormous keyspace that frequency analysis still cracks.
weight: 20
---

# Substitution Cipher

A **substitution cipher** replaces every letter with another according to a fixed but *arbitrary* permutation of the alphabet. Unlike the [[Caesar Cipher]], which only rotates the alphabet, here the mapping can be any scramble at all:

```
plain:  A B C D E F G H I J K L M N O P Q R S T U V W X Y Z
cipher: Q W E R T Y U I O P A S D F G H J K L Z X C V B N M
```

So `HELLO` becomes `ITSSG`. Decryption just runs the table backwards. Because each plaintext letter always maps to the same ciphertext letter, this is called a **monoalphabetic** cipher.

## A staggering keyspace

How many keys are there? Any of the 26 letters can map to the first slot, any of the remaining 25 to the next, and so on — a full permutation:

$$26! = 403{,}291{,}461{,}126{,}605{,}635{,}584{,}000{,}000 \approx 4 \times 10^{26}.$$

That is about **88 bits** of key, more than $10^{26}$ possibilities. Brute force is utterly hopeless — you could try a billion keys a second since the Big Bang and not scratch it. Surely *this* one is secure?

{{< note kind="warning" title="A big keyspace is not security" >}}
It is not. A huge keyspace defeats brute force, but says nothing about *cleverer* attacks. The substitution cipher hides **which** letter is which, yet it faithfully preserves **how often** each one appears. That leak is all an attacker needs — and it has nothing to do with trying keys one by one.
{{< /note >}}

## The crack: frequency analysis

In English, letters are wildly uneven. **E** alone is about 12.7% of text; **T**, **A**, **O**, **I**, **N** follow; **J**, **Q**, **X**, **Z** are rare. A monoalphabetic cipher renames the letters but cannot change their counts: whatever symbol stands for *E* will be the most common symbol in the ciphertext. Match the ciphertext's histogram against English's, throw in common pairs (`TH`, `HE`) and one-letter words (`A`, `I`), and the message unravels.

{{< chart type="bar" data="8.2,1.5,2.8,4.3,12.7,2.2,2.0,6.1,7.0,0.15,0.77,4.0,2.4,6.7,7.5,1.9,0.095,6.0,6.3,9.1,2.8,0.98,2.4,0.15,2.0,0.074" labels="A,B,C,D,E,F,G,H,I,J,K,L,M,N,O,P,Q,R,S,T,U,V,W,X,Y,Z" title="English letter frequencies (%)" ylabel="% of letters" xlabel="letter" caption="The fingerprint of English. A substitution cipher relabels the bars but never flattens them — so the tallest ciphertext bar is almost certainly E." >}}

This is exactly why the cipher is weak despite its colossal keyspace: it has high **key entropy** but leaks almost everything about the plaintext's structure. The connection to [[Entropy]] is the lesson — what matters is not how many keys exist, but how much the ciphertext still tells you about the message. The only way to truly plug the leak is to make the key as informative as the message itself, which leads to the [[One-Time Pad]].

{{< quiz question="Why is a substitution cipher breakable even though 26! keys make brute force impossible?" options="The keys are easy to guess|It preserves each letter's frequency, so the plaintext's statistics leak through|26! is actually a small number|It only works on short messages" answer="2" explain="Monoalphabetic substitution renames letters but keeps their counts. Frequency analysis exploits that statistical leak directly, never touching the keyspace." >}}

## See also

- [[Caesar Cipher]]
- [[Entropy]]
- [[One-Time Pad]]
