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
		var n, k, x int64
		fmt.Fscan(in, &n, &k, &x)

		if n < k || x < k-1 {
			fmt.Fprintln(out, -1)
			continue
		}

		base := k * (k - 1) / 2
		val := x
		if val == k {
			val--
		}
		sum := base + (n-k)*val
		fmt.Fprintln(out, sum)
	}
}
