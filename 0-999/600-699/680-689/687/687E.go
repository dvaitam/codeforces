package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	n, m    int
	adj     [][]int
	dfn     []int
	low     []int
	comp    []int
	onstack []bool
	stack   []int
	idx     int
	compCnt int
	comps   [][]int
	outdeg  []int
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func tarjan(v int) {
	idx++
	dfn[v] = idx
	low[v] = idx
	stack = append(stack, v)
	onstack[v] = true
	for _, to := range adj[v] {
		if dfn[to] == 0 {
			tarjan(to)
			if low[to] < low[v] {
				low[v] = low[to]
			}
		} else if onstack[to] {
			if dfn[to] < low[v] {
				low[v] = dfn[to]
			}
		}
	}
	if low[v] == dfn[v] {
		compCnt++
		var arr []int
		for {
			x := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			onstack[x] = false
			comp[x] = compCnt
			arr = append(arr, x)
			if x == v {
				break
			}
		}
		comps = append(comps, arr)
	}
}

func bfsCycle(nodes []int) int {
	best := 1 << 30
	nodeSet := make(map[int]bool)
	for _, x := range nodes {
		nodeSet[x] = true
	}
	for _, start := range nodes {
		dist := make([]int, n+1)
		for i := range dist {
			dist[i] = -1
		}
		q := make([]int, 0, len(nodes))
		dist[start] = 0
		q = append(q, start)
		for head := 0; head < len(q); head++ {
			u := q[head]
			if dist[u]+1 >= best { // prune
				continue
			}
			for _, v := range adj[u] {
				if !nodeSet[v] {
					continue
				}
				if v == start {
					if dist[u]+1 < best {
						best = dist[u] + 1
					}
				}
				if dist[v] == -1 {
					dist[v] = dist[u] + 1
					q = append(q, v)
				}
			}
		}
		if best == 2 {
			return 2
		}
	}
	return best
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n, &m)
	adj = make([][]int, n+1)
	outdeg = make([]int, n+1)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		outdeg[u]++
	}
	dfn = make([]int, n+1)
	low = make([]int, n+1)
	comp = make([]int, n+1)
	onstack = make([]bool, n+1)
	for i := 1; i <= n; i++ {
		if dfn[i] == 0 {
			tarjan(i)
		}
	}

	hasOut := make([]bool, compCnt+1)
	hasZeroOut := make([]bool, compCnt+1)
	for u := 1; u <= n; u++ {
		if outdeg[u] == 0 {
			hasZeroOut[comp[u]] = true
		}
		for _, v := range adj[u] {
			if comp[u] != comp[v] {
				hasOut[comp[u]] = true
			}
		}
	}

	ans := n
	for id := 1; id <= compCnt; id++ {
		if hasOut[id] || hasZeroOut[id] {
			continue
		}
		g := bfsCycle(comps[id-1])
		ans += 998*g + 1
	}
	fmt.Println(ans)
}
