package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	sum int
	i   int
	j   int
}

type pairHeap []pair

func (h pairHeap) Len() int            { return len(h) }
func (h pairHeap) Less(a, b int) bool  { return h[a].sum < h[b].sum }
func (h pairHeap) Swap(a, b int)       { h[a], h[b] = h[b], h[a] }
func (h *pairHeap) Push(x interface{}) { *h = append(*h, x.(pair)) }
func (h *pairHeap) Pop() interface{} {
	old := *h
	v := old[len(old)-1]
	*h = old[:len(old)-1]
	return v
}

func kSmallestPairSums(arr []int, k int) []int {
	n := len(arr)
	if k <= 0 || n == 0 {
		return nil
	}
	h := &pairHeap{}
	heap.Init(h)
	for i := 0; i < n; i++ {
		heap.Push(h, pair{arr[i] + arr[0], i, 0})
	}
	res := make([]int, 0, k)
	for h.Len() > 0 && len(res) < k {
		p := heap.Pop(h).(pair)
		res = append(res, p.sum)
		if p.j+1 < n {
			heap.Push(h, pair{arr[p.i] + arr[p.j+1], p.i, p.j + 1})
		}
	}
	return res
}

const mod = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	for len(arr) > 1 {
		sort.Ints(arr)
		m := len(arr)
		next := kSmallestPairSums(arr, m-1)
		arr = next
	}
	ans := arr[0] % mod
	fmt.Println(ans)
}
