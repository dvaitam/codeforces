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
		var n int64
		fmt.Fscan(in, &n)
		switch {
		case n < 1000:
			fmt.Fprintln(out, n)
		case n < 1_000_000:
			k := (n + 500) / 1000
			if k == 1000 {
				fmt.Fprintln(out, "1M")
			} else {
				fmt.Fprintf(out, "%dK\n", k)
			}
		default:
			m := (n + 500_000) / 1_000_000
			fmt.Fprintf(out, "%dM\n", m)
		}
	}
}
