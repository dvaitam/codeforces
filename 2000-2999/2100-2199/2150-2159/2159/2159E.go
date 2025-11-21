package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	mod       int64 = 1000000007
	nMax            = 300000
	kMax            = 600000
	blockSize       = 1024
)

var (
	aVal int64
	bVal int64
	cVal int64
	invC int64
	powC []int64
	invK []int64
)

func modPow(base, exp int64) int64 {
	res := int64(1)
	b := base % mod
	e := exp
	for e > 0 {
		if e&1 == 1 {
			res = res * b % mod
		}
		b = b * b % mod
		e >>= 1
	}
	return res
}

func buildPowC() {
	powC = make([]int64, nMax+1)
	powC[0] = 1
	for i := 1; i <= nMax; i++ {
		powC[i] = powC[i-1] * cVal % mod
	}
}

func buildInvK() {
	invK = make([]int64, kMax+2)
	invK[1] = 1
	for i := 2; i < len(invK); i++ {
		invK[i] = (mod - (mod/int64(i))*invK[int(mod%int64(i))]%mod) % mod
	}
}

func nextCoeff(m, k int, prev1, prev2 int64) int64 {
	t1 := int64(m - k + 1)
	t1 %= mod
	if t1 < 0 {
		t1 += mod
	}
	t2 := int64(2*m - k + 2)
	t2 %= mod
	if t2 < 0 {
		t2 += mod
	}
	val := (t1*bVal%mod*prev1%mod + t2*aVal%mod*prev2%mod) % mod
	val = val * invC % mod
	val = val * invK[k] % mod
	return val
}

func buildSmallPolys() [][]uint32 {
	size := blockSize
	if size > nMax+1 {
		size = nMax + 1
	}
	res := make([][]uint32, size)
	for m := 0; m < size; m++ {
		deg := 2 * m
		if deg > kMax {
			deg = kMax
		}
		coeff := make([]uint32, deg+1)
		c0 := powC[m]
		coeff[0] = uint32(c0)
		if deg >= 1 && m > 0 {
			c1 := int64(m) * bVal % mod * powC[m-1] % mod
			coeff[1] = uint32(c1)
			prev2 := c0
			prev1 := c1
			for k := 2; k <= deg; k++ {
				cur := nextCoeff(m, k, prev1, prev2)
				coeff[k] = uint32(cur)
				prev2 = prev1
				prev1 = cur
			}
		}
		res[m] = coeff
	}
	return res
}

func buildBlockPrefixes() [][]uint32 {
	maxQ := nMax / blockSize
	res := make([][]uint32, maxQ+1)
	for q := 0; q <= maxQ; q++ {
		m := q * blockSize
		deg := 2 * m
		if deg > kMax {
			deg = kMax
		}
		pref := make([]uint32, deg+1)
		c0 := powC[m]
		sum := c0 % mod
		pref[0] = uint32(sum)
		if deg >= 1 && m > 0 {
			c1 := int64(m) * bVal % mod * powC[m-1] % mod
			sum = (sum + c1) % mod
			pref[1] = uint32(sum)
			prev2 := c0
			prev1 := c1
			for k := 2; k <= deg; k++ {
				cur := nextCoeff(m, k, prev1, prev2)
				sum = (sum + cur) % mod
				pref[k] = uint32(sum)
				prev2 = prev1
				prev1 = cur
			}
		}
		res[q] = pref
	}
	return res
}

func solveQuery(n, k int, small [][]uint32, blocks [][]uint32) int64 {
	q := n / blockSize
	r := n - q*blockSize
	smallPoly := small[r]
	limit := len(smallPoly) - 1
	if limit > k {
		limit = k
	}
	pref := blocks[q]
	degBig := len(pref) - 1
	totalBig := int64(pref[degBig])
	ans := int64(0)
	for i := 0; i <= limit; i++ {
		t := k - i
		var prefVal int64
		if t >= degBig {
			prefVal = totalBig
		} else {
			prefVal = int64(pref[t])
		}
		contrib := int64(smallPoly[i]) * prefVal % mod
		ans += contrib
		if ans >= mod {
			ans -= mod
		}
	}
	return ans
}

func bruteForce(n, k int) int64 {
	coeff := []int64{1}
	for i := 0; i < n; i++ {
		next := make([]int64, len(coeff)+2)
		for idx, val := range coeff {
			next[idx] = (next[idx] + val*cVal) % mod
			next[idx+1] = (next[idx+1] + val*bVal) % mod
			next[idx+2] = (next[idx+2] + val*aVal) % mod
		}
		coeff = next
	}
	if k >= len(coeff) {
		k = len(coeff) - 1
	}
	sum := int64(0)
	for i := 0; i <= k; i++ {
		sum += coeff[i]
		if sum >= mod {
			sum -= mod
		}
	}
	return sum
}

func preciseAnswer(n, k int) int64 {
	deg := 2 * n
	coeff := make([]int64, deg+1)
	coeff[0] = powC[n]
	if deg >= 1 && n > 0 {
		coeff[1] = int64(n) * bVal % mod * powC[n-1] % mod
		prev2 := coeff[0]
		prev1 := coeff[1]
		for idx := 2; idx <= deg; idx++ {
			curr := nextCoeff(n, idx, prev1, prev2)
			coeff[idx] = curr
			prev2 = prev1
			prev1 = curr
		}
	}
	if k > deg {
		k = deg
	}
	sum := int64(0)
	for i := 0; i <= k; i++ {
		sum += coeff[i]
		if sum >= mod {
			sum -= mod
		}
	}
	return sum
}

func selfTest(small [][]uint32, blocks [][]uint32) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for n := 0; n <= 60; n++ {
		for k := 0; k <= 2*n; k++ {
			got := solveQuery(n, k, small, blocks)
			want := bruteForce(n, k)
			if got != want {
				panic(fmt.Sprintf("mismatch n=%d k=%d got=%d want=%d", n, k, got, want))
			}
		}
	}
	for t := 0; t < 2000; t++ {
		n := rnd.Intn(200)
		k := rnd.Intn(2*n + 1)
		got := solveQuery(n, k, small, blocks)
		want := bruteForce(n, k)
		if got != want {
			panic(fmt.Sprintf("rand mismatch n=%d k=%d got=%d want=%d", n, k, got, want))
		}
	}
	for t := 0; t < 200; t++ {
		n := rnd.Intn(4000)
		k := rnd.Intn(2*n + 1)
		got := solveQuery(n, k, small, blocks)
		want := preciseAnswer(n, k)
		if got != want {
			panic(fmt.Sprintf("large mismatch n=%d k=%d got=%d want=%d", n, k, got, want))
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var aInput, bInput, cInput int64
	if _, err := fmt.Fscan(reader, &aInput, &bInput, &cInput); err != nil {
		return
	}
	aVal = aInput % mod
	bVal = bInput % mod
	cVal = cInput % mod
	invC = modPow(cVal, mod-2)
	buildPowC()
	buildInvK()
	small := buildSmallPolys()
	blocks := buildBlockPrefixes()
	if os.Getenv("CF_DEBUG") == "1" {
		selfTest(small, blocks)
		return
	}
	var q int
	fmt.Fscan(reader, &q)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var prevAns int64
	for ; q > 0; q-- {
		var nPrime, kPrime int
		fmt.Fscan(reader, &nPrime, &kPrime)
		actualN := nPrime ^ int(prevAns)
		actualK := kPrime ^ int(prevAns)
		ans := solveQuery(actualN, actualK, small, blocks)
		fmt.Fprintln(writer, ans)
		prevAns = ans
	}
}
