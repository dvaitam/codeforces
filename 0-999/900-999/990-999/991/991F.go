package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

var memo = make(map[int64]string)

const LIMIT int64 = 10000000000

func powInt(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if a > LIMIT/res {
			return LIMIT + 1
		}
		res *= a
		b--
	}
	return res
}

func root(n int64, b int) int64 {
	if b <= 1 {
		return -1
	}
	x := int64(math.Pow(float64(n), 1.0/float64(b)) + 0.5)
	if x < 2 {
		return -1
	}
	for {
		p := powInt(x, int64(b))
		if p == n {
			return x
		}
		if p > n {
			x--
		} else {
			x++
		}
		if x < 2 {
			return -1
		}
		if powInt(x, int64(b)) == n {
			return x
		}
		if powInt(x, int64(b)) > n && powInt(x-1, int64(b)) < n {
			break
		}
	}
	return -1
}

func solve(n int64) string {
	if v, ok := memo[n]; ok {
		return v
	}
	best := strconv.FormatInt(n, 10)
	bestLen := len(best)

	// exponent representation
	for b := 2; b <= 34; b++ {
		base := root(n, b)
		if base > 1 {
			s1 := solve(base)
			cand := s1 + "^" + strconv.Itoa(b)
			if len(cand) < bestLen {
				best = cand
				bestLen = len(cand)
			}
		}
	}

	// multiplication
	for i := int64(2); i*i <= n; i++ {
		if n%i == 0 {
			s1 := solve(i)
			s2 := solve(n / i)
			cand := s1 + "*" + s2
			if len(cand) < bestLen {
				best = cand
				bestLen = len(cand)
			}
		}
	}

	// using powers of ten
	for k := 1; k <= 10; k++ {
		p := int64(math.Pow10(k))
		if p > n {
			break
		}
		q := n / p
		r := n % p
		if q == 0 {
			continue
		}
		left := solve(q) + "*" + fmt.Sprintf("10^%d", k)
		if r == 0 {
			if len(left) < bestLen {
				best = left
				bestLen = len(left)
			}
		} else {
			right := solve(r)
			cand := left + "+" + right
			if len(cand) < bestLen {
				best = cand
				bestLen = len(cand)
			}
		}
	}

	// addition with small numbers
	for r := int64(1); r <= 9 && r < n; r++ {
		s1 := solve(n - r)
		cand := s1 + "+" + strconv.FormatInt(r, 10)
		if len(cand) < bestLen {
			best = cand
			bestLen = len(cand)
		}
	}

	memo[n] = best
	return best
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int64
	fmt.Fscan(in, &n)
	fmt.Println(solve(n))
}
