package main

import (
	"bufio"
	"fmt"
	"os"
)

func canDraw(n int, edges [][2]int) bool {
	adj := make([][]int, n)
	for _, e := range edges {
		a, b := e[0], e[1]
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	if n <= 2 {
		return true
	}
	deg := make([]int, n)
	for i := 0; i < n; i++ {
		deg[i] = len(adj[i])
	}
	isLeaf := make([]bool, n)
	for i := 0; i < n; i++ {
		if deg[i] == 1 {
			isLeaf[i] = true
		}
	}
	cnt := make([]int, n)
	for v := 0; v < n; v++ {
		for _, u := range adj[v] {
			if !isLeaf[u] {
				cnt[v]++
			}
		}
	}
	var spine []int
	for i := 0; i < n; i++ {
		if !isLeaf[i] {
			if cnt[i] > 2 {
				return false
			}
			spine = append(spine, i)
		}
	}
	if len(spine) == 0 || len(spine) == 1 {
		return true
	}
	ends := 0
	for _, v := range spine {
		if cnt[v] == 1 {
			ends++
		}
	}
	if ends != 2 {
		return false
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	edges := make([][2]int, n-1)
	for i := 0; i < n-1; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		edges[i] = [2]int{a, b}
	}
	if canDraw(n, edges) {
		fmt.Println("Yes")
	} else {
		fmt.Println("No")
	}
}
