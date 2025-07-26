package main

import (
	"bufio"
	"fmt"
	"os"
)

func segmentDiff(prefix []int64, k int) int64 {
	n := len(prefix) - 1
	var maxSum int64 = -1 << 63
	var minSum int64 = 1<<63 - 1
	for i := k; i <= n; i += k {
		sum := prefix[i] - prefix[i-k]
		if sum > maxSum {
			maxSum = sum
		}
		if sum < minSum {
			minSum = sum
		}
	}
	return maxSum - minSum
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		prefix := make([]int64, n+1)
		for i := 0; i < n; i++ {
			prefix[i+1] = prefix[i] + a[i]
		}
		var ans int64
		for d := 1; d*d <= n; d++ {
			if n%d == 0 {
				diff := segmentDiff(prefix, d)
				if diff > ans {
					ans = diff
				}
				if d != n/d {
					diff2 := segmentDiff(prefix, n/d)
					if diff2 > ans {
						ans = diff2
					}
				}
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
