package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Segment tree for range chmax (a2[v] = max(a2[v], x)) and querying min(v + a2[v])
const INF int64 = 1e18

type segNode struct {
	l, r     int
	min1     int   // minimal a2
	min2     int   // second minimal a2
	cnt1     int   // count of elements == min1
	minv1    int   // minimal v among elements == min1
	minB2    int64 // minimal b2 = v + a2
	minB2Alt int64 // minimal b2 among elements with a2 > min1
}

type segTree struct {
	n    int
	tree []segNode
}

func newSegTree(n int) *segTree {
	size := 4 * n
	st := &segTree{n: n, tree: make([]segNode, size)}
	st.build(1, 0, n-1)
	return st
}

func (st *segTree) build(idx, l, r int) {
	node := &st.tree[idx]
	node.l, node.r = l, r
	if l == r {
		node.min1 = 0
		node.min2 = int(1e9)
		node.cnt1 = 1
		node.minv1 = l
		node.minB2 = int64(l)
		node.minB2Alt = INF
		return
	}
	mid := (l + r) >> 1
	st.build(idx<<1, l, mid)
	st.build(idx<<1|1, mid+1, r)
	st.pushUp(idx)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func (st *segTree) pushUp(idx int) {
	a, b := &st.tree[idx<<1], &st.tree[idx<<1|1]
	node := &st.tree[idx]
	// min1
	if a.min1 < b.min1 {
		node.min1 = a.min1
		node.cnt1 = a.cnt1
		node.minv1 = a.minv1
	} else if b.min1 < a.min1 {
		node.min1 = b.min1
		node.cnt1 = b.cnt1
		node.minv1 = b.minv1
	} else {
		node.min1 = a.min1
		node.cnt1 = a.cnt1 + b.cnt1
		if a.minv1 < b.minv1 {
			node.minv1 = a.minv1
		} else {
			node.minv1 = b.minv1
		}
	}
	// min2
	node.min2 = int(1e9)
	if a.min1 != node.min1 {
		node.min2 = minInt(node.min2, a.min1)
	} else {
		node.min2 = minInt(node.min2, a.min2)
	}
	if b.min1 != node.min1 {
		node.min2 = minInt(node.min2, b.min1)
	} else {
		node.min2 = minInt(node.min2, b.min2)
	}
	// minB2
	node.minB2 = min64(a.minB2, b.minB2)
	// minB2Alt: merge alt from children
	m := INF
	// child a
	if a.min1 > node.min1 {
		m = min64(m, a.minB2)
	} else {
		m = min64(m, a.minB2Alt)
	}
	// child b
	if b.min1 > node.min1 {
		m = min64(m, b.minB2)
	} else {
		m = min64(m, b.minB2Alt)
	}
	node.minB2Alt = m
}

// apply raising a2 == min1 to x (x < min2)
func (st *segTree) applyChmax(idx, x int) {
	node := &st.tree[idx]
	if node.min1 >= x {
		return
	}
	// x < min2 guaranteed when called
	_ = x - node.min1
	node.min1 = x
	// minv1 unchanged
	// update minB2: if it was from min1 group or alt
	// new b2 from min1 group = minv1 + x
	newB2_1 := int64(node.minv1 + x)
	node.minB2 = min64(newB2_1, node.minB2Alt)
}

func (st *segTree) pushDown(idx int) {
	node := &st.tree[idx]
	for _, c := range []int{node.min1} {
		// propagate to children
		child := &st.tree[idx<<1]
		if child.min1 < c {
			st.applyChmax(idx<<1, c)
		}
		child = &st.tree[idx<<1|1]
		if child.min1 < c {
			st.applyChmax(idx<<1|1, c)
		}
	}
}

func (st *segTree) rangeChmax(lq, rq, x int) {
	st._chmax(1, lq, rq, x)
}
func (st *segTree) _chmax(idx, lq, rq, x int) {
	node := &st.tree[idx]
	if node.r < lq || node.l > rq || node.min1 >= x {
		return
	}
	if lq <= node.l && node.r <= rq && node.min2 > x {
		st.applyChmax(idx, x)
		return
	}
	st.pushDown(idx)
	st._chmax(idx<<1, lq, rq, x)
	st._chmax(idx<<1|1, lq, rq, x)
	st.pushUp(idx)
}

func (st *segTree) minB2() int64 {
	return st.tree[1].minB2
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var n int
	fmt.Fscan(reader, &n)
	a := make([]int, n)
	b := make([]int, n)
	c := make([]int, n)
	vals := make([]int, 0, 3*n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		vals = append(vals, a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &b[i])
		vals = append(vals, b[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &c[i])
		vals = append(vals, c[i])
	}
	sort.Ints(vals)
	mp := make(map[int]int, len(vals))
	m := 0
	for _, v := range vals {
		if _, ok := mp[v]; !ok {
			mp[v] = m
			m++
		}
	}
	inf := n + 1
	iA := make([]int, m)
	iB := make([]int, m)
	iC := make([]int, m)
	for i := 0; i < m; i++ {
		iA[i], iB[i], iC[i] = inf, inf, inf
	}
	for i, v := range a {
		id := mp[v]
		if i+1 < iA[id] {
			iA[id] = i + 1
		}
	}
	for i, v := range b {
		id := mp[v]
		if i+1 < iB[id] {
			iB[id] = i + 1
		}
	}
	for i, v := range c {
		id := mp[v]
		if i+1 < iC[id] {
			iC[id] = i + 1
		}
	}
	addAt := make([][]struct{ b, c int }, n+2)
	for id := 0; id < m; id++ {
		ua := iA[id]
		addAt[ua] = append(addAt[ua], struct{ b, c int }{iB[id], iC[id]})
	}
	st := newSegTree(n + 1)
	ans := int64(3 * n)
	for u := n; u >= 0; u-- {
		if u+1 <= n+1 {
			for _, p := range addAt[u+1] {
				l, r := 0, p.b-1
				if r >= 0 {
					st.rangeChmax(l, r, p.c)
				}
			}
		}
		cur := st.minB2()
		tot := int64(u) + cur
		if tot < ans {
			ans = tot
		}
	}
	fmt.Fprintln(writer, ans)
}
