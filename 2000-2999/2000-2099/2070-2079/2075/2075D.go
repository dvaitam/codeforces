package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	maxShift = 60
	inf      = int64(4e18)
)

// precompute minimal cost to achieve exact shifts (a, b) with disjoint k choices
func buildDP() [][]int64 {
	dp := make([][]int64, maxShift+1)
	for i := range dp {
		dp[i] = make([]int64, maxShift+1)
		for j := range dp[i] {
			dp[i][j] = inf
		}
	}
	dp[0][0] = 0

	for k := 1; k <= maxShift; k++ {
		cost := int64(1) << k
		for a := maxShift; a >= 0; a-- {
			for b := maxShift; b >= 0; b-- {
				cur := dp[a][b]
				if cur == inf {
					continue
				}
				if a+k <= maxShift && dp[a+k][b] > cur+cost {
					dp[a+k][b] = cur + cost
				}
				if b+k <= maxShift && dp[a][b+k] > cur+cost {
					dp[a][b+k] = cur + cost
				}
			}
		}
	}
	return dp
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	dp := buildDP()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var x, y uint64
		fmt.Fscan(in, &x, &y)

		ans := inf
		for a := 0; a <= maxShift; a++ {
			vx := x >> a
			for b := 0; b <= maxShift; b++ {
				if vx == (y >> b) {
					if dp[a][b] < ans {
						ans = dp[a][b]
					}
				}
			}
		}

		fmt.Fprintln(out, ans)
	}
}
