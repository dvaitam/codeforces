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

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, a, b, p, q int64
	if _, err := fmt.Fscan(in, &n, &a, &b, &p, &q); err != nil {
		return
	}
	g := gcd(a, b)
	lcm := a / g * b
	countA := n / a
	countB := n / b
	countBoth := int64(0)
	if lcm <= n {
		countBoth = n / lcm
	}
	result := (countA-countBoth)*p + (countB-countBoth)*q + countBoth*max(p, q)
	fmt.Println(result)
}
