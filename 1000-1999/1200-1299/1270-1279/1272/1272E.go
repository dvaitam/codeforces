package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		if i-a[i] >= 0 {
			adj[i-a[i]] = append(adj[i-a[i]], i)
		}
		if i+a[i] < n {
			adj[i+a[i]] = append(adj[i+a[i]], i)
		}
	}
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	q := make([]int, 0)
	for i := 0; i < n; i++ {
		found := false
		if i-a[i] >= 0 && (a[i-a[i]]%2 != a[i]%2) {
			found = true
		}
		if !found && i+a[i] < n && (a[i+a[i]]%2 != a[i]%2) {
			found = true
		}
		if found {
			dist[i] = 1
			q = append(q, i)
		}
	}
	for head := 0; head < len(q); head++ {
		v := q[head]
		for _, to := range adj[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				q = append(q, to)
			}
		}
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(dist[i])
	}
	fmt.Println()
}
