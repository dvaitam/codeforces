package main

import (
	"bufio"
	"fmt"
	"os"
)

type Interval struct {
	l int
	r int
}

type Edge struct {
	to int
	l  int
	r  int
}

func insertInterval(arr []Interval, L, R int) ([]Interval, []Interval) {
	if L > R {
		return arr, nil
	}
	newArr := make([]Interval, 0, len(arr)+1)
	var added []Interval
	i := 0
	for i < len(arr) && arr[i].r < L-1 {
		newArr = append(newArr, arr[i])
		i++
	}
	start, end := L, R
	addStart := L
	for i < len(arr) && arr[i].l <= R+1 {
		if arr[i].l > addStart {
			limit := arr[i].l - 1
			if limit > R {
				limit = R
			}
			if addStart <= limit {
				added = append(added, Interval{addStart, limit})
			}
		}
		if arr[i].l < start {
			start = arr[i].l
		}
		if arr[i].r > end {
			end = arr[i].r
		}
		addStart = arr[i].r + 1
		i++
	}
	if addStart <= R {
		added = append(added, Interval{addStart, R})
	}
	newArr = append(newArr, Interval{start, end})
	for i < len(arr) {
		newArr = append(newArr, arr[i])
		i++
	}
	return newArr, added
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	l := make([]int, n+1)
	r := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &l[i], &r[i])
	}
	adj := make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		L := max(l[u], l[v])
		R := min(r[u], r[v])
		if L <= R {
			adj[u] = append(adj[u], Edge{to: v, l: L, r: R})
			adj[v] = append(adj[v], Edge{to: u, l: L, r: R})
		}
	}

	intervals := make([][]Interval, n+1)
	intervals[1] = []Interval{{l[1], r[1]}}
	type Item struct {
		v int
		l int
		r int
	}
	queue := []Item{{1, l[1], r[1]}}

	for len(queue) > 0 {
		it := queue[0]
		queue = queue[1:]
		v := it.v
		L := it.l
		R := it.r
		for _, e := range adj[v] {
			L2 := max(L, e.l)
			R2 := min(R, e.r)
			if L2 > R2 {
				continue
			}
			arr, added := insertInterval(intervals[e.to], L2, R2)
			if len(added) > 0 {
				intervals[e.to] = arr
				for _, seg := range added {
					queue = append(queue, Item{e.to, seg.l, seg.r})
				}
			}
		}
	}

	ans := []int{}
	for i := 1; i <= n; i++ {
		if len(intervals[i]) > 0 {
			ans = append(ans, i)
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
