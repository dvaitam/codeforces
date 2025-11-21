package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		var mod int64
		fmt.Fscan(in, &n, &k, &mod)
		fmt.Fprintln(out, solveCase(n, k, mod))
	}
}

func solveCase(n, k int, mod int64) int64 {
	if k == 0 {
		return 1 % mod
	}
	size := k + 1
	totPrev := make([]int64, size)
	for i := 0; i <= k; i++ {
		totPrev[i] = 1 % mod
	}
	aPrev := make([][]int64, size)
	bPrev := make([][]int64, size)
	for s := 0; s <= k; s++ {
		aPrev[s] = make([]int64, k+2)
		bPrev[s] = make([]int64, k+2)
		aPrev[s][0] = 1 % mod
		bPrev[s][0] = 1 % mod
	}

	for h := 2; h <= n; h++ {
		prefixTot := make([]int64, size)
		var running int64
		for i := 0; i <= k; i++ {
			running += totPrev[i]
			if running >= mod {
				running -= mod
			}
			prefixTot[i] = running
		}

		totCur := make([]int64, size)
		for s := 0; s <= k; s++ {
			var val int64
			for l := 0; l <= s; l++ {
				rem := s - l
				prod := totPrev[l] * prefixTot[rem] % mod
				val += prod
				if val >= mod {
					val -= mod
				}
			}
			totCur[s] = val
		}

		A1 := make([]int64, size)
		prefixA := make([][]int64, size)
		suffixB := make([][]int64, size)
		for s := 0; s <= k; s++ {
			prefixA[s] = make([]int64, k+2)
			suffixB[s] = make([]int64, k+3)
			var prefixSum int64
			for idx := 0; idx <= k+1; idx++ {
				prefixSum += aPrev[s][idx]
				if prefixSum >= mod {
					prefixSum -= mod
				}
				prefixA[s][idx] = prefixSum
			}
			A1[s] = prefixSum
			var suffixSum int64
			for idx := k + 1; idx >= 0; idx-- {
				suffixSum += bPrev[s][idx]
				if suffixSum >= mod {
					suffixSum -= mod
				}
				suffixB[s][idx] = suffixSum
			}
			suffixB[s][k+2] = 0
		}

		S1 := make([][]int64, size)
		S2 := make([][]int64, size)
		for t := 0; t <= k; t++ {
			S1[t] = make([]int64, size)
			S2[t] = make([]int64, size)
			var run1, run2 int64
			for r := 0; r <= k; r++ {
				idx := r + 2
				var sumGT int64
				if idx <= k+1 {
					sumGT = suffixB[t][idx]
				}
				if sumGT != 0 {
					prod := sumGT * totPrev[r] % mod
					run1 += prod
					if run1 >= mod {
						run1 -= mod
					}
				}
				S1[t][r] = run1

				sumLT := prefixA[t][r]
				if sumLT != 0 {
					prod := sumLT * A1[r] % mod
					run2 += prod
					if run2 >= mod {
						run2 -= mod
					}
				}
				S2[t][r] = run2
			}
		}

		aCur := make([][]int64, size)
		bCur := make([][]int64, size)
		for s := 0; s <= k; s++ {
			aCur[s] = make([]int64, k+2)
			bCur[s] = make([]int64, k+2)
		}

		for s := 0; s <= k; s++ {
			for t := 1; t <= s; t++ {
				R := t - 1
				if s-t < R {
					R = s - t
				}
				if R < 0 {
					continue
				}
				sumTot := prefixTot[R]
				if sumTot != 0 && A1[t] != 0 {
					valA := 2 * (A1[t] % mod) % mod
					valA = valA * sumTot % mod
					aCur[s][t+1] = valA
				}
				valB := (S1[t][R] + S2[t][R]) % mod
				if valB != 0 {
					valB = (valB * 2) % mod
					bCur[s][t+1] = valB
				}
			}
		}

		totPrev = totCur
		aPrev = aCur
		bPrev = bCur
	}

	var ans int64
	for idx := 0; idx <= k+1; idx++ {
		ans += bPrev[k][idx]
		ans %= mod
	}
	return ans
}
