package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

// Event represents a building edge: start or end
type Event struct {
	x   int64
	h   int64
	typ int // 1=start, 0=end
}

// MaxHeap is a max-heap of int64
type MaxHeap []int64

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *MaxHeap) Pop() interface{} {
	a := *h
	v := a[len(a)-1]
	*h = a[:len(a)-1]
	return v
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	events := make([]Event, 0, 2*n)
	for i := 0; i < n; i++ {
		var h, l, r int64
		fmt.Fscan(reader, &h, &l, &r)
		events = append(events, Event{x: l, h: h, typ: 1})
		events = append(events, Event{x: r, h: h, typ: 0})
	}
	sort.Slice(events, func(i, j int) bool {
		if events[i].x != events[j].x {
			return events[i].x < events[j].x
		}
		if events[i].typ != events[j].typ {
			return events[i].typ > events[j].typ // start before end
		}
		if events[i].typ == 1 {
			return events[i].h > events[j].h // higher start first
		}
		return events[i].h < events[j].h // lower end first
	})
	// max-heap and removal map for lazy deletion
	hq := &MaxHeap{}
	heap.Init(hq)
	heap.Push(hq, 0)
	rem := make(map[int64]int)
	prevMax := int64(0)
	// key points
	var keys []Event
	for i := 0; i < len(events); {
		x := events[i].x
		// process all events at x
		j := i
		for j < len(events) && events[j].x == x {
			e := events[j]
			if e.typ == 1 {
				heap.Push(hq, e.h)
			} else {
				rem[e.h]++
			}
			j++
		}
		// clean up heap
		for hq.Len() > 0 {
			top := (*hq)[0]
			if cnt, ok := rem[top]; ok && cnt > 0 {
				heap.Pop(hq)
				rem[top]--
			} else {
				break
			}
		}
		curr := (*hq)[0]
		if curr != prevMax {
			keys = append(keys, Event{x: x, h: curr})
			prevMax = curr
		}
		i = j
	}
	// build polyline vertices
	var verts [][2]int64
	if len(keys) > 0 {
		// start at ground at first x
		prevX := keys[0].x
		prevH := int64(0)
		verts = append(verts, [2]int64{prevX, 0})
		for _, kp := range keys {
			x, h := kp.x, kp.h
			if h != prevH {
				if x != prevX {
					verts = append(verts, [2]int64{x, prevH})
				}
				verts = append(verts, [2]int64{x, h})
				prevX, prevH = x, h
			}
		}
	}
	// output
	fmt.Fprintln(writer, len(verts))
	for _, v := range verts {
		fmt.Fprintln(writer, v[0], v[1])
	}
}
