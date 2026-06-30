---
title: Convexity
aliases: [convex, convex optimization]
tags: [optimization]
summary: The property that a function is a single bowl with no false bottoms — the case where gradient descent is guaranteed to find the one true minimum.
weight: 30
---

# Convexity

A function is **convex** if it curves the same way everywhere — like a bowl. Pin down two points on its graph, draw the straight chord between them, and the function never pokes above that chord. Formally, for any two inputs $a, b$ and any blend $t \in [0,1]$,

$$L\big(t\,a + (1-t)\,b\big) \;\le\; t\,L(a) + (1-t)\,L(b).$$

This one inequality has an enormous consequence: **a convex function has no local minima other than the global one.** There is a single bottom, and *every* downhill direction eventually leads there. For an optimizer this is paradise — [[Gradient Descent|gradient descent]] cannot get stuck, cannot be fooled by a false valley, and (with a sensible step size) is guaranteed to converge to the best possible answer. This is why a whole field, **convex optimization**, exists: when you can phrase a problem convexly, you can essentially declare it solved.

## One bowl versus many

Contrast the two functions below. The smooth parabola is convex: a single basin, and a ball released anywhere slides to the same point. The wiggly one is **non-convex**: same overall upward sweep, but riddled with local dips where descent can come to rest short of the true minimum. Real learning problems — especially [[Neural Network|neural networks]] — live in the second picture.

{{< plot fn="0.4*x*x;;0.4*x*x + 1.3*Math.sin(3*x)" xmin="-4" xmax="4" title="Convex (smooth bowl) vs non-convex (bowl + ripples)" caption="Left curve: convex — one minimum, descent always wins. Right curve: non-convex — the same bowl plus sinusoidal ripples create local minima that can trap a gradient method." >}}

{{< note kind="key" title="Why convexity tames optimization" >}}
A twice-differentiable function of one variable is convex exactly when $L''(x) \ge 0$ everywhere — it never curves downward. In many dimensions the analog is that the **Hessian** (the matrix of second derivatives) is positive semidefinite: it has no negative [[Eigenvalues and Eigenvectors|eigenvalues]], so there is no direction of downward curvature in which a saddle or false minimum could form.
{{< /note >}}

## When the world is not convex

Most interesting models are non-convex, and that is the source of nearly every training headache: sensitivity to initialization, the need for [[Momentum|momentum]] and random restarts, and the gnawing doubt about whether you found *the* answer or just *an* answer. The practical wisdom is encouraging, though — in very high dimensions the many minima of a large network tend to be nearly as good as one another, so descent usually lands somewhere fine even without convexity's guarantee.

{{< quiz question="Why is gradient descent guaranteed to find the global minimum of a convex function?" options="Convex functions have no minimum at all|A convex function has only one minimum, so any local minimum is the global one|The gradient is always zero|Convex functions are always linear" answer="2" explain="Convexity rules out local minima distinct from the global minimum. Since descent halts only at a stationary point and the sole stationary minimum is global, it cannot get trapped anywhere worse." >}}

## See also

- [[Gradient Descent]]
- [[Loss Landscape]]
- [[Newton's Method]]
