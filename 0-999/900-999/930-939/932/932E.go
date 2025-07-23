package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func powMod(a, e int64) int64 {
	res := int64(1)
	a %= mod
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int64
	var k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	// compute Stirling numbers of the second kind S(k, i)
	Sprev := make([]int64, k+2)
	Scurr := make([]int64, k+2)
	Sprev[0] = 1
	for i := 1; i <= k; i++ {
		for j := 1; j <= i; j++ {
			Scurr[j] = (Sprev[j-1] + int64(j)*Sprev[j]) % mod
		}
		for j := 0; j <= i; j++ {
			Sprev[j] = Scurr[j]
			Scurr[j] = 0
		}
	}

	pow2 := powMod(2, n) // 2^n
	inv2 := (mod + 1) / 2
	under := int64(1) // n^{underline{i}}
	ans := int64(0)

	maxI := k
	if int64(maxI) > n {
		maxI = int(n)
	}
	for i := 1; i <= maxI; i++ {
		pow2 = pow2 * inv2 % mod                         // 2^{n-i}
		under = under * ((n - int64(i) + 1) % mod) % mod // n^{underline{i}}
		term := Sprev[i] * under % mod
		term = term * pow2 % mod
		ans += term
		if ans >= mod {
			ans -= mod
		}
	}

	fmt.Fprintln(out, ans)
}
