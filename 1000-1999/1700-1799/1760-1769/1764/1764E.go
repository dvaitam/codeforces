package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type Pair struct {
	a int
	b int
}

type MaxHeap []int

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i], &b[i])
		}

		if k <= a[0] {
			fmt.Fprintln(out, "YES")
			continue
		}
		if k > a[0]+b[0] {
			fmt.Fprintln(out, "NO")
			continue
		}
		target := k - b[0]
		pairs := make([]Pair, 0, n-1)
		for i := 1; i < n; i++ {
			pairs = append(pairs, Pair{a[i], b[i]})
		}
		sort.Slice(pairs, func(i, j int) bool { return pairs[i].a < pairs[j].a })
		cur := 0
		for _, p := range pairs {
			if p.a <= a[0] && p.a > cur {
				cur = p.a
			}
		}
		if cur == 0 {
			cur = a[0]
		}
		idx := 0
		h := &MaxHeap{}
		heap.Init(h)
		for idx < len(pairs) && pairs[idx].a <= cur {
			heap.Push(h, pairs[idx].b)
			idx++
		}
		for h.Len() > 0 && cur < target {
			cur += heap.Pop(h).(int)
			if cur > a[0] {
				cur = a[0]
			}
			for idx < len(pairs) && pairs[idx].a <= cur {
				heap.Push(h, pairs[idx].b)
				idx++
			}
		}
		if cur >= target {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
