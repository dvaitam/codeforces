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

	var q int
	if _, err := fmt.Fscan(in, &q); err != nil {
		return
	}
	for ; q > 0; q-- {
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
		w := make([]int, n)
		for i := 0; i < n; i++ {
			w[i] = len(g[i]) - 1
		}
		down := make([]int, n)
		ans := 0
		var dfs func(v, p int)
		dfs = func(v, p int) {
			max1, max2 := 0, 0
			for _, u := range g[v] {
				if u == p {
					continue
				}
				dfs(u, v)
				val := down[u]
				if val > max1 {
					max2 = max1
					max1 = val
				} else if val > max2 {
					max2 = val
				}
			}
			if cur := w[v] + max1 + max2; cur > ans {
				ans = cur
			}
			down[v] = w[v] + max1
		}
		dfs(0, -1)
		fmt.Fprintln(out, ans+2)
	}
}
