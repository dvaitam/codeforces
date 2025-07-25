package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353
const MAXA = 200005

var fact [MAXA]int64
var invFact [MAXA]int64
var invNum [MAXA]int64

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func init() {
	fact[0] = 1
	for i := 1; i < MAXA; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[MAXA-1] = modPow(fact[MAXA-1], MOD-2)
	for i := MAXA - 1; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}
	invNum[1] = 1
	for i := 2; i < MAXA; i++ {
		invNum[i] = MOD - MOD/int64(i)*invNum[int(MOD%int64(i))]%MOD
	}
}

type Fenwick struct {
	n    int
	tree []int
}

func (f *Fenwick) init(n int) {
	f.n = n
	f.tree = make([]int, n+2)
}

func (f *Fenwick) add(i, delta int) {
	for i <= f.n {
		f.tree[i] += delta
		i += i & -i
	}
}

func (f *Fenwick) sum(i int) int {
	res := 0
	for i > 0 {
		res += f.tree[i]
		i -= i & -i
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	freq := make([]int, MAXA)
	for i := 0; i < n; i++ {
		var v int
		fmt.Fscan(in, &v)
		if v < MAXA {
			freq[v]++
		}
	}
	t := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &t[i])
	}

	fen := Fenwick{}
	fen.init(MAXA - 1)
	for i := 1; i < MAXA; i++ {
		if freq[i] > 0 {
			fen.add(i, freq[i])
		}
	}

	total := fact[n]
	for i := 1; i < MAXA; i++ {
		if freq[i] > 0 {
			total = total * invFact[freq[i]] % MOD
		}
	}

	ans := int64(0)
	limit := n
	if m < limit {
		limit = m
	}

	for i := 0; i < limit; i++ {
		remaining := n - i
		if remaining <= 0 {
			break
		}
		less := fen.sum(t[i] - 1)
		if less > 0 {
			contrib := total * int64(less) % MOD * invNum[remaining] % MOD
			ans = (ans + contrib) % MOD
		}
		if freq[t[i]] == 0 {
			fmt.Fprintln(out, ans%MOD)
			return
		}
		total = total * int64(freq[t[i]]) % MOD * invNum[remaining] % MOD
		freq[t[i]]--
		fen.add(t[i], -1)
	}

	if n < m {
		ans = (ans + 1) % MOD
	}
	fmt.Fprintln(out, ans%MOD)
}
