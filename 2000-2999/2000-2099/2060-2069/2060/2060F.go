package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353
const maxComb = 200

var fact [maxComb + 1]int
var invFact [maxComb + 1]int
var spf []int

func modPow(a, e int) int {
	res := 1
	base := a % mod
	for e > 0 {
		if e&1 == 1 {
			res = int(int64(res) * int64(base) % mod)
		}
		base = int(int64(base) * int64(base) % mod)
		e >>= 1
	}
	return res
}

func initFactorials() {
	fact[0] = 1
	for i := 1; i <= maxComb; i++ {
		fact[i] = int(int64(fact[i-1]) * int64(i) % mod)
	}
	invFact[maxComb] = modPow(fact[maxComb], mod-2)
	for i := maxComb; i >= 1; i-- {
		invFact[i-1] = int(int64(invFact[i]) * int64(i) % mod)
	}
}

func comb(n, k int) int {
	if n < 0 || k < 0 || k > n {
		return 0
	}
	return int(int64(fact[n]) * int64(invFact[k]) % mod * int64(invFact[n-k]) % mod)
}

func initSPF(limit int) {
	spf = make([]int, limit+1)
	for i := 2; i <= limit; i++ {
		if spf[i] == 0 {
			spf[i] = i
			if i <= limit/i {
				for j := i * i; j <= limit; j += i {
					if spf[j] == 0 {
						spf[j] = i
					}
				}
			}
		}
	}
}

func factorize(x int) [][2]int {
	res := make([][2]int, 0)
	if x == 1 {
		return res
	}
	for x > 1 {
		p := spf[x]
		if p == 0 {
			p = x
		}
		cnt := 0
		for x%p == 0 {
			x /= p
			cnt++
		}
		res = append(res, [2]int{p, cnt})
	}
	return res
}

func lagrange(y []int, n int) int {
	m := len(y) - 1
	if n <= m {
		return y[n]
	}

	pref := make([]int, m+2)
	suff := make([]int, m+2)
	pref[0] = 1
	for i := 0; i <= m; i++ {
		val := n - i
		val %= mod
		if val < 0 {
			val += mod
		}
		pref[i+1] = int(int64(pref[i]) * int64(val) % mod)
	}
	suff[m+1] = 1
	for i := m; i >= 0; i-- {
		val := n - i
		val %= mod
		if val < 0 {
			val += mod
		}
		suff[i] = int(int64(suff[i+1]) * int64(val) % mod)
	}

	ans := 0
	for i := 0; i <= m; i++ {
		if y[i] == 0 {
			continue
		}
		numer := int(int64(pref[i]) * int64(suff[i+1]) % mod)
		term := int(int64(y[i]) * int64(numer) % mod)
		term = int(int64(term) * int64(invFact[i]) % mod)
		term = int(int64(term) * int64(invFact[m-i]) % mod)
		if (m-i)%2 == 1 {
			term = (mod - term) % mod
		}
		ans += term
		if ans >= mod {
			ans -= mod
		}
	}
	return ans
}

func solveCase(k, n int) []int {
	results := make([]int, k)
	for x := 1; x <= k; x++ {
		factors := factorize(x)
		deg := 0
		for _, pr := range factors {
			deg += pr[1]
		}
		m := deg + 1
		y := make([]int, m+1)
		sumVal := 0
		for L := 1; L <= m; L++ {
			term := 1
			for _, pr := range factors {
				c := pr[1]
				term = int(int64(term) * int64(comb(c+L-1, c)) % mod)
			}
			sumVal += term
			if sumVal >= mod {
				sumVal -= mod
			}
			y[L] = sumVal
		}
		results[x-1] = lagrange(y, n)
	}
	return results
}

func main() {
	initFactorials()
	initSPF(100000)
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var k, n int
		fmt.Fscan(in, &k, &n)
		ans := solveCase(k, n)
		for i := 0; i < k; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[i])
		}
		fmt.Fprintln(out)
	}
}
