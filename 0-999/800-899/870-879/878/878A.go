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

	var n int
	fmt.Fscan(reader, &n)
	var zero, one [10]int
	for i := 0; i < 10; i++ {
		zero[i] = 0
		one[i] = 1
	}

	for i := 0; i < n; i++ {
		var op string
		var x int
		fmt.Fscan(reader, &op, &x)
		for j := 0; j < 10; j++ {
			b := (x >> j) & 1
			switch op {
			case "|":
				zero[j] |= b
				one[j] |= b
			case "&":
				zero[j] &= b
				one[j] &= b
			case "^":
				zero[j] ^= b
				one[j] ^= b
			}
		}
	}

	andMask, orMask, xorMask := 0, 0, 0
	for j := 0; j < 10; j++ {
		a0 := zero[j]
		a1 := one[j]
		if a0 == 0 && a1 == 0 {
			// force 0: leave bit cleared in andMask
		} else if a0 == 1 && a1 == 1 {
			andMask |= 1 << j
			orMask |= 1 << j
		} else if a0 == 0 && a1 == 1 {
			andMask |= 1 << j
		} else { // a0 == 1 && a1 == 0
			andMask |= 1 << j
			xorMask |= 1 << j
		}
	}

	fmt.Fprintln(writer, 3)
	fmt.Fprintln(writer, "&", andMask)
	fmt.Fprintln(writer, "|", orMask)
	fmt.Fprintln(writer, "^", xorMask)
}
