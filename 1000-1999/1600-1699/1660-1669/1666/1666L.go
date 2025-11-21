package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m, s int
	if _, err := fmt.Fscan(in, &n, &m, &s); err != nil {
		return
	}

	adj := make([][]int, n+1)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
	}

	root := make([]int, n+1)
	parent := make([]int, n+1)
	queue := make([]int, 0)

	for _, v := range adj[s] {
		if root[v] == 0 {
			root[v] = v
			parent[v] = s
			queue = append(queue, v)
		}
	}

	var meetNode int
	var meetFrom int

	for head := 0; head < len(queue) && meetNode == 0; head++ {
		u := queue[head]
		for _, v := range adj[u] {
			if v == s {
				continue
			}
			if root[v] == 0 {
				root[v] = root[u]
				parent[v] = u
				queue = append(queue, v)
			} else if root[v] != root[u] {
				meetNode = v
				meetFrom = u
				break
			}
		}
	}

	if meetNode == 0 {
		fmt.Println("Impossible")
		return
	}

	path1 := buildPath(meetNode, parent, s)
	path2 := buildPath(meetFrom, parent, s)
	path2 = append(path2, meetNode)

	fmt.Println("Possible")
	fmt.Println(len(path1))
	printPath(path1)
	fmt.Println(len(path2))
	printPath(path2)
}

func buildPath(to int, parent []int, s int) []int {
	path := make([]int, 0)
	for v := to; ; v = parent[v] {
		path = append(path, v)
		if v == s {
			break
		}
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func printPath(path []int) {
	for i, v := range path {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(v)
	}
	fmt.Println()
}
