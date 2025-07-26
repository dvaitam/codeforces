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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		max1, max2 := 0, 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			if x > max1 {
				max2 = max1
				max1 = x
			} else if x > max2 {
				max2 = x
			}
		}
		fmt.Fprintln(out, max1+max2)
	}
}
