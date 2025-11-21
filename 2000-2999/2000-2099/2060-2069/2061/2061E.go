package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math/bits"
	"os"
)

type item struct {
	gain int64
	idx  int
}

type maxHeap []item

func (h maxHeap) Len() int            { return len(h) }
func (h maxHeap) Less(i, j int) bool  { return h[i].gain > h[j].gain }
func (h maxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(item)) }
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
	fullMask := (1 << 30) - 1

	for ; t > 0; t-- {
		var n, m int
		var k int
		fmt.Fscan(in, &n, &m, &k)

		a := make([]int, n)
		var total int64
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			total += int64(a[i])
		}

		b := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &b[i])
		}

		limit := 1 << m
		maskVal := make([]int, limit)
		pop := make([]int, limit)
		maskVal[0] = fullMask
		for mask := 1; mask < limit; mask++ {
			lowbit := mask & -mask
			bit := bits.TrailingZeros(uint(lowbit))
			maskVal[mask] = maskVal[mask^lowbit] & b[bit]
			pop[mask] = pop[mask^lowbit] + 1
		}

		best := make([][]int64, n)
		cur := make([]int, n)
		h := maxHeap{}
		for i := 0; i < n; i++ {
			best[i] = make([]int64, m+1)
			for c := 0; c <= m; c++ {
				best[i][c] = int64(a[i])
			}
			for mask := 1; mask < limit; mask++ {
				c := pop[mask]
				val := int64(a[i] & maskVal[mask])
				if val < best[i][c] {
					best[i][c] = val
				}
			}
			for c := 1; c <= m; c++ {
				if best[i][c] > best[i][c-1] {
					best[i][c] = best[i][c-1]
				}
			}
			if m >= 1 && best[i][1] < best[i][0] {
				heap.Push(&h, item{best[i][0] - best[i][1], i})
			}
		}

		ops := 0
		for ops < k && h.Len() > 0 {
			it := heap.Pop(&h).(item)
			if it.gain <= 0 {
				break
			}
			total -= it.gain
			idx := it.idx
			cur[idx]++
			ops++
			if cur[idx] < m {
				nextGain := best[idx][cur[idx]] - best[idx][cur[idx]+1]
				if nextGain > 0 {
					heap.Push(&h, item{nextGain, idx})
				}
			}
		}

		fmt.Fprintln(out, total)
	}
}
