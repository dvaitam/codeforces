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
		var a, b, c int64
		fmt.Fscan(in, &a, &b, &c)
		total := a + 2*b + 3*c
		if total%2 == 0 {
			fmt.Fprintln(out, 0)
		} else {
			fmt.Fprintln(out, 1)
		}
	}
}
