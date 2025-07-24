package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type BIT struct {
	n    int
	data []int
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, data: make([]int, n+2)}
}

func (b *BIT) Add(idx, val int) {
	for idx <= b.n {
		b.data[idx] += val
		idx += idx & -idx
	}
}

func (b *BIT) Sum(idx int) int {
	res := 0
	for idx > 0 {
		res += b.data[idx]
		idx -= idx & -idx
	}
	return res
}

func (b *BIT) RangeSum(l, r int) int {
	if l > r {
		return 0
	}
	return b.Sum(r) - b.Sum(l-1)
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
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		next := make([]int, n)
		rev := make([][]int, n)
		for i := 0; i < n; i++ {
			next[i] = i + 1 + a[i]
			if next[i] >= 1 && next[i] <= n {
				rev[next[i]-1] = append(rev[next[i]-1], i)
			}
		}
		// BFS to determine nodes that lead outside
		isExit := make([]bool, n)
		q := make([]int, 0)
		for i := 0; i < n; i++ {
			if next[i] < 1 || next[i] > n {
				isExit[i] = true
				q = append(q, i)
			}
		}
		for p := 0; p < len(q); p++ {
			u := q[p]
			for _, v := range rev[u] {
				if !isExit[v] {
					isExit[v] = true
					q = append(q, v)
				}
			}
		}

		// build tree of exit nodes
		parent := make([]int, n)
		children := make([][]int, n)
		roots := make([]int, 0)
		for i := 0; i < n; i++ {
			if isExit[i] {
				v := next[i]
				if v >= 1 && v <= n && isExit[v-1] {
					parent[i] = v - 1
					children[v-1] = append(children[v-1], i)
				} else {
					parent[i] = -1
					roots = append(roots, i)
				}
			} else {
				parent[i] = -1
			}
		}
		tin := make([]int, n)
		tout := make([]int, n)
		timer := 0
		type item struct{ v, idx int }
		for _, r := range roots {
			stack := []item{{r, 0}}
			timer++
			tin[r] = timer
			for len(stack) > 0 {
				top := &stack[len(stack)-1]
				if top.idx < len(children[top.v]) {
					child := children[top.v][top.idx]
					top.idx++
					timer++
					tin[child] = timer
					stack = append(stack, item{child, 0})
				} else {
					tout[top.v] = timer
					stack = stack[:len(stack)-1]
				}
			}
		}

		// prefix of exit nodes by index
		pref := make([]int, n+1)
		for i := 1; i <= n; i++ {
			pref[i] = pref[i-1]
			if isExit[i-1] {
				pref[i]++
			}
		}

		// path from 1 in original array
		visited := make([]bool, n)
		path := make([]int, 0)
		cur := 1
		for cur >= 1 && cur <= n && !visited[cur-1] {
			visited[cur-1] = true
			path = append(path, cur)
			cur = next[cur-1]
		}
		pathExit := cur < 1 || cur > n

		// queries for subtree-range counts
		type Query struct{ L, R, tl, tr, id, pos int }
		queries := make([]Query, 0)
		idMap := make(map[int]int)
		for pos, x := range path {
			if isExit[x-1] {
				L := x - n
				if L < 1 {
					L = 1
				}
				R := x + n
				if R > n {
					R = n
				}
				queries = append(queries, Query{L, R, tin[x-1], tout[x-1], len(queries), pos})
				idMap[pos] = len(queries) - 1
			}
		}

		type Point struct{ idx, tin int }
		points := make([]Point, 0)
		for i := 0; i < n; i++ {
			if isExit[i] {
				points = append(points, Point{i + 1, tin[i]})
			}
		}
		sort.Slice(points, func(i, j int) bool { return points[i].idx < points[j].idx })

		type Event struct{ z, id, typ int }
		events := make([]Event, 0, len(queries)*2)
		for _, qv := range queries {
			events = append(events, Event{qv.R, qv.id, 1})
			events = append(events, Event{qv.L - 1, qv.id, -1})
		}
		sort.Slice(events, func(i, j int) bool { return events[i].z < events[j].z })

		bit := NewBIT(timer + 2)
		res := make([]int, len(queries))
		pt := 0
		for _, e := range events {
			for pt < len(points) && points[pt].idx <= e.z {
				bit.Add(points[pt].tin, 1)
				pt++
			}
			qv := queries[e.id]
			count := bit.RangeSum(qv.tl, qv.tr)
			res[e.id] += e.typ * count
		}

		subInside := make([]int, len(path))
		for pos, x := range path {
			if isExit[x-1] {
				subInside[pos] = res[idMap[pos]]
			}
		}

		ans := 0
		if pathExit {
			ans += (n - len(path)) * (2*n + 1)
		}
		for pos, x := range path {
			L := x - n
			if L < 1 {
				L = 1
			}
			R := x + n
			if R > n {
				R = n
			}
			total := pref[R] - pref[L-1]
			sub := subInside[pos]
			ans += n + 1 + total - sub
		}
		fmt.Fprintln(out, ans)
	}
}
