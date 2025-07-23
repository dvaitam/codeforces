package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

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

func modInv(a int64) int64 {
	return modPow(a, mod-2)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var k int
	if _, err := fmt.Fscan(in, &k); err != nil {
		return
	}
	c := make([]int, k)
	sum := 0
	for i := 0; i < k; i++ {
		fmt.Fscan(in, &c[i])
		sum += c[i]
	}

	maxN := sum
	fact := make([]int64, maxN+1)
	invFact := make([]int64, maxN+1)
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[maxN] = modInv(fact[maxN])
	for i := maxN; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}

	res := int64(1)
	pref := c[0]
	for i := 1; i < k; i++ {
		n := pref + c[i] - 1
		choose := fact[n] * invFact[c[i]-1] % mod * invFact[pref] % mod
		res = res * choose % mod
		pref += c[i]
	}
	fmt.Println(res)
}
