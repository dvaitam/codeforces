package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// record represents a tower projected onto the plane orthogonal to the
// viewing direction. The segment [l, r] describes all rays intersecting the
// unit square of the tower, k is the order of the tower along the viewing
// direction and h is its height.
type record struct {
	k, l, r, h int64
}

const inf = int64(1e8)

type node struct {
	s, tag int64
	lc, rc int
}

var tree []node

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func ensure(idx *int) {
	if *idx == 0 {
		*idx = len(tree)
		tree = append(tree, node{})
	}
}

func change(idx int, l, r, ll, rr, v int64) int {
	ensure(&idx)
	nd := &tree[idx]
	if ll <= l && r <= rr {
		nd.s = max(nd.s, v)
		nd.tag = max(nd.tag, v)
		return idx
	}
	mid := (l + r) >> 1
	if ll <= mid {
		nd.lc = change(nd.lc, l, mid, ll, rr, v)
	}
	if rr > mid {
		nd.rc = change(nd.rc, mid+1, r, ll, rr, v)
	}
	left := int64(0)
	if nd.lc != 0 {
		left = tree[nd.lc].s
	}
	right := int64(0)
	if nd.rc != 0 {
		right = tree[nd.rc].s
	}
	nd.s = max(nd.tag, min(left, right))
	return idx
}

func query(idx int, l, r, ll, rr int64) int64 {
	if idx == 0 {
		return 0
	}
	nd := tree[idx]
	if ll <= l && r <= rr {
		return nd.s
	}
	mid := (l + r) >> 1
	res := inf * 20
	if ll <= mid {
		res = min(res, query(nd.lc, l, mid, ll, rr))
	}
	if rr > mid {
		res = min(res, query(nd.rc, mid+1, r, ll, rr))
	}
	return max(res, nd.tag)
}

func add(records *[]record, i, j int64, vx, vy, h int64) {
	k := int64(1 << 60)
	l := int64(1 << 60)
	r := int64(-1 << 60)
	for p := i - 1; p <= i; p++ {
		for q := j - 1; q <= j; q++ {
			k = min(k, vx*p+vy*q)
			t := vy*p - vx*q
			l = min(l, t)
			r = max(r, t)
		}
	}
	*records = append(*records, record{k: k, l: l, r: r - 1, h: h})
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int64
	var vx, vy int64
	if _, err := fmt.Fscan(reader, &n, &vx, &vy); err != nil {
		return
	}
	records := make([]record, 0, n*n)
	for i := int64(1); i <= n; i++ {
		for j := int64(1); j <= n; j++ {
			var h int64
			fmt.Fscan(reader, &h)
			add(&records, i, j, vx, vy, h)
		}
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].k < records[j].k
	})

	tree = make([]node, 2, 1<<20) // index 0 unused, 1 is root
	root := 1
	var ans int64
	for _, rec := range records {
		ll := rec.l + inf
		rr := rec.r + inf
		cur := query(root, 1, inf*2, ll, rr)
		if rec.h > cur {
			ans += rec.h - cur
		}
		root = change(root, 1, inf*2, ll, rr, rec.h)
	}
	fmt.Fprint(writer, ans)
}
