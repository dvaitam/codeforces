package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	n       int
	letters []byte
	adj     [][]int
	nxt     [][]int
	dp      [][]int
)

func lps(u, v int) int {
	if dp[u][v] != -1 {
		return dp[u][v]
	}
	if u == v {
		dp[u][v] = 1
		return 1
	}
	nu := nxt[u][v]
	nv := nxt[v][u]
	if nu == v { // adjacent
		if letters[u] == letters[v] {
			dp[u][v] = 2
		} else {
			dp[u][v] = 1
		}
		return dp[u][v]
	}
	res := lps(nu, v)
	if t := lps(u, nv); t > res {
		res = t
	}
	if letters[u] == letters[v] {
		if t := lps(nu, nv) + 2; t > res {
			res = t
		}
	}
	dp[u][v] = res
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)
		letters = []byte(s)
		adj = make([][]int, n)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		nxt = make([][]int, n)
		for i := 0; i < n; i++ {
			nxt[i] = make([]int, n)
			for j := 0; j < n; j++ {
				nxt[i][j] = -1
			}
			queue := make([]int, 0, n)
			visited := make([]bool, n)
			queue = append(queue, i)
			nxt[i][i] = i
			visited[i] = true
			for head := 0; head < len(queue); head++ {
				v := queue[head]
				for _, to := range adj[v] {
					if !visited[to] {
						visited[to] = true
						if v == i {
							nxt[i][to] = to
						} else {
							nxt[i][to] = nxt[i][v]
						}
						queue = append(queue, to)
					}
				}
			}
		}

		dp = make([][]int, n)
		for i := 0; i < n; i++ {
			dp[i] = make([]int, n)
			for j := 0; j < n; j++ {
				dp[i][j] = -1
			}
		}

		ans := 1
		for u := 0; u < n; u++ {
			for v := u + 1; v < n; v++ {
				if val := lps(u, v); val > ans {
					ans = val
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
