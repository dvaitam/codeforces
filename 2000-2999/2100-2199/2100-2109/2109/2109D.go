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
		minOdd := -1
		for i := 0; i < l; i++ {
			var x int
			fmt.Fscan(in, &x)
			totalSum += int64(x)
			if x%2 == 1 {
				if minOdd == -1 || x < minOdd {
					minOdd = x
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

		dist := make([]int, n)
		for i := range dist {
			dist[i] = -1
		}
		color := make([]int8, n)
		for i := range color {
			color[i] = -1
		}

		queue := make([]int, 0, n)
		dist[0] = 0
		color[0] = 0
		queue = append(queue, 0)
		isBipartite := true
		for head := 0; head < len(queue); head++ {
			v := queue[head]
			for _, to := range adj[v] {
				if dist[to] == -1 {
					dist[to] = dist[v] + 1
					color[to] = 1 - color[v]
					queue = append(queue, to)
				} else if color[to] == color[v] {
					isBipartite = false
				}
			}
		}

		maxEven, maxOdd := int64(-1), int64(-1)
		if isBipartite {
			if totalSum%2 == 0 {
				maxEven = totalSum
				if minOdd != -1 {
					maxOdd = totalSum - int64(minOdd)
				}
			} else {
				maxOdd = totalSum
				if minOdd != -1 {
					maxEven = totalSum - int64(minOdd)
				}
			}
		}

		ans := make([]byte, n)
		for i := 0; i < n; i++ {
			d := int64(dist[i])
			reachable := false
			if dist[i] == -1 {
				reachable = false
			} else if i == 0 {
				reachable = true
			} else if !isBipartite {
				reachable = d <= totalSum
			} else {
				if d%2 == 0 {
					reachable = maxEven >= d && maxEven != -1
				} else {
					reachable = maxOdd >= d && maxOdd != -1
				}
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
