package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// Candidate represents a node with its fortress-neighbor ratio.
type Candidate struct {
	rate float64
	u    int
}

// MaxHeap implements a max-heap of Candidates based on rate.
type MaxHeap []Candidate

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i].rate > h[j].rate }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) {
	*h = append(*h, x.(Candidate))
}
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var N, M, K int
	fmt.Fscan(reader, &N, &M, &K)

	fortress := make([]bool, N+1)
	for i := 0; i < K; i++ {
		var u int
		fmt.Fscan(reader, &u)
		fortress[u] = true
	}

	gph := make([][]int, N+1)
	dg := make([]int, N+1)
	dgf := make([]int, N+1)

	for i := 0; i < M; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		gph[u] = append(gph[u], v)
		dg[u]++
		gph[v] = append(gph[v], u)
		dg[v]++
		if fortress[u] {
			dgf[v]++
		}
		if fortress[v] {
			dgf[u]++
		}
	}

	used := make([]bool, N+1)
	h := &MaxHeap{}
	heap.Init(h)

	for u := 1; u <= N; u++ {
		if fortress[u] {
			continue
		}
		heap.Push(h, Candidate{rate: float64(dgf[u]) / float64(dg[u]), u: u})
	}

	ord := make([]int, 0, N)
	ans := Candidate{rate: 1000.0, u: -1}

	for h.Len() > 0 {
		c := heap.Pop(h).(Candidate)
		if used[c.u] {
			continue
		}
		if c.rate < ans.rate {
			ans = c
		}
		u := c.u
		used[u] = true
		ord = append(ord, u)
		for _, v := range gph[u] {
			if fortress[v] || used[v] {
				continue
			}
			dgf[v]++
			heap.Push(h, Candidate{rate: float64(dgf[v]) / float64(dg[v]), u: v})
		}
	}

	// Trim ord to start from the weakest-point removal
	start := 0
	for i, u := range ord {
		if u == ans.u {
			start = i
			break
		}
	}
	ord = ord[start:]

	fmt.Fprintln(writer, len(ord))
	for i, u := range ord {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, u)
	}
	fmt.Fprintln(writer)
}
