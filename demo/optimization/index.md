---
title: Optimization & Learning
tags: [optimization]
summary: How a system searches a landscape of possibilities for its best point — the common thread linking gradient descent, evolution, and learning machines.
weight: 65
---

# Optimization & Learning

Almost every learning system, fitted model, and engineered design is secretly answering the same question: **out of all the choices I could make, which one is best?** Phrase the quality of a choice as a number — call it a *loss* you want small, or a *fitness* you want large — and the space of all choices becomes a [[Loss Landscape|landscape]] of hills and valleys. **Optimization** is the art of searching that terrain for its lowest point, and **learning** is what we call it when the terrain is carved out by data.

The remarkable thing is how few ideas you need. Stand somewhere on the landscape, feel which way is downhill, and take a step. Repeat. That single recipe — [[Gradient Descent]] — trains everything from a one-line regression to a billion-parameter network. Everything else in this section is a refinement of it: how to step faster, how to escape traps, how to use curvature, and how to wire up a model whose landscape is worth descending.

{{< note kind="key" title="Descent is a dynamical system" >}}
Gradient descent is not just an algorithm — it is a [[Dynamical System]]. Each update $\theta \leftarrow \theta - \eta\nabla L(\theta)$ is a discrete map on parameter space, the continuous limit $\dot\theta = -\nabla L$ is a gradient flow, and the points it settles onto are exactly the [[Attractor|attractors]] of that flow: the **minima** of the loss. To optimize is to let a trajectory fall into a basin of attraction.
{{< /note >}}

## The landscape and how to move on it

- [[Gradient Descent]] — step downhill along $-\nabla L$, the workhorse of all learning.
- [[Loss Landscape]] — the terrain itself: minima, saddles, plateaus, and basins.
- [[Convexity]] — the lucky case of a single bowl, where descent cannot fail.
- [[Momentum]] — accumulate velocity to glide through ravines and over small bumps.
- [[Newton's Method]] — use curvature to leap, not crawl, toward the bottom.

## Searching when the landscape fights back

- [[Simulated Annealing]] — occasionally step *uphill* to escape local minima, cooling as you go.
- [[Genetic Algorithm]] — evolve a population by mutation and selection, no gradient required.

## Machines that learn

- [[Perceptron]] — a single weighted threshold that learns to draw a separating line.
- [[Neural Network]] — stacked layers that bend straight boundaries into curved ones.
- [[Backpropagation]] — the chain rule that sends error backward to teach every weight.

## See also

- [[Dynamical System]]
- [[Attractor]]
- [[Monte Carlo Method]]
