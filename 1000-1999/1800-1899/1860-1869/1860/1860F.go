package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type Item struct {
	a    int64
	b    int64
	sign int
}

type Event struct {
	r     float64
	id    int
	left  int
	right int
}

type EventHeap []Event

func (h EventHeap) Len() int { return len(h) }
func (h EventHeap) Less(i, j int) bool {
	if h[i].r == h[j].r {
		return h[i].id < h[j].id
	}
	return h[i].r < h[j].r
}
func (h EventHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *EventHeap) Push(x interface{}) { *h = append(*h, x.(Event)) }
func (h *EventHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

type SegTree struct {
	n   int
	sum []int
	min []int
}

func NewSegTree(arr []int) *SegTree {
	n := 1
	for n < len(arr) {
		n <<= 1
	}
	sum := make([]int, 2*n)
	minv := make([]int, 2*n)
	for i := 0; i < len(arr); i++ {
		sum[n+i] = arr[i]
		minv[n+i] = arr[i]
	}
	for i := n - 1; i > 0; i-- {
		l, r := i<<1, i<<1|1
		sum[i] = sum[l] + sum[r]
		tmp := minv[l]
		if sum[l]+minv[r] < tmp {
			tmp = sum[l] + minv[r]
		}
		minv[i] = tmp
	}
	return &SegTree{n: n, sum: sum, min: minv}
}

func (st *SegTree) Update(pos int, val int) {
	i := st.n + pos
	st.sum[i] = val
	st.min[i] = val
	for i >>= 1; i > 0; i >>= 1 {
		l, r := i<<1, i<<1|1
		st.sum[i] = st.sum[l] + st.sum[r]
		tmp := st.min[l]
		if st.sum[l]+st.min[r] < tmp {
			tmp = st.sum[l] + st.min[r]
		}
		st.min[i] = tmp
	}
}

func (st *SegTree) MinPrefix() int {
	return st.min[1]
}

func (st *SegTree) Total() int {
	return st.sum[1]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		m := 2 * n
		items := make([]Item, m)
		for i := 0; i < m; i++ {
			var a, b int64
			var c string
			fmt.Fscan(reader, &a, &b, &c)
			sign := -1
			if c == "(" {
				sign = 1
			}
			items[i] = Item{a: a, b: b, sign: sign}
		}

		perm := make([]int, m)
		for i := 0; i < m; i++ {
			perm[i] = i
		}
		sort.Slice(perm, func(i, j int) bool {
			ai, aj := items[perm[i]].a, items[perm[j]].a
			if ai == aj {
				bi, bj := items[perm[i]].b, items[perm[j]].b
				if bi == bj {
					return items[perm[i]].sign > items[perm[j]].sign
				}
				return bi < bj
			}
			return ai < aj
		})

		pos := make([]int, m)
		arr := make([]int, m)
		for idx, id := range perm {
			pos[id] = idx
			arr[idx] = items[id].sign
		}

		seg := NewSegTree(arr)
		if seg.Total() == 0 && seg.MinPrefix() >= 0 {
			fmt.Fprintln(writer, "YES")
			continue
		}

		h := &EventHeap{}
		heap.Init(h)
		eventID := 0
		var addEvent func(int, int, float64)
		addEvent = func(leftIdx, rightIdx int, curR float64) {
			if leftIdx < 0 || rightIdx >= m {
				return
			}
			i := perm[leftIdx]
			j := perm[rightIdx]
			if items[i].b == items[j].b {
				return
			}
			r := float64(items[j].a-items[i].a) / float64(items[i].b-items[j].b)
			if r > curR && r > 0 {
				heap.Push(h, Event{r: r, id: eventID, left: i, right: j})
				eventID++
			}
		}

		currentR := 0.0
		for i := 0; i+1 < m; i++ {
			addEvent(i, i+1, currentR)
		}

		success := false
		for h.Len() > 0 {
			ev := heap.Pop(h).(Event)
			if ev.r <= currentR {
				continue
			}
			i := ev.left
			j := ev.right
			pi := pos[i]
			pj := pos[j]
			if pj-pi != 1 {
				continue
			}
			currentR = ev.r
			perm[pi], perm[pj] = perm[pj], perm[pi]
			pos[i], pos[j] = pos[j], pos[i]
			arr[pi], arr[pj] = arr[pj], arr[pi]
			seg.Update(pi, arr[pi])
			seg.Update(pj, arr[pj])
			if seg.Total() == 0 && seg.MinPrefix() >= 0 {
				success = true
				break
			}
			addEvent(pi-1, pi, currentR)
			addEvent(pj, pj+1, currentR)
		}

		if !success {
			if seg.Total() == 0 && seg.MinPrefix() >= 0 {
				success = true
			}
		}
		if success {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
