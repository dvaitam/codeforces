package main

import (
	"bufio"
	"fmt"
	"os"
)

var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

func flush() {
	writer.Flush()
}

func queryA(a int) int {
	fmt.Fprintf(writer, "A %d\n", a)
	flush()
	var cnt int
	fmt.Fscan(reader, &cnt)
	return cnt
}

func queryB(a int) int {
	fmt.Fprintf(writer, "B %d\n", a)
	flush()
	var cnt int
	fmt.Fscan(reader, &cnt)
	return cnt
}

func answer(x int) {
	fmt.Fprintf(writer, "C %d\n", x)
	flush()
	os.Exit(0)
}

func sieve(n int) []int {
	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}
	primes := make([]int, 0)
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
			if i <= n/i {
				for j := i * i; j <= n; j += i {
					isPrime[j] = false
				}
			}
		}
	}
	return primes
}

func main() {
	defer flush()
	var n int
	fmt.Fscan(reader, &n)
	// get all primes up to n
	primes := sieve(n)
	limit := 1
	for limit*limit <= n {
		limit++
	}
	limit-- // limit = floor(sqrt(n))
	res := 1
	// handle small primes
	for _, p := range primes {
		if p > limit {
			break
		}
		cnt := queryA(p)
		if cnt > 0 {
			// find max power
			power := p
			for power*p <= n {
				if queryA(power*p) > 0 {
					power *= p
				} else {
					break
				}
			}
			res *= power
		}
		queryB(p)
	}
	// if result already has factor >1, maybe more large prime
	if res > 1 {
		// check if any large prime remains
		for _, p := range primes {
			if p <= limit {
				continue
			}
			if res*p > n {
				break
			}
			if queryA(p) > 0 {
				res *= p
				break
			}
		}
		answer(res)
	}
	// res == 1, x might be 1 or large prime
	for _, p := range primes {
		if p <= limit {
			continue
		}
		if queryA(p) > 0 {
			res = p
			answer(res)
		}
	}
	// otherwise x is 1
	answer(1)
}
