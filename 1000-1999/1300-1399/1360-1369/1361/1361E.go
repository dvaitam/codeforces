package main

import (
	"bufio"
	"fmt"
	"os"
)

func isAncestor(tin, tout []int, u, v int) bool {
	return tin[v] <= tin[u] && tout[u] <= tout[v]
}

func checkInteresting(n int, g [][]int, root int) bool {
	visited := make([]bool, n)
	parent := make([]int, n)
	tin := make([]int, n)
	tout := make([]int, n)
	for i := range parent {
		parent[i] = -2
	}
	time := 0
	type item struct{ v, idx int }
	stack := []item{{root, 0}}
	parent[root] = -1
	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		v := top.v
		if top.idx == 0 {
			visited[v] = true
			tin[v] = time
			time++
		}
		if top.idx < len(g[v]) {
			to := g[v][top.idx]
			top.idx++
			if !visited[to] {
				parent[to] = v
				stack = append(stack, item{to, 0})
			}
			continue
		}
		tout[v] = time
		time++
		stack = stack[:len(stack)-1]
	}
	for _, ok := range visited {
		if !ok {
			return false
		}
	}
	for u := 0; u < n; u++ {
		for _, v := range g[u] {
			if parent[v] == u {
				continue
			}
			if !isAncestor(tin, tout, u, v) {
				return false
			}
		}
	}
	return true
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
		var n, m int
		fmt.Fscan(reader, &n, &m)
		g := make([][]int, n)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			u--
			v--
			g[u] = append(g[u], v)
		}
		var res []int
		for r := 0; r < n; r++ {
			if checkInteresting(n, g, r) {
				res = append(res, r+1)
			}
		}
		if len(res)*5 < n {
			fmt.Fprintln(writer, -1)
		} else {
			for i, v := range res {
				if i > 0 {
					writer.WriteByte(' ')
				}
				fmt.Fprint(writer, v)
			}
			fmt.Fprintln(writer)
		}
	}
}
