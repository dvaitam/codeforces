package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Edge struct {
	to  int
	t0  int64
	mod int64
	ok  bool
}

type Item struct {
	t int64
	v int
}

type MinHeap []Item

func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool  { return h[i].t < h[j].t }
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func exgcd(a, b int64) (int64, int64, int64) {
	if b == 0 {
		return a, 1, 0
	}
	g, x1, y1 := exgcd(b, a%b)
	x := y1
	y := x1 - (a/b)*y1
	return g, x, y
}

func modInv(a, m int64) int64 {
	g, x, _ := exgcd(a, m)
	if g != 1 {
		return -1
	}
	x %= m
	if x < 0 {
		x += m
	}
	return x
}

func buildEdge(u, v int, l, s []int64, H int64) Edge {
	a := (s[v] - s[u]) % H
	if a < 0 {
		a += H
	}
	b := (l[u] - l[v]) % H
	if b < 0 {
		b += H
	}
	if a == 0 {
		if b == 0 {
			return Edge{to: v, t0: 0, mod: 1, ok: true}
		}
		return Edge{ok: false}
	}
	g := gcd(a, H)
	if b%g != 0 {
		return Edge{ok: false}
	}
	A := a / g
	B := b / g
	M := H / g
	invA := modInv(A%M, M)
	if invA == -1 {
		return Edge{ok: false}
	}
	t0 := (B * invA) % M
	return Edge{to: v, t0: t0, mod: M, ok: true}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		var H int64
		fmt.Fscan(in, &n, &m, &H)
		l := make([]int64, n)
		s := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &l[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &s[i])
		}
		adj := make([][]Edge, n)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			e1 := buildEdge(u, v, l, s, H)
			e2 := buildEdge(v, u, l, s, H)
			if e1.ok {
				adj[u] = append(adj[u], e1)
			}
			if e2.ok {
				adj[v] = append(adj[v], e2)
			}
		}
		const INF int64 = 1<<63 - 1
		dist := make([]int64, n)
		for i := range dist {
			dist[i] = INF
		}
		pq := &MinHeap{}
		heap.Init(pq)
		dist[0] = 0
		heap.Push(pq, Item{0, 0})
		for pq.Len() > 0 {
			it := heap.Pop(pq).(Item)
			if it.t != dist[it.v] {
				continue
			}
			if it.v == n-1 {
				break
			}
			tCur := it.t
			for _, e := range adj[it.v] {
				ready := e.t0
				if tCur > ready {
					k := (tCur - ready + e.mod - 1) / e.mod
					ready += k * e.mod
				}
				arrive := ready + 1
				if arrive < dist[e.to] {
					dist[e.to] = arrive
					heap.Push(pq, Item{arrive, e.to})
				}
			}
		}
		if dist[n-1] == INF {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, dist[n-1])
		}
	}
}
