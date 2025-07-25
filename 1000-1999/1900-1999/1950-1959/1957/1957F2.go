package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	VALMAX = 100000
	LOG    = 17 + 1
)

type Node struct {
	left, right int
	sum         int
}

var (
	seg   []Node
	roots []int
	adj   [][]int
	up    [][]int
	depth []int
	val   []int
)

func update(prev, l, r, pos int) int {
	cur := len(seg)
	seg = append(seg, seg[prev])
	if l == r {
		seg[cur].sum++
		return cur
	}
	mid := (l + r) >> 1
	if pos <= mid {
		seg[cur].left = update(seg[prev].left, l, mid, pos)
	} else {
		seg[cur].right = update(seg[prev].right, mid+1, r, pos)
	}
	seg[cur].sum = seg[seg[cur].left].sum + seg[seg[cur].right].sum
	return cur
}

func dfsIterative(root int) {
	type frame struct{ v, p, i int }
	stack := []frame{{root, 0, 0}}
	depth[root] = 0
	roots[root] = update(0, 1, VALMAX, val[root])
	for len(stack) > 0 {
		f := &stack[len(stack)-1]
		v := f.v
		if f.i < len(adj[v]) {
			to := adj[v][f.i]
			f.i++
			if to == f.p {
				continue
			}
			up[0][to] = v
			depth[to] = depth[v] + 1
			roots[to] = update(roots[v], 1, VALMAX, val[to])
			stack = append(stack, frame{to, v, 0})
		} else {
			stack = stack[:len(stack)-1]
		}
	}
}

func lca(u, v int) int {
	if depth[u] < depth[v] {
		u, v = v, u
	}
	diff := depth[u] - depth[v]
	for i := 0; diff > 0; i++ {
		if diff&1 != 0 {
			u = up[i][u]
		}
		diff >>= 1
	}
	if u == v {
		return u
	}
	for i := LOG - 1; i >= 0; i-- {
		if up[i][u] != up[i][v] {
			u = up[i][u]
			v = up[i][v]
		}
	}
	return up[0][u]
}

func getSum(a, b, c, d int) int {
	return seg[a].sum + seg[b].sum - seg[c].sum - seg[d].sum
}

func collect(a1, b1, c1, d1, a2, b2, c2, d2, l, r int, k *int, res *[]int) {
	if *k == 0 {
		return
	}
	diff := getSum(a1, b1, c1, d1) - getSum(a2, b2, c2, d2)
	if diff == 0 {
		return
	}
	if l == r {
		*res = append(*res, l)
		*k--
		return
	}
	mid := (l + r) >> 1
	la1, lb1, lc1, ld1 := seg[a1].left, seg[b1].left, seg[c1].left, seg[d1].left
	la2, lb2, lc2, ld2 := seg[a2].left, seg[b2].left, seg[c2].left, seg[d2].left
	if getSum(la1, lb1, lc1, ld1) != getSum(la2, lb2, lc2, ld2) {
		collect(la1, lb1, lc1, ld1, la2, lb2, lc2, ld2, l, mid, k, res)
		if *k == 0 {
			return
		}
	}
	ra1, rb1, rc1, rd1 := seg[a1].right, seg[b1].right, seg[c1].right, seg[d1].right
	ra2, rb2, rc2, rd2 := seg[a2].right, seg[b2].right, seg[c2].right, seg[d2].right
	if getSum(ra1, rb1, rc1, rd1) != getSum(ra2, rb2, rc2, rd2) {
		collect(ra1, rb1, rc1, rd1, ra2, rb2, rc2, rd2, mid+1, r, k, res)
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	val = make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &val[i])
	}
	adj = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	seg = make([]Node, 1)
	roots = make([]int, n+1)
	up = make([][]int, LOG)
	for i := range up {
		up[i] = make([]int, n+1)
	}
	depth = make([]int, n+1)

	dfsIterative(1)
	for i := 1; i < LOG; i++ {
		for v := 1; v <= n; v++ {
			up[i][v] = up[i-1][up[i-1][v]]
		}
	}

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var u1, v1, u2, v2, k int
		fmt.Fscan(in, &u1, &v1, &u2, &v2, &k)
		l1 := lca(u1, v1)
		l2 := lca(u2, v2)
		p1 := up[0][l1]
		p2 := up[0][l2]
		res := make([]int, 0, k)
		kval := k
		collect(roots[u1], roots[v1], roots[l1], roots[p1], roots[u2], roots[v2], roots[l2], roots[p2], 1, VALMAX, &kval, &res)
		fmt.Fprint(out, len(res))
		for _, v := range res {
			fmt.Fprint(out, " ", v)
		}
		fmt.Fprintln(out)
	}
}
