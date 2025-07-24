package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// interval represents an active router covering up to r with cost value val.
type interval struct {
	r   int   // rightmost room this router covers
	val int64 // total cost if this router is used for rooms up to r
}

// priority queue ordered by val, then r.
type pq []interval

func (h pq) Len() int { return len(h) }
func (h pq) Less(i, j int) bool {
	if h[i].val == h[j].val {
		return h[i].r < h[j].r
	}
	return h[i].val < h[j].val
}
func (h pq) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *pq) Push(x interface{}) { *h = append(*h, x.(interval)) }
func (h *pq) Pop() interface{} {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	fmt.Fscan(reader, &n, &k)
	var s string
	fmt.Fscan(reader, &s)

	// group routers by the leftmost room they cover
	starts := make([][]int, n+2)
	for i := 1; i <= n; i++ {
		if s[i-1] == '1' {
			L := i - k
			if L < 1 {
				L = 1
			}
			starts[L] = append(starts[L], i)
		}
	}

	dp := make([]int64, n+1)
	active := &pq{}
	heap.Init(active)

	for i := 1; i <= n; i++ {
		// add all routers whose coverage starts at i
		for _, j := range starts[i] {
			r := j + k
			if r > n {
				r = n
			}
			val := dp[i-1] + int64(j)
			heap.Push(active, interval{r, val})
		}
		// remove routers that no longer cover position i
		for active.Len() > 0 && (*active)[0].r < i {
			heap.Pop(active)
		}
		best := int64(1 << 60)
		if active.Len() > 0 {
			best = (*active)[0].val
		}
		direct := dp[i-1] + int64(i)
		if direct < best {
			dp[i] = direct
		} else {
			dp[i] = best
		}
	}

	fmt.Fprintln(writer, dp[n])
}
