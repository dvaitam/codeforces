package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

func modpow(a, e int64) int64 {
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
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	if n < 3 || n-1 > m {
		fmt.Fprintln(writer, 0)
		return
	}

	fac := make([]int64, m+1)
	ifac := make([]int64, m+1)
	fac[0] = 1
	for i := 1; i <= m; i++ {
		fac[i] = fac[i-1] * int64(i) % mod
	}
	ifac[m] = modpow(fac[m], mod-2)
	for i := m; i > 0; i-- {
		ifac[i-1] = ifac[i] * int64(i) % mod
	}

    // Formula: nCr(m, n-1) * (n-2) * 2^(n-3)
    cmb := fac[m] * ifac[n-1] % mod * ifac[m-(n-1)] % mod
	term2 := int64(n - 2)
	pow2 := modpow(2, int64(n-3))

	ans := cmb * term2 % mod * pow2 % mod
	fmt.Fprintln(writer, ans)
}
