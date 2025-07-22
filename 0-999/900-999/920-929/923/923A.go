package main

import (
	"bufio"
	"fmt"
	"os"
)

func linearSieve(n int) []int {
	spf := make([]int, n+1)
	primes := make([]int, 0)
	for i := 2; i <= n; i++ {
		if spf[i] == 0 {
			spf[i] = i
			primes = append(primes, i)
		}
		for _, p := range primes {
			if p > spf[i] || i*p > n {
				break
			}
			spf[i*p] = p
		}
	}
	return spf
}

func primeFactors(n int, spf []int) []int {
	res := make([]int, 0)
	for n > 1 {
		p := spf[n]
		res = append(res, p)
		for n%p == 0 {
			n /= p
		}
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var X2 int
	if _, err := fmt.Fscan(in, &X2); err != nil {
		return
	}

	const limit = 1000000
	spf := linearSieve(limit)

	factors2 := primeFactors(X2, spf)
	best := X2

	for _, p2 := range factors2 {
		start := X2 - p2 + 1
		if start < 3 {
			start = 3
		}
		for X1 := start; X1 <= X2; X1++ {
			if p2 >= X1 {
				continue
			}
			if X2-X1 >= p2 || X1 > X2 {
				continue
			}
			factors1 := primeFactors(X1, spf)
			for _, p1 := range factors1 {
				if p1 >= X1 {
					continue
				}
				candidate := X1 - p1 + 1
				if candidate < p1+1 {
					candidate = p1 + 1
				}
				if candidate < 3 {
					candidate = 3
				}
				if candidate <= X1 && candidate < best {
					best = candidate
				}
			}
		}
	}

	fmt.Fprintln(out, best)
}
