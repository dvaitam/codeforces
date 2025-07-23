package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func powmod(base, exp int64) int64 {
	res := int64(1)
	base %= mod
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}

	// compute mobius function up to k using linear sieve
	mu := make([]int, k+1)
	lp := make([]int, k+1)
	primes := make([]int, 0)
	mu[1] = 1
	for i := 2; i <= k; i++ {
		if lp[i] == 0 {
			lp[i] = i
			primes = append(primes, i)
			mu[i] = -1
		}
		for _, p := range primes {
			if p > lp[i] || p*i > k {
				break
			}
			lp[p*i] = p
			if i%p == 0 {
				mu[p*i] = 0
				break
			} else {
				mu[p*i] = -mu[i]
			}
		}
	}

	powArr := make([]int64, k+1)
	for i := 1; i <= k; i++ {
		powArr[i] = powmod(int64(i), int64(n))
	}

	b := make([]int64, k+1)
	for d := 1; d <= k; d++ {
		if mu[d] == 0 {
			continue
		}
		md := int64(mu[d])
		for m := d; m <= k; m += d {
			b[m] += md * powArr[m/d]
		}
	}

	ans := int64(0)
	for i := 1; i <= k; i++ {
		bi := b[i] % mod
		if bi < 0 {
			bi += mod
		}
		ans = (ans + int64(int(bi)^i)) % mod
	}
	fmt.Fprintln(writer, ans)
}
