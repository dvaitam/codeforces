package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)

	minNM := min(n, m)
	maxNM := n + m - 1

	pow := make([]int, k+2)
	pow[0] = 1
	for i := 1; i <= k; i++ {
		pow[i] = (pow[i-1] * 2) % MOD
	}

	for i := minNM; i <= maxNM; i++ {
		if i < k {
			fmt.Fprint(out, 0, " ")
			continue
		}
		val := pow[k-1]
		fmt.Fprint(out, val, " ")
	}
	fmt.Fprintln(out)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

