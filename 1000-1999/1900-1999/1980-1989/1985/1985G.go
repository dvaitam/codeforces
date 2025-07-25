package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1_000_000_007

func modPow(base, exp int64) int64 {
	result := int64(1)
	for exp > 0 {
		if exp&1 == 1 {
			result = result * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return result
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var l, r, k int64
		fmt.Fscan(in, &l, &r, &k)

		if k > 9 {
			fmt.Fprintln(out, 0)
			continue
		}

		limit := int64(9 / k)
		if limit == 0 {
			fmt.Fprintln(out, 0)
			continue
		}
		p := limit + 1
		powL := modPow(p, l)
		powSpan := modPow(p, r-l)
		numerator := limit % mod
		numerator = numerator * powL % mod
		numerator = numerator * ((powSpan - 1 + mod) % mod) % mod
		inv := modPow(p-1, mod-2)
		ans := numerator * inv % mod
		fmt.Fprintln(out, ans)
	}
}
