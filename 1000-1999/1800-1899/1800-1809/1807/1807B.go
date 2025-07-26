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
		evenSum := 0
		oddSum := 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			if x%2 == 0 {
				evenSum += x
			} else {
				oddSum += x
			}
		}
		if evenSum > oddSum {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
