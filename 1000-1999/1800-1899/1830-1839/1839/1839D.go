package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		c := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &c[i])
		}

		incEnd := make([]int, n+2)
		incEnd[n] = n
		for i := n - 1; i >= 1; i-- {
			if c[i] < c[i+1] {
				incEnd[i] = incEnd[i+1]
			} else {
				incEnd[i] = i
			}
		}

		nextGreater := make([]int, n+1)
		nextPos := make([]int, n+2)
		for i := n; i >= 1; i-- {
			best := 0
			for v := c[i] + 1; v <= n; v++ {
				if nextPos[v] != 0 {
					if best == 0 || nextPos[v] < best {
						best = nextPos[v]
					}
				}
			}
			nextGreater[i] = best
			nextPos[c[i]] = i
		}

		dp := make([][]int, n+2)
		for i := 0; i <= n+1; i++ {
			dp[i] = make([]int, n+1)
		}
		for i := n; i >= 1; i-- {
			for k := 1; k <= n; k++ {
				if dp[i+1][k] > dp[i][k] {
					dp[i][k] = dp[i+1][k]
				}
			}
			end := i
			for end <= incEnd[i] {
				length := end - i + 1
				nxt := nextGreater[end]
				for k := 1; k <= n; k++ {
					if nxt == 0 {
						if length > dp[i][k] {
							dp[i][k] = length
						}
					} else if k > 1 {
						val := length + dp[nxt][k-1]
						if val > dp[i][k] {
							dp[i][k] = val
						}
					}
				}
				end++
			}
		}
		for k := 1; k <= n; k++ {
			if k > 1 {
				out.WriteByte(' ')
			}
			fmt.Fprint(out, n-dp[1][k])
		}
		fmt.Fprintln(out)
	}
}
