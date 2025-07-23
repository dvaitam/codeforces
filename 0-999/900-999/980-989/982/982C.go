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
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	if n%2 == 1 {
		fmt.Println(-1)
		return
	}
	parent := make([]int, n+1)
	order := make([]int, 0, n)
	stack := []int{1}
	parent[1] = 0
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		for _, to := range adj[v] {
			if to != parent[v] {
				parent[to] = v
				stack = append(stack, to)
			}
		}
	}
	size := make([]int, n+1)
	ans := 0
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		size[v] = 1
		for _, to := range adj[v] {
			if to != parent[v] {
				size[v] += size[to]
			}
		}
		if v != 1 && size[v]%2 == 0 {
			ans++
		}
	}
	fmt.Println(ans)
}
