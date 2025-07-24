package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

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

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	n := len(s)
	// count total '?' and prefix counts
	qpref := make([]int, n+1)
	for i := 0; i < n; i++ {
		qpref[i+1] = qpref[i]
		if s[i] == '?' {
			qpref[i+1]++
		}
	}
	totalQ := qpref[n]
	pow2 := make([]int64, totalQ+1)
	invpow2 := make([]int64, totalQ+1)
	pow2[0] = 1
	for i := 1; i <= totalQ; i++ {
		pow2[i] = pow2[i-1] * 2 % MOD
	}
	invpow2[totalQ] = modPow(pow2[totalQ], MOD-2)
	for i := totalQ; i > 0; i-- {
		invpow2[i-1] = invpow2[i] * 2 % MOD
	}
	// prefix dp for '(' counts
	pref := make([][]int64, n+1)
	for i := range pref {
		pref[i] = make([]int64, n+1)
	}
	pref[0][0] = 1
	for i := 1; i <= n; i++ {
		ch := s[i-1]
		for a := 0; a <= i; a++ {
			var v int64
			if ch == '(' {
				if a > 0 {
					v = pref[i-1][a-1]
				}
			} else if ch == ')' {
				v = pref[i-1][a]
			} else {
				v = pref[i-1][a]
				if a > 0 {
					v = (v + pref[i-1][a-1]) % MOD
				}
			}
			pref[i][a] = v % MOD
		}
	}
	// suffix dp for ')' counts
	suf := make([][]int64, n+2)
	for i := range suf {
		suf[i] = make([]int64, n+1)
	}
	suf[n+1][0] = 1
	for i := n; i >= 1; i-- {
		ch := s[i-1]
		limit := n - i + 1
		for b := 0; b <= limit; b++ {
			var v int64
			if ch == ')' {
				if b > 0 {
					v = suf[i+1][b-1]
				}
			} else if ch == '(' {
				v = suf[i+1][b]
			} else {
				v = suf[i+1][b]
				if b > 0 {
					v = (v + suf[i+1][b-1]) % MOD
				}
			}
			suf[i][b] = v % MOD
		}
	}
	var ans int64
	for k := 1; k <= n; k++ {
		A := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			if s[i-1] == ')' {
				continue
			}
			if k-1 <= i-1 {
				A[i] = pref[i-1][k-1] * invpow2[qpref[i]] % MOD
			}
		}
		var prefix int64
		for j := 1; j <= n; j++ {
			prefix = (prefix + A[j-1]) % MOD
			if s[j-1] == '(' {
				continue
			}
			if k-1 <= n-j {
				B := suf[j+1][k-1] * pow2[qpref[j-1]] % MOD
				ans = (ans + prefix*B) % MOD
			}
		}
	}
	fmt.Println(ans % MOD)
}
