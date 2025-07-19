package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for tc := 0; tc < t; tc++ {
		var n int
		var s int64
		fmt.Fscan(reader, &n, &s)
		var sum int64
		var maxVal int64
		var maxIdx int
		bestLen := 0
		bestRemove := 0
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(reader, &x)
			sum += x
			if x > maxVal {
				maxVal = x
				maxIdx = i
			}
			// without removal
			if sum <= s {
				if i+1 > bestLen {
					bestLen = i + 1
					bestRemove = 0
				}
			}
			// with removal of max element
			if sum-maxVal <= s {
				if i > bestLen {
					bestLen = i
					// convert to 1-based index
					bestRemove = maxIdx + 1
				}
			}
		}
		fmt.Fprintln(writer, bestRemove)
	}
}
