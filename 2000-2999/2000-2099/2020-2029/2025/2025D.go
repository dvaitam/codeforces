package main

import (
	"bufio"
	"fmt"
	"os"
)

const negInf = -1 << 60

func applySegment(dp []int, k int, strCnt, intCnt []int, m int) {
	// prefix sums for strength
	acc := 0
	for i := 1; i <= m; i++ {
		acc += strCnt[i]
		strCnt[i] = acc
	}
	acc = 0
	for i := 1; i <= m; i++ {
		acc += intCnt[i]
		intCnt[i] = acc
	}

	for s := 0; s <= k; s++ {
		if dp[s] == negInf {
			continue
		}
		dp[s] += strCnt[s] + intCnt[k-s]
	}

	for i := 1; i <= m; i++ {
		strCnt[i] = 0
		intCnt[i] = 0
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	dp := make([]int, m+1)
	tmp := make([]int, m+1)
	for i := 0; i <= m; i++ {
		dp[i] = negInf
		tmp[i] = negInf
	}
	dp[0] = 0

	strCnt := make([]int, m+1)
	intCnt := make([]int, m+1)

	pointsUsed := 0

	applySeg := func() {
		applySegment(dp, pointsUsed, strCnt, intCnt, m)
	}

	for i := 0; i < n; i++ {
		var v int
		fmt.Fscan(in, &v)
		if v == 0 {
			applySeg()
			// distribute new point
			newK := pointsUsed + 1
			for j := 0; j <= newK; j++ {
				tmp[j] = negInf
			}
			for s := 0; s <= pointsUsed; s++ {
				val := dp[s]
				if val == negInf {
					continue
				}
				if val > tmp[s] {
					tmp[s] = val
				}
				if val > tmp[s+1] {
					tmp[s+1] = val
				}
			}
			dp, tmp = tmp, dp
			pointsUsed = newK
		} else if v > 0 {
			if v > m {
				v = m
			}
			intCnt[v]++
		} else {
			v = -v
			if v > m {
				v = m
			}
			strCnt[v]++
		}
	}

	applySeg()

	best := 0
	upper := pointsUsed
	for s := 0; s <= upper; s++ {
		if dp[s] > best {
			best = dp[s]
		}
	}
	fmt.Fprintln(out, best)
}
