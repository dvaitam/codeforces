package main

import (
	"bufio"
	"fmt"
	"os"
)

func powmod(a, b, mod int64) int64 {
	a %= mod
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k, m int64
	if _, err := fmt.Fscan(in, &n, &k, &m); err != nil {
		return
	}

	invK := powmod(k, m-2, m)

	if k%2 == 1 { // k is odd
		powK1 := powmod(k, n-1, m)
		powKMinus1 := powmod(k-1, n, m)
		minus1 := int64(1)
		if n%2 == 1 {
			minus1 = m - 1
		}
		ans := int64(0)
		for x := int64(0); x < k; x++ {
			var nod int64
			if ((n-2)*x)%k == 0 {
				nod = (powKMinus1 + minus1*(k-1)%m) % m
			} else {
				nod = (powKMinus1 - minus1 + m) % m
			}
			nod = nod * invK % m
			lucky := (powK1 - nod + m) % m
			ans = (ans + lucky) % m
		}
		fmt.Fprintln(out, ans)
	} else { // k is even
		powK1 := powmod(k, n-1, m)
		powKMinus2 := powmod(k-2, n, m)
		pow2n := powmod(2, n, m)
		minus1 := int64(1)
		if n%2 == 1 {
			minus1 = m - 1
		}
		half := k / 2
		ans := int64(0)
		for r := int64(0); r < k; r += 2 {
			var C int64 = -1
			if (r*(n-2))%k == 0 {
				C = half - 1
			}
			nod := (powKMinus2 + (minus1*pow2n%m)*C%m + m) % m
			nod = nod * invK % m
			lucky := (powK1 - nod + m) % m
			ans = (ans + lucky) % m
		}
		fmt.Fprintln(out, ans)
	}
}
