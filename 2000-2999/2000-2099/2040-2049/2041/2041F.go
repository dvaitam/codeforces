package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const MOD int64 = 998244353
const baseLimit = 2000000

var basePrimes []int

func initBasePrimes() {
	sieve := make([]bool, baseLimit+1)
	for i := 2; i*i <= baseLimit; i++ {
		if !sieve[i] {
			for j := i * i; j <= baseLimit; j += i {
				sieve[j] = true
			}
		}
	}
	for i := 2; i <= baseLimit; i++ {
		if !sieve[i] {
			basePrimes = append(basePrimes, i)
		}
	}
}

func segmentedPrimes(lo, hi int64) []int64 {
	if lo > hi {
		return nil
	}
	length := hi - lo + 1
	mark := make([]bool, length)
	lim := int64(math.Sqrt(float64(hi))) + 1
	for _, p := range basePrimes {
		if int64(p) > lim {
			break
		}
		start := (lo + int64(p) - 1) / int64(p) * int64(p)
		if start < int64(p)*int64(p) {
			start = int64(p) * int64(p)
		}
		for j := start; j <= hi; j += int64(p) {
			mark[j-lo] = true
		}
	}
	primes := make([]int64, 0)
	for i := int64(0); i < length; i++ {
		if !mark[i] && (lo+i) > 1 {
			primes = append(primes, lo+i)
		}
	}
	return primes
}

type result struct {
	length int64
	ways   int64
}

var s []int64
var memo map[int64]result

func keyify(l, r int) int64 {
	return (int64(l) << 32) | int64(r)
}

func solve(l, r int) result {
	if r == l+1 {
		return result{length: s[r] - s[l], ways: 1}
	}
	key := keyify(l, r)
	if val, ok := memo[key]; ok {
		return val
	}
	sum := s[l] + s[r]
	lo := l + 1
	hi := r - 1
	j := l
	for lo <= hi {
		mid := (lo + hi) >> 1
		if 2*s[mid] <= sum {
			j = mid
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	lo = l + 1
	hi = r - 1
	k := r
	for lo <= hi {
		mid := (lo + hi) >> 1
		if 2*s[mid] >= sum {
			k = mid
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	options := make([]result, 0, 2)
	if j > l {
		options = append(options, solve(j, r))
	}
	if k < r {
		options = append(options, solve(l, k))
	}
	var ans result
	if len(options) == 0 {
		ans = result{length: s[r] - s[l], ways: 1}
	} else {
		ans = options[0]
		ans.ways %= MOD
		for i := 1; i < len(options); i++ {
			cur := options[i]
			cur.ways %= MOD
			if cur.length < ans.length {
				ans = cur
			} else if cur.length == ans.length {
				ans.ways = (ans.ways + cur.ways) % MOD
			}
		}
	}
	memo[key] = ans
	return ans
}

func main() {
	initBasePrimes()
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var l, r int64
		fmt.Fscan(in, &l, &r)
		L := 2 * l
		R := 2 * r
		primes := segmentedPrimes(L+1, R-1)
		s = make([]int64, 0, len(primes)+2)
		s = append(s, L)
		s = append(s, primes...)
		s = append(s, R)
		memo = make(map[int64]result)
		res := solve(0, len(s)-1)
		fmt.Fprintln(out, res.ways%MOD)
	}
}
