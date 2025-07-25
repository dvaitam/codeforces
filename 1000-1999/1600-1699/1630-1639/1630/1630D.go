package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		g := 0
		for i := 0; i < m; i++ {
			var b int
			fmt.Fscan(in, &b)
			g = gcd(g, b)
		}

		sumAbs := make([]int64, g)
		minAbs := make([]int64, g)
		negParity := make([]int, g)
		for i := 0; i < g; i++ {
			minAbs[i] = math.MaxInt64
		}
		for i := 0; i < n; i++ {
			r := i % g
			v := a[i]
			if v < 0 {
				negParity[r] ^= 1
				v = -v
			}
			sumAbs[r] += v
			if v < minAbs[r] {
				minAbs[r] = v
			}
		}

		totalAbs := int64(0)
		for i := 0; i < g; i++ {
			totalAbs += sumAbs[i]
		}

		ans1 := totalAbs
		for i := 0; i < g; i++ {
			if negParity[i] == 1 {
				ans1 -= 2 * minAbs[i]
			}
		}

		ans2 := totalAbs
		for i := 0; i < g; i++ {
			if negParity[i] == 0 {
				ans2 -= 2 * minAbs[i]
			}
		}

		if ans1 > ans2 {
			fmt.Fprintln(out, ans1)
		} else {
			fmt.Fprintln(out, ans2)
		}
	}
}
