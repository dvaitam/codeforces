package main

import (
	"bufio"
	"fmt"
	"os"
)

// SegTree supports point updates and range minimum queries.
type SegTree struct {
	n    int
	tree []int64
}

func NewSegTree(n int) *SegTree {
	size := 1
	for size < n {
		size <<= 1
	}
	t := &SegTree{n: size, tree: make([]int64, size<<1)}
	const INF = int64(1 << 60)
	for i := range t.tree {
		t.tree[i] = INF
	}
	return t
}

func (t *SegTree) Update(pos int, val int64) {
	idx := pos + t.n
	t.tree[idx] = val
	for idx >>= 1; idx > 0; idx >>= 1 {
		if t.tree[idx<<1] < t.tree[idx<<1|1] {
			t.tree[idx] = t.tree[idx<<1]
		} else {
			t.tree[idx] = t.tree[idx<<1|1]
		}
	}
}

func (t *SegTree) Query(l, r int) int64 {
	const INF = int64(1 << 60)
	res := INF
	l += t.n
	r += t.n + 1
	for l < r {
		if l&1 == 1 {
			if t.tree[l] < res {
				res = t.tree[l]
			}
			l++
		}
		if r&1 == 1 {
			r--
			if t.tree[r] < res {
				res = t.tree[r]
			}
		}
		l >>= 1
		r >>= 1
	}
	return res
}

var (
	g       [][]int
	tin     []int
	tout    []int
	depth   []int
	order   int
	byDepth [][]int
)

func dfs(u, p, d int) {
	depth[u] = d
	if d >= len(byDepth) {
		tmp := make([][]int, d+1)
		copy(tmp, byDepth)
		byDepth = tmp
	}
	byDepth[d] = append(byDepth[d], u)
	order++
	tin[u] = order
	for _, v := range g[u] {
		if v == p {
			continue
		}
		dfs(v, u, d+1)
	}
	tout[u] = order
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, r int
	if _, err := fmt.Fscan(reader, &n, &r); err != nil {
		return
	}
	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	g = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		g[x] = append(g[x], y)
		g[y] = append(g[y], x)
	}
	tin = make([]int, n+1)
	tout = make([]int, n+1)
	depth = make([]int, n+1)
	byDepth = make([][]int, 1)
	order = 0
	dfs(r, 0, 0)
	maxDepth := len(byDepth) - 1

	seg := NewSegTree(n + 2)
	curDepth := -1

	var m int
	fmt.Fscan(reader, &m)
	last := 0
	for i := 0; i < m; i++ {
		var p, q int
		fmt.Fscan(reader, &p, &q)
		x := (p+last)%n + 1
		k := (q + last) % n
		t := depth[x] + k
		if t > maxDepth {
			t = maxDepth
		}
		for curDepth < t {
			curDepth++
			for _, v := range byDepth[curDepth] {
				seg.Update(tin[v], a[v])
			}
		}
		ans := seg.Query(tin[x], tout[x])
		fmt.Fprintln(writer, ans)
		last = int(ans)
	}
}
