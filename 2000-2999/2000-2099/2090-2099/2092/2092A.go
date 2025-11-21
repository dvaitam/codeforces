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
		var minVal, maxVal int64
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(in, &x)
			if i == 0 || x < minVal {
				minVal = x
			}
			if i == 0 || x > maxVal {
				maxVal = x
			}
		}
		fmt.Fprintln(out, maxVal-minVal)
	}
}
