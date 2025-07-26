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

	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	n := len(s)
	ones := 0
	for i := 0; i < n; i++ {
		if s[i] == '1' {
			ones++
		}
	}
	m := ones
	const maxN = 100
	maxSum := maxN * maxN
	shift := maxSum
	inf := 1 << 30

	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, 2*maxSum+1)
		for j := range dp[i] {
			dp[i][j] = inf
		}
	}
	dp[0][shift] = 0

	rangePrev := 0
	for i := 0; i < n; i++ {
		w := 2*(i+1) - n - 1
		newdp := make([][]int, m+1)
		for k := range newdp {
			newdp[k] = make([]int, 2*maxSum+1)
			for j := range newdp[k] {
				newdp[k][j] = inf
			}
		}
		rangeCurr := rangePrev + 99
		if rangeCurr > maxSum {
			rangeCurr = maxSum
		}
		for k := 0; k <= m; k++ {
			start := shift - rangePrev
			end := shift + rangePrev
			if start < 0 {
				start = 0
			}
			if end > 2*maxSum {
				end = 2 * maxSum
			}
			for b := start; b <= end; b++ {
				val := dp[k][b]
				if val == inf {
					continue
				}
				// place 0 at position i
				cost0 := 0
				if s[i] == '1' {
					cost0 = 1
				}
				if val+cost0 < newdp[k][b] {
					newdp[k][b] = val + cost0
				}
				// place 1 at position i
				if k < m {
					nb := b + w
					if nb >= 0 && nb <= 2*maxSum {
						cost1 := 0
						if s[i] == '0' {
							cost1 = 1
						}
						if val+cost1 < newdp[k+1][nb] {
							newdp[k+1][nb] = val + cost1
						}
					}
				}
			}
		}
		dp = newdp
		rangePrev = rangeCurr
	}
	mismatches := dp[m][shift]
	fmt.Fprintln(writer, mismatches/2)
}
