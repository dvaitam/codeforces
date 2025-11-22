package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	mod      int64 = 998244353
	maxLimit int   = 1_000_000 + 5
)

var fact [maxLimit]int64
var invFact [maxLimit]int64
var pow2 [maxLimit]int64

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

func initPrecalc() {
	fact[0] = 1
	pow2[0] = 1
	for i := 1; i < maxLimit; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
		pow2[i] = pow2[i-1] * 2 % mod
	}
	invFact[maxLimit-1] = modPow(fact[maxLimit-1], mod-2)
	for i := maxLimit - 1; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}
}

func comb(n, k int) int64 {
	if k < 0 || k > n || n < 0 {
		return 0
	}
	return fact[n] * invFact[k] % mod * invFact[n-k] % mod
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	initPrecalc()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)

		cntL := make([]int, n+2)   // l = i
		cntR := make([]int, n+2)   // r = i
		prefRcnt := make([]int, n+2)
		prefRsum := make([]int64, n+2)
		suffLcnt := make([]int, n+3)
		suffLsum := make([]int64, n+3)
		prefixB := make([]int64, n+1)
		prefixB1 := make([]int64, n+1)
		prefixC := make([]int64, n+1)

		for i := 0; i < n; i++ {
			var l, r int
			fmt.Fscan(in, &l, &r)
			cntL[l]++
			cntR[r]++
		}

		for i := 1; i <= n; i++ {
			prefRcnt[i] = prefRcnt[i-1] + cntR[i]
			prefRsum[i] = prefRsum[i-1] + int64(cntR[i])*int64(i)
		}
		for i := n; i >= 1; i-- {
			suffLcnt[i] = suffLcnt[i+1] + cntL[i]
			suffLsum[i] = suffLsum[i+1] + int64(cntL[i])*int64(i)
		}

		var ans int64
		for x := 1; x <= n; x++ {
			// category sizes
			a := 0
			if x >= 2 {
				a = prefRcnt[x-2] // r <= x-2
			}
			b := 0
			if x >= 1 {
				b = cntR[x-1] // r = x-1
			}
			c := cntL[x] // l = x
			d := 0
			if x+1 <= n {
				d = suffLcnt[x+1] // l >= x+1
			}
			e := n - (a + b + c + d) // middle intervals, free choices

			powE := pow2[e]
			powB := pow2[b]
			powC := pow2[c]

			// prefix sums of combinations for B (r = x-1)
			for i := 0; i <= b; i++ {
				val := comb(b, i)
				if i == 0 {
					prefixB[i] = val
				} else {
					prefixB[i] = (prefixB[i-1] + val) % mod
				}
			}
			if b > 0 {
				for i := 0; i <= b-1; i++ {
					val := comb(b-1, i)
					if i == 0 {
						prefixB1[i] = val
					} else {
						prefixB1[i] = (prefixB1[i-1] + val) % mod
					}
				}
			}

			// prefix sums of combinations for C (l = x)
			for i := 0; i <= c; i++ {
				val := comb(c, i)
				if i == 0 {
					prefixC[i] = val
				} else {
					prefixC[i] = (prefixC[i-1] + val) % mod
				}
			}

			sumCostA := int64(0)
			if a > 0 {
				sumCostA = int64(a)*int64(x) - prefRsum[x-2]
			}
			sumCostD := int64(0)
			if d > 0 && x+1 <= n {
				sumCostD = suffLsum[x+1] - int64(d)*int64(x)
			}

			coeffA := int64(0)
			coeffD := int64(0)
			coeffB := int64(0)

			// k ranges where conditions can hold: k in [-c, b]
			minK := -c
			if minK < -a {
				minK = -a
			}
			maxK := b
			if maxK > d {
				maxK = d
			}

			for k := minK; k <= maxK; k++ {
				// S_b
				sb := int64(0)
				if k <= 0 {
					sb = powB
				} else {
					if k-1 <= b-1 {
						sb = (powB - prefixB[k-1] + mod) % mod
					} else {
						continue
					}
				}
				if sb == 0 {
					continue
				}

				// S_c
				sc := int64(0)
				if k > 0 {
					sc = powC
				} else {
					idx := -k
					if idx <= c {
						sc = (powC - prefixC[idx] + mod) % mod
					} else {
						continue
					}
				}
				if sc == 0 {
					continue
				}

				common := sb * sc % mod
				base := comb(a+d, a+k)
					if base != 0 {
						if b > 0 {
							// subsets containing a fixed interval of B
							sbIncl := int64(0)
							if k <= 0 {
							sbIncl = pow2[b-1]
						} else {
							if k-2 >= 0 && k-2 < len(prefixB1) {
								sbIncl = (pow2[b-1] - prefixB1[k-2] + mod) % mod
							} else if k-1 <= b-1 {
								sbIncl = pow2[b-1]
							} else {
								sbIncl = 0
							}
						}
						if sbIncl != 0 {
							coeffB = (coeffB + base*sbIncl%mod*sc%mod) % mod
						}
					}
					if a > 0 {
						na := comb(a+d-1, a+k)
						if na != 0 {
							coeffA = (coeffA + na*common) % mod
						}
					}
					if d > 0 {
						nd := comb(a+d-1, a+k-1)
						if nd != 0 {
							coeffD = (coeffD + nd*common) % mod
						}
					}
				}
			}

			part := (sumCostA%mod*coeffA%mod + sumCostD%mod*coeffD%mod + int64(b)%mod*coeffB%mod) % mod
			part = part * powE % mod
			ans += part
			if ans >= mod {
				ans -= mod
			}
		}

		fmt.Fprintln(out, ans%mod)
	}
}
