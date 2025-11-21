package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)

	ns := make([]int, t)
	maxN := 0
	for i := 0; i < t; i++ {
		fmt.Fscan(in, &ns[i])
		if ns[i] > maxN {
			maxN = ns[i]
		}
	}

	if maxN < 2 {
		maxN = 2
	}
	dp := make([]int64, maxN+1)
	base := []int64{0, 0, 1, 2, 5, 19, 102}
	for i := 0; i < len(base) && i <= maxN; i++ {
		dp[i] = base[i] % mod
	}

	coeff := []int64{875214362, 181553509, 815613431, 916864433, 193014013}
	for n := 7; n <= maxN; n++ {
		val := int64(0)
		for i := 0; i < 5; i++ {
			val = (val + coeff[i]*dp[n-i-1]) % mod
		}
		dp[n] = val
	}

	for _, n := range ns {
		fmt.Fprintln(out, dp[n]%mod)
	}
}
