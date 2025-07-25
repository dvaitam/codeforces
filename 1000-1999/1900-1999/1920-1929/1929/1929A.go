package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt.
// The beauty of an array after rearranging is maximized when
// the smallest element is first and the largest element is last.
// Therefore the result is simply max(a) - min(a).
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
		fmt.Fscan(reader, &n)
		var minVal, maxVal int64
		minVal = 1<<63 - 1
		maxVal = -1 << 63
		for i := 0; i < n; i++ {
			var v int64
			fmt.Fscan(reader, &v)
			if v < minVal {
				minVal = v
			}
			if v > maxVal {
				maxVal = v
			}
		}
		fmt.Fprintln(writer, maxVal-minVal)
	}
}
