package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math/bits"
	"os"
)

type node struct {
	cost uint64
	idx  int
}

type minHeap []node

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i].cost < h[j].cost }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(node)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

func nextCost(v uint64) uint64 {
	// smallest power of two corresponding to the lowest zero bit of v
	return 1 << bits.TrailingZeros64(^v)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		var k uint64
		fmt.Fscan(in, &n, &k)
		vals := make([]uint64, n)
		var beauty uint64
		h := &minHeap{}
		heap.Init(h)

		for i := 0; i < n; i++ {
			fmt.Fscan(in, &vals[i])
			beauty += uint64(bits.OnesCount64(vals[i]))
			heap.Push(h, node{cost: nextCost(vals[i]), idx: i})
		}

		for h.Len() > 0 {
			top := heap.Pop(h).(node)
			if top.cost > k {
				break
			}
			k -= top.cost
			beauty++

			vals[top.idx] += top.cost
			heap.Push(h, node{cost: nextCost(vals[top.idx]), idx: top.idx})
		}

		fmt.Fprintln(out, beauty)
	}
}
