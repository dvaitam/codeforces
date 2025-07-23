package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Node struct {
	left, right int
	sum         int64
}

type Event struct {
	x   int
	idx int
	da  int64
	dc  int64
}

var slopeNodes []Node
var constNodes []Node

func newNode(nodes *[]Node) int {
	*nodes = append(*nodes, Node{})
	return len(*nodes) - 1
}

func update(nodes *[]Node, prev int, l, r, pos int, delta int64) int {
	cur := newNode(nodes)
	(*nodes)[cur] = (*nodes)[prev]
	if l == r {
		(*nodes)[cur].sum += delta
		return cur
	}
	mid := (l + r) >> 1
	if pos <= mid {
		left := update(nodes, (*nodes)[prev].left, l, mid, pos, delta)
		(*nodes)[cur].left = left
	} else {
		right := update(nodes, (*nodes)[prev].right, mid+1, r, pos, delta)
		(*nodes)[cur].right = right
	}
	(*nodes)[cur].sum = (*nodes)[(*nodes)[cur].left].sum + (*nodes)[(*nodes)[cur].right].sum
	return cur
}

func query(nodes []Node, root, l, r, ql, qr int) int64 {
	if root == 0 || ql > r || qr < l {
		return 0
	}
	if ql <= l && r <= qr {
		return nodes[root].sum
	}
	mid := (l + r) >> 1
	if qr <= mid {
		return query(nodes, nodes[root].left, l, mid, ql, qr)
	}
	if ql > mid {
		return query(nodes, nodes[root].right, mid+1, r, ql, qr)
	}
	return query(nodes, nodes[root].left, l, mid, ql, qr) + query(nodes, nodes[root].right, mid+1, r, ql, qr)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	type Func struct {
		x1, x2       int
		y1, a, b, y2 int64
	}
	funcs := make([]Func, n+1)
	events := make([]Event, 0, 2*n)

	for i := 1; i <= n; i++ {
		var x1, x2 int
		var y1, a, b, y2 int64
		fmt.Fscan(in, &x1, &x2, &y1, &a, &b, &y2)
		funcs[i] = Func{x1, x2, y1, a, b, y2}
		events = append(events, Event{x1 + 1, i, a, b - y1})
		events = append(events, Event{x2 + 1, i, -a, y2 - b})
	}

	slopeNodes = append(slopeNodes, Node{})
	constNodes = append(constNodes, Node{})
	slopeRoot := 0
	constRoot := 0
	for i := 1; i <= n; i++ {
		constRoot = update(&constNodes, constRoot, 1, n, i, funcs[i].y1)
	}

	sort.Slice(events, func(i, j int) bool { return events[i].x < events[j].x })
	times := []int{0}
	slopeRoots := []int{slopeRoot}
	constRoots := []int{constRoot}
	idx := 0
	for idx < len(events) {
		x := events[idx].x
		curSlope := slopeRoots[len(slopeRoots)-1]
		curConst := constRoots[len(constRoots)-1]
		for idx < len(events) && events[idx].x == x {
			e := events[idx]
			curSlope = update(&slopeNodes, curSlope, 1, n, e.idx, e.da)
			curConst = update(&constNodes, curConst, 1, n, e.idx, e.dc)
			idx++
		}
		times = append(times, x)
		slopeRoots = append(slopeRoots, curSlope)
		constRoots = append(constRoots, curConst)
	}

	var m int
	fmt.Fscan(in, &m)
	last := int64(0)
	const mod int64 = 1000000000

	for i := 0; i < m; i++ {
		var l, r int
		var x int64
		fmt.Fscan(in, &l, &r, &x)
		xi := (x + last) % mod
		j := sort.Search(len(times), func(k int) bool { return int64(times[k]) > xi })
		ver := j - 1
		slopeSum := query(slopeNodes, slopeRoots[ver], 1, n, l, r)
		constSum := query(constNodes, constRoots[ver], 1, n, l, r)
		ans := slopeSum*xi + constSum
		fmt.Fprintln(out, ans)
		last = ans
	}
}
