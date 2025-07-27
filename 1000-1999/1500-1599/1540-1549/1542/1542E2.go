package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution to Codeforces problem 1542E2 - hard version of Abnormal Permutation Pairs.
// Counts the number of permutation pairs (p,q) of 1..n where p is lexicographically
// smaller than q and has strictly more inversions, modulo mod.
// The algorithm uses dynamic programming on permutation inversion differences.

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var mod int64
	if _, err := fmt.Fscan(in, &n, &mod); err != nil {
		return
	}

	if mod == 1 {
		fmt.Fprintln(out, 0)
		return
	}

	g := computeG(n, mod)

	ans := int64(0)
	prefix := int64(1)
	for j := 1; j <= n; j++ {
		if j > 1 {
			prefix = prefix * int64(n-j+2) % mod
		}
		m := n - j + 1
		L := m - 1
		for delta := 1; delta <= m-1; delta++ {
			ans = (ans + prefix*int64(m-delta)%mod*g[L][delta]) % mod
		}
	}
	fmt.Fprintln(out, ans%mod)
}

func computeG(n int, mod int64) [][]int64 {
	g := make([][]int64, n+1)
	g[0] = []int64{0}

	arrPrev := []int64{1}
	basePrev := 0

	for L := 1; L <= n; L++ {
		maxCurr := L * (L - 1) / 2
		baseCurr := maxCurr
		arrCurr := make([]int64, 2*maxCurr+1)

		sizePrev := len(arrPrev)
		prefix0 := make([]int64, sizePrev+1)
		prefix1 := make([]int64, sizePrev+1)
		for i := 0; i < sizePrev; i++ {
			prefix0[i+1] = (prefix0[i] + arrPrev[i]) % mod
			prefix1[i+1] = (prefix1[i] + int64(i-basePrev)*arrPrev[i]) % mod
		}
		getSum := func(a, b int) (int64, int64) {
			if a < -basePrev {
				a = -basePrev
			}
			if b > basePrev {
				b = basePrev
			}
			if a > b {
				return 0, 0
			}
			idxA := a + basePrev
			idxB := b + basePrev
			sum0 := (prefix0[idxB+1] - prefix0[idxA]) % mod
			if sum0 < 0 {
				sum0 += mod
			}
			sum1 := (prefix1[idxB+1] - prefix1[idxA]) % mod
			if sum1 < 0 {
				sum1 += mod
			}
			return sum0, sum1
		}
		for d := -maxCurr; d <= maxCurr; d++ {
			s1, s2 := getSum(d-(L-1), d)
			s3, s4 := getSum(d+1, d+(L-1))
			val := (int64(L) - int64(d)) % mod * s1 % mod
			val = (val + s2) % mod
			val = (val + (int64(L)+int64(d))%mod*s3%mod) % mod
			val = (val - s4) % mod
			if val < 0 {
				val += mod
			}
			arrCurr[d+baseCurr] = val
		}
		gL := make([]int64, L+1)
		prefixPos := make([]int64, maxCurr+2)
		for d := maxCurr; d >= 0; d-- {
			prefixPos[d] = (prefixPos[d+1] + arrCurr[d+baseCurr]) % mod
		}
		for delta := 0; delta <= L; delta++ {
			if delta+1 <= maxCurr {
				gL[delta] = prefixPos[delta+1]
			} else {
				gL[delta] = 0
			}
		}
		g[L] = gL
		arrPrev = arrCurr
		basePrev = baseCurr
	}
	return g
}
