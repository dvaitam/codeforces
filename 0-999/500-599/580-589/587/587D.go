package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	u, v  int
	color int
	t     int
	idx   int
}

type TwoSat struct {
	n  int
	g  [][]int
	gr [][]int
}

func NewTwoSat(n int) *TwoSat {
	g := make([][]int, 2*n)
	gr := make([][]int, 2*n)
	return &TwoSat{n: n, g: g, gr: gr}
}

func (ts *TwoSat) addImp(a, b int) {
	ts.g[a] = append(ts.g[a], b)
	ts.gr[b] = append(ts.gr[b], a)
}

func (ts *TwoSat) addOr(x int, xv bool, y int, yv bool) {
	a := 2 * x
	if !xv {
		a ^= 1
	}
	b := 2 * y
	if !yv {
		b ^= 1
	}
	ts.addImp(a^1, b)
	ts.addImp(b^1, a)
}

func (ts *TwoSat) addTrue(x int, val bool) {
	ts.addOr(x, val, x, val)
}

func (ts *TwoSat) newVar() int {
	id := ts.n
	ts.n++
	ts.g = append(ts.g, nil, nil)
	ts.gr = append(ts.gr, nil, nil)
	return id
}

func (ts *TwoSat) solve() ([]bool, bool) {
	n := ts.n
	order := make([]int, 0, 2*n)
	used := make([]bool, 2*n)
	var dfs1 func(int)
	dfs1 = func(v int) {
		used[v] = true
		for _, to := range ts.g[v] {
			if !used[to] {
				dfs1(to)
			}
		}
		order = append(order, v)
	}
	for i := 0; i < 2*n; i++ {
		if !used[i] {
			dfs1(i)
		}
	}
	comp := make([]int, 2*n)
	for i := range comp {
		comp[i] = -1
	}
	var dfs2 func(int, int)
	dfs2 = func(v, c int) {
		comp[v] = c
		for _, to := range ts.gr[v] {
			if comp[to] == -1 {
				dfs2(to, c)
			}
		}
	}
	cid := 0
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		if comp[v] == -1 {
			dfs2(v, cid)
			cid++
		}
	}
	res := make([]bool, n)
	for i := 0; i < n; i++ {
		if comp[2*i] == comp[2*i^1] {
			return nil, false
		}
		res[i] = comp[2*i] > comp[2*i^1]
	}
	return res, true
}

func addAtMostOne(ts *TwoSat, vars []int) {
	k := len(vars)
	if k <= 1 {
		return
	}
	prev := ts.newVar()
	ts.addOr(vars[0], false, prev, true) // ~x1 or s1
	for i := 1; i < k-1; i++ {
		cur := ts.newVar()
		ts.addOr(vars[i], false, cur, true)   // ~xi or si
		ts.addOr(prev, false, cur, true)      // ~prev or cur
		ts.addOr(vars[i], false, prev, false) // ~xi or ~prev
		prev = cur
	}
	ts.addOr(vars[k-1], false, prev, false) // ~xk or ~prev
}

func check(edges []Edge, n int, limit int) (bool, []int) {
	colorFixed := make(map[int]map[int]bool)
	maxVar := 0
	varIndex := make([]int, len(edges))
	for i := range edges {
		varIndex[i] = -1
	}
	for _, e := range edges {
		if e.t > limit {
			m := colorFixed[e.color]
			if m == nil {
				m = make(map[int]bool)
				colorFixed[e.color] = m
			}
			if m[e.u] || m[e.v] {
				return false, nil
			}
			m[e.u] = true
			m[e.v] = true
		}
	}
	for i, e := range edges {
		if e.t <= limit {
			varIndex[i] = maxVar
			maxVar++
		}
	}
	ts := NewTwoSat(maxVar)
	vertexEdges := make(map[int][]int)
	for i, e := range edges {
		if e.t <= limit {
			idx := varIndex[i]
			vertexEdges[e.u] = append(vertexEdges[e.u], idx)
			vertexEdges[e.v] = append(vertexEdges[e.v], idx)
		}
	}
	for _, list := range vertexEdges {
		addAtMostOne(ts, list)
	}
	colorVertexVars := make(map[int]map[int][]int)
	for i, e := range edges {
		if e.t <= limit {
			if m, ok := colorFixed[e.color]; ok && (m[e.u] || m[e.v]) {
				ts.addTrue(varIndex[i], true)
				continue
			}
			mm := colorVertexVars[e.color]
			if mm == nil {
				mm = make(map[int][]int)
				colorVertexVars[e.color] = mm
			}
			mm[e.u] = append(mm[e.u], varIndex[i])
			mm[e.v] = append(mm[e.v], varIndex[i])
		}
	}
	for _, vm := range colorVertexVars {
		for _, list := range vm {
			l := len(list)
			for i := 0; i < l; i++ {
				for j := i + 1; j < l; j++ {
					ts.addOr(list[i], true, list[j], true) // xi OR xj
				}
			}
		}
	}
	assign, ok := ts.solve()
	if !ok {
		return false, nil
	}
	var removed []int
	for i, e := range edges {
		if e.t <= limit && assign[varIndex[i]] {
			removed = append(removed, e.idx)
		}
	}
	return true, removed
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	edges := make([]Edge, m)
	maxT := 0
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &edges[i].u, &edges[i].v, &edges[i].color, &edges[i].t)
		edges[i].u--
		edges[i].v--
		edges[i].idx = i + 1
		if edges[i].t > maxT {
			maxT = edges[i].t
		}
	}
	lo, hi := 0, maxT
	resT := -1
	for lo <= hi {
		mid := (lo + hi) / 2
		ok, _ := check(edges, n, mid)
		if ok {
			resT = mid
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	if resT == -1 {
		fmt.Fprintln(out, "No")
		return
	}
	ok, removed := check(edges, n, resT)
	if !ok {
		fmt.Fprintln(out, "No")
		return
	}
	fmt.Fprintln(out, "Yes")
	fmt.Fprintf(out, "%d %d\n", resT, len(removed))
	for i, id := range removed {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, id)
	}
	if len(removed) > 0 {
		fmt.Fprintln(out)
	}
}
