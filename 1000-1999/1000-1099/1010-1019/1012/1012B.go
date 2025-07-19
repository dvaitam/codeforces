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

	var n, m, q int
	fmt.Fscan(reader, &n, &m, &q)
	// Union-Find over n rows and m columns (columns offset by n)
	parent := make([]int, n+m+1)
	for i := 1; i <= n+m; i++ {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		for parent[x] != x {
			parent[x] = parent[parent[x]]
			x = parent[x]
		}
		return x
	}
	union := func(x, y int) bool {
		x = find(x)
		y = find(y)
		if x == y {
			return false
		}
		parent[x] = y
		return true
	}

	ans := n + m - 1
	for i := 0; i < q; i++ {
		var r, c int
		fmt.Fscan(reader, &r, &c)
		if union(r, c+n) {
			ans--
		}
	}
	fmt.Fprintln(writer, ans)
}
