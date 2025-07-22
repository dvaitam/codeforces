package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007
const INV2 int64 = 500000004

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int64, n+1)
	for i := 1; i <= n-1; i++ {
		fmt.Fscan(reader, &a[i])
	}

	// N[i] = number of nodes on level i
	N := make([]int64, n+1)
	N[1] = 1
	for i := 2; i <= n; i++ {
		N[i] = N[i-1] * (a[i-1] % MOD) % MOD
	}

	// prefix sums of node counts from level i to n
	pref := make([]int64, n+2)
	for i := n; i >= 1; i-- {
		pref[i] = (pref[i+1] + N[i]) % MOD
	}

	ans := make([]int64, 2*n)

	// ancestor-descendant pairs
	for k := 1; k <= n-1; k++ {
		ans[k] = (ans[k] + pref[k+1]) % MOD
	}

	// arrays describing subtree structure
	childLevels := []int64{1}
	childConv := []int64{1}

	for level := n - 1; level >= 1; level-- {
		if a[level] >= 2 {
			comb := a[level] % MOD * ((a[level] - 1) % MOD) % MOD * INV2 % MOD
			factor := N[level] % MOD * comb % MOD
			for i := 0; i < len(childConv); i++ {
				dist := i + 2
				ans[dist] = (ans[dist] + factor*childConv[i]%MOD) % MOD
			}
		}

		// compute current level distribution
		cur := make([]int64, len(childLevels)+1)
		cur[0] = 1
		mul := a[level] % MOD
		for i := 1; i < len(cur); i++ {
			cur[i] = mul * childLevels[i-1] % MOD
		}

		// convolution for next iteration
		conv := make([]int64, 2*(len(cur)-1)+1)
		conv[0] = 1
		if len(conv) > 1 {
			conv[1] = 2 * mul % MOD
		}
		mul2 := mul * mul % MOD
		for d := 2; d < len(conv); d++ {
			var v1 int64
			if d-1 < len(childLevels) {
				v1 = 2 * mul % MOD * childLevels[d-1] % MOD
			}
			var v2 int64
			if d-2 < len(childConv) {
				v2 = mul2 * childConv[d-2] % MOD
			}
			conv[d] = (v1 + v2) % MOD
		}

		childLevels = cur
		childConv = conv
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for k := 1; k <= 2*n-2; k++ {
		if k > 1 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, ans[k]%MOD)
	}
}
