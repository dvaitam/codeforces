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
	if a < 0 {
		return -a
	}
	return a
}

func phi(n int64) int64 {
	result := n
	m := n
	for i := int64(2); i*i <= m; i++ {
		if m%i == 0 {
			for m%i == 0 {
				m /= i
			}
			result = result / i * (i - 1)
		}
	}
	if m > 1 {
		result = result / m * (m - 1)
	}
	return result
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var a, m int64
		fmt.Fscan(reader, &a, &m)
		g := gcd(a, m)
		res := phi(m / g)
		fmt.Fprintln(writer, res)
	}
}
