package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const LOG = 17

var (
	g     [][]int
	up    [LOG][]int
	depth []int
	tin   []int
	tout  []int
	timer int
)

type Pair struct{ a, b int }

func dfs(v, p int) {
	timer++
	tin[v] = timer
	up[0][v] = p
	for i := 1; i < LOG; i++ {
		up[i][v] = up[i-1][up[i-1][v]]
	}
	for _, to := range g[v] {
		if to == p {
			continue
		}
		depth[to] = depth[v] + 1
		dfs(to, v)
	}
	tout[v] = timer
}

func isAncestor(u, v int) bool {
	return tin[u] <= tin[v] && tout[v] <= tout[u]
}

func lca(u, v int) int {
	if isAncestor(u, v) {
		return u
	}
	if isAncestor(v, u) {
		return v
	}
	for i := LOG - 1; i >= 0; i-- {
		if !isAncestor(up[i][u], v) {
			u = up[i][u]
		}
	}
	return up[0][u]
}

func buildVirtualTree(nodes []int) ([]int, [][2]int) {
	sort.Slice(nodes, func(i, j int) bool { return tin[nodes[i]] < tin[nodes[j]] })
	lcas := make([]int, 0, len(nodes)-1)
	for i := 0; i+1 < len(nodes); i++ {
		l := lca(nodes[i], nodes[i+1])
		lcas = append(lcas, l)
	}
	nodes = append(nodes, lcas...)
	sort.Slice(nodes, func(i, j int) bool { return tin[nodes[i]] < tin[nodes[j]] })
	uniq := nodes[:0]
	for _, v := range nodes {
		if len(uniq) == 0 || uniq[len(uniq)-1] != v {
			uniq = append(uniq, v)
		}
	}
	nodes = uniq
	if len(nodes) == 1 {
		return nodes, nil
	}
	st := []int{nodes[0]}
	edges := make([][2]int, 0)
	for i := 1; i < len(nodes); i++ {
		v := nodes[i]
		l := lca(v, st[len(st)-1])
		for len(st) >= 2 && depth[st[len(st)-2]] >= depth[l] {
			a := st[len(st)-1]
			st = st[:len(st)-1]
			b := st[len(st)-1]
			edges = append(edges, [2]int{b, a})
		}
		if st[len(st)-1] != l {
			a := st[len(st)-1]
			edges = append(edges, [2]int{l, a})
			st[len(st)-1] = l
		}
		st = append(st, v)
	}
	for len(st) > 1 {
		a := st[len(st)-1]
		st = st[:len(st)-1]
		b := st[len(st)-1]
		edges = append(edges, [2]int{b, a})
	}
	return nodes, edges
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		g = make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
		}
		for i := 0; i < LOG; i++ {
			up[i] = make([]int, n+1)
		}
		depth = make([]int, n+1)
		tin = make([]int, n+1)
		tout = make([]int, n+1)
		timer = 0
		dfs(1, 1)
		var k int
		fmt.Fscan(reader, &k)
		pairs := make([]Pair, k)
		nodes := make([]int, 0, 2*k+1)
		nodes = append(nodes, 1)
		for i := 0; i < k; i++ {
			fmt.Fscan(reader, &pairs[i].a, &pairs[i].b)
			nodes = append(nodes, pairs[i].a, pairs[i].b)
		}
		for i := 0; i < k; i++ {
			l := lca(pairs[i].a, pairs[i].b)
			nodes = append(nodes, l)
		}
		_, vtEdges := buildVirtualTree(nodes)
		maskSet := make(map[int]struct{})
		for _, e := range vtEdges {
			parent := e[0]
			child := e[1]
			if depth[parent] > depth[child] {
				parent, child = child, parent
			}
			mask := 0
			for j := 0; j < k; j++ {
				aIn := isAncestor(child, pairs[j].a)
				bIn := isAncestor(child, pairs[j].b)
				if aIn != bIn {
					mask |= 1 << uint(j)
				}
			}
			if mask != 0 {
				maskSet[mask] = struct{}{}
			}
		}
		uniqueMasks := make([]int, 0, len(maskSet))
		for m := range maskSet {
			uniqueMasks = append(uniqueMasks, m)
		}
		full := (1 << uint(k)) - 1
		INF := k + 1
		dp := make([]int, full+1)
		for i := range dp {
			dp[i] = INF
		}
		dp[0] = 0
		for _, m := range uniqueMasks {
			for s := full; ; s-- {
				ns := s | m
				if dp[s]+1 < dp[ns] {
					dp[ns] = dp[s] + 1
				}
				if s == 0 {
					break
				}
			}
		}
		fmt.Fprintln(writer, dp[full])
	}
}
