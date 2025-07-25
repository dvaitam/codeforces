package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// Item is used in max-heap for generating top k sums.
type Item struct {
	val int64
	j   int
	idx int
}

type MaxHeap []Item

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i].val > h[j].val }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

func mergeTopK(a, b []int64, k int) []int64 {
	res := make([]int64, 0, k)
	i, j := 0, 0
	for len(res) < k && (i < len(a) || j < len(b)) {
		if j >= len(b) || (i < len(a) && a[i] >= b[j]) {
			res = append(res, a[i])
			i++
		} else {
			res = append(res, b[j])
			j++
		}
	}
	return res
}

func topKUnion(dps [][]int64, weights []int64, k int) []int64 {
	h := &MaxHeap{}
	heap.Init(h)
	for idx := range weights {
		arr := dps[idx]
		if len(arr) > 0 {
			heap.Push(h, Item{val: arr[0] + weights[idx], j: idx, idx: 0})
		}
	}
	res := make([]int64, 0, k)
	for len(res) < k && h.Len() > 0 {
		cur := heap.Pop(h).(Item)
		res = append(res, cur.val)
		arr := dps[cur.j]
		if cur.idx+1 < len(arr) {
			heap.Push(h, Item{val: arr[cur.idx+1] + weights[cur.j], j: cur.j, idx: cur.idx + 1})
		}
	}
	return res
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
		a := make([][]int64, n+1)
		for i := 1; i <= n; i++ {
			a[i] = make([]int64, n+1)
			for j := i; j <= n; j++ {
				fmt.Fscan(in, &a[i][j])
			}
		}

		dp := make([][]int64, n+1)
		dp[0] = []int64{0}
		for i := 1; i <= n; i++ {
			// collect weights for segments ending at i
			weights := make([]int64, i)
			for j := 1; j <= i; j++ {
				weights[j-1] = a[j][i]
			}
			segs := topKUnion(dp[:i], weights, k)
			dp[i] = mergeTopK(dp[i-1], segs, k)
		}

		ans := dp[n]
		if len(ans) < k {
			// pad with zeros if needed (should not happen per constraints)
			for len(ans) < k {
				ans = append(ans, 0)
			}
		}
		for i := 0; i < k; i++ {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(out, ans[i])
		}
		out.WriteByte('\n')
	}
}
