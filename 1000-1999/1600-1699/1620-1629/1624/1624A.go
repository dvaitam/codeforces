package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt for contest 1624A.
// We can increment any subset of indices by 1 in one operation. The optimal
// strategy is to raise all elements to the current maximum value, since we can
// only increase elements. The minimum number of operations is therefore equal to
// max(a) - min(a).
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
		if n <= 0 {
			fmt.Fprintln(writer, 0)
			continue
		}
		minVal := int(1<<63 - 1)
		maxVal := -minVal - 1
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			if x < minVal {
				minVal = x
			}
			if x > maxVal {
				maxVal = x
			}
		}
		fmt.Fprintln(writer, maxVal-minVal)
	}
}
