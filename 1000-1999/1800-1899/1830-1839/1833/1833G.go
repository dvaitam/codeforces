package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct{ u, v int }

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		adj := make([][]int, n+1)
		deg := make([]int, n+1)
		size := make([]int, n+1)
		visited := make([]bool, n+1)
		mp2 := make(map[pair]int)

		for i := 1; i <= n; i++ {
			size[i] = 1
		}
		// read edges
		for i := 1; i < n; i++ {
			var x, y int
			fmt.Fscan(reader, &x, &y)
			adj[x] = append(adj[x], y)
			adj[y] = append(adj[y], x)
			deg[x]++
			deg[y]++
			u, v := min(x, y), max(x, y)
			mp2[pair{u, v}] = i
		}
		if n%3 != 0 {
			fmt.Fprintln(writer, -1)
			continue
		}
		// initialize queue with leaves
		queue := make([]int, 0, n)
		head := 0
		for i := 1; i <= n; i++ {
			if deg[i] == 1 {
				queue = append(queue, i)
			}
		}
		st := make(map[pair]bool)
		check := false
		// process
		for head < len(queue) {
			it := queue[head]
			head++
			if size[it] > 3 {
				check = true
				break
			}
			visited[it] = true
			for _, child := range adj[it] {
				if !visited[child] {
					deg[child]--
					if size[it] == 3 {
						u, v := min(it, child), max(it, child)
						st[pair{u, v}] = true
					} else {
						size[child] += size[it]
					}
					if deg[child] == 1 {
						queue = append(queue, child)
					}
				}
			}
		}
		if check {
			fmt.Fprintln(writer, -1)
			continue
		}
		// output
		count := len(st)
		fmt.Fprintln(writer, count)
		if count == 0 {
			fmt.Fprintln(writer)
			continue
		}
		// collect and sort pairs
		keys := make([]pair, 0, count)
		for k := range st {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			if keys[i].u != keys[j].u {
				return keys[i].u < keys[j].u
			}
			return keys[i].v < keys[j].v
		})
		// print edge indices
		for i, k := range keys {
			if i > 0 {
				writer.WriteByte(' ')
			}
			writer.WriteString(fmt.Sprintf("%d", mp2[k]))
		}
		fmt.Fprintln(writer)
	}
}
