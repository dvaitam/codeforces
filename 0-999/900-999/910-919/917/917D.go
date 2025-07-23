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

func modInv(a int64) int64 { return modPow((a%MOD+MOD)%MOD, MOD-2) }

func det(mat [][]int64) int64 {
	n := len(mat)
	res := int64(1)
	for i := 0; i < n; i++ {
		pivot := -1
		for j := i; j < n; j++ {
			if mat[j][i] != 0 {
				pivot = j
				break
			}
		}
		if pivot == -1 {
			return 0
		}
		if pivot != i {
			mat[i], mat[pivot] = mat[pivot], mat[i]
			res = (MOD - res) % MOD
		}
		res = res * mat[i][i] % MOD
		inv := modInv(mat[i][i])
		for j := i + 1; j < n; j++ {
			if mat[j][i] == 0 {
				continue
			}
			factor := mat[j][i] * inv % MOD
			for k := i; k < n; k++ {
				mat[j][k] = (mat[j][k] - factor*mat[i][k]) % MOD
				if mat[j][k] < 0 {
					mat[j][k] += MOD
				}
			}
		}
	}
	return res % MOD
}

func polyMul(a, b []int64) []int64 {
	res := make([]int64, len(a)+len(b)-1)
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b); j++ {
			res[i+j] = (res[i+j] + a[i]*b[j]) % MOD
		}
	}
	return res
}

func polyDiv(poly []int64, r int64) []int64 {
	n := len(poly) - 1
	res := make([]int64, n)
	res[n-1] = poly[n]
	for i := n - 1; i >= 1; i-- {
		res[i-1] = (poly[i] + r*res[i]) % MOD
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	edges := make([][]bool, n)
	for i := 0; i < n; i++ {
		edges[i] = make([]bool, n)
	}
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		edges[u][v] = true
		edges[v][u] = true
	}

	// precompute values P(x) for x = 0..n-1
	vals := make([]int64, n)
	for x := 0; x < n; x++ {
		// build Laplacian
		lap := make([][]int64, n)
		for i := range lap {
			lap[i] = make([]int64, n)
		}
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				w := int64(1)
				if edges[i][j] {
					w = int64(x)
				}
				lap[i][i] = (lap[i][i] + w) % MOD
				lap[j][j] = (lap[j][j] + w) % MOD
				lap[i][j] = (lap[i][j] - w) % MOD
				lap[j][i] = (lap[j][i] - w) % MOD
			}
		}
		// remove last row and column
		m := n - 1
		mat := make([][]int64, m)
		for i := 0; i < m; i++ {
			mat[i] = make([]int64, m)
			for j := 0; j < m; j++ {
				v := lap[i][j]
				if v < 0 {
					v += MOD
				}
				mat[i][j] = v
			}
		}
		vals[x] = det(mat)
	}

	// polynomial interpolation with points x=0..n-1
	fact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}

	// build total polynomial prod (x - i)
	poly := []int64{1}
	for i := 0; i < n; i++ {
		term := []int64{(MOD - int64(i)%MOD) % MOD, 1}
		poly = polyMul(poly, term)
	}

	res := make([]int64, n)
	for i := 0; i < n; i++ {
		denom := fact[i] * fact[n-1-i] % MOD
		if (n-1-i)%2 == 1 {
			denom = (MOD - denom) % MOD
		}
		invDenom := modInv(denom)
		q := polyDiv(poly, int64(i))
		for j := 0; j < n; j++ {
			res[j] = (res[j] + vals[i]*q[j]%MOD*invDenom) % MOD
		}
	}

	out := bufio.NewWriter(os.Stdout)
	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, res[i])
	}
	fmt.Fprintln(out)
	out.Flush()
}
