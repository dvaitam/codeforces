package main

import (
	"bufio"
	"fmt"
	"math"
	"math/bits"
	"os"
	"sort"
)

type primeOption struct {
	sigma int64
	power int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var A int64
	fmt.Fscan(in, &A)
	if A == 1 {
		fmt.Fprintln(out, 1)
		return
	}

	divs := divisors(A)
	primeMap := make(map[int64][]primeOption)
	for _, d := range divs {
		if d == 1 {
			continue
		}
		p, k, pow, ok := findPrimePower(d)
		if !ok {
			continue
		}
		if _, exists := primeMap[p]; exists {
			// avoid duplicate sigma values
			duplicate := false
			for _, opt := range primeMap[p] {
				if opt.sigma == d {
					duplicate = true
					break
				}
			}
			if duplicate {
				continue
			}
		}
		primeMap[p] = append(primeMap[p], primeOption{sigma: d, power: pow})
		_ = k
	}

	primes := make([]int64, 0, len(primeMap))
	for p := range primeMap {
		primes = append(primes, p)
	}
	sort.Slice(primes, func(i, j int) bool { return primes[i] < primes[j] })
	for _, p := range primes {
		opts := primeMap[p]
		sort.Slice(opts, func(i, j int) bool { return opts[i].sigma < opts[j].sigma })
		primeMap[p] = opts
	}

	type key struct {
		idx int
		rem int64
	}
	memo := make(map[key]int64)

	var dfs func(int, int64) int64
	dfs = func(idx int, rem int64) int64 {
		if rem == 1 {
			return 1
		}
		if idx == len(primes) {
			if rem == 1 {
				return 1
			}
			return 0
		}
		kkey := key{idx, rem}
		if v, ok := memo[kkey]; ok {
			return v
		}
		res := dfs(idx+1, rem)
		for _, opt := range primeMap[primes[idx]] {
			if rem%opt.sigma == 0 {
				res += dfs(idx+1, rem/opt.sigma)
			}
		}
		memo[kkey] = res
		return res
	}

	ans := dfs(0, A)
	fmt.Fprintln(out, ans)
}

func divisors(n int64) []int64 {
	factors := factorize(n)
	divs := []int64{1}
	for p, cnt := range factors {
		cur := make([]int64, 0, len(divs)*(cnt+1))
		mul := int64(1)
		for i := 0; i <= cnt; i++ {
			for _, d := range divs {
				cur = append(cur, d*mul)
			}
			mul *= p
		}
		divs = cur
	}
	return divs
}

func factorize(n int64) map[int64]int {
	res := make(map[int64]int)
	for n%2 == 0 {
		res[2]++
		n /= 2
	}
	limit := int64(math.Sqrt(float64(n))) + 1
	for i := int64(3); i <= limit && n > 1; i += 2 {
		for n%i == 0 {
			res[i]++
			n /= i
		}
	}
	if n > 1 {
		res[n]++
	}
	return res
}

func findPrimePower(d int64) (int64, int, int64, bool) {
	maxK := 60
	sigma2 := int64(0)
	for k := 1; k <= maxK; k++ {
		if k == 1 {
			sigma2 = 3
		} else {
			sigma2 = sigma2*2 + 1
		}
		if sigma2 > d {
			break
		}
		low := int64(2)
		high := int64(2)
		val := calcSigma(high, k, d)
		for val < d {
			high *= 2
			if high > d*2 {
				break
			}
			val = calcSigma(high, k, d)
		}
		if val < d {
			continue
		}
		for low <= high {
			mid := (low + high) / 2
			valMid := calcSigma(mid, k, d)
			if valMid == d {
				if isPrime(mid) {
					pow := int64(1)
					for i := 0; i < k; i++ {
						pow *= mid
					}
					return mid, k, pow, true
				}
				break
			} else if valMid < d {
				low = mid + 1
			} else {
				high = mid - 1
			}
		}
	}
	return 0, 0, 0, false
}

func calcSigma(p int64, k int, limit int64) int64 {
	sum := int64(1)
	term := int64(1)
	for i := 0; i < k; i++ {
		if term > limit/p {
			return limit + 1
		}
		term *= p
		sum += term
		if sum > limit {
			return sum
		}
	}
	return sum
}

func isPrime(n int64) bool {
	if n < 2 {
		return false
	}
	smallPrimes := []int64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37}
	for _, p := range smallPrimes {
		if n == p {
			return true
		}
		if n%p == 0 {
			return false
		}
	}
	return millerRabin(n)
}

func millerRabin(n int64) bool {
	d := n - 1
	s := 0
	for d%2 == 0 {
		d /= 2
		s++
	}
	for _, a := range []int64{2, 3, 5, 7, 11} {
		if a >= n {
			continue
		}
		x := powMod(a, d, n)
		if x == 1 || x == n-1 {
			continue
		}
		composite := true
		for r := 1; r < s; r++ {
			x = mulMod(x, x, n)
			if x == n-1 {
				composite = false
				break
			}
		}
		if composite {
			return false
		}
	}
	return true
}

func powMod(a, e, mod int64) int64 {
	res := int64(1)
	base := a % mod
	for e > 0 {
		if e&1 == 1 {
			res = mulMod(res, base, mod)
		}
		base = mulMod(base, base, mod)
		e >>= 1
	}
	return res
}

func mulMod(a, b, mod int64) int64 {
	hi, lo := bits.Mul64(uint64(a), uint64(b))
	_, rem := bits.Div64(hi, lo, uint64(mod))
	return int64(rem)
}
