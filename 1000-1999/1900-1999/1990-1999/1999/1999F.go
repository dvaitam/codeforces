package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1_000_000_007

var fact []int64
var invFact []int64

func modPow(a, e int64) int64 {
	res := int64(1)
	base := a % mod
	for e > 0 {
		if e&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		e >>= 1
	}
	return res
}

func comb(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * invFact[k] % mod * invFact[n-k] % mod
}

func precompute(maxN int) {
	fact = make([]int64, maxN+1)
	invFact = make([]int64, maxN+1)
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[maxN] = modPow(fact[maxN], mod-2)
	for i := maxN; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const maxN = 200000
	precompute(maxN)

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)

		ones := 0
		for i := 0; i < n; i++ {
			var val int
			fmt.Fscan(in, &val)
			if val == 1 {
				ones++
			}
		}

		zeros := n - ones
		tReq := (k + 1) / 2
		if tReq > ones {
			fmt.Fprintln(out, 0)
			continue
		}

		upper := k
		if ones < upper {
			upper = ones
		}
		var ans int64
		for j := tReq; j <= upper; j++ {
			z := k - j
			if z > zeros {
				continue
			}
			ans += comb(ones, j) * comb(zeros, z) % mod
			if ans >= mod {
				ans -= mod
			}
		}
		fmt.Fprintln(out, ans%mod)
	}
}
