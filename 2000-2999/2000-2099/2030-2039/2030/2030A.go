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
		minVal := int(1e9)
		maxVal := -1
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			if x < minVal {
				minVal = x
			}
			if x > maxVal {
				maxVal = x
			}
		}
		if n == 1 {
			fmt.Fprintln(out, 0)
			continue
		}
		score := (n - 1) * (maxVal - minVal)
		fmt.Fprintln(out, score)
	}
}
