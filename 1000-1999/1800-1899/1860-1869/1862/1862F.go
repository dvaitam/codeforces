package main

import (
	"bufio"
	"fmt"
	"os"
)

func can(t int64, w, f int64, sum int, dp []bool) bool {
	maxW := w * t
	if maxW > int64(sum) {
		maxW = int64(sum)
	}
	maxF := f * t
	for i := int64(0); i <= maxW; i++ {
		if dp[i] {
			if int64(sum)-i <= maxF {
				return true
			}
		}
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var w, f int64
		fmt.Fscan(reader, &w, &f)
		var n int
		fmt.Fscan(reader, &n)
		s := make([]int, n)
		sum := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &s[i])
			sum += s[i]
		}
		dp := make([]bool, sum+1)
		dp[0] = true
		for _, v := range s {
			for j := sum; j >= v; j-- {
				if dp[j-v] {
					dp[j] = true
				}
			}
		}
		low := int64(0)
		high := int64(sum)
		if w < f {
			if int64(sum)%w != 0 {
				if int64(sum)/w+1 > high {
					high = int64(sum)/w + 1
				}
			} else {
				if int64(sum)/w > high {
					high = int64(sum) / w
				}
			}
		} else {
			if int64(sum)%f != 0 {
				if int64(sum)/f+1 > high {
					high = int64(sum)/f + 1
				}
			} else {
				if int64(sum)/f > high {
					high = int64(sum) / f
				}
			}
		}
		for low < high {
			mid := (low + high) / 2
			if can(mid, w, f, sum, dp) {
				high = mid
			} else {
				low = mid + 1
			}
		}
		fmt.Fprintln(writer, low)
	}
}
