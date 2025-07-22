package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

const MOD = 1009

// PolyItem is used in priority queue, sorted by polynomial length
type PolyItem struct {
	poly []int
}

type PolyHeap []*PolyItem

func (h PolyHeap) Len() int           { return len(h) }
func (h PolyHeap) Less(i, j int) bool { return len(h[i].poly) < len(h[j].poly) }
func (h PolyHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *PolyHeap) Push(x interface{}) {
	*h = append(*h, x.(*PolyItem))
}

func (h *PolyHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// multiply polynomials a and b modulo MOD, keeping only terms up to degree k
func polyMul(a, b []int, k int) []int {
	maxDeg := len(a) + len(b) - 2
	if maxDeg > k {
		maxDeg = k
	}
	res := make([]int, maxDeg+1)
	for i := 0; i < len(a); i++ {
		if a[i] == 0 {
			continue
		}
		limit := len(b)
		if i+limit-1 > k {
			limit = k - i + 1
			if limit <= 0 {
				break
			}
		}
		for j := 0; j < limit; j++ {
			res[i+j] += a[i] * b[j]
			res[i+j] %= MOD
		}
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}
	freq := make([]int, m)
	for i := 0; i < n; i++ {
		var c int
		fmt.Fscan(reader, &c)
		freq[c-1]++
	}

	pq := &PolyHeap{}
	heap.Init(pq)
	for i := 0; i < m; i++ {
		c := freq[i]
		if c > k {
			c = k
		}
		poly := make([]int, c+1)
		for j := 0; j <= c; j++ {
			poly[j] = 1
		}
		heap.Push(pq, &PolyItem{poly: poly})
	}

	for pq.Len() > 1 {
		a := heap.Pop(pq).(*PolyItem).poly
		b := heap.Pop(pq).(*PolyItem).poly
		c := polyMul(a, b, k)
		heap.Push(pq, &PolyItem{poly: c})
	}

	if pq.Len() == 0 {
		fmt.Fprintln(writer, 0)
		return
	}
	resPoly := heap.Pop(pq).(*PolyItem).poly
	if k < len(resPoly) {
		fmt.Fprintln(writer, resPoly[k]%MOD)
	} else {
		fmt.Fprintln(writer, 0)
	}
}
