package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF = int(1e9)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	l := make([]int, n+1)
	bases := make([]int, 0)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &l[i])
		if l[i] == 1 {
			bases = append(bases, i)
		}
	}
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	// compute height to nearest base
	h := make([]int, n+1)
	for i := 1; i <= n; i++ {
		h[i] = -1
	}
	q := make([]int, len(bases))
	copy(q, bases)
	for _, b := range bases {
		h[b] = 0
	}
	for head := 0; head < len(q); head++ {
		v := q[head]
		for _, u := range adj[v] {
			if h[u] == -1 {
				h[u] = h[v] + 1
				q = append(q, u)
			}
		}
	}

	// minimal height of a flat edge along some path to a base
	m := make([]int, n+1)
	for i := 1; i <= n; i++ {
		m[i] = INF
	}
	q2 := make([]int, 0)
	for u := 1; u <= n; u++ {
		for _, v := range adj[u] {
			if h[u] == h[v] {
				if m[u] > h[u] {
					m[u] = h[u]
					q2 = append(q2, u)
				}
				if m[v] > h[v] {
					m[v] = h[v]
					q2 = append(q2, v)
				}
			}
		}
	}
	for head := 0; head < len(q2); head++ {
		v := q2[head]
		for _, u := range adj[v] {
			if h[u] >= h[v] && m[v] < m[u] {
				m[u] = m[v]
				q2 = append(q2, u)
			}
		}
	}

	ans := make([]int, n+1)
	for u := 1; u <= n; u++ {
		ans[u] = h[u]
		if h[u] == 0 {
			continue
		}
		for _, v := range adj[u] {
			if h[v] == h[u]-1 {
				cand := h[u]
				if m[v] != INF {
					cand = 2*h[u] - m[v]
				}
				if cand > ans[u] {
					ans[u] = cand
				}
			}
		}
	}

	for i := 1; i <= n; i++ {
		if i > 1 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, ans[i])
	}
	out.WriteByte('\n')
}
