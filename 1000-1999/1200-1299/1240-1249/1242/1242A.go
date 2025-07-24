package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt.
// We need the maximum number of colors such that tiles i and j
// share a color whenever |i-j| is a divisor of n greater than 1.
// The pattern from small cases shows that the tiles form one of
// three possibilities:
//   - If n = 1, only one tile exists so the answer is 1.
//   - If n is a prime power p^k, the components correspond to the
//     prime p (for p=2 it is always 2). Hence the answer is p.
//   - Otherwise n has at least two distinct prime factors and all
//     tiles become connected, so the answer is 1.
//
// Since n \le 10^12, simple trial division up to sqrt(n) is enough
// to check its prime factors.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int64
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	if n == 1 {
		fmt.Fprintln(out, 1)
		return
	}

	orig := n
	var p int64 = -1
	for i := int64(2); i*i <= n; i++ {
		if n%i == 0 {
			p = i
			break
		}
	}
	if p == -1 {
		// n is prime
		fmt.Fprintln(out, orig)
		return
	}
	for orig%p == 0 {
		orig /= p
	}
	if orig == 1 {
		// n is a power of p
		fmt.Fprintln(out, p)
	} else {
		// n has at least two distinct primes
		fmt.Fprintln(out, 1)
	}
}
