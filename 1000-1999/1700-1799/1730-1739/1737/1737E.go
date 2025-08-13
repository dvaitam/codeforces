package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1_000_000_007
const iv2 int64 = (mod + 1) / 2
const maxN = 1_000_007

var pw2 [maxN]int64

func solve(n int, out *bufio.Writer) {
	r := make([]int64, n+2)

	r[n] = pw2[(n-1)/2]

	i := n - 1
	for ; i > n-i; i-- {
		r[i] = pw2[1+(i-1)/2]
	}

	tot := int64(1)
	for ; i >= 1; i-- {
		tot = (tot + mod - r[i*2] + mod - r[i*2+1]) % mod
		r[i] = pw2[1+(i-1)/2] * tot % mod
	}

	for j := 1; j <= n; j++ {
		fmt.Fprintln(out, r[j])
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	pw2[0] = 1
	for i := 1; i < maxN; i++ {
		pw2[i] = pw2[i-1] * iv2 % mod
	}

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		solve(n, out)
	}
}
