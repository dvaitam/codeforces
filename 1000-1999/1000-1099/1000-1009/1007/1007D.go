package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	MaxN  = 100200
	MaxTJ = 4000000 // 2-SAT nodes: needs to be large for segment tree optimization
)

var (
	// Tree data structures
	h, nxt, to         []int
	num                int
	dep, size, son, fa [MaxN]int
	dfn, top           [MaxN]int
	hldCounter         int
	
	// Segment tree nodes storing 2-SAT variable indices
	vc [MaxN * 4][]int
	
	n, m, totalTJ int
	tjSolver      *Tarjan
)

type Tarjan struct {
	h, dfn, low, st, bel []int
	wayTo, wayNxt        []int
	in                   []bool
	num, tot, top, idx   int
}

func NewTarjan(nodes int, edges int) *Tarjan {
	return &Tarjan{
		h:      make([]int, nodes),
		dfn:    make([]int, nodes),
		low:    make([]int, nodes),
		st:     make([]int, nodes),
		bel:    make([]int, nodes),
		in:     make([]bool, nodes),
		wayTo:  make([]int, edges),
		wayNxt: make([]int, edges),
	}
}

func (tj *Tarjan) Link(x, y int) {
	if x <= 0 || y <= 0 { return }
	tj.num++
	if tj.num >= len(tj.wayTo) { return } // Safety check
	tj.wayTo[tj.num] = y
	tj.wayNxt[tj.num] = tj.h[x]
	tj.h[x] = tj.num
}

func (tj *Tarjan) DFS(x int) {
	tj.tot++
	tj.dfn[x] = tj.tot
	tj.low[x] = tj.tot
	tj.top++
	tj.st[tj.top] = x
	tj.in[x] = true

	for i := tj.h[x]; i > 0; i = tj.wayNxt[i] {
		v := tj.wayTo[i]
		if tj.dfn[v] == 0 {
			tj.DFS(v)
			if tj.low[v] < tj.low[x] {
				tj.low[x] = tj.low[v]
			}
		} else if tj.in[v] {
			if tj.dfn[v] < tj.low[x] {
				tj.low[x] = tj.dfn[v]
			}
		}
	}

	if tj.low[x] == tj.dfn[x] {
		tj.idx++
		for {
			node := tj.st[tj.top]
			tj.top--
			tj.in[node] = false
			tj.bel[node] = tj.idx
			if node == x {
				break
			}
		}
	}
}

func addEdge(x, y int) {
	num++
	to[num] = y
	nxt[num] = h[x]
	h[x] = num
	num++
	to[num] = x
	nxt[num] = h[y]
	h[y] = num
}

func dfs0(x, f int) {
	dep[x] = dep[f] + 1
	fa[x] = f
	size[x] = 1
	for i := h[x]; i > 0; i = nxt[i] {
		v := to[i]
		if v != f {
			dfs0(v, x)
			size[x] += size[v]
			if size[v] > size[son[x]] {
				son[x] = v
			}
		}
	}
}

func dfs1(x, tp int) {
	hldCounter++
	dfn[x] = hldCounter
	top[x] = tp
	if son[x] != 0 {
		dfs1(son[x], tp)
	}
	for i := h[x]; i > 0; i = nxt[i] {
		v := to[i]
		if v != fa[x] && v != son[x] {
			dfs1(v, v)
		}
	}
}

func ins(ql, qr, v, cur, l, r int) {
	if ql <= l && r <= qr {
		vc[cur] = append(vc[cur], v)
		return
	}
	mid := (l + r) >> 1
	if ql <= mid {
		ins(ql, qr, v, cur<<1, l, mid)
	}
	if qr > mid {
		ins(ql, qr, v, cur<<1|1, mid+1, r)
	}
}

func addPath(x, y, v int) {
	for top[x] != top[y] {
		if dep[top[x]] < dep[top[y]] {
			x, y = y, x
		}
		ins(dfn[top[x]], dfn[x], v, 1, 2, n)
		x = fa[top[x]]
	}
	if x != y {
		l, r := dfn[x], dfn[y]
		if l > r {
			l, r = r, l
		}
		ins(l+1, r, v, 1, 2, n)
	}
}

func dfs2(cur, now0, now1, l, r int) {
	for _, i := range vc[cur] {
		// Clause: if this edge choice is made, 
		// it implies the previous choice in this segment tree node must also be consistent
		tjSolver.Link(i, now0)
		totalTJ++
		tjSolver.Link(totalTJ, now0)
		now0 = totalTJ
		tjSolver.Link(now0, i^1)

		tjSolver.Link(now1, i^1)
		totalTJ++
		tjSolver.Link(now1, totalTJ)
		now1 = totalTJ
		tjSolver.Link(i, now1)
	}
	if l == r {
		return
	}
	mid := (l + r) >> 1
	dfs2(cur<<1, now0, now1, l, mid)
	dfs2(cur<<1|1, now0, now1, mid+1, r)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	
	h = make([]int, n+1)
	nxt = make([]int, n*2+1)
	to = make([]int, n*2+1)

	for i := 2; i <= n; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		addEdge(x, y)
	}

	dfs0(1, 0)
	dfs1(1, 1)

	fmt.Fscan(reader, &m)
	for i := 1; i <= m; i++ {
		var a, b, c, d int
		fmt.Fscan(reader, &a, &b, &c, &d)
		// 2-SAT indices: choice 1 is 2i, choice 2 is 2i+1
		addPath(a, b, i<<1)
		addPath(c, d, i<<1|1)
	}

	// 2-SAT nodes 2 to 2m+1 are for the M restrictions.
	// totalTJ will expand as we build the prefix/suffix optimization in dfs2.
	tjSolver = NewTarjan(MaxTJ, MaxTJ*3)
	totalTJ = (m + 1) << 1 | 1
	
	dfs2(1, (m+1)<<1, totalTJ, 2, n)

	for i := 2; i <= totalTJ; i++ {
		if tjSolver.dfn[i] == 0 {
			tjSolver.DFS(i)
		}
	}

	for i := 1; i <= m; i++ {
		if tjSolver.bel[i<<1] == tjSolver.bel[i<<1|1] {
			fmt.Fprintln(writer, "NO")
			return
		}
	}

	fmt.Fprintln(writer, "YES")
	for i := 1; i <= m; i++ {
		res := 1
		if tjSolver.bel[i<<1] > tjSolver.bel[i<<1|1] {
			res = 2
		}
		fmt.Fprintln(writer, res)
	}
}
