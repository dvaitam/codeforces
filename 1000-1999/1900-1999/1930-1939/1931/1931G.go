package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353
const MAX int = 4000005

var fact [MAX]int64
var invfact [MAX]int64

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

func init() {
	fact[0] = 1
	for i := 1; i < MAX; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invfact[MAX-1] = modPow(fact[MAX-1], MOD-2)
	for i := MAX - 2; i >= 0; i-- {
		invfact[i] = invfact[i+1] * int64(i+1) % MOD
	}
}

func comb(n, k int64) int64 {
	if k < 0 || k > n {
		return 0
	}
	ni := int(n)
	ki := int(k)
	return fact[ni] * invfact[ki] % MOD * invfact[ni-ki] % MOD
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var c1, c2, c3, c4 int64
		fmt.Fscan(in, &c1, &c2, &c3, &c4)

		e01 := c1 // 0->1
		e00 := c2 // 0->0
		e11 := c3 // 1->1
		e10 := c4 // 1->0

		// no cross edges
		if e01+e10 == 0 {
			if e00 > 0 && e11 > 0 {
				fmt.Fprintln(out, 0)
			} else {
				fmt.Fprintln(out, 1)
			}
			continue
		}
		if e01-e10 > 1 || e10-e01 > 1 {
			fmt.Fprintln(out, 0)
			continue
		}
		var ans int64
		if e01 == e10 {
			// start from 0
			n0 := e10 + 1
			n1 := e01
			ways0 := comb(e00+n0-1, n0-1)
			ways1 := comb(e11+n1-1, n1-1)
			ans = ways0 * ways1 % MOD
			// start from 1
			n0 = e10
			n1 = e01 + 1
			ways0 = comb(e00+n0-1, n0-1)
			ways1 = comb(e11+n1-1, n1-1)
			ans = (ans + ways0*ways1) % MOD
		} else if e01 == e10+1 { // start at 0
			n0 := e10 + 1
			n1 := e01
			ans = comb(e00+n0-1, n0-1) * comb(e11+n1-1, n1-1) % MOD
		} else { // e10 == e01+1 start at 1
			n0 := e10
			n1 := e01 + 1
			ans = comb(e00+n0-1, n0-1) * comb(e11+n1-1, n1-1) % MOD
		}
		fmt.Fprintln(out, ans)
	}
}
