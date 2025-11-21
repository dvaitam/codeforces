package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

var fact []int64
var invFact []int64

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func initComb(limit int) {
	fact = make([]int64, limit+1)
	invFact = make([]int64, limit+1)
	fact[0] = 1
	for i := 1; i <= limit; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[limit] = modPow(fact[limit], mod-2)
	for i := limit; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}
}

func combMod(n, k int) int64 {
	if n < 0 || k < 0 || k > n {
		return 0
	}
	return fact[n] * invFact[k] % mod * invFact[n-k] % mod
}

func countInteriorWays(m int, rp int64) int64 {
	if m == 0 {
		if rp == 0 {
			return 1
		}
		return 0
	}
	if rp < 0 || rp > int64(2*m) {
		return 0
	}
	r := int(rp)
	tMin := 0
	if r-m > tMin {
		tMin = r - m
	}
	tMax := r / 2
	if tMax > m {
		tMax = m
	}
	res := int64(0)
	for t := tMin; t <= tMax; t++ {
		s := r - 2*t
		if s < 0 || s > m-t {
			continue
		}
		term := combMod(m, t)
		term = term * combMod(m-t, s) % mod
		res += term
		if res >= mod {
			res -= mod
		}
	}
	return res
}

func solveCase(n int64, k int) (int64, int64) {
	S := n - int64(k+1)
	if S == 0 {
		return 0, 1
	}
	if S <= int64(k-1) {
		return 0, combMod(k-1, int(S))
	}
	numerator := S - int64(k-1)
	denom := int64(2 * k)
	C := (numerator + denom - 1) / denom
	if C < 1 {
		C = 1
	}
	cntLess := int64(k-1) + 2*int64(k)*(C-1)
	costLess := int64(k) * (C - 1) * C
	R := S - cntLess
	best := costLess + R*C
	m := k - 1
	ways := int64(0)
	maxEdges := 2
	if int64(maxEdges) > R {
		maxEdges = int(R)
	}
	for x := 0; x <= maxEdges; x++ {
		Rp := R - int64(x)
		if Rp < 0 {
			continue
		}
		waysEdges := combMod(2, x)
		waysInterior := countInteriorWays(m, Rp)
		ways = (ways + waysEdges*waysInterior) % mod
	}
	return best, ways
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	fmt.Fscan(in, &t)
	nArr := make([]int64, t)
	kArr := make([]int, t)
	maxK := 0
	for i := 0; i < t; i++ {
		fmt.Fscan(in, &nArr[i], &kArr[i])
		if kArr[i] > maxK {
			maxK = kArr[i]
		}
	}
	limit := maxK
	if limit < 2 {
		limit = 2
	}
	initComb(limit)

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for i := 0; i < t; i++ {
		best, ways := solveCase(nArr[i], kArr[i])
		fmt.Fprintf(out, "%d %d\n", best, ways%mod)
	}
}
