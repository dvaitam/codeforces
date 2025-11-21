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
		var n, x int64
		fmt.Fscan(in, &n, &x)

		var sum int64
		var maxA int64
		for i := int64(0); i < n; i++ {
			var val int64
			fmt.Fscan(in, &val)
			sum += val
			if val > maxA {
				maxA = val
			}
		}

		customers := (sum + x - 1) / x
		if maxA > customers {
			customers = maxA
		}
		fmt.Fprintln(out, customers)
	}
}
