package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var a, b, c int
		fmt.Fscan(in, &a, &b, &c)
		diff := abs(a - b)
		denom := 2 * c
		moves := (diff + denom - 1) / denom
		fmt.Fprintln(out, moves)
	}
}
