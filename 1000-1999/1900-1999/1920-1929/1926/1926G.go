package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf = int(1e9)

func min(a, b int) int {
	if a < b {
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
		g := make([][]int, n)
		parent := make([]int, n)
		parent[0] = -1
		for i := 1; i < n; i++ {
			var a int
			fmt.Fscan(reader, &a)
			a--
			g[a] = append(g[a], i)
			g[i] = append(g[i], a)
		}
		var s string
		fmt.Fscan(reader, &s)

		order := make([]int, 0, n)
		stack := []int{0}
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			order = append(order, v)
			for _, to := range g[v] {
				if to == parent[v] {
					continue
				}
				parent[to] = v
				stack = append(stack, to)
			}
		}

		dp0 := make([]int, n)
		dp1 := make([]int, n)
		for i := n - 1; i >= 0; i-- {
			v := order[i]
			var val0, val1 int
			switch s[v] {
			case 'S':
				val0 = inf
				val1 = 0
			case 'P':
				val0 = 0
				val1 = inf
			default:
				val0 = 0
				val1 = 0
			}
			for _, to := range g[v] {
				if to == parent[v] {
					continue
				}
				val0 += min(dp0[to], dp1[to]+1)
				val1 += min(dp1[to], dp0[to]+1)
			}
			dp0[v] = val0
			dp1[v] = val1
		}

		ans := min(dp0[0], dp1[0])
		fmt.Fprintln(writer, ans)
	}
}
