package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, q int
	fmt.Fscan(reader, &n, &m, &q)

	adj := make([][]int, n+1)
	rev := make([][]int, n+1)
	for i := 0; i < m; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		adj[x] = append(adj[x], y)
		rev[y] = append(rev[y], x)
	}
	for i := 1; i <= n; i++ {
		sort.Ints(adj[i])
	}

	// Precompute reachability from every node to every target using reverse BFS
	reachable := make([][]bool, n+1)
	for t := 1; t <= n; t++ {
		reachable[t] = make([]bool, n+1)
		queue := []int{t}
		reachable[t][t] = true
		for head := 0; head < len(queue); head++ {
			v := queue[head]
			for _, u := range rev[v] {
				if !reachable[t][u] {
					reachable[t][u] = true
					queue = append(queue, u)
				}
			}
		}
	}

	type pair struct{ s, t int }
	cache := make(map[pair][]int)

	for ; q > 0; q-- {
		var s, t, k int
		fmt.Fscan(reader, &s, &t, &k)
		key := pair{s, t}
		path, ok := cache[key]
		if !ok {
			if !reachable[t][s] {
				cache[key] = nil
			} else {
				visited := make([]bool, n+1)
				curr := s
				visited[curr] = true
				res := []int{curr}
				valid := true
				for curr != t {
					nextFound := false
					for _, v := range adj[curr] {
						if reachable[t][v] {
							curr = v
							nextFound = true
							if visited[curr] {
								valid = false
								break
							}
							visited[curr] = true
							res = append(res, curr)
							break
						}
					}
					if !nextFound || !valid {
						valid = false
						break
					}
				}
				if valid {
					path = res
					cache[key] = path
				} else {
					cache[key] = nil
					path = nil
				}
			}
		}
		if path == nil || k <= 0 || k > len(path) {
			fmt.Fprintln(writer, -1)
		} else {
			fmt.Fprintln(writer, path[k-1])
		}
	}
}
