package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxN = 2000000

var (
	mu   [maxN + 1]int
	spf  [maxN + 1]int
	divs [][]int
)

// precompute initializes the smallest prime factor array, the Möbius function
// and the list of divisors for every number up to maxN.
func precompute() {
	mu[1] = 1
	primes := []int{}
	for i := 2; i <= maxN; i++ {
		if spf[i] == 0 {
			spf[i] = i
			primes = append(primes, i)
			mu[i] = -1
		}
		for _, p := range primes {
			if p*i > maxN {
				break
			}
			spf[p*i] = p
			if i%p == 0 {
				mu[p*i] = 0
				break
			}
			mu[p*i] = -mu[i]
		}
	}
	divs = make([][]int, maxN+1)
	for d := 1; d <= maxN; d++ {
		for m := d; m <= maxN; m += d {
			divs[m] = append(divs[m], d)
		}
	}
}

// countCoprimeRange returns the number of integers x with L <= x <= R
// that are coprime with k.
func countCoprimeRange(k, L, R int) int {
	if L > R {
		return 0
	}
	s := 0
	for _, d := range divs[k] {
		s += mu[d] * (R/d - (L-1)/d)
	}
	return s
}

func solve(n, m int) int64 {
	limit := n
	if m < limit {
		limit = m
	}
	var ans int64
	for g := 1; g <= limit; g++ {
		n1 := n / g
		m1 := m / g
		for _, k := range divs[g] {
			if k <= 1 {
				continue
			}
			L := 1
			if k-m1 > L {
				L = k - m1
			}
			R := k - 1
			if n1 < R {
				R = n1
			}
			if L <= R {
				ans += int64(countCoprimeRange(k, L, R))
			}
		}
	}
	return ans
}

func main() {
	precompute()
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		fmt.Fprintln(out, solve(n, m))
	}
}
