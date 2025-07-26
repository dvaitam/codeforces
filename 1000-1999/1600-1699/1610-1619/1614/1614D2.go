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
	arr := make([]int, n)
	maxVal := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
		if arr[i] > maxVal {
			maxVal = arr[i]
		}
	}

	freq := make([]int, maxVal+1)
	for _, v := range arr {
		freq[v]++
	}

	cnt := make([]int, maxVal+1)
	for i := 1; i <= maxVal; i++ {
		for j := i; j <= maxVal; j += i {
			if freq[j] > 0 {
				cnt[i] += freq[j]
			}
		}
	}

	dp := make([]int64, maxVal+1)
	for i := maxVal; i >= 1; i-- {
		if cnt[i] == 0 {
			continue
		}
		best := int64(cnt[i]) * int64(i)
		for j := i * 2; j <= maxVal; j += i {
			if cnt[j] == 0 {
				continue
			}
			val := dp[j] + int64(cnt[i]-cnt[j])*int64(i)
			if val > best {
				best = val
			}
		}
		dp[i] = best
	}

	fmt.Fprintln(out, dp[1])
}
