package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Event struct {
	x, y, k int
}

type Query struct {
	x, y, s int
	idx     int
}

type BIT2D struct {
	n    int
	ys   [][]int
	tree [][]int
}

func newBIT2D(n int, events []Event, queries []Query) *BIT2D {
	ys := make([][]int, n+2)
	for _, e := range events {
		for i := e.x; i <= n; i += i & -i {
			ys[i] = append(ys[i], e.y)
		}
	}
	for _, q := range queries {
		for i := q.x; i > 0; i -= i & -i {
			ys[i] = append(ys[i], q.y)
		}
	}
	tree := make([][]int, n+2)
	for i := 1; i <= n; i++ {
		if len(ys[i]) == 0 {
			continue
		}
		sort.Ints(ys[i])
		pos := 1
		for j := 1; j < len(ys[i]); j++ {
			if ys[i][j] != ys[i][pos-1] {
				ys[i][pos] = ys[i][j]
				pos++
			}
		}
		ys[i] = ys[i][:pos]
		tree[i] = make([]int, pos+2)
	}
	return &BIT2D{n: n, ys: ys, tree: tree}
}

func (b *BIT2D) Add(x, y, delta int) {
	for i := x; i <= b.n; i += i & -i {
		if len(b.ys[i]) == 0 {
			continue
		}
		idx := sort.SearchInts(b.ys[i], y) + 1
		for j := idx; j < len(b.tree[i]); j += j & -j {
			b.tree[i][j] += delta
		}
	}
}

func (b *BIT2D) Query(x, y int) int {
	res := 0
	for i := x; i > 0; i -= i & -i {
		if len(b.ys[i]) == 0 {
			continue
		}
		idx := sort.Search(len(b.ys[i]), func(j int) bool { return b.ys[i][j] > y })
		for j := idx; j > 0; j -= j & -j {
			res += b.tree[i][j]
		}
	}
	return res
}

func handleOrientation(n int, events []Event, queries []Query, ans []int) {
	if len(queries) == 0 || len(events) == 0 {
		if len(events) == 0 {
			if len(queries) == 0 {
				return
			}
			// even if no events, still need to process queries (result 0)
			return
		}
	}
	bit := newBIT2D(n, events, queries)
	sort.Slice(events, func(i, j int) bool { return events[i].k > events[j].k })
	sort.Slice(queries, func(i, j int) bool { return queries[i].s > queries[j].s })
	e := 0
	for _, q := range queries {
		for e < len(events) && events[e].k >= q.s {
			bit.Add(events[e].x, events[e].y, 1)
			e++
		}
		ans[q.idx] += bit.Query(q.x, q.y)
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, Q int
	fmt.Fscan(in, &N, &Q)

	events := make([][]Event, 4)
	queries := make([][]Query, 4)
	ans := []int{}
	qIdx := 0

	for i := 0; i < Q; i++ {
		var tp int
		fmt.Fscan(in, &tp)
		if tp == 1 {
			var dir, x, y, l int
			fmt.Fscan(in, &dir, &x, &y, &l)
			switch dir {
			case 1:
				events[0] = append(events[0], Event{x, y, x + y + l})
			case 2:
				y1 := N + 1 - y
				events[1] = append(events[1], Event{x, y1, x + y1 + l})
			case 3:
				x1 := N + 1 - x
				events[2] = append(events[2], Event{x1, y, x1 + y + l})
			case 4:
				x1 := N + 1 - x
				y1 := N + 1 - y
				events[3] = append(events[3], Event{x1, y1, x1 + y1 + l})
			}
		} else {
			var x, y int
			fmt.Fscan(in, &x, &y)
			ans = append(ans, 0)
			queries[0] = append(queries[0], Query{x, y, x + y, qIdx})
			y1 := N + 1 - y
			queries[1] = append(queries[1], Query{x, y1, x + y1, qIdx})
			x1 := N + 1 - x
			queries[2] = append(queries[2], Query{x1, y, x1 + y, qIdx})
			queries[3] = append(queries[3], Query{x1, y1, x1 + y1, qIdx})
			qIdx++
		}
	}

	for t := 0; t < 4; t++ {
		handleOrientation(N, events[t], queries[t], ans)
	}

	for _, v := range ans {
		fmt.Fprintln(out, v)
	}
}
