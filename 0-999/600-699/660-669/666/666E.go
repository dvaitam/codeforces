package main

import (
	"bufio"
	"fmt"
	"os"
)

const ALPHA = 26

type state struct {
	next   [ALPHA]int
	link   int
	length int
}

var sam []state
var last int
var sz int

func saInit(maxLen int) {
	sam = make([]state, 2*maxLen+5)
	last = 1
	sz = 2
	sam[1].link = 0
	sam[1].length = 0
}

func saExtend(c int) {
	cur := sz
	sz++
	sam[cur].length = sam[last].length + 1
	p := last
	for p > 0 && sam[p].next[c] == 0 {
		sam[p].next[c] = cur
		p = sam[p].link
	}
	if p == 0 {
		sam[cur].link = 1
	} else {
		q := sam[p].next[c]
		if sam[p].length+1 == sam[q].length {
			sam[cur].link = q
		} else {
			clone := sz
			sz++
			sam[clone] = sam[q]
			sam[clone].length = sam[p].length + 1
			for p > 0 && sam[p].next[c] == q {
				sam[p].next[c] = clone
				p = sam[p].link
			}
			sam[q].link = clone
			sam[cur].link = clone
		}
	}
	last = cur
}

// persistent/dynamic segment tree node

type node struct {
	left, right *node
	val         int
	idx         int
}

func update(n *node, l, r, pos int) *node {
	if n == nil {
		n = &node{idx: l}
	}
	if l == r {
		n.val++
		return n
	}
	mid := (l + r) / 2
	if pos <= mid {
		n.left = update(n.left, l, mid, pos)
	} else {
		n.right = update(n.right, mid+1, r, pos)
	}
	v1, i1 := 0, l
	if n.left != nil {
		v1 = n.left.val
		i1 = n.left.idx
	}
	v2, i2 := 0, mid+1
	if n.right != nil {
		v2 = n.right.val
		i2 = n.right.idx
	}
	if v1 > v2 || (v1 == v2 && i1 < i2) {
		n.val = v1
		n.idx = i1
	} else {
		n.val = v2
		n.idx = i2
	}
	return n
}

func merge(a, b *node, l, r int) *node {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if l == r {
		a.val += b.val
		return a
	}
	mid := (l + r) / 2
	a.left = merge(a.left, b.left, l, mid)
	a.right = merge(a.right, b.right, mid+1, r)
	v1, i1 := 0, l
	if a.left != nil {
		v1 = a.left.val
		i1 = a.left.idx
	}
	v2, i2 := 0, mid+1
	if a.right != nil {
		v2 = a.right.val
		i2 = a.right.idx
	}
	if v1 > v2 || (v1 == v2 && i1 < i2) {
		a.val = v1
		a.idx = i1
	} else {
		a.val = v2
		a.idx = i2
	}
	return a
}

func query(n *node, l, r, L, R int) (int, int) {
	if L > r || R < l {
		return 0, 0
	}
	if n == nil {
		if L <= l && r <= R {
			return 0, l
		}
		mid := (l + r) / 2
		if R <= mid {
			return query(nil, l, mid, L, R)
		}
		if L > mid {
			return query(nil, mid+1, r, L, R)
		}
		v1, i1 := query(nil, l, mid, L, R)
		v2, i2 := query(nil, mid+1, r, L, R)
		if v1 > v2 || (v1 == v2 && i1 < i2) {
			return v1, i1
		}
		return v2, i2
	}
	if L <= l && r <= R {
		return n.val, n.idx
	}
	mid := (l + r) / 2
	if R <= mid {
		return query(n.left, l, mid, L, R)
	}
	if L > mid {
		return query(n.right, mid+1, r, L, R)
	}
	v1, i1 := query(n.left, l, mid, L, R)
	v2, i2 := query(n.right, mid+1, r, L, R)
	if v1 > v2 || (v1 == v2 && i1 < i2) {
		return v1, i1
	}
	return v2, i2
}

type Query struct {
	l, r int
	id   int
}

var (
	children [][]int
	seg      []*node
	qlist    [][]Query
	ansIdx   []int
	ansCnt   []int
	mval     int
)

func dfs(v int) *node {
	root := seg[v]
	for _, to := range children[v] {
		child := dfs(to)
		root = merge(root, child, 1, mval)
	}
	seg[v] = root
	for _, q := range qlist[v] {
		cnt, idx := query(root, 1, mval, q.l, q.r)
		ansIdx[q.id] = idx
		ansCnt[q.id] = cnt
	}
	return root
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	saInit(len(s))
	pos := make([]int, len(s)+1)
	for i := 0; i < len(s); i++ {
		saExtend(int(s[i] - 'a'))
		pos[i+1] = last
	}

	var m int
	fmt.Fscan(in, &m)
	texts := make([]string, m+1)
	for i := 1; i <= m; i++ {
		fmt.Fscan(in, &texts[i])
	}
	mval = m
	children = make([][]int, sz)
	for i := 2; i < sz; i++ {
		p := sam[i].link
		children[p] = append(children[p], i)
	}
	seg = make([]*node, sz)
	for idx := 1; idx <= m; idx++ {
		cur := 1
		l := 0
		t := texts[idx]
		for j := 0; j < len(t); j++ {
			c := int(t[j] - 'a')
			for cur > 1 && sam[cur].next[c] == 0 {
				cur = sam[cur].link
				l = sam[cur].length
			}
			if sam[cur].next[c] != 0 {
				cur = sam[cur].next[c]
				l++
			} else {
				cur = 1
				l = 0
			}
			seg[cur] = update(seg[cur], 1, m, idx)
		}
	}

	var q int
	fmt.Fscan(in, &q)
	qlist = make([][]Query, sz)
	ansIdx = make([]int, q)
	ansCnt = make([]int, q)
	for i := 0; i < q; i++ {
		var l, r, pl, pr int
		fmt.Fscan(in, &l, &r, &pl, &pr)
		length := pr - pl + 1
		state := pos[pr]
		for sam[sam[state].link].length >= length {
			state = sam[state].link
		}
		qlist[state] = append(qlist[state], Query{l, r, i})
	}

	dfs(1)
	for i := 0; i < q; i++ {
		fmt.Fprintf(out, "%d %d\n", ansIdx[i], ansCnt[i])
	}
}
