package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1_000_000_007

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 {
	return a / gcd(a, b) * b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var b, q, y int64
		var c, r, z int64
		fmt.Fscan(in, &b, &q, &y)
		fmt.Fscan(in, &c, &r, &z)
		lastB := b + (y-1)*q
		lastC := c + (z-1)*r
		if (c-b)%q != 0 || r%q != 0 || c < b || lastC > lastB {
			fmt.Fprintln(out, 0)
			continue
		}
		if c-r < b || c+z*r > lastB {
			fmt.Fprintln(out, -1)
			continue
		}
		ans := int64(0)
		for d := int64(1); d*d <= r; d++ {
			if r%d == 0 {
				if lcm(d, q) == r {
					x := r / d
					ans = (ans + x*x) % MOD
				}
				d2 := r / d
				if d2 != d && lcm(d2, q) == r {
					x := r / d2
					ans = (ans + x*x) % MOD
				}
			}
		}
		fmt.Fprintln(out, ans%MOD)
	}
}
