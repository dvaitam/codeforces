package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
)

type maxHeap []float64

func (h maxHeap) Len() int            { return len(h) }
func (h maxHeap) Less(i, j int) bool  { return h[i] > h[j] } // reverse for max-heap
func (h maxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(float64)) }
func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var xEnc, yEnc int
	var m int
	fmt.Fscan(in, &xEnc, &yEnc, &m)
	p := float64(xEnc) / 1000.0
	q := float64(yEnc) / 1000.0

	k := make([]float64, n)
	b := make([]float64, n)
	for i := 0; i < n; i++ {
		var ai, bi int
		fmt.Fscan(in, &ai, &bi)
		k[i] = float64(ai) / 1000.0
		b[i] = float64(bi) / 1000.0
	}

	h := &maxHeap{}
	heap.Init(h)

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if k[i] == k[j] {
				// parallel lines do not intersect
				continue
			}
			x := (b[j] - b[i]) / (k[i] - k[j])
			y := k[i]*x + b[i]
			d := math.Hypot(x-p, y-q)
			if h.Len() < m {
				heap.Push(h, d)
			} else if d < (*h)[0] {
				heap.Pop(h)
				heap.Push(h, d)
			}
		}
	}

	sum := 0.0
	for h.Len() > 0 {
		sum += heap.Pop(h).(float64)
	}

	fmt.Printf("%.10f\n", sum)
}
