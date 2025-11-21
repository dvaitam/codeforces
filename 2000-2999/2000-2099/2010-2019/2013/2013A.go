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
		var n, x, y int64
		fmt.Fscan(in, &n)
		fmt.Fscan(in, &x, &y)
		if n == 0 {
			fmt.Fprintln(out, 0)
			continue
		}
		limit := x
		if y < limit {
			limit = y
		}
		ans := (n + limit - 1) / limit
		fmt.Fprintln(out, ans)
	}
}
