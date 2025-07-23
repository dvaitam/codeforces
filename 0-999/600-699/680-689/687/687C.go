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

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	coins := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &coins[i])
	}

	units := (k + 64) / 64
	dp := make([][]uint64, k+1)
	for i := range dp {
		dp[i] = make([]uint64, units)
	}
	dp[0][0] = 1

	for _, c := range coins {
		for j := k; j >= c; j-- {
			src := dp[j-c]
			dst := dp[j]
			// OR without taking the coin into the subset for x
			for idx := 0; idx < units; idx++ {
				dst[idx] |= src[idx]
			}
			// OR when taking the coin into the subset for x (shift by c)
			w := c / 64
			b := uint(c % 64)
			if b == 0 {
				for idx := units - 1; idx >= w; idx-- {
					dst[idx] |= src[idx-w]
				}
			} else {
				for idx := units - 1; idx > w; idx-- {
					dst[idx] |= src[idx-w]<<b | src[idx-w-1]>>(64-b)
				}
				if w < units {
					dst[w] |= src[0] << b
				}
			}
		}
	}

	var ans []int
	for x := 0; x <= k; x++ {
		if (dp[k][x/64]>>uint(x%64))&1 == 1 {
			ans = append(ans, x)
		}
	}

	fmt.Fprintln(out, len(ans))
	for i, v := range ans {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
}
