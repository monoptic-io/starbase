---
title: Matrix
aliases: [matrices]
tags: [linear-algebra]
summary: A rectangular grid of numbers that doubles as a recipe for transforming space.
weight: 20
---

# Matrix

A **matrix** is a rectangular grid of numbers, written with rows and columns:

$$A = \begin{bmatrix} a & b \\ c & d \end{bmatrix}.$$

On the surface it is just bookkeeping — a table of coefficients. Its real meaning is that **a matrix is a [[Linear Transformation]]**: a rule for moving every [[Vector]] in space at once. Feeding a vector $\vec v = (x, y)$ to the matrix produces a new vector $A\vec v$, and the whole power of the subject comes from reading the grid of numbers and the geometric warp it encodes as two faces of one thing.

## How a matrix eats a vector

Matrix–vector multiplication is a weighted sum of the matrix's **columns**:

{{< eq number="1" >}}
A\vec v = \begin{bmatrix} a & b \\ c & d \end{bmatrix}\begin{bmatrix} x \\ y \end{bmatrix}
= x\begin{bmatrix} a \\ c \end{bmatrix} + y\begin{bmatrix} b \\ d \end{bmatrix}.
{{< /eq >}}

This is the key to *seeing* a matrix. The first column is where the matrix sends the basis vector $(1,0)$; the second column is where it sends $(0,1)$. Every other vector is carried along as the same linear combination of those two landing spots. So you can read a $2\times2$ matrix at a glance: its columns are the images of the coordinate axes.

## Special matrices

- The **identity** $\begin{bmatrix} 1 & 0 \\ 0 & 1 \end{bmatrix}$ leaves the axes — and therefore every vector — exactly where they are.
- A **diagonal** matrix $\begin{bmatrix} s & 0 \\ 0 & t \end{bmatrix}$ stretches the $x$-axis by $s$ and the $y$-axis by $t$.
- A **rotation** $\begin{bmatrix} \cos\theta & -\sin\theta \\ \sin\theta & \cos\theta \end{bmatrix}$ turns every vector by $\theta$ without changing its length.

{{< note kind="tip" title="Shape rules" >}}
A matrix with $m$ rows and $n$ columns maps $n$-dimensional input to $m$-dimensional output. A square ($n \times n$) matrix keeps you in the same space, which is why those are the ones with [[Determinant|determinants]] and [[Eigenvalues and Eigenvectors|eigenvectors]].
{{< /note >}}

## See also

- [[Linear Transformation]]
- [[Matrix Multiplication]]
- [[Determinant]]
