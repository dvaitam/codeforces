package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		g := make([][]int, n)
		for i := 0; i < n-1; i++ {
			var x, y int
			fmt.Fscan(in, &x, &y)
			x--
			y--
			g[x] = append(g[x], y)
			g[y] = append(g[y], x)
		}

		if n == 1 {
			fmt.Fprintln(out, 0)
			continue
		}

		deg := make([]int, n)
		hasMajor := false
		for i := 0; i < n; i++ {
			deg[i] = len(g[i])
			if deg[i] >= 3 {
				hasMajor = true
			}
		}
		if !hasMajor {
			fmt.Fprintln(out, 1)
			continue
		}

		counts := make([]int, n)
		for i := 0; i < n; i++ {
			if deg[i] == 1 { // leaf
				curr := i
				prev := -1
				for deg[curr] < 3 {
					next := -1
					for _, v := range g[curr] {
						if v != prev {
							next = v
							break
						}
					}
					prev = curr
					curr = next
				}
				counts[curr]++
			}
		}

		ans := 0
		for i := 0; i < n; i++ {
			if counts[i] > 0 {
				ans += counts[i] - 1
			}
		}
		fmt.Fprintln(out, ans)
	}
}
