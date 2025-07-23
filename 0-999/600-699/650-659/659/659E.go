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

	var n, m int
	fmt.Fscan(in, &n, &m)
	adj := make([][]int, n+1)
	for i := 0; i < m; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		adj[x] = append(adj[x], y)
		adj[y] = append(adj[y], x)
	}

	vis := make([]bool, n+1)
	ans := 0
	stack := make([]int, 0)
	for i := 1; i <= n; i++ {
		if !vis[i] {
			countV := 0
			countE := 0
			stack = append(stack[:0], i)
			vis[i] = true
			for len(stack) > 0 {
				v := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				countV++
				countE += len(adj[v])
				for _, to := range adj[v] {
					if !vis[to] {
						vis[to] = true
						stack = append(stack, to)
					}
				}
			}
			if countE/2 == countV-1 {
				ans++
			}
		}
	}
	fmt.Fprintln(out, ans)
}
