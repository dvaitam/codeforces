package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	n       int
	arr     []int64
	adj     map[int64][]int64
	visited map[int64]bool
	res     []int64
)

func dfs(x int64) bool {
	res = append(res, x)
	visited[x] = true
	if len(res) == n {
		return true
	}
	for _, y := range adj[x] {
		if !visited[y] {
			if dfs(y) {
				return true
			}
		}
	}
	res = res[:len(res)-1]
	visited[x] = false
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n)
	arr = make([]int64, n)
	set := make(map[int64]bool, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
		set[arr[i]] = true
	}

	adj = make(map[int64][]int64, n)
	inDeg := make(map[int64]int, n)
	for _, v := range arr {
		if v%3 == 0 && set[v/3] {
			adj[v] = append(adj[v], v/3)
			inDeg[v/3]++
		}
		if set[v*2] {
			adj[v] = append(adj[v], v*2)
			inDeg[v*2]++
		}
	}

	var start int64
	for _, v := range arr {
		if inDeg[v] == 0 {
			start = v
			break
		}
	}

	visited = make(map[int64]bool, n)
	res = make([]int64, 0, n)
	dfs(start)

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for i, v := range res {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	if len(res) > 0 {
		fmt.Fprintln(out)
	}
}
