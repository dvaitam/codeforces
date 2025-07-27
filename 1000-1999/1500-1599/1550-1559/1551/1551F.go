package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		g := make([][]int, n)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			u--
			v--
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
		}

		if k == 2 {
			ans := int64(n) * int64(n-1) / 2 % MOD
			fmt.Fprintln(writer, ans)
			continue
		}

		ans := int64(0)

		for r := 0; r < n; r++ {
			m := len(g[r])
			if m < k {
				continue
			}
			// BFS from r to classify vertices by first step (branch) and distance
			dist := make([]int, n)
			branch := make([]int, n)
			for i := 0; i < n; i++ {
				dist[i] = -1
				branch[i] = -1
			}
			queue := make([]int, 0, n)
			dist[r] = 0
			for idx, v := range g[r] {
				dist[v] = 1
				branch[v] = idx
				queue = append(queue, v)
			}
			for front := 0; front < len(queue); front++ {
				u := queue[front]
				for _, w := range g[u] {
					if w == r {
						continue
					}
					if dist[w] == -1 {
						dist[w] = dist[u] + 1
						branch[w] = branch[u]
						queue = append(queue, w)
					}
				}
			}

			// collect counts: for each branch, for each distance
			maxD := 0
			cnt := make([][]int, m)
			for i := range cnt {
				cnt[i] = make([]int, n+1)
			}
			for v := 0; v < n; v++ {
				if v == r || dist[v] == -1 {
					continue
				}
				b := branch[v]
				d := dist[v]
				cnt[b][d]++
				if d > maxD {
					maxD = d
				}
			}

			// for each distance compute combinations from different branches
			for d := 1; d <= maxD; d++ {
				dp := make([]int64, k+1)
				dp[0] = 1
				for i := 0; i < m; i++ {
					c := cnt[i][d]
					if c == 0 {
						continue
					}
					for j := k; j >= 1; j-- {
						dp[j] = (dp[j] + dp[j-1]*int64(c)) % MOD
					}
				}
				ans = (ans + dp[k]) % MOD
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
