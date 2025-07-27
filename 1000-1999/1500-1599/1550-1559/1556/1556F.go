package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1_000_000_007

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

func modInv(a int64) int64 { return modPow(a, mod-2) }

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	// win[i][j] is probability i defeats j modulo mod
	win := make([][]int64, n)
	for i := range win {
		win[i] = make([]int64, n)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				win[i][j] = 1
			} else {
				denom := a[i] + a[j]
				win[i][j] = a[i] % mod * modInv(denom%mod) % mod
			}
		}
	}

	m := 1 << n
	// prod[i][mask] = product of win[i][j] for all j in mask
	prod := make([][]int64, n)
	for i := range prod {
		prod[i] = make([]int64, m)
		prod[i][0] = 1
		for mask := 1; mask < m; mask++ {
			b := mask & -mask
			j := 0
			for b>>j&1 == 0 {
				j++
			}
			prod[i][mask] = prod[i][mask^b] * win[i][j] % mod
		}
	}

	dp := make([]int64, m)
	for mask := 1; mask < m; mask++ {
		if mask&(mask-1) == 0 {
			dp[mask] = 1
			continue
		}
		val := int64(1)
		for sub := (mask - 1) & mask; sub > 0; sub = (sub - 1) & mask {
			if sub == mask {
				continue
			}
			comp := mask ^ sub
			tmp := dp[sub]
			for i := 0; i < n; i++ {
				if sub>>i&1 == 1 {
					tmp = tmp * prod[i][comp] % mod
				}
			}
			val = (val - tmp + mod) % mod
		}
		dp[mask] = val
	}

	full := m - 1
	ans := int64(0)
	for mask := 1; mask <= full; mask++ {
		comp := full ^ mask
		prob := dp[mask]
		for i := 0; i < n; i++ {
			if mask>>i&1 == 1 {
				prob = prob * prod[i][comp] % mod
			}
		}
		cnt := int64(0)
		for i := 0; i < n; i++ {
			if mask>>i&1 == 1 {
				cnt++
			}
		}
		ans = (ans + prob*cnt) % mod
	}

	fmt.Println(ans)
}
