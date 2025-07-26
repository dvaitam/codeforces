package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type MaxHeap []int

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func possible(a []int, k int, m int64) bool {
	n := len(a)
	pre := make([]int, n+1)
	h := &MaxHeap{}
	heap.Init(h)
	var sum int64
	for i := 1; i <= n; i++ {
		x := a[i-1]
		if int64(x) <= m {
			heap.Push(h, x)
			sum += int64(x)
			if sum > m {
				largest := heap.Pop(h).(int)
				sum -= int64(largest)
			}
		}
		pre[i] = h.Len()
	}

	suf := make([]int, n+2)
	h = &MaxHeap{}
	heap.Init(h)
	sum = 0
	for i := n; i >= 1; i-- {
		x := a[i-1]
		if int64(x) <= m {
			heap.Push(h, x)
			sum += int64(x)
			if sum > m {
				largest := heap.Pop(h).(int)
				sum -= int64(largest)
			}
		}
		suf[i] = h.Len()
	}

	for t := 0; t <= n; t++ {
		if pre[t]+suf[t+1] >= k {
			return true
		}
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var tc int
	fmt.Fscan(reader, &tc)
	for ; tc > 0; tc-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		arr := make([]int, n)
		var hi int64
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
			hi += int64(arr[i])
		}
		lo := int64(0)
		for lo < hi {
			mid := (lo + hi) / 2
			if possible(arr, k, mid) {
				hi = mid
			} else {
				lo = mid + 1
			}
		}
		fmt.Fprintln(writer, lo)
	}
}
