package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Pair holds total profit and selected jobs bitset
type Pair struct {
	profit int64
	used   []bool
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	tt := make([]int, n)
	dd := make([]int, n)
	pp := make([]int64, n)
	for i := 0; i < n; i++ {
		var t, d int
		var p int64
		fmt.Fscan(reader, &t, &d, &p)
		tt[i] = t
		dd[i] = d
		pp[i] = p
	}
	// sorted job indices by deadline
	idx := make([]int, n)
	for i := 0; i < n; i++ {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool {
		return dd[idx[i]] < dd[idx[j]]
	})
	// find maximum deadline
	maxD := 0
	for i := 0; i < n; i++ {
		if dd[i] > maxD {
			maxD = dd[i]
		}
	}
	// dp[t] = best Pair finishing at time t
	dp := make([]Pair, maxD+1)
	for t := range dp {
		dp[t].profit = 0
		dp[t].used = make([]bool, n)
	}
	best := Pair{profit: 0, used: make([]bool, n)}
	// process jobs in increasing deadline
	for _, x := range idx {
		t := tt[x]
		d := dd[x]
		p := pp[x]
		// possible start times: 0 <= tim <= d-t-1
		for tim := d - t - 1; tim >= 0; tim-- {
			cur := dp[tim]
			newProfit := cur.profit + p
			newT := tim + t
			// if better than current dp[newT]
			if newProfit > dp[newT].profit {
				// copy used set
				newUsed := make([]bool, n)
				copy(newUsed, cur.used)
				newUsed[x] = true
				dp[newT].profit = newProfit
				dp[newT].used = newUsed
				// update global best
				if newProfit > best.profit {
					best.profit = newProfit
					best.used = newUsed
				}
			}
		}
	}
	// collect result jobs
	lis := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if best.used[i] {
			lis = append(lis, i)
		}
	}
	// sort by deadline
	sort.Slice(lis, func(i, j int) bool {
		return dd[lis[i]] < dd[lis[j]]
	})
	// output
	fmt.Fprintln(writer, best.profit)
	fmt.Fprintln(writer, len(lis))
	for _, e := range lis {
		fmt.Fprintf(writer, "%d ", e+1)
	}
	fmt.Fprintln(writer)
}
