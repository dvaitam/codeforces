package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var x, y int64
	if _, err := fmt.Fscan(reader, &n, &x, &y); err != nil {
		return
	}

	adj := make([][]int, n)
	deg := make([]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		u--
		v--
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		deg[u]++
		deg[v]++
	}

	if x >= y {
		star := false
		for _, d := range deg {
			if d == n-1 {
				star = true
				break
			}
		}
		if star {
			ans := int64(n-2)*y + x
			fmt.Fprintln(writer, ans)
		} else {
			ans := int64(n-1) * y
			fmt.Fprintln(writer, ans)
		}
		return
	}

	parent := make([]int, n)
	for i := range parent {
		parent[i] = -1
	}
	order := make([]int, 0, n)
	stack := []int{0}
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		for _, w := range adj[v] {
			if w == parent[v] {
				continue
			}
			parent[w] = v
			stack = append(stack, w)
		}
	}

	dp0 := make([]int, n)
	dp1 := make([]int, n)
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		base0 := 0
		best1, best2 := 0, 0
		for _, u := range adj[v] {
			if u == parent[v] {
				continue
			}
			base0 += dp0[u]
			diff := dp1[u] - dp0[u]
			if diff > best1 {
				best2 = best1
				best1 = diff
			} else if diff > best2 {
				best2 = diff
			}
		}
		if best1 < 0 {
			best1 = 0
		}
		if best2 < 0 {
			best2 = 0
		}
		dp0[v] = base0 + best1 + best2
		if parent[v] != -1 {
			if best1 < 0 {
				best1 = 0
			}
			dp1[v] = 1 + base0 + best1
		}
	}
	s := int64(dp0[0])
	ans := s*x + int64(n-1-int(s))*y
	fmt.Fprintln(writer, ans)
}
