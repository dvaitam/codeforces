package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf int = int(1e9)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		dp := make([]int, n+1)
		rmin := make([]int, n+1)
		rmax := make([]int, n+1)
		sufCost := make([]int, n+2)
		sufMin := make([]int, n+2)
		sufMax := make([]int, n+2)

		for i := n - 1; i >= 0; i-- {
			keepCost := inf
			keepMin, keepMax := inf, inf
			L := a[i] + 1
			if i+L <= n {
				keepCost = dp[i+L]
				keepMin = rmin[i+L] + 1
				keepMax = rmax[i+L] + 1
			}
			bestCost := sufCost[i+1]
			modCost := 1 + bestCost
			modMin := sufMin[i+1] + 1
			modMax := sufMax[i+1] + 1
			if keepCost < modCost {
				dp[i] = keepCost
				rmin[i] = keepMin
				rmax[i] = keepMax
			} else if keepCost > modCost {
				dp[i] = modCost
				rmin[i] = modMin
				rmax[i] = modMax
			} else {
				dp[i] = keepCost
				if keepMin < modMin {
					rmin[i] = keepMin
				} else {
					rmin[i] = modMin
				}
				if keepMax > modMax {
					rmax[i] = keepMax
				} else {
					rmax[i] = modMax
				}
			}
			if dp[i] < sufCost[i+1] {
				sufCost[i] = dp[i]
				sufMin[i] = rmin[i]
				sufMax[i] = rmax[i]
			} else if dp[i] > sufCost[i+1] {
				sufCost[i] = sufCost[i+1]
				sufMin[i] = sufMin[i+1]
				sufMax[i] = sufMax[i+1]
			} else {
				sufCost[i] = dp[i]
				if rmin[i] < sufMin[i+1] {
					sufMin[i] = rmin[i]
				} else {
					sufMin[i] = sufMin[i+1]
				}
				if rmax[i] > sufMax[i+1] {
					sufMax[i] = rmax[i]
				} else {
					sufMax[i] = sufMax[i+1]
				}
			}
		}

		for i := 0; i < n-1; i++ {
			cost := dp[i+1]
			if rmin[i+1] <= a[i] && a[i] <= rmax[i+1] {
				if i > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, cost)
			} else {
				if i > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, cost+1)
			}
		}
		fmt.Fprintln(out)
	}
}
