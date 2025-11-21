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
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)
		missing := make([]int, m)
		for i := range missing {
			fmt.Fscan(in, &missing[i])
		}
		known := make([]bool, n+1)
		for i := 0; i < k; i++ {
			var q int
			fmt.Fscan(in, &q)
			known[q] = true
		}
		ans := make([]byte, m)
		knownCount := 0
		for i := 1; i <= n; i++ {
			if known[i] {
				knownCount++
			}
		}
		for i := 0; i < m; i++ {
			miss := missing[i]
			totalUnknown := n - knownCount
			if totalUnknown == 0 || (totalUnknown == 1 && !known[miss]) {
				if known[miss] {
					ans[i] = '1'
				} else {
					ans[i] = '1'
				}
			} else {
				ans[i] = '0'
			}
		}
		fmt.Fprintln(out, string(ans))
	}
}
