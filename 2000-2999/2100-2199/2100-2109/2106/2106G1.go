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
		var n, root int
		fmt.Fscan(in, &n, &root)
		_ = root // root is irrelevant for offline reconstruction

		val := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &val[i])
		}
		// consume edges
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			_, _ = u, v
		}

		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, val[i])
		}
		if t > 1 {
			fmt.Fprintln(out)
		}
	}
}
