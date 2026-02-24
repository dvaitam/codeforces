package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	rdr := bufio.NewScanner(os.Stdin)
	rdr.Split(bufio.ScanWords)
	rdr.Buffer(make([]byte, 1024*1024), 1024*1024)

	readInt := func() int {
		if !rdr.Scan() {
			return 0
		}
		res := 0
		for _, b := range rdr.Bytes() {
			res = res*10 + int(b-'0')
		}
		return res
	}

	n := readInt()
	m := readInt()
	if n == 0 {
		return
	}

	type Edge struct {
		to, id int
	}
	graph := make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		u := readInt()
		v := readInt()
		graph[u] = append(graph[u], Edge{v, i})
		graph[v] = append(graph[v], Edge{u, i})
	}

	tin := make([]int, n+1)
	low := make([]int, n+1)
	timer := 0
	isBridge := make([]bool, m)

	var dfsBridge func(u, pEdge int)
	dfsBridge = func(u, pEdge int) {
		timer++
		tin[u] = timer
		low[u] = timer
		for _, e := range graph[u] {
			if e.id == pEdge {
				continue
			}
			if tin[e.to] != 0 {
				if tin[e.to] < low[u] {
					low[u] = tin[e.to]
				}
			} else {
				dfsBridge(e.to, e.id)
				if low[e.to] < low[u] {
					low[u] = low[e.to]
				}
				if low[e.to] > tin[u] {
					isBridge[e.id] = true
				}
			}
		}
	}
	if n > 0 {
		dfsBridge(1, -1)
	}

	comp := make([]int, n+1)
	compCnt := 0
	var dfsComp func(u, c int)
	dfsComp = func(u, c int) {
		comp[u] = c
		for _, e := range graph[u] {
			if !isBridge[e.id] && comp[e.to] == 0 {
				dfsComp(e.to, c)
			}
		}
	}
	for i := 1; i <= n; i++ {
		if comp[i] == 0 {
			compCnt++
			dfsComp(i, compCnt)
		}
	}

	W := make([]int, compCnt+1)
	for i := 1; i <= n; i++ {
		W[comp[i]]++
	}

	tree := make([][]int, compCnt+1)
	for u := 1; u <= n; u++ {
		for _, e := range graph[u] {
			if isBridge[e.id] {
				cu := comp[u]
				cv := comp[e.to]
				tree[cu] = append(tree[cu], cv)
			}
		}
	}

	maxPairs := int64(0)

	for root := 1; root <= compCnt; root++ {
		S := make([]int, compCnt+1)
		var dfs func(u, p int)
		dfs = func(u, p int) {
			S[u] = W[u]
			for _, v := range tree[u] {
				if v != p {
					dfs(v, u)
					S[u] += S[v]
				}
			}
		}
		dfs(root, 0)

		V := int64(0)
		for i := 1; i <= compCnt; i++ {
			if i != root {
				V += int64(W[i]) * int64(S[i])
			}
		}

		dp := make([]bool, n+1)
		dp[0] = true
		for _, v := range tree[root] {
			sz := S[v]
			for j := n; j >= sz; j-- {
				if dp[j-sz] {
					dp[j] = true
				}
			}
		}

		bestXY := int64(0)
		rem := n - W[root]
		for i := 0; i <= rem; i++ {
			if dp[i] {
				xy := int64(i) * int64(rem-i)
				if xy > bestXY {
					bestXY = xy
				}
			}
		}

		pairs := V + int64(W[root])*int64(n) + bestXY
		if pairs > maxPairs {
			maxPairs = pairs
		}
	}

	fmt.Println(maxPairs)
}
