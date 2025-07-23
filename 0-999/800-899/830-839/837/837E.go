package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func primeFactors(n int64) []int64 {
	factors := []int64{}
	for i := int64(2); i*i <= n; i++ {
		if n%i == 0 {
			factors = append(factors, i)
			for n%i == 0 {
				n /= i
			}
		}
	}
	if n > 1 {
		factors = append(factors, n)
	}
	return factors
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var x, y int64
	if _, err := fmt.Fscan(in, &x, &y); err != nil {
		return
	}
	primes := primeFactors(x)
	var ans int64
	for y > 0 {
		g := gcd(x, y)
		x1 := x / g
		if x1 == 1 {
			ans += y / g
			break
		}
		yg := y / g
		rmin := yg
		for _, p := range primes {
			if x1%p == 0 {
				r := yg % p
				if r < rmin {
					rmin = r
				}
			}
		}
		ans += rmin
		y -= rmin * g
	}
	fmt.Println(ans)
}
