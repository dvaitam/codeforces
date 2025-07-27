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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		g := make([][]int, n)
		for i := 0; i < m; i++ {
			var x, y int
			fmt.Fscan(in, &x, &y)
			x--
			y--
			g[x] = append(g[x], y)
		}

		dp := make([]int, n)
		removed := make([]bool, n)
		ans := make([]int, 0)
		for i := 0; i < n; i++ {
			if dp[i] >= 2 {
				removed[i] = true
				ans = append(ans, i+1)
			} else {
				for _, v := range g[i] {
					if dp[i]+1 > dp[v] {
						dp[v] = dp[i] + 1
					}
				}
			}
		}
		fmt.Fprintln(out, len(ans))
		for i, v := range ans {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		if len(ans) > 0 {
			fmt.Fprintln(out)
		} else {
			fmt.Fprintln(out)
		}
	}
}
