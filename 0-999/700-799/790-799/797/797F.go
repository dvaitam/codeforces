package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Hole struct {
	p int64
	c int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	mice := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &mice[i])
	}
	holes := make([]Hole, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &holes[i].p, &holes[i].c)
	}

	sort.Slice(mice, func(i, j int) bool { return mice[i] < mice[j] })
	sort.Slice(holes, func(i, j int) bool { return holes[i].p < holes[j].p })

	const INF int64 = 1<<60 - 1
	dp := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		dp[i] = INF
	}

	prefix := make([]int64, n+1)
	for _, h := range holes {
		prefix[0] = 0
		for i := 1; i <= n; i++ {
			diff := mice[i-1] - h.p
			if diff < 0 {
				diff = -diff
			}
			prefix[i] = prefix[i-1] + diff
		}

		newDP := make([]int64, n+1)
		dequeIdx := make([]int, 0, n+1)
		dequeVal := make([]int64, 0, n+1)
		front := 0
		for j := 0; j <= n; j++ {
			for front < len(dequeIdx) && dequeIdx[front] < j-h.c {
				front++
			}

			val := dp[j] - prefix[j]
			for len(dequeIdx) > front && dequeVal[len(dequeVal)-1] >= val {
				dequeIdx = dequeIdx[:len(dequeIdx)-1]
				dequeVal = dequeVal[:len(dequeVal)-1]
			}
			dequeIdx = append(dequeIdx, j)
			dequeVal = append(dequeVal, val)

			best := dequeVal[front]
			newDP[j] = prefix[j] + best
			if newDP[j] > INF {
				newDP[j] = INF
			}
		}
		dp = newDP
	}

	if dp[n] >= INF {
		fmt.Fprintln(out, -1)
	} else {
		fmt.Fprintln(out, dp[n])
	}
}
