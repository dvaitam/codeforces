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
		var n, k int
		fmt.Fscan(in, &n, &k)
		adj := make([][]int, n+1)
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
			edges[i] = [2]int{u, v}
		}
		if k == 1 {
			fmt.Fprintln(out, n-1)
			continue
		}

		deg := make([]int, n+1)
		for i := 1; i <= n; i++ {
			deg[i] = len(adj[i])
		}

		leafVec := make([][]int, n+1)
		for i := 1; i <= n; i++ {
			if deg[i] == 1 {
				p := adj[i][0]
				leafVec[p] = append(leafVec[p], i)
			}
		}

		leafCnt := make([]int, n+1)
		for i := 1; i <= n; i++ {
			leafCnt[i] = len(leafVec[i])
		}

		queue := make([]int, 0)
		for i := 1; i <= n; i++ {
			if leafCnt[i] >= k {
				queue = append(queue, i)
			}
		}
		head := 0
		ans := 0
		for head < len(queue) {
			v := queue[head]
			head++
			if leafCnt[v] < k {
				continue
			}
			times := leafCnt[v] / k
			ans += times
			for i := 0; i < times*k; i++ {
				u := leafVec[v][len(leafVec[v])-1]
				leafVec[v] = leafVec[v][:len(leafVec[v])-1]
				deg[u] = 0
				leafCnt[v]--
			}
			deg[v] -= times * k
			if deg[v] == 1 {
				var parent int
				for _, to := range adj[v] {
					if deg[to] > 0 {
						parent = to
						break
					}
				}
				if parent > 0 {
					leafVec[parent] = append(leafVec[parent], v)
					leafCnt[parent]++
					if leafCnt[parent] >= k {
						queue = append(queue, parent)
					}
				}
			}
			if leafCnt[v] >= k {
				queue = append(queue, v)
			}
		}
		fmt.Fprintln(out, ans)
	}
}
