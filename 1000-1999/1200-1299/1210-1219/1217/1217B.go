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
		var n int
		var x int64
		fmt.Fscan(reader, &n, &x)
		var maxD, bestDiff int64
		for i := 0; i < n; i++ {
			var d, h int64
			fmt.Fscan(reader, &d, &h)
			if d > maxD {
				maxD = d
			}
			if d-h > bestDiff {
				bestDiff = d - h
			}
		}
		if maxD >= x {
			fmt.Fprintln(writer, 1)
			continue
		}
		if bestDiff <= 0 {
			fmt.Fprintln(writer, -1)
			continue
		}
		remaining := x - maxD
		steps := (remaining + bestDiff - 1) / bestDiff
		fmt.Fprintln(writer, steps+1)
	}
}
