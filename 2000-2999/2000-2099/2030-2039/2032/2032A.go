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
		total := 0
		for i := 0; i < 2*n; i++ {
			var x int
			fmt.Fscan(in, &x)
			total += x
		}
		minOn := total % 2
		maxOn := total
		if maxOn > 2*n-total {
			maxOn = 2*n - total
		}
		fmt.Fprintf(out, "%d %d\n", minOn, maxOn)
	}
}
