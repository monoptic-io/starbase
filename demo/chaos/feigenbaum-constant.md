---
title: Feigenbaum Constant
aliases: [Feigenbaum delta, Feigenbaum number, period-doubling constant]
tags: [chaos]
summary: The universal ratio δ ≈ 4.6692 by which the gaps between successive period-doubling bifurcations shrink — the same number on the road to chaos in countless unrelated systems.
weight: 55
---

# Feigenbaum Constant

Turn the control knob of a [[Logistic Map]] and its steady state period-doubles: a 1-cycle becomes a 2-cycle, then 4, 8, 16, the doublings arriving faster and faster until the period is infinite and the motion is chaotic. The remarkable thing is not *that* this happens but *how fast*. The parameter windows between successive doublings shrink geometrically, and the ratio of one gap to the next homes in on a single number,

$$\delta = \lim_{n\to\infty}\frac{r_{n}-r_{n-1}}{r_{n+1}-r_{n}} = 4.6692016\ldots$$

This is the **Feigenbaum constant**. Mitchell Feigenbaum discovered in the 1970s that the *same* δ governs the period-doubling cascade in an enormous range of systems — dripping faucets, convecting fluids, nonlinear circuits, cardiac cells — provided only that the underlying map has a smooth quadratic maximum. The details of the system wash out; the *rate* at which order dissolves into chaos is **universal**, a constant of nature in the same spirit as $\pi$ or $e$. It is one of the cleanest pieces of evidence that [[Chaos|chaos]] has a rigid, quantitative structure rather than being mere messiness.

## The cascade, located numerically

We can pin the cascade down from first principles, with no data and no fitting — just the arithmetic of $x \mapsto r\,x(1-x)$. The trick is that a period-$2^{\,n-1}$ cycle loses [[Stability|stability]] *exactly* where its **multiplier** — the product of the slope $f'(x)=r(1-2x)$ taken around the cycle — passes through $-1$. That $r$ is the [[Bifurcation|bifurcation]] $r_n$ at which the period doubles. Solving the cycle equations with Newton's method and driving the multiplier to $-1$ with a secant step locates each $r_n$ to machine precision:

{{< claim check="feigenbaum-delta" asof="2026-06-30" >}}
The first six period-doubling thresholds are $r_1=3$ (period $1\to2$), $r_2=1+\sqrt6\approx3.4495$, $r_3\approx3.5441$, $r_4\approx3.5644$, $r_5\approx3.5688$, $r_6\approx3.5697$. The successive gap ratios $\delta_n$ march 4.75, 4.66, 4.668, 4.669 — already converging on Feigenbaum's $\delta \approx$ **{{< val check="feigenbaum-delta" field="delta" >}}** by the fifth doubling.

```python
#!/usr/bin/env python3
# A period-2^(n-1) cycle loses stability where its multiplier (product of
# f'(x_i) around the cycle) reaches -1; that r is the bifurcation b_n. Pin
# each b_n with a Newton-polished cycle + a secant on the multiplier, then
# form delta_n = (b_n - b_{n-1}) / (b_{n+1} - b_n) -> 4.6692016...
def f(r, x):  return r * x * (1.0 - x)
def fp(r, x): return r * (1.0 - 2.0 * x)
def cycle_mult(r, p, x):
    for _ in range(200):
        y, deriv = x, 1.0
        for _ in range(p):
            deriv *= fp(r, y); y = f(r, y)
        if deriv == 1.0: break
        dx = (y - x) / (deriv - 1.0); x -= dx
        if abs(dx) < 1e-15: break
    y, M = x, 1.0
    for _ in range(p):
        M *= fp(r, y); y = f(r, y)
    return x, M
def bifurcation(p, r_seed):
    x = 0.5
    for _ in range(200000): x = f(r_seed, x)          # land on the cycle
    r0, (x0, M0) = r_seed, cycle_mult(r_seed, p, x)
    r1, (x1, M1) = r_seed + 1e-4, cycle_mult(r_seed + 1e-4, p, x)
    for _ in range(100):                              # secant: multiplier -> -1
        if M1 == M0: break
        r2 = r1 - (M1 + 1.0) * (r1 - r0) / (M1 - M0)
        x2, M2 = cycle_mult(r2, p, x1)
        r0, M0, x0, r1, M1, x1 = r1, M1, x1, r2, M2, x2
        if abs(r1 - r0) < 1e-15: break
    return r1
periods = [1, 2, 4, 8, 16, 32]
b = [bifurcation(p, s) for p, s in zip(periods, [2.9,3.4,3.52,3.560,3.5675,3.5694])]
print("n  period  r_n (bifurcation)      gap             delta_n")
for i in range(len(b)):
    gap   = f"{b[i]-b[i-1]:.10f}" if i >= 1 else ""
    delta = f"{(b[i-1]-b[i-2])/(b[i]-b[i-1]):.5f}" if i >= 2 else ""
    print(f"{i+1}  {periods[i]:6d}  {b[i]:.10f}   {gap:>13}   {delta:>9}")
best = (b[-2] - b[-3]) / (b[-1] - b[-2])
print(f"best ratio delta_5 = {best:.5f}  ->  delta ~ 4.669  (Feigenbaum 4.6692016...)")
```
{{< /claim >}}

The convergence is not monotone-from-one-side or accidental: it is the geometric signature of the cascade. Each gap is very nearly $\delta^{-1}$ times the one before, so the doublings accumulate at a *finite* parameter rather than running off to $r=4$.

## Where the doublings run out

Because the gaps shrink by a constant factor, the infinite tally of bifurcations converges. Summing the remaining geometric tail of gaps past the ones we computed gives the **accumulation point** $r_\infty$ — the exact parameter where the period first becomes infinite and the orbit, no longer settling onto any finite cycle, lands on a [[Strange Attractor|strange attractor]] of [[Fractal|fractal]] (Cantor-set) structure:

{{< claim check="feigenbaum-accumulation" asof="2026-06-30" >}}
Extrapolating past the sixth doubling, the cascade accumulates at $r_\infty \approx$ **{{< val check="feigenbaum-accumulation" field="r_inf" >}}**. Cross this threshold and the logistic map exhibits [[Sensitive Dependence on Initial Conditions|sensitive dependence on initial conditions]]: nearby orbits separate exponentially and long-term prediction collapses.

```python
#!/usr/bin/env python3
# The bifurcations b_n pile up geometrically (gaps shrink by ~delta each step).
# Locate the last two we can resolve, then sum the remaining geometric tail of
# gaps to reach the accumulation point r_inf, where the period goes infinite.
# (bifurcation() as in the delta check)
b5 = bifurcation(16, 3.5675)
b6 = bifurcation(32, 3.5694)
delta = 4.669201609102990
last_gap = b6 - b5
tail = (last_gap / delta) / (1.0 - 1.0 / delta)   # b7-b6 + b8-b7 + ...
r_inf = b6 + tail
print(f"b5 (period 16 -> 32) = {b5:.10f}")
print(f"b6 (period 32 -> 64) = {b6:.10f}")
print(f"last gap b6-b5       = {last_gap:.10f}")
print(f"geometric tail / delta = {tail:.10f}")
print(f"accumulation r_inf   = {r_inf:.5f}  (true 3.5699456...)")
```
{{< /claim >}}

Beyond $r_\infty$ lies the chaotic regime, shot through with narrow periodic **windows** — but the universal number $\delta$ has already done its work by the time we arrive.

## Why universal?

The depth of Feigenbaum's discovery is that $\delta$ does not depend on the logistic map at all. Any smooth one-dimensional map with a single quadratic hump — $\sin$, a cubic, a tent rounded at the tip — period-doubles its way to chaos with the *same* δ and the *same* spatial scaling constant $\alpha\approx2.5029$. The explanation is a **renormalization** argument borrowed from statistical physics: zoom in on the map near its peak after each doubling and rescale, and the rescaled maps converge to a fixed shape, a fixed point of a "doubling operator." δ is simply the rate at which that operator stretches parameter space — an eigenvalue, blind to which particular map you started from. That is why a [[Bifurcation|bifurcation]] cascade in a dripping tap and one in the [[Logistic Map]] share a digit-for-digit constant.

## See also

- [[Logistic Map]]
- [[Bifurcation]]
- [[Strange Attractor]]
- [[Sensitive Dependence on Initial Conditions]]
- [[Chaos]]
