---
title: Linear Algebra
tags: [linear-algebra]
summary: The mathematics of vectors and the linear transformations that move them — the quiet engine beneath stability analysis, Fourier methods, and network ranking.
weight: 45
---

# Linear Algebra

**Linear algebra** is the study of [[Vector|vectors]] and the [[Linear Transformation|linear transformations]] that act on them. Its central insight is that an enormous range of operations — rotating a shape, blurring an image, advancing a physical system by one timestep, ranking web pages — are all the *same kind* of object: a [[Matrix]] that stretches and shears space while keeping straight lines straight and the origin fixed. Once a problem is written this way, a small, reusable toolkit takes over.

That toolkit is why linear algebra sits underneath so much of the rest of this knowledge base. The [[Stability]] of an equilibrium is decided by the eigenvalues of a linearized matrix. A [[Fourier Series]] is a change of [[Basis]] into sinusoids. [[PageRank]] is the dominant eigenvector of a link matrix. Learn the vocabulary once and these connections light up everywhere.

## The arc of the subject

We build up in three movements:

{{< columns count="2" >}}
**Objects** — a [[Vector]] is magnitude and direction; a [[Matrix]] packages a transformation as a grid of numbers; the [[Dot Product]] measures angle and projection; a [[Basis]] is the coordinate frame that turns geometry into numbers.

**Actions** — a [[Linear Transformation]] warps space; [[Matrix Multiplication]] composes two warps into one; the [[Determinant]] reports how much area or volume the warp scales by.
{{< /columns >}}

**Structure** — the deepest results expose the hidden skeleton of a transformation: [[Eigenvalues and Eigenvectors]] are the directions it merely stretches, and the [[Singular Value Decomposition]] factors *any* matrix into a rotation, a scaling, and another rotation.

## Why it underpins everything else

{{< note kind="key" title="One idea, many disguises" >}}
Linearize a hard nonlinear problem and you get a matrix. The eigenvalues of that matrix then tell you almost everything: whether a [[Fixed Point]] is [[Stability|stable]], what the normal modes of [[Coupled Oscillators]] are, which frequencies a signal contains, and where a random walker on a [[Graph]] spends its time. The pages here are the foundation; the payoff is scattered across the whole field guide.
{{< /note >}}

## Start here

- [[Vector]] — magnitude, direction, and components.
- [[Matrix]] — a grid of numbers that is secretly a transformation.
- [[Linear Transformation]] — how a matrix warps the plane.
- [[Matrix Multiplication]] — composing transformations.
- [[Determinant]] — the area/volume scale factor.
- [[Dot Product]] — projection, angle, and orthogonality.
- [[Basis]] — coordinates depend on the axes you pick.
- [[Eigenvalues and Eigenvectors]] — the directions a map only stretches.
- [[Singular Value Decomposition]] — rotate, scale, rotate.

## See also

- [[Stability]]
- [[Fourier Series]]
- [[PageRank]]
