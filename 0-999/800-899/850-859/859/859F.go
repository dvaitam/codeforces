package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// Item represents a starting index in the dynamic programming.
type Item struct {
	t   int64 // threshold when interval saturates
	f   int64 // coefficient for unsaturated intervals
	idx int
}

// min-heap ordered by threshold t
type minHeap []Item

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].t < h[j].t }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x any)        { *h = append(*h, x.(Item)) }
func (h *minHeap) Pop() any          { n := len(*h); v := (*h)[n-1]; *h = (*h)[:n-1]; return v }

// max-heap ordered by coefficient f
type maxHeap []Item

func (h maxHeap) Len() int           { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i].f > h[j].f }
func (h maxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x any)        { *h = append(*h, x.(Item)) }
func (h *maxHeap) Pop() any          { n := len(*h); v := (*h)[n-1]; *h = (*h)[:n-1]; return v }

// minimalShirts computes the minimal total number of T-shirts needed.
func minimalShirts(n int, C int64, single []int64, multi []int64) int64 {
	// extend multi so that multiExt[i] corresponds to pair (i,i+1)
	multiExt := make([]int64, n+1)
	for i := 1; i < n; i++ {
		multiExt[i] = multi[i-1]
	}

	// prefix sums for singles and pairs
	S1 := make([]int64, n+1)
	S2 := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		S1[i] = S1[i-1] + single[i-1]
	}
	for i := 1; i <= n; i++ {
		S2[i] = S2[i-1] + multiExt[i]
	}
	// PS[i] = singles(1..i) + pairs(1..i-1)
	PS := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		PS[i] = S1[i] + S2[i-1]
	}

	dp := make([]int64, n+1)
	satMax := int64(-1 << 62)
	uheap := &maxHeap{}
	theap := &minHeap{}
	heap.Init(uheap)
	heap.Init(theap)
	active := make([]bool, n+1)

	for i := 1; i <= n; i++ {
		// add starting position i for future intervals
		F1 := dp[i-1] - PS[i-1] - multiExt[i-1]
		T := PS[i-1] + multiExt[i-1] + C
		it := Item{t: T, f: F1, idx: i}
		heap.Push(uheap, it)
		heap.Push(theap, it)
		active[i] = true

		// move candidates whose threshold is reached to the saturated set
		for theap.Len() > 0 && (*theap)[0].t <= PS[i] {
			v := heap.Pop(theap).(Item)
			if !active[v.idx] {
				continue
			}
			active[v.idx] = false
			val := dp[v.idx-1] + C
			if val > satMax {
				satMax = val
			}
		}
		for uheap.Len() > 0 && !active[(*uheap)[0].idx] {
			heap.Pop(uheap)
		}
		bestUnsat := int64(-1 << 62)
		if uheap.Len() > 0 {
			bestUnsat = (*uheap)[0].f + PS[i]
		}
		if satMax > bestUnsat {
			dp[i] = satMax
		} else {
			dp[i] = bestUnsat
		}
	}
	return dp[n]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var C int64
	if _, err := fmt.Fscan(in, &n, &C); err != nil {
		return
	}
	values := make([]int64, 2*n-1)
	for i := range values {
		fmt.Fscan(in, &values[i])
	}
	single := make([]int64, n)
	multi := make([]int64, n-1)
	for i := 0; i < n; i++ {
		single[i] = values[2*i]
		if i < n-1 {
			multi[i] = values[2*i+1]
		}
	}
	ans := minimalShirts(n, C, single, multi)
	fmt.Println(ans)
}
