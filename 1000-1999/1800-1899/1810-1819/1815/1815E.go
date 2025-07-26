package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func sieve(n int) []int {
	primes := []int{}
	mark := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		if !mark[i] {
			primes = append(primes, i)
			for j := i * i; j <= n; j += i {
				mark[j] = true
			}
		}
	}
	return primes
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	// pre-sieve primes up to 1e6
	primes := sieve(1000000)
	primeMax := make(map[int]int)

	r := n + 1
	lengths := make([]int, 0)
	for i := 1; i <= n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		if r == n+1 {
			allOne := true
			for j := 0; j < len(s); j++ {
				if s[j] != '1' {
					allOne = false
					break
				}
			}
			if allOne {
				r = i
				continue
			}
		}
		if i < r {
			lengths = append(lengths, len(s))
		}
	}

	if r == n+1 { // no constant 1 oscillator
		r = n + 1
	}

	for _, l := range lengths {
		x := l
		tmp := x
		for _, p := range primes {
			if p*p > x {
				break
			}
			if x%p == 0 {
				cnt := 0
				for x%p == 0 {
					x /= p
					cnt++
				}
				if cnt > primeMax[p] {
					primeMax[p] = cnt
				}
			}
		}
		if x > 1 {
			if primeMax[x] < 1 {
				primeMax[x] = 1
			}
		}
		_ = tmp
	}

	lcm := int64(1)
	for p, e := range primeMax {
		pp := int64(1)
		base := int64(p)
		for ; e > 0; e-- {
			pp = (pp * base) % MOD
		}
		lcm = (lcm * pp) % MOD
	}

	ans := (int64(2*r) % MOD) * lcm % MOD
	fmt.Fprintln(writer, ans)
}
