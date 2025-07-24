package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Item struct {
	sum int
	i   int
	j   int
}

type MaxHeap []Item

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i].sum > h[j].sum }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MaxHeap) Push(x interface{}) {
	*h = append(*h, x.(Item))
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
	var n int
	var k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}
	pref := make([]int, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1] + a[i]
	}

	maxPrefIdx := make([]int, n+2)
	maxPrefIdx[n+1] = n
	for i := n; i >= 1; i-- {
		if pref[i] >= pref[maxPrefIdx[i+1]] {
			maxPrefIdx[i] = i
		} else {
			maxPrefIdx[i] = maxPrefIdx[i+1]
		}
	}

	h := &MaxHeap{}
	heap.Init(h)
	for i := 1; i <= n; i++ {
		j := maxPrefIdx[i]
		if j >= i {
			heap.Push(h, Item{sum: pref[j] - pref[i-1], i: i, j: j})
		}
	}
	ans := int64(0)
	for t := 0; t < k && h.Len() > 0; t++ {
		it := heap.Pop(h).(Item)
		ans += int64(it.sum)
		if it.j > it.i {
			j2 := it.j - 1
			heap.Push(h, Item{sum: pref[j2] - pref[it.i-1], i: it.i, j: j2})
		}
	}
	fmt.Println(ans)
}
