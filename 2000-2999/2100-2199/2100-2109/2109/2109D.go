package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m, l int
		fmt.Fscan(in, &n, &m, &l)

		totalSum := int64(0)
		minOdd := int64(-1)
		for i := 0; i < l; i++ {
			var x int
			fmt.Fscan(in, &x)
			totalSum += int64(x)
			if x%2 == 1 {
				if minOdd == -1 || int64(x) < minOdd {
					minOdd = int64(x)
				}
			}
		}

		adj := make([][]int, n)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		const INF = int64(1e18)
		dist := make([][2]int64, n)
		for i := range dist {
			dist[i] = [2]int64{INF, INF}
		}
		dist[0][0] = 0

		type state struct {
			v, p int
		}
		queue := []state{{0, 0}}
		for head := 0; head < len(queue); head++ {
			s := queue[head]
			c := dist[s.v][s.p]
			np := 1 - s.p
			for _, to := range adj[s.v] {
				if dist[to][np] <= c+1 {
					continue
				}
				dist[to][np] = c + 1
				queue = append(queue, state{to, np})
			}
		}

		ans := make([]byte, n)
		for i := 0; i < n; i++ {
			reachable := false
			sp := int(totalSum % 2)
			if dist[i][sp] <= totalSum {
				reachable = true
			}
			if minOdd != -1 && dist[i][1-sp] <= totalSum-minOdd {
				reachable = true
			}
			if reachable {
				ans[i] = '1'
			} else {
				ans[i] = '0'
			}
		}
		fmt.Fprintln(out, string(ans))
	}
}
