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
		c0, c1 := 0, 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			if x == 0 {
				c0++
			} else if x == 1 {
				c1++
			}
		}
		ans := int64(c1) * (int64(1) << c0)
		fmt.Fprintln(out, ans)
	}
}
