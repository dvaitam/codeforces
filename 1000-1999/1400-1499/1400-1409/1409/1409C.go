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
		var x, y int64
		fmt.Fscan(reader, &n, &x, &y)
		diff := y - x
		var step int64
		for i := n - 1; i >= 1; i-- {
			if diff%int64(i) == 0 {
				step = diff / int64(i)
				break
			}
		}
		// select start to minimize maximum term
		bestMax := int64(1 << 62)
		var start int64
		for i := 0; i < n; i++ {
			cand := y - int64(i)*step
			if cand <= 0 {
				continue
			}
			maxTerm := cand + int64(n-1)*step
			if maxTerm < bestMax {
				bestMax = maxTerm
				start = cand
			}
		}
		// output progression
		for i := 0; i < n; i++ {
			val := start + int64(i)*step
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, val)
		}
		fmt.Fprintln(writer)
	}
}
