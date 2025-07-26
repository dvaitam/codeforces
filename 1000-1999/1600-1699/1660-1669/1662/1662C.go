package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

type Matrix [][]int64

func makeMatrix(n, m int) Matrix {
	mat := make(Matrix, n)
	for i := range mat {
		mat[i] = make([]int64, m)
	}
	return mat
}

func matMul(a, b Matrix) Matrix {
	n := len(a)
	m := len(a[0])
	p := len(b[0])
	res := makeMatrix(n, p)
	for i := 0; i < n; i++ {
		ai := a[i]
		for k := 0; k < m; k++ {
			if ai[k] == 0 {
				continue
			}
			bk := b[k]
			val := ai[k]
			for j := 0; j < p; j++ {
				res[i][j] = (res[i][j] + val*bk[j]) % MOD
			}
		}
	}
	return res
}

func matPow(base Matrix, exp int) Matrix {
	n := len(base)
	res := makeMatrix(n, n)
	for i := 0; i < n; i++ {
		res[i][i] = 1
	}
	for exp > 0 {
		if exp&1 == 1 {
			res = matMul(res, base)
		}
		base = matMul(base, base)
		exp >>= 1
	}
	return res
}

func trace(mat Matrix) int64 {
	n := len(mat)
	var s int64
	for i := 0; i < n; i++ {
		s = (s + mat[i][i]) % MOD
	}
	return s
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}
	A := makeMatrix(n, n)
	deg := make([]int, n)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		A[a][b] = 1
		A[b][a] = 1
		deg[a]++
		deg[b]++
	}

	if k == 1 {
		fmt.Println(0)
		return
	}

	// B1 = A
	B1 := makeMatrix(n, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			B1[i][j] = A[i][j]
		}
	}
	// B2 = A*A - D
	B2 := matMul(A, A)
	for i := 0; i < n; i++ {
		B2[i][i] = (B2[i][i] - int64(deg[i])%MOD + MOD) % MOD
	}
	if k == 2 {
		fmt.Println(trace(B2))
		return
	}

	size := 2 * n
	M := makeMatrix(size, size)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			M[i][j] = A[i][j]
		}
		// -(D - I) -> (MOD - (deg[i]-1)) % MOD
		M[i][n+i] = (MOD - int64(deg[i]-1)%MOD) % MOD
		M[n+i][i] = 1
	}

	Pow := matPow(M, k-2)

	X2 := makeMatrix(size, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			X2[i][j] = B2[i][j]
			X2[n+i][j] = B1[i][j]
		}
	}

	Y := matMul(Pow, X2)
	Bk := Y[:n]
	fmt.Println(trace(Bk))
}
