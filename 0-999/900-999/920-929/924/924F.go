package main

import (
	"bufio"
	"fmt"
	"os"
)

func isKBeautiful(x int64, k int) bool {
	if x == 0 {
		return k >= 0 // 0 is not positive though, but we won't call with 0
	}
	digits := []int{}
	for x > 0 {
		digits = append(digits, int(x%10))
		x /= 10
	}
	sum := 0
	for _, d := range digits {
		sum += d
	}
	dp := make([]bool, sum+1)
	dp[0] = true
	for _, d := range digits {
		next := make([]bool, sum+1)
		for s := 0; s <= sum; s++ {
			if dp[s] {
				next[s] = true
				if s+d <= sum {
					next[s+d] = true
				}
			}
		}
		dp = next
	}
	for s := 0; s <= sum; s++ {
		if dp[s] {
			diff := sum - 2*s
			if diff < 0 {
				diff = -diff
			}
			if diff <= k {
				return true
			}
		}
	}
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for i := 0; i < n; i++ {
		var l, r int64
		var k int
		fmt.Fscan(in, &l, &r, &k)
		count := int64(0)
		for x := l; x <= r; x++ {
			if isKBeautiful(x, k) {
				count++
			}
		}
		fmt.Fprintln(writer, count)
	}
}
