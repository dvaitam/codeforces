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
		var n, k int
		var s string
		fmt.Fscan(in, &n, &k)
		fmt.Fscan(in, &s)

		pairs := n / 2
		balance := 0
		segments := 0
		for _, ch := range s {
			if ch == '(' {
				balance++
			} else {
				balance--
				if balance == 0 {
					segments++
				}
			}
		}
		extra := pairs - segments
		if extra < 0 {
			extra = 0
		}
		if k > extra {
			k = extra
		}
		fmt.Fprintln(out, segments+k)
	}
}
