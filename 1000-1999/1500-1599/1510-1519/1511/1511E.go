package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 998244353

func modPow(a, e int) int {
	res := 1
	base := a % MOD
	for e > 0 {
		if e&1 == 1 {
			res = res * base % MOD
		}
		base = base * base % MOD
		e >>= 1
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
	grid := make([]string, n)
	totalWhite := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &grid[i])
		for j := 0; j < m; j++ {
			if grid[i][j] == 'o' {
				totalWhite++
			}
		}
	}

	pow2 := make([]int, totalWhite+1)
	pow2[0] = 1
	for i := 1; i <= totalWhite; i++ {
		pow2[i] = pow2[i-1] * 2 % MOD
	}
	inv18 := modPow(18, MOD-2)

	calc := func(L int) int {
		if L <= 1 {
			return 0
		}
		term := (3*L - 2) % MOD
		term = term * pow2[L] % MOD
		if L%2 == 0 {
			term = (term + 2) % MOD
		} else {
			term = (term + MOD - 2) % MOD
		}
		return term * inv18 % MOD
	}

	ans := 0
	// rows
	for i := 0; i < n; i++ {
		j := 0
		for j < m {
			if grid[i][j] == 'o' {
				k := j
				for k < m && grid[i][k] == 'o' {
					k++
				}
				L := k - j
				add := calc(L)
				add = add * pow2[totalWhite-L] % MOD
				ans = (ans + add) % MOD
				j = k
			} else {
				j++
			}
		}
	}
	// columns
	for j := 0; j < m; j++ {
		i := 0
		for i < n {
			if grid[i][j] == 'o' {
				k := i
				for k < n && grid[k][j] == 'o' {
					k++
				}
				L := k - i
				add := calc(L)
				add = add * pow2[totalWhite-L] % MOD
				ans = (ans + add) % MOD
				i = k
			} else {
				i++
			}
		}
	}

	fmt.Fprintln(out, ans%MOD)
}
