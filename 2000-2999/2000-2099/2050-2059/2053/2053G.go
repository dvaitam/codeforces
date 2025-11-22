package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	maxLen = 5_000_005

	mod1  int64 = 1_000_000_007
	mod2  int64 = 1_000_000_009
	base1 int64 = 911382323
	base2 int64 = 972663749
)

var (
	pw1  [maxLen]int64
	pw2  [maxLen]int64
	inv1 [maxLen]int64
	inv2 [maxLen]int64
)

type hashPair struct {
	a int64
	b int64
}

func modPow(base, exp, mod int64) int64 {
	res := int64(1)
	cur := base % mod
	for exp > 0 {
		if exp&1 == 1 {
			res = res * cur % mod
		}
		cur = cur * cur % mod
		exp >>= 1
	}
	return res
}

func precompute() {
	pw1[0], pw2[0] = 1, 1
	for i := 1; i < maxLen; i++ {
		pw1[i] = pw1[i-1] * base1 % mod1
		pw2[i] = pw2[i-1] * base2 % mod2
	}

	prod1 := make([]int64, maxLen)
	prod2 := make([]int64, maxLen)
	prod1[0], prod2[0] = 1, 1
	for i := 1; i < maxLen; i++ {
		prod1[i] = prod1[i-1] * ((pw1[i] - 1 + mod1) % mod1) % mod1
		prod2[i] = prod2[i-1] * ((pw2[i] - 1 + mod2) % mod2) % mod2
	}

	ip1 := modPow(prod1[maxLen-1], mod1-2, mod1)
	ip2 := modPow(prod2[maxLen-1], mod2-2, mod2)
	for i := maxLen - 1; i >= 1; i-- {
		inv1[i] = ip1 * prod1[i-1] % mod1
		inv2[i] = ip2 * prod2[i-1] % mod2
		ip1 = ip1 * ((pw1[i] - 1 + mod1) % mod1) % mod1
		ip2 = ip2 * ((pw2[i] - 1 + mod2) % mod2) % mod2
	}
}

func subHash(h []int64, l, r int, pw []int64, mod int64) int64 {
	res := h[r] - h[l]*pw[r-l]%mod
	if res < 0 {
		res += mod
	}
	return res
}

func rangeHash(h1, h2 []int64, l, r int) hashPair {
	return hashPair{
		a: subHash(h1, l, r, pw1[:], mod1),
		b: subHash(h2, l, r, pw2[:], mod2),
	}
}

func repeatHash(h hashPair, length, times int) hashPair {
	if times == 0 {
		return hashPair{}
	}
	total := length * times

	coef1 := pw1[total] - 1
	if coef1 < 0 {
		coef1 += mod1
	}
	var res1 int64
	if coef1 == 0 {
		res1 = (h.a * int64(times)) % mod1
	} else {
		res1 = h.a * coef1 % mod1
		res1 = res1 * inv1[length] % mod1
	}

	coef2 := pw2[total] - 1
	if coef2 < 0 {
		coef2 += mod2
	}
	var res2 int64
	if coef2 == 0 {
		res2 = (h.b * int64(times)) % mod2
	} else {
		res2 = h.b * coef2 % mod2
		res2 = res2 * inv2[length] % mod2
	}

	return hashPair{a: res1, b: res2}
}

func getMax(h1, h2 []int64, l, r int, pat hashPair, length int) int {
	hi := (r - l) / length
	lo := 0
	for lo < hi {
		mid := (lo + hi + 1) >> 1
		sub := rangeHash(h1, h2, l, l+mid*length)
		if sub == repeatHash(pat, length, mid) {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return lo
}

func prefixFunction(s string) []int {
	n := len(s)
	pi := make([]int, n)
	for i := 1; i < n; i++ {
		j := pi[i-1]
		for j > 0 && s[i] != s[j] {
			j = pi[j-1]
		}
		if s[i] == s[j] {
			j++
		}
		pi[i] = j
	}
	return pi
}

func exgcd(a, b int) (g, x, y int) {
	if b == 0 {
		return a, 1, 0
	}
	g, x1, y1 := exgcd(b, a%b)
	x = y1
	y = x1 - a/b*y1
	return
}

func solveCase(in *bufio.Reader, out *bufio.Writer) {
	var n, m int
	fmt.Fscan(in, &n, &m)

	var s, t string
	fmt.Fscan(in, &s)
	fmt.Fscan(in, &t)

	hs1 := make([]int64, n+1)
	hs2 := make([]int64, n+1)
	for i := 0; i < n; i++ {
		hs1[i+1] = (hs1[i]*base1 + int64(s[i])) % mod1
		hs2[i+1] = (hs2[i]*base2 + int64(s[i])) % mod2
	}
	ht1 := make([]int64, m+1)
	ht2 := make([]int64, m+1)
	for i := 0; i < m; i++ {
		ht1[i+1] = (ht1[i]*base1 + int64(t[i])) % mod1
		ht2[i+1] = (ht2[i]*base2 + int64(t[i])) % mod2
	}

	pi := prefixFunction(s)
	T := n
	per := n - pi[n-1]
	if n%per == 0 {
		T = per
	}

	ans := make([]byte, n-1)
	for i := 0; i < n-1; i++ {
		ans[i] = '0'
	}

	check := func(split int) bool {
		if split%T == 0 {
			if m%T != 0 {
				return false
			}
			baseHash := rangeHash(hs1, hs2, 0, T)
			if rangeHash(ht1, ht2, 0, m) != repeatHash(baseHash, T, m/T) {
				return false
			}
			a := split
			b := n - split
			g, x, _ := exgcd(a, b)
			if m%g != 0 {
				return false
			}
			x = (int(int64(x)*(int64(m)/int64(g))) % (b / g))
			if x < 0 {
				x += b / g
			}
			return int64(x)*int64(a) <= int64(m)
		}

		la, lb := split, n-split
		ha := rangeHash(hs1, hs2, 0, split)
		hb := rangeHash(hs1, hs2, split, n)

		times := 0
		if la <= lb {
			times = getMax(hs1, hs2, split, n, ha, la)
		} else {
			times = getMax(hs1, hs2, 0, split, hb, lb)
			la, lb = lb, la
			ha, hb = hb, ha
		}

		for pos := 0; pos < m; {
			x := getMax(ht1, ht2, pos, m, ha, la)
			if pos+x*la == m {
				return true
			}
			if x < times {
				return false
			}
			pos1 := pos + (x-times)*la + lb
			if pos1 <= m && rangeHash(ht1, ht2, pos+(x-times)*la, pos1) == hb {
				pos = pos1
				continue
			}
			if x > times {
				pos2 := pos + (x-times-1)*la + lb
				if pos2 <= m && rangeHash(ht1, ht2, pos+(x-times-1)*la, pos2) == hb {
					pos = pos2
					continue
				}
			}
			return false
		}
		return true
	}

	for i := 1; i < n; i++ {
		if check(i) {
			ans[i-1] = '1'
		}
	}
	fmt.Fprintln(out, string(ans))
}

func main() {
	precompute()

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		solveCase(in, out)
	}
}
