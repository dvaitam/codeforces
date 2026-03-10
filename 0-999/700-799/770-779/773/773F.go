package main

import (
	"fmt"
)

func mobiusOf(n int) int {
	cnt := 0
	for p := 2; p*p <= n; p++ {
		if n%p == 0 {
			cnt++
			n /= p
			if n%p == 0 {
				return 0
			}
		}
	}
	if n > 1 {
		cnt++
	}
	if cnt%2 == 0 {
		return 1
	}
	return -1
}

func comb(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	result := int64(1)
	for i := 0; i < k; i++ {
		result = result * int64(n-i) / int64(i+1)
	}
	return result
}

// count (n-1)-subsets of {1,...,m-1} whose gcd with m equals 1
// via Möbius: sum_{d|m} mu(d) * C((m-1)/d, n-1)
func countSubsets(m, n int) int64 {
	total := int64(0)
	for d := 1; d*d <= m; d++ {
		if m%d != 0 {
			continue
		}
		if mu := mobiusOf(d); mu != 0 {
			total += int64(mu) * comb((m-1)/d, n-1)
		}
		if d2 := m / d; d2 != d {
			if mu := mobiusOf(d2); mu != 0 {
				total += int64(mu) * comb((m-1)/d2, n-1)
			}
		}
	}
	return total
}

func main() {
	var maxn, maxa, q int
	fmt.Scan(&maxn, &maxa, &q)

	ans := int64(0)
	// Interesting test case conditions: n odd, g even, an/g = m odd
	// Elements: n distinct multiples of g in [1,maxa], largest = g*m
	// Equivalently: choose n-1 elements from {1,...,m-1} with gcd(all,m)=1
	for n := 1; n <= maxn; n += 2 {
		for g := 2; g <= maxa/n; g += 2 {
			M := maxa / g
			for m := n; m <= M; m += 2 { // m odd, m >= n
				c := countSubsets(m, n)
				ans = (ans + c%int64(q) + int64(q)) % int64(q)
			}
		}
	}
	fmt.Println(ans)
}
