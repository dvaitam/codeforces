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
		var a, b, c int
		fmt.Fscan(in, &a, &b, &c)
		if a == b {
			fmt.Fprintln(out, c)
		} else if a == c {
			fmt.Fprintln(out, b)
		} else {
			fmt.Fprintln(out, a)
		}
	}
}
