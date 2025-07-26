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
		var n, r, b int
		fmt.Fscan(in, &n, &r, &b)
		groups := b + 1
		q := r / groups
		extra := r % groups
		for i := 0; i < groups; i++ {
			for j := 0; j < q; j++ {
				fmt.Fprint(out, "R")
			}
			if extra > 0 {
				fmt.Fprint(out, "R")
				extra--
			}
			if i < b {
				fmt.Fprint(out, "B")
			}
		}
		fmt.Fprintln(out)
	}
}
