package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Union-find
type UF struct{ p []int }

func NewUF(n int) *UF {
	p := make([]int, n)
	for i := range p {
		p[i] = i
	}
	return &UF{p}
}
func (u *UF) Find(x int) int {
	if u.p[x] != x {
		u.p[x] = u.Find(u.p[x])
	}
	return u.p[x]
}
func (u *UF) Union(a, b int) {
	ra, rb := u.Find(a), u.Find(b)
	if ra != rb {
		u.p[ra] = rb
	}
}

// Item for max-heap
type Item struct {
	node, val int
}

// MaxPQ implements heap.Interface (max-heap by val)
type MaxPQ []Item

func (h MaxPQ) Len() int            { return len(h) }
func (h MaxPQ) Less(i, j int) bool  { return h[i].val > h[j].val }
func (h MaxPQ) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxPQ) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *MaxPQ) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func main() {
	rand.Seed(time.Now().UnixNano())
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var T int
	fmtFscan := func(a ...interface{}) { _, _ = fmt.Fscan(in, a...) }
	fmtFprint := func(a ...interface{}) { fmt.Fprint(out, a...) }
	fmtFprintln := func(a ...interface{}) { fmt.Fprintln(out, a...) }

	fmtFscan(&T)
	for T > 0 {
		T--
		var N1, N2, M int
		fmtFscan(&N1, &N2, &M)
		N := N1 + N2
		uf := NewUF(N + 1)
		edges := make([][2]int, M)
		g := make([][]int, N+1)
		for i := 0; i < M; i++ {
			u, v := 0, 0
			fmtFscan(&u, &v)
			edges[i][0], edges[i][1] = u, v
			if uf.Find(u) != uf.Find(v) {
				g[u] = append(g[u], v)
				g[v] = append(g[v], u)
				uf.Union(u, v)
			}
		}
		flag := false
		if N1 > N2 {
			flag = true
			N1, N2 = N2, N1
		}
		// try random partitions
		for {
			// random values
			val := make([]int, N+1)
			for i := 1; i <= N; i++ {
				val[i] = rand.Intn(3*N) + 1
			}
			visited := make([]bool, N+1)
			need := N1
			// priority queue
			pq := &MaxPQ{}
			heap.Init(pq)
			seed := rand.Intn(N) + 1
			heap.Push(pq, Item{seed, val[seed]})
			for pq.Len() > 0 {
				it := heap.Pop(pq).(Item)
				u := it.node
				if visited[u] {
					continue
				}
				visited[u] = true
				need--
				if need == 0 {
					break
				}
				for _, v := range g[u] {
					if !visited[v] {
						heap.Push(pq, Item{v, val[v]})
					}
				}
			}
			// check complement connectivity
			uf2 := NewUF(N + 1)
			for i := 0; i < M; i++ {
				u, v := edges[i][0], edges[i][1]
				if !visited[u] && !visited[v] {
					uf2.Union(u, v)
				}
			}
			iv := -1
			ok := true
			for i := 1; i <= N; i++ {
				if !visited[i] {
					fi := uf2.Find(i)
					if iv == -1 {
						iv = fi
					} else if iv != fi {
						ok = false
						break
					}
				}
			}
			if ok {
				// output two groups
				var tail []int
				for i := 1; i <= N; i++ {
					if visited[i] == flag {
						tail = append(tail, i)
					} else {
						fmtFprint(i, " ")
					}
				}
				fmtFprintln()
				for _, x := range tail {
					fmtFprint(x, " ")
				}
				fmtFprintln()
				break
			}
		}
	}
}
