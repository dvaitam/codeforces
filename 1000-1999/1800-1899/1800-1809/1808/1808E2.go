package main

import (
	"bufio"
	"fmt"
	"os"
)

func modPow(base, exp, mod int64) int64 {
	res := int64(1 % mod)
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return res
}

func polyMul(a, b []int64, k int, mod int64) []int64 {
	res := make([]int64, k)
	for i := 0; i < k; i++ {
		if a[i] == 0 {
			continue
		}
		ai := a[i]
		for j := 0; j < k; j++ {
			if b[j] == 0 {
				continue
			}
			idx := (i + j) % k
			res[idx] = (res[idx] + ai*b[j]) % mod
		}
	}
	return res
}

func polyPow(base []int64, exp int64, k int, mod int64) []int64 {
	res := make([]int64, k)
	res[0] = 1 % mod
	for exp > 0 {
		if exp&1 == 1 {
			res = polyMul(res, base, k, mod)
		}
		base = polyMul(base, base, k, mod)
		exp >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int64
	var k int
	var m int64
	if _, err := fmt.Fscan(in, &n, &k, &m); err != nil {
		return
	}

	total := modPow(int64(k), n, m)
	complement := int64(0)

	for s := 0; s < k; s++ {
		base := make([]int64, k)
		for d := 0; d < k; d++ {
			if (2*d)%k != s {
				base[d] = 1
			}
		}
		coeff := polyPow(base, n, k, m)
		complement = (complement + coeff[s]) % m
	}

	ans := (total - complement) % m
	if ans < 0 {
		ans += m
	}
	fmt.Println(ans)
}
