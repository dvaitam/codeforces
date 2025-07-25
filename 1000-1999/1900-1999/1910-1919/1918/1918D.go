package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// Pair represents a value with its index for the priority queue.
type Pair struct {
	val int64
	idx int
}

// MinHeap implements a min-heap for Pair based on val.
type MinHeap []Pair

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].val < h[j].val }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x any)        { *h = append(*h, x.(Pair)) }
func (h *MinHeap) Pop() any {
	old := *h
	v := old[len(old)-1]
	*h = old[:len(old)-1]
	return v
}

// minBlockedSum computes the minimal sum of blocked elements required so that
// every unblocked contiguous segment has sum at most M. A zero-valued sentinel
// is appended to simplify handling of the last segment.
func minBlockedSum(a []int64, M int64) int64 {
	n := len(a)
	prefix := make([]int64, n+2)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1] + a[i-1]
	}
	prefix[n+1] = prefix[n]

	arr := append(a, 0) // sentinel element

	dp := make([]int64, n+2)
	h := &MinHeap{}
	heap.Init(h)
	heap.Push(h, Pair{0, 0}) // dp[0]
	L := 0
	for i := 1; i <= n+1; i++ {
		thresh := prefix[i-1] - M
		for L <= i-1 && prefix[L] < thresh {
			L++
		}
		for h.Len() > 0 && (*h)[0].idx < L {
			heap.Pop(h)
		}
		minVal := int64(1 << 62)
		if h.Len() > 0 {
			minVal = (*h)[0].val
		}
		dp[i] = arr[i-1] + minVal
		heap.Push(h, Pair{dp[i], i})
	}
	return dp[n+1]
}

func solveOne(reader *bufio.Reader, writer *bufio.Writer) {
	var n int
	fmt.Fscan(reader, &n)
	a := make([]int64, n)
	var mx, sum int64
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		if a[i] > mx {
			mx = a[i]
		}
		sum += a[i]
	}
	l, r := mx, sum
	for l < r {
		mid := (l + r) / 2
		if minBlockedSum(a, mid) <= mid {
			r = mid
		} else {
			l = mid + 1
		}
	}
	fmt.Fprintln(writer, l)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		solveOne(reader, writer)
	}
}
