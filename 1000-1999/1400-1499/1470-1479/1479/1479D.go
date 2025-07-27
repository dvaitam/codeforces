package main

import (
	"bufio"
	"fmt"
	"os"
)

const LOG int = 20

type Node struct {
	lson, rson int
	val        int
}

var (
	seg   []Node
	root  []int
	up    [][]int
	depth []int
	a     []int
	n     int
	adj   [][]int
)

func update(prev, l, r, pos int) int {
	idx := len(seg)
	seg = append(seg, Node{})
	if prev != 0 {
		seg[idx] = seg[prev]
	}
	if l == r {
		seg[idx].val ^= 1
		return idx
	}
	mid := (l + r) >> 1
	if pos <= mid {
		seg[idx].lson = update(seg[prev].lson, l, mid, pos)
	} else {
		seg[idx].rson = update(seg[prev].rson, mid+1, r, pos)
	}
	seg[idx].val = seg[seg[idx].lson].val ^ seg[seg[idx].rson].val
	return idx
}

func query(aIdx, bIdx, cIdx, dIdx, l, r, L, R int) int {
	if L <= l && r <= R {
		val := seg[aIdx].val ^ seg[bIdx].val ^ seg[cIdx].val ^ seg[dIdx].val
		if val == 0 {
			return -1
		}
		if l == r {
			return l
		}
	}
	mid := (l + r) >> 1
	if R <= mid {
		return query(seg[aIdx].lson, seg[bIdx].lson, seg[cIdx].lson, seg[dIdx].lson, l, mid, L, R)
	}
	if L > mid {
		return query(seg[aIdx].rson, seg[bIdx].rson, seg[cIdx].rson, seg[dIdx].rson, mid+1, r, L, R)
	}
	left := query(seg[aIdx].lson, seg[bIdx].lson, seg[cIdx].lson, seg[dIdx].lson, l, mid, L, R)
	if left != -1 {
		return left
	}
	return query(seg[aIdx].rson, seg[bIdx].rson, seg[cIdx].rson, seg[dIdx].rson, mid+1, r, L, R)
}

func dfs(u, p int) {
	up[0][u] = p
	for i := 1; i < LOG; i++ {
		up[i][u] = up[i-1][up[i-1][u]]
	}
	depth[u] = depth[p] + 1
	root[u] = update(root[p], 1, n, a[u])
	for _, v := range adj[u] {
		if v == p {
			continue
		}
		dfs(v, u)
	}
}

func lca(u, v int) int {
	if depth[u] < depth[v] {
		u, v = v, u
	}
	diff := depth[u] - depth[v]
	for i := LOG - 1; i >= 0; i-- {
		if diff>>i&1 == 1 {
			u = up[i][u]
		}
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

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}

	a = make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}

	adj = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		adj[x] = append(adj[x], y)
		adj[y] = append(adj[y], x)
	}

	seg = make([]Node, 1)
	root = make([]int, n+1)
	depth = make([]int, n+1)
	up = make([][]int, LOG)
	for i := 0; i < LOG; i++ {
		up[i] = make([]int, n+1)
	}

	dfs(1, 0)

	for ; q > 0; q-- {
		var u, v, lq, rq int
		fmt.Fscan(in, &u, &v, &lq, &rq)
		w := lca(u, v)
		pw := up[0][w]
		if pw == 0 {
			pw = 0
		}
		res := query(root[u], root[v], root[w], root[pw], 1, n, lq, rq)
		fmt.Fprintln(out, res)
	}
}
