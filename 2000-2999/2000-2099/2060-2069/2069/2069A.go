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
		m := n - 2
		b := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &b[i])
		}

		xLen := n - 1
		dp := make([][2]bool, xLen)
		dp[0][0], dp[0][1] = true, true

		for i := 0; i < xLen-1; i++ {
			bVal := b[i]
			for cur := 0; cur <= 1; cur++ {
				if !dp[i][cur] {
					continue
				}
				for next := 0; next <= 1; next++ {
					if (cur & next) == bVal {
						dp[i+1][next] = true
					}
				}
			}
		}

		if dp[xLen-1][0] || dp[xLen-1][1] {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
