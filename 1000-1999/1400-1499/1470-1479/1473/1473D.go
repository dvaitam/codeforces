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
		var n, m int
		fmt.Fscan(in, &n, &m)
		var s string
		fmt.Fscan(in, &s)

		pref := make([]int, n+1)
		prefMin := make([]int, n+1)
		prefMax := make([]int, n+1)
		for i := 1; i <= n; i++ {
			if s[i-1] == '+' {
				pref[i] = pref[i-1] + 1
			} else {
				pref[i] = pref[i-1] - 1
			}
			prefMin[i] = min(prefMin[i-1], pref[i])
			prefMax[i] = max(prefMax[i-1], pref[i])
		}

		sufMin := make([]int, n+1)
		sufMax := make([]int, n+1)
		sufMin[n] = pref[n]
		sufMax[n] = pref[n]
		for i := n - 1; i >= 0; i-- {
			if pref[i] < sufMin[i+1] {
				sufMin[i] = pref[i]
			} else {
				sufMin[i] = sufMin[i+1]
			}
			if pref[i] > sufMax[i+1] {
				sufMax[i] = pref[i]
			} else {
				sufMax[i] = sufMax[i+1]
			}
		}

		for i := 0; i < m; i++ {
			var l, r int
			fmt.Fscan(in, &l, &r)
			leftMin := prefMin[l-1]
			leftMax := prefMax[l-1]

			rightMin := pref[l-1] + (sufMin[r] - pref[r])
			rightMax := pref[l-1] + (sufMax[r] - pref[r])

			overallMin := min(leftMin, rightMin)
			overallMax := max(leftMax, rightMax)
			fmt.Fprintln(out, overallMax-overallMin+1)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
