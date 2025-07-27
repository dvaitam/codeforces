package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 1000000007

type edge struct{ u, v int }

func firstPlayerWins(n, d int, edges []edge, portals [][2]int) bool {
	total := (d + 1) * n
	adj := make([][]int, total)
	for i := 0; i <= d; i++ {
		for _, e := range edges {
			u := i*n + (e.u - 1)
			v := i*n + (e.v - 1)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}
	}
	for i := 0; i < d; i++ {
		a := portals[i][0] - 1
		b := portals[i][1] - 1
		u := i*n + a
		v := (i+1)*n + b
		adj[u] = append(adj[u], v)
	}
	type state struct {
		cur  int
		mask uint64
	}
	memo := make(map[state]bool)
	var dfs func(int, uint64) bool
	dfs = func(cur int, mask uint64) bool {
		st := state{cur, mask}
		if val, ok := memo[st]; ok {
			return val
		}
		for _, nx := range adj[cur] {
			if mask&(1<<uint(nx)) == 0 {
				if !dfs(nx, mask|1<<uint(nx)) {
					memo[st] = true
					return true
				}
			}
		}
		memo[st] = false
		return false
	}
	return dfs(0, 1)
}

func bruteForce(n, d int, edges []edge) int64 {
	portals := make([][2]int, d)
	var ans int64
	var rec func(int)
	rec = func(i int) {
		if i == d {
			if firstPlayerWins(n, d, edges, portals) {
				ans++
			}
			return
		}
		for a := 1; a <= n; a++ {
			for b := 1; b <= n; b++ {
				portals[i] = [2]int{a, b}
				rec(i + 1)
			}
		}
	}
	rec(0)
	return ans % mod
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var d int
	if _, err := fmt.Fscan(in, &n, &d); err != nil {
		return
	}
	edges := make([]edge, n-1)
	for i := 0; i < n-1; i++ {
		fmt.Fscan(in, &edges[i].u, &edges[i].v)
	}
	ans := bruteForce(n, d, edges)
	fmt.Println(ans)
}
