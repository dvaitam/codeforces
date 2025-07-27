package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		zeros := make([]int, n+m-1)
		ones := make([]int, n+m-1)
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				var x int
				fmt.Fscan(in, &x)
				if x == 0 {
					zeros[i+j]++
				} else {
					ones[i+j]++
				}
			}
		}
		total := n + m - 2
		ans := 0
		for l, r := 0, total; l < r; l, r = l+1, r-1 {
			ans += min(zeros[l]+zeros[r], ones[l]+ones[r])
		}
		fmt.Fprintln(out, ans)
	}
}
