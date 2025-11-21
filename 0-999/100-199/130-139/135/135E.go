package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 1000000007

var fac, ifac []int64

func modPow(base, exp int64) int64 {
	base %= mod
	if base < 0 {
		base += mod
	}
	res := int64(1)
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return res
}

func initFactorials(limit int) {
	fac = make([]int64, limit+1)
	ifac = make([]int64, limit+1)
	fac[0] = 1
	for i := 1; i <= limit; i++ {
		fac[i] = fac[i-1] * int64(i) % mod
	}
	ifac[limit] = modPow(fac[limit], mod-2)
	for i := limit; i >= 1; i-- {
		ifac[i-1] = ifac[i] * int64(i) % mod
	}
}

func fall(n, r int) int64 {
	if r < 0 || r > n {
		return 0
	}
	return fac[n] * ifac[n-r] % mod
}

func getPow(powVals []int64, idx int, k int, w int64) int64 {
	if idx <= 0 {
		return 1
	}
	if idx < len(powVals) && powVals[idx] != 0 {
		return powVals[idx]
	}
	exp := w - int64(idx)
	if exp < 0 {
		return 0
	}
	return modPow(int64(k), exp)
}

func countLen(k int, w int64, t int, powVals []int64) int64 {
	if t > k {
		return 0
	}
	if int64(t) <= w {
		ff := fall(k, t)
		powVal := getPow(powVals, t, k, w)
		return ff * ff % mod * powVal % mod
	}
	d := int(int64(t) - w)
	if d <= 0 || d > k {
		return 0
	}
	restSize := k - d
	if int64(restSize) < w {
		return 0
	}
	rest := fall(restSize, int(w))
	overlap := fall(k, d)
	return overlap * rest % mod * rest % mod
}

func countLenNext(k int, w int64, t int, powVals []int64) int64 {
	if t+1 > k && int64(t) <= w-2 {
		return 0
	}
	if int64(t) <= w-2 {
		ff := fall(k, t+1)
		powVal := getPow(powVals, t+2, k, w)
		return ff * ff % mod * powVal % mod
	}
	d := int(int64(t) - w + 2)
	if d <= 0 || d > k {
		return 0
	}
	restSize := k - d
	target := w - 1
	if target < 0 || int64(restSize) < target {
		return 0
	}
	rest := fall(restSize, int(target))
	overlap := fall(k, d)
	return overlap * rest % mod * rest % mod
}

func solve(k int, w int64) int64 {
	if k == 0 {
		return 0
	}
	powLimit := k + 2
	if int64(powLimit) > w {
		powLimit = int(w)
	}
	powVals := make([]int64, powLimit+2)
	if powLimit >= 1 {
		powVals[1] = modPow(int64(k), w-1)
		inv := modPow(int64(k), mod-2)
		for i := 2; i <= powLimit; i++ {
			powVals[i] = powVals[i-1] * inv % mod
		}
	}
	ans := int64(0)
	for t := 1; t <= k; t++ {
		ans += countLen(k, w, t, powVals)
		ans -= countLenNext(k, w, t, powVals)
	}
	ans %= mod
	if ans < 0 {
		ans += mod
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var k int
	var w int64
	if _, err := fmt.Fscan(in, &k, &w); err != nil {
		return
	}
	initFactorials(k + 2)
	fmt.Println(solve(k, w))
}
