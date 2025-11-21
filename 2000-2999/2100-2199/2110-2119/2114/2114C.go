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
		last := int(-2e9)
		ans := 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			if x >= last+2 {
				ans++
				last = x
			}
		}
		fmt.Fprintln(out, ans)
	}
}
