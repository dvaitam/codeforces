package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func polyMul(a, b []int64, d int, c []int) []int64 {
	tmp := make([]int64, len(a)+len(b)-1)
	for i, av := range a {
		if av == 0 {
			continue
		}
		for j, bv := range b {
			if bv == 0 {
				continue
			}
			tmp[i+j] = (tmp[i+j] + av*bv) % mod
		}
	}
	for len(tmp) > d {
		k := len(tmp) - 1
		coeff := tmp[k] % mod
		tmp = tmp[:k]
		if coeff != 0 {
			for _, cj := range c {
				idx := k - cj
				tmp[idx] = (tmp[idx] + coeff) % mod
			}
		}
	}
	if len(tmp) < d {
		out := make([]int64, d)
		copy(out, tmp)
		tmp = out
	}
	if len(tmp) > d {
		tmp = tmp[:d]
	}
	return tmp
}

func polyPow(base []int64, exp int64, d int, c []int) []int64 {
	res := []int64{1}
	b := base
	e := exp
	for e > 0 {
		if e&1 == 1 {
			res = polyMul(res, b, d, c)
		}
		b = polyMul(b, b, d, c)
		e >>= 1
	}
	if len(res) < d {
		out := make([]int64, d)
		copy(out, res)
		res = out
	}
	if len(res) > d {
		res = res[:d]
	}
	return res
}

func polyPowX(exp int64, d int, c []int) []int64 {
	base := []int64{0, 1}
	base = polyMul([]int64{1}, base, d, c)
	return polyPow(base, exp, d, c)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var F, B int
	var N int64
	if _, err := fmt.Fscan(in, &F, &B, &N); err != nil {
		return
	}
	p := make([]int64, F)
	for i := 0; i < F; i++ {
		fmt.Fscan(in, &p[i])
	}
	c := make([]int, B)
	d := 0
	for i := 0; i < B; i++ {
		fmt.Fscan(in, &c[i])
		if c[i] > d {
			d = c[i]
		}
	}
	P := make([]int64, d)
	for _, pi := range p {
		poly := polyPowX(pi, d, c)
		for i := 0; i < d; i++ {
			P[i] = (P[i] + poly[i]) % mod
		}
	}
	res := polyPow(P, N, d, c)
	Bseq := make([]int64, d)
	Bseq[0] = 1
	for t := 1; t < d; t++ {
		var val int64
		for _, cj := range c {
			if t-cj >= 0 {
				val = (val + Bseq[t-cj]) % mod
			}
		}
		Bseq[t] = val
	}
	var ans int64
	for i := 0; i < d; i++ {
		ans = (ans + res[i]*Bseq[i]) % mod
	}
	fmt.Println(ans % mod)
}
