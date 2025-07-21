package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		f := make([]int64, n+2)
		s := make([]int64, n+2)
		for i := 1; i <= n; i++ {
			var x int64
			fmt.Fscan(in, &x)
			f[i] = (s[i-1] + x*int64(n-i+1)) % MOD
			s[i] = (s[i-1] + f[i]) % MOD
		}
		fmt.Fprintln(out, f[n])
	}
}
