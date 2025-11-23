package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for i := 0; i < t; i++ {
		var n int
		fmt.Fscan(in, &n)
		// For any prime n >= 2, the sum n + n = 2n is composite (divisible by 2 and n).
		// Since n is prime, m = n is a valid prime number to output.
		fmt.Fprintln(out, n)
	}
}