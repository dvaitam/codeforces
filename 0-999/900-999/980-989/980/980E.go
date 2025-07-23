package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}

	// root the tree at node n
	parent := make([]int, n+1)
	queue := make([]int, 0, n)
	queue = append(queue, n)
	parent[n] = 0
	for head := 0; head < len(queue); head++ {
		v := queue[head]
		for _, u := range adj[v] {
			if u == parent[v] {
				continue
			}
			parent[u] = v
			queue = append(queue, u)
		}
	}

	need := n - k
	good := make([]bool, n+1)
	good[n] = true
	count := 1
	for i := n - 1; i >= 1 && count < need; i-- {
		if good[i] {
			continue
		}
		path := []int{}
		x := i
		for !good[x] {
			path = append(path, x)
			x = parent[x]
		}
		if count+len(path) <= need {
			for _, v := range path {
				good[v] = true
			}
			count += len(path)
		}
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	res := make([]string, 0, k)
	for i := 1; i <= n; i++ {
		if !good[i] {
			res = append(res, strconv.Itoa(i))
		}
	}
	fmt.Fprintln(out, strings.Join(res, " "))
}
