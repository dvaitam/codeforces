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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		var x, y int64
		fmt.Fscan(in, &n, &m, &x, &y)
		for i := 0; i < n; i++ {
			var tmp int64
			fmt.Fscan(in, &tmp)
		}
		for i := 0; i < m; i++ {
			var tmp int64
			fmt.Fscan(in, &tmp)
		}
		fmt.Fprintln(out, n+m)
	}
}
