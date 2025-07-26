package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

var p2 []int64
var xVal int

func pow2(n int) int64 {
	if n < len(p2) {
		return p2[n]
	}
	for len(p2) <= n {
		p2 = append(p2, p2[len(p2)-1]*2%MOD)
	}
	return p2[n]
}

func solveAll(arr []int, bit int) int64 {
	if len(arr) == 0 {
		return 0
	}
	if bit < 0 {
		return (pow2(len(arr)) + MOD - 1) % MOD
	}
	var zero, one []int
	mask := 1 << bit
	for _, v := range arr {
		if v&mask == 0 {
			zero = append(zero, v)
		} else {
			one = append(one, v)
		}
	}
	if (xVal>>bit)&1 == 0 {
		left := solveAll(zero, bit-1)
		right := solveAll(one, bit-1)
		return (left + right) % MOD
	}
	cross := solveCross(zero, one, bit-1)
	return cross % MOD
}

func solveCross(a, b []int, bit int) int64 {
	if len(a) == 0 || len(b) == 0 {
		return (pow2(len(a)+len(b)) + MOD - 1) % MOD
	}
	if bit < 0 {
		res := (pow2(len(a)) - 1) % MOD
		res = res * ((pow2(len(b)) - 1) % MOD) % MOD
		return res
	}
	var a0, a1, b0, b1 []int
	mask := 1 << bit
	for _, v := range a {
		if v&mask == 0 {
			a0 = append(a0, v)
		} else {
			a1 = append(a1, v)
		}
	}
	for _, v := range b {
		if v&mask == 0 {
			b0 = append(b0, v)
		} else {
			b1 = append(b1, v)
		}
	}
	if (xVal>>bit)&1 == 0 {
		res := (solveCross(a0, b0, bit-1) + solveCross(a1, b1, bit-1)) % MOD
		return res
	}
	res := solveCross(a0, b0, bit-1)
	res += solveCross(a0, b1, bit-1)
	res %= MOD
	res += solveCross(a1, b0, bit-1)
	res %= MOD
	res += solveCross(a1, b1, bit-1)
	res %= MOD
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n, &xVal)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	p2 = []int64{1}
	ans := solveAll(arr, 29)
	fmt.Fprintln(out, ans%MOD)
}
