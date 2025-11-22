package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxN = 200000

func sieve(limit int) []int {
	isPrime := make([]bool, limit+1)
	for i := 2; i <= limit; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= limit; i++ {
		if isPrime[i] {
			for j := i * i; j <= limit; j += i {
				isPrime[j] = false
			}
		}
	}
	primes := make([]int, 0)
	for i := 2; i <= limit; i++ {
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

	primes := sieve(maxN)

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)

		perm := make([]int, n+1) // 1-indexed
		assigned := make([]bool, n+1)

		// iterate primes in descending order
		for pi := len(primes) - 1; pi >= 0; pi-- {
			p := primes[pi]
			if p > n {
				continue
			}
			vec := make([]int, 0)
			for m := p; m <= n; m += p {
				if !assigned[m] {
					vec = append(vec, m)
				}
			}
			if len(vec) >= 2 {
				for i := 0; i < len(vec); i++ {
					next := vec[(i+1)%len(vec)]
					perm[vec[i]] = next
					assigned[vec[i]] = true
				}
			}
		}

		// fill remaining positions with fixed points
		for i := 1; i <= n; i++ {
			if perm[i] == 0 {
				perm[i] = i
			}
		}

		for i := 1; i <= n; i++ {
			if i > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, perm[i])
		}
		fmt.Fprintln(out)
	}
}
