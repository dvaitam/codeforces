package main

import (
	"bufio"
	"fmt"
	"os"
)

// The optimal strategy is to wait until the timer drops to 1 before using
// each tool. Using a tool when the timer is 1 increases the remaining time
// by min(x_i, a-1) seconds. Therefore the bomb will explode after b plus the
// sum of min(x_i, a-1) seconds.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var a, b int64
		var n int
		fmt.Fscan(in, &a, &b, &n)
		ans := b
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(in, &x)
			// each tool effectively extends time by min(x, a-1)
			if x > a-1 {
				ans += a - 1
			} else {
				ans += x
			}
		}
		fmt.Fprintln(out, ans)
	}
}
