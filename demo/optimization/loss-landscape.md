---
title: Loss Landscape
aliases: [loss surface, objective function, cost function]
tags: [optimization]
summary: The terrain an optimizer walks — a surface over parameter space whose height is the loss, dotted with minima, saddles, and plateaus.
weight: 20
---

# Loss Landscape

When you train a model, every possible setting of its parameters $\theta$ earns a score: the **loss** $L(\theta)$, a single number measuring how badly it does. Imagine plotting that score as a *height* above the space of all parameter settings. The result is a surface — the **loss landscape** — and optimization is nothing more than a hike across it toward the lowest valley. The objective function, the cost function, the loss surface: these are all the same terrain seen from different fields.

The landscape's *shape* decides whether learning is easy or agonizing. A smooth bowl is a gift; a jagged range of peaks, ditches, and dead-flat plains is a nightmare. Understanding the features below is how you reason about why [[Gradient Descent|descent]] sometimes sails to the answer and sometimes gets hopelessly stuck.

## The features of the terrain

- **Global minimum** — the genuinely lowest point; the best the model can possibly do.
- **Local minimum** — a valley that is lowest *in its neighborhood* but not overall. Descent that wanders in can never climb out on gradient alone.
- **Saddle point** — a spot that is downhill in one direction and uphill in another, like a mountain pass. The gradient vanishes there, so naive descent slows to a crawl even though it is not a minimum.
- **Plateau** — a vast near-flat region where the gradient is tiny and progress stalls.
- **Basin of attraction** — the catchment of a minimum: every start inside it rolls to the same bottom. In the language of [[Dynamical System|dynamics]], the minima are [[Attractor|attractors]] and the basins are their domains.

{{< note kind="warning" title="Saddles, not minima, are the real enemy" >}}
In a landscape of many dimensions, a true local minimum requires the surface to curve *upward in every direction at once* — increasingly unlikely as dimensions grow. **Saddle points**, where some directions go up and others down, are exponentially more common, and they are what bog down high-dimensional training far more often than spurious minima.
{{< /note >}}

## A surface you can read

The plot below is a one-dimensional slice through a rugged loss landscape. Trace it with your eye: the deepest dip is the global minimum, the shallower dip is a local minimum that could trap a careless optimizer, and the gentle rise between them is a small ridge to be crossed.

{{< plot fn="0.5*x*x + Math.sin(2.4*x) + 0.4*Math.cos(4.1*x)" xmin="-4" xmax="4" title="A 1-D slice of a non-convex loss L(θ)" caption="Valleys are minima, the overall upward sweep is a confining bowl, and the little wiggles are local traps. Descent from the left settles in a different valley than descent from the right." >}}

Where the curve is a single clean bowl, optimization is trivial — that special, lucky case is [[Convexity|convexity]]. Real models are usually **non-convex**, full of the bumps you see above, which is exactly why momentum, annealing, and clever initialization exist.

{{< quiz question="At a saddle point of a loss landscape, what is true of the gradient?" options="It is large and points to the global minimum|It is zero, yet the point is not a minimum|It does not exist|It points straight uphill" answer="2" explain="At a saddle the gradient vanishes (it is a stationary point), so gradient descent stalls — but the surface curves downward in some direction, so it is not a minimum. Escaping requires noise, momentum, or curvature information." >}}

## See also

- [[Gradient Descent]]
- [[Convexity]]
- [[Attractor]]
