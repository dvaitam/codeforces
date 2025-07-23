package main

import (
	"bufio"
	"fmt"
	"os"
)

const MaxA = 5000

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}

	first := make([]int, MaxA+1)
	last := make([]int, MaxA+1)
	for i := 0; i <= MaxA; i++ {
		first[i] = n + 1
		last[i] = 0
	}
	for i := 1; i <= n; i++ {
		v := a[i]
		if first[v] == n+1 {
			first[v] = i
		}
		last[v] = i
	}

	dp := make([]int, n+1)
	visited := make([]bool, MaxA+1)

	for l := 1; l <= n; l++ {
		for i := 0; i <= MaxA; i++ {
			visited[i] = false
		}
		xorVal := 0
		r := l - 1
		for j := l; j <= n; j++ {
			v := a[j]
			if first[v] < l {
				break
			}
			if !visited[v] {
				visited[v] = true
				xorVal ^= v
			}
			if last[v] > r {
				r = last[v]
			}
			if j == r {
				if cand := dp[l-1] + xorVal; cand > dp[r] {
					dp[r] = cand
				}
			}
		}
	}
	for i := 1; i <= n; i++ {
		if dp[i] < dp[i-1] {
			dp[i] = dp[i-1]
		}
	}
	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, dp[n])
	out.Flush()
}
