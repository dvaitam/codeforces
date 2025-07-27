package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k, m int
	if _, err := fmt.Fscan(in, &n, &k, &m); err != nil {
		return
	}

	L := make([]int, m)
	R := make([]int, m)
	X := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &L[i], &R[i], &X[i])
	}

	ans := int64(1)

	for bit := 0; bit < k; bit++ {
		diff := make([]int, n+2)
		left := make([]int, n+1)
		zeroSegL := make([]int, 0)
		zeroSegR := make([]int, 0)
		for i := 0; i < m; i++ {
			if (X[i]>>bit)&1 == 1 {
				diff[L[i]]++
				diff[R[i]+1]--
			} else {
				if left[R[i]] < L[i] {
					left[R[i]] = L[i]
				}
				zeroSegL = append(zeroSegL, L[i])
				zeroSegR = append(zeroSegR, R[i])
			}
		}

		forced := make([]bool, n+1)
		prefOnes := make([]int, n+1)
		cur := 0
		for i := 1; i <= n; i++ {
			cur += diff[i]
			if cur > 0 {
				forced[i] = true
				prefOnes[i] = prefOnes[i-1] + 1
			} else {
				prefOnes[i] = prefOnes[i-1]
			}
		}

		valid := true
		for idx := range zeroSegL {
			l := zeroSegL[idx]
			r := zeroSegR[idx]
			if prefOnes[r]-prefOnes[l-1] == r-l+1 {
				valid = false
				break
			}
			if left[r] < l {
				left[r] = l
			}
		}
		if !valid {
			ans = 0
			break
		}

		dp := make([]int64, n+1)
		dp[0] = 1
		pointer := 0
		sumValid := int64(1)
		requirement := 0
		for i := 1; i <= n; i++ {
			validPrev := sumValid
			if !forced[i] {
				dp[i] = validPrev % mod
			}
			sumValid = (sumValid + dp[i]) % mod
			if left[i] > requirement {
				requirement = left[i]
			}
			for pointer < requirement {
				sumValid -= dp[pointer]
				sumValid %= mod
				if sumValid < 0 {
					sumValid += mod
				}
				pointer++
			}
		}

		ans = ans * (sumValid % mod) % mod
	}

	fmt.Fprintln(out, ans%mod)
}
