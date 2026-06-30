---
title: Data Compression
aliases: [compression, source coding]
tags: [information]
summary: Encoding data in fewer bits by removing redundancy — with entropy as the unbreakable floor.
weight: 30
---

# Data Compression

**Data compression** (formally, *source coding*) is the art of describing the same data in fewer bits. It works by exploiting **redundancy**: real data is not random. Letters in text follow predictable frequencies, pixels resemble their neighbors, audio samples vary smoothly. Wherever a source is predictable, a clever code can name its outcomes more cheaply than the naive fixed-width encoding.

There are two regimes. **Lossless** compression (ZIP, PNG, FLAC) reconstructs the original bit-for-bit by recoding redundancy. **Lossy** compression (JPEG, MP3) goes further by discarding detail the recipient will not miss. Information theory governs the lossless case exactly.

## Entropy is the floor

The central fact, Shannon's source-coding theorem, is brutal and beautiful: **no lossless code can use fewer than $H$ bits per symbol on average**, where $H$ is the [[Entropy]] of the source. You can approach that floor as closely as you like, but never beat it.

{{< note kind="key" title="The compression bound" >}}
For a source with entropy $H$ bits/symbol, any uniquely decodable code has average length $\bar{L} \ge H$. The redundancy you can remove is exactly the gap between the data's apparent size and its entropy. Once you reach $H$, what remains is *incompressible* — it looks, bit for bit, like pure randomness.
{{< /note >}}

A skewed [[Probability Distribution]] has low entropy and compresses well; a uniform one is already at the maximum and cannot be squeezed. This is why an already-compressed file barely shrinks when you ZIP it again — its redundancy is already gone.

## How codes beat the fixed-width baseline

Spend short codewords on frequent symbols and long ones on rare symbols, and the *average* length drops below the fixed-width $\log_2 n$. The chart contrasts a flat 3-bit-per-symbol encoding with a variable-length code matched to symbol frequencies — the same trick [[Huffman Coding]] makes optimal.

{{< chart type="bar" data="3,3,3,3,3,3" labels="A,B,C,D,E,F" title="Fixed-width code: every symbol costs 3 bits" ylabel="bits" color="#e0a458" caption="A naive code ignores frequencies and pays the worst case for everything." >}}

{{< chart type="bar" data="1,2,3,4,5,5" labels="A,B,C,D,E,F" title="Variable-length code: lengths matched to frequency" ylabel="bits" caption="Frequent symbols (A, B) get short codewords; rare ones (E, F) absorb the long codes. The frequency-weighted average can fall well below 3 bits." >}}

## The ultimate limit

Entropy bounds compression *for a given probability model*. But what is the shortest description of one specific string, with no model at all? That is its [[Kolmogorov Complexity]] — the length of the smallest program that prints it. A string is **incompressible** exactly when no program shorter than the string itself can generate it, which is the algorithmic definition of randomness. Practical compressors are finite, fast approximations to this uncomputable ideal.

{{< quiz question="A file is already compressed to its entropy. What happens if you run it through a lossless compressor again?" options="It roughly halves in size|It barely changes — its redundancy is already gone|It becomes uncompressible noise and is lost|It always grows to exactly twice the size" answer="2" explain="Once data sits at its entropy floor it has no remaining redundancy to exploit, so a second pass can shave essentially nothing (and may even add a few bytes of overhead)." >}}

## See also

- [[Entropy]]
- [[Huffman Coding]]
- [[Kolmogorov Complexity]]
