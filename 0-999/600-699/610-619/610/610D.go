package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Interval struct{ l, r int64 }

func mergeIntervals(arr []Interval) []Interval {
	if len(arr) == 0 {
		return nil
	}
	sort.Slice(arr, func(i, j int) bool {
		if arr[i].l == arr[j].l {
			return arr[i].r < arr[j].r
		}
		return arr[i].l < arr[j].l
	})
	res := make([]Interval, 0, len(arr))
	cur := arr[0]
	for _, v := range arr[1:] {
		if v.l <= cur.r+1 {
			if v.r > cur.r {
				cur.r = v.r
			}
		} else {
			res = append(res, cur)
			cur = v
		}
	}
	res = append(res, cur)
	return res
}

func unionLen(arr []Interval) int64 {
	var sum int64
	for _, v := range arr {
		sum += v.r - v.l + 1
	}
	return sum
}

type BIT struct {
	n    int
	tree []int
}

func NewBIT(n int) *BIT { return &BIT{n: n, tree: make([]int, n+2)} }
func (b *BIT) Add(idx, delta int) {
	for idx <= b.n {
		b.tree[idx] += delta
		idx += idx & -idx
	}
}
func (b *BIT) Sum(idx int) int {
	s := 0
	for idx > 0 {
		s += b.tree[idx]
		idx &= idx - 1
	}
	return s
}
func (b *BIT) Range(l, r int) int {
	if r < l {
		return 0
	}
	return b.Sum(r) - b.Sum(l-1)
}

const (
	REM = iota
	ADD
	QUERY
)

type Event struct {
	x   int64
	typ int
	y1  int64
	y2  int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	vert := make(map[int64][]Interval)
	horiz := make(map[int64][]Interval)

	for i := 0; i < n; i++ {
		var x1, y1, x2, y2 int64
		fmt.Fscan(in, &x1, &y1, &x2, &y2)
		if x1 == x2 {
			if y1 > y2 {
				y1, y2 = y2, y1
			}
			vert[x1] = append(vert[x1], Interval{y1, y2})
		} else {
			if x1 > x2 {
				x1, x2 = x2, x1
			}
			horiz[y1] = append(horiz[y1], Interval{x1, x2})
		}
	}

	var vCells, hCells int64
	for x, arr := range vert {
		merged := mergeIntervals(arr)
		vert[x] = merged
		vCells += unionLen(merged)
	}
	for y, arr := range horiz {
		merged := mergeIntervals(arr)
		horiz[y] = merged
		hCells += unionLen(merged)
	}

	ys := make([]int64, 0, len(horiz))
	for y := range horiz {
		ys = append(ys, y)
	}
	sort.Slice(ys, func(i, j int) bool { return ys[i] < ys[j] })
	pos := make(map[int64]int, len(ys))
	for i, v := range ys {
		pos[v] = i + 1
	}

	events := make([]Event, 0)
	for y, arr := range horiz {
		for _, iv := range arr {
			events = append(events, Event{iv.l, ADD, y, 0})
			events = append(events, Event{iv.r + 1, REM, y, 0})
		}
	}
	for x, arr := range vert {
		for _, iv := range arr {
			events = append(events, Event{x, QUERY, iv.l, iv.r})
		}
	}

	sort.Slice(events, func(i, j int) bool {
		if events[i].x == events[j].x {
			return events[i].typ < events[j].typ
		}
		return events[i].x < events[j].x
	})

	bit := NewBIT(len(ys) + 2)
	var inter int64
	for _, e := range events {
		switch e.typ {
		case REM:
			if idx, ok := pos[e.y1]; ok {
				bit.Add(idx, -1)
			}
		case ADD:
			if idx, ok := pos[e.y1]; ok {
				bit.Add(idx, 1)
			}
		case QUERY:
			l := sort.Search(len(ys), func(i int) bool { return ys[i] >= e.y1 })
			r := sort.Search(len(ys), func(i int) bool { return ys[i] > e.y2 }) - 1
			if l <= r {
				inter += int64(bit.Range(l+1, r+1))
			}
		}
	}

	res := vCells + hCells - inter
	fmt.Fprintln(out, res)
}
