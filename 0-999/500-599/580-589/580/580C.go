package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	cats := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &cats[i])
	}
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	type node struct {
		v   int
		cnt int
	}
	stack := []node{{1, cats[1]}}
	visited := make([]bool, n+1)
	visited[1] = true
	ans := 0
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if cur.cnt > m {
			continue
		}
		isLeaf := true
		for _, to := range adj[cur.v] {
			if !visited[to] {
				visited[to] = true
				nextCnt := 0
				if cats[to] == 1 {
					nextCnt = cur.cnt + 1
				}
				stack = append(stack, node{to, nextCnt})
				isLeaf = false
			}
		}
		if isLeaf {
			ans++
		}
	}
	fmt.Println(ans)
}
