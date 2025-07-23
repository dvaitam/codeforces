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

	var n, c int
	if _, err := fmt.Fscan(reader, &n, &c); err != nil {
		return
	}
	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	if c == 1 {
		fmt.Fprintln(writer, 0)
		return
	}
	pref := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1] + a[i]
	}
	type pair struct {
		idx int
		val int64
	}
	dp := make([]int64, n+1)
	dq := []pair{{0, 0}} // monotonic deque for dp[j]-pref[j]
	dqFront := 0
	minD := []pair{}
	minFront := 0
	for i := 1; i <= n; i++ {
		// maintain deque of minimum for last c elements
		for len(minD) > minFront && minD[len(minD)-1].val >= a[i] {
			minD = minD[:len(minD)-1]
		}
		minD = append(minD, pair{i, a[i]})
		for len(minD) > minFront && minD[minFront].idx <= i-c {
			minFront++
		}
		// maintain deque for range of dp values
		for len(dq) > dqFront && dq[dqFront].idx < i-c+1 {
			dqFront++
		}
		best := dq[dqFront].val
		dp[i] = pref[i] + best
		if i >= c {
			minVal := minD[minFront].val
			val := dp[i-c] + pref[i] - pref[i-c] - minVal
			if val < dp[i] {
				dp[i] = val
			}
		}
		nv := dp[i] - pref[i]
		for len(dq) > dqFront && dq[len(dq)-1].val >= nv {
			dq = dq[:len(dq)-1]
		}
		dq = append(dq, pair{i, nv})
	}
	fmt.Fprintln(writer, dp[n])
}
