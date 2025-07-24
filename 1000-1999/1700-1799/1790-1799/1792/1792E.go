package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// factorize performs trial division and returns prime factorization of x.
func factorize(x int64) map[int64]int {
	res := make(map[int64]int)
	for p := int64(2); p*p <= x; p++ {
		for x%p == 0 {
			res[p]++
			x /= p
		}
	}
	if x > 1 {
		res[x]++
	}
	return res
}

// factor d using the primes of m.
func factorWithPrimes(d int64, primes []int64) []int {
	exps := make([]int, len(primes))
	for i, p := range primes {
		for d%p == 0 {
			exps[i]++
			d /= p
		}
	}
	return exps
}

// findMinDiv searches for the minimal divisor in [lo, hi].
func findMinDiv(primes []int64, exps []int, idx int, cur, lo, hi int64, best *int64) {
	if cur > hi || cur >= *best {
		return
	}
	if idx == len(primes) {
		if cur >= lo && cur < *best {
			*best = cur
		}
		return
	}
	val := int64(1)
	for i := 0; i <= exps[idx]; i++ {
		findMinDiv(primes, exps, idx+1, cur*val, lo, hi, best)
		val *= primes[idx]
		if cur*val > hi {
			break
		}
	}
}

// minimalRow returns the minimal row index containing d in n x n table.
func minimalRow(n int64, primes []int64, d int64) int64 {
	if d > n*n {
		return 0
	}
	if d <= n {
		return 1
	}
	lo := (d + n - 1) / n
	exps := factorWithPrimes(d, primes)
	best := int64(1<<63 - 1)
	findMinDiv(primes, exps, 0, 1, lo, n, &best)
	if best == int64(1<<63-1) {
		return 0
	}
	return best
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, m1, m2 int64
		fmt.Fscan(reader, &n, &m1, &m2)

		fac := factorize(m1)
		for p, e := range factorize(m2) {
			fac[p] += e
		}
		type pe struct {
			p int64
			e int
		}
		pes := make([]pe, 0, len(fac))
		for p, e := range fac {
			pes = append(pes, pe{p, e})
		}
		sort.Slice(pes, func(i, j int) bool { return pes[i].p < pes[j].p })
		primes := make([]int64, len(pes))
		exps := make([]int, len(pes))
		for i, v := range pes {
			primes[i] = v.p
			exps[i] = v.e
		}

		// generate all divisors of m
		divisors := []int64{1}
		for i, p := range primes {
			cnt := exps[i]
			curSize := len(divisors)
			pow := int64(1)
			for e := 1; e <= cnt; e++ {
				pow *= p
				for j := 0; j < curSize; j++ {
					divisors = append(divisors, divisors[j]*pow)
				}
			}
		}

		var present int
		var xorVal int64
		for _, d := range divisors {
			row := minimalRow(n, primes, d)
			if row != 0 {
				present++
			}
			xorVal ^= row
		}
		fmt.Fprintln(writer, present, xorVal)
	}
}
