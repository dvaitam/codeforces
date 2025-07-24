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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		adj := make([][]int, n+1)
		deg := make([]int, n+1)
		ok := true
		for i := 0; i < n; i++ {
			var a, b int
			fmt.Fscan(in, &a, &b)
			if a == b {
				ok = false
			}
			deg[a]++
			deg[b]++
			adj[a] = append(adj[a], b)
			adj[b] = append(adj[b], a)
		}
		if ok {
			for i := 1; i <= n; i++ {
				if deg[i] > 2 {
					ok = false
					break
				}
			}
		}
		if ok {
			vis := make([]bool, n+1)
			for i := 1; i <= n && ok; i++ {
				if vis[i] || deg[i] == 0 {
					continue
				}
				stack := []int{i}
				vis[i] = true
				nodes, edges := 0, 0
				for len(stack) > 0 {
					v := stack[len(stack)-1]
					stack = stack[:len(stack)-1]
					nodes++
					edges += len(adj[v])
					for _, to := range adj[v] {
						if !vis[to] {
							vis[to] = true
							stack = append(stack, to)
						}
					}
				}
				edges /= 2
				if edges == nodes && nodes%2 == 1 {
					ok = false
					break
				}
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
