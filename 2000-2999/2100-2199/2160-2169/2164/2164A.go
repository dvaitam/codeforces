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
		var n int
		fmt.Fscan(in, &n)

		var minVal, maxVal int64
		for i := 0; i < n; i++ {
			var v int64
			fmt.Fscan(in, &v)
			if i == 0 {
				minVal, maxVal = v, v
			} else {
				if v < minVal {
					minVal = v
				}
				if v > maxVal {
					maxVal = v
				}
			}
		}

		var x int64
		fmt.Fscan(in, &x)

		if x >= minVal && x <= maxVal {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
