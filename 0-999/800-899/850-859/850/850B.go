package main

import (
	"bufio"
	"fmt"
	"os"
)

// sieve returns all primes up to n using the Sieve of Eratosthenes.
func sieve(n int) []int {
	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			for j := i * i; j <= n; j += i {
				isPrime[j] = false
			}
		}
	}
	primes := make([]int, 0)
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}
	return primes
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var x, y int64
	if _, err := fmt.Fscan(in, &n, &x, &y); err != nil {
		return
	}

	const limit = 2000000
	cnt := make([]int, limit+1)
	maxA := 0
	for i := 0; i < n; i++ {
		var v int
		fmt.Fscan(in, &v)
		if v > maxA {
			maxA = v
		}
		cnt[v]++
	}

	prefCnt := make([]int64, limit+1)
	prefSum := make([]int64, limit+1)
	for i := 1; i <= limit; i++ {
		prefCnt[i] = prefCnt[i-1] + int64(cnt[i])
		prefSum[i] = prefSum[i-1] + int64(cnt[i])*int64(i)
	}

	primes := sieve(limit)
	deleteAll := int64(n) * x
	minCost := deleteAll

	T := x / y
	if T > int64(limit) {
		T = int64(limit)
	}
	t := int(T)

	for _, p := range primes {
		var cost int64
		for m := p; m <= limit; m += p {
			l := m - p + 1
			if l > maxA {
				break
			}
			r := m
			if r > limit {
				r = limit
			}
			if r > maxA {
				r = maxA
			}
			if l < 1 {
				l = 1
			}
			if l > r {
				continue
			}
			keepStart := r - t
			if keepStart < l {
				keepStart = l
			}
			// cost to increment numbers in [keepStart, r]
			if keepStart <= r {
				cntKeep := prefCnt[r] - prefCnt[keepStart-1]
				sumKeep := prefSum[r] - prefSum[keepStart-1]
				cost += (int64(m)*cntKeep - sumKeep) * y
			}
			// cost to delete numbers in [l, keepStart-1]
			if keepStart-1 >= l {
				delCnt := prefCnt[keepStart-1] - prefCnt[l-1]
				cost += delCnt * x
			}
			if cost >= minCost {
				break
			}
		}
		if cost < minCost {
			minCost = cost
		}
	}

	fmt.Fprintln(out, minCost)
}
