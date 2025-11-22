package main

import (
	"bufio"
	"fmt"
	"os"
)

// extGCD returns g, x, y such that ax + by = g = gcd(a, b).
func extGCD(a, b int64) (int64, int64, int64) {
	if b == 0 {
		return a, 1, 0
	}
	g, x1, y1 := extGCD(b, a%b)
	return g, y1, x1 - (a/b)*y1
}

func ceilDiv(a, b int64) int64 {
	if a >= 0 {
		return (a + b - 1) / b
	}
	// Division in Go truncates toward zero; for negative numbers that is already a ceiling.
	return a / b
}

func solveCase(n, x, y, vx, vy int64) int64 {
	// Equation: n (vy*a - vx*b) = vy*x - vx*y
	S := vy*x - vx*y
	g := gcd64(vx, vy)
	den := n * g
	if S%den != 0 {
		return -1
	}

	target := S / den
	vx1 := vx / g
	vy1 := vy / g // vx1 and vy1 are coprime

	_, p, q := extGCD(vy1, vx1) // vy1*p + vx1*q = 1
	a0 := p * target
	b0 := -q * target

	// Shift along the solution line to make a, b >= 1.
	k := ceilDiv(1-a0, vx1)
	k2 := ceilDiv(1-b0, vy1)
	if k2 > k {
		k = k2
	}

	a := a0 + vx1*k
	b := b0 + vy1*k
	if a < 1 || b < 1 {
		return -1
	}

	// Number of boundary hits:
	// vertical (x = kn): a-1
	// horizontal (y = kn): b-1
	// lines x+y = n*(2k+1): floor((a+b)/2)
	// lines x-y = n*(2k+1): floor(|a-b|/2)
	ans := a + b + max64(a, b) - 2
	return ans
}

func gcd64(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, x, y, vx, vy int64
		fmt.Fscan(in, &n, &x, &y, &vx, &vy)
		fmt.Fprintln(out, solveCase(n, x, y, vx, vy))
	}
}
