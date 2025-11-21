package main

import (
	"bufio"
	"fmt"
	"os"
)

type factor struct {
	p, e int
}

var primes []int

func sieve(limit int) {
	mark := make([]bool, limit+1)
	for i := 2; i <= limit; i++ {
		if !mark[i] {
			primes = append(primes, i)
			if i*i <= limit {
				for j := i * i; j <= limit; j += i {
					mark[j] = true
				}
			}
		}
	}
}

func factorize(x int) []factor {
	res := make([]factor, 0)
	tmp := x
	for _, p := range primes {
		if p*p > tmp {
			break
		}
		if tmp%p == 0 {
			cnt := 0
			for tmp%p == 0 {
				cnt++
				tmp /= p
			}
			res = append(res, factor{p, cnt})
		}
	}
	if tmp > 1 {
		res = append(res, factor{tmp, 1})
	}
	return res
}

func getExp(factors []factor, prime int) int {
	if factors == nil {
		return 0
	}
	for _, f := range factors {
		if f.p == prime {
			return f.e
		}
	}
	return 0
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	sieve(40000)

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}
		s := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &s[i])
		}

		valid := true
		if p[n-1] != s[0] {
			valid = false
		}
		for i := 1; i < n && valid; i++ {
			if p[i-1]%p[i] != 0 {
				valid = false
			}
		}
		for i := 0; i < n-1 && valid; i++ {
			if s[i+1]%s[i] != 0 {
				valid = false
			}
		}
		if !valid {
			fmt.Fprintln(out, "No")
			continue
		}

		facP := make([][]factor, n)
		facS := make([][]factor, n)
		for i := 0; i < n; i++ {
			facP[i] = factorize(p[i])
			facS[i] = factorize(s[i])
		}

		for i := 0; i < n && valid; i++ {
			primesMap := make(map[int]struct{})
			addFactors := func(fs []factor) {
				for _, f := range fs {
					primesMap[f.p] = struct{}{}
				}
			}
			addFactors(facP[i])
			addFactors(facS[i])
			if i > 0 {
				addFactors(facP[i-1])
			}
			if i < n-1 {
				addFactors(facS[i+1])
			}
			for prime := range primesMap {
				A := 0
				if i > 0 {
					A = getExp(facP[i-1], prime)
				}
				B := getExp(facP[i], prime)
				C := getExp(facS[i], prime)
				D := 0
				if i < n-1 {
					D = getExp(facS[i+1], prime)
				}
				lb := B
				if C > lb {
					lb = C
				}
				ub := 60
				if i == 0 {
					if B < ub {
						ub = B
					}
				} else if A > B {
					if B < ub {
						ub = B
					}
				}
				if i == n-1 {
					if C < ub {
						ub = C
					}
				} else if D > C {
					if C < ub {
						ub = C
					}
				}
				if lb > ub {
					valid = false
					break
				}
			}
		}
		if valid {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}
