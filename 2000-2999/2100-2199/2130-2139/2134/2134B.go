package main

import (
	"bufio"
	"fmt"
	"os"
)

func modPow(a, e, mod int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func pickPrime(k int64) int64 {
	primes := []int64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97}
	for _, p := range primes {
		if k%p != 0 {
			return p
		}
	}
	// Fallback (shouldn't happen for k <= 1e9), search forward.
	for p := int64(101); ; p += 2 {
		isPrime := true
		for d := int64(3); d*d <= p; d += 2 {
			if p%d == 0 {
				isPrime = false
				break
			}
		}
		if isPrime && k%p != 0 {
			return p
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	for ; T > 0; T-- {
		var n int
		var k int64
		fmt.Fscan(in, &n, &k)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		p := pickPrime(k)
		inv := modPow(k%p, p-2, p) // modular inverse of k modulo p

		for i := 0; i < n; i++ {
			need := (p - (a[i] % p)) % p
			times := need * inv % p // 0..p-1 <= k
			a[i] += times * k
		}

		for i, v := range a {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		if T > 1 {
			fmt.Fprintln(out)
		}
	}
}
