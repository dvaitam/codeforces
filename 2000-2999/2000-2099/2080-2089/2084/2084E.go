package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1_000_000_007

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
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		missing := make([]bool, n)
		for i := range missing {
			missing[i] = true
		}
		freeCnt := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			if a[i] == -1 {
				freeCnt++
			} else {
				missing[a[i]] = false
			}
		}

		// factorials up to freeCnt
		fac := make([]int64, freeCnt+1)
		invFac := make([]int64, freeCnt+1)
		fac[0] = 1
		for i := 1; i <= freeCnt; i++ {
			fac[i] = fac[i-1] * int64(i) % mod
		}
		invFac[freeCnt] = modPow(fac[freeCnt], mod-2)
		for i := freeCnt; i >= 1; i-- {
			invFac[i-1] = invFac[i] * int64(i) % mod
		}

		comb := func(n, k int) int64 {
			if k < 0 || k > n {
				return 0
			}
			return fac[n] * invFac[k] % mod * invFac[n-k] % mod
		}

		invComb := make([]int64, freeCnt+1)
		for t := 0; t <= freeCnt; t++ {
			invComb[t] = modPow(comb(freeCnt, t), mod-2)
		}

		// prefix of free positions
		prefFree := make([]int, n+1)
		for i := 1; i <= n; i++ {
			prefFree[i] = prefFree[i-1]
			if a[i-1] == -1 {
				prefFree[i]++
			}
		}

		// min known value up to position i from the left/right
		const infVal = 1 << 30
		minLeft := make([]int, n+1)
		minLeft[0] = infVal
		for i := 1; i <= n; i++ {
			minLeft[i] = minLeft[i-1]
			if a[i-1] != -1 && a[i-1] < minLeft[i] {
				minLeft[i] = a[i-1]
			}
		}
		minRight := make([]int, n+2)
		minRight[n+1] = infVal
		for i := n; i >= 1; i-- {
			minRight[i] = minRight[i+1]
			if a[i-1] != -1 && a[i-1] < minRight[i] {
				minRight[i] = a[i-1]
			}
		}

		// counts of intervals grouped by (limit, freeInside)
		cnt := make([][]int, n+1)
		for i := range cnt {
			cnt[i] = make([]int, freeCnt+1)
		}

		for l := 1; l <= n; l++ {
			for r := l; r <= n; r++ {
				A := prefFree[r] - prefFree[l-1]
				limit := minLeft[l-1]
				if minRight[r+1] < limit {
					limit = minRight[r+1]
				}
				if limit > n {
					limit = n
				}
				cnt[limit][A]++
			}
		}

		s := make([]int64, freeCnt+1) // expected mex for current limit per A
		var ans int64
		kMissing := 0 // number of missing values < current L

		for L := 0; L <= n; L++ {
			// add contributions of intervals whose limit is L
			row := cnt[L]
			for A := 0; A <= freeCnt; A++ {
				if row[A] == 0 {
					continue
				}
				ans = (ans + int64(row[A])*s[A]) % mod
			}

			if L == n {
				break
			}

			if missing[L] { // encountering a missing value, starts a new segment with length 0
				kMissing++
			} else {
				// length of current segment increases by 1 -> w_k += invComb[k]
				if kMissing <= freeCnt {
					delta := invComb[kMissing]
					for A := kMissing; A <= freeCnt; A++ {
						s[A] = (s[A] + comb(A, kMissing)*delta) % mod
					}
				}
			}
		}

		ans = ans % mod
		ans = ans * fac[freeCnt] % mod
		fmt.Fprintln(out, ans)
	}
}
