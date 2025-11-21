package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353
const maxN = 5000

var fact [maxN + 1]int64
var invFact [maxN + 1]int64

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

func precompute() {
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[maxN] = modPow(fact[maxN], mod-2)
	for i := maxN; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}
}

func usedElements(j, m int) int {
	if j == 0 {
		return 0
	}
	if j == m {
		return 2 * m
	}
	return 2*j + 1
}

func requiredOdds(j, m int) int {
	if j == 0 {
		if m == 1 {
			return 1
		}
		return 2
	}
	if j <= m-2 {
		return 1
	}
	return 0
}

func main() {
	precompute()
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		m := n / 2
		// extract distinct values and counts (input already non-decreasing)
		vals := make([]int, 0)
		counts := make([]int, 0)
		for i := 0; i < n; {
			j := i
			for j < n && arr[j] == arr[i] {
				j++
			}
			vals = append(vals, arr[i])
			counts = append(counts, j-i)
			i = j
		}

		dp := make([]int64, m+1)
		dp[0] = 1
		processed := 0

		for idx := range vals {
			c := counts[idx]
			newDp := make([]int64, m+1)
			copy(newDp, dp)
			for j := m - 1; j >= 0; j-- {
				if dp[j] == 0 {
					continue
				}
				r := requiredOdds(j, m)
				avail := processed - usedElements(j, m)
				if avail < r {
					continue
				}
				var ways int64 = 1
				if r == 1 {
					ways = int64(avail)
				} else if r == 2 {
					ways = int64(avail)
					ways = ways * int64(avail-1) % mod
				}
				add := dp[j] * int64(c) % mod
				add = add * ways % mod
				newDp[j+1] = (newDp[j+1] + add) % mod
			}
			dp = newDp
			processed += c
		}

		ans := dp[m]
		for _, c := range counts {
			ans = ans * invFact[c] % mod
		}
		fmt.Fprintln(out, ans)
	}
}
