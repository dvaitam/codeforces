package main

import (
	"bufio"
	"fmt"
	"os"
)

func modPow(base, exp, mod int64) int64 {
	result := int64(1)
	for exp > 0 {
		if exp&1 == 1 {
			result = (result * base) % mod
		}
		base = (base * base) % mod
		exp >>= 1
	}
	return result
}

func maxSubarraySum(arr []int64) int64 {
	maxSum := arr[0]
	cur := arr[0]
	for i := 1; i < len(arr); i++ {
		if cur+arr[i] > arr[i] {
			cur = cur + arr[i]
		} else {
			cur = arr[i]
		}
		if cur > maxSum {
			maxSum = cur
		}
	}
	return maxSum
}

func normMod(x, mod int64) int64 {
	x %= mod
	if x < 0 {
		x += mod
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	const MOD int64 = 1_000_000_007

	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		arr := make([]int64, n)
		var sum int64
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
			sum += arr[i]
		}

		maxSub := maxSubarraySum(arr)
		if maxSub <= 0 {
			fmt.Fprintln(out, normMod(sum, MOD))
			continue
		}

		pow2 := modPow(2, int64(k), MOD)
		if sum >= maxSub {
			ans := normMod(sum, MOD)
			ans = (ans * pow2) % MOD
			fmt.Fprintln(out, ans)
		} else {
			inc := normMod(maxSub, MOD)
			inc = (inc * ((pow2 - 1 + MOD) % MOD)) % MOD
			ans := (normMod(sum, MOD) + inc) % MOD
			fmt.Fprintln(out, ans)
		}
	}
}
