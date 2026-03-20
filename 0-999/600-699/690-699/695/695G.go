package main

import (
	"bufio"
	"fmt"
	"os"
)

type SegNode struct {
	l, r int32
}

var seg []SegNode

func newNode() int32 {
	seg = append(seg, SegNode{})
	return int32(len(seg) - 1)
}

func insert(l, r, val int32) int32 {
	z := newNode()
	if l == r {
		return z
	}
	mid := (l + r) / 2
	if val <= mid {
		seg[z].l = insert(l, mid, val)
	} else {
		seg[z].r = insert(mid+1, r, val)
	}
	return z
}

func merge(x, y, l, r int32) int32 {
	if x == 0 || y == 0 {
		return x ^ y
	}
	z := newNode()
	if l == r {
		return z
	}
	mid := (l + r) / 2
	seg[z].l = merge(seg[x].l, seg[y].l, l, mid)
	seg[z].r = merge(seg[x].r, seg[y].r, mid+1, r)
	return z
}

func query(rt, l, r, ql, qr int32) bool {
	if rt == 0 || ql > r || qr < l {
		return false
	}
	if ql <= l && r <= qr {
		return true
	}
	mid := (l + r) / 2
	return query(seg[rt].l, l, mid, ql, qr) || query(seg[rt].r, mid+1, r, ql, qr)
}

type State struct {
	len, link int32
	next      [26]int32
	pos       int32
}

var st []State
var sz, last int32
var root []int32

func samInit(n int32) {
	st = make([]State, 1, 2*n+1)
	st[0].len = 0
	st[0].link = -1
	sz = 1
	last = 0
	root = make([]int32, 2*n+1)
}

func samExtend(c int32, pos int32, n int32) {
	cur := sz
	sz++
	st = append(st, State{})
	st[cur].len = st[last].len + 1
	st[cur].pos = pos
	root[cur] = insert(1, n, pos)

	p := last
	for p != -1 && st[p].next[c] == 0 {
		st[p].next[c] = cur
		p = st[p].link
	}

	if p == -1 {
		st[cur].link = 0
	} else {
		q := st[p].next[c]
		if st[p].len+1 == st[q].len {
			st[cur].link = q
		} else {
			clone := sz
			sz++
			st = append(st, State{})
			st[clone].len = st[p].len + 1
			st[clone].next = st[q].next
			st[clone].link = st[q].link
			st[clone].pos = st[q].pos

			for p != -1 && st[p].next[c] == q {
				st[p].next[c] = clone
				p = st[p].link
			}
			st[q].link = clone
			st[cur].link = clone
		}
	}
	last = cur
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(reader, &n)
	var w string
	fmt.Fscan(reader, &w)

	n32 := int32(n)

	seg = make([]SegNode, 1, 10000000)
	samInit(n32)

	for i := 0; i < n; i++ {
		samExtend(int32(w[i]-'a'), int32(i+1), n32)
	}

	order := make([]int32, sz)
	cnt := make([]int32, n+1)
	for i := int32(0); i < sz; i++ {
		cnt[st[i].len]++
	}
	for i := 1; i <= n; i++ {
		cnt[i] += cnt[i-1]
	}
	for i := sz - 1; i >= 0; i-- {
		cnt[st[i].len]--
		order[cnt[st[i].len]] = i
	}

	for i := sz - 1; i > 0; i-- {
		u := order[i]
		p := st[u].link
		if p != -1 {
			root[p] = merge(root[p], root[u], 1, n32)
		}
	}

	dp := make([]int32, sz)
	top := make([]int32, sz)
	ans := int32(1)

	for i := int32(1); i < sz; i++ {
		u := order[i]
		p := st[u].link
		if p == 0 {
			dp[u] = 1
			top[u] = u
		} else {
			v := top[p]
			pos := st[u].pos
			if query(root[v], 1, n32, pos-st[u].len+st[v].len, pos-1) {
				dp[u] = dp[p] + 1
				top[u] = u
			} else {
				dp[u] = dp[p]
				top[u] = top[p]
			}
		}
		if dp[u] > ans {
			ans = dp[u]
		}
	}

	fmt.Println(ans)
}
