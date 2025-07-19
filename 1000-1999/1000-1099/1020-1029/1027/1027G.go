package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// pair holds a factor and its count
type pair struct {
	first  int64
	second int
}

// mulMod computes (a * b) % mod without overflow
func mulMod(a, b, mod int64) int64 {
	var result int64 = 0
	a %= mod
	for b > 0 {
		if b&1 == 1 {
			result = (result + a) % mod
		}
		a = (a << 1) % mod
		b >>= 1
	}
	return result
}

// powMod computes a^pw % mod
func powMod(a, pw, mod int64) int64 {
	var res int64 = 1
	a %= mod
	for pw > 0 {
		if pw&1 == 1 {
			res = mulMod(res, a, mod)
		}
		a = mulMod(a, a, mod)
		pw >>= 1
	}
	return res
}

// get returns prime factorization of n as slice of pairs
func get(n int64) []pair {
	var ans []pair
	for i := int64(2); i*i <= n; i++ {
		if n%i == 0 {
			cnt := 0
			for n%i == 0 {
				cnt++
				n /= i
			}
			ans = append(ans, pair{first: i, second: cnt})
		}
	}
	if n > 1 {
		ans = append(ans, pair{first: n, second: 1})
	}
	return ans
}

// gcd computes greatest common divisor
func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// lcm computes least common multiple
func lcm(a, b int64) int64 {
	return a / gcd(a, b) * b
}

// findOrder finds multiplicative order of x modulo q
func findOrder(x, q int64, factPrs []pair) int64 {
	// extract primes of p-1
	ps := make([]int64, len(factPrs))
	for i, pr := range factPrs {
		ps[i] = pr.first
	}
	// compute phi = product of ps[i]^cnt
	var phi int64 = 1
	for _, pr := range factPrs {
		for j := 0; j < pr.second; j++ {
			phi *= pr.first
		}
	}
	x %= q
	// reduce phi
	for {
		reduced := false
		for _, pp := range ps {
			if phi%pp == 0 && powMod(x, phi/pp, q) == 1 {
				phi /= pp
				reduced = true
			}
		}
		if !reduced {
			break
		}
	}
	return phi
}

// rec recursively counts valid y
func rec(pos int, curLCM, curPhi int64, orders [][]int64, degs []int, ps []int64) int64 {
	if pos == len(degs) {
		// curPhi divisible by curLCM
		return curPhi / curLCM
	}
	var ans int64 = 0
	phi := curPhi
	for i := 0; i <= degs[pos]; i++ {
		ans += rec(pos+1, lcm(curLCM, orders[pos][i]), phi, orders, degs, ps)
		// update phi for next power
		if i == 0 {
			phi *= (ps[pos] - 1)
		} else {
			phi *= ps[pos]
		}
	}
	return ans
}

// solve processes a single test case
func solve(m, x int64) int64 {
	fact := get(m)
	// prepare factorization of p-1 for each prime factor
	factPrsArr := make([][]pair, len(fact))
	for i, pr := range fact {
		factPrsArr[i] = get(pr.first - 1)
	}
	// compute orders
	orders := make([][]int64, len(fact))
   for i, pr := range fact {
       p := pr.first
       deg := pr.second
       orders[i] = make([]int64, deg+1)
       orders[i][0] = 1
       q := p
       for j := 1; j <= deg; j++ {
           ord := findOrder(x, q, factPrsArr[i])
           orders[i][j] = ord
           q *= p
           // mimic C++: append p-factor at j==1, then increment its exponent
           if j == 1 {
               factPrsArr[i] = append(factPrsArr[i], pair{first: p, second: 0})
           }
           // increase exponent of latest p factor
           idx := len(factPrsArr[i]) - 1
           factPrsArr[i][idx].second++
       }
   }
	// prepare degs and ps
	degs := make([]int, len(fact))
	ps := make([]int64, len(fact))
	for i, pr := range fact {
		ps[i] = pr.first
		degs[i] = pr.second
	}
	return rec(0, 1, 1, orders, degs, ps)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for {
		var m, x int64
		if _, err := fmt.Fscan(reader, &m, &x); err != nil {
			break
		}
		ans := solve(m, x)
		fmt.Fprintln(writer, ans)
	}
}
