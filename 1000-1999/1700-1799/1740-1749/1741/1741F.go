package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Segment struct {
	l, r int
	c    int
	id   int
}

const INF = int(1e18)

type Event struct {
	x int
	c int
	d int
}

type BITMax struct {
	n    int
	tree []int
}

func NewBITMax(n int) *BITMax {
	b := &BITMax{n: n, tree: make([]int, n+2)}
	for i := range b.tree {
		b.tree[i] = -INF
	}
	return b
}

func (b *BITMax) Update(idx, val int) {
	for idx <= b.n {
		if val > b.tree[idx] {
			b.tree[idx] = val
		}
		idx += idx & -idx
	}
}

func (b *BITMax) Query(idx int) int {
	res := -INF
	for idx > 0 {
		if b.tree[idx] > res {
			res = b.tree[idx]
		}
		idx -= idx & -idx
	}
	return res
}

type BITMin struct {
	n    int
	tree []int
}

func NewBITMin(n int) *BITMin {
	b := &BITMin{n: n, tree: make([]int, n+2)}
	for i := range b.tree {
		b.tree[i] = INF
	}
	return b
}

func (b *BITMin) Update(idx, val int) {
	for idx <= b.n {
		if val < b.tree[idx] {
			b.tree[idx] = val
		}
		idx += idx & -idx
	}
}

func (b *BITMin) Query(idx int) int {
	res := INF
	for idx > 0 {
		if b.tree[idx] < res {
			res = b.tree[idx]
		}
		idx -= idx & -idx
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		segs := make([]Segment, n)
		maxC := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &segs[i].l, &segs[i].r, &segs[i].c)
			segs[i].id = i
			if segs[i].c > maxC {
				maxC = segs[i].c
			}
		}

		ans := make([]int, n)
		for i := range ans {
			ans[i] = INF
		}

		// prepare hot zones intervals more efficiently
		events := make([]Event, 0, 2*n)
		for _, s := range segs {
			events = append(events, Event{s.l, s.c, 1})
			events = append(events, Event{s.r + 1, s.c, -1})
		}
		sort.Slice(events, func(i, j int) bool {
			if events[i].x == events[j].x {
				return events[i].d > events[j].d
			}
			return events[i].x < events[j].x
		})
		cnt := make([]int, maxC+1)
		active := 0
		hot := make([][2]int, 0)
		prev := events[0].x
		for i := 0; i < len(events); {
			x := events[i].x
			if prev < x && active >= 2 {
				hot = append(hot, [2]int{prev, x - 1})
			}
			for ; i < len(events) && events[i].x == x; i++ {
				e := events[i]
				if e.d == 1 {
					if cnt[e.c] == 0 {
						active++
					}
					cnt[e.c]++
				} else {
					cnt[e.c]--
					if cnt[e.c] == 0 {
						active--
					}
				}
			}
			prev = x
		}
		// check intersection with hot zones
		sort.Slice(segs, func(i, j int) bool { return segs[i].l < segs[j].l })
		p := 0
		for _, s := range segs {
			for p < len(hot) && hot[p][1] < s.l {
				p++
			}
			if p < len(hot) && hot[p][0] <= s.r {
				if ans[s.id] > 0 {
					ans[s.id] = 0
				}
			}
		}

		// left to right sweep
		sort.Slice(segs, func(i, j int) bool { return segs[i].l < segs[j].l })
		bitL := NewBITMax(maxC)
		bitLRev := NewBITMax(maxC)
		maxR := make([]int, maxC+1)
		for _, s := range segs {
			leftMax := bitL.Query(s.c - 1)
			revIdx := maxC - s.c
			tmp := bitLRev.Query(revIdx)
			if tmp > leftMax {
				leftMax = tmp
			}
			if leftMax > -INF {
				d := s.l - leftMax
				if d < ans[s.id] {
					ans[s.id] = d
				}
			}
			if s.r > maxR[s.c] {
				maxR[s.c] = s.r
				bitL.Update(s.c, s.r)
				bitLRev.Update(maxC-s.c+1, s.r)
			}
		}

		// right to left sweep
		sort.Slice(segs, func(i, j int) bool { return segs[i].r > segs[j].r })
		bitR := NewBITMin(maxC)
		bitRRev := NewBITMin(maxC)
		minL := make([]int, maxC+1)
		for i := range minL {
			minL[i] = INF
		}
		for _, s := range segs {
			rightMin := bitR.Query(s.c - 1)
			tmp := bitRRev.Query(maxC - s.c)
			if tmp < rightMin {
				rightMin = tmp
			}
			if rightMin < INF {
				d := rightMin - s.r
				if d < ans[s.id] {
					ans[s.id] = d
				}
			}
			if s.l < minL[s.c] {
				minL[s.c] = s.l
				bitR.Update(s.c, s.l)
				bitRRev.Update(maxC-s.c+1, s.l)
			}
		}

		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, ans[i])
		}
		fmt.Fprintln(writer)
	}
}
