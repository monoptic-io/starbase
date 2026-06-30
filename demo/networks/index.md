---
title: Graph Theory & Networks
aliases: [graphs, graph theory and networks, networks]
tags: [networks]
summary: Nodes joined by edges are the universal skeleton of connection — friendships, roads, neurons, and the web all reduce to the same simple picture.
weight: 50
---

# Graph Theory & Networks

Draw some dots. Join a few with lines. You have just built a **graph** — and with it a model flexible enough to describe friendships, road maps, food webs, the wiring of the brain, the citation trails of science, and the hyperlinks of the web. A graph keeps only what matters about a system of relationships: *who is connected to whom*. Strip away everything else and a startling amount of structure survives.

This section is a tour of that structure. We start with the [[Graph]] itself — vertices and edges, directed or weighted — and the two ways to walk one: [[Breadth-First Search]] in patient layers and [[Depth-First Search]] plunging deep. We add weights and ask for the cheapest route with [[Dijkstra's Algorithm]], or the cheapest way to connect everything at once with a [[Minimum Spanning Tree]]. Then we step back and ask *which nodes matter*: [[Centrality]] measures importance many ways, and [[PageRank]] answers it with a beautiful trick borrowed from linear algebra.

## Why networks belong here

Networks are [[Complex Systems]] in their purest combinatorial form. There are no forces and no differential equations — just nodes and edges — yet the same lessons return: local connections produce global behavior no single node contains. Two structural patterns dominate the real world. [[Small-World Network]]s are tightly clustered yet astonishingly shallow ("six degrees of separation"), while [[Scale-Free Network]]s grow lopsided hubs because the rich get richer. Both arise from rules you can state in a sentence.

{{< note kind="key" title="A graph is a relation made visible" >}}
Everything in this section is built from one object: a set of **nodes** and a set of **edges** between them. Algorithms *traverse* it, metrics *score* it, and growth rules *shape* it — but the underlying picture never changes.
{{< /note >}}

## The linear-algebra connection

A graph can be written as a **matrix**: row $i$, column $j$ is $1$ when an edge runs from $i$ to $j$. Once a network is a matrix, the whole machinery of linear algebra applies. [[PageRank]] is the headline example — a page's importance is its share of the **dominant eigenvector** of the link matrix, exactly the [[Eigenvalues and Eigenvectors]] idea from the [[Linear Algebra]] section. In fact, *this very site* ranks its related-topic suggestions with a PageRank-style computation over the wiki-link graph you are browsing right now.

## Explore

{{< columns count="2" >}}
**Foundations** — [[Graph]] · [[Breadth-First Search]] · [[Depth-First Search]]

**Weighted problems** — [[Dijkstra's Algorithm]] · [[Minimum Spanning Tree]]

**Importance** — [[PageRank]] · [[Centrality]]

**Network shapes** — [[Small-World Network]] · [[Scale-Free Network]]
{{< /columns >}}

## See also

- [[Complex Systems]]
- [[Linear Algebra]]
- [[Eigenvalues and Eigenvectors]]
