package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement the correct solution for problem E.
// This placeholder outputs 0 for the first ant and assumes the rest have equal
// probability to survive.

const mod int64 = 1_000_000_007

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
		if n <= 1 {
			fmt.Fprintln(out, 1)
			continue
		}
		inv := modInv(int64(n - 1))
		fmt.Fprintln(out, 0)
		for i := 2; i <= n; i++ {
			fmt.Fprintln(out, inv)
		}
	}
}

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func modInv(a int64) int64 {
	return modPow(a%mod, mod-2)
}
