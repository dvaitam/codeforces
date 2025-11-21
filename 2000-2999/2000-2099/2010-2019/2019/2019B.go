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
		var n, q int
		fmt.Fscan(in, &n, &q)

		x := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &x[i])
		}

		counts := make(map[int64]int64, 2*n)
		for idx := 0; idx < n; idx++ {
			left := int64(idx)                 // i-1
			right := int64(n - idx)            // n - i + 1
			tail := int64(n - idx - 1)         // n - i
			val := left*right + tail           // coverage at x_i
			counts[val]++
		}

		for idx := 0; idx+1 < n; idx++ {
			gap := x[idx+1] - x[idx] - 1
			if gap <= 0 {
				continue
			}
			i := int64(idx + 1)
			val := i * (int64(n) - i)
			counts[val] += gap
		}

		for i := 0; i < q; i++ {
			var k int64
			fmt.Fscan(in, &k)
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, counts[k])
		}
		fmt.Fprintln(out)
	}
}
