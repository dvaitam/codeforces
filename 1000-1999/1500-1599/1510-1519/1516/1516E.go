package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func combBig(n int64, r int, invFact []int64) int64 {
	if r < 0 || int64(r) > n {
		return 0
	}
	res := int64(1)
	for i := 0; i < r; i++ {
		res = res * ((n - int64(i)) % MOD) % MOD
	}
	res = res * invFact[r] % MOD
	return res
}

func evalPoly(y []int64, n int64, invFact []int64) int64 {
	d := len(y) - 1
	if n <= int64(d) {
		return y[n] % MOD
	}
	ans := int64(0)
	for i := 0; i <= d; i++ {
		c1 := combBig(n, i, invFact)
		c2 := combBig(n-1-int64(i), d-i, invFact)
		sign := int64(1)
		if (d-i)%2 == 1 {
			sign = MOD - 1
		}
		ans = (ans + y[i]*c1%MOD*c2%MOD*sign) % MOD
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int64
	var k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	maxN := 2 * k
	st := make([][]int64, maxN+1)
	for i := range st {
		st[i] = make([]int64, maxN+1)
	}
	st[0][0] = 1
	for i := 1; i <= maxN; i++ {
		for j := 1; j <= i; j++ {
			st[i][j] = (st[i-1][j-1] + int64(i-1)*st[i-1][j]) % MOD
		}
	}

	fact := make([]int64, maxN+1)
	invFact := make([]int64, maxN+1)
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[maxN] = modPow(fact[maxN], MOD-2)
	for i := maxN; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}

	val := make([]int64, k+1)
	for m := 0; m <= k; m++ {
		d := 2 * m
		y := make([]int64, d+1)
		for x := 0; x <= d; x++ {
			if x < m {
				y[x] = 0
			} else {
				y[x] = st[x][x-m]
			}
		}
		val[m] = evalPoly(y, n, invFact)
	}

	for j := 1; j <= k; j++ {
		tot := int64(0)
		limit := j
		if int64(limit) > n-1 {
			limit = int(n - 1)
		}
		for m := 0; m <= limit; m++ {
			if (j-m)%2 == 0 {
				tot = (tot + val[m]) % MOD
			}
		}
		if j > 1 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, tot)
	}
	fmt.Fprintln(out)
}
