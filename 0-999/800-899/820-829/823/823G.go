package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"sort"
)

type Edge struct {
	to int
	l  int
	r  int
}

type Interval struct {
	l int
	r int
}

var adj [][]Edge
var reach [][]Interval

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func insertInterval(arr []Interval, iv Interval) []Interval {
	n := len(arr)
	pos := sort.Search(n, func(i int) bool { return arr[i].r >= iv.l-1 })
	start := iv.l
	end := iv.r
	i := pos
	for i < n && arr[i].l <= end+1 {
		if arr[i].l < start {
			start = arr[i].l
		}
		if arr[i].r > end {
			end = arr[i].r
		}
		i++
	}
	res := make([]Interval, 0, n-(i-pos)+1)
	res = append(res, arr[:pos]...)
	res = append(res, Interval{start, end})
	res = append(res, arr[i:]...)
	return res
}

func subtractInterval(iv Interval, arr []Interval) []Interval {
	var res []Interval
	l := iv.l
	r := iv.r
	i := sort.Search(len(arr), func(i int) bool { return arr[i].r >= l })
	for i < len(arr) && l <= r {
		cur := arr[i]
		if cur.l > r {
			break
		}
		if l < cur.l {
			end := min(r, cur.l-1)
			if end >= l {
				res = append(res, Interval{l, end})
			}
		}
		if cur.r >= r {
			l = r + 1
			break
		}
		l = cur.r + 1
		i++
	}
	if l <= r {
		res = append(res, Interval{l, r})
	}
	return res
}

func addReach(v int, seg Interval, q *list.List) {
	diff := subtractInterval(seg, reach[v])
	if len(diff) == 0 {
		return
	}
	for _, d := range diff {
		reach[v] = insertInterval(reach[v], d)
		q.PushBack(struct {
			node int
			iv   Interval
		}{v, d})
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	adj = make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		var a, b, l, r int
		fmt.Fscan(reader, &a, &b, &l, &r)
		adj[a] = append(adj[a], Edge{to: b, l: l, r: r})
		adj[b] = append(adj[b], Edge{to: a, l: l, r: r})
	}
	reach = make([][]Interval, n+1)
	q := list.New()
	addReach(1, Interval{0, 0}, q)

	for q.Len() > 0 {
		front := q.Front()
		q.Remove(front)
		item := front.Value.(struct {
			node int
			iv   Interval
		})
		v := item.node
		seg := item.iv
		for _, e := range adj[v] {
			start := max(seg.l, e.l)
			end := min(seg.r, e.r-1)
			if start <= end {
				addReach(e.to, Interval{start + 1, end + 1}, q)
			}
		}
	}
	if len(reach[n]) == 0 {
		fmt.Fprintln(writer, -1)
	} else {
		fmt.Fprintln(writer, reach[n][0].l)
	}
}
