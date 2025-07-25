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
		sum := int64(0)
		odd := 0
		for i := 1; i <= n; i++ {
			var x int64
			fmt.Fscan(in, &x)
			sum += x
			if x%2 == 1 {
				odd++
			}
			if i == 1 {
				fmt.Fprint(out, x)
			} else {
				loss := int64(odd / 3)
				if odd%3 == 1 {
					loss++
				}
				fmt.Fprint(out, sum-loss)
			}
			if i == n {
				fmt.Fprintln(out)
			} else {
				fmt.Fprint(out, " ")
			}
		}
	}
}
