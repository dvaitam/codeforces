package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD = 1000000007

var fact []int64
var invFact []int64

func power(base, exp int64) int64 {
	res := int64(1)
	base %= MOD
	for exp > 0 {
		if exp%2 == 1 {
			res = (res * base) % MOD
		}
		base = (base * base) % MOD
		exp /= 2
	}
	return res
}

func precompute(maxN int) {
	fact = make([]int64, maxN+1)
	invFact = make([]int64, maxN+1)
	fact[0] = 1
	invFact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = (fact[i-1] * int64(i)) % MOD
	}
	invFact[maxN] = power(fact[maxN], MOD-2)
	for i := maxN - 1; i >= 1; i-- {
		invFact[i] = (invFact[i+1] * int64(i+1)) % MOD
	}
}

func nCr(n, r int) int64 {
	if r < 0 || r > n {
		return 0
	}
	num := fact[n]
	den := (invFact[r] * invFact[n-r]) % MOD
	return (num * den) % MOD
}

func F(r1, c1, r2, c2 int) int64 {
	if r1 > r2 || c1 > c2 {
		return 0
	}
	dr := r2 - r1
	dc := c2 - c1
	return nCr(dr+dc, dr)
}

type Point struct {
	r, c int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m, k, s int
	if _, err := fmt.Fscan(reader, &n, &m, &k, &s); err != nil {
		return
	}

	precompute(n + m)

	V := make([]Point, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(reader, &V[i].r, &V[i].c)
	}

	sort.Slice(V, func(i, j int) bool {
		if V[i].r != V[j].r {
			return V[i].r < V[j].r
		}
		return V[i].c < V[j].c
	})

	dp := make([][]int64, k)
	for i := 0; i < k; i++ {
		dp[i] = make([]int64, 21)
	}

	T := make([]int64, 21)
	T[0] = F(1, 1, n, m)

	for c := 1; c <= 20; c++ {
		for i := 0; i < k; i++ {
			var waysC int64 = 0
			if c == 1 {
				waysC = F(1, 1, V[i].r, V[i].c)
			} else {
				for x := 0; x < i; x++ {
					if dp[x][c-1] > 0 {
						waysC = (waysC + dp[x][c-1]*F(V[x].r, V[x].c, V[i].r, V[i].c)) % MOD
					}
				}
			}

			var sumDp int64 = 0
			for x := 0; x < i; x++ {
				if dp[x][c] > 0 {
					sumDp = (sumDp + dp[x][c]*F(V[x].r, V[x].c, V[i].r, V[i].c)) % MOD
				}
			}

			dp[i][c] = (waysC - sumDp) % MOD
			if dp[i][c] < 0 {
				dp[i][c] += MOD
			}
		}

		var tc int64 = 0
		for x := 0; x < k; x++ {
			if dp[x][c] > 0 {
				tc = (tc + dp[x][c]*F(V[x].r, V[x].c, n, m)) % MOD
			}
		}
		T[c] = tc
	}

	charges := make([]int64, 21)
	curr := int64(s)
	charges[0] = curr % MOD
	for c := 1; c <= 20; c++ {
		curr = curr - (curr / 2)
		charges[c] = curr % MOD
	}

	var expectedNum int64 = 0
	for c := 0; c < 20; c++ {
		exact := (T[c] - T[c+1]) % MOD
		if exact < 0 {
			exact += MOD
		}
		term := (exact * charges[c]) % MOD
		expectedNum = (expectedNum + term) % MOD
	}
	expectedNum = (expectedNum + T[20]*charges[20]) % MOD

	ans := (expectedNum * power(T[0], MOD-2)) % MOD
	fmt.Println(ans)
}
