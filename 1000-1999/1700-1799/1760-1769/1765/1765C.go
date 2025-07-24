package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}

	maxM := k
	if maxM > 4*n-1 {
		maxM = 4*n - 1
	}
	maxX := k / 4
	if maxX > n {
		maxX = n
	}

	// factorials for combinations
	fact := make([]int64, 4*n+1)
	invFact := make([]int64, 4*n+1)
	fact[0] = 1
	for i := 1; i <= 4*n; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[4*n] = modPow(fact[4*n], MOD-2)
	for i := 4 * n; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}
	comb := func(a, b int) int64 {
		if b < 0 || b > a {
			return 0
		}
		return fact[a] * invFact[b] % MOD * invFact[a-b] % MOD
	}

	// DP for coefficients of B_x^i
	b1 := []int64{1}
	b2 := []int64{1}
	b3 := []int64{1}
	b4 := []int64{1}

	F := make([][]int64, maxX+1)
	for i := range F {
		F[i] = make([]int64, maxM+1)
	}

	for x := n; x >= 0; x-- {
		if x <= maxX {
			for m := 4 * x; m <= maxM && m-4*x < len(b4); m++ {
				F[x][m] = b4[m-4*x]
			}
		}
		if x == 0 {
			break
		}
		c := comb(n, x-1)
		c2 := c * c % MOD
		c3 := c2 * c % MOD
		c4 := c3 * c % MOD
		maxLen := maxM - 4*(x-1)

		nb1 := make([]int64, min(len(b1)+1, maxLen+1))
		if len(nb1) > 0 {
			nb1[0] = c
			for i := 1; i < len(nb1); i++ {
				nb1[i] = b1[i-1]
			}
		}

		nb2 := make([]int64, min(len(b2)+2, maxLen+1))
		if len(nb2) > 0 {
			nb2[0] = c2
			for i := 1; i < len(nb2); i++ {
				if i-1 < len(b1) {
					nb2[i] = (nb2[i] + 2*c%MOD*b1[i-1]) % MOD
				}
				if i-2 >= 0 && i-2 < len(b2) {
					nb2[i] = (nb2[i] + b2[i-2]) % MOD
				}
			}
		}

		nb3 := make([]int64, min(len(b3)+3, maxLen+1))
		if len(nb3) > 0 {
			nb3[0] = c3
			for i := 1; i < len(nb3); i++ {
				val := int64(0)
				if i-1 < len(b1) {
					val = (val + 3*c2%MOD*b1[i-1]) % MOD
				}
				if i-2 >= 0 && i-2 < len(b2) {
					val = (val + 3*c%MOD*b2[i-2]) % MOD
				}
				if i-3 >= 0 && i-3 < len(b3) {
					val = (val + b3[i-3]) % MOD
				}
				nb3[i] = val
			}
		}

		nb4 := make([]int64, min(len(b4)+4, maxLen+1))
		if len(nb4) > 0 {
			nb4[0] = c4
			for i := 1; i < len(nb4); i++ {
				val := int64(0)
				if i-1 < len(b1) {
					val = (val + 4*c3%MOD*b1[i-1]) % MOD
				}
				if i-2 >= 0 && i-2 < len(b2) {
					val = (val + 6*c2%MOD*b2[i-2]) % MOD
				}
				if i-3 >= 0 && i-3 < len(b3) {
					val = (val + 4*c%MOD*b3[i-3]) % MOD
				}
				if i-4 >= 0 && i-4 < len(b4) {
					val = (val + b4[i-4]) % MOD
				}
				nb4[i] = val % MOD
			}
		}

		b1, b2, b3, b4 = nb1, nb2, nb3, nb4
	}

	invDen := make([]int64, maxM+1)
	for m := 0; m <= maxM; m++ {
		invDen[m] = modPow(comb(4*n, m), MOD-2)
	}

	pCorr := make([]int64, maxM+1)
	for m := 0; m <= maxM; m++ {
		sumProb := int64(0)
		for x := 1; x <= maxX && 4*x <= m; x++ {
			sumProb = (sumProb + F[x][m]*invDen[m]) % MOD
		}
		val := (int64(n)%MOD - sumProb + MOD) % MOD
		pCorr[m] = val * modPow(int64(4*n-m), MOD-2) % MOD
	}

	total := int64(0)
	if k >= 4*n {
		for m := 0; m <= 4*n-1; m++ {
			total = (total + pCorr[m]) % MOD
		}
	} else {
		for m := 0; m <= k; m++ {
			total = (total + pCorr[m]) % MOD
		}
		remain := 4*n - k - 1
		if remain > 0 {
			total = (total + pCorr[k]*int64(remain)) % MOD
		}
	}

	fmt.Fprintln(writer, total%MOD)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
