package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type item struct {
	cnt int
	val int
}

type minHeap []item

func (h minHeap) Len() int { return len(h) }

func (h minHeap) Less(i, j int) bool {
	if h[i].cnt == h[j].cnt {
		return h[i].val < h[j].val
	}
	return h[i].cnt < h[j].cnt
}

func (h minHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *minHeap) Push(x interface{}) {
	*h = append(*h, x.(item))
}

func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	it := old[n-1]
	*h = old[:n-1]
	return it
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		counts := make([]int, n+2)
		ans := make([]int, n)
		h := make(minHeap, n)
		for i := 0; i < n; i++ {
			h[i] = item{cnt: 0, val: i + 1}
		}
		heap.Init(&h)
		maxCount := 0

		assign := func(idx int, val int) {
			counts[val]++
			if counts[val] > maxCount {
				maxCount = counts[val]
			}
			heap.Push(&h, item{cnt: counts[val], val: val})
			ans[idx] = val
		}

		findCandidate := func(forbidden int, limit int) int {
			temp := make([]item, 0)
			for h.Len() > 0 {
				it := heap.Pop(&h).(item)
				if counts[it.val] != it.cnt {
					continue
				}
				temp = append(temp, it)
				if it.val == forbidden {
					continue
				}
				if it.cnt <= limit-1 {
					for i := 0; i < len(temp)-1; i++ {
						heap.Push(&h, temp[i])
					}
					return it.val
				}
				break
			}
			for _, it := range temp {
				heap.Push(&h, it)
			}
			return -1
		}

		for i := 0; i < n; i++ {
			v := a[i]
			cntV := counts[v]
			if cntV < maxCount || cntV == 0 {
				assign(i, v)
				continue
			}

			candidate := findCandidate(v, cntV)
			if candidate == -1 {
				assign(i, v)
			} else {
				assign(i, candidate)
			}
		}

		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[i])
		}
		fmt.Fprintln(out)
	}
}
