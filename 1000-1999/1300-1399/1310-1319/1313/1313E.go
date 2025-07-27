package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007
const base int64 = 911382323

func buildHash(s []byte) ([]int64, []int64) {
	n := len(s)
	pow := make([]int64, n+1)
	h := make([]int64, n+1)
	pow[0] = 1
	for i := 0; i < n; i++ {
		pow[i+1] = (pow[i] * base) % mod
		h[i+1] = (h[i]*base + int64(s[i])) % mod
	}
	return pow, h
}

func getHash(h, pow []int64, l, r int) int64 {
	res := (h[r] - h[l]*pow[r-l]) % mod
	if res < 0 {
		res += mod
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	var aStr, bStr, sStr string
	fmt.Fscan(in, &aStr)
	fmt.Fscan(in, &bStr)
	fmt.Fscan(in, &sStr)
	a := []byte(aStr)
	b := []byte(bStr)
	s := []byte(sStr)

	powA, hashA := buildHash(a)
	powB, hashB := buildHash(b)
	powS, hashS := buildHash(s)

	// precompute prefix and suffix hashes of s
	prefixHash := make([]int64, m+1)
	for i := 0; i <= m; i++ {
		prefixHash[i] = getHash(hashS, powS, 0, i)
	}
	suffixHash := make([]int64, m+1)
	for i := 0; i <= m; i++ {
		suffixHash[i] = getHash(hashS, powS, m-i, m)
	}

	var ans int64 = 0
	for l1 := 0; l1 < n; l1++ {
		maxLen1 := n - l1
		if maxLen1 > m-1 {
			maxLen1 = m - 1
		}
		for len1 := 1; len1 <= maxLen1; len1++ {
			if getHash(hashA, powA, l1, l1+len1) != prefixHash[len1] {
				continue
			}
			len2 := m - len1
			if len2 <= 0 {
				continue
			}
			for l2 := l1 - len2 + 1; l2 <= l1+len1-1; l2++ {
				if l2 < 0 || l2+len2 > n {
					continue
				}
				if getHash(hashB, powB, l2, l2+len2) == suffixHash[len2] {
					ans++
				}
			}
		}
	}
	fmt.Println(ans)
}
