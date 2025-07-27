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
	for ; t > 0; t-- {
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)
		per := n / k
		x := m
		if x > per {
			x = per
		}
		remaining := m - x
		y := 0
		if k > 1 {
			y = (remaining + k - 2) / (k - 1)
		}
		fmt.Fprintln(out, x-y)
	}
}
