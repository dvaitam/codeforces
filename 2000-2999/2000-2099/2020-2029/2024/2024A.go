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
		var a, b int64
		fmt.Fscan(in, &a, &b)
		if a >= b {
			fmt.Fprintln(out, a)
			continue
		}
		ans := int64(0)
		x1 := b - a
		if x1 >= 0 && x1 <= a && x1*2 < b {
			candidate := a - x1
			if candidate > ans {
				ans = candidate
			}
		}
		x2 := (b + 1) / 2
		if x2 <= a {
			candidate := a - x2
			if candidate > ans {
				ans = candidate
			}
		}
		fmt.Fprintln(out, ans)
	}
}
