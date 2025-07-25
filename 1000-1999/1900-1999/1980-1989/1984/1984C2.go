package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func modPow(a int64, b int64) int64 {
	res := int64(1)
	a %= mod
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func solve(n int, arr []int64) int64 {
	pref := make([]int64, n+1)
	cntNonNeg := make([]int, n+1)
	var sum int64
	minPref := int64(0)
	for i := 1; i <= n; i++ {
		sum += arr[i-1]
		pref[i] = sum
		cntNonNeg[i] = cntNonNeg[i-1]
		if sum >= 0 {
			cntNonNeg[i]++
		}
		if sum < minPref {
			minPref = sum
		}
	}
	if minPref >= 0 {
		return modPow(2, int64(n))
	}
	ans := int64(0)
	for j := 1; j <= n; j++ {
		if pref[j] == minPref {
			nonNegBefore := cntNonNeg[j-1]
			exp := nonNegBefore + (n - j)
			ans = (ans + modPow(2, int64(exp))) % mod
		}
	}
	return ans
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		fmt.Fprintln(writer, solve(n, arr))
	}
}
