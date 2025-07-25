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
		var a, b int64
		fmt.Fscan(in, &a, &b)
		var x int64
		if b%a == 0 {
			x = b * (b / a)
		} else {
			x = a * b
		}
		fmt.Fprintln(out, x)
	}
}
