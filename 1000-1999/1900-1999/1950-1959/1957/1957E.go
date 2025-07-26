package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	N   = 1000000
	MOD = 1000000007
)

func sieve(n int) []int {
	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			for j := i * i; j <= n; j += i {
				isPrime[j] = false
			}
		}
	}
	primes := []int{}
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}
	return primes
}

var prefix []int64

func init() {
	primes := sieve(N)
	diff := make([]int64, N+2)

	for _, p := range primes {
		for q := 1; q <= N/p; q++ {
			val := int64((p - (q % p)) % p)
			l := q * p
			r := (q+1)*p - 1
			if r > N {
				r = N
			}
			diff[l] += val
			diff[r+1] -= val
		}
	}

	for q := 1; q <= N/4; q++ {
		if q%2 == 1 {
			l := 4 * q
			r := l + 3
			if r > N {
				r = N
			}
			diff[l] += 2
			diff[r+1] -= 2
		}
	}

	prefix = make([]int64, N+1)
	var cur int64
	for i := 1; i <= N; i++ {
		cur += diff[i]
		curMod := cur % MOD
		if curMod < 0 {
			curMod += MOD
		}
		prefix[i] = (prefix[i-1] + curMod) % MOD
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t++ {
		var n int
		fmt.Fscan(reader, &n)
		fmt.Fprintln(writer, prefix[n]%MOD)
	}
}
