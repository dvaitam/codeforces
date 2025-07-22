package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1e9 + 7

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var k int
	var pa, pb int64
	if _, err := fmt.Fscan(reader, &k, &pa, &pb); err != nil {
		return
	}

	denom := (pa + pb) % mod
	invDen := modPow(denom, mod-2)
	pA := pa % mod * invDen % mod
	pB := pb % mod * invDen % mod
	ratio := pa % mod * modPow(pb%mod, mod-2) % mod

	dp := make([][]int64, k+1)
	for i := range dp {
		dp[i] = make([]int64, k+1)
	}

	for j := k - 1; j >= 0; j-- {
		dp[k][j] = (int64(j) + int64(k) + ratio) % mod
		for i := k - 1; i >= 1; i-- {
			valA := pA * dp[i+1][j] % mod
			var valB int64
			if j+i >= k {
				valB = pB * int64(j+i) % mod
			} else {
				valB = pB * dp[i][j+i] % mod
			}
			dp[i][j] = (valA + valB) % mod
		}
		dp[0][j] = dp[1][j]
	}

	fmt.Println(dp[0][0] % mod)
}
