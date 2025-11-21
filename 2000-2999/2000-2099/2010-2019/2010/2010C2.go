package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	mod1  int64 = 1_000_000_007
	mod2  int64 = 1_000_000_009
	base1 int64 = 911_382_323
	base2 int64 = 972_663_749
)

func buildHashes(s string, base, mod int64) ([]int64, []int64) {
	n := len(s)
	pref := make([]int64, n+1)
	pow := make([]int64, n+1)
	pow[0] = 1
	for i := 0; i < n; i++ {
		val := int64(s[i]-'a') + 1
		pref[i+1] = (pref[i]*base + val) % mod
		pow[i+1] = (pow[i] * base) % mod
	}
	return pref, pow
}

func getHash(pref, pow []int64, l, r int, mod int64) int64 {
	res := (pref[r] - (pref[l]*pow[r-l])%mod) % mod
	if res < 0 {
		res += mod
	}
	return res
}

func equalSub(pref1, pow1 []int64, pref2, pow2 []int64, l1, r1, l2, r2 int) bool {
	if r1-l1 != r2-l2 {
		return false
	}
	if getHash(pref1, pow1, l1, r1, mod1) != getHash(pref1, pow1, l2, r2, mod1) {
		return false
	}
	if getHash(pref2, pow2, l1, r1, mod2) != getHash(pref2, pow2, l2, r2, mod2) {
		return false
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	t, _ := in.ReadString('\n')
	if len(t) == 0 {
		fmt.Println("NO")
		return
	}
	t = trimNewline(t)
	n := len(t)

	pref1, pow1 := buildHashes(t, base1, mod1)
	pref2, pow2 := buildHashes(t, base2, mod2)

	for x := 1; x < n; x++ {
		if (n-x)%2 != 0 {
			continue
		}
		L := (n + x) / 2
		if L >= n { // needs to be strictly less than n
			continue
		}
		// condition 1: suffix length x equals prefix length x
		if !equalSub(pref1, pow1, pref2, pow2, L-x, L, 0, x) {
			continue
		}
		// condition 2: tail equals middle segment
		lenSecond := (n - x) / 2
		if !equalSub(pref1, pow1, pref2, pow2, L, n, x, x+lenSecond) {
			continue
		}
		fmt.Println("YES")
		fmt.Println(t[:L])
		return
	}
	fmt.Println("NO")
}

func trimNewline(s string) string {
	n := len(s)
	for n > 0 {
		c := s[n-1]
		if c == '\n' || c == '\r' {
			n--
		} else {
			break
		}
	}
	return s[:n]
}
