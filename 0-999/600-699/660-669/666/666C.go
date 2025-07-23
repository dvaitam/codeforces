package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007
const MAXN = 100000

var fact [MAXN + 1]int64
var invfact [MAXN + 1]int64
var pow25 [MAXN + 1]int64

func modPow(a, e int64) int64 {
	res := int64(1)
	a %= MOD
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func initPrecalc() {
	fact[0] = 1
	for i := 1; i <= MAXN; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invfact[MAXN] = modPow(fact[MAXN], MOD-2)
	for i := MAXN; i >= 1; i-- {
		invfact[i-1] = invfact[i] * int64(i) % MOD
	}
	pow25[0] = 1
	for i := 1; i <= MAXN; i++ {
		pow25[i] = pow25[i-1] * 25 % MOD
	}
}

func C(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * invfact[k] % MOD * invfact[n-k] % MOD
}

func main() {
	initPrecalc()
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var m int
	if _, err := fmt.Fscan(reader, &m); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)

	curLen := len(s)
	dp := make(map[int][]int64)
	dp[curLen] = make([]int64, curLen+1)
	dp[curLen][curLen] = 1

	for ; m > 0; m-- {
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			fmt.Fscan(reader, &s)
			curLen = len(s)
			if _, ok := dp[curLen]; !ok {
				arr := make([]int64, curLen+1)
				arr[curLen] = 1
				dp[curLen] = arr
			}
		} else if t == 2 {
			var n int
			fmt.Fscan(reader, &n)
			if n < curLen {
				fmt.Fprintln(writer, 0)
				continue
			}
			arr := dp[curLen]
			for len(arr) <= n {
				i := len(arr)
				val := (arr[i-1]*26 + C(i-1, curLen-1)*pow25[i-curLen]) % MOD
				arr = append(arr, val)
			}
			dp[curLen] = arr
			fmt.Fprintln(writer, arr[n]%MOD)
		}
	}
}
