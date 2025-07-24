package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Car struct {
	id        int
	availFrom int64
}

type CarHeap []Car

func (h CarHeap) Len() int { return len(h) }
func (h CarHeap) Less(i, j int) bool {
	if h[i].availFrom == h[j].availFrom {
		return h[i].id < h[j].id
	}
	return h[i].availFrom < h[j].availFrom
}
func (h CarHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *CarHeap) Push(x interface{}) { *h = append(*h, x.(Car)) }
func (h *CarHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// Busy event for a car currently servicing a ride
type Busy struct {
	time int64
	id   int
	pos  int
}

type BusyHeap []Busy

func (h BusyHeap) Len() int { return len(h) }
func (h BusyHeap) Less(i, j int) bool {
	if h[i].time == h[j].time {
		return h[i].id < h[j].id
	}
	return h[i].time < h[j].time
}
func (h BusyHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *BusyHeap) Push(x interface{}) { *h = append(*h, x.(Busy)) }
func (h *BusyHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// Segment tree to find nearest positions with available cars
type SegTree struct {
	n, size int
	tree    []int
}

func NewSegTree(n int) *SegTree {
	size := 1
	for size < n {
		size <<= 1
	}
	return &SegTree{n: n, size: size, tree: make([]int, 2*size)}
}

func (st *SegTree) add(pos, delta int) {
	p := pos + st.size - 1
	st.tree[p] += delta
	for p > 1 {
		p >>= 1
		st.tree[p] = st.tree[p<<1] + st.tree[p<<1|1]
	}
}

func (st *SegTree) findPrevRec(l, r, pos, idx int) int {
	if pos < l || st.tree[idx] == 0 {
		return -1
	}
	if l == r {
		if l <= st.n {
			return l
		}
		return -1
	}
	mid := (l + r) / 2
	if pos > mid {
		res := st.findPrevRec(mid+1, r, pos, idx*2+1)
		if res != -1 {
			return res
		}
	}
	return st.findPrevRec(l, mid, pos, idx*2)
}

func (st *SegTree) findNextRec(l, r, pos, idx int) int {
	if pos > r || st.tree[idx] == 0 {
		return -1
	}
	if l == r {
		if l <= st.n {
			return l
		}
		return -1
	}
	mid := (l + r) / 2
	if pos <= mid {
		res := st.findNextRec(l, mid, pos, idx*2)
		if res != -1 {
			return res
		}
	}
	return st.findNextRec(mid+1, r, pos, idx*2+1)
}

func (st *SegTree) findPrev(pos int) int { return st.findPrevRec(1, st.size, pos, 1) }
func (st *SegTree) findNext(pos int) int { return st.findNextRec(1, st.size, pos, 1) }

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k, m int
	if _, err := fmt.Fscan(in, &n, &k, &m); err != nil {
		return
	}

	heaps := make([]*CarHeap, n+1)
	seg := NewSegTree(n)

	for i := 1; i <= k; i++ {
		var x int
		fmt.Fscan(in, &x)
		if heaps[x] == nil {
			h := &CarHeap{}
			heap.Init(h)
			heaps[x] = h
		}
		heap.Push(heaps[x], Car{id: i, availFrom: 0})
		seg.add(x, 1)
	}

	type Req struct {
		t    int64
		a, b int
	}

	reqs := make([]Req, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &reqs[i].t, &reqs[i].a, &reqs[i].b)
	}

	busy := &BusyHeap{}
	heap.Init(busy)

	for i := 0; i < m; i++ {
		t := reqs[i].t
		a := reqs[i].a
		b := reqs[i].b
		// release cars finished by time t
		for busy.Len() > 0 && (*busy)[0].time <= t {
			e := heap.Pop(busy).(Busy)
			pos := e.pos
			id := e.id
			avail := e.time
			if heaps[pos] == nil {
				h := &CarHeap{}
				heap.Init(h)
				heaps[pos] = h
			}
			heap.Push(heaps[pos], Car{id: id, availFrom: avail})
			seg.add(pos, 1)
		}

		// if no cars available, wait for next available time
		if seg.tree[1] == 0 && busy.Len() > 0 {
			t = (*busy)[0].time
			for busy.Len() > 0 && (*busy)[0].time <= t {
				e := heap.Pop(busy).(Busy)
				pos := e.pos
				id := e.id
				avail := e.time
				if heaps[pos] == nil {
					h := &CarHeap{}
					heap.Init(h)
					heaps[pos] = h
				}
				heap.Push(heaps[pos], Car{id: id, availFrom: avail})
				seg.add(pos, 1)
			}
		}

		left := seg.findPrev(a)
		right := seg.findNext(a)

		choosePos := -1
		if left == -1 && right == -1 {
			// should not happen
			fmt.Fprintln(out, -1, 0)
			continue
		} else if left == -1 {
			choosePos = right
		} else if right == -1 {
			choosePos = left
		} else {
			distL := abs(a - left)
			distR := abs(right - a)
			lc := (*heaps[left])[0]
			rc := (*heaps[right])[0]
			if distL < distR {
				choosePos = left
			} else if distL > distR {
				choosePos = right
			} else {
				if lc.availFrom < rc.availFrom {
					choosePos = left
				} else if lc.availFrom > rc.availFrom {
					choosePos = right
				} else if lc.id < rc.id {
					choosePos = left
				} else {
					choosePos = right
				}
			}
		}

		h := heaps[choosePos]
		car := heap.Pop(h).(Car)
		seg.add(choosePos, -1)
		wait := t + int64(abs(choosePos-a)) - reqs[i].t
		fmt.Fprintln(out, car.id, wait)
		finish := t + int64(abs(choosePos-a)+abs(a-b))
		heap.Push(busy, Busy{time: finish, id: car.id, pos: b})
	}
}
