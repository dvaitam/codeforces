package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type pair struct {
	val int
	l   int
}

type event struct {
	end int
	val int
}

// max-heap for active segments
type item struct {
	val int
	end int
}
type maxHeap []item

func (h maxHeap) Len() int           { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i].val > h[j].val }
func (h maxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) {
	*h = append(*h, x.(item))
}
func (h *maxHeap) Pop() interface{} {
	old := *h
	v := old[len(old)-1]
	*h = old[:len(old)-1]
	return v
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		mask := (1 << k) - 1

		events := make([][]event, n)
		curr := make([]pair, 0)

		for r := 0; r < n; r++ {
			tmp := make([]pair, 0, len(curr)+1)
			// subarray of length 1
			tmp = append(tmp, pair{val: a[r], l: r})

			for _, p := range curr {
				nv := ^(p.val | a[r]) & mask
				if len(tmp) > 0 && tmp[len(tmp)-1].val == nv {
					if p.l < tmp[len(tmp)-1].l {
						tmp[len(tmp)-1].l = p.l
					}
				} else {
					tmp = append(tmp, pair{val: nv, l: p.l})
				}
			}

			curr = tmp
			for _, p := range curr {
				events[p.l] = append(events[p.l], event{end: r, val: p.val})
			}
		}

	ans := make([]int, n)
		h := &maxHeap{}
		heap.Init(h)
		for i := 0; i < n; i++ {
			for _, ev := range events[i] {
				heap.Push(h, item{val: ev.val, end: ev.end})
			}
			for h.Len() > 0 && (*h)[0].end < i {
				heap.Pop(h)
			}
			if h.Len() > 0 {
				ans[i] = (*h)[0].val
			} else {
				ans[i] = 0
			}
		}

		for i, v := range ans {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
