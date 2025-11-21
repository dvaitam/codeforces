package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const mod int64 = 1000000007

func revDirect(n, p int) int64 {
	if n == 0 {
		return 0
	}
	digits := make([]int, 0, 20)
	value := n
	for value > 0 {
		digits = append(digits, value%p)
		value /= p
	}
	var res int64
	for _, d := range digits {
		res = res*int64(p) + int64(d)
	}
	return res
}

func sumRange(l, r int64) int64 {
	cnt := r - l + 1
	return (l + r) * cnt / 2
}

func prefixSquares(x int64) int64 {
	if x <= 0 {
		return 0
	}
	return x * (x + 1) * (2*x + 1) / 6
}

func sumSquaresRange(l, r int64) int64 {
	return prefixSquares(r) - prefixSquares(l-1)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)

	for ; t > 0; t-- {
		var n int
		var k int64
		fmt.Fscan(in, &n, &k)

		limit := k
		if limit > int64(n) {
			limit = int64(n)
		}
		limitInt := int(limit)

		ans := int64(0)

		if limitInt >= 2 {
			B := int(math.Sqrt(float64(n)))
			if B < 1 {
				B = 1
			}
			directEnd := B
			if directEnd > limitInt {
				directEnd = limitInt
			}
			for p := 2; p <= directEnd; p++ {
				ans += revDirect(n, p) % mod
				if ans >= mod {
					ans -= mod
				}
			}

			start := B + 1
			if start < 2 {
				start = 2
			}
			if start <= limitInt {
				p := start
				nMod := int64(n) % mod
				for p <= limitInt {
					a := n / p
					if a == 0 {
						break
					}
					r := n / a
					if r > limitInt {
						r = limitInt
					}
					l64 := int64(p)
					r64 := int64(r)
					count := r64 - l64 + 1
					sumP := sumRange(l64, r64) % mod
					sumP2 := sumSquaresRange(l64, r64) % mod
					aMod := int64(a) % mod
					cntMod := count % mod

					term := nMod * sumP % mod
					term = (term - aMod*sumP2%mod + mod) % mod
					term = (term + aMod*cntMod) % mod

					ans += term
					ans %= mod

					p = r + 1
				}
			}
		}

		if k > int64(n) {
			tail := (k - int64(n)) % mod
			ans = (ans + tail*(int64(n)%mod)) % mod
		}

		fmt.Fprintln(out, ans%mod)
	}
}
