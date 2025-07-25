package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1_000_000_007

type Matrix struct{ a, b, c, d int64 }

func mulMat(x, y Matrix) Matrix {
	return Matrix{
		(x.a*y.a + x.b*y.c) % MOD,
		(x.a*y.b + x.b*y.d) % MOD,
		(x.c*y.a + x.d*y.c) % MOD,
		(x.c*y.b + x.d*y.d) % MOD,
	}
}
func mulVec(m Matrix, x, y int64) (int64, int64) {
	return (m.a*x + m.b*y) % MOD, (m.c*x + m.d*y) % MOD
}
func main() {
	in := bufio.NewReader(os.Stdin)
	var m int
	var n int64
	fmt.Fscan(in, &m, &n)
	s := make([]int64, m)
	l := make([]int64, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &s[i])
	}
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &l[i])
	}
	t := make([]int64, m)
	var sum_ts, sum_ss, sum_tl, sum_sl int64
	var sum_s, sum_t int64
	for i := 0; i < m; i++ {
		t[i] = (s[i] + l[i]) % MOD
		s_mod := s[i] % MOD
		l_mod := l[i] % MOD
		t_mod := t[i]
		sum_ts = (sum_ts + t_mod*s_mod) % MOD
		sum_ss = (sum_ss + s_mod*s_mod) % MOD
		sum_tl = (sum_tl + t_mod*l_mod) % MOD
		sum_sl = (sum_sl + s_mod*l_mod) % MOD
		sum_s = (sum_s + s_mod) % MOD
		sum_t = (sum_t + t_mod) % MOD
	}
	if n == 1 {
		ans := (s[0]%MOD)*sum_t%MOD + (l[0]%MOD)*sum_s%MOD
		ans %= MOD
		fmt.Println(ans)
		return
	}
	M := Matrix{sum_ts, sum_ss, sum_tl, sum_sl}
	p := n - 1
	res := Matrix{1, 0, 0, 1}
	for p > 0 {
		if p&1 == 1 {
			res = mulMat(res, M)
		}
		M = mulMat(M, M)
		p >>= 1
	}
	a, b := mulVec(res, s[0]%MOD, l[0]%MOD)
	ans := (a*sum_t + b*sum_s) % MOD
	fmt.Println(ans)
}
