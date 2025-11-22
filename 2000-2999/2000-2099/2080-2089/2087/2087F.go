package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type pair struct {
	t   int // time of the last kill (stage index)
	sum int // cumulative sum of kill times
}

// priority queue ordered by (sum, time)
type node struct {
	sum int
	t   int
	p   int
	e   int
	idx int // index in state list to validate
}

type pq []node

func (h pq) Len() int { return len(h) }
func (h pq) Less(i, j int) bool {
	if h[i].sum == h[j].sum {
		return h[i].t < h[j].t
	}
	return h[i].sum < h[j].sum
}
func (h pq) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *pq) Push(x interface{}) { *h = append(*h, x.(node)) }
func (h *pq) Pop() interface{} {
	old := *h
	v := old[len(old)-1]
	*h = old[:len(old)-1]
	return v
}

const inf = int(1e9)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &b[i])
	}

	// posCache[p][e] = earliest stage when at least k+1 monsters are killable, where k=p+e
	posCache := make([][]int, n+1)
	for i := range posCache {
		posCache[i] = make([]int, n-i+1)
		for j := range posCache[i] {
			posCache[i][j] = -1
		}
	}

	// helper to get earliest position for state (p,e)
	getPos := func(p, e int) int {
		k := p + e
		if posCache[p][e] != -1 {
			return posCache[p][e]
		}
		need := k + 1
		P := p + 1
		E := e + 1
		cnt := 0
		pos := inf
		for i := 0; i < n; i++ {
			if a[i] <= P || b[i] <= E {
				cnt++
				if cnt == need {
					pos = i + 1 // stages are 1-indexed
					break
				}
			}
		}
		if cnt < need {
			pos = inf
		}
		posCache[p][e] = pos
		return pos
	}

	// For each state (p,e) store Pareto frontier of (time,sum)
	states := make([][][]pair, n+1)
	for p := 0; p <= n; p++ {
		states[p] = make([][]pair, n-p+1)
	}

	pushState := func(p, e, t, sum int, h *pq) {
		lst := states[p][e]
		// If an existing pair dominates the new one, skip it.
		for _, pr := range lst {
			if pr.t <= t && pr.sum <= sum {
				return
			}
		}
		states[p][e] = append(lst, pair{t: t, sum: sum})
		idx := len(states[p][e]) - 1
		heap.Push(h, node{sum: sum, t: t, p: p, e: e, idx: idx})
	}

	h := &pq{}
	heap.Init(h)
	pushState(0, 0, 0, 0, h)

	bestSum := inf
	target := n * (n + 1) / 2 // sum of spawn times, used later

	for h.Len() > 0 {
		cur := heap.Pop(h).(node)
		lst := states[cur.p][cur.e]
		if cur.idx >= len(lst) {
			continue
		}
		pr := lst[cur.idx]
		if pr.t != cur.t || pr.sum != cur.sum {
			continue // outdated
		}
		k := cur.p + cur.e
		if k == n {
			if pr.sum < bestSum {
				bestSum = pr.sum
			}
			continue
		}
		pos := getPos(cur.p, cur.e)
		if pos == inf {
			continue
		}
		nt := pos
		if cur.t+1 > nt {
			nt = cur.t + 1
		}
		if nt > 1337 { // cannot act after game ends
			continue
		}
		nsum := pr.sum + nt
		if cur.p+1 <= n {
			pushState(cur.p+1, cur.e, nt, nsum, h)
		}
		if cur.e+1 <= n {
			pushState(cur.p, cur.e+1, nt, nsum, h)
		}
	}

	if bestSum == inf {
		fmt.Println(-1)
		return
	}
	damage := bestSum - target
	fmt.Println(damage)
}
