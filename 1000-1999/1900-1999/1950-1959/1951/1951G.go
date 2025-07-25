package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int64 = 1_000_000_007

func modPow(a, e int64) int64 {
	res := int64(1)
	a %= MOD
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func modInv(a int64) int64 { return modPow(a, MOD-2) }

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		var m int64
		fmt.Fscan(reader, &n, &m)
		pos := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &pos[i])
		}
		sort.Slice(pos, func(i, j int) bool { return pos[i] < pos[j] })

		var prefix, prefix2 int64
		var sum1, sum2 int64
		for j, v := range pos {
			x := v % MOD
			jm := int64(j) % MOD
			sum1 = (sum1 + (x*jm%MOD-MOD+prefix)%MOD) % MOD
			t := (x*x%MOD*jm%MOD - 2*x%MOD*prefix%MOD + prefix2) % MOD
			if t < 0 {
				t += MOD
			}
			sum2 = (sum2 + t) % MOD
			prefix = (prefix + x) % MOD
			prefix2 = (prefix2 + x*x%MOD) % MOD
		}
		total := (2 * ((m%MOD)*sum1%MOD - sum2 + MOD) % MOD) % MOD
		ans := total * modInv(int64(n)%MOD) % MOD
		fmt.Fprintln(writer, ans)
	}
}
