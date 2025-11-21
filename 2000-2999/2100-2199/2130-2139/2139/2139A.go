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
		var a, b int
		fmt.Fscan(in, &a, &b)
		if a == b {
			fmt.Fprintln(out, 0)
			continue
		}
		if a < b {
			a, b = b, a
		}
		if a%b == 0 {
			fmt.Fprintln(out, 1)
		} else {
			fmt.Fprintln(out, 2)
		}
	}
}
