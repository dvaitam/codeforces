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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	const maxK = 5000
	freq := make([]int, maxK+1)
	for i := 0; i < n; i++ {
		var k int
		fmt.Fscan(in, &k)
		if k >= 0 && k <= maxK {
			freq[k]++
		}
	}

	// Sieve for smallest prime factors up to maxK
	spf := make([]int, maxK+1)
	for i := 2; i <= maxK; i++ {
		if spf[i] == 0 {
			for j := i; j <= maxK; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
			}
		}
	}

	primes := []int{}
	primeIndex := make(map[int]int)
	for i := 2; i <= maxK; i++ {
		if spf[i] == i {
			primeIndex[i] = len(primes)
			primes = append(primes, i)
		}
	}
	P := len(primes)

	// Precompute exponent arrays for k!
	exp := make([][]int, maxK+1)
	lens := make([]int, maxK+1)
	exp[0] = make([]int, P)
	for k := 1; k <= maxK; k++ {
		exp[k] = make([]int, P)
		copy(exp[k], exp[k-1])
		x := k
		cntFactors := 0
		for x > 1 {
			p := spf[x]
			c := 0
			for x%p == 0 {
				x /= p
				c++
			}
			idx := primeIndex[p]
			exp[k][idx] += c
			cntFactors += c
		}
		lens[k] = lens[k-1] + cntFactors
	}

	// Histogram of exponents for each prime
	total := 0
	for _, v := range freq {
		total += v
	}
	maxExp := make([]int, P)
	hist := make([][]int, P)
	for idx := 0; idx < P; idx++ {
		m := 0
		for k := 0; k <= maxK; k++ {
			if freq[k] > 0 && exp[k][idx] > m {
				m = exp[k][idx]
			}
		}
		maxExp[idx] = m
		arr := make([]int, m+1)
		for k := 0; k <= maxK; k++ {
			if freq[k] > 0 {
				e := exp[k][idx]
				if e <= m {
					arr[e] += freq[k]
				}
			}
		}
		for t := m - 1; t >= 0; t-- {
			arr[t] += arr[t+1]
		}
		hist[idx] = arr
	}

	prefix := make([]int, P)
	depth := 0
	half := total / 2
	for idx := P - 1; idx >= 0; idx-- {
		arr := hist[idx]
		m := maxExp[idx]
		for t := 1; t <= m && arr[t] > half; t++ {
			prefix[idx]++
			depth++
		}
	}

	lcp := func(a []int, b []int) int {
		s := 0
		for i := P - 1; i >= 0; i-- {
			if a[i] == b[i] {
				s += a[i]
			} else {
				if a[i] < b[i] {
					s += a[i]
				} else {
					s += b[i]
				}
				break
			}
		}
		return s
	}

	ans := 0
	for k := 0; k <= maxK; k++ {
		if freq[k] == 0 {
			continue
		}
		l := lcp(prefix, exp[k])
		dist := lens[k] + depth - 2*l
		ans += dist * freq[k]
	}

	fmt.Fprintln(out, ans)
}
