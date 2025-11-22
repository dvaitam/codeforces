package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type bit struct {
	n    int
	tree []int
}

func newBit(n int) *bit {
	return &bit{n: n, tree: make([]int, n+2)}
}

func (b *bit) add(idx, delta int) {
	for idx <= b.n {
		b.tree[idx] += delta
		idx += idx & -idx
	}
}

func (b *bit) kth(k int) int {
	// returns smallest idx such that prefix sum >= k, assumes k in [1, total]
	pos := 0
	bitMask := 1
	for bitMask<<1 <= b.n {
		bitMask <<= 1
	}
	for bitMask > 0 {
		next := pos + bitMask
		if next <= b.n && b.tree[next] < k {
			k -= b.tree[next]
			pos = next
		}
		bitMask >>= 1
	}
	return pos + 1
}

type interval struct {
	low, high int
	l, r      int
}

// max-heap by high
type intervalHeap []interval

func (h intervalHeap) Len() int            { return len(h) }
func (h intervalHeap) Less(i, j int) bool  { return h[i].high > h[j].high }
func (h intervalHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *intervalHeap) Push(x interface{}) { *h = append(*h, x.(interval)) }
func (h *intervalHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func collectIntervals(a []int, k int) []interval {
	n := len(a)
	lengths := []int{k}
	if k < n {
		lengths = append(lengths, k+1)
	}

	res := make([]interval, 0, 2*n)
	for _, L := range lengths {
		k1 := (L + 1) / 2
		k2 := L/2 + 1
		b := newBit(n)
		for i := 0; i < L; i++ {
			b.add(a[i], 1)
		}
		low := b.kth(k1)
		high := b.kth(k2)
		res = append(res, interval{low, high, 1, L})
		for i := L; i < n; i++ {
			b.add(a[i-L], -1)
			b.add(a[i], 1)
			low = b.kth(k1)
			high = b.kth(k2)
			res = append(res, interval{low, high, i - L + 2, i + 1})
		}
	}
	return res
}

func solveCase(n, k int, a []int) (ansV, ansL, ansR []int) {
	ints := collectIntervals(a, k)
	sort.Slice(ints, func(i, j int) bool {
		if ints[i].low == ints[j].low {
			return ints[i].high < ints[j].high
		}
		return ints[i].low < ints[j].low
	})

	h := intervalHeap{}
	heap.Init(&h)
	idx := 0

	for v := 1; v <= n; v++ {
		for idx < len(ints) && ints[idx].low <= v {
			heap.Push(&h, ints[idx])
			idx++
		}
		for h.Len() > 0 && h[0].high < v {
			heap.Pop(&h)
		}
		if h.Len() > 0 {
			cur := h[0]
			ansV = append(ansV, v)
			ansL = append(ansL, cur.l)
			ansR = append(ansR, cur.r)
		}
	}
	return
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		v, l, r := solveCase(n, k, a)
		fmt.Fprintln(out, len(v))
		for i := range v {
			fmt.Fprintf(out, "%d %d %d\n", v[i], l[i], r[i])
		}
	}
}
