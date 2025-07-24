package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

// Item heaps for costs used when buying water
// MinHeap provides access to the cheapest cost
// MaxHeap provides access to the most expensive cost

type MinHeap []int64

func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

type MaxHeap []int64

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// Event represents friend arrival
// and is processed in chronological order

type Event struct {
	t int64
	a int64
	b int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var n int
		var m, c, c0 int64
		fmt.Fscan(in, &n, &m, &c, &c0)
		events := make([]Event, n)
		var total int64 = c0
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &events[i].t, &events[i].a, &events[i].b)
			total += events[i].a
		}
		if total < m {
			// not enough water overall
			// read remaining input for this query if any
			fmt.Fprintln(out, -1)
			continue
		}
		sort.Slice(events, func(i, j int) bool { return events[i].t < events[j].t })
		events = append(events, Event{t: m})

		minH := &MinHeap{}
		maxH := &MaxHeap{}
		heap.Init(minH)
		heap.Init(maxH)
		counts := make(map[int64]int64)

		// initial water has cost 0
		counts[0] = c0
		heap.Push(minH, int64(0))
		heap.Push(maxH, int64(0))

		volume := c0
		var costSum int64
		prev := int64(0)
		possible := true

		consume := func(need int64) bool {
			for need > 0 {
				if minH.Len() == 0 {
					return false
				}
				cost := heap.Pop(minH).(int64)
				cnt := counts[cost]
				if cnt == 0 {
					continue
				}
				use := need
				if cnt < use {
					use = cnt
				}
				cnt -= use
				counts[cost] = cnt
				need -= use
				volume -= use
				costSum += cost * use
				if cnt > 0 {
					heap.Push(minH, cost)
					heap.Push(maxH, cost)
				}
			}
			return true
		}

		discard := func(rem int64) {
			for rem > 0 {
				if maxH.Len() == 0 {
					break
				}
				cost := heap.Pop(maxH).(int64)
				cnt := counts[cost]
				if cnt == 0 {
					continue
				}
				use := rem
				if cnt < use {
					use = cnt
				}
				cnt -= use
				counts[cost] = cnt
				rem -= use
				volume -= use
				if cnt > 0 {
					heap.Push(minH, cost)
					heap.Push(maxH, cost)
				}
			}
		}

		for _, e := range events {
			delta := e.t - prev
			if delta > 0 {
				if volume < delta {
					possible = false
					break
				}
				if !consume(delta) {
					possible = false
					break
				}
			}
			if e.t == m {
				break
			}
			counts[e.b] += e.a
			heap.Push(minH, e.b)
			heap.Push(maxH, e.b)
			volume += e.a
			if volume > c {
				discard(volume - c)
			}
			prev = e.t
		}

		if possible {
			fmt.Fprintln(out, costSum)
		} else {
			fmt.Fprintln(out, -1)
		}
	}
}
