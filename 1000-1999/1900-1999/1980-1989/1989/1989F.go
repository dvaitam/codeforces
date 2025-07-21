package main

import (
	"bufio"
	"fmt"
	"os"
)

const N = 400000 + 5

type edge struct {
	x, y, t int
}

var (
	n, m, q int
	fa      [N]int
	sz      [N]int
	ans     int64

	e   [N][]int
	dfn [N]int
	low [N]int
	st  [N]int
	tp  int
	co  [N]int
	tot int
	vis [N]bool
)

func gf(x int) int {
	if fa[x] == x {
		return x
	}
	fa[x] = gf(fa[x])
	return fa[x]
}

func merge(x, y int) {
	x = gf(x)
	y = gf(y)
	if x == y {
		return
	}
	if sz[x] > 1 {
		ans -= int64(sz[x]) * int64(sz[x])
	}
	if sz[y] > 1 {
		ans -= int64(sz[y]) * int64(sz[y])
	}
	fa[y] = x
	sz[x] += sz[y]
	ans += int64(sz[x]) * int64(sz[x])
}

func tarjan(x int) {
	tot++
	dfn[x] = tot
	low[x] = tot
	tp++
	st[tp] = x
	vis[x] = true
	for _, v := range e[x] {
		if dfn[v] == 0 {
			tarjan(v)
			if low[v] < low[x] {
				low[x] = low[v]
			}
		} else if vis[v] {
			if dfn[v] < low[x] {
				low[x] = dfn[v]
			}
		}
	}
	if low[x] == dfn[x] {
		co[x] = x
		for st[tp] != x {
			co[st[tp]] = x
			vis[st[tp]] = false
			tp--
		}
		tp--
		vis[x] = false
	}
}

func solve(l, r int, ed []edge) {
	if l == r {
		if l > q {
			return
		}
		for _, v := range ed {
			merge(v.x, v.y)
		}
		fmt.Println(ans)
		return
	}
	mid := (l + r) >> 1
	tot = 0
	var el, er []edge
	for i := range ed {
		v := &ed[i]
		v.x = gf(v.x)
		v.y = gf(v.y)
		e[v.x] = nil
		e[v.y] = nil
		dfn[v.x] = 0
		dfn[v.y] = 0
	}
	for _, v := range ed {
		if v.t <= mid {
			e[v.x] = append(e[v.x], v.y)
		}
	}
	for _, v := range ed {
		if v.t <= mid {
			if dfn[v.x] == 0 {
				tarjan(v.x)
			}
			if dfn[v.y] == 0 {
				tarjan(v.y)
			}
			if co[v.x] == co[v.y] {
				el = append(el, v)
				continue
			}
		}
		er = append(er, v)
	}
	solve(l, mid, el)
	solve(mid+1, r, er)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n, &m, &q)
	for i := 1; i <= n+m; i++ {
		fa[i] = i
		sz[i] = 1
	}
	edges := make([]edge, q)
	for i := 0; i < q; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		var c string
		fmt.Fscan(in, &c)
		if c == "R" {
			edges[i] = edge{y + n, x, i + 1}
		} else {
			edges[i] = edge{x, y + n, i + 1}
		}
	}
	solve(1, q+1, edges)
}
