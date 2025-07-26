package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const MOD int64 = 998244353

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	fmt.Fscan(reader, &n)
	fmt.Fscan(reader, &m)
	g := [10][10]bool{}
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		g[u][v] = true
		g[v][u] = true
	}
	c := make([]int, 11)
	for mask := 0; mask < (1 << 10); mask++ {
		ok := true
		for i := 0; i < 10 && ok; i++ {
			if mask&(1<<i) == 0 {
				continue
			}
			for j := i + 1; j < 10; j++ {
				if mask&(1<<j) == 0 {
					continue
				}
				if !g[i][j] {
					ok = false
					break
				}
			}
		}
		if ok {
			size := bits.OnesCount(uint(mask))
			c[size]++
		}
	}

	h := make([]int64, n+1)
	h[0] = 1
	for i := 1; i <= n; i++ {
		val := int64(0)
		for k := 1; k <= 10 && k <= i; k++ {
			if c[k] == 0 {
				continue
			}
			tmp := int64(c[k]) * h[i-k] % MOD
			if k%2 == 1 {
				val += tmp
			} else {
				val -= tmp
			}
		}
		val %= MOD
		if val < 0 {
			val += MOD
		}
		h[i] = val
	}
	fmt.Println(h[n] % MOD)
}
