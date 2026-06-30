---
title: Channel Capacity
aliases: [shannon capacity, noisy channel coding]
tags: [information]
summary: The maximum rate at which information can cross a noisy channel with arbitrarily small error.
weight: 60
---

# Channel Capacity

**Channel capacity** is the top speed of a communication channel — the highest rate, in bits per use, at which you can send information through noise and still recover it essentially perfectly. Shannon's **noisy-channel coding theorem** (1948) is the stunning result that this top speed $C$ is not just a practical limit but a sharp threshold:

- Transmit at any rate **below** $C$ and, with a clever enough [[Error-Correcting Code]], you can drive the error probability **arbitrarily close to zero**.
- Transmit **above** $C$ and reliable communication is **impossible**, no matter how much redundancy you add.

Before Shannon, engineers assumed noise forced a grim trade-off: go faster, get more errors. Shannon showed that below capacity you can have *both* speed and reliability — you just need long, well-designed codes.

## Capacity is maximized mutual information

The capacity is defined through [[Mutual Information]]. The channel turns an input $X$ into a noisy output $Y$; the information that actually gets through is $I(X;Y)$. Capacity is the best you can do over all ways of choosing the input:

$$C = \max_{p(x)} I(X;Y).$$

{{< note kind="key" title="A worked channel" >}}
For a **binary symmetric channel** that flips each bit independently with probability $p$, the capacity is

$$C = 1 - H(p) \;=\; 1 + p\log_2 p + (1-p)\log_2(1-p) \quad \text{bits/use.}$$

A clean channel ($p = 0$) gives the full 1 bit. A maximally noisy one ($p = \tfrac12$) gives $C = 0$ — the output is independent of the input, so nothing crosses. The cost of noise is exactly the binary [[Entropy]] $H(p)$ it injects.
{{< /note >}}

## How noise eats the rate

The plot shows $C = 1 - H(p)$ for the binary symmetric channel. Capacity is highest for a quiet channel, plunges as the flip probability climbs toward $\tfrac12$, then — strikingly — **recovers** for $p > \tfrac12$, because a channel that almost always flips is just as informative as one that almost never does (you simply invert its output).

{{< plot fn="1-(-x*Math.log2(x)-(1-x)*Math.log2(1-x))" xmin="0.001" xmax="0.999" ymin="0" ymax="1.05" title="Capacity of a binary symmetric channel, C = 1 − H(p)" caption="Reliable bits per channel use versus bit-flip probability p. Zero at p = 0.5, where input and output are independent." >}}

## Bandwidth, power, and the analog version

For continuous channels Shannon gave an equally famous formula. A channel of bandwidth $B$ hertz with signal-to-noise ratio $\mathrm{SNR}$ has capacity

$$C = B\,\log_2(1 + \mathrm{SNR}) \quad \text{bits/second.}$$

Capacity grows linearly with the slice of [[Spectrum]] you occupy and only logarithmically with power — which is why grabbing more **bandwidth** is the cheaper road to higher data rates, and why every modem, Wi-Fi link, and deep-space probe is engineered right up against this line.

{{< quiz question="A binary symmetric channel flips each bit with probability p = 0.5. What is its capacity?" options="1 bit per use|0.5 bits per use|0 bits per use|Infinite" answer="3" explain="At p = 0.5 the output is statistically independent of the input, so I(X;Y) = 0 and no information gets through. Capacity is zero." >}}

## See also

- [[Mutual Information]]
- [[Error-Correcting Code]]
- [[Spectrum]]
