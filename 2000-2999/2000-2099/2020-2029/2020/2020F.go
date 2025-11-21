package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const MOD int64 = 1_000_000_007

// ---------- combinatorics ----------
var fact, invFact []int64

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func modInv(a int64) int64 { return modPow(a, MOD-2) }

func prepareComb(limit int) {
	fact = make([]int64, limit+1)
	invFact = make([]int64, limit+1)
	fact[0] = 1
	for i := 1; i <= limit; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[limit] = modInv(fact[limit])
	for i := limit; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}
}

func comb(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * invFact[k] % MOD * invFact[n-k] % MOD
}

// ---------- sieve & prime counting (Lehmer) ----------
const MAXSIEVE = 5_000_000 // enough for lehmer up to 1e9

var primes []int
var piSmall []int

func sieve() {
	not := make([]bool, MAXSIEVE+1)
	piSmall = make([]int, MAXSIEVE+1)
	for i := 2; i <= MAXSIEVE; i++ {
		if !not[i] {
			primes = append(primes, i)
			if int64(i)*int64(i) <= MAXSIEVE {
				for j := i * i; j <= MAXSIEVE; j += i {
					not[j] = true
				}
			}
		}
		piSmall[i] = piSmall[i-1]
		if !not[i] {
			piSmall[i]++
		}
	}
}

// memoization for phi in Lehmer
var phiCache map[uint64]int64

func phi(x int64, s int) int64 {
	if s == 0 {
		return x
	}
	if s == 1 {
		return x - x/int64(primes[0])
	}
	if x == 0 {
		return 0
	}
	if x <= int64(primes[s-1]) {
		return 1
	}
	if x < int64(MAXSIEVE) && s < 100 { // memoize only small s to limit map size
		key := (uint64(x) << 7) | uint64(s)
		if v, ok := phiCache[key]; ok {
			return v
		}
		res := phi(x, s-1) - phi(x/int64(primes[s-1]), s-1)
		phiCache[key] = res
		return res
	}
	return phi(x, s-1) - phi(x/int64(primes[s-1]), s-1)
}

func lehmerPi(x int64) int64 {
	if x < int64(MAXSIEVE) {
		return int64(piSmall[x])
	}
	y := int64(math.Sqrt(float64(x)))
	z := int64(math.Pow(float64(x), 1.0/3.0))
	a := lehmerPi(int64(math.Pow(float64(x), 0.25)))
	b := lehmerPi(y)
	c := lehmerPi(z)

	res := phi(x, int(a)) + (b+a-2)*(b-a+1)/2
	for i := a; i < b; i++ {
		p := int64(primes[i])
		w := x / p
		res -= lehmerPi(w)
		if i < c {
			lim := lehmerPi(int64(math.Sqrt(float64(w))))
			for j := i; j < lim; j++ {
				res -= lehmerPi(w/int64(primes[j])) - (j)
			}
		}
	}
	return res
}

// ---------- main solver per test ----------
var curK, curD int
var c1 int64

// memo for dfs
var dfsCache map[uint64]int64

func keyDFS(n int64, idx int) uint64 { return (uint64(n) << 20) | uint64(idx) }

func dfs(n int64, idx int) int64 {
	if n < 2 || idx >= len(primes) || int64(primes[idx]) > n {
		return 1 // only number 1 available
	}
	k := keyDFS(n, idx)
	if v, ok := dfsCache[k]; ok {
		return v
	}
	piN := lehmerPi(n)
	if int64(idx) > piN {
		return 1
	}
	primeOnly := (c1 * ((piN - int64(idx)) % MOD)) % MOD
	res := (1 + primeOnly) % MOD

	sqrtN := int64(math.Sqrt(float64(n)))
	for i := idx; i < len(primes); i++ {
		p := int64(primes[i])
		if p > sqrtN {
			break
		}
		// exponent 1 with further factors
		more := dfs(n/p, i+1) - 1
		if more < 0 {
			more += MOD
		}
		res = (res + c1%MOD*more) % MOD

		pe := p * p
		e := 2
		for pe <= n {
			term := comb(curK*e+curD, curD)
			res = (res + term*dfs(n/pe, i+1)) % MOD
			pe *= p
			e++
		}
	}

	dfsCache[k] = res
	return res
}

func solveCase(n int64) int64 {
	dfsCache = make(map[uint64]int64)
	return dfs(n, 0) % MOD
}

// ---------- main ----------
type testCase struct {
	n int64
	k int
	d int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	cases := make([]testCase, t)
	maxK, maxD := 0, 0
	for i := 0; i < t; i++ {
		fmt.Fscan(in, &cases[i].n, &cases[i].k, &cases[i].d)
		if cases[i].k > maxK {
			maxK = cases[i].k
		}
		if cases[i].d > maxD {
			maxD = cases[i].d
		}
	}

	sieve()
	phiCache = make(map[uint64]int64)

	maxCombN := maxK*30 + maxD + 5
	prepareComb(maxCombN)

	for _, cs := range cases {
		curK = cs.k
		curD = cs.d
		c1 = comb(curK+curD, curD)
		ans := solveCase(cs.n)
		fmt.Fprintln(out, ans)
	}
}
