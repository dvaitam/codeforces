package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int64 = 998244353
const maxN int = 5000
const maxFact int = 10000 // enough for 2 * maxN

var fact [maxFact + 1]int64
var invFact [maxFact + 1]int64

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

func initFact() {
	fact[0] = 1
	for i := 1; i <= maxFact; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[maxFact] = modPow(fact[maxFact], mod-2)
	for i := maxFact; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}
}

func comb(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * invFact[k] % mod * invFact[n-k] % mod
}

func main() {
	in := bufio.NewReader(os.Stdin)
	initFact()

	var T int
	fmt.Fscan(in, &T)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		l := make([]int, n)
		r := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &l[i], &r[i])
		}

		orderL := make([]int, n)
		orderR := make([]int, n)
		for i := 0; i < n; i++ {
			orderL[i] = i
			orderR[i] = i
		}

		sort.Slice(orderL, func(i, j int) bool {
			li, lj := l[orderL[i]], l[orderL[j]]
			if li != lj {
				return li > lj
			}
			return orderL[i] < orderL[j]
		})

		sort.Slice(orderR, func(i, j int) bool {
			ri, rj := r[orderR[i]], r[orderR[j]]
			if ri != rj {
				return ri < rj
			}
			return orderR[i] < orderR[j]
		})

		posL := make([]int, n)
		posR := make([]int, n)
		for idx, v := range orderL {
			posL[v] = idx + 1
		}
		for idx, v := range orderR {
			posR[v] = idx + 1
		}

		// pref[x][y] = #elements with posL <= x and posR <= y
		pref := make([][]int, n+1)
		for i := 0; i <= n; i++ {
			pref[i] = make([]int, n+1)
		}
		for i := 0; i < n; i++ {
			x, y := posL[i], posR[i]
			pref[x][y]++
		}
		for i := 1; i <= n; i++ {
			row := pref[i]
			prev := pref[i-1]
			for j := 1; j <= n; j++ {
				row[j] += row[j-1] + prev[j] - prev[j-1]
			}
		}

		// powers of two up to n (since F + C <= n-1 at most)
		pow2 := make([]int64, n+1)
		pow2[0] = 1
		for i := 1; i <= n; i++ {
			pow2[i] = pow2[i-1] * 2 % mod
		}

		var ans int64
		for i := 0; i < n; i++ {
			pL := posL[i]
			pRofI := posR[i]
			for j := 0; j < n; j++ {
				diff := l[i] - r[j]
				if diff <= 0 {
					continue
				}
				pR := posR[j]
				addL := 0
				if posL[j] < pL {
					addL = 1
				}
				addR := 0
				if pRofI < pR {
					addR = 1
				}

				cBoth := pref[pL-1][pR-1]
				dOnlyL := (pL - 1) - cBoth - addL
				eOnlyR := (pR - 1) - cBoth - addR
				if dOnlyL < 0 || eOnlyR < 0 {
					continue
				}

				totalOthers := n - 2
				if i == j {
					totalOthers = n - 1
				}
				fFree := totalOthers - cBoth - dOnlyL - eOnlyR
				if fFree < 0 {
					continue
				}

				delta := addR - addL
				nC := dOnlyL + eOnlyR
				kC := eOnlyR + delta
				if kC < 0 || kC > nC {
					continue
				}

				ways := pow2[fFree+cBoth] * comb(nC, kC) % mod
				contrib := int64(diff%int(mod)) * ways % mod
				ans += contrib
				if ans >= mod {
					ans -= mod
				}
			}
		}

		fmt.Fprintln(writer, ans%mod)
	}
}
