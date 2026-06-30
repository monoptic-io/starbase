---
title: Benford's Law
aliases: [first-digit law, Benford, leading-digit law]
tags: [probability]
summary: In many natural collections of numbers the leading digit is 1 about 30% of the time, not 11% — a logarithmic bias that holds across populations, areas, and physical constants.
weight: 65
---

# Benford's Law

Pick a number at random from a long list of real-world quantities — the populations of countries, the lengths of rivers, the areas of nations, the figures buried in a tax return — and look only at its **first digit**. Naively you might expect each of the digits 1 through 9 to lead about $\tfrac{1}{9} \approx 11\%$ of the time. They do not. The digit **1** leads roughly **30%** of the time, **2** about **18%**, and the frequencies fall away smoothly to **9**, which leads under **5%** of the time. This lopsided pattern is **Benford's law**, also called the *first-digit law*.

The rule is precise. The probability that the leading digit equals $d$ is

{{< eq number="1" >}}
P(d) = \log_{10}\!\left(1 + \frac{1}{d}\right), \qquad d \in \{1,2,\ldots,9\}.
{{< /eq >}}

Because $\sum_{d=1}^{9} \log_{10}(1 + 1/d) = \log_{10}(10) = 1$, this is a genuine [[Probability Distribution]] over the nine possible leading digits. Plugging in $d=1$ gives $\log_{10} 2 = 0.301$ — the famous "30%."

{{< chart type="bar" data="30.1 17.6 12.5 9.7 7.9 6.7 5.8 5.1 4.6" labels="1,2,3,4,5,6,7,8,9" title="Benford's law: probability of each leading digit" xlabel="leading digit" ylabel="probability (%)" >}}

## Why a logarithm?

The deep reason is **scale invariance**. If a collection of numbers obeys some universal first-digit law at all, that law cannot care whether the quantities are measured in miles or kilometres, dollars or euros — multiplying every value by a constant must leave the digit distribution unchanged. The *only* distribution with that property is the logarithmic one above.

The intuition: lay the numbers out on a logarithmic ruler, where the distance from 1 to 2 is much wider than from 8 to 9. A quantity that grows by steady *percentages* — the way populations, prices, and bank balances do — spreads itself **uniformly across that logarithmic ruler**. The stretch of the ruler labelled with a leading 1 occupies $\log_{10} 2 \approx 30\%$ of each decade; the stretch labelled 9 occupies only $\log_{10}(10/9) \approx 4.6\%$. The digits inherit the widths of their logarithmic homes. This is why the law loves data that spans **many orders of magnitude** and arises from multiplicative growth — the same engine behind a [[Scale-Free Network]]'s power-law degrees.

{{< note kind="key" title="When to expect Benford" >}}
- Data spanning **several orders of magnitude** (populations from thousands to billions).
- Quantities that grow **multiplicatively** rather than additively.
- Numbers **not** artificially bounded or assigned (heights of adults, phone numbers, and lottery draws do *not* obey it).
- A large enough sample that the [[Law of Large Numbers]] lets observed frequencies settle onto the predicted ones.
{{< /note >}}

## Testing it on real data

Talk is cheap; the law is checkable. The dataset `data/world-figures.csv` holds **292** genuine reference figures — the 2023 population and the land area in km² of 146 countries. These span from Luxembourg's 2,586 km² to Russia's 17 million km², and from Iceland's 375 thousand people to China's 1.4 billion: many orders of magnitude, exactly the regime where Benford should bite. Stripping each value down to its leading digit and tallying gives:

{{< claim value="maxdev 1.7%" check="benford-leading-digits" source="data/world-figures.csv" asof="2026-06-30" >}}
Across all **292** figures the observed leading-digit distribution tracks Benford's prediction strikingly well — a chi-square statistic of just **3.22** (well under the 15.5 critical value for 8 degrees of freedom at the 5% level) and a largest single-digit gap of only **1.7** percentage points.
```sh
# evidence/benford-leading-digits/{inputs: data/world-figures.csv, run}
LC_ALL=C awk -F, '
NR>1 { v=$3; sub(/^[^1-9]*/, "", v); d = substr(v,1,1) + 0
       if (d>=1 && d<=9) { c[d]++; n++ } }
END { split("0.301 0.176 0.125 0.097 0.079 0.067 0.058 0.051 0.046", b, " ")
      print "digit  count  observed  benford"; chi2=0; maxdev=0
      for (d=1; d<=9; d++) { obs=c[d]/n; ex=n*b[d]
        chi2 += (c[d]-ex)*(c[d]-ex)/ex
        dev=obs*100-b[d]*100; if(dev<0)dev=-dev; if(dev>maxdev)maxdev=dev
        printf "%d      %5d   %6.1f%%   %5.1f%%\n", d, c[d], obs*100, b[d]*100 }
      printf "N=%d  chi2=%.2f  maxdev=%.1f%%\n", n, chi2, maxdev }' world-figures.csv
```
```result
digit  count  observed  benford
1         85     29.1%    30.1%
2         53     18.2%    17.6%
3         37     12.7%    12.5%
4         25      8.6%     9.7%
5         26      8.9%     7.9%
6         23      7.9%     6.7%
7         12      4.1%     5.8%
8         16      5.5%     5.1%
9         15      5.1%     4.6%
N=292  chi2=3.22  maxdev=1.7%
```
{{< /claim >}}

The single most visible signature of the law is the dominance of the digit 1:

{{< claim value="29.1%" check="benford-digit-one" source="data/world-figures.csv" asof="2026-06-30" >}}
**85 of the 292 figures — 29.1%** — begin with the digit 1, almost exactly Benford's predicted $\log_{10} 2 = 30.1\%$ and nearly triple the **11.1%** a "uniform" guess would give.
```sh
# evidence/benford-digit-one/{inputs: data/world-figures.csv, run}
LC_ALL=C awk -F, '
NR>1 { v=$3; sub(/^[^1-9]*/, "", v); d=substr(v,1,1)+0
       if (d>=1 && d<=9) { n++; if (d==1) ones++ } }
END { printf "%d of %d figures (%.1f%%) lead with digit 1\n", ones, n, 100*ones/n }' world-figures.csv
```
```result
85 of 292 figures (29.1%) lead with digit 1
```
{{< /claim >}}

The fit is not a coincidence of these particular numbers — pool populations and areas measured in different units and the same shape appears, exactly as scale invariance demands.

## Where it is used

Because *honest* natural data obeys Benford so reliably, **departures from it are suspicious**. Forensic accountants and auditors run Benford tests on ledgers, expense reports, and tax filings: fabricated figures, invented by people who unconsciously spread their leading digits too evenly, tend to show a deficit of 1s and a surplus of middle digits. The same test has been turned on reported election tallies, scientific datasets, and macroeconomic statistics as a cheap first-pass screen for manipulation. It is not proof — some clean datasets fail the law and some doctored ones pass — but a strong Benford violation is a flag worth pulling.

The law also reaches into pure mathematics: the leading digits of the powers of 2, of the Fibonacci numbers, and of many sequences tied to the [[Prime Number]] counting function are Benford-distributed, again because their logarithms fill the unit interval uniformly.

## See also

- [[Probability Distribution]]
- [[Law of Large Numbers]]
- [[Scale-Free Network]]
- [[Entropy]]
