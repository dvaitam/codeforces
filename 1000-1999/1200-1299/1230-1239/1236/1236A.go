package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var a, b, c int
		fmt.Fscan(reader, &a, &b, &c)
		maxStones := 0
		// try all possible numbers of type2 operations
		for y := 0; y*2 <= c && y <= b; y++ {
			x := b - y
			x /= 2
			if x > a {
				x = a
			}
			stones := 3 * (x + y)
			if stones > maxStones {
				maxStones = stones
			}
		}
		fmt.Fprintln(writer, maxStones)
	}
}
