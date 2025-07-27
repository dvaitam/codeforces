package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

var n int

type Interval struct{ l, r int }

type IntHeap []*Interval

func (h IntHeap) Len() int { return len(h) }

func distance(iv *Interval) int {
	if iv.l == 1 || iv.r == n {
		return iv.r - iv.l + 1
	}
	return (iv.r - iv.l) / 2
}

func seatOf(iv *Interval) int {
	if iv.l == 1 {
		return 1
	}
	if iv.r == n {
		return n
	}
	return (iv.l + iv.r) / 2
}

func (h IntHeap) Less(i, j int) bool {
	di, dj := distance(h[i]), distance(h[j])
	if di != dj {
		return di > dj
	}
	si, sj := seatOf(h[i]), seatOf(h[j])
	if si != sj {
		return si < sj
	}
	return h[i].l < h[j].l
}

func (h IntHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) { *h = append(*h, x.(*Interval)) }

func (h *IntHeap) Pop() interface{} {
	old := *h
	last := old[len(old)-1]
	*h = old[:len(old)-1]
	return last
}

func allocateSeat(h *IntHeap, start, end map[int]*Interval) int {
	for h.Len() > 0 {
		iv := heap.Pop(h).(*Interval)
		if cur, ok := start[iv.l]; !ok || cur != iv {
			continue
		}
		seat := seatOf(iv)
		delete(start, iv.l)
		delete(end, iv.r)

		if iv.l <= seat-1 {
			left := &Interval{iv.l, seat - 1}
			start[left.l] = left
			end[left.r] = left
			heap.Push(h, left)
		}
		if seat+1 <= iv.r {
			right := &Interval{seat + 1, iv.r}
			start[right.l] = right
			end[right.r] = right
			heap.Push(h, right)
		}
		return seat
	}
	return -1
}

func freeSeat(pos int, h *IntHeap, start, end map[int]*Interval) {
	l, r := pos, pos
	if left, ok := end[pos-1]; ok {
		l = left.l
		delete(start, left.l)
		delete(end, left.r)
	}
	if right, ok := start[pos+1]; ok {
		r = right.r
		delete(start, right.l)
		delete(end, right.r)
	}
	iv := &Interval{l, r}
	start[l] = iv
	end[r] = iv
	heap.Push(h, iv)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var m int
	fmt.Fscan(in, &n, &m)

	h := &IntHeap{}
	heap.Init(h)

	start := make(map[int]*Interval)
	end := make(map[int]*Interval)

	initInterval := &Interval{1, n}
	start[1] = initInterval
	end[n] = initInterval
	heap.Push(h, initInterval)

	carPos := make(map[int]int)

	for i := 0; i < m; i++ {
		var t, id int
		fmt.Fscan(in, &t, &id)
		if t == 1 {
			seat := allocateSeat(h, start, end)
			carPos[id] = seat
			fmt.Fprintln(out, seat)
		} else {
			seat := carPos[id]
			delete(carPos, id)
			freeSeat(seat, h, start, end)
		}
	}
}