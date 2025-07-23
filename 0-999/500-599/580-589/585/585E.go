package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	maxA := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		if a[i] > maxA {
			maxA = a[i]
		}
	}

	freq := make([]int32, maxA+1)
	for _, v := range a {
		freq[v]++
	}

	// count of numbers divisible by d
	cnt := make([]int32, maxA+1)
	for d := 1; d <= maxA; d++ {
		for m := d; m <= maxA; m += d {
			cnt[d] += freq[m]
		}
	}

	pow2 := make([]int64, n+1)
	pow2[0] = 1
	for i := 1; i <= n; i++ {
		pow2[i] = (pow2[i-1] * 2) % mod
	}

	// Mobius function
	mu := make([]int8, maxA+1)
	primes := make([]int, 0)
	isComp := make([]bool, maxA+1)
	mu[1] = 1
	for i := 2; i <= maxA; i++ {
		if !isComp[i] {
			primes = append(primes, i)
			mu[i] = -1
		}
		for _, p := range primes {
			if i*p > maxA {
				break
			}
			isComp[i*p] = true
			if i%p == 0 {
				mu[i*p] = 0
				break
			} else {
				mu[i*p] = -mu[i]
			}
		}
	}

	// numbers coprime with g
	coprime := make([]int64, maxA+1)
	for d := 1; d <= maxA; d++ {
		if mu[d] == 0 {
			continue
		}
		val := int64(mu[d]) * int64(cnt[d])
		for m := d; m <= maxA; m += d {
			coprime[m] += val
		}
	}

	F := make([]int64, maxA+1)
	for g := maxA; g >= 2; g-- {
		c := int(cnt[g])
		if c == 0 {
			continue
		}
		val := (pow2[c] - 1) % mod
		for m := g * 2; m <= maxA; m += g {
			val -= F[m]
			if val < 0 {
				val += mod
			}
		}
		F[g] = val % mod
	}

	var ans int64
	for g := 2; g <= maxA; g++ {
		if F[g] == 0 {
			continue
		}
		co := coprime[g]
		if co <= 0 {
			continue
		}
		ans = (ans + F[g]*(co%mod)) % mod
	}

	fmt.Fprintln(out, ans)
}
