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

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	for ; n > 0; n-- {
		var p, q, b int64
		fmt.Fscan(reader, &p, &q, &b)

		if p == 0 {
			fmt.Fprintln(writer, "Finite")
			continue
		}
		g := gcd(p, q)
		q /= g
		g = gcd(q, b)
		for g > 1 {
			for q%g == 0 {
				q /= g
			}
			g = gcd(q, b)
		}
		if q == 1 {
			fmt.Fprintln(writer, "Finite")
		} else {
			fmt.Fprintln(writer, "Infinite")
		}
	}
}
