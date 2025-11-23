package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD = 998244353
const MAX = 1000005

var fact [MAX]int64
var invFact [MAX]int64

func power(a, b int64) int64 {
	res := int64(1)
	a %= MOD
	for b > 0 {
		if b%2 == 1 {
			res = (res * a) % MOD
		}
		a = (a * a) % MOD
		b /= 2
	}
	return res
}

func initComb() {
	fact[0] = 1
	for i := 1; i < MAX; i++ {
		fact[i] = (fact[i-1] * int64(i)) % MOD
	}

	invFact[MAX-1] = power(fact[MAX-1], MOD-2)
	for i := MAX - 2; i >= 0; i-- {
		invFact[i] = (invFact[i+1] * int64(i+1)) % MOD
	}
}

func nCr(n, r int) int64 {
	if r < 0 || r > n {
		return 0
	}
	return fact[n] * invFact[r] % MOD * invFact[n-r] % MOD
}

func main() {
	initComb()
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)

	for i := 0; i < t; i++ {
		var l, n int
		fmt.Fscan(reader, &l, &n)

		S := l - 2*n
		
		// Total configurations for one fixed ordering (e.g., FJ's cow first)
		// is equivalent to choosing 2n positions from l: C(l, 2n)
		total := nCr(l, 2*n)

		// We subtract configurations where FJ loses.
		// FJ loses if all gaps between adjacent opponents (FJ-FN or FN-FJ pairs) are even.
		// Let these gaps be x_1, x_2, ..., x_n. We require x_i = 2*k_i.
		// Let other gaps be y_0, ..., y_n.
		// Sum(x_i) + Sum(y_j) = S.
		// Sum(2*k_i) + Sum(y_j) = S.
		// Let K = Sum(k_i). Then 2K + Sum(y_j) = S.
		// Iterate over possible values of K.

		bad := int64(0)
		maxK := S / 2
		for K := 0; K <= maxK; K++ {
			// Ways to choose k_i such that Sum(k_i) = K. 
			// Stars and bars with n variables: C(K + n - 1, n - 1)
			waysK := nCr(K+n-1, n-1)
			
			// Ways to choose y_j such that Sum(y_j) = S - 2K.
			// Stars and bars with n+1 variables: C((S - 2K) + (n + 1) - 1, (n + 1) - 1)
			// = C(S - 2K + n, n)
			waysY := nCr(S-2*K+n, n)
			
			term := (waysK * waysY) % MOD
			bad = (bad + term) % MOD
		}

		// There are two orderings: FJ-FN-... and FN-FJ-...
		// The count is identical for both.
		ans := (total - bad + MOD) % MOD
		ans = (ans * 2) % MOD
		fmt.Fprintln(writer, ans)
	}
}