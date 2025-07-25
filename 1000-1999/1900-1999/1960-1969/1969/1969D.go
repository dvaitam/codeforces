package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type item struct{ a, b int64 }

type maxHeap []int64

func (h maxHeap) Len() int            { return len(h) }
func (h maxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h maxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *maxHeap) Pop() interface{} {
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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		items := make([]item, n)
		for i := 0; i < n; i++ {
			items[i] = item{a[i], b[i]}
		}
		sort.Slice(items, func(i, j int) bool {
			if items[i].b == items[j].b {
				return items[i].a < items[j].a
			}
			return items[i].b > items[j].b
		})

		suff := make([]int64, n+1)
		for i := n - 1; i >= 0; i-- {
			suff[i] = suff[i+1]
			diff := items[i].b - items[i].a
			if diff > 0 {
				suff[i] += diff
			}
		}

		h := &maxHeap{}
		sum := int64(0)
		ans := int64(0)
		for i := 0; i < n; i++ {
			if len(*h) >= k {
				profit := suff[i] - sum
				if profit > ans {
					ans = profit
				}
			}
			heap.Push(h, items[i].a)
			sum += items[i].a
			if len(*h) > k {
				sum -= heap.Pop(h).(int64)
			}
		}
		fmt.Fprintln(out, ans)
	}
}
