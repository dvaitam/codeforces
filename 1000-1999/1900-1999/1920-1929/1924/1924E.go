package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func modPow(a, b int64) int64 {
	res := int64(1)
	a %= MOD
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func inv(x int64) int64 { return modPow(x, MOD-2) }

var kGlobal int64
var memo map[[2]int64]int64

func solve(h, w int64) int64 {
	if h*w < kGlobal {
		return 0
	}
	key := [2]int64{h, w}
	if val, ok := memo[key]; ok {
		return val
	}
	var sum int64
	for i := int64(1); i < w; i++ {
		sum = (sum + solve(h, i)) % MOD
	}
	for j := int64(1); j < h; j++ {
		sum = (sum + solve(j, w)) % MOD
	}
	res := (1 + sum*inv(h+w-2)%MOD) % MOD
	memo[key] = res
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m, k int64
		fmt.Fscan(in, &n, &m, &k)
		kGlobal = k
		memo = make(map[[2]int64]int64)
		ans := solve(n, m)
		fmt.Fprintln(out, ans)
	}
}
