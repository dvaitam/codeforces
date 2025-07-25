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
		var sum, maxVal int64
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(in, &x)
			sum += x
			if x > maxVal {
				maxVal = x
			}
		}
		if sum == 0 {
			fmt.Fprintln(out, 0)
			continue
		}
		diff := maxVal - (sum - maxVal)
		if diff <= 0 {
			fmt.Fprintln(out, 1)
		} else {
			fmt.Fprintln(out, diff)
		}
	}
}
