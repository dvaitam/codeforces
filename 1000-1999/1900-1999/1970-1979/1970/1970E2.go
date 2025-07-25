package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

// matMul multiplies two matrices of size m x m under modulo.
func matMul(a, b [][]int64) [][]int64 {
	m := len(a)
	c := make([][]int64, m)
	for i := 0; i < m; i++ {
		c[i] = make([]int64, m)
		ai := a[i]
		ci := c[i]
		for k := 0; k < m; k++ {
			if ai[k] == 0 {
				continue
			}
			aik := ai[k]
			bk := b[k]
			for j := 0; j < m; j++ {
				ci[j] = (ci[j] + aik*bk[j]) % mod
			}
		}
	}
	return c
}

// vecMul multiplies row vector v (size m) by matrix a (m x m).
func vecMul(v []int64, a [][]int64) []int64 {
	m := len(v)
	res := make([]int64, m)
	for i := 0; i < m; i++ {
		if v[i] == 0 {
			continue
		}
		vi := v[i]
		ai := a[i]
		for j := 0; j < m; j++ {
			res[j] = (res[j] + vi*ai[j]) % mod
		}
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var m int
	var n int64
	if _, err := fmt.Fscan(in, &m, &n); err != nil {
		return
	}
	s := make([]int64, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &s[i])
	}
	l := make([]int64, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &l[i])
	}

	total := make([]int64, m)
	for i := 0; i < m; i++ {
		total[i] = s[i] + l[i]
	}

	A := make([][]int64, m)
	for i := 0; i < m; i++ {
		row := make([]int64, m)
		si := s[i]
		li := l[i]
		for j := 0; j < m; j++ {
			row[j] = (si*total[j] + li*s[j]) % mod
		}
		A[i] = row
	}

	// res represents row vector after applying powers of A.
	res := make([]int64, m)
	res[0] = 1 // start at cabin 1 (index 0)

	for n > 0 {
		if n&1 == 1 {
			res = vecMul(res, A)
		}
		n >>= 1
		if n > 0 {
			A = matMul(A, A)
		}
	}

	var ans int64
	for _, v := range res {
		ans = (ans + v) % mod
	}
	fmt.Fprintln(out, ans)
}
