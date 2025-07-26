package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)

		maxVal := int64(-1 << 63)
		rMax, cMax := 0, 0
		for i := 1; i <= n; i++ {
			for j := 1; j <= m; j++ {
				var x int64
				fmt.Fscan(reader, &x)
				if x > maxVal {
					maxVal = x
					rMax = i
					cMax = j
				}
			}
		}
		h := max(rMax, n-rMax+1)
		w := max(cMax, m-cMax+1)
		fmt.Fprintln(writer, h*w)
	}
}
