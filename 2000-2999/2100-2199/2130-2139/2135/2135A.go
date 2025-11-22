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
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n+1) // 1-indexed
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
		}

		pos := make([][]int, n+1)
		occIdx := make([]int, n+1)
		for i := 1; i <= n; i++ {
			v := a[i]
			pos[v] = append(pos[v], i)
			occIdx[i] = len(pos[v]) - 1
		}

		dp := make([]int, n+1) // dp[i] = best length using elements up to i
		for i := 1; i <= n; i++ {
			dp[i] = dp[i-1]
			v := a[i]
			idx := occIdx[i] // which occurrence of v is at position i
			if idx+1 >= v {  // enough occurrences to form a block ending here
				startPos := pos[v][idx-v+1] // first occurrence in this block
				cand := dp[startPos-1] + v  // previous best + block size
				if cand > dp[i] {           // maximize
					dp[i] = cand
				}
			}
		}

		fmt.Fprintln(out, dp[n])
	}
}
