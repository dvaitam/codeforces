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
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
		}

		adj := make([][]int, n+1)
		for i := 1; i <= n; i++ {
			j := a[i]
			exists := false
			for _, v := range adj[i] {
				if v == j {
					exists = true
					break
				}
			}
			if !exists {
				adj[i] = append(adj[i], j)
			}
			exists = false
			for _, v := range adj[j] {
				if v == i {
					exists = true
					break
				}
			}
			if !exists {
				adj[j] = append(adj[j], i)
			}
		}

		visited := make([]bool, n+1)
		components := 0
		cycleComp := 0

		for i := 1; i <= n; i++ {
			if visited[i] {
				continue
			}
			components++
			queue := []int{i}
			visited[i] = true
			isCycle := true
			for len(queue) > 0 {
				v := queue[0]
				queue = queue[1:]
				if len(adj[v]) != 2 {
					isCycle = false
				}
				for _, to := range adj[v] {
					if !visited[to] {
						visited[to] = true
						queue = append(queue, to)
					}
				}
			}
			if isCycle {
				cycleComp++
			}
		}

		pathComp := components - cycleComp
		maxCycles := components
		minCycles := cycleComp
		if pathComp > 0 {
			minCycles++
		}
		fmt.Fprintf(out, "%d %d\n", minCycles, maxCycles)
	}
}
