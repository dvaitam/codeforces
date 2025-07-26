package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 998244353
const K int = 11
const N int = 22

func make2D() [][]int {
	a := make([][]int, K+1)
	for i := range a {
		a[i] = make([]int, K+1)
	}
	return a
}

func computeDP(n int) [K + 1][K + 1]int {
	dpPrev := make([][][]int, 1)
	dpPrev[0] = make2D()
	dpPrev[0][0][0] = 1
	for pos := n - 1; pos >= 1; pos-- {
		limit := n - pos
		dpCurr := make([][][]int, limit+1)
		for i := 0; i <= limit; i++ {
			dpCurr[i] = make2D()
		}
		for last := 0; last < len(dpPrev); last++ {
			table := dpPrev[last]
			for inv := 0; inv <= K; inv++ {
				row := table[inv]
				for desc := 0; desc <= K; desc++ {
					val := row[desc]
					if val == 0 {
						continue
					}
					for valNew := 0; valNew <= limit; valNew++ {
						ni := inv + valNew
						if ni > K {
							continue
						}
						nd := desc
						if valNew > last {
							nd++
						}
						if nd > K {
							continue
						}
						dpCurr[valNew][ni][nd] = (dpCurr[valNew][ni][nd] + val) % MOD
					}
				}
			}
		}
		dpPrev = dpCurr
	}
	var res [K + 1][K + 1]int
	for _, table := range dpPrev {
		for inv := 0; inv <= K; inv++ {
			for desc := 0; desc <= K; desc++ {
				res[inv][desc] = (res[inv][desc] + table[inv][desc]) % MOD
			}
		}
	}
	return res
}

func modInverse(a int) int {
	return powMod(a, MOD-2)
}

func powMod(a, e int) int {
	res := 1
	base := a % MOD
	for e > 0 {
		if e&1 == 1 {
			res = res * base % MOD
		}
		base = base * base % MOD
		e >>= 1
	}
	return res
}

var dp [N + 1][K + 1][K + 1]int
var bases [K + 1][K + 1]int
var coeffs [K + 1][K + 1][]int
var invFact [K + 1]int

func precompute() {
	// factorials for small k
	fact := make([]int, K+1)
	fact[0] = 1
	for i := 1; i <= K; i++ {
		fact[i] = fact[i-1] * i % MOD
	}
	invFact[K] = modInverse(fact[K])
	for i := K; i >= 1; i-- {
		invFact[i-1] = invFact[i] * i % MOD
	}

	for n := 1; n <= N; n++ {
		res := computeDP(n)
		for k := 0; k <= K; k++ {
			for x := 0; x <= K; x++ {
				dp[n][k][x] = res[k][x]
			}
		}
	}
	// Precompute polynomial coefficients
	for k := 0; k <= K; k++ {
		for x := 0; x <= K && x <= k; x++ {
			base := N - x
			bases[k][x] = base
			y := make([]int, x+1)
			for i := 0; i <= x; i++ {
				y[i] = dp[base+i][k][x]
			}
			coeff := make([]int, x+1)
			cur := make([]int, len(y))
			copy(cur, y)
			for i := 0; i <= x; i++ {
				coeff[i] = (cur[0]%MOD + MOD) % MOD
				if i == x {
					break
				}
				next := make([]int, len(cur)-1)
				for j := 0; j < len(cur)-1; j++ {
					val := cur[j+1] - cur[j]
					if val < 0 {
						val += MOD
					}
					next[j] = val % MOD
				}
				cur = next
			}
			coeffs[k][x] = coeff
		}
	}
}

func choose(n, r int) int {
	if r < 0 || r > n {
		return 0
	}
	res := 1
	for i := 0; i < r; i++ {
		res = res * ((n - i) % MOD) % MOD
	}
	res = res * invFact[r] % MOD
	return res
}

func getAnswer(n, k, x int) int {
	if x > k {
		return 0
	}
	if n <= N {
		return dp[n][k][x] % MOD
	}
	base := bases[k][x]
	diff := coeffs[k][x]
	m := n - base
	res := 0
	for i, c := range diff {
		res = (res + c*choose(m, i)%MOD) % MOD
	}
	return res
}

func main() {
	precompute()
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, k, x int
		fmt.Fscan(reader, &n, &k, &x)
		ans := getAnswer(n, k, x)
		fmt.Fprintln(writer, ans)
	}
}
