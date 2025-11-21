package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func modPow(a, e int64) int64 {
	res := int64(1)
	base := a % mod
	for e > 0 {
		if e&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		e >>= 1
	}
	return res
}

func solveCase(arr []int64, k int64, m int64) int64 {
	n := int64(len(arr))
	sum := arr[n-1]
	restSum := sum - k
	if restSum < 0 {
		return 0
	}

	ans := int64(0)
	for i := int64(0); i < n-1; i++ {
		gaps := arr[i]
		part := comb(restSum+gaps-1, gaps-1) % m
		ans = (ans + part) % m
	}
	return ans
}

func comb(n, k int64) int64 {
	if n < 0 || k < 0 || k > n {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	num := int64(1)
	den := int64(1)
	for i := int64(1); i <= k; i++ {
		num = num * (n - i + 1) % mod
		den = den * i % mod
	}
	return num * modPow(den, mod-2) % mod
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int64
		fmt.Fscan(in, &n)
		arr := make([]int64, n)
		for i := int64(0); i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		var k, m int64
		fmt.Fscan(in, &k, &m)
		fmt.Fprintln(out, solveCase(arr, k, m))
	}
}

