package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int64
		fmt.Fscan(in, &n, &m)
		g := gcd(n, m)
		n1 := n / g
		m1 := m / g
		if m1&(m1-1) != 0 { // denominator not power of two
			fmt.Fprintln(out, -1)
			continue
		}
		q := n1 / m1
		r := n1 % m1
		needed := m*q + m*int64(bits.OnesCount64(uint64(r)))
		ops := needed - n
		if ops < 0 {
			ops = 0
		}
		fmt.Fprintln(out, ops)
	}
}
