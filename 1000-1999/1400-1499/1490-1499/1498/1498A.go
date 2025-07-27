package main

import (
	"bufio"
	"fmt"
	"os"
)

func sumDigits(x int64) int64 {
	var s int64
	for x > 0 {
		s += x % 10
		x /= 10
	}
	return s
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
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
		var n int64
		fmt.Fscan(reader, &n)
		for {
			g := gcd(n, sumDigits(n))
			if g > 1 {
				fmt.Fprintln(writer, n)
				break
			}
			n++
		}
	}
}
