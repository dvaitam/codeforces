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
		var s int64
		fmt.Fscan(in, &n, &s)
		ans := 0
		for i := 0; i < n; i++ {
			var dx, dy int64
			var x, y int64
			fmt.Fscan(in, &dx, &dy, &x, &y)
			if dx*(s-2*x) == dy*(s-2*y) {
				ans++
			}
		}
		fmt.Fprintln(out, ans)
	}
}
