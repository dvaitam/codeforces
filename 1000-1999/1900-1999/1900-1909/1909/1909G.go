package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007
const BASE int64 = 911382323

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func getHash(h []int64, l, r int, pow []int64) int64 {
	v := (h[r] - h[l]*pow[r-l]) % MOD
	if v < 0 {
		v += MOD
	}
	return v
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	var s, t string
	fmt.Fscan(reader, &s)
	fmt.Fscan(reader, &t)

	diff := m - n
	if diff <= 0 {
		fmt.Fprintln(writer, 0)
		return
	}

	pow := make([]int64, m+1)
	pow[0] = 1
	for i := 1; i <= m; i++ {
		pow[i] = pow[i-1] * BASE % MOD
	}

	hs := make([]int64, n+1)
	for i := 0; i < n; i++ {
		hs[i+1] = (hs[i]*BASE + int64(s[i]-'a'+1)) % MOD
	}
	ht := make([]int64, m+1)
	for i := 0; i < m; i++ {
		ht[i+1] = (ht[i]*BASE + int64(t[i]-'a'+1)) % MOD
	}

	lcp := 0
	for lcp < n && lcp < m && s[lcp] == t[lcp] {
		lcp++
	}
	lcs := 0
	for lcs < n && lcs < m && s[n-1-lcs] == t[m-1-lcs] {
		lcs++
	}

	divisors := []int{}
	for d := 1; d*d <= diff; d++ {
		if diff%d == 0 {
			if d <= n {
				divisors = append(divisors, d)
			}
			if d*d != diff {
				other := diff / d
				if other <= n {
					divisors = append(divisors, other)
				}
			}
		}
	}

	ans := 0
	for _, L := range divisors {
		r := diff / L
		powL := pow[L]
		den := (powL - 1) % MOD
		if den < 0 {
			den += MOD
		}
		invDen := modPow(den, MOD-2)
		totalPow := pow[L*(r+1)]
		num := (totalPow - 1) % MOD
		if num < 0 {
			num += MOD
		}
		minJ := L
		if n-lcs > minJ {
			minJ = n - lcs
		}
		maxJ := lcp
		if maxJ > n {
			maxJ = n
		}
		for j := minJ; j <= maxJ; j++ {
			i := j - L
			if getHash(ht, 0, i, pow) != getHash(hs, 0, i, pow) {
				continue
			}
			if getHash(ht, j+diff, m, pow) != getHash(hs, j, n, pow) {
				continue
			}
			hy := getHash(hs, i, j, pow)
			expected := hy * num % MOD * invDen % MOD
			if getHash(ht, i, i+L*(r+1), pow) == expected {
				ans++
			}
		}
	}

	fmt.Fprintln(writer, ans)
}
