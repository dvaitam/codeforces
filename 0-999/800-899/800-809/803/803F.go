package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1e9 + 7
const MAXA = 100000

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	freq := make([]int, MAXA+1)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		if x <= MAXA {
			freq[x]++
		}
	}

	cnt := make([]int, MAXA+1)
	for d := 1; d <= MAXA; d++ {
		for j := d; j <= MAXA; j += d {
			cnt[d] += freq[j]
		}
	}

	pow2 := make([]int64, n+1)
	pow2[0] = 1
	for i := 1; i <= n; i++ {
		pow2[i] = pow2[i-1] * 2 % MOD
	}

	g := make([]int64, MAXA+1)
	for d := MAXA; d >= 1; d-- {
		if cnt[d] == 0 {
			g[d] = 0
			continue
		}
		val := pow2[cnt[d]] - 1
		if val < 0 {
			val += MOD
		}
		res := val
		for m := d * 2; m <= MAXA; m += d {
			res -= g[m]
			if res < 0 {
				res += MOD
			}
		}
		g[d] = res % MOD
	}
	fmt.Println(g[1] % MOD)
}
