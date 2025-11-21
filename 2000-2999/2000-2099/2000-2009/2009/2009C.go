package main

import (
	"bufio"
	"fmt"
	"os"
)

func ceilDiv(a, b int64) int64 {
	if a == 0 {
		return 0
	}
	return (a + b - 1) / b
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var x, y, k int64
		fmt.Fscan(in, &x, &y, &k)

		nx := ceilDiv(x, k)
		ny := ceilDiv(y, k)

		even := 2 * max64(nx, ny)

		oddBase := max64(ny, nx-1)
		if oddBase < 0 {
			oddBase = 0
		}
		odd := 2*oddBase + 1

		ans := even
		if odd < ans {
			ans = odd
		}

		fmt.Fprintln(out, ans)
	}
}
