package main

import (
	"bufio"
	"fmt"
	"os"
)

// Helper for building implication graph

type rule struct {
	p1, t1 int
	p2, t2 int
}

var (
	n      int
	baseG  [][]int
	baseRG [][]int
)

func lit(pos, val int) int {
	return pos*2 + val
}

func addEdge(g, rg [][]int, u, v int) {
	g[u] = append(g[u], v)
	rg[v] = append(rg[v], u)
}

func buildGraphs(rules []rule) {
	baseG = make([][]int, 2*n)
	baseRG = make([][]int, 2*n)
	for _, r := range rules {
		u1 := lit(r.p1, r.t1)
		v1 := lit(r.p2, r.t2)
		addEdge(baseG, baseRG, u1, v1)
		u2 := lit(r.p2, 1-r.t2)
		v2 := lit(r.p1, 1-r.t1)
		addEdge(baseG, baseRG, u2, v2)
	}
}

func check(assign []int) bool {
	g := make([][]int, len(baseG))
	rg := make([][]int, len(baseRG))
	for i := range g {
		if len(baseG[i]) > 0 {
			g[i] = append([]int(nil), baseG[i]...)
		}
	}
	for i := range rg {
		if len(baseRG[i]) > 0 {
			rg[i] = append([]int(nil), baseRG[i]...)
		}
	}
	for i, v := range assign {
		if v == -1 {
			continue
		}
		if v == 1 {
			u := lit(i, 0)
			w := lit(i, 1)
			g[u] = append(g[u], w)
			rg[w] = append(rg[w], u)
		} else {
			u := lit(i, 1)
			w := lit(i, 0)
			g[u] = append(g[u], w)
			rg[w] = append(rg[w], u)
		}
	}

	m := len(g)
	order := make([]int, 0, m)
	vis := make([]bool, m)
	var dfs1 func(int)
	dfs1 = func(v int) {
		vis[v] = true
		for _, to := range g[v] {
			if !vis[to] {
				dfs1(to)
			}
		}
		order = append(order, v)
	}
	for i := 0; i < m; i++ {
		if !vis[i] {
			dfs1(i)
		}
	}
	comp := make([]int, m)
	for i := range comp {
		comp[i] = -1
	}
	var dfs2 func(int, int)
	dfs2 = func(v, c int) {
		comp[v] = c
		for _, to := range rg[v] {
			if comp[to] == -1 {
				dfs2(to, c)
			}
		}
	}
	cid := 0
	for i := m - 1; i >= 0; i-- {
		v := order[i]
		if comp[v] == -1 {
			dfs2(v, cid)
			cid++
		}
	}
	for i := 0; i < n; i++ {
		if comp[lit(i, 0)] == comp[lit(i, 1)] {
			return false
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var typeStr string
	fmt.Fscan(reader, &typeStr)
	l := len(typeStr)
	letterType := make([]int, l)
	for i, ch := range typeStr {
		if ch == 'V' {
			letterType[i] = 1
		} else {
			letterType[i] = 0
		}
	}
	fmt.Fscan(reader, &n)
	var mRules int
	fmt.Fscan(reader, &mRules)
	rules := make([]rule, mRules)
	for i := 0; i < mRules; i++ {
		var p1, p2 int
		var t1c, t2c string
		fmt.Fscan(reader, &p1, &t1c, &p2, &t2c)
		t1 := 0
		if t1c[0] == 'V' {
			t1 = 1
		}
		t2 := 0
		if t2c[0] == 'V' {
			t2 = 1
		}
		rules[i] = rule{p1 - 1, t1, p2 - 1, t2}
	}
	var s string
	fmt.Fscan(reader, &s)
	buildGraphs(rules)

	assign := make([]int, n)
	for i := range assign {
		assign[i] = -1
	}

	equal := true
	res := make([]byte, n)

	for i := 0; i < n; i++ {
		start := byte('a')
		if equal {
			start = s[i]
		}
		found := false
		for ch := start; ch < byte('a'+l); ch++ {
			t := letterType[int(ch-'a')]
			assign[i] = t
			if check(assign) {
				res[i] = ch
				if ch > s[i] {
					equal = false
				}
				found = true
				break
			}
			assign[i] = -1
		}
		if !found {
			fmt.Println("-1")
			return
		}
	}
	fmt.Println(string(res))
}
