package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type IntMinHeap []int

func (h IntMinHeap) Len() int            { return len(h) }
func (h IntMinHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h IntMinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *IntMinHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *IntMinHeap) Pop() interface{} {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

type IntMaxHeap []int

func (h IntMaxHeap) Len() int            { return len(h) }
func (h IntMaxHeap) Less(i, j int) bool  { return h[i] > h[j] } // max-heap
func (h IntMaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *IntMaxHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *IntMaxHeap) Pop() interface{} {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

type pair struct {
	cap, dist int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		posByVal := make([][]int, n+1)
		for i := 1; i <= n; i++ {
			var v int
			fmt.Fscan(in, &v)
			if v > n {
				v = n
			}
			posByVal[v] = append(posByVal[v], i)
		}

		minH := &IntMinHeap{}
		maxH := &IntMaxHeap{}
		heap.Init(minH)
		heap.Init(maxH)
		alive := make([]bool, n+1) // 1-based positions
		cnt := 0

		popMin := func() int {
			for minH.Len() > 0 {
				x := heap.Pop(minH).(int)
				if alive[x] {
					return x
				}
			}
			return -1
		}
		popMax := func() int {
			for maxH.Len() > 0 {
				x := heap.Pop(maxH).(int)
				if alive[x] {
					return x
				}
			}
			return -1
		}

		var pairs []pair
		for v := n; v >= 1; v-- {
			for _, p := range posByVal[v] {
				heap.Push(minH, p)
				heap.Push(maxH, p)
				alive[p] = true
				cnt++
			}
			for cnt >= 2 {
				l := popMin()
				r := popMax()
				if l == -1 || r == -1 {
					break
				}
				if l == r {
					// put back if only one distinct element
					heap.Push(minH, l)
					heap.Push(maxH, r)
					break
				}
				alive[l] = false
				alive[r] = false
				cnt -= 2
				pairs = append(pairs, pair{cap: v, dist: r - l})
			}
		}

		sort.Slice(pairs, func(i, j int) bool {
			if pairs[i].cap == pairs[j].cap {
				return pairs[i].dist > pairs[j].dist
			}
			return pairs[i].cap < pairs[j].cap
		})

		// Select feasible pairs with maximum total distance.
		distH := &IntMinHeap{}
		heap.Init(distH)
		var total int64
		for _, p := range pairs {
			heap.Push(distH, p.dist)
			total += int64(p.dist)
			if distH.Len() > p.cap {
				removed := heap.Pop(distH).(int)
				total -= int64(removed)
			}
		}

		fmt.Fprintln(out, total)
	}
}
