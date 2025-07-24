package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(n, x int64) int64 {
	if x&^n != 0 { // x has bits not set in n
		return -1
	}
	const LIMIT int64 = 5000000000000000000 // 5e18
	lower := n
	upper := LIMIT
	for b := 0; b < 61; b++ {
		bit := int64(1) << uint(b)
		if n&bit != 0 {
			nextZero := (n | (bit - 1)) + 1
			if x&bit == 0 {
				if nextZero > lower {
					lower = nextZero
				}
			} else { // need bit to remain 1
				if nextZero-1 < upper {
					upper = nextZero - 1
				}
			}
		} else {
			if x&bit != 0 {
				return -1
			}
		}
	}
	if lower <= upper {
		return lower
	}
	return -1
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, x int64
		fmt.Fscan(reader, &n, &x)
		fmt.Fprintln(writer, solve(n, x))
	}
}
