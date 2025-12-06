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
		var s string
		fmt.Fscan(in, &n, &s)

		m := n / 2
		found := false
		for i := 0; i < n; i++ {
			if s[i] == '0' {
				pos := i + 1
				if pos <= m {
					fmt.Fprintf(out, "%d %d %d %d\n", pos, n, pos+1, n)
				} else {
					fmt.Fprintf(out, "%d %d %d %d\n", 1, pos, 1, pos-1)
				}
				found = true
				break
			}
		}
		if !found {
			fmt.Fprintf(out, "%d %d %d %d\n", 1, n-1, 2, n)
		}
	}
}