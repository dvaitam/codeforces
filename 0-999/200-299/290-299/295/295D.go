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
			res = (res * a) % MOD
		}
		a = (a * a) % MOD
		e >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	if m < 2 {
		fmt.Println(0)
		return
	}

	maxN := m + 2*n + 5
	fac := make([]int64, maxN+1)
	ifac := make([]int64, maxN+1)
	fac[0] = 1
	for i := 1; i <= maxN; i++ {
		fac[i] = (fac[i-1] * int64(i)) % MOD
	}
	ifac[maxN] = modPow(fac[maxN], MOD-2)
	for i := maxN - 1; i >= 0; i-- {
		ifac[i] = (ifac[i+1] * int64(i+1)) % MOD
	}

	comb := func(nv, kv int) int64 {
		if nv < 0 || kv < 0 || kv > nv {
			return 0
		}
		return (((fac[nv] * ifac[kv]) % MOD) * ifac[nv-kv]) % MOD
	}

	A := make([]int64, n)
	Astar := make([]int64, n)
	S0 := make([]int64, n)
	S1 := make([]int64, n)

	var ans int64 = 0
	for W := 0; W <= m-2; W++ {
		for k := 0; k < n; k++ {
			top := W + 2*k
			A[k] = comb(top, 2*k)
		}
		for k := 0; k < n; k++ {
			if k == 0 {
				// No rows before apex, unique choice
				Astar[k] = 1
			} else {
				t1 := comb(W-1+2*k, 2*k)
				t2 := comb(W-2+2*k, 2*k)
				val := (2*t1 - t2) % MOD
				if val < 0 {
					val += MOD
				}
				Astar[k] = val
			}
		}
		// Prefix sums for A
		var s0, s1 int64 = 0, 0
		for i := 0; i < n; i++ {
			s0 += A[i]
			if s0 >= MOD {
				s0 -= MOD
			}
			S0[i] = s0
			s1 = (s1 + int64(i)*A[i]) % MOD
			S1[i] = s1
		}

		var T int64 = 0
		for k := 0; k < n; k++ {
			Nk := n - 1 - k
			rk := (int64(n-k)*S0[Nk] - S1[Nk]) % MOD
			if rk < 0 {
				rk += MOD
			}
			T = (T + Astar[k]*rk) % MOD
		}
		mult := int64(m - W - 1)
		ans = (ans + T*mult) % MOD
	}
	if ans < 0 {
		ans += MOD
	}
	fmt.Println(ans % MOD)
}
