package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD = 998244353

var (
	n     int
	adj   [][]int
	edges [][2]int
	k     int
	s     []int
	ans   int
)

func bfs(start, blockA, blockB int, alive []bool) ([]int, int) {
	q := make([]int, 0)
	q = append(q, start)
	visited := make([]bool, n+1)
	visited[start] = true
	comp := make([]int, 0)
	comp = append(comp, start)
	for idx := 0; idx < len(q); idx++ {
		v := q[idx]
		for _, to := range adj[v] {
			if (v == blockA && to == blockB) || (v == blockB && to == blockA) {
				continue
			}
			if !alive[to] || visited[to] {
				continue
			}
			visited[to] = true
			q = append(q, to)
			comp = append(comp, to)
		}
	}
	return comp, len(comp)
}

func dfs(alive []bool, curSize int, step int) {
	if step == k {
		ans++
		if ans >= MOD {
			ans -= MOD
		}
		return
	}
	target := s[step]
	for _, e := range edges {
		u, v := e[0], e[1]
		if !alive[u] || !alive[v] {
			continue
		}
		comp, size1 := bfs(u, u, v, alive)
		size2 := curSize - size1
		if size1 == target {
			newAlive := make([]bool, n+1)
			for _, node := range comp {
				newAlive[node] = true
			}
			dfs(newAlive, target, step+1)
		}
		if size2 == target {
			newAlive := make([]bool, n+1)
			copy(newAlive, alive)
			for _, node := range comp {
				newAlive[node] = false
			}
			dfs(newAlive, target, step+1)
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &n)
	adj = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
		edges = append(edges, [2]int{a, b})
	}
	fmt.Fscan(reader, &k)
	s = make([]int, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(reader, &s[i])
	}
	alive := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		alive[i] = true
	}
	ans = 0
	dfs(alive, n, 0)
	fmt.Println(ans % MOD)
}
