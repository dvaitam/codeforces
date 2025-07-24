package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func modPow(a, e int64) int64 {
	res := int64(1)
	a %= mod
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	w := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &w[i])
	}

	m := n / 3 // number of triangles
	ans := int64(1)
	for i := 0; i < m; i++ {
		a := w[3*i]
		b := w[3*i+1]
		c := w[3*i+2]
		s1 := a + b
		s2 := a + c
		s3 := b + c
		mx := s1
		if s2 > mx {
			mx = s2
		}
		if s3 > mx {
			mx = s3
		}
		cnt := 0
		if s1 == mx {
			cnt++
		}
		if s2 == mx {
			cnt++
		}
		if s3 == mx {
			cnt++
		}
		ans = ans * int64(cnt) % mod
	}

	// compute C(m, m/2)
	fact := make([]int64, m+1)
	fact[0] = 1
	for i := 1; i <= m; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact := make([]int64, m+1)
	invFact[m] = modPow(fact[m], mod-2)
	for i := m; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}
	half := m / 2
	comb := fact[m] * invFact[half] % mod * invFact[m-half] % mod
	ans = ans * comb % mod

	fmt.Fprintln(out, ans)
}
