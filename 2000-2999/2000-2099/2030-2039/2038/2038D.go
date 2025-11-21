package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const (
	mod  = 998244353
	maxB = 31
)

type Event struct {
	w     int32
	start int32
	l     int32
	r     int32
}

type BIT struct {
	n    int
	tree []int64
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int64, n+2)}
}

func (b *BIT) add(idx int, delta int64) {
	for idx <= b.n {
		b.tree[idx] += delta
		if b.tree[idx] >= mod {
			b.tree[idx] -= mod
		} else if b.tree[idx] < 0 {
			b.tree[idx] += mod
		}
		idx += idx & -idx
	}
}

func (b *BIT) rangeAdd(l, r int, delta int64) {
	if l > r {
		return
	}
	b.add(l, delta)
	b.add(r+1, -delta)
}

func (b *BIT) pointQuery(idx int) int64 {
	res := int64(0)
	for idx > 0 {
		res += b.tree[idx]
		idx -= idx & -idx
	}
	res %= mod
	if res < 0 {
		res += mod
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int32, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	nextPos := make([][maxB]int32, n+2)
	inf := int32(n + 1)
	for b := 0; b < maxB; b++ {
		nextPos[n+1][b] = inf
	}
	for i := n; i >= 1; i-- {
		nextPos[i] = nextPos[i+1]
		val := a[i-1]
		for b := 0; b < maxB; b++ {
			if (val>>uint(b))&1 == 1 {
				nextPos[i][b] = int32(i)
			}
		}
	}

	events := make([]Event, 0, n*maxB)
	type pair struct {
		pos  int32
		mask int32
	}
	var buf [maxB]pair

	for s := 1; s <= n; s++ {
		cnt := 0
		for b := 0; b < maxB; b++ {
			pos := nextPos[s][b]
			if pos <= int32(n) {
				buf[cnt] = pair{pos: pos, mask: int32(1) << uint(b)}
				cnt++
			}
		}
		pairs := buf[:cnt]
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].pos < pairs[j].pos
		})
		current := int32(0)
		last := int32(s)
		idx := 0
		for idx < cnt {
			pos := pairs[idx].pos
			if last <= pos-1 {
				events = append(events, Event{w: current, start: int32(s - 1), l: last, r: pos - 1})
			}
			for idx < cnt && pairs[idx].pos == pos {
				current |= pairs[idx].mask
				idx++
			}
			last = pos
		}
		if last <= int32(n) {
			events = append(events, Event{w: current, start: int32(s - 1), l: last, r: int32(n)})
		}
	}

	sort.Slice(events, func(i, j int) bool {
		if events[i].w == events[j].w {
			if events[i].start == events[j].start {
				if events[i].l == events[j].l {
					return events[i].r < events[j].r
				}
				return events[i].l < events[j].l
			}
			return events[i].start < events[j].start
		}
		return events[i].w < events[j].w
	})

	bit := NewBIT(n + 2)
	bit.rangeAdd(1, 1, 1)

	for _, e := range events {
		startIdx := int(e.start) + 1
		lIdx := int(e.l) + 1
		rIdx := int(e.r) + 1
		val := bit.pointQuery(startIdx)
		if val == 0 {
			continue
		}
		bit.rangeAdd(lIdx, rIdx, val)
	}

	ans := bit.pointQuery(n + 1)
	fmt.Fprintln(out, ans%mod)
}
