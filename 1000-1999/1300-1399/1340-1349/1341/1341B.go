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

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		// mark peaks
		isPeak := make([]int, n)
		for i := 1; i < n-1; i++ {
			if a[i] > a[i-1] && a[i] > a[i+1] {
				isPeak[i] = 1
			}
		}
		// prefix sums of peaks
		pref := make([]int, n)
		for i := 1; i < n; i++ {
			pref[i] = pref[i-1] + isPeak[i]
		}
		bestPeaks := -1
		bestL := 0
		for l := 0; l+k-1 < n; l++ {
			// peaks strictly inside segment (l, l+k-1)
			peaks := pref[l+k-2] - pref[l]
			if peaks > bestPeaks {
				bestPeaks = peaks
				bestL = l
			}
		}
		fmt.Fprintf(writer, "%d %d\n", bestPeaks+1, bestL+1)
	}
}
