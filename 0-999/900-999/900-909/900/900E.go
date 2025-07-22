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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var s string
	fmt.Fscan(in, &s)
	var m int
	fmt.Fscan(in, &m)

	// precompute expected characters for two global patterns
	expectedA := make([]byte, n+1)
	expectedB := make([]byte, n+1)
	for i := 1; i <= n; i++ {
		if i%2 == 1 {
			expectedA[i] = 'a'
			expectedB[i] = 'b'
		} else {
			expectedA[i] = 'b'
			expectedB[i] = 'a'
		}
	}

	mismA := make([]int, n+1)
	mismB := make([]int, n+1)
	qpref := make([]int, n+1)
	for i := 1; i <= n; i++ {
		ch := s[i-1]
		mismA[i] = mismA[i-1]
		mismB[i] = mismB[i-1]
		qpref[i] = qpref[i-1]
		if ch == '?' {
			qpref[i]++
		}
		if ch != '?' && ch != expectedA[i] {
			mismA[i]++
		}
		if ch != '?' && ch != expectedB[i] {
			mismB[i]++
		}
	}

	const INF int = 1 << 30
	dpMax := make([]int, n+1)
	dpCost := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dpMax[i] = -INF
		dpCost[i] = INF
	}

	for i := 0; i < n; i++ {
		if dpMax[i] == -INF {
			continue
		}
		// Skip current character
		if dpMax[i] > dpMax[i+1] || (dpMax[i] == dpMax[i+1] && dpCost[i] < dpCost[i+1]) {
			dpMax[i+1] = dpMax[i]
			dpCost[i+1] = dpCost[i]
		}
		// Try to place an occurrence starting at position i+1
		if i+m <= n {
			var mism int
			if (i+1)%2 == 1 {
				mism = mismA[i+m] - mismA[i]
			} else {
				mism = mismB[i+m] - mismB[i]
			}
			if mism == 0 {
				cost := qpref[i+m] - qpref[i]
				occ := dpMax[i] + 1
				newCost := dpCost[i] + cost
				if occ > dpMax[i+m] || (occ == dpMax[i+m] && newCost < dpCost[i+m]) {
					dpMax[i+m] = occ
					dpCost[i+m] = newCost
				}
			}
		}
	}

	fmt.Fprintln(out, dpCost[n])
}
