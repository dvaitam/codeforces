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
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		c, d, e := int64(0), int64(0), int64(0)
		for i := 0; i < n; i++ {
			var b int64
			fmt.Fscan(in, &b)
			if a[i] > b {
				c += a[i]
			} else if a[i] < b {
				d += b
			} else if a[i] == 1 {
				e += a[i]
			} else if a[i] == -1 {
				e++
				c--
				d--
			}
		}
		x1 := c + e
		x2 := d + e
		x3 := (c + e + d) >> 1
		ans := x1
		if x2 < ans {
			ans = x2
		}
		if x3 < ans {
			ans = x3
		}
		fmt.Fprintln(out, ans)
	}
}
