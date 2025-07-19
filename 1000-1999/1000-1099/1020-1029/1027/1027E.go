package main

import (
	"bufio"
	"fmt"
	"os"
)

const P = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	// dp array a[i][j]: prefix values
	a := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		a[i] = make([]int, n+2)
	}
	// base cases
	for i := 1; i <= n; i++ {
		a[0][i] = 2
		a[1][i] = 2
	}
	// fill dp
	for i := 2; i <= n; i++ {
		for j := 1; j <= i; j++ {
			v := 2 * a[i-1][j]
			if i-j-1 >= 0 {
				v -= a[i-j-1][j]
			}
			v %= P
			if v < 0 {
				v += P
			}
			a[i][j] = v
		}
		for j := i + 1; j <= n; j++ {
			a[i][j] = a[i][j-1]
		}
	}
	// compute differences b[i]
	b := make([]int, n+1)
	for i := 1; i <= n; i++ {
		v := a[n][i] - a[n][i-1]
		if v < 0 {
			v += P
		}
		b[i] = v
	}
	// accumulate answer
	var ans int64
	for i := 1; i <= n; i++ {
		t := (k - 1) / i
		if t > n {
			t = n
		}
		ans = (ans + int64(b[i])*int64(a[n][t])) % P
	}
	// multiply by inverse of 2
	inv2 := (P + 1) / 2
	ans = ans * int64(inv2) % P
	fmt.Fprintln(out, ans)
}
