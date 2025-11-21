package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	maxA = 1000000
	mod  = 998244353
)

var (
	mu      [maxA + 1]int
	spf     [maxA + 1]int
	sumVals [maxA + 1]int64
	seen    [maxA + 1]int
	testID  int
)

func sieve() {
	mu[1] = 1
	primes := make([]int, 0)
	for i := 2; i <= maxA; i++ {
		if spf[i] == 0 {
			spf[i] = i
			mu[i] = -1
			primes = append(primes, i)
		}
		for _, p := range primes {
			if p*i > maxA {
				break
			}
			spf[p*i] = p
			if i%p == 0 {
				mu[p*i] = 0
				break
			} else {
				mu[p*i] = -mu[i]
			}
		}
	}
}

func getDivisors(x int) []int {
	res := []int{1}
	for x > 1 {
		p := spf[x]
		cnt := 0
		for x%p == 0 {
			x /= p
			cnt++
		}
		curSize := len(res)
		mul := 1
		for c := 0; c < cnt; c++ {
			mul *= p
			for i := 0; i < curSize; i++ {
				res = append(res, res[i]*mul)
			}
		}
	}
	return res
}

func main() {
	sieve()
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}

	testID++
	curID := testID

	dpSum := int64(0)
	ans := int64(0)

	for i := 1; i <= n; i++ {
		divs := getDivisors(a[i])
		var dp int64
		if i == 1 {
			dp = 1
		} else {
			total := dpSum % mod
			g := total
			for _, d := range divs {
				if d == 1 {
					continue
				}
				if seen[d] != curID {
					seen[d] = curID
					sumVals[d] = 0
				}
				g += int64(mu[d]) * sumVals[d]
				g %= mod
			}
			dp = total - g
			dp %= mod
			if dp < 0 {
				dp += mod
			}
		}

		if i == n {
			ans = dp % mod
		}

		dpSum += dp
		dpSum %= mod

		for _, d := range divs {
			if d == 1 {
				continue
			}
			if seen[d] != curID {
				seen[d] = curID
				sumVals[d] = 0
			}
			sumVals[d] += dp
			sumVals[d] %= mod
		}
	}

	fmt.Fprintln(out, ans%mod)
}
