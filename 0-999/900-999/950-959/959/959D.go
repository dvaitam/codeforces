package main

import (
	"bufio"
	"fmt"
	"os"
)

const LIMIT = 2000005

func linearSieve(n int) ([]int, []int) {
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
	return primes, spf
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	primes, spf := linearSieve(LIMIT)
	used := make([]bool, LIMIT+1)
	ans := make([]int, n)
	modified := false
	primeIdx := 0

	for i := 0; i < n; i++ {
		if !modified {
			x := a[i]
			for ; x <= LIMIT; x++ {
				ok := true
				t := x
				for t > 1 {
					p := spf[t]
					if used[p] {
						ok = false
						break
					}
					for t%p == 0 {
						t /= p
					}
				}
				if ok {
					break
				}
			}
			ans[i] = x
			t := x
			for t > 1 {
				p := spf[t]
				used[p] = true
				for t%p == 0 {
					t /= p
				}
			}
			if x > a[i] {
				modified = true
			}
		} else {
			for primeIdx < len(primes) && used[primes[primeIdx]] {
				primeIdx++
			}
			if primeIdx < len(primes) {
				ans[i] = primes[primeIdx]
				used[primes[primeIdx]] = true
				primeIdx++
			} else {
				ans[i] = 2
			}
		}
	}

	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, ans[i])
	}
	fmt.Fprintln(out)
}
