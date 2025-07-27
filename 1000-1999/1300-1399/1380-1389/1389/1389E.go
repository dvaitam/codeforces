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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var m, d, w int64
		fmt.Fscan(reader, &m, &d, &w)

		if d == 1 {
			fmt.Fprintln(writer, 0)
			continue
		}

		if m > d {
			m = d
		}
		g := gcd(w, d-1)
		q := w / g
		a := m / q
		b := m % q
		ans := b*(a+1)*a/2 + (q-b)*a*(a-1)/2
		fmt.Fprintln(writer, ans)
	}
}
