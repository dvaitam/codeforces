package main

import (
	"bufio"
	"fmt"
	"os"
)

// gcd returns the greatest common divisor of a and b.
// gcd returns the greatest common divisor of a and b.
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

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for i := 0; i < t; i++ {
		var flag int
		fmt.Fscan(reader, &flag)

		pgcd := gcd(360, flag*2)
		if pgcd == 1 {
			fmt.Fprintln(writer, 1)
			continue
		}
		r := 360 / pgcd
		if flag < 90 {
			fmt.Fprintln(writer, r)
		} else {
			m := 2 * flag
			if 360-flag*2 < m {
				m = 360 - flag*2
			}
			if m == 360/r {
				r *= 2
			}
			fmt.Fprintln(writer, r)
		}
	}
}
