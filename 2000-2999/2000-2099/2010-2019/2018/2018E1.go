package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
	"sort"
)

type segment struct {
	l, r int
}

type minHeap []int

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}
func (h minHeap) Peek() int { return h[0] }

func countGroups(g int, segs []segment) int {
	if g <= 0 {
		return 0
	}
	n := len(segs)
	idx := 0
	pos := 0
	groups := 0
	h := &minHeap{}
	heap.Init(h)
	for idx < n || h.Len() > 0 {
		if h.Len() == 0 {
			if idx >= n {
				break
			}
			pos = segs[idx].l
		}
		for idx < n && segs[idx].l <= pos {
			heap.Push(h, segs[idx].r)
			idx++
		}
		for h.Len() > 0 && h.Peek() < pos {
			heap.Pop(h)
		}
		if h.Len() < g {
			if h.Len() == 0 {
				continue
			}
			nextStart := math.MaxInt32
			if idx < n {
				nextStart = segs[idx].l
			}
			earliestEnd := h.Peek()
			if nextStart <= earliestEnd {
				pos = nextStart
			} else {
				pos = earliestEnd + 1
			}
			continue
		}
		Rmax := 0
		for i := 0; i < g; i++ {
			r := heap.Pop(h).(int)
			if r > Rmax {
				Rmax = r
			}
		}
		groups++
		*h = (*h)[:0]
		for idx < n && segs[idx].l <= Rmax {
			idx++
		}
		pos = Rmax + 1
	}
	return groups
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		l := make([]int, n)
		r := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &l[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &r[i])
		}
		segs := make([]segment, n)
		for i := 0; i < n; i++ {
			segs[i] = segment{l[i], r[i]}
		}
		sort.Slice(segs, func(i, j int) bool {
			if segs[i].l == segs[j].l {
				return segs[i].r < segs[j].r
			}
			return segs[i].l < segs[j].l
		})
		cache := make(map[int]int)
		var compute func(int) int
		compute = func(g int) int {
			if val, ok := cache[g]; ok {
				return val
			}
			val := countGroups(g, segs)
			cache[g] = val
			return val
		}
		best := 0
		maxG := n
		B := int(math.Sqrt(float64(n))) + 2
		if B > maxG {
			B = maxG
		}
		for g := 1; g <= B; g++ {
			groups := compute(g)
			if groups > 0 {
				total := g * groups
				if total > best {
					best = total
				}
			}
		}
		if B < maxG {
			for target := 1; target <= B; target++ {
				low := B + 1
				if compute(low) < target {
					continue
				}
				high := maxG
				res := low
				for low <= high {
					mid := (low + high) / 2
					if compute(mid) >= target {
						res = mid
						low = mid + 1
					} else {
						high = mid - 1
					}
				}
				total := res * target
				if total > best {
					best = total
				}
			}
		}
		fmt.Fprintln(out, best)
	}
}
