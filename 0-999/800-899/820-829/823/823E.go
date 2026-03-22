package main

import (
	"fmt"
)

func main() {
	var K int
	fmt.Scan(&K)

	mod := int64(1000000007)

	dp := make([]int64, K+2)
	dp[0] = 1
	dp[1] = 1

	for k := 2; k <= K; k++ {
		limit := K - k + 1
		t := make([]int64, limit+2)

		for c := 0; c <= limit+1; c++ {
			var sum int64 = 0
			for i := 0; i <= c; i++ {
				val := (dp[i] * dp[c-i]) % mod
				sum = (sum + val) % mod
			}
			t[c] = sum
		}

		nextDp := make([]int64, limit+1)
		for c := 0; c <= limit; c++ {
			var res int64 = 0
			if c > 0 {
				res = (res + t[c-1]) % mod
			}
			term2 := (int64(2*c+1) * t[c]) % mod
			res = (res + term2) % mod

			term3 := (int64(c+1) * int64(c)) % mod
			term3 = (term3 * t[c+1]) % mod
			res = (res + term3) % mod

			nextDp[c] = res
		}
		dp = nextDp
	}

	fmt.Println(dp[1])
}
