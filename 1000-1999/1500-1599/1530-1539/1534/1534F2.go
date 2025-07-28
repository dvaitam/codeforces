package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct{ l, r int }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		grid[i] = []byte(s)
	}
	a := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &a[i])
	}

	id := make([][]int, n)
	for i := range id {
		id[i] = make([]int, m)
	}
	row := make([]int, 0)
	col := make([]int, 0)
	idx := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '#' {
				id[i][j] = idx
				row = append(row, i)
				col = append(col, j)
				idx++
			} else {
				id[i][j] = -1
			}
		}
	}
	nodes := idx
	if nodes == 0 {
		fmt.Fprintln(out, 0)
		return
	}

	nextBelow := make([][]int, n+1)
	for i := range nextBelow {
		nextBelow[i] = make([]int, m)
	}
	for j := 0; j < m; j++ {
		next := -1
		for i := n - 1; i >= 0; i-- {
			if grid[i][j] == '#' {
				next = i
			}
			nextBelow[i][j] = next
		}
		nextBelow[n][j] = -1
	}

	g := make([][]int, nodes)
	rg := make([][]int, nodes)
	addEdge := func(u, v int) {
		g[u] = append(g[u], v)
		rg[v] = append(rg[v], u)
	}

	for v := 0; v < nodes; v++ {
		i := row[v]
		j := col[v]
		if i+1 < n {
			r := nextBelow[i+1][j]
			if r != -1 {
				addEdge(v, id[r][j])
			}
		}
		if j-1 >= 0 {
			r := nextBelow[i][j-1]
			if r != -1 {
				addEdge(v, id[r][j-1])
			}
		}
		if j+1 < m {
			r := nextBelow[i][j+1]
			if r != -1 {
				addEdge(v, id[r][j+1])
			}
		}
	}

	visited := make([]bool, nodes)
	order := make([]int, 0, nodes)
	var dfs1 func(int)
	dfs1 = func(v int) {
		visited[v] = true
		for _, to := range g[v] {
			if !visited[to] {
				dfs1(to)
			}
		}
		order = append(order, v)
	}
	for i := 0; i < nodes; i++ {
		if !visited[i] {
			dfs1(i)
		}
	}

	comp := make([]int, nodes)
	for i := range comp {
		comp[i] = -1
	}
	compCnt := 0
	var dfs2 func(int)
	dfs2 = func(v int) {
		comp[v] = compCnt
		for _, to := range rg[v] {
			if comp[to] == -1 {
				dfs2(to)
			}
		}
	}
	for i := nodes - 1; i >= 0; i-- {
		v := order[i]
		if comp[v] == -1 {
			dfs2(v)
			compCnt++
		}
	}

	dag := make([][]int, compCnt)
	indeg := make([]int, compCnt)
	for u := 0; u < nodes; u++ {
		cu := comp[u]
		for _, v := range g[u] {
			cv := comp[v]
			if cu != cv {
				dag[cu] = append(dag[cu], cv)
				indeg[cv]++
			}
		}
	}
	for i := 0; i < compCnt; i++ {
		if len(dag[i]) > 1 {
			mp := make(map[int]struct{})
			nw := dag[i][:0]
			for _, v := range dag[i] {
				if _, ok := mp[v]; !ok {
					mp[v] = struct{}{}
					nw = append(nw, v)
				}
			}
			dag[i] = nw
		}
	}

	L := make([]int, compCnt)
	R := make([]int, compCnt)
	for i := 0; i < compCnt; i++ {
		L[i] = 1 << 30
		R[i] = -1
	}

	blocks := make([][]int, m)
	for j := 0; j < m; j++ {
		for i := 0; i < n; i++ {
			if grid[i][j] == '#' {
				blocks[j] = append(blocks[j], i)
			}
		}
	}

	reqComp := make([]int, m)
	for j := 0; j < m; j++ {
		reqComp[j] = -1
		if a[j] > 0 {
			pos := len(blocks[j]) - a[j]
			r := blocks[j][pos]
			c := comp[id[r][j]]
			reqComp[j] = c
			if L[c] > j {
				L[c] = j
			}
			if R[c] < j {
				R[c] = j
			}
		}
	}

	queue := make([]int, 0, compCnt)
	for i := 0; i < compCnt; i++ {
		if indeg[i] == 0 {
			queue = append(queue, i)
		}
	}
	topo := make([]int, 0, compCnt)
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		topo = append(topo, v)
		for _, to := range dag[v] {
			indeg[to]--
			if indeg[to] == 0 {
				queue = append(queue, to)
			}
		}
	}

	for i := len(topo) - 1; i >= 0; i-- {
		v := topo[i]
		for _, to := range dag[v] {
			if L[to] < L[v] {
				L[v] = L[to]
			}
			if R[to] > R[v] {
				R[v] = R[to]
			}
		}
	}

	mp := make(map[int]pair)
	for j := 0; j < m; j++ {
		c := reqComp[j]
		if c != -1 {
			mp[c] = pair{L[c], R[c]}
		}
	}

	intervals := make([]pair, 0, len(mp))
	for _, p := range mp {
		intervals = append(intervals, p)
	}

	sort.Slice(intervals, func(i, j int) bool {
		if intervals[i].r == intervals[j].r {
			return intervals[i].l < intervals[j].l
		}
		return intervals[i].r < intervals[j].r
	})

	ans := 0
	cur := -1
	for _, seg := range intervals {
		if seg.l > cur {
			ans++
			cur = seg.r
		}
	}

	fmt.Fprintln(out, ans)
}
