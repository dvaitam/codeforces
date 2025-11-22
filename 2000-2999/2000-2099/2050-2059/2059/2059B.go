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

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		evenLeft := k / 2
		cur := 0         // index of the first unassigned element (0-based)
		need := int64(1) // expected value at current position in b
		answer := int64(-1)

		for evenLeft > 0 {
			lenRem := n - cur
			// minimal elements needed from here: current odd (1) + remaining even/odd pairs
			maxSkip := lenRem - (2*evenLeft - 1) // extra elements we may assign to the current odd subarray
			if maxSkip < 0 {
				maxSkip = 0
			}

			startRange := cur + 1         // first possible position of current even subarray
			endRange := cur + maxSkip + 1 // last position we can start and still fit everything
			if endRange >= n {
				endRange = n - 1
			}

			mismatchFound := false
			for i := startRange; i <= endRange; i++ {
				if a[i] != need {
					answer = need
					mismatchFound = true
					break
				}
			}
			if mismatchFound {
				break
			}

			start := startRange
			limit := n - 2*(evenLeft-1) - 1 // last index we can include in this even subarray
			for i := start; i <= limit; i++ {
				if a[i] != need {
					answer = need
					mismatchFound = true
					break
				}
				need++
			}
			if mismatchFound {
				break
			}

			cur = limit + 1
			evenLeft--
		}

		if answer == -1 {
			answer = need // mismatch happens at the appended 0
		}

		fmt.Fprintln(out, answer)
	}
}
