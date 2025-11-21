package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const mod int64 = 998244353

func modPow(a, e int64) int64 {
	res := int64(1)
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
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	prob := make([][]int64, n)
	for i := 0; i < n; i++ {
		prob[i] = make([]int64, n)
	}

	for i := 0; i < m; i++ {
		var u, v int
		var p, q int64
		fmt.Fscan(in, &u, &v, &p, &q)
		u--
		v--
		val := p * modPow(q, mod-2) % mod
		prob[u][v] = val
		prob[v][u] = val
	}

	factor := make([][]int64, n)
	invFactor := make([][]int64, n)
	for i := 0; i < n; i++ {
		factor[i] = make([]int64, n)
		invFactor[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			f := (1 - prob[i][j] + mod) % mod
			factor[i][j] = f
			invFactor[i][j] = modPow(f, mod-2)
		}
	}

	size := 1 << n
	prod := make([][]int64, n)
	invProd := make([][]int64, n)
	for v := 0; v < n; v++ {
		prod[v] = make([]int64, size)
		invProd[v] = make([]int64, size)
		prod[v][0] = 1
		invProd[v][0] = 1
		for mask := 1; mask < size; mask++ {
			lsb := bits.TrailingZeros(uint(mask))
			prev := mask & (mask - 1)
			prod[v][mask] = prod[v][prev] * factor[lsb][v] % mod
			invProd[v][mask] = invProd[v][prev] * invFactor[lsb][v] % mod
		}
	}

	fullMask := size - 1
	E := make([]int64, size)
	F := make([]int64, size)
	E[fullMask] = 0
	F[fullMask] = 0

	for mask := fullMask - 1; mask >= 0; mask-- {
		if mask == 0 {
			// Unreachable from the initial state; skipped to avoid division by zero when P0==1.
			E[mask] = 0
			F[mask] = 0
			continue
		}
		compMask := fullMask ^ mask
		if compMask == 0 {
			E[mask] = 0
			F[mask] = 0
			continue
		}
		P0 := int64(1)
		for v := 0; v < n; v++ {
			if (mask>>v)&1 == 0 {
				P0 = P0 * prod[v][mask] % mod
			}
		}
		oneMinusP0 := (1 - P0 + mod) % mod
		invOneMinus := modPow(oneMinusP0, mod-2)

		sumAlphaF := int64(0)
		for v := 0; v < n; v++ {
			if (mask>>v)&1 != 0 {
				continue
			}
			alpha := (1 - prod[v][mask] + mod) % mod
			alpha = alpha * invProd[v][mask] % mod
			sumAlphaF = (sumAlphaF + alpha*F[mask|(1<<v)]) % mod
		}

		E[mask] = (1 + P0*sumAlphaF%mod) % mod
		E[mask] = E[mask] * invOneMinus % mod
		F[mask] = (E[mask] + sumAlphaF) % mod
	}

	startMask := 1 << 0
	fmt.Println(E[startMask] % mod)
}
