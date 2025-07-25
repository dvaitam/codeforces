package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func powMod(a, e int64) int64 {
	res := int64(1)
	a %= mod
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func factorPrimes(x int64) []int64 {
	primes := []int64{}
	for p := int64(2); p*p <= x; p++ {
		if x%p == 0 {
			primes = append(primes, p)
			for x%p == 0 {
				x /= p
			}
		}
	}
	if x > 1 {
		primes = append(primes, x)
	}
	return primes
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var x, n int64
	if _, err := fmt.Fscan(reader, &x, &n); err != nil {
		return
	}

	primes := factorPrimes(x)
	ans := int64(1)
	for _, p := range primes {
		exp := int64(0)
		t := n
		for t > 0 {
			t /= p
			exp += t
		}
		ans = ans * powMod(p, exp) % mod
	}

	fmt.Fprintln(writer, ans)
}
