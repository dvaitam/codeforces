package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type interval struct {
	l, r int
	id   int
}

type minHeap []interval

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i].r < h[j].r }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(interval)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	val := old[n-1]
	*h = old[:n-1]
	return val
}

// buildComponents returns disjoint segments where the coverage count is at least 2.
func buildComponents(ints []interval) [][2]int {
	events := make([][2]int, 0, 2*len(ints))
	for _, iv := range ints {
		events = append(events, [2]int{iv.l, 1})
		events = append(events, [2]int{iv.r + 1, -1})
	}
	sort.Slice(events, func(i, j int) bool {
		if events[i][0] == events[j][0] {
			return events[i][1] < events[j][1]
		}
		return events[i][0] < events[j][0]
	})
	if len(events) == 0 {
		return nil
	}
	comp := make([][2]int, 0)
	cnt := 0
	prev := events[0][0]
	for _, ev := range events {
		pos, delta := ev[0], ev[1]
		if cnt >= 2 && prev <= pos-1 {
			comp = append(comp, [2]int{prev, pos - 1})
		}
		cnt += delta
		prev = pos
	}
	return comp
}

// tryTwoColors attempts to find a subset of pairwise non-overlapping intervals
// that covers every component with coverage >= 2.
// If successful, it marks them with color 2 and returns true.
func tryTwoColors(ints []interval, comps [][2]int, colors []int) bool {
	if len(comps) == 0 {
		return true
	}
	// sort intervals by start
	sorted := make([]interval, len(ints))
	copy(sorted, ints)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].l == sorted[j].l {
			return sorted[i].r < sorted[j].r
		}
		return sorted[i].l < sorted[j].l
	})

	h := &minHeap{}
	heap.Init(h)
	idx := 0
	prevEnd := -1 << 60
	activeID := -1
	for _, c := range comps {
		L, R := c[0], c[1]
		if activeID != -1 {
			iv := ints[activeID]
			if iv.l <= L && R <= iv.r {
				continue // already covered by current chosen interval
			}
		}
		for idx < len(sorted) && sorted[idx].l <= L {
			if sorted[idx].l > prevEnd {
				heap.Push(h, sorted[idx])
			}
			idx++
		}
		for h.Len() > 0 && (*h)[0].r < R {
			heap.Pop(h)
		}
		if h.Len() == 0 {
			return false
		}
		iv := heap.Pop(h).(interval)
		colors[iv.id] = 2
		prevEnd = iv.r
		activeID = iv.id
		// All intervals currently in heap overlap the chosen one at point L, so discard them.
		*h = (*h)[:0]
	}
	return true
}

// buildThreeColors paints intervals with colors 2 and 3 so that the union of these
// intervals covers every point with coverage >= 2. Intervals of the same special
// color are disjoint by construction.
func buildThreeColors(ints []interval, comps [][2]int, colors []int) {
	if len(comps) == 0 {
		return
	}
	sorted := make([]interval, len(ints))
	copy(sorted, ints)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].l == sorted[j].l {
			return sorted[i].r < sorted[j].r
		}
		return sorted[i].l < sorted[j].l
	})
	h := &minHeap{}
	heap.Init(h)
	idx := 0
	for _, c := range comps {
		pos := c[0]
		end := c[1]
		curColor := 2
		for pos <= end {
			for idx < len(sorted) && sorted[idx].l <= pos {
				heap.Push(h, sorted[idx])
				idx++
			}
			for h.Len() > 0 && (*h)[0].r < pos {
				heap.Pop(h)
			}
			if h.Len() == 0 {
				// Should not happen because coverage is >=2 on this component.
				break
			}
			iv := heap.Pop(h).(interval)
			if colors[iv.id] == 1 {
				colors[iv.id] = curColor
			}
			pos = iv.r + 1
			if curColor == 2 {
				curColor = 3
			} else {
				curColor = 2
			}
		}
	}
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
		ints := make([]interval, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &ints[i].l, &ints[i].r)
			ints[i].id = i
		}

		colors := make([]int, n)
		for i := range colors {
			colors[i] = 1
		}

		comps := buildComponents(ints)

		if len(comps) == 0 {
			fmt.Fprintln(out, 1)
			for i, c := range colors {
				if i > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, c)
			}
			fmt.Fprintln(out)
			continue
		}

		if tryTwoColors(ints, comps, colors) {
			fmt.Fprintln(out, 2)
		} else {
			buildThreeColors(ints, comps, colors)
			fmt.Fprintln(out, 3)
		}
		for i, c := range colors {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, c)
		}
		fmt.Fprintln(out)
	}
}
