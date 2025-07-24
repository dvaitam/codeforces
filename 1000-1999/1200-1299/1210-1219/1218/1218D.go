package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

type Edge struct {
	to int
	id int
}

type EdgeInfo struct {
	u, v, w int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	edges := make([]EdgeInfo, m)
	g := make([][]Edge, n)
	xorTotal := 0
	for i := 0; i < m; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		u--
		v--
		edges[i] = EdgeInfo{u, v, w}
		g[u] = append(g[u], Edge{v, i})
		g[v] = append(g[v], Edge{u, i})
		xorTotal ^= w
	}
	deg := make([]int, n)
	for i := 0; i < n; i++ {
		deg[i] = len(g[i])
	}
	edgeRemoved := make([]bool, m)
	queue := make([]int, 0)
	for i := 0; i < n; i++ {
		if deg[i] == 1 {
			queue = append(queue, i)
		}
	}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, e := range g[v] {
			if edgeRemoved[e.id] {
				continue
			}
			edgeRemoved[e.id] = true
			deg[e.to]--
			if deg[e.to] == 1 {
				queue = append(queue, e.to)
			}
		}
	}
	visitedV := make([]bool, n)
	visitedE := make([]bool, m)
	cycles := make([][]int, 0)
	for i := 0; i < n; i++ {
		if deg[i] > 0 && !visitedV[i] {
			cur := i
			prev := -1
			cycle := []int{}
			for {
				visitedV[cur] = true
				var next Edge
				for _, e := range g[cur] {
					if edgeRemoved[e.id] {
						continue
					}
					if e.id == prev {
						continue
					}
					next = e
					break
				}
				if visitedE[next.id] {
					break
				}
				visitedE[next.id] = true
				cycle = append(cycle, next.id)
				cur = next.to
				prev = next.id
				if cur == i {
					break
				}
			}
			cycles = append(cycles, cycle)
		}
	}
	limit := 1
	for limit <= 100000 {
		limit <<= 1
	}
	dp := make([]int64, limit)
	dp[0] = 1
	for _, cyc := range cycles {
		cnt := make(map[int]int)
		for _, id := range cyc {
			w := edges[id].w
			cnt[w]++
		}
		ndp := make([]int64, limit)
		for x := 0; x < limit; x++ {
			if dp[x] == 0 {
				continue
			}
			for w, c := range cnt {
				nx := x ^ w
				ndp[nx] = (ndp[nx] + dp[x]*int64(c)) % MOD
			}
		}
		dp = ndp
	}
	minX := 0
	for ; minX < limit; minX++ {
		if dp[minX] > 0 {
			break
		}
	}
	cost := xorTotal ^ minX
	fmt.Printf("%d %d\n", cost, dp[minX]%MOD)
}
