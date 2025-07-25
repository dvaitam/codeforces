package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	mod  int64 = 1000000007
	base int64 = 911382323
)

func preprocess(s string) ([]int64, []int64) {
	n := len(s)
	pref := make([]int64, n+1)
	pow := make([]int64, n+1)
	pow[0] = 1
	for i := 0; i < n; i++ {
		pref[i+1] = (pref[i]*base + int64(s[i])) % mod
		pow[i+1] = (pow[i] * base) % mod
	}
	return pref, pow
}

func getHash(pref, pow []int64, l, r int) int64 {
	h := (pref[r] - pref[l]*pow[r-l]) % mod
	if h < 0 {
		h += mod
	}
	return h
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		var s string
		fmt.Fscan(reader, &n, &s)
		prefS, pow := preprocess(s)
		// reversed string
		rs := []byte(s)
		for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
			rs[i], rs[j] = rs[j], rs[i]
		}
		prefR, _ := preprocess(string(rs))

		answer := n
		for L := 1; L <= n; L++ {
			tHash := getHash(prefS, pow, 0, L)
			tRevHash := getHash(prefR, pow, n-L, n)
			ok := true
			for start := 0; start < n; start += L {
				end := start + L
				if end > n {
					end = n
				}
				segLen := end - start
				segHash := getHash(prefS, pow, start, end)
				if segLen == L {
					if segHash != tHash && segHash != tRevHash {
						ok = false
						break
					}
				} else {
					prefHash := getHash(prefS, pow, 0, segLen)
					prefRevHash := getHash(prefR, pow, n-L, n-L+segLen)
					if segHash != prefHash && segHash != prefRevHash {
						ok = false
						break
					}
				}
			}
			if ok {
				answer = L
				break
			}
		}
		fmt.Fprintln(writer, answer)
	}
}
