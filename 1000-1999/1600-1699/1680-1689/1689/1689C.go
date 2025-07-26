package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	g    [][]int
	size []int
	dp   []int
)

func solve(n int, edges [][2]int, writer *bufio.Writer) {
	g = make([][]int, n)
	for _, e := range edges {
		u := e[0]
		v := e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	size = make([]int, n)
	dp = make([]int, n)
	parent := make([]int, n)
	order := make([]int, 0, n)
	stack := []int{0}
	parent[0] = -1
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, u)
		for _, v := range g[u] {
			if v == parent[u] {
				continue
			}
			parent[v] = u
			stack = append(stack, v)
		}
	}
	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		size[u] = 1
		children := make([]int, 0, 2)
		for _, v := range g[u] {
			if v == parent[u] {
				continue
			}
			size[u] += size[v]
			children = append(children, v)
		}
		switch len(children) {
		case 0:
			dp[u] = 0
		case 1:
			c := children[0]
			opt1 := dp[c]
			opt2 := size[c] - 1
			if opt2 > opt1 {
				dp[u] = opt2
			} else {
				dp[u] = opt1
			}
		default: // two children
			a := children[0]
			b := children[1]
			best := dp[a] + dp[b]
			if tmp := dp[a] + size[b] - 1; tmp > best {
				best = tmp
			}
			if tmp := dp[b] + size[a] - 1; tmp > best {
				best = tmp
			}
			dp[u] = best
		}
	}
	fmt.Fprintln(writer, dp[0])
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			u--
			v--
			edges[i] = [2]int{u, v}
		}
		solve(n, edges, writer)
	}
}
