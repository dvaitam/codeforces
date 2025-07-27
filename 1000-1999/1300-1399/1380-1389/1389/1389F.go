package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type Event struct {
	x     int
	typ   int // 0=start, 1=end
	idx   int
	color int
	r     int
}

type Item struct {
	end int
	idx int
}

type MinHeap []Item

func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool  { return h[i].end < h[j].end }
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *MinHeap) Pop() interface{} {
	old := *h
	v := old[0]
	*h = old[1:]
	return v
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	events := make([]Event, 0, n*2)
	for i := 0; i < n; i++ {
		var l, r, t int
		fmt.Fscan(in, &l, &r, &t)
		color := t - 1
		events = append(events, Event{l, 0, i, color, r})
		events = append(events, Event{r, 1, i, color, r})
	}

	sort.Slice(events, func(i, j int) bool {
		if events[i].x != events[j].x {
			return events[i].x < events[j].x
		}
		return events[i].typ < events[j].typ
	})

	heaps := [2]*MinHeap{&MinHeap{}, &MinHeap{}}
	heap.Init(heaps[0])
	heap.Init(heaps[1])
	active := make([]bool, n)
	matches := 0

	for _, e := range events {
		if e.typ == 0 {
			active[e.idx] = true
			heap.Push(heaps[e.color], Item{end: e.r, idx: e.idx})
		} else {
			if !active[e.idx] {
				continue
			}
			opp := 1 - e.color
			for heaps[opp].Len() > 0 {
				top := (*heaps[opp])[0]
				if active[top.idx] {
					break
				}
				heap.Pop(heaps[opp])
			}
			if heaps[opp].Len() > 0 {
				top := heap.Pop(heaps[opp]).(Item)
				active[top.idx] = false
				active[e.idx] = false
				matches++
			} else {
				active[e.idx] = false
			}
		}
	}

	fmt.Fprintln(out, n-matches)
}
