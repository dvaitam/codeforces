package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	if n == 1 {
		fmt.Println("Yes")
		fmt.Println(0)
		return
	}

	deg := make([]int, n+1)
	for i := 1; i <= n; i++ {
		deg[i] = len(adj[i])
	}

	hub := -1
	for i := 1; i <= n; i++ {
		if deg[i] > 2 {
			if hub == -1 {
				hub = i
			} else {
				fmt.Println("No")
				return
			}
		}
	}

	fmt.Println("Yes")
	if hub == -1 {
		// tree is a simple path
		var ends []int
		for i := 1; i <= n; i++ {
			if deg[i] == 1 {
				ends = append(ends, i)
			}
		}
		fmt.Println(1)
		fmt.Printf("%d %d\n", ends[0], ends[1])
		return
	}

	fmt.Println(deg[hub])
	for i := 1; i <= n; i++ {
		if deg[i] == 1 {
			fmt.Printf("%d %d\n", hub, i)
		}
	}
}
