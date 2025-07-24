package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const INF int64 = 1e18

var digits = []int64{2, 3, 4, 5, 6, 7, 8, 9}

func comb(n, k int64) int64 {
	if k < 0 || k > n {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	res := int64(1)
	for i := int64(1); i <= k; i++ {
		num := n - k + i
		if res > INF/num {
			return INF
		}
		res = res * num / i
		if res > INF {
			return INF
		}
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var m, k int64
	if _, err := fmt.Fscan(reader, &m, &k); err != nil {
		return
	}

	primes := []int64{2, 3, 5, 7}
	exps := make([]int, 4)
	temp := m
	for i, p := range primes {
		for temp%p == 0 {
			exps[i]++
			temp /= p
		}
	}
	if temp != 1 {
		fmt.Fprintln(writer, -1)
		return
	}

	// generate all divisors of m
	divisors := []int64{}
	var gen func(int, int64)
	gen = func(idx int, cur int64) {
		if idx == 4 {
			divisors = append(divisors, cur)
			return
		}
		p := primes[idx]
		val := int64(1)
		for e := 0; e <= exps[idx]; e++ {
			gen(idx+1, cur*val)
			val *= p
		}
	}
	gen(0, 1)
	sort.Slice(divisors, func(i, j int) bool { return divisors[i] < divisors[j] })
	idx := make(map[int64]int)
	for i, d := range divisors {
		idx[d] = i
	}

	D := len(divisors)
	T := exps[0] + exps[1] + exps[2] + exps[3]
	g := make([][]int64, T+1)
	for i := range g {
		g[i] = make([]int64, D)
	}
	g[0][idx[1]] = 1
	for t := 1; t <= T; t++ {
		for i, div := range divisors {
			var total int64
			for _, d := range digits {
				if div%d == 0 {
					total += g[t-1][idx[div/d]]
					if total > INF {
						total = INF
					}
				}
			}
			g[t][i] = total
		}
	}

	mIdx := idx[m]

	countLen := func(L int64) int64 {
		var tot int64
		for t := 1; t <= T && int64(t) <= L; t++ {
			val := g[t][mIdx]
			if val == 0 {
				continue
			}
			cmb := comb(L, int64(t))
			if cmb == 0 {
				continue
			}
			if cmb == INF || val > INF/cmb {
				tot = INF
				break
			} else {
				tot += cmb * val
				if tot > INF {
					tot = INF
					break
				}
			}
		}
		return tot
	}

	// find length containing k-th number
	var L int64 = 1
	for {
		cnt := countLen(L)
		if k > cnt {
			k -= cnt
			L++
			if L > 200000 { // safety limit
				break
			}
		} else {
			break
		}
	}

	// precompute function to count sequences of length len with product rem
	var countRest func(int64, int64) int64
	countRest = func(rem, len int64) int64 {
		if rem == 1 && len == 0 {
			return 1
		}
		if len == 0 {
			if rem == 1 {
				return 1
			} // though should be handled above
			return 0
		}
		idxRem := idx[rem]
		var tot int64
		for t := int64(0); t <= int64(T) && t <= len; t++ {
			val := g[t][idxRem]
			if val == 0 {
				continue
			}
			cmb := comb(len, t)
			if cmb == 0 {
				continue
			}
			if cmb == INF || val > INF/cmb {
				tot = INF
				break
			} else {
				tot += cmb * val
				if tot > INF {
					tot = INF
					break
				}
			}
		}
		return tot
	}

	// build k-th number of length L
	prefixMul := int64(1)
	var result []byte
	for pos := int64(0); pos < L; pos++ {
		for d := int64(1); d <= 9; d++ {
			newMul := prefixMul
			if d > 1 {
				if m%(prefixMul*d) != 0 {
					continue
				}
				newMul = prefixMul * d
			}
			rem := m / newMul
			cnt := countRest(rem, L-pos-1)
			if k > cnt {
				k -= cnt
			} else {
				result = append(result, byte('0'+d))
				prefixMul = newMul
				break
			}
		}
	}

	fmt.Fprintln(writer, string(result))
}
