package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

// fast exponentiation a^e mod mod
func powMod(a, e int64) int64 {
	a %= mod
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

// memoization map for XY values
var memo = map[int64][2]int64{}

func xy(n int64) (int64, int64) {
	if n == 1 {
		return 1, 0
	}
	if val, ok := memo[n]; ok {
		return val[0], val[1]
	}
	right := n / 2
	left := n - right
	XL, YL := xy(left)
	XR, YR := xy(right)
	f := powMod(2, n) - powMod(2, left) - powMod(2, right) + 1
	f %= mod
	if f < 0 {
		f += mod
	}
	X := (f + 2*XL + 2*XR) % mod
	Y := (YL + XR + YR) % mod
	memo[n] = [2]int64{X, Y}
	return X, Y
}

func solve(n int64) int64 {
	X, Y := xy(n)
	return (X + Y) % mod
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
