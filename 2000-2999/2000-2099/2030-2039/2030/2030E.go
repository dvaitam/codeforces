package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353
const MAXN = 200000

var fact [MAXN + 1]int64
var invFact [MAXN + 1]int64
var pow2 [MAXN + 1]int64

func modPow(a int64, e int64) int64 {
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

func nCr(n, r int) int64 {
	if r < 0 || r > n {
		return 0
	}
	return fact[n] * invFact[r] % MOD * invFact[n-r] % MOD
}

func init() {
	fact[0] = 1
	for i := 1; i <= MAXN; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[MAXN] = modPow(fact[MAXN], MOD-2)
	for i := MAXN; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}
	pow2[0] = 1
	for i := 1; i <= MAXN; i++ {
		pow2[i] = (pow2[i-1] << 1) % MOD
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		cnt := make([]int, n)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			cnt[x]++
		}
		var ans int64
		prefix := 0
		currentMin := n
		prod := make([]int64, 0)
		for v := 0; v < n; v++ {
			c := cnt[v]
			if c == 0 {
				break
			}
			prefix += c
			if currentMin > c {
				currentMin = c
			}
			tails := make([]int64, currentMin)
			tailVal := pow2[c]
			for m := 1; m <= currentMin; m++ {
				tailVal = (tailVal - nCr(c, m-1)) % MOD
				if tailVal < 0 {
					tailVal += MOD
				}
				tails[m-1] = tailVal
			}
			if len(prod) == 0 {
				prod = tails
			} else {
				for i := 0; i < currentMin; i++ {
					prod[i] = prod[i] * tails[i] % MOD
				}
			}
			var sumProd int64
			for i := 0; i < currentMin; i++ {
				sumProd += prod[i]
				if sumProd >= MOD {
					sumProd -= MOD
				}
			}
			rem := pow2[n-prefix]
			ans = (ans + sumProd*rem) % MOD
		}
		if ans < 0 {
			ans += MOD
		}
		fmt.Fprintln(out, ans)
	}
}
