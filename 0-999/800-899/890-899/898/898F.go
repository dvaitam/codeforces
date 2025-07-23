package main

import (
	"bufio"
	"fmt"
	"os"
)

func addStrings(a, b string) string {
	i, j := len(a)-1, len(b)-1
	carry := 0
	res := make([]byte, 0, max(len(a), len(b))+1)
	for i >= 0 || j >= 0 || carry > 0 {
		sum := carry
		if i >= 0 {
			sum += int(a[i] - '0')
			i--
		}
		if j >= 0 {
			sum += int(b[j] - '0')
			j--
		}
		res = append(res, byte(sum%10)+'0')
		carry = sum / 10
	}
	// reverse
	for l, r := 0, len(res)-1; l < r; l, r = l+1, r-1 {
		res[l], res[r] = res[r], res[l]
	}
	return string(res)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func substringMod(pre, pow []int64, l, r int, mod int64) int64 {
	val := (pre[r] - pre[l]*pow[r-l]) % mod
	if val < 0 {
		val += mod
	}
	return val
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	s, _ := reader.ReadString('\n')
	if len(s) > 0 && s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}
	n := len(s)
	if n == 0 {
		return
	}

	const mod1 int64 = 1000000007
	const mod2 int64 = 1000000009
	pre1 := make([]int64, n+1)
	pre2 := make([]int64, n+1)
	pow1 := make([]int64, n+1)
	pow2 := make([]int64, n+1)
	pow1[0], pow2[0] = 1, 1
	for i := 0; i < n; i++ {
		d := int64(s[i] - '0')
		pre1[i+1] = (pre1[i]*10 + d) % mod1
		pre2[i+1] = (pre2[i]*10 + d) % mod2
		pow1[i+1] = (pow1[i] * 10) % mod1
		pow2[i+1] = (pow2[i] * 10) % mod2
	}

	check := func(i, j int) bool {
		if i <= 0 || j <= i || j >= n {
			return false
		}
		a := s[:i]
		b := s[i:j]
		c := s[j:]
		if len(a) > 1 && a[0] == '0' {
			return false
		}
		if len(b) > 1 && b[0] == '0' {
			return false
		}
		if len(c) > 1 && c[0] == '0' {
			return false
		}
		la, lb, lc := len(a), len(b), len(c)
		mx := max(la, lb)
		if lc < mx || lc > mx+1 {
			return false
		}
		va1 := substringMod(pre1, pow1, 0, i, mod1)
		vb1 := substringMod(pre1, pow1, i, j, mod1)
		vc1 := substringMod(pre1, pow1, j, n, mod1)
		if (va1+vb1)%mod1 != vc1 {
			return false
		}
		va2 := substringMod(pre2, pow2, 0, i, mod2)
		vb2 := substringMod(pre2, pow2, i, j, mod2)
		vc2 := substringMod(pre2, pow2, j, n, mod2)
		if (va2+vb2)%mod2 != vc2 {
			return false
		}
		sum := addStrings(a, b)
		return sum == c
	}

	// iterate possible length of c
	for lc := n / 3; lc <= n/2; lc++ {
		j := n - lc
		candidates := []int{lc, lc - 1, j - lc, j - (lc - 1)}
		for _, i := range candidates {
			if check(i, j) {
				fmt.Printf("%s+%s=%s\n", s[:i], s[i:j], s[j:])
				return
			}
		}
	}
}
