package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = int64(998244353)

var (
	n   int
	adj [][]int
	way [][]int64
	ans [][]int64
)

func dfs(u, p int) {
	way[u][0] = 1
	ans[u][0] = 0
	for _, v := range adj[u] {
		if v == p {
			continue
		}
		dfs(v, u)
		tmp1 := [2]int64{way[u][0], way[u][1]}
		tmp2 := [2]int64{ans[u][0], ans[u][1]}
		way[u][0], way[u][1], ans[u][0], ans[u][1] = 0, 0, 0, 0
		// combine without extra distance
		for a := 0; a < 2; a++ {
			for c := 0; c < 2; c++ {
				way[u][c] = (way[u][c] + tmp1[c]*way[v][a]) % mod
				ans[u][c] = (ans[u][c] + tmp1[c]*ans[v][a] + tmp2[c]*way[v][a]) % mod
			}
		}
		// combine with extra distance
		for a := 0; a < 2; a++ {
			for c := 0; c < 2; c++ {
				if a != 0 || c != 0 {
					way[u][c] = (way[u][c] + tmp1[c]*way[v][a]) % mod
					ans[u][c] = (ans[u][c] + tmp1[c]*ans[v][a] + tmp2[c]*way[v][a]) % mod
				} else {
					way[u][1] = (way[u][1] + tmp1[c]*way[v][a]) % mod
					ans[u][1] = (ans[u][1] + tmp1[c]*ans[v][a] + tmp2[c]*way[v][a] + 2*tmp1[c]*way[v][a]) % mod
				}
			}
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n)
	adj = make([][]int, n+1)
	for i := 1; i < n; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	way = make([][]int64, n+1)
	ans = make([][]int64, n+1)
	for i := 1; i <= n; i++ {
		way[i] = make([]int64, 2)
		ans[i] = make([]int64, 2)
	}
	dfs(1, 0)
	res := (ans[1][0] + ans[1][1]) % mod
	fmt.Println(res)
}
