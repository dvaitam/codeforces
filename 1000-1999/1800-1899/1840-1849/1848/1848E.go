package main

import (
	"bufio"
	"fmt"
	"os"
)

// Sieve of Eratosthenes to generate primes up to limit
func sieve(limit int) []int {
	isPrime := make([]bool, limit+1)
	for i := range isPrime {
		isPrime[i] = true
	}
	isPrime[0] = false
	isPrime[1] = false
	for i := 2; i*i <= limit; i++ {
		if isPrime[i] {
			for j := i * i; j <= limit; j += i {
				isPrime[j] = false
			}
		}
	}
	primes := []int{}
	for i, v := range isPrime {
		if v {
			primes = append(primes, i)
		}
	}
	return primes
}

// fast exponentiation modulo mod
func powMod(a, b, mod int64) int64 {
	res := int64(1)
	a %= mod
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func modInv(a, mod int64) int64 {
	// mod is prime
	return powMod(a, mod-2, mod)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var x int64
	var q int
	var M int64
	fmt.Fscan(reader, &x, &q, &M)

	primes := sieve(1000000)

	exps := make(map[int64]int)
	result := int64(1)
	zeroCnt := 0

	update := func(p int64, add int) {
		old := exps[p]
		if p != 2 {
			oldVal := int64(old+1) % M
			if oldVal == 0 {
				if old > 0 {
					zeroCnt--
				}
			} else {
				result = result * modInv(oldVal, M) % M
			}
		}
		exps[p] = old + add
		if p != 2 {
			newVal := int64(exps[p]+1) % M
			if newVal == 0 {
				zeroCnt++
			} else {
				result = result * newVal % M
			}
		}
	}

	factorize := func(n int64) {
		for _, p := range primes {
			pp := int64(p)
			if pp*pp > n {
				break
			}
			if n%pp == 0 {
				cnt := 0
				for n%pp == 0 {
					n /= pp
					cnt++
				}
				update(pp, cnt)
			}
		}
		if n > 1 {
			update(n, 1)
		}
	}

	factorize(x)

	outputs := make([]int64, q)
	for i := 0; i < q; i++ {
		var mul int64
		fmt.Fscan(reader, &mul)
		factorize(mul)
		if zeroCnt > 0 {
			outputs[i] = 0
		} else {
			outputs[i] = result % M
		}
	}

	for i := 0; i < q; i++ {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, outputs[i])
	}
	fmt.Fprintln(writer)
}
