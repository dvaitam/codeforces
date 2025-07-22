package main

import (
	"bufio"
	"fmt"
	"os"
)

// primeFactors returns the distinct prime factors of n.
func primeFactors(n int) []int {
	factors := []int{}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			factors = append(factors, i)
			for n%i == 0 {
				n /= i
			}
		}
	}
	if n > 1 {
		factors = append(factors, n)
	}
	return factors
}

// countCoprime returns number of integers <= limit that are coprime with p,
// given the prime factors of p.
func countCoprime(limit int64, factors []int) int64 {
	if limit <= 0 {
		return 0
	}
	m := len(factors)
	var bad int64
	for mask := 1; mask < (1 << m); mask++ {
		prod := int64(1)
		bits := 0
		for i := 0; i < m; i++ {
			if mask&(1<<i) != 0 {
				prod *= int64(factors[i])
				bits++
			}
		}
		if bits%2 == 1 {
			bad += limit / prod
		} else {
			bad -= limit / prod
		}
	}
	return limit - bad
}

func kthCoprime(x, p, k int) int64 {
	factors := primeFactors(p)
	base := countCoprime(int64(x), factors)
	low := int64(x) + 1
	high := int64(x) + int64(k*p) + 100 // safe upper bound
	for low < high {
		mid := (low + high) / 2
		if countCoprime(mid, factors)-base >= int64(k) {
			high = mid
		} else {
			low = mid + 1
		}
	}
	return low
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for i := 0; i < t; i++ {
		var x, p, k int
		fmt.Fscan(reader, &x, &p, &k)
		ans := kthCoprime(x, p, k)
		fmt.Fprintln(writer, ans)
	}
}
