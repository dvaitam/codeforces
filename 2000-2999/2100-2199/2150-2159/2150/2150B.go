package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353
const maxN = 200000

var fact [maxN + 1]int
var invFact [maxN + 1]int

func modPow(base, exp int) int {
	res := 1
	for exp > 0 {
		if exp&1 == 1 {
			res = int(int64(res) * int64(base) % mod)
		}
		base = int(int64(base) * int64(base) % mod)
		exp >>= 1
	}
	return res
}

func prepareFactorials() {
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = int(int64(fact[i-1]) * int64(i) % mod)
	}
	invFact[maxN] = modPow(fact[maxN], mod-2)
	for i := maxN; i >= 1; i-- {
		invFact[i-1] = int(int64(invFact[i]) * int64(i) % mod)
	}
}

func nCr(n, r int) int {
	if r < 0 || r > n {
		return 0
	}
	return int(int64(fact[n]) * int64(invFact[r]) % mod * int64(invFact[n-r]) % mod)
}

type fastScanner struct {
	r *bufio.Reader
}

func newFastScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) nextInt() int {
	sign := 1
	val := 0
	c, err := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, err = fs.r.ReadByte()
		if err != nil {
			return 0
		}
	}
	if c == '-' {
		sign = -1
		c, err = fs.r.ReadByte()
		if err != nil {
			return 0
		}
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

func main() {
	prepareFactorials()
	fs := newFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := fs.nextInt()
	for ; t > 0; t-- {
		n := fs.nextInt()
		a := make([]int, n)
		sum := 0
		for i := 0; i < n; i++ {
			a[i] = fs.nextInt()
			sum += a[i]
		}

		m := (n + 1) / 2
		ans := 0
		ok := true

		if sum != n {
			ok = false
		}
		for i := m; i < n && ok; i++ {
			if a[i] != 0 {
				ok = false
			}
		}

		if ok {
			taken := 0
			ans = 1
			for i := m; i >= 1; i-- {
				length := n - 2*i + 2
				avail := length - taken
				need := a[i-1]
				if avail < need || need < 0 {
					ok = false
					break
				}
				ans = int(int64(ans) * int64(nCr(avail, need)) % mod)
				taken += need
			}
			if !ok {
				ans = 0
			}
		}

		fmt.Fprintln(out, ans)
	}
}
