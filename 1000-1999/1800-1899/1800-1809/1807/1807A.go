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
		var a, b, c int
		fmt.Fscan(in, &a, &b, &c)
		if a+b == c {
			fmt.Fprintln(out, "+")
		} else {
			fmt.Fprintln(out, "-")
		}
	}
}
