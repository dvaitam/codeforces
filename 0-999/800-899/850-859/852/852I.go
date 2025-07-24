package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type Query struct {
	l, r  int
	idx   int
	lca   int
	block int
}

var (
	n      int
	gender []int
	fav    []int
	compID []int
	adj    [][]int
	euler  []int
	first  []int
	last   []int
	up     [][]int
	depth  []int
)

func dfs(u, p int) {
	up[0][u] = p
	first[u] = len(euler)
	euler = append(euler, u)
	for _, v := range adj[u] {
		if v == p {
			continue
		}
		depth[v] = depth[u] + 1
		dfs(v, u)
	}
	last[u] = len(euler)
	euler = append(euler, u)
}

func buildLCA() {
	LOG := len(up)
	for k := 1; k < LOG; k++ {
		for v := 1; v <= n; v++ {
			if up[k-1][v] != 0 {
				up[k][v] = up[k-1][up[k-1][v]]
			}
		}
	}
}

func lca(a, b int) int {
	if depth[a] < depth[b] {
		a, b = b, a
	}
	diff := depth[a] - depth[b]
	for k := len(up) - 1; k >= 0; k-- {
		if diff&(1<<uint(k)) != 0 {
			a = up[k][a]
		}
	}
	if a == b {
		return a
	}
	for k := len(up) - 1; k >= 0; k-- {
		if up[k][a] != up[k][b] {
			a = up[k][a]
			b = up[k][b]
		}
	}
	return up[0][a]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	gender = make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &gender[i])
	}
	fav = make([]int, n+1)
	vals := make([]int, n)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &fav[i])
		vals[i-1] = fav[i]
	}
	sort.Ints(vals)
	vals = unique(vals)
	compID = make([]int, n+1)
	for i := 1; i <= n; i++ {
		compID[i] = sort.SearchInts(vals, fav[i])
	}
	m := len(vals)

	adj = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	LOG := 1
	for (1 << LOG) <= n {
		LOG++
	}
	up = make([][]int, LOG)
	for i := 0; i < LOG; i++ {
		up[i] = make([]int, n+1)
	}
	depth = make([]int, n+1)
	euler = make([]int, 0, 2*n)
	first = make([]int, n+1)
	last = make([]int, n+1)
	dfs(1, 0)
	buildLCA()

	var q int
	fmt.Fscan(reader, &q)
	queries := make([]Query, q)
	blockSize := int(math.Sqrt(float64(len(euler)))) + 1
	for i := 0; i < q; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		if first[a] > first[b] {
			a, b = b, a
		}
		w := lca(a, b)
		if w == a {
			queries[i] = Query{l: first[a], r: first[b], idx: i, lca: -1}
		} else {
			queries[i] = Query{l: last[a], r: first[b], idx: i, lca: w}
		}
		queries[i].block = queries[i].l / blockSize
	}
	sort.Slice(queries, func(i, j int) bool {
		qi, qj := queries[i], queries[j]
		if qi.block != qj.block {
			return qi.block < qj.block
		}
		if qi.block%2 == 0 {
			return qi.r < qj.r
		}
		return qi.r > qj.r
	})

	countB := make([]int, m)
	countG := make([]int, m)
	visited := make([]bool, n+1)
	var total int64

	toggle := func(pos int) {
		node := euler[pos]
		val := compID[node]
		if visited[node] {
			if gender[node] == 1 {
				countB[val]--
				total -= int64(countG[val])
			} else {
				countG[val]--
				total -= int64(countB[val])
			}
			visited[node] = false
		} else {
			if gender[node] == 1 {
				total += int64(countG[val])
				countB[val]++
			} else {
				total += int64(countB[val])
				countG[val]++
			}
			visited[node] = true
		}
	}

	answers := make([]int64, q)
	curL, curR := 0, -1
	for _, qu := range queries {
		L, R := qu.l, qu.r
		for curL > L {
			curL--
			toggle(curL)
		}
		for curR < R {
			curR++
			toggle(curR)
		}
		for curL < L {
			toggle(curL)
			curL++
		}
		for curR > R {
			toggle(curR)
			curR--
		}
		res := total
		if qu.lca != -1 {
			node := qu.lca
			val := compID[node]
			if gender[node] == 1 {
				res += int64(countG[val])
			} else {
				res += int64(countB[val])
			}
		}
		answers[qu.idx] = res
	}
	for i := 0; i < q; i++ {
		fmt.Fprintln(writer, answers[i])
	}
}

func unique(a []int) []int {
	if len(a) == 0 {
		return a
	}
	j := 1
	for i := 1; i < len(a); i++ {
		if a[i] != a[i-1] {
			a[j] = a[i]
			j++
		}
	}
	return a[:j]
}
