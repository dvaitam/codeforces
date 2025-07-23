package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type rawEvent struct {
	x  int
	y1 int
	y2 int
	d  int
}

type event struct {
	x int
	l int
	r int
	d int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	rawEvents := make([]rawEvent, 0, n*2)
	yBounds := make([]int, 0, n*2)
	points := make([][2]int, n)
	for i := 0; i < n; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		points[i][0] = x
		points[i][1] = y
		l := x - k + 1
		r := x + 1
		b1 := y - k + 1
		b2 := y
		rawEvents = append(rawEvents, rawEvent{l, b1, b2, 1})
		rawEvents = append(rawEvents, rawEvent{r, b1, b2, -1})
		yBounds = append(yBounds, b1, b2+1)
	}

	sort.Ints(yBounds)
	yBounds = uniqueInts(yBounds)
	m := len(yBounds)
	// map boundary value to index
	idx := make(map[int]int, m)
	for i, v := range yBounds {
		idx[v] = i
	}

	events := make([]event, len(rawEvents))
	for i, e := range rawEvents {
		l := idx[e.y1]
		r := idx[e.y2+1]
		events[i] = event{e.x, l, r, e.d}
	}
	sort.Slice(events, func(i, j int) bool { return events[i].x < events[j].x })

	segLen := make([]int, m-1)
	for i := 0; i < m-1; i++ {
		segLen[i] = yBounds[i+1] - yBounds[i]
	}

	cover := make([]int, m-1)
	hist := make([]int64, n+1)
	res := make([]int64, n+1)

	update := func(l, r, delta int) {
		for i := l; i < r; i++ {
			old := cover[i]
			if old > 0 {
				hist[old] -= int64(segLen[i])
			}
			cover[i] = old + delta
			if cover[i] > 0 {
				hist[cover[i]] += int64(segLen[i])
			}
		}
	}

	xPrev := events[0].x
	i := 0
	for i < len(events) {
		xCur := events[i].x
		width := xCur - xPrev
		if width > 0 {
			for c := 1; c <= n; c++ {
				if hist[c] != 0 {
					res[c] += hist[c] * int64(width)
				}
			}
		}
		for i < len(events) && events[i].x == xCur {
			ev := events[i]
			update(ev.l, ev.r, ev.d)
			i++
		}
		xPrev = xCur
	}

	for c := 1; c <= n; c++ {
		if c > 1 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, res[c])
	}
	fmt.Fprintln(out)
}

func uniqueInts(a []int) []int {
	if len(a) == 0 {
		return a
	}
	j := 1
	for i := 1; i < len(a); i++ {
		if a[i] != a[i-1] {
			a[j] = a[i]
			j++
		}
	}
	return a[:j]
}
