package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007
const maxN = 100 * 100

var fact [maxN + 1]int64
var invFact [maxN + 1]int64

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func initFacts() {
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[maxN] = modPow(fact[maxN], mod-2)
	for i := maxN; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}
}

func comb(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * invFact[k] % mod * invFact[n-k] % mod
}

func main() {
	initFacts()

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, c, m int
		fmt.Fscan(in, &n, &c, &m)

		totalSlots := c * (n - 1)
		need := m - c
		ans := comb(totalSlots, need)
		fmt.Fprintln(out, ans)
	}
}
