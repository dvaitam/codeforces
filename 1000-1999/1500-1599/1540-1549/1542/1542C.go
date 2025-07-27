package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1e9 + 7

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b, limit int64) int64 {
	g := gcd(a, b)
	res := a / g * b
	if res > limit {
		return limit + 1
	}
	return res
}

func solve(n int64) int64 {
	ans := int64(0)
	l := int64(1)
	for m := int64(2); l <= n; m++ {
		nl := lcm(l, m, n)
		cnt := n/l - n/nl
		ans = (ans + cnt*m) % mod
		l = nl
	}
	return ans % mod
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int64
		fmt.Fscan(in, &n)
		fmt.Fprintln(out, solve(n))
	}
}
