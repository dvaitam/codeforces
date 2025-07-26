package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, x, y int64
		fmt.Fscan(reader, &n, &x, &y)

		if x == y {
			fmt.Fprintln(writer, 0)
			continue
		}

		g := gcd(x, y)
		lcm := x / g * y

		cntX := n/x - n/lcm
		cntY := n/y - n/lcm

		sumX := cntX * (2*n - cntX + 1) / 2
		sumY := cntY * (cntY + 1) / 2

		fmt.Fprintln(writer, sumX-sumY)
	}
}
