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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		minVal := int(1<<31 - 1)
		maxVal := -1
		minIdx := 1
		maxIdx := 1
		for i := 1; i <= n; i++ {
			var v int
			fmt.Fscan(in, &v)
			if v < minVal {
				minVal = v
				minIdx = i
			}
			if v > maxVal {
				maxVal = v
				maxIdx = i
			}
		}
		fmt.Fprintf(out, "%d %d\n", minIdx, maxIdx)
	}
}
