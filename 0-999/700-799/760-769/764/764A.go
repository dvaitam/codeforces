package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, z int
	if _, err := fmt.Fscan(reader, &n, &m, &z); err != nil {
		return
	}
	g := gcd(n, m)
	lcm := n / g * m
	ans := z / lcm
	fmt.Fprintln(writer, ans)
}
