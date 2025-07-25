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
		var n, k int
		fmt.Fscan(in, &n, &k)
		if k == n {
			for i := 0; i < n; i++ {
				if i > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, 1)
			}
			fmt.Fprintln(out)
			continue
		}
		if k == 1 {
			for i := 1; i <= n; i++ {
				if i > 1 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, i)
			}
			fmt.Fprintln(out)
			continue
		}
		fmt.Fprintln(out, -1)
	}
}
